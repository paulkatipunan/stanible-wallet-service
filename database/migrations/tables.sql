CREATE TABLE accounts (
	pk_account_id UUID DEFAULT uuid_generate_v4(),
	user_id VARCHAR UNIQUE NOT NULL,
	type ramp_type,
	description VARCHAR,
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_account_id)
);

CREATE TABLE transaction_types (
	pk_transaction_type_id UUID DEFAULT uuid_generate_v4(),
	ramp tx_type,
	description VARCHAR,
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_transaction_type_id)
);

CREATE TABLE fiat_currencies (
	pk_fiat_currency_id UUID DEFAULT uuid_generate_v4(),
	name VARCHAR NOT NULL,
	symbol VARCHAR NOT NULL,
	description VARCHAR,
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_fiat_currency_id)
);

CREATE TABLE fiat_transactions (
	pk_fiat_transaction_id UUID DEFAULT uuid_generate_v4(),
	fk_account_id UUID NOT NULL REFERENCES accounts(pk_account_id),
	fk_transaction_type_id UUID NOT NULL REFERENCES transaction_types(pk_transaction_type_id),
	fk_fiat_currency_id UUID NOT NULL REFERENCES fiat_currencies(pk_fiat_currency_id),
	
	amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,
	
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_fiat_transaction_id)
);