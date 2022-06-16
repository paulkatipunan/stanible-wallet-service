INSERT INTO
	transaction_types(pk_transaction_type_id, type)
VALUES
	('96c6b1c9-6c5a-493e-9645-64d2582a478b', 'deposit'),
	('242934fd-b51b-448e-ba31-74189172f056', 'withdraw'),
	('a457cb69-f670-4bf1-bf13-b99c82b0d170', 'buy'),
	('5eecd051-f51a-40b7-aaa3-81afb44150cc', 'refund'),
	('4f6116ee-9271-4b72-9a4c-70eba27fadcd', 'fee');

INSERT INTO
	fiat_currencies(pk_fiat_currency_id, name, symbol)
VALUES
	('76eb713b-5bd1-4a54-968b-3897e88a50fb', 'Philippine Peso', 'PHP');

INSERT INTO
	accounts(user_id, type)
VALUES
	('2681d82e-dd66-4357-96dc-ee5c7b7a6797', 'gcash_deposit'),
	('55e34c28-d3f9-4161-a910-dc643adfffd3', 'gcash_withdraw'),
	('ae582449-ec97-4c41-812b-d5c25d26c882', 'grabpay_withdraw'),
	('2c5d991b-909c-4cee-a652-8e5c3a0ccbf2', 'grabpay_deposit'),
	('d5e4698f-0cf4-4ab3-80da-5aa0d14876d4', 'paymaya_deposit'),
	('9b955c81-17ce-487d-9bf2-9686349ca652', 'paymaya_withdraw'),
	('37fda4f6-acdb-4411-995e-305e226dd4c9', 'treasury'),
	('9aa8ed53-dc51-448c-82fe-5f017f1c18fb', 'regular_user'),
	('c59c8ca5-8d67-4d03-ab12-8a824ceb754e', 'creator'),
	('121ab07d-bcf8-46a8-a111-ce053bc0eb69', 'regular_user'),
	('aa9b25dd-951c-4845-bb93-30c9d4bb4ca1', 'admin');

INSERT INTO fiat_fee_types(fee_name) VALUES ('buy'), ('deposit');

-- CSV file:
COPY accounts(pk_account_id, user_id, type, description, active, created_at, updated_at)
FROM '/Users/barrylavides/Documents/stanible/stanible-wallet-service/database/migrations/accounts-1655355602612.csv'
DELIMITER ','
CSV HEADER;