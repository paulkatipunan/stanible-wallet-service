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

type Fiat_transaction_list_model struct {
	Pk_fiat_transaction_id string `json:"transaction_id"`
	Total_amount           int32  `json:"amount"`
	Actual_amount          int32  `json:"actual_amount"`
	Type                   string `json:"transation_type"`
	Ramp_tx_id             string `json:"ramp_tx_id"`
	Status                 string `json:"status"`
	Created_at             string `json:"created_at"`
}

type Transaction_payload struct {
	Transaction_type_id string `json:"transaction_type_id"`
	Sender_user_id      string `json:"sender_user_id"`
	Receiver_user_id    string `json:"receiver_user_id"`
	Fiat_currency_id    string `json:"fiat_currency_id"`
	Amount              int32  `json:"amount"`
	Reference_number    string `json:"reference_number"`
}

type Transaction_data struct {
	Transaction_type_id string `json:"transaction_type_id"`
	Sender_user_id      string `json:"sender_user_id"`
	Receiver_user_id    string `json:"receiver_user_id"`
	Fiat_currency_id    string `json:"fiat_currency_id"`
	Total_amount        int32  `json:"total_amount"`
	Actual_amount       int32  `json:"actual_amount"`
	Ramp_tx_id          string `json:"ramp_tx_id"`
}

type RequestApprove_payload struct {
	Request_id string `json:"request_id"`
	Status     string `json:"status"`
}

type Transaction_types struct {
	Pk_transaction_type_id string `json:"pk_transaction_type_id"`
	Type_name              string `json:"type_name"`
}

type Fiat_currencies struct {
	Pk_fiat_currency_id string `json:"pk_fiat_currency_id"`
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
}

type RequestListModel struct {
	Pk_fiat_transactions_assoc_id string `json:"request_id"`
	Fk_user_id                    string `json:"requesting_user_id"`
	Amount                        int32  `json:"amount"`
	Type_name                     string `json:"type_name"`
	Status                        string `json:"status"`
	Created_at                    string `json:"created_at"`
}
