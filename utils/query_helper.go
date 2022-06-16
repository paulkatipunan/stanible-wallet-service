package utils

import (
	"context"
	"fmt"

	"api.stanible.com/wallet/database"
	"api.stanible.com/wallet/enums"
	"api.stanible.com/wallet/models"
	"github.com/google/uuid"
)

func UserUpdate(user_id string, type_name string) (string, error) {
	db := database.CreateConnection()
	defer db.Close()

	sql := `
		UPDATE accounts SET type=$1 WHERE user_id=$2 RETURNING user_id;
	`

	var returned_user_id string
	err := db.QueryRow(sql, type_name, user_id).Scan(&returned_user_id)

	return returned_user_id, err
}

func GetTransactionType(txType string) (string, string) {
	db := database.CreateConnection()
	defer db.Close()

	var pk_transaction_type_id, tx_type string
	sql := `SELECT pk_transaction_type_id, type as type_name FROM transaction_types WHERE type=$1`
	db.QueryRow(sql, txType).Scan(&pk_transaction_type_id, &tx_type)

	return pk_transaction_type_id, tx_type
}

func GetTransactionFeeType(txType string) (string, string) {
	db := database.CreateConnection()
	defer db.Close()

	var pk_fiat_fee_type_id, percentage string
	sql := `SELECT pk_fiat_fee_type_id, percentage FROM fiat_fee_types WHERE fee_name=$1`
	db.QueryRow(sql, txType).Scan(&pk_fiat_fee_type_id, &percentage)

	return pk_fiat_fee_type_id, percentage
}

