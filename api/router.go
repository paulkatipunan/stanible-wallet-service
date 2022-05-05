package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

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
		json.NewEncoder(w).Encode(utils.Response("error", queryErr.Error(), nil))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", "", nil))
	}
}

func fiatTransaction(w http.ResponseWriter, r *http.Request) {
	// Get payloads and assign fiat_transaction model
	var errorFound bool = true
	var transactionPayload models.Transaction_payload
	json.NewDecoder(r.Body).Decode(&transactionPayload)

	// Begin tx
	ctx := context.Background()
	db := database.CreateConnection()
	tx, err := db.BeginTx(ctx, nil)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	sqlGetUsers := `
		SELECT
			COUNT(*) as total_users
		FROM
			accounts
		WHERE
			user_id IN (
				$1,
				$2
			) AND
			active = true
	`

	// Check if sender and receiver addresses are both active
	var total_users int
	tx.QueryRow(sqlGetUsers, transactionPayload.Sender_user_id, transactionPayload.Receiver_user_id).Scan(&total_users)

	if total_users < 2 {
		// If not, rollback and return error
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Sender or receiver addresses are not active", nil))
		return
	} else {
		errorFound = false
	}

	// Insert fiat_transaction record
	sqlInsertFiatTransaction := `
		INSERT INTO fiat_transactions (
			pk_fiat_transaction_id, fk_user_id, fk_transaction_type_id, fk_fiat_currency_id, amount
		) VALUES
			($1, $2, $3, $4, $5),
			($6, $7, $8, $9, $10)
	`
	pkSenderId := uuid.New()
	pkReceiverId := uuid.New()
	rows, err := tx.Query(
		sqlInsertFiatTransaction,

		pkSenderId,
		transactionPayload.Sender_user_id,
		transactionPayload.Transaction_type_id,
		transactionPayload.Fiat_currency_id,
		-transactionPayload.Amount,

		pkReceiverId,
		transactionPayload.Receiver_user_id,
		transactionPayload.Transaction_type_id,
		transactionPayload.Fiat_currency_id,
		transactionPayload.Amount,
	)
	if err != nil {
		// Rollback if error
		errorFound = true
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	} else {
		errorFound = false
		rows.Close()
	}

	// Insert fiat_transations_assoc record
	sqlInsertFiatTransactionAssoc := `
		INSERT INTO fiat_transactions_assoc (
			pk_sender_fiat_transaction_id, pk_receiver_fiat_transaction_id, ramp_tx_id
		) VALUES
			($1, $2, $3)
	`
	row, err := tx.Query(
		sqlInsertFiatTransactionAssoc,
		pkSenderId,
		pkReceiverId,
		transactionPayload.Ramp_tx_id,
	)
	if err != nil {
		// Rollback if error
		errorFound = true
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	} else {
		errorFound = false
		row.Close()
	}

	// Commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		errorFound = true
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	} else {
		errorFound = false
	}

	if !errorFound {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", "", nil))
	}
}

func walletBalance(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	defer db.Close()

	vars := mux.Vars(r)
	user_id := vars["user_id"]

	// fmt.Println("user_id", user_id)

	var balance []uint8
	sqlGetUserBalance := `
		SELECT
			coalesce(SUM(ft.amount), 0) + (
			SELECT
				coalesce(SUM(ft.amount), 0) as balance
			FROM
				fiat_transactions ft
			INNER JOIN
				transaction_types tt
				ON
					ft.fk_transaction_type_id = tt.pk_transaction_type_id
			WHERE
				ft.fk_user_id = $1 AND
				tt.type IN ('withdraw', 'buy', 'refund')
			) as balance
		FROM
			fiat_transactions ft
		INNER JOIN
			transaction_types tt
			ON
				ft.fk_transaction_type_id = tt.pk_transaction_type_id
		WHERE
			ft.fk_user_id = $2 AND
			tt.type = 'deposit';
	`
	err := db.QueryRow(sqlGetUserBalance, user_id, user_id).Scan(&balance)
	new_balance := []string{string(balance)}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", "", new_balance))
	}
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
