package spending_handlers

import (
	"encoding/json"
	"net/http"
	"spending/mappers"
	"spending/repositories/spending_repo"
	"spending/utils"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type GetSpendingListHandler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
}
type getSpendingListHandler struct {
	spending_repo spending_repo.SpendingRepository
}

func NewGetSpendingListHandler(spendingRepo spending_repo.SpendingRepository) GetSpendingListHandler {
	return &getSpendingListHandler{
		spending_repo: spendingRepo,
	}
}

func (handler *getSpendingListHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "GetSpendingListHandler")
	defer span.End()

	log.Info().Msg("Fetching spending list...")
	records, err := handler.spending_repo.GetSpendingList(context)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handler.spending_repo.LoadSpendingListCategory(context, records)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mappers.MapSpendingList(records)

	err = json.NewEncoder(writer).Encode(response)
	utils.TraceError(span, err)
}
