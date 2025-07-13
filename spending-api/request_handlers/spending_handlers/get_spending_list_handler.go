package spending_handlers

import (
	"encoding/json"
	"net/http"
	"spending/mappers"
	"spending/repositories/spending_repo"
	"spending/utils"

	"github.com/rs/zerolog/log"
)

func GetSpendingListHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Fetching spending list...")
	records := spending_repo.GetSpendingList()
	response := mappers.MapSpendingList(records)

	err := json.NewEncoder(w).Encode(response)
	utils.CheckError(err)
}
