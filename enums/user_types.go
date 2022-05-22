package enums

func mergeUserTypes(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func AllUserTypes() map[string]string {
	return mergeUserTypes(NonPaymentUserTypes, PaymentTypes)
}

var NonPaymentUserTypes = map[string]string{
	"TREASURY":     "treasury",
	"ADMIN":        "admin",
	"CREATOR":      "creator",
	"REGULAR_USER": "regular_user",
}

var PaymentTypes = map[string]string{
	"GCASH":     "gcash",
	"PAYMAYA":   "paymaya",
	"DRAGONPAY": "dragonpay",
}
