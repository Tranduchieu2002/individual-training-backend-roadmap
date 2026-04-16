package domain

import "github.com/shopspring/decimal"

/*
Rate is the exchange rate of the currency to the base currency (USD)
Precision is the number of decimal places of the currency
*/
const (
	CurrencyUSD CurrencyType = "USD"
	CurrencyVND CurrencyType = "VND"
	CurrencyEUR CurrencyType = "EUR"
	CurrencyGBP CurrencyType = "GBP"
	CurrencyJPY CurrencyType = "JPY"
	CurrencyKRW CurrencyType = "KRW"
	CurrencySGD CurrencyType = "SGD"
)

type CurrencyType string

/*
*

	type Currency struct {
		Code      CurrencyType    `json:"code"`
		Symbol    string          `json:"symbol"`
		Rate      decimal.Decimal `json:"rate"`
		Precision int             `json:"precision"`
		}

		func NewCurrency(code CurrencyType, symbol string, rate decimal.Decimal, precision int) *Currency {
			return &Currency{
				Code:      code,
				Symbol:    symbol,
				Rate:      rate,
				Precision: precision,
				}
		}
*/
type Money struct {
	Amount   decimal.Decimal
	Currency CurrencyType
}

func NewMoney(amount decimal.Decimal, currency CurrencyType) *Money {
	return &Money{
		Amount:   amount,
		Currency: currency,
	}
}

func (m *Money) Add(other *Money) *Money {
	return NewMoney(m.Amount.Add(other.Amount), m.Currency)
}

func (m *Money) Subtract(other *Money) *Money {
	return NewMoney(m.Amount.Sub(other.Amount), m.Currency)
}

func (m *Money) Multiply(factor float64) *Money {
	return NewMoney(m.Amount.Mul(decimal.NewFromFloat(factor)), m.Currency)
}

func (m *Money) Divide(divisor float64) *Money {
	return NewMoney(m.Amount.Div(decimal.NewFromFloat(divisor)), m.Currency)
}

func ZeroMoney(currency CurrencyType) *Money {
	return &Money{Amount: decimal.Zero, Currency: currency}
}

func (m *Money) IsZero() bool {
	return m.Amount.IsZero()
}

func (m *Money) IsPositive() bool {
	return m.Amount.IsPositive()
}

func (m *Money) IsNegative() bool {
	return m.Amount.IsNegative()
}

func (m *Money) GreaterThan(other *Money) bool {
	return m.Amount.GreaterThan(other.Amount)
}

func (m *Money) LessThan(other *Money) bool {
	return m.Amount.LessThan(other.Amount)
}

func (m *Money) LessThanOrEqual(other *Money) bool {
	return m.Amount.LessThanOrEqual(other.Amount)
}
