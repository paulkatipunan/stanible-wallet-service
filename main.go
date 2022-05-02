package main

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

func main() {
	port := ":8080"
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	router.HandleFunc("/wallet/register", register).Methods("POST")
	router.HandleFunc("/wallet/fiat", fiatTransaction).Methods("POST")
	router.HandleFunc("/wallet/fiat/balance/{user_id}", walletBalance).Methods("GET")
	http.ListenAndServe(port, router)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
