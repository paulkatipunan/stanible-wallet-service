package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"api.stanible.com/wallet/database"
	"api.stanible.com/wallet/models"
	"api.stanible.com/wallet/utils"
	"github.com/gorilla/mux"
)

func register(w http.ResponseWriter, r *http.Request) {
	var account models.Accounts

	payloadErr := json.NewDecoder(r.Body).Decode(&account)

	if payloadErr != nil {
		log.Fatalf("Unable to decode the request body.  %v", payloadErr)
	}

	db := database.CreateConnection()
	defer db.Close()

	var pk_account_id string
	sqlStatement := `INSERT INTO accounts (user_id, type, description) VALUES ($1, $2, $3) RETURNING pk_account_id`
	queryErr := db.QueryRow(sqlStatement, account.User_id, account.Type, account.Description).Scan(&pk_account_id)

	if queryErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", queryErr.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", ""))
	}
}

func fiatTransaction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response("success", ""))
}

func walletBalance(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	defer db.Close()

	var account_list []models.Accounts

	sqlStatement := `SELECT * FROM accounts`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var accounts models.Accounts

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
