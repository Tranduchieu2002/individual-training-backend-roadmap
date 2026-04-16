package payment

import (
	"context"

	"individual-training-backend-roadmap/week1/solution/excercise1/domain"
)

type UserService interface {
	GetUser(ctx context.Context, id domain.UserID) (*domain.User, error)
	IsUserBanned(ctx context.Context, id domain.UserID) (bool, error)
	UpdateLoyaltyPoints(ctx context.Context, id domain.UserID, points int) error
	UpdateStats(ctx context.Context, id domain.UserID, amount *domain.Money) error
}

type OrderService interface {
	GetOrder(ctx context.Context, id domain.OrderID) (*domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id domain.OrderID, status domain.OrderStatus, txnID domain.TransactionID) error
}

type PromoCodeService interface {
	Validate(ctx context.Context, code domain.PromoCode) error
	CalculateDiscount(ctx context.Context, code domain.PromoCode, amount *domain.Money, orderID domain.OrderID) (*domain.Money, error)
	MarkAsUsed(ctx context.Context, code domain.PromoCode) error
}

type NotificationService interface {
	SendPaymentConfirmation(ctx context.Context, user *domain.User, order *domain.Order, txnID domain.TransactionID) error
}

type PaymentGateway interface {
	Charge(ctx context.Context, card *domain.CreditCard, amount *domain.Money) (domain.TransactionID, error)
}

type WalletService interface {
	GetBalance(ctx context.Context, userID domain.UserID) (*domain.Money, error)
	Deduct(ctx context.Context, userID domain.UserID, amount *domain.Money) error
}

type TransactionRepository interface {
	Save(ctx context.Context, txn *domain.Transaction) error
}

type CacheService interface {
	Invalidate(ctx context.Context, keys ...string) error
}
