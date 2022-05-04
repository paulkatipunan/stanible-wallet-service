INSERT INTO
	transaction_types(type)
VALUES
	('deposit'),
	('withdraw'),
	('buy'),
	('refund');

INSERT INTO
	fiat_currencies(name, symbol)
VALUES
	('Philippine Peso', 'PHP');

INSERT INTO
	accounts(user_id, type)
VALUES
	('2681d82e-dd66-4357-96dc-ee5c7b7a6797', 'gcash'),
	('37fda4f6-acdb-4411-995e-305e226dd4c9', 'treasury'),
	('9aa8ed53-dc51-448c-82fe-5f017f1c18fb', 'regular_user'),
	('121ab07d-bcf8-46a8-a111-ce053bc0eb69', 'regular_user'),
	('aa9b25dd-951c-4845-bb93-30c9d4bb4ca1', 'admin');