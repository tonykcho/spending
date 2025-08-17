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
	_, span := tracer.Start(r.Context(), "GetSpendingListHandler")
	defer span.End()

	log.Info().Msg("Fetching spending list...")
	records := spending_repo.GetSpendingList()
	response := mappers.MapSpendingList(records)

	err := json.NewEncoder(w).Encode(response)
	utils.CheckError(err)
}
