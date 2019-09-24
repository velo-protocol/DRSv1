package vxdr

type Currency string

const (
	CurrencyTHB Currency = "THB"
	CurrencySGD Currency = "SGD"
	CurrencyUSD Currency = "USD"
)

func (currency Currency) IsValid() bool {
	return currency == CurrencyTHB ||
		currency == CurrencySGD ||
		currency == CurrencyUSD
}
