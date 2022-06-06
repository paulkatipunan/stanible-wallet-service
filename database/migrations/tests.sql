-- Get latest fiat transactions w/ a fee
SELECT
	ft.pk_fiat_transaction_id,
	ft.fk_user_id,
	tt.type,
	ft.total_amount,
	ft.actual_amount,
	ft.created_at
FROM
	fiat_transactions ft
LEFT JOIN
	transaction_types tt
	ON
		tt.pk_transaction_type_id = ft.fk_transaction_type_id
GROUP BY
	ft.pk_fiat_transaction_id,
	tt.type
ORDER BY
	created_at DESC
LIMIT 4;

SELECT * FROM fiat_transactions_assoc ORDER BY created_at DESC LIMIT 1;
SELECT * FROM fiat_transactions_fee_assoc ORDER BY created_at DESC LIMIT 1;