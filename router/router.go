package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type StdResponse struct {
	Status  string  `json:"status"`
	Message *string `json:"message"`
	Data    *string `json:"data"`
}

var response = StdResponse{
	Status:  "success",
	Message: nil,
	Data:    nil,
}

func register(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func fiatTransaction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func walletBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	router.HandleFunc("/wallet/register", register).Methods("POST")
	router.HandleFunc("/wallet/fiat", fiatTransaction).Methods("POST")
	router.HandleFunc("/wallet/fiat/balance/{user_id}", walletBalance).Methods("GET")

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
