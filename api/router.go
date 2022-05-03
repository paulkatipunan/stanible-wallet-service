package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"api.stanible.com/wallet/database"
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

type Accounts struct {
	Pk_account_id string `json:"pk_account_id"`
	User_id       string `json:"user_id"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	Active        bool   `json:"active"`
	Created_at    string `json:"created_at"`
	Updated_at    string `json:"updated_at"`
}

func walletBalance(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	// close the db connection
	defer db.Close()

	var account_list []Accounts

	sqlStatement := `SELECT * FROM accounts`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// close the db connection
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var accounts Accounts

		// unmarshal the row object to accounts
		err = rows.Scan(
			&accounts.Pk_account_id,
			&accounts.User_id,
			&accounts.Type,
			&accounts.Description,
			&accounts.Active,
			&accounts.Created_at,
			&accounts.Updated_at,
		)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		account_list = append(account_list, accounts)

		fmt.Println(accounts)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account_list)
}

func Router() *mux.Router {
	routers := mux.NewRouter().StrictSlash(true)
	routers.Use(commonMiddleware)

	routers.HandleFunc("/wallet/register", register).Methods("POST")
	routers.HandleFunc("/wallet/fiat", fiatTransaction).Methods("POST")
	routers.HandleFunc("/wallet/fiat/balance/{user_id}", walletBalance).Methods("GET")

	return routers
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
