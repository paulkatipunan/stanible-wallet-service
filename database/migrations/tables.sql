CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE accounts (
	pk_account_id UUID DEFAULT uuid_generate_v4(),

	user_id UUID UNIQUE NOT NULL,
	type user_type NOT NULL,
	description VARCHAR,
	
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_account_id)
);
CREATE TRIGGER update_table_modtime BEFORE UPDATE ON accounts FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();

CREATE TABLE transaction_types (
	pk_transaction_type_id UUID DEFAULT uuid_generate_v4(),

	type tx_type UNIQUE NOT NULL,
	description VARCHAR,

	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_transaction_type_id)
);
CREATE TRIGGER update_table_modtime BEFORE UPDATE ON transaction_types FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();

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
CREATE TRIGGER update_table_modtime BEFORE UPDATE ON fiat_currencies FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();

CREATE TABLE fiat_transactions (
	pk_fiat_transaction_id UUID DEFAULT uuid_generate_v4(),
	fk_user_id UUID NOT NULL REFERENCES accounts(user_id),
	fk_transaction_type_id UUID NOT NULL REFERENCES transaction_types(pk_transaction_type_id),
	fk_fiat_currency_id UUID NOT NULL REFERENCES fiat_currencies(pk_fiat_currency_id),
	
	total_amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,
	actual_amount NUMERIC(12, 2) DEFAULT 0 NOT NULL,
	status tx_status NOT NULL DEFAULT 'pending',
	
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_fiat_transaction_id)
);
CREATE TRIGGER update_table_modtime BEFORE UPDATE ON fiat_transactions FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();

CREATE TABLE fiat_transactions_assoc (
	pk_fiat_transactions_assoc_id UUID DEFAULT uuid_generate_v4(),

	fk_sender_fiat_transaction_id UUID NOT NULL REFERENCES fiat_transactions(pk_fiat_transaction_id),
	fk_receiver_fiat_transaction_id UUID NOT NULL REFERENCES fiat_transactions(pk_fiat_transaction_id),

	ramp_tx_id VARCHAR UNIQUE NOT NULL,
	description VARCHAR,
	status tx_status NOT NULL DEFAULT 'pending',

	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	UNIQUE (fk_sender_fiat_transaction_id, fk_receiver_fiat_transaction_id, ramp_tx_id),
	PRIMARY KEY (pk_fiat_transactions_assoc_id)
);
CREATE TRIGGER update_table_modtime BEFORE UPDATE ON fiat_transactions_assoc FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();

CREATE TABLE fiat_fee_types (
	pk_fiat_fee_type_id UUID DEFAULT uuid_generate_v4(),
	fee_name VARCHAR NOT NULL,
	percentage decimal(5,2) DEFAULT 20.00 NOT NULL,
	description VARCHAR,
	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (pk_fiat_fee_type_id)
);
CREATE TRIGGER update_table_modtime BEFORE UPDATE ON fiat_fee_types FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();

CREATE TABLE fiat_transactions_fee_assoc (
	pk_fiat_transactions_fee_assoc_id UUID DEFAULT uuid_generate_v4(),

	fk_fiat_transactions_assoc_id UUID NOT NULL REFERENCES fiat_transactions_assoc(pk_fiat_transactions_assoc_id),

	fk_sender_fiat_transaction_id UUID NOT NULL REFERENCES fiat_transactions(pk_fiat_transaction_id),
	fk_receiver_fiat_transaction_id UUID NOT NULL REFERENCES fiat_transactions(pk_fiat_transaction_id),

	description VARCHAR,
	status tx_status NOT NULL DEFAULT 'pending',

	active BOOLEAN DEFAULT TRUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	UNIQUE (
		fk_fiat_transactions_assoc_id,
		fk_sender_fiat_transaction_id,
		fk_receiver_fiat_transaction_id
	),
	PRIMARY KEY (pk_fiat_transactions_fee_assoc_id)
);
CREATE TRIGGER update_table_modtime BEFORE UPDATE ON fiat_transactions_fee_assoc FOR EACH ROW EXECUTE PROCEDURE on_update_trigger();