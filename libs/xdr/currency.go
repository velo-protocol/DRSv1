package vxdr

// Currency is a constant which defined supported fiat currency.
type Currency string

const (
	CurrencyTHB Currency = "THB"
	CurrencySGD Currency = "SGD"
	CurrencyUSD Currency = "USD"
)

// IsValid checks if the given currency is supported/valid or not.
func (currency Currency) IsValid() bool {
	return currency == CurrencyTHB ||
		currency == CurrencySGD ||
		currency == CurrencyUSD
}
