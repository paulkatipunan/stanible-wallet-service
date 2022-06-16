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
	('aa9b25dd-951c-4845-bb93-30c9d4bb4ca1', 'admin'),
	('2da7d4b4-8949-4b27-b44f-4354a9f4858b','creator'),
    ('9b4e595c-d450-425a-bda9-c2595a1868e8','regular_user'),
    ('e61ed038-9507-4e64-961f-da95f9a161d9','creator'),
    ('d6709db3-f97d-485d-b33b-646bac7e118c','creator'),
    ('8c9519e3-1da7-4766-9e28-ce4638638c67','creator'),
    ('3c8eb348-4aca-443f-a3ce-420b27df44cb','regular_user'),
    ('01d91e53-c2db-42ca-b69c-e8a114fe9c7f','creator'),
    ('746e0f9a-d54f-496d-a2ec-bd32d44e08c6','regular_user'),
    ('95ba627f-17d0-4c4e-8779-f3ff39701119','regular_user'),
    ('362a253a-d141-482c-93bf-0dc03bd8354d','creator'),
    ('94fef150-be9e-4c95-b4e8-3477c2050bbd','creator'),
    ('61da0e84-ed93-4ac6-bc75-4dfea69353cd','regular_user'),
    ('889184a8-eda4-4565-84da-b57652541950','creator'),
    ('9b341a38-5745-42d2-84a4-3c47b6c9e97d','creator'),
    ('4df651f5-d842-4ba9-8eb7-ad2b40cc588f','creator'),
    ('859f980b-3874-4c12-ba9e-bdfe4a6857c6','creator'),
    ('c143c152-486f-4b52-8aef-ad0c9a5f7996','creator'),
    ('179478f8-aa3d-4186-b21a-d7841e5642d0','regular_user'),
    ('799a240c-acd8-4836-8a10-547d7678ee1a','regular_user'),
    ('d889a031-b03b-4a06-bb62-e9d130e88dc1','regular_user'),
    ('106aa487-cc01-43c6-ac14-ee3015e506cb','creator'),
    ('b8dc06b2-4b9e-4283-a80c-e5346c74f9b6','regular_user'),
    ('604cdfd1-0f35-4068-9a3c-f68027f57f61','creator'),
    ('bed8535e-c4aa-4cc6-867d-484f45e7d147','creator'),
    ('ed6c9d7b-7b7f-4629-b7e9-1cb9af116a5e','regular_user'),
    ('76eb5209-1fee-4951-9784-605f512ff584','regular_user'),
    ('9123f0a4-25ae-4ab1-91b6-80b49213dab1','regular_user'),
    ('48a0939f-920c-4e79-8ee1-6aba2df95ffa','regular_user'),
    ('c6563216-9727-4508-89bf-329e2c217c53','regular_user'),
    ('b2761fe1-c52b-4b60-bae9-07858799f84a','regular_user'),
    ('e3a45cea-b481-443e-a7d8-9515f153d735','regular_user'),
    ('302e3d22-5d4a-463f-80ca-4e3881547755','regular_user'),
    ('b78facac-3513-488f-8cc9-a932c770f003','regular_user'),
    ('6f9d19e9-ce50-42d1-9511-92bd8d18574c','regular_user'),
    ('ebfa359f-21ef-43e0-a613-d0458d67fecd','regular_user'),
    ('527dc87a-2662-48cf-8c2d-a2c2965f4e7d','regular_user'),
    ('d74d691a-f7b6-403d-aea3-c4fd9e61969c','regular_user'),
    ('362b38b6-a8dc-4bb4-a8cc-64feb27af480','regular_user'),
    ('aed1d636-87f1-4b4f-8e62-8e413d03a1e3','regular_user'),
    ('11404a16-c848-4542-93d6-91b494c5126a','regular_user'),
    ('be2e84b2-2b2b-4e6a-87fc-c20c1b3c0f93','regular_user'),
    ('6c3501b9-e4b1-4aab-bf03-83cbc375368b','regular_user'),
    ('e3c3dcc9-527d-4ce8-9585-07fc32b92243','creator'),
    ('7071b9e0-ac57-4f15-817f-6ca82110bc33','creator'),
    ('9116b826-b7f6-424b-a014-ed236e17015e','creator'),
    ('ae86a6cd-c5f8-4615-aa8f-e6ab3b94a1a4','creator'),
    ('a5ee415c-9ff7-4514-90f9-2641bdce05c2','creator');

INSERT INTO fiat_fee_types(fee_name) VALUES ('buy'), ('deposit');

-- CSV file:
COPY accounts(pk_account_id, user_id, type, description, active, created_at, updated_at)
FROM '/Users/barrylavides/Documents/stanible/stanible-wallet-service/database/migrations/accounts-1655355602612.csv'
DELIMITER ','
CSV HEADER;