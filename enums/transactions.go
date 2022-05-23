package enums

const (
	DEPOSIT  = "deposit"
	BUY      = "buy"
	REFUND   = "refund"
	WITHDRAW = "withdraw"
)

var TX_STATUS = map[string]string{
	"PENDING":   "pending",
	"SUCCESS":   "success",
	"FAILED":    "failed",
	"CANCELLED": "cancelled",
}
