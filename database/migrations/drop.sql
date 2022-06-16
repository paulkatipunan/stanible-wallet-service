DELETE FROM fiat_transactions_fee_assoc;
DELETE FROM fiat_transactions_assoc;
DELETE FROM fiat_transactions;
DELETE FROM transaction_types;
DELETE FROM fiat_currencies;
DELETE FROM accounts;
DELETE FROM fiat_fee_types;

DROP TABLE fiat_transactions_fee_assoc;
DROP TABLE fiat_transactions_assoc;
DROP TABLE fiat_transactions;
DROP TABLE transaction_types;
DROP TABLE fiat_currencies;
DROP TABLE accounts;
DROP TABLE fiat_fee_types;

DROP TYPE user_type;
DROP TYPE tx_status;
DROP TYPE tx_type;
DROP TYPE numeric_precision_type;