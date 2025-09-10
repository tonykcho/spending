package store_handlers

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
	"spending/utils"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type createStoreHandler struct {
	store_repo    store_repo.StoreRepository
	category_repo category_repo.CategoryRepository
	unit_of_work  repositories.UnitOfWork
}

func NewCreateStoreHandler(storeRepo store_repo.StoreRepository, categoryRepo category_repo.CategoryRepository, unitOfWork repositories.UnitOfWork) *createStoreHandler {
	return &createStoreHandler{
		store_repo:    storeRepo,
		category_repo: categoryRepo,
		unit_of_work:  unitOfWork,
	}
}

type CreateStoreRequest struct {
	Name       string    `json:"name"`
	CategoryId uuid.UUID `json:"categoryId"`
}

func (request CreateStoreRequest) Valid(ctx context.Context) error {
	if request.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if request.CategoryId == uuid.Nil {
		return fmt.Errorf("categoryId must be a valid UUID")
	}
	return nil
}

func (handler *createStoreHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "CreateStoreHandler")
	defer span.End()

	command, err := utils.DecodeValid[CreateStoreRequest](context, request)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var store *models.Store

	err = handler.unit_of_work.WithTransaction(func(tx *sql.Tx) error {
		category, txErr := handler.category_repo.GetCategoryByUUId(context, tx, command.CategoryId)
		if txErr != nil {
			return txErr
		}

		if category == nil {
			return utils.ErrNotFound
		}

		existingStore, txErr := handler.store_repo.GetStoreByCategoryAndName(context, tx, category.Id, command.Name)
		if txErr != nil {
			return txErr
		}

		if existingStore != nil {
			return utils.ErrResourceExists
		}

		newStore := &models.Store{
			Name:       command.Name,
			CategoryId: category.Id,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		store, txErr = handler.store_repo.InsertStore(context, tx, newStore)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), utils.MapErrorToStatusCode(err))
		return
	}

	writer.WriteHeader(http.StatusCreated)
	response := mappers.MapStore(store)
	err = utils.Encode(context, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
