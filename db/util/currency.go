package util

const (
	USD = "USD"
	EUR = "EUR"
	INR = "INR"
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, INR, CAD:
		return true
	}
	return false
}
