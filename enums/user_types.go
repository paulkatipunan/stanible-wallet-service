package enums

func mergeUserTypes(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func AllUserTypes() map[string]string {
	non_customer_users := mergeUserTypes(SystemUserTypes, PaymentUserTypes)
	return mergeUserTypes(non_customer_users, CustomerUserTypes)
}

func NonPaymentUserTypes() map[string]string {
	return mergeUserTypes(SystemUserTypes, CustomerUserTypes)
}

func WithdrawalUserTypes() map[string]string {
	return map[string]string{
		"TREASURY": SystemUserTypes["TREASURY"],
		"CREATOR":  CustomerUserTypes["CREATOR"],
	}
}

var SystemUserTypes = map[string]string{
	"TREASURY": "treasury",
	"ADMIN":    "admin",
}

var CustomerUserTypes = map[string]string{
	"CREATOR":      "creator",
	"REGULAR_USER": "regular_user",
}

var PaymentUserTypes = map[string]string{
	"GCASH_DEPOSIT":    "gcash_deposit",
	"GCASH_WITHDRAW":   "gcash_withdraw",
	"PAYMAYA_DEPOSIT":  "paymaya_deposit",
	"PAYMAYA_WITHDRAW": "paymaya_withdraw",
	"GRABPAY_WITHDRAW": "grabpay_withdraw",
	"GRABPAY_DEPOSIT":  "grabpay_deposit",
}
