package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"api.stanible.com/wallet/database"
	"api.stanible.com/wallet/enums"
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

func fiatDeposit(w http.ResponseWriter, r *http.Request) {
	// Get payloads and assign fiat_transaction model
	var transactionPayload models.Transaction_payload
	json.NewDecoder(r.Body).Decode(&transactionPayload)

	// Validation#01
	// Sender and receiver cannot be the same
	if transactionPayload.Sender_user_id == transactionPayload.Receiver_user_id {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Sender and receiver cannot be the same", nil))
		return
	}

	// Validation#02
	// Check if sender and receiver addresses are both active
	total_users := utils.ActiveSenderReceiver(transactionPayload.Sender_user_id, transactionPayload.Receiver_user_id)
	if total_users < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Sender or receiver addresses are not active", nil))
		return
	}

	// Validation#03
	// Validate user types
	row, err := utils.UserTypes(transactionPayload.Sender_user_id, transactionPayload.Receiver_user_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	}
	sender_type := row[0]
	receiver_type := row[1]
	if enums.PaymentUserTypes[strings.ToUpper(sender_type)] == "" ||
		enums.NonPaymentUserTypes()[strings.ToUpper(receiver_type)] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Invalid user types", nil))
		return
	}

	// Validation#04
	// Check balance cap
	// Balance is from the receiving account
	bal, err := utils.AccountBalance(transactionPayload.Receiver_user_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	}
	if (transactionPayload.Amount + bal) >= enums.BALANCE_CAP {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Balance cap exceeded", nil))
		return
	}

	// Get transaction type and transaction_type_id
	pk_transaction_type_id, _ := utils.GetTransactionType(enums.DEPOSIT)
	transactionPayload.Transaction_type_id = pk_transaction_type_id

	// Insert fiat_transaction and fiat_transaction_assoc records
	txResponse := utils.InsertFiatTransactionRecord(transactionPayload, enums.TX_STATUS["SUCCESS"])

	if txResponse != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", txResponse.Error(), nil))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", "", nil))
	}
}

func fiatBuy(w http.ResponseWriter, r *http.Request) {
	// Get payloads and assign fiat_transaction model
	var transactionPayload models.Transaction_payload
	json.NewDecoder(r.Body).Decode(&transactionPayload)

	// Check balance
	// Balance is from the sending account
	bal, err := utils.AccountBalance(transactionPayload.Sender_user_id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	}

	if bal <= 0 || transactionPayload.Amount > int32(bal) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Insuffucient balance", nil))
		return
	}

	// Get transaction type and transaction_type_id
	pk_transaction_type_id, _ := utils.GetTransactionType(enums.BUY)
	transactionPayload.Transaction_type_id = pk_transaction_type_id

	// Check if sender and receiver addresses are both active
	total_users := utils.ActiveSenderReceiver(transactionPayload.Sender_user_id, transactionPayload.Receiver_user_id)
	if total_users < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Sender or receiver addresses are not active", nil))
		return
	}

	// Insert fiat_transaction and fiat_transaction_assoc records
	txResponse := utils.InsertFiatTransactionRecord(transactionPayload, enums.TX_STATUS["SUCCESS"])

	if txResponse != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", txResponse.Error(), nil))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", "", nil))
	}
}

func fiatRefundRequest(w http.ResponseWriter, r *http.Request) {
	// Get payloads and assign fiat_transaction model
	var transactionPayload models.Transaction_payload
	json.NewDecoder(r.Body).Decode(&transactionPayload)

	// Validation#01
	// Sender and receiver cannot be the same
	if transactionPayload.Sender_user_id == transactionPayload.Receiver_user_id {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Sender and receiver cannot be the same", nil))
		return
	}

	// Validation#02
	// Check if sender and receiver addresses are both active
	total_users := utils.ActiveSenderReceiver(transactionPayload.Sender_user_id, transactionPayload.Receiver_user_id)
	if total_users < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Sender or receiver addresses are not active", nil))
		return
	}

	// Validation#03
	// Validate user types
	row, err := utils.UserTypes(transactionPayload.Sender_user_id, transactionPayload.Receiver_user_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	}
	sender_type := row[0]
	receiver_type := row[1]
	if sender_type != enums.SystemUserTypes["TREASURY"] ||
		enums.CustomerUserTypes[strings.ToUpper(receiver_type)] == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Invalid user types", nil))
		return
	}

	// Validation#04
	// Check balance cap
	// Balance is from the receiving account
	bal, err := utils.AccountBalance(transactionPayload.Receiver_user_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	}

	// NOTE:
	// Balance 0 should be allowed as long as there was a deposit and buy made before,
	// and refund should be less than or equal to the buy amount
	// if bal == 0 {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(utils.Response("error", "Insufficient balance", nil))
	// 	return
	// }

	if (transactionPayload.Amount + bal) >= enums.BALANCE_CAP {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", "Balance cap exceeded", nil))
		return
	}

	// Get transaction type and transaction_type_id
	pk_transaction_type_id, _ := utils.GetTransactionType(enums.REFUND)
	transactionPayload.Transaction_type_id = pk_transaction_type_id

	// Insert fiat_transaction and fiat_transaction_assoc records
	txResponse := utils.InsertFiatTransactionRecord(transactionPayload, enums.TX_STATUS["PENDING"])

	if txResponse != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", txResponse.Error(), nil))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", "", nil))
	}
}

