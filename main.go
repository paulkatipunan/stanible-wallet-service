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

func returnBalance(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/wallet", returnBalance).Methods("GET")
	http.ListenAndServe(":8081", router)
}
