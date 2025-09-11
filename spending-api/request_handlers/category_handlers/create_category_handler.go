package category_handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
	"spending/repositories"
	"spending/repositories/category_repo"
	"spending/repositories/store_repo"
	"spending/request_handlers"
	"spending/request_handlers/store_handlers"
	"spending/utils"

	"go.opentelemetry.io/otel"
)

type createCategoryHandler struct {
	category_repo category_repo.CategoryRepository
	store_repo    store_repo.StoreRepository
	unit_of_work  repositories.UnitOfWork
}

func NewCreateCategoryHandler(categoryRepo category_repo.CategoryRepository, storeRepo store_repo.StoreRepository, unitOfWork repositories.UnitOfWork) request_handlers.RequestHandler {
	return &createCategoryHandler{
		category_repo: categoryRepo,
		store_repo:    storeRepo,
		unit_of_work:  unitOfWork,
	}
}

type CreateCategoryRequest struct {
	Name   string                              `json:"name"`
	Stores []store_handlers.CreateStoreRequest `json:"stores"`
}

func (request CreateCategoryRequest) Valid(context context.Context) error {
	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	for _, store := range request.Stores {
		if err := store.Valid(context); err != nil {
			return fmt.Errorf("invalid store: %w", err)
		}
	}
	return nil
}

func (handler *createCategoryHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	ctx, span := tracer.Start(request.Context(), "CreateCategoryHandler")
	defer span.End()

	command, err := utils.DecodeValid[CreateCategoryRequest](ctx, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var category *models.Category

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		existingCategory, txErr := handler.category_repo.GetCategoryByName(ctx, tx, command.Name)
		if txErr != nil {
			return txErr
		}

		if existingCategory != nil {
			txErr = utils.ErrResourceExists
			return txErr
		}

		newCategory := models.NewCategory(command.Name)
		category, txErr = handler.category_repo.InsertCategory(ctx, tx, newCategory)
		if txErr != nil {
			return txErr
		}

		if len(command.Stores) > 0 {
			var stores []*models.Store
			for _, storeReq := range command.Stores {
				store := models.NewStore(storeReq.Name, category.Id)
				stores = append(stores, store)
			}

			createdStores, txErr := handler.store_repo.InsertStores(ctx, tx, stores)
			if txErr != nil {
				return txErr
			}
			category.Stores = createdStores
		}

		return nil
	})

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), utils.MapErrorToStatusCode(err))
		return
	}

	response := mappers.MapCategory(category)
	writer.Header().Set("Location", fmt.Sprintf("/categories/%s", category.UUId))
	err = utils.Encode(ctx, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
