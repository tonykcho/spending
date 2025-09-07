package spending_handlers

import (
	"net/http"
	"spending/mappers"
	"spending/repositories/spending_repo"
	"spending/request_handlers"
	"spending/utils"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type getSpendingListHandler struct {
	spending_repo spending_repo.SpendingRepository
}

func NewGetSpendingListHandler(spendingRepo spending_repo.SpendingRepository) request_handlers.RequestHandler {
	return &getSpendingListHandler{
		spending_repo: spendingRepo,
	}
}

func (handler *getSpendingListHandler) Handle(writer http.ResponseWriter, request *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(request.Context(), "GetSpendingListHandler")
	defer span.End()

	log.Info().Msg("Fetching spending list...")
	records, err := handler.spending_repo.GetSpendingList(context, nil)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handler.spending_repo.LoadSpendingListCategory(context, nil, records)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mappers.MapSpendingList(records)

	err = utils.Encode(context, writer, http.StatusOK, response)
	utils.TraceError(span, err)
}