func fiatRefundApprove(w http.ResponseWriter, r *http.Request) {
	// Get payloads and assign fiat_transaction model
	var transactionPayload models.RefundApprove_payload
	json.NewDecoder(r.Body).Decode(&transactionPayload)

	// Validate status payload

	utils.RefundApprove(transactionPayload)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response("success", "", nil))
}

func fiatRefundList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.RefundRequest("success", "", utils.RefundRequestList()))
}

func fiatWithdraw(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(utils.Response("success", "", nil))
}

func fiatWalletBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_id := vars["user_id"]

	bal, err := utils.AccountBalance(user_id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.Response("success", "", []string{strconv.Itoa(int(bal))}))
	}
}

func transactionTypes(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	defer db.Close()

	sqlGetTransactionTypes := `SELECT pk_transaction_type_id, type as type_name FROM transaction_types`
	rows, err := db.Query(sqlGetTransactionTypes)

	var tx_types_list []models.Transaction_types
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	} else {
		defer rows.Close()
		// iterate over the rows
		for rows.Next() {
			var tx_types models.Transaction_types

			// unmarshal the row object to accounts
			err = rows.Scan(
				&tx_types.Pk_transaction_type_id,
				&tx_types.Type_name,
			)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
				return
			}

			tx_types_list = append(tx_types_list, tx_types)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.TxType("success", "", tx_types_list))
	}
}

func fiatCurrencies(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	defer db.Close()

	sqlGetFiatCurrencies := `SELECT pk_fiat_currency_id, name, symbol FROM fiat_currencies`
	rows, err := db.Query(sqlGetFiatCurrencies)

	var fiat_currencies_list []models.Fiat_currencies
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
		return
	} else {
		defer rows.Close()
		// iterate over the rows
		for rows.Next() {
			var fiat_currencies models.Fiat_currencies

			// unmarshal the row object to accounts
			err = rows.Scan(
				&fiat_currencies.Pk_fiat_currency_id,
				&fiat_currencies.Name,
				&fiat_currencies.Symbol,
			)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(utils.Response("error", err.Error(), nil))
				return
			}

			fiat_currencies_list = append(fiat_currencies_list, fiat_currencies)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(utils.FiatCurrencies("success", "", fiat_currencies_list))
	}
}

func Router() *mux.Router {
	routers := mux.NewRouter().StrictSlash(true)
	routers.Use(commonMiddleware)

	routers.HandleFunc("/wallet/register", register).Methods("POST")
	routers.HandleFunc("/wallet/fiat/deposit", fiatDeposit).Methods("POST")
	routers.HandleFunc("/wallet/fiat/buy", fiatBuy).Methods("POST")
	routers.HandleFunc("/wallet/fiat/refund/request", fiatRefundRequest).Methods("POST")
	routers.HandleFunc("/wallet/fiat/refund/approve", fiatRefundApprove).Methods("POST")
	routers.HandleFunc("/wallet/fiat/refund/list", fiatRefundList).Methods("GET")
	routers.HandleFunc("/wallet/fiat/withdraw", fiatWithdraw).Methods("POST")
	routers.HandleFunc("/wallet/fiat/balance/{user_id}", fiatWalletBalance).Methods("GET")

	routers.HandleFunc("/wallet/transaction_types", transactionTypes).Methods("GET")
	routers.HandleFunc("/wallet/fiat_currencies", fiatCurrencies).Methods("GET")

	return routers
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
