package payment

import (
	"strings"

	"github.com/shopspring/decimal"
	"individual-training-backend-roadmap/week1/solution/excercise1/domain"
)

const (
	MaxPaymentAmount    = 10_000_000
	LoyaltyPointsPerUnit = 1000
	MinCardNumberLength = 13
	MaxCardNumberLength = 19
)

var taxRates = map[domain.CurrencyType]float64{
	domain.CurrencyUSD: 0.10,
	domain.CurrencyVND: 0.08,
	domain.CurrencySGD: 0.07,
}

const defaultTaxRate = 0.10

// PaymentValidator validates payment request fields.
type PaymentValidator struct{}

func NewPaymentValidator() *PaymentValidator {
	return &PaymentValidator{}
}

func (v *PaymentValidator) ValidateAmount(amount *domain.Money) error {
	if amount == nil || !amount.IsPositive() {
		return domain.ErrInvalidAmount
	}
	maxAmount := domain.NewMoney(decimal.NewFromInt(MaxPaymentAmount), amount.Currency)
	if amount.GreaterThan(maxAmount) {
		return domain.ErrAmountTooLarge
	}
	return nil
}

func (v *PaymentValidator) ValidateCreditCard(card *domain.CreditCard) error {
	if card == nil {
		return domain.ErrInvalidCardNumber
	}
	if len(card.Number) < MinCardNumberLength || len(card.Number) > MaxCardNumberLength {
		return domain.ErrInvalidCardNumber
	}
	if len(card.CVV) != 3 && len(card.CVV) != 4 {
		return domain.ErrInvalidCVV
	}
	parts := strings.Split(card.Expiry, "/")
	if len(parts) != 2 {
		return domain.ErrInvalidExpiry
	}
	return nil
}

// TaxCalculator computes tax based on currency-specific rates.
type TaxCalculator struct{}

func NewTaxCalculator() *TaxCalculator {
	return &TaxCalculator{}
}

func (t *TaxCalculator) Calculate(amount *domain.Money) *domain.Money {
	rate, ok := taxRates[amount.Currency]
	if !ok {
		rate = defaultTaxRate
	}
	return amount.Multiply(rate)
}
