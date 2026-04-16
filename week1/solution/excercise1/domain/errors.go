package domain

import "errors"

var (
	ErrInvalidAmount            = errors.New("invalid amount: must be greater than zero")
	ErrAmountTooLarge           = errors.New("amount exceeds maximum limit of 10,000,000")
	ErrInvalidCardNumber        = errors.New("invalid credit card number length")
	ErrInvalidCVV               = errors.New("invalid CVV length")
	ErrInvalidExpiry            = errors.New("invalid expiry date format")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserBanned               = errors.New("user is banned")
	ErrInsufficientBalance      = errors.New("insufficient wallet balance")
	ErrOrderNotFound            = errors.New("order not found")
	ErrOrderAlreadyProcessed    = errors.New("order has already been processed")
	ErrUnsupportedPaymentMethod = errors.New("unsupported payment method")
)
