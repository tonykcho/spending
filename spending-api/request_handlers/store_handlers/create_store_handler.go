package store_handlers

import (
	"context"
	"fmt"
	"net/http"
	"spending/mappers"
	"spending/models"
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
}

func NewCreateStoreHandler(storeRepo store_repo.StoreRepository, categoryRepo category_repo.CategoryRepository) *createStoreHandler {
	return &createStoreHandler{
		store_repo:    storeRepo,
		category_repo: categoryRepo,
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

	category, err := handler.category_repo.GetCategoryByUUId(context, command.CategoryId)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if category == nil {
		utils.TraceError(span, fmt.Errorf("category not found"))
		http.Error(writer, "category not found", http.StatusBadRequest)
		return
	}

	existingStore, err := handler.store_repo.GetStoreByCategoryAndName(context, category.Id, command.Name)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if existingStore != nil {
		utils.TraceError(span, fmt.Errorf("store already exists"))
		http.Error(writer, "store already exists", http.StatusBadRequest)
		return
	}

	newStore := &models.Store{
		Name:       command.Name,
		CategoryId: category.Id,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	id, err := handler.store_repo.InsertStore(context, newStore)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	createdStore, err := handler.store_repo.GetStoreById(context, id)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	response := mappers.MapStore(createdStore)

	err = utils.Encode(context, writer, http.StatusCreated, response)
	utils.TraceError(span, err)
}
