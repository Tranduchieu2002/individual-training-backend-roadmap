package domain

type UserID string
type OrderID string
type PromoCode string
type TransactionID string

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusPaid    OrderStatus = "paid"
)

type PaymentMethodType string

const (
	PaymentMethodCreditCard   PaymentMethodType = "credit_card"
	PaymentMethodWallet       PaymentMethodType = "wallet"
	PaymentMethodBankTransfer PaymentMethodType = "bank_transfer"
	PaymentMethodCash         PaymentMethodType = "cash"
	PaymentMethodEwallet      PaymentMethodType = "ewallet"
	PaymentMethodCOD          PaymentMethodType = "cod"
)

type CreditCard struct {
	Number string
	CVV    string
	Expiry string
}

type User struct {
	ID    UserID
	Email string
	Phone string
}

type Order struct {
	ID           OrderID
	Status       OrderStatus
	ShippingCost *Money
}

type PaymentRequest struct {
	UserID     UserID
	OrderID    OrderID
	Amount     *Money
	Method     PaymentMethodType
	PromoCode  PromoCode
	CreditCard *CreditCard
}

type Transaction struct {
	ID             TransactionID
	UserID         UserID
	OrderID        OrderID
	Amount         *Money
	Method         PaymentMethodType
	Status         string
	VirtualAccount string
}
