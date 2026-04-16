package payment

import (
	"context"
	"fmt"
	"time"

	"individual-training-backend-roadmap/week1/solution/excercise1/domain"
)

type PaymentMethod interface {
	Process(ctx context.Context, req PaymentMethodRequest) (domain.TransactionID, error)
}

type PaymentMethodRequest struct {
	UserID  domain.UserID
	OrderID domain.OrderID
	Amount  *domain.Money
	Card    *domain.CreditCard
}

// CreditCardPayment charges via an external payment gateway.
type CreditCardPayment struct {
	gateway PaymentGateway
	txnRepo TransactionRepository
}

func (p *CreditCardPayment) Process(ctx context.Context, req PaymentMethodRequest) (domain.TransactionID, error) {
	txnID, err := p.gateway.Charge(ctx, req.Card, req.Amount)
	if err != nil {
		return "", fmt.Errorf("credit card charge failed: %w", err)
	}

	txn := &domain.Transaction{
		ID:      txnID,
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Method:  domain.PaymentMethodCreditCard,
		Status:  "success",
	}
	if err := p.txnRepo.Save(ctx, txn); err != nil {
		return "", fmt.Errorf("failed to save transaction: %w", err)
	}

	return txnID, nil
}

// WalletPayment deducts from the user's internal wallet.
type WalletPayment struct {
	walletService WalletService
	txnRepo       TransactionRepository
}

func (p *WalletPayment) Process(ctx context.Context, req PaymentMethodRequest) (domain.TransactionID, error) {
	balance, err := p.walletService.GetBalance(ctx, req.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to get wallet balance: %w", err)
	}

	if balance.LessThan(req.Amount) {
		return "", domain.ErrInsufficientBalance
	}

	if err := p.walletService.Deduct(ctx, req.UserID, req.Amount); err != nil {
		return "", fmt.Errorf("failed to deduct wallet balance: %w", err)
	}

	txnID := domain.TransactionID(fmt.Sprintf("TXN_%d", time.Now().Unix()))
	txn := &domain.Transaction{
		ID:      txnID,
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Method:  domain.PaymentMethodWallet,
		Status:  "success",
	}
	if err := p.txnRepo.Save(ctx, txn); err != nil {
		return "", fmt.Errorf("failed to save transaction: %w", err)
	}

	return txnID, nil
}

// BankTransferPayment creates a pending transaction with a virtual account.
type BankTransferPayment struct {
	txnRepo TransactionRepository
}

func (p *BankTransferPayment) Process(ctx context.Context, req PaymentMethodRequest) (domain.TransactionID, error) {
	virtualAccount := fmt.Sprintf("VA_%s_%d", req.UserID, time.Now().Unix())
	txnID := domain.TransactionID(fmt.Sprintf("TXN_%d", time.Now().Unix()))

	txn := &domain.Transaction{
		ID:             txnID,
		UserID:         req.UserID,
		OrderID:        req.OrderID,
		Amount:         req.Amount,
		Method:         domain.PaymentMethodBankTransfer,
		Status:         "pending",
		VirtualAccount: virtualAccount,
	}
	if err := p.txnRepo.Save(ctx, txn); err != nil {
		return "", fmt.Errorf("failed to save transaction: %w", err)
	}

	return txnID, nil
}

// CashPayment creates a pending transaction for cash-on-delivery.
type CashPayment struct {
	txnRepo TransactionRepository
}

func (p *CashPayment) Process(ctx context.Context, req PaymentMethodRequest) (domain.TransactionID, error) {
	txnID := domain.TransactionID(fmt.Sprintf("TXN_%d", time.Now().Unix()))

	txn := &domain.Transaction{
		ID:      txnID,
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Method:  domain.PaymentMethodCash,
		Status:  "pending",
	}
	if err := p.txnRepo.Save(ctx, txn); err != nil {
		return "", fmt.Errorf("failed to save transaction: %w", err)
	}

	return txnID, nil
}

// PaymentMethodFactory creates the appropriate PaymentMethod for a given type.
type PaymentMethodFactory struct {
	gateway       PaymentGateway
	walletService WalletService
	txnRepo       TransactionRepository
}

func NewPaymentMethodFactory(
	gateway PaymentGateway,
	walletService WalletService,
	txnRepo TransactionRepository,
) *PaymentMethodFactory {
	return &PaymentMethodFactory{
		gateway:       gateway,
		walletService: walletService,
		txnRepo:       txnRepo,
	}
}

func (f *PaymentMethodFactory) Create(method domain.PaymentMethodType) (PaymentMethod, error) {
	switch method {
	case domain.PaymentMethodCreditCard:
		return &CreditCardPayment{gateway: f.gateway, txnRepo: f.txnRepo}, nil
	case domain.PaymentMethodWallet:
		return &WalletPayment{walletService: f.walletService, txnRepo: f.txnRepo}, nil
	case domain.PaymentMethodBankTransfer:
		return &BankTransferPayment{txnRepo: f.txnRepo}, nil
	case domain.PaymentMethodCash:
		return &CashPayment{txnRepo: f.txnRepo}, nil
	default:
		return nil, domain.ErrUnsupportedPaymentMethod
	}
}
