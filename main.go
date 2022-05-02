package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Balance struct {
	Id     int
	UserId int
	Bal    float32
}

var balance = []Balance{
	{1, 11, 10000},
	{2, 22, 20000},
	{3, 33, 30000},
}

func register(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

func fiatTransaction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

func walletBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
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