func ActiveSenderReceiver(sender_user_id string, receiver_user_id string) int {
	db := database.CreateConnection()
	defer db.Close()

	// Check if sender and receiver addresses are both active
	sql := `
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

	var total_users int
	db.QueryRow(sql, sender_user_id, receiver_user_id).Scan(&total_users)

	return total_users
}

func UserTypes(sender_user_id string, receiver_user_id string) ([]string, error) {
	db := database.CreateConnection()
	defer db.Close()

	// Get user types
	sql := `
		SELECT (
			SELECT
				type
			FROM
				accounts
			WHERE
				user_id = $1 AND
				active = true
			) as sender_type, (
			SELECT
				type
			FROM
				accounts
			WHERE
				user_id = $2 AND
				active = true
			) as receiver_type
	`
	rows, err := db.Query(sql, sender_user_id, receiver_user_id)

	if err != nil {
		return []string{}, err
	} else {
		defer rows.Close()
		var sender_type, receiver_type string

		for rows.Next() {
			rows.Scan(&sender_type, &receiver_type)
		}
		return []string{sender_type, receiver_type}, nil
	}
}

func InsertFiatTransactionRecord(transactionData models.Transaction_data, status string) error {
	// Begin tx
	ctx := context.Background()
	db := database.CreateConnection()
	tx, err := db.BeginTx(ctx, nil)
	defer db.Close()

	if err != nil {
		return err
	}

	status_value := status

	// Insert fiat_transaction record
	sqlInsertFiatTransaction := `
		INSERT INTO fiat_transactions (
			pk_fiat_transaction_id,
			fk_user_id,
			fk_transaction_type_id,
			fk_fiat_currency_id,
			total_amount,
			actual_amount,
			status
		) VALUES
			($1, $2, $3, $4, $5, $6, $7),
			($8, $9, $10, $11, $12, $13, $14)
	`
	pkSenderId := uuid.New()
	pkReceiverId := uuid.New()
	rows, err := tx.Query(
		sqlInsertFiatTransaction,

		pkSenderId,
		transactionData.Sender_user_id,
		transactionData.Transaction_type_id,
		transactionData.Fiat_currency_id,
		-transactionData.Total_amount,
		-transactionData.Actual_amount,
		status_value,

		pkReceiverId,
		transactionData.Receiver_user_id,
		transactionData.Transaction_type_id,
		transactionData.Fiat_currency_id,
		transactionData.Total_amount,
		transactionData.Actual_amount,
		status_value,
	)
	if err != nil {
		// Rollback if error
		tx.Rollback()
		return err
	} else {
		rows.Close()
	}

	// Insert fiat_transactions_assoc record
	sqlInsertFiatTransactionAssoc := `
		INSERT INTO fiat_transactions_assoc (
			fk_sender_fiat_transaction_id, fk_receiver_fiat_transaction_id, ramp_tx_id, status
		) VALUES
			($1, $2, $3, $4)
	`
	row, err := tx.Query(
		sqlInsertFiatTransactionAssoc,
		pkSenderId,
		pkReceiverId,
		transactionData.Ramp_tx_id,
		status_value,
	)
	if err != nil {
		// Rollback if error
		tx.Rollback()
		return err
	} else {
		row.Close()
	}

	// Commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Insert two transactions
// 1. normal insert to fiat_transactions and fiat_transactions_assoc
//	- sender buyer
//	- receiver seller
//	- actual_amount from fee calculator
// 2.1 new insert to fiat_transactions
// 	- transaction_type_id = fee id
//	- sender seller
//	- receiver stanible
//	- actual_amount from fee calculator
// 2.2 new insert to fiat_transactions_fee_assoc
//	- fk_fiat_transactions_assoc_id
//	- fk_sender_fiat_transaction_id
//	- fk_receiver_fiat_transaction_id
func InsertFiatTransactionWithFeeRecord(
	transactionData models.Transaction_data,
	status string,
	feeAmount int32,
	actualAmount int32,
	feeRecipient string,
) error {
	// Begin tx
	ctx := context.Background()
	db := database.CreateConnection()
	tx, err := db.BeginTx(ctx, nil)
	defer db.Close()

	if err != nil {
		return err
	}

	// Insert fiat_transaction record
	sqlInsertFiatTransaction := `
		INSERT INTO fiat_transactions (
			pk_fiat_transaction_id,
			fk_user_id,
			fk_transaction_type_id,
			fk_fiat_currency_id,
			total_amount,
			actual_amount,
			status
		) VALUES
			($1, $2, $3, $4, $5, $6, $7),
			($8, $9, $10, $11, $12, $13, $14)
	`
	pkSenderId := uuid.New()
	pkReceiverId := uuid.New()
	rows, err := tx.Query(
		sqlInsertFiatTransaction,

		pkSenderId,
		transactionData.Sender_user_id,
		transactionData.Transaction_type_id,
		transactionData.Fiat_currency_id,
		-transactionData.Total_amount,
		-actualAmount,
		status,

		pkReceiverId,
		transactionData.Receiver_user_id,
		transactionData.Transaction_type_id,
		transactionData.Fiat_currency_id,
		transactionData.Total_amount,
		actualAmount,
		status,
	)
	if err != nil {
		// Rollback if error
		tx.Rollback()

		fmt.Println("error 01")
		return err
	} else {
		rows.Close()
	}

	// Insert fiat_transactions_assoc record
	pkTransactionAssoccId := uuid.New()
	sqlInsertFiatTransactionAssoc := `
		INSERT INTO fiat_transactions_assoc (
			pk_fiat_transactions_assoc_id,
			fk_sender_fiat_transaction_id,
			fk_receiver_fiat_transaction_id,
			ramp_tx_id,
			status
		) VALUES
			($1, $2, $3, $4, $5)
	`
	row, err := tx.Query(
		sqlInsertFiatTransactionAssoc,

		pkTransactionAssoccId,
		pkSenderId,
		pkReceiverId,
		transactionData.Ramp_tx_id,
		status,
	)
	if err != nil {
		// Rollback if error
		tx.Rollback()

		fmt.Println("error 02")
		return err
	} else {
		row.Close()
	}

	// Get transactioin_type_id of type fee
	var tx_fee_id, tx_type string
	sqlTransactioinTypeId := `SELECT pk_transaction_type_id as tx_fee_id, type as tx_type FROM transaction_types WHERE type = $1`
	db.QueryRow(sqlTransactioinTypeId, enums.FEE).Scan(&tx_fee_id, &tx_type)

	// Insert insert to fiat_transactions for fee record
	sqlInsertFiatFeeTransaction := `
		INSERT INTO fiat_transactions (
			pk_fiat_transaction_id,
			fk_user_id,
			fk_transaction_type_id,
			fk_fiat_currency_id,
			total_amount,
			actual_amount,
			status
		) VALUES
			($1, $2, $3, $4, $5, $6, $7),
			($8, $9, $10, $11, $12, $13, $14)
	`
	pkSenderIdFee := uuid.New()
	pkReceiverIdFee := uuid.New()
	rowsFee, err := tx.Query(
		sqlInsertFiatFeeTransaction,

		pkSenderIdFee,
		transactionData.Receiver_user_id,
		tx_fee_id,
		transactionData.Fiat_currency_id,
		-transactionData.Total_amount,
		-feeAmount,
		status,

		pkReceiverIdFee,
		feeRecipient,
		tx_fee_id,
		transactionData.Fiat_currency_id,
		transactionData.Total_amount,
		feeAmount,
		status,
	)
	if err != nil {
		// Rollback if error
		tx.Rollback()

		fmt.Println("error 03")
		return err
	} else {
		rowsFee.Close()
	}
	// Insert fiat_transactions_fee_assoc record
	sqlInsertFiatTransactionFeeAssoc := `
		INSERT INTO fiat_transactions_fee_assoc (
			fk_fiat_transactions_assoc_id,
			fk_sender_fiat_transaction_id,
			fk_receiver_fiat_transaction_id,
			status
		) VALUES
			($1, $2, $3, $4)
	`
	rowFeeAssoc, err := tx.Query(
		sqlInsertFiatTransactionFeeAssoc,

		pkTransactionAssoccId,
		pkSenderIdFee,
		pkReceiverIdFee,
		status,
	)
	if err != nil {
		// Rollback if error
		tx.Rollback()

		fmt.Println("error 04")
		return err
	} else {
		rowFeeAssoc.Close()
	}

	// Commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func AccountBalance(userId string, fiat_currency_id string) (int32, string, string, error) {
	db := database.CreateConnection()
	defer db.Close()

	var balance int32
	var pk_account_id string
	var symbol string
	var numeric_precision string

	sqlUser := `SELECT pk_account_id FROM accounts WHERE user_id = $1`
	errUser := db.QueryRow(sqlUser, userId).Scan(&pk_account_id)

	if errUser != nil {
		return balance, "", "", errUser
	}

	sqlCurrency := `SELECT symbol, numeric_precision FROM fiat_currencies WHERE pk_fiat_currency_id = $1`
	errCurrency := db.QueryRow(sqlCurrency, fiat_currency_id).Scan(&symbol, &numeric_precision)
	if errUser != nil {
		return balance, "", "", errCurrency
	}

	sql := `
		SELECT
			COALESCE(CAST(SUM(ft.total_amount) as Integer), 0) + (
				SELECT
					COALESCE(CAST(SUM(ft.total_amount) as Integer), 0) as balance
				FROM
					fiat_transactions ft
				INNER JOIN
					transaction_types tt
					ON
						ft.fk_transaction_type_id = tt.pk_transaction_type_id
				LEFT JOIN
					accounts a
					ON
						a.user_id = ft.fk_user_id
				WHERE
					ft.fk_user_id = $1 AND
					tt.type IN ('withdraw', 'buy') AND
					a.active = true AND
					ft.status = 'success'
			) as balance
		FROM
			fiat_transactions ft
		INNER JOIN
			transaction_types tt
			ON
				ft.fk_transaction_type_id = tt.pk_transaction_type_id
		LEFT JOIN
			accounts a
			ON
				a.user_id = ft.fk_user_id
		WHERE
			ft.fk_user_id = $2 AND
			tt.type IN ('deposit', 'refund') AND
			a.active = true AND
			ft.status = 'success'
	`
	db.QueryRow(sql, userId, userId).Scan(&balance)
	return balance, symbol, numeric_precision, nil
}

func RequestList(transaction_type string) []models.RequestListModel {
	db := database.CreateConnection()
	defer db.Close()

	var request_list []models.RequestListModel

	sql := `
		SELECT
			fta.pk_fiat_transactions_assoc_id,
			ft_receiver.fk_user_id,
			CAST(ft_receiver.actual_amount as Integer),
			tt.type as type_name,
			fta.status,
			fta.created_at
		FROM
			fiat_transactions_assoc fta
		LEFT JOIN
			fiat_transactions ft_sender
			ON
				fta.fk_sender_fiat_transaction_id = ft_sender.pk_fiat_transaction_id
		LEFT JOIN
			fiat_transactions ft_receiver
			ON
				fta.fk_receiver_fiat_transaction_id = ft_receiver.pk_fiat_transaction_id
		LEFT JOIN
			transaction_types tt
			ON
				ft_sender.fk_transaction_type_id = tt.pk_transaction_type_id
		WHERE
			fta.status = 'pending' AND
			ft_sender.status = 'pending' AND
			tt.type = $1
	`
	rows, _ := db.Query(sql, transaction_type)

	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var refundRequest models.RequestListModel

		_ = rows.Scan(
			&refundRequest.Pk_fiat_transactions_assoc_id,
			&refundRequest.Fk_user_id,
			&refundRequest.Amount,
			&refundRequest.Type_name,
			&refundRequest.Status,
			&refundRequest.Created_at,
		)

		request_list = append(request_list, refundRequest)
	}

	return request_list
}

func RequestApprove(transactionPayload models.RequestApprove_payload) error {
	// Begin tx
	ctx := context.Background()
	db := database.CreateConnection()
	tx, err := db.BeginTx(ctx, nil)
	defer db.Close()

	if err != nil {
		return err
	}

	var status_value string

	if transactionPayload.Status == "approve" {
		status_value = enums.TX_STATUS["SUCCESS"]
	} else if transactionPayload.Status == "cancel" {
		status_value = enums.TX_STATUS["CANCELLED"]
	} else {
		status_value = enums.TX_STATUS["FAILED"]
	}

	// Update fiat_transactions_assoc record
	sqlUpdateFiatTransactionAssoc := `
		UPDATE
			fiat_transactions_assoc
		SET
			status=$1
		WHERE
			pk_fiat_transactions_assoc_id=$2
		RETURNING
			fk_sender_fiat_transaction_id,
			fk_receiver_fiat_transaction_id;
	`
	var fk_sender_fiat_transaction_id, fk_receiver_fiat_transaction_id string
	tx.QueryRow(sqlUpdateFiatTransactionAssoc, status_value, transactionPayload.Request_id).Scan(&fk_sender_fiat_transaction_id, &fk_receiver_fiat_transaction_id)

	// Update fiat_transactions record
	sqlUpdateFiatTransaction := `
		UPDATE
			fiat_transactions
		SET
			status=$1
		WHERE
			pk_fiat_transaction_id IN (
				$2,
				$3
			)
	`
	tx.QueryRow(sqlUpdateFiatTransaction, status_value, fk_sender_fiat_transaction_id, fk_receiver_fiat_transaction_id)

	// Commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func TransactionList(
	user_id string,
	offset string,
	limit string,
	tx_type string,
	dateStart string,
	dateEnd string,
) ([]models.Fiat_transaction_list_model, error) {
	db := database.CreateConnection()
	defer db.Close()

	var fiat_transaction_list []models.Fiat_transaction_list_model
	var pk_account_id string

	sqlUser := `SELECT pk_account_id FROM accounts WHERE user_id = $1`
	errUser := db.QueryRow(sqlUser, user_id).Scan(&pk_account_id)

	if errUser != nil {
		return fiat_transaction_list, errUser
	}

	var sql string = `
		SELECT
			ft.pk_fiat_transaction_id as transaction_id,
			fc.symbol,
			fc.numeric_precision,
			CAST(ft.total_amount as Integer),
			tt.type as transaction_type,
			COALESCE(fta_sender.ramp_tx_id, fta_receiver.ramp_tx_id) as reference_number,
			ft.status,
			ft.created_at
		FROM
			fiat_transactions ft
		INNER JOIN
			transaction_types tt
			ON
				ft.fk_transaction_type_id = tt.pk_transaction_type_id
		LEFT JOIN
			fiat_transactions_assoc fta_sender
			ON
				fta_sender.fk_sender_fiat_transaction_id = ft.pk_fiat_transaction_id
		LEFT JOIN
			fiat_transactions_assoc fta_receiver
			ON
				fta_receiver.fk_receiver_fiat_transaction_id = ft.pk_fiat_transaction_id
		LEFT JOIN
			accounts a
			ON
				a.user_id = ft.fk_user_id
		LEFT JOIN
			fiat_currencies fc
			ON
				fc.pk_fiat_currency_id = ft.fk_fiat_currency_id
		WHERE
			ft.fk_user_id = $1 AND
			a.active = true AND
			TO_DATE($2,'YYYY/MM/DD') <= TO_DATE(CAST(ft.created_at as TEXT),'YYYY-MM-DD') AND
			TO_DATE($3,'YYYY/MM/DD') >= TO_DATE(CAST(ft.created_at as TEXT),'YYYY-MM-DD') AND
	`

	if tx_type == enums.ALL {
		sql += `tt.type IN ('deposit', 'refund', 'withdraw', 'buy')`
	} else {
		sql += `tt.type = '` + tx_type + `'`
	}

	sql += `
		ORDER BY
			ft.created_at DESC
		OFFSET ($4-1)*$5 LIMIT $6
	`

	rows, err := db.Query(sql, user_id, dateStart, dateEnd, offset, limit, limit)

	if err != nil {
		defer rows.Close()
	}
	// iterate over the rows
	for rows.Next() {
		var fiatTransaction models.Fiat_transaction_list_model
		_ = rows.Scan(
			&fiatTransaction.Pk_fiat_transaction_id,
			&fiatTransaction.Fiat_currency_symbol,
			&fiatTransaction.Numeric_precision,
			&fiatTransaction.Total_amount,
			&fiatTransaction.Type,
			&fiatTransaction.Ramp_tx_id,
			&fiatTransaction.Status,
			&fiatTransaction.Created_at,
		)

		fiat_transaction_list = append(fiat_transaction_list, fiatTransaction)
	}

	return fiat_transaction_list, err
}
