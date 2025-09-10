package category_handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"spending/models"
	"spending/repositories"
	"spending/repositories/category_repo"
	"spending/repositories/store_repo"
	"spending/request_handlers"
	"spending/request_handlers/store_handlers"
	"spending/utils"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type updateCategoryHandler struct {
	category_repo category_repo.CategoryRepository
	store_repo    store_repo.StoreRepository
	unit_of_work  repositories.UnitOfWork
}

func NewUpdateCategoryHandler(categoryRepo category_repo.CategoryRepository, storeRepo store_repo.StoreRepository, unitOfWork repositories.UnitOfWork) request_handlers.RequestHandler {
	return &updateCategoryHandler{
		category_repo: categoryRepo,
		store_repo:    storeRepo,
		unit_of_work:  unitOfWork,
	}
}

type UpdateCategoryRequest struct {
	Id            uuid.UUID                            `json:"id"`
	Name          string                               `json:"name"`
	AddedStores   []*store_handlers.CreateStoreRequest `json:"addedStores"`
	EditedStores  []*store_handlers.UpdateStoreRequest `json:"editedStores"`
	DeletedStores []uuid.UUID                          `json:"deletedStores"`
}

func (request UpdateCategoryRequest) Valid(context context.Context) error {
	if request.Id == uuid.Nil {
		return fmt.Errorf("id cannot be empty")
	}

	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if request.AddedStores != nil {
		for _, store := range request.AddedStores {
			if err := store.Valid(context); err != nil {
				return fmt.Errorf("invalid added store: %w", err)
			}
		}
	}

	if request.EditedStores != nil {
		for _, store := range request.EditedStores {
			if err := store.Valid(); err != nil {
				return fmt.Errorf("invalid edited store: %w", err)
			}
		}
	}

	if request.DeletedStores != nil {
		for _, storeID := range request.DeletedStores {
			if storeID == uuid.Nil {
				return fmt.Errorf("invalid deleted store ID: %s", storeID)
			}
		}
	}

	return nil
}

func (handler *updateCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "UpdateCategoryHandler")
	defer span.End()

	routerVars := mux.Vars(request)
	categoryUUId, err := uuid.Parse(routerVars["id"])
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	command, err := utils.DecodeValid[UpdateCategoryRequest](ctx, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if command.Id != categoryUUId {
		err = fmt.Errorf("id in path and body do not match")
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		existingCategory, txErr := handler.category_repo.GetCategoryByName(ctx, tx, command.Name)
		if txErr != nil {
			return txErr
		}

		if existingCategory != nil && existingCategory.UUId != categoryUUId {
			txErr = utils.ErrConflict
			return txErr
		}

		category, txErr := handler.category_repo.GetCategoryByUUId(ctx, tx, categoryUUId)
		if txErr != nil {
			return txErr
		}

		if category == nil {
			txErr = utils.ErrNotFound
			return txErr
		}

		category.Name = command.Name
		category.UpdatedAt = time.Now().UTC()
		txErr = handler.category_repo.UpdateCategory(ctx, tx, category)
		if txErr != nil {
			return txErr
		}

		if len(command.AddedStores) > 0 {
			txErr = handler.addStores(ctx, tx, category, command.AddedStores)
			if txErr != nil {
				return txErr
			}
		}

		if len(command.EditedStores) > 0 {
			txErr = handler.updateStores(ctx, tx, command.EditedStores)
			if txErr != nil {
				return txErr
			}
		}

		if len(command.DeletedStores) > 0 {
			txErr = handler.deleteStores(ctx, tx, command.DeletedStores)
			if txErr != nil {
				return txErr
			}
		}

		return nil
	})

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), utils.MapErrorToStatusCode(err))
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *updateCategoryHandler) addStores(ctx context.Context, tx *sql.Tx, category *models.Category, stores []*store_handlers.CreateStoreRequest) error {
	var storeModels []*models.Store
	for _, storeReq := range stores {
		store := models.NewStore(storeReq.Name, category.Id)
		storeModels = append(storeModels, store)
	}

	_, txErr := handler.store_repo.InsertStores(ctx, tx, storeModels)
	if txErr != nil {
		return txErr
	}

	return nil
}

func (handler *updateCategoryHandler) updateStores(ctx context.Context, tx *sql.Tx, stores []*store_handlers.UpdateStoreRequest) error {
	for _, storeReq := range stores {
		store, txErr := handler.store_repo.GetStoreByUUId(ctx, tx, storeReq.Id)
		if txErr != nil {
			return txErr
		}

		store.Name = storeReq.Name
		store.UpdatedAt = time.Now().UTC()
		txErr = handler.store_repo.UpdateStore(ctx, tx, store)
		if txErr != nil {
			return txErr
		}
	}

	return nil
}

func (handler *updateCategoryHandler) deleteStores(ctx context.Context, tx *sql.Tx, storeIDs []uuid.UUID) error {
	if len(storeIDs) > 0 {
		return handler.store_repo.DeleteStores(ctx, tx, storeIDs)
	}

	return nil
}
