package utils

import (
	"context"
	"strconv"
	"strings"

	"api.stanible.com/wallet/database"
	"api.stanible.com/wallet/enums"
	"api.stanible.com/wallet/models"
	"github.com/google/uuid"
)

func GetTransactionType(txType string) (string, string) {
	db := database.CreateConnection()
	defer db.Close()

	var pk_transaction_type_id, tx_type string
	sql := `SELECT pk_transaction_type_id, type as type_name FROM transaction_types WHERE type=$1`
	db.QueryRow(sql, txType).Scan(&pk_transaction_type_id, &tx_type)

	return pk_transaction_type_id, tx_type
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

func InsertFiatTransactionRecord(transactionPayload models.Transaction_payload, status string) error {
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
			pk_fiat_transaction_id, fk_user_id, fk_transaction_type_id, fk_fiat_currency_id, amount, status
		) VALUES
			($1, $2, $3, $4, $5, $6),
			($7, $8, $9, $10, $11, $12)
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
		status_value,

		pkReceiverId,
		transactionPayload.Receiver_user_id,
		transactionPayload.Transaction_type_id,
		transactionPayload.Fiat_currency_id,
		transactionPayload.Amount,
		status_value,
	)
	if err != nil {
		// Rollback if error
		tx.Rollback()
		return err
	} else {
		rows.Close()
	}

	// Insert fiat_transations_assoc record
	sqlInsertFiatTransactionAssoc := `
		INSERT INTO fiat_transactions_assoc (
			pk_sender_fiat_transaction_id, pk_receiver_fiat_transaction_id, ramp_tx_id, status
		) VALUES
			($1, $2, $3, $4)
	`
	row, err := tx.Query(
		sqlInsertFiatTransactionAssoc,
		pkSenderId,
		pkReceiverId,
		transactionPayload.Ramp_tx_id,
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

func AccountBalance(userId string) (int32, error) {
	db := database.CreateConnection()
	defer db.Close()

	var balance []uint8
	sql := `
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

	bal, err := strconv.Atoi(strings.Split(string(balance), ".")[0])

	return int32(bal), err
}

func RefundRequestList() []models.RefundRequestListModel {
	db := database.CreateConnection()
	defer db.Close()

	var refund_request_list []models.RefundRequestListModel

	sql := `
		SELECT
			fta.pk_fiat_transations_assoc_id,
			ft_receiver.fk_user_id,
			CAST(ft_receiver.amount as Integer),
			tt.type as type_name,
			fta.status,
			fta.created_at
		FROM
			fiat_transactions_assoc fta
		LEFT JOIN
			fiat_transactions ft_sender
			ON
				fta.pk_sender_fiat_transaction_id = ft_sender.pk_fiat_transaction_id
		LEFT JOIN
			fiat_transactions ft_receiver
			ON
				fta.pk_receiver_fiat_transaction_id = ft_receiver.pk_fiat_transaction_id
		LEFT JOIN
			transaction_types tt
			ON
				ft_sender.fk_transaction_type_id = tt.pk_transaction_type_id
		WHERE
			fta.status = 'pending' AND
			ft_sender.status = 'pending' AND
			tt.type = 'refund'
	`
	rows, _ := db.Query(sql)

	defer rows.Close()
	// iterate over the rows
	for rows.Next() {
		var refundRequest models.RefundRequestListModel

		_ = rows.Scan(
			&refundRequest.Pk_fiat_transations_assoc_id,
			&refundRequest.Fk_user_id,
			&refundRequest.Amount,
			&refundRequest.Type_name,
			&refundRequest.Status,
			&refundRequest.Created_at,
		)

		refund_request_list = append(refund_request_list, refundRequest)
	}

	return refund_request_list
}

func RefundApprove(transactionPayload models.RefundApprove_payload) error {
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
			pk_fiat_transations_assoc_id=$2
		RETURNING
			pk_sender_fiat_transaction_id,
			pk_receiver_fiat_transaction_id;
	`
	var pk_sender_fiat_transaction_id, pk_receiver_fiat_transaction_id string
	tx.QueryRow(sqlUpdateFiatTransactionAssoc, status_value, transactionPayload.Refund_id).Scan(&pk_sender_fiat_transaction_id, &pk_receiver_fiat_transaction_id)

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
	tx.QueryRow(sqlUpdateFiatTransaction, status_value, pk_sender_fiat_transaction_id, pk_receiver_fiat_transaction_id)

	// Commit the change if all queries ran successfully
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
