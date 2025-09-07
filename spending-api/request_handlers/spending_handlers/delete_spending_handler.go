package spending_handlers

import (
	"net/http"
	"spending/repositories"
	"spending/repositories/spending_repo"
	"spending/request_handlers"
	"spending/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type deleteSpendingHandler struct {
	spending_repo spending_repo.SpendingRepository
	unit_of_work  repositories.UnitOfWork
}

func NewDeleteSpendingHandler(spendingRepo spending_repo.SpendingRepository, unitOfWork repositories.UnitOfWork) request_handlers.RequestHandler {
	return &deleteSpendingHandler{
		spending_repo: spendingRepo,
		unit_of_work:  unitOfWork,
	}
}

func (handler *deleteSpendingHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "DeleteSpendingHandler")
	defer span.End()

	routerVars := mux.Vars(request)
	spendingUUId, err := uuid.Parse(routerVars["id"])
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// Start a new transaction
	tx, err := handler.unit_of_work.BeginTx()
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer handler.unit_of_work.CommitOrRollback(tx, err)

	err = handler.spending_repo.DeleteSpending(context, tx, spendingUUId)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
