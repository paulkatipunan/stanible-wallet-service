package utils

import (
	"strconv"
)

func FeeCalculator(transactionType string, totalAmount int32) (float64, float64) {
	// get percentage_fee
	_, percentage := GetTransactionFeeType(transactionType)
	res, _ := strconv.ParseFloat(percentage, 64)
	// percent to decimal
	decimal := res / 100
	fee := float64(totalAmount) * decimal
	actualAmount := float64(totalAmount) - fee

	return actualAmount, fee
}
