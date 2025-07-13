package spending_handlers

import (
	"github.com/gorilla/mux"
)

func RegisterSpendingHandlers(router *mux.Router) {
	router.HandleFunc("/spending", GetSpendingListHandler).Methods("GET")
	router.HandleFunc("/spending/{id}", GetSpendingRequestHandler).Methods("GET")
	router.HandleFunc("/spending", CreateSpendingRequestHandler).Methods("POST")
}
