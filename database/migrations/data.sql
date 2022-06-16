INSERT INTO
	transaction_types(pk_transaction_type_id, type)
VALUES
	('96c6b1c9-6c5a-493e-9645-64d2582a478b', 'deposit'),
	('242934fd-b51b-448e-ba31-74189172f056', 'withdraw'),
	('a457cb69-f670-4bf1-bf13-b99c82b0d170', 'buy'),
	('5eecd051-f51a-40b7-aaa3-81afb44150cc', 'refund'),
	('4f6116ee-9271-4b72-9a4c-70eba27fadcd', 'fee');

INSERT INTO
	fiat_currencies(pk_fiat_currency_id, name, symbol, numeric_precision)
VALUES
	('76eb713b-5bd1-4a54-968b-3897e88a50fb', 'Philippine Peso', 'PHP', '12,2');

INSERT INTO fiat_fee_types(fee_name) VALUES ('buy'), ('deposit');

-- CSV file:
COPY accounts(pk_account_id, user_id, type, description, active, created_at, updated_at)
FROM '/Users/barrylavides/Documents/stanible/stanible-wallet-service/database/migrations/accounts-1655355602612.csv'
DELIMITER ','
CSV HEADER;