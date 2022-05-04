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
