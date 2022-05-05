INSERT INTO
	transaction_types(pk_transaction_type_id, type)
VALUES
	('96c6b1c9-6c5a-493e-9645-64d2582a478b', 'deposit'),
	('242934fd-b51b-448e-ba31-74189172f056', 'withdraw'),
	('a457cb69-f670-4bf1-bf13-b99c82b0d170', 'buy'),
	('5eecd051-f51a-40b7-aaa3-81afb44150cc', 'refund');

INSERT INTO
	fiat_currencies(pk_fiat_currency_id, name, symbol)
VALUES
	('76eb713b-5bd1-4a54-968b-3897e88a50fb', 'Philippine Peso', 'PHP');

INSERT INTO
	accounts(user_id, type)
VALUES
	('2681d82e-dd66-4357-96dc-ee5c7b7a6797', 'gcash'),
	('37fda4f6-acdb-4411-995e-305e226dd4c9', 'treasury'),
	('9aa8ed53-dc51-448c-82fe-5f017f1c18fb', 'regular_user'),
	('121ab07d-bcf8-46a8-a111-ce053bc0eb69', 'regular_user'),
	('aa9b25dd-951c-4845-bb93-30c9d4bb4ca1', 'admin');