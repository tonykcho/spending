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

func GetSpendingListHandler(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("spending-api")
	context, span := tracer.Start(r.Context(), "GetSpendingListHandler")
	defer span.End()

	log.Info().Msg("Fetching spending list...")
	records, err := spending_repo.GetSpendingList(context)

	if err != nil {
		utils.TraceError(span, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = spending_repo.LoadSpendingListCategory(context, records)
	if err != nil {
		utils.TraceError(span, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mappers.MapSpendingList(records)

	err = json.NewEncoder(w).Encode(response)
	utils.TraceError(span, err)
}
