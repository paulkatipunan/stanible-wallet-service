package utils

import (
	"context"
	"strconv"
	"strings"

	"api.stanible.com/wallet/database"
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

func InsertFiatTransactionRecord(transactionPayload models.Transaction_payload) error {
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
		tx.Rollback()
		return err
	} else {
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
				a.active = true
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
			a.active = true
	`
	// NOTE:
	// WHERE
	// 		tt.type IN ('deposit', 'refund') AND
	// 		ft.status = 'success'
	db.QueryRow(sql, userId, userId).Scan(&balance)

	bal, err := strconv.Atoi(strings.Split(string(balance), ".")[0])

	return int32(bal), err
}
