package enums

const (
	DEPOSIT  = "deposit"
	BUY      = "buy"
	REFUND   = "refund"
	WITHDRAW = "withdraw"
	FEE      = "fee"
)

var TX_STATUS = map[string]string{
	"PENDING":   "pending",
	"SUCCESS":   "success",
	"FAILED":    "failed",
	"CANCELLED": "cancelled",
}

var FE_TX_STATUS = map[string]string{
	"APPROVE": "approve",
	"CANCEL":  "cancel",
}
