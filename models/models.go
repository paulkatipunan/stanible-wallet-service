package models

import "database/sql"

type Accounts struct {
	Pk_account_id string         `json:"pk_account_id"`
	User_id       string         `json:"user_id"`
	Type          string         `json:"type"`
	Description   sql.NullString `json:"description"`
	Active        bool           `json:"active"`
	Created_at    string         `json:"created_at"`
	Updated_at    string         `json:"updated_at"`
}

type Fiat_transactions struct {
	Pk_fiat_transaction_id string  `json:"pk_fiat_transaction_id"`
	Fk_account_id          string  `json:"fk_account_id"`
	Fk_transaction_type_id string  `json:"fk_transaction_type_id"`
	Fk_fiat_currency_id    float32 `json:"fk_fiat_currency_id"`

	Active     bool   `json:"active"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type Fiat_ramp_logs struct {
	Pk_fiat_ramp_logs_id   string `json:"pk_fiat_ramp_logs_id"`
	Fk_fiat_transaction_id string `json:"fk_fiat_transaction_id"`
	Ramp_tx_id             string `json:"ramp_tx_id"`

	Description sql.NullString `json:"description"`

	Active     bool   `json:"active"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type Transaction_payload struct {
	Transaction_type_id string `json:"transaction_type_id"`
	Sender_user_id      string `json:"sender_user_id"`
	Receiver_user_id    string `json:"receiver_user_id"`
	Fiat_currency_id    string `json:"fiat_currency_id"`
	Amount              int32  `json:"amount"`
	Ramp_tx_id          string `json:"ramp_tx_id"`
}
