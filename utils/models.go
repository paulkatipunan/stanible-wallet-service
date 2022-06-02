package utils

import "api.stanible.com/wallet/models"

func FiatPayloadConverter(transactionPayload models.Transaction_payload, actualAmount int32) models.Transaction_data {
	var transactionData models.Transaction_data

	transactionData.Transaction_type_id = transactionPayload.Transaction_type_id
	transactionData.Sender_user_id = transactionPayload.Sender_user_id
	transactionData.Receiver_user_id = transactionPayload.Receiver_user_id
	transactionData.Fiat_currency_id = transactionPayload.Fiat_currency_id
	transactionData.Total_amount = transactionPayload.Amount
	transactionData.Actual_amount = actualAmount
	transactionData.Ramp_tx_id = transactionPayload.Ramp_tx_id

	return transactionData
}
