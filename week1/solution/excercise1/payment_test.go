package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"individual-training-backend-roadmap/week1/solution/excercise1/domain"
)

// --- Mock implementations ---

type mockUserService struct {
	user      *domain.User
	banned    bool
	getUserErr error
	banErr    error
}

func (m *mockUserService) GetUser(_ context.Context, _ domain.UserID) (*domain.User, error) {
	if m.getUserErr != nil {
		return nil, m.getUserErr
	}
	return m.user, nil
}

func (m *mockUserService) IsUserBanned(_ context.Context, _ domain.UserID) (bool, error) {
	return m.banned, m.banErr
}

func (m *mockUserService) UpdateLoyaltyPoints(_ context.Context, _ domain.UserID, _ int) error {
	return nil
}

func (m *mockUserService) UpdateStats(_ context.Context, _ domain.UserID, _ *domain.Money) error {
	return nil
}

type mockOrderService struct {
	order    *domain.Order
	getErr   error
	updateErr error
}

func (m *mockOrderService) GetOrder(_ context.Context, _ domain.OrderID) (*domain.Order, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.order, nil
}

func (m *mockOrderService) UpdateOrderStatus(_ context.Context, _ domain.OrderID, _ domain.OrderStatus, _ domain.TransactionID) error {
	return m.updateErr
}

type mockPromoCodeService struct {
	validateErr error
	discount    *domain.Money
	discountErr error
	markErr     error
}

func (m *mockPromoCodeService) Validate(_ context.Context, _ domain.PromoCode) error {
	return m.validateErr
}

func (m *mockPromoCodeService) CalculateDiscount(_ context.Context, _ domain.PromoCode, _ *domain.Money, _ domain.OrderID) (*domain.Money, error) {
	if m.discountErr != nil {
		return nil, m.discountErr
	}
	return m.discount, nil
}

func (m *mockPromoCodeService) MarkAsUsed(_ context.Context, _ domain.PromoCode) error {
	return m.markErr
}

type mockNotificationService struct{}

func (m *mockNotificationService) SendPaymentConfirmation(_ context.Context, _ *domain.User, _ *domain.Order, _ domain.TransactionID) error {
	return nil
}

type mockPaymentGateway struct {
	txnID domain.TransactionID
	err   error
}

func (m *mockPaymentGateway) Charge(_ context.Context, _ *domain.CreditCard, _ *domain.Money) (domain.TransactionID, error) {
	return m.txnID, m.err
}

type mockWalletService struct {
	balance   *domain.Money
	getErr    error
	deductErr error
}

func (m *mockWalletService) GetBalance(_ context.Context, _ domain.UserID) (*domain.Money, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.balance, nil
}

func (m *mockWalletService) Deduct(_ context.Context, _ domain.UserID, _ *domain.Money) error {
	return m.deductErr
}

type mockTransactionRepo struct {
	err error
}

func (m *mockTransactionRepo) Save(_ context.Context, _ *domain.Transaction) error {
	return m.err
}

type mockCacheService struct{}

func (m *mockCacheService) Invalidate(_ context.Context, _ ...string) error {
	return nil
}

// --- Helpers ---

func newTestProcessor(opts ...func(*testDeps)) *PaymentProcessor {
	deps := defaultTestDeps()
	for _, opt := range opts {
		opt(deps)
	}

	factory := NewPaymentMethodFactory(deps.gateway, deps.wallet, deps.txnRepo)

	return NewPaymentProcessor(
		deps.userSvc,
		deps.orderSvc,
		deps.promoSvc,
		deps.notifySvc,
		factory,
		deps.cacheSvc,
	)
}

type testDeps struct {
	userSvc   UserService
	orderSvc  OrderService
	promoSvc  PromoCodeService
	notifySvc NotificationService
	gateway   PaymentGateway
	wallet    WalletService
	txnRepo   TransactionRepository
	cacheSvc  CacheService
}

func defaultTestDeps() *testDeps {
	return &testDeps{
		userSvc: &mockUserService{
			user: &domain.User{ID: "user-1", Email: "test@example.com", Phone: "0123456789"},
		},
		orderSvc: &mockOrderService{
			order: &domain.Order{ID: "order-1", Status: domain.OrderStatusPending},
		},
		promoSvc:  &mockPromoCodeService{},
		notifySvc: &mockNotificationService{},
		gateway: &mockPaymentGateway{
			txnID: "TXN_GATEWAY_001",
		},
		wallet: &mockWalletService{
			balance: domain.NewMoney(decimal.NewFromInt(1_000_000), domain.CurrencyUSD),
		},
		txnRepo:  &mockTransactionRepo{},
		cacheSvc: &mockCacheService{},
	}
}

func usdMoney(amount int64) *domain.Money {
	return domain.NewMoney(decimal.NewFromInt(amount), domain.CurrencyUSD)
}

func validCreditCard() *domain.CreditCard {
	return &domain.CreditCard{
		Number: "4111111111111111",
		CVV:    "123",
		Expiry: "12/26",
	}
}

// --- Tests ---

func TestProcessPayment_CreditCard_Success(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:     "user-1",
		OrderID:    "order-1",
		Amount:     usdMoney(100_000),
		Method:     domain.PaymentMethodCreditCard,
		CreditCard: validCreditCard(),
	}

	txnID, err := processor.ProcessPayment(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txnID == "" {
		t.Fatal("expected transaction ID, got empty string")
	}
}

func TestProcessPayment_Wallet_Success(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(50_000),
		Method:  domain.PaymentMethodWallet,
	}

	txnID, err := processor.ProcessPayment(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txnID == "" {
		t.Fatal("expected transaction ID, got empty string")
	}
}

func TestProcessPayment_BankTransfer_Success(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(200_000),
		Method:  domain.PaymentMethodBankTransfer,
	}

	txnID, err := processor.ProcessPayment(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txnID == "" {
		t.Fatal("expected transaction ID, got empty string")
	}
}

func TestProcessPayment_Cash_Success(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(30_000),
		Method:  domain.PaymentMethodCash,
	}

	txnID, err := processor.ProcessPayment(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txnID == "" {
		t.Fatal("expected transaction ID, got empty string")
	}
}

func TestProcessPayment_InvalidAmount_Zero(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(0),
		Method:  domain.PaymentMethodCreditCard,
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestProcessPayment_InvalidAmount_Negative(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  domain.NewMoney(decimal.NewFromInt(-100), domain.CurrencyUSD),
		Method:  domain.PaymentMethodCreditCard,
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestProcessPayment_InvalidAmount_TooLarge(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:     "user-1",
		OrderID:    "order-1",
		Amount:     usdMoney(20_000_000),
		Method:     domain.PaymentMethodCreditCard,
		CreditCard: validCreditCard(),
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrAmountTooLarge) {
		t.Fatalf("expected ErrAmountTooLarge, got %v", err)
	}
}

func TestProcessPayment_UserNotFound(t *testing.T) {
	processor := newTestProcessor(func(d *testDeps) {
		d.userSvc = &mockUserService{getUserErr: errors.New("not found")}
	})
	req := domain.PaymentRequest{
		UserID:     "unknown",
		OrderID:    "order-1",
		Amount:     usdMoney(100_000),
		Method:     domain.PaymentMethodCreditCard,
		CreditCard: validCreditCard(),
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestProcessPayment_UserBanned(t *testing.T) {
	processor := newTestProcessor(func(d *testDeps) {
		d.userSvc = &mockUserService{
			user:   &domain.User{ID: "user-banned"},
			banned: true,
		}
	})
	req := domain.PaymentRequest{
		UserID:     "user-banned",
		OrderID:    "order-1",
		Amount:     usdMoney(100_000),
		Method:     domain.PaymentMethodCreditCard,
		CreditCard: validCreditCard(),
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrUserBanned) {
		t.Fatalf("expected ErrUserBanned, got %v", err)
	}
}

func TestProcessPayment_OrderNotFound(t *testing.T) {
	processor := newTestProcessor(func(d *testDeps) {
		d.orderSvc = &mockOrderService{getErr: errors.New("not found")}
	})
	req := domain.PaymentRequest{
		UserID:     "user-1",
		OrderID:    "unknown",
		Amount:     usdMoney(100_000),
		Method:     domain.PaymentMethodCreditCard,
		CreditCard: validCreditCard(),
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestProcessPayment_OrderAlreadyProcessed(t *testing.T) {
	processor := newTestProcessor(func(d *testDeps) {
		d.orderSvc = &mockOrderService{
			order: &domain.Order{ID: "order-1", Status: domain.OrderStatusPaid},
		}
	})
	req := domain.PaymentRequest{
		UserID:     "user-1",
		OrderID:    "order-1",
		Amount:     usdMoney(100_000),
		Method:     domain.PaymentMethodCreditCard,
		CreditCard: validCreditCard(),
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrOrderAlreadyProcessed) {
		t.Fatalf("expected ErrOrderAlreadyProcessed, got %v", err)
	}
}

func TestProcessPayment_InvalidCreditCard_ShortNumber(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(100_000),
		Method:  domain.PaymentMethodCreditCard,
		CreditCard: &domain.CreditCard{
			Number: "411",
			CVV:    "123",
			Expiry: "12/26",
		},
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrInvalidCardNumber) {
		t.Fatalf("expected ErrInvalidCardNumber, got %v", err)
	}
}

func TestProcessPayment_InvalidCVV(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(100_000),
		Method:  domain.PaymentMethodCreditCard,
		CreditCard: &domain.CreditCard{
			Number: "4111111111111111",
			CVV:    "12",
			Expiry: "12/26",
		},
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrInvalidCVV) {
		t.Fatalf("expected ErrInvalidCVV, got %v", err)
	}
}

func TestProcessPayment_InvalidExpiry(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(100_000),
		Method:  domain.PaymentMethodCreditCard,
		CreditCard: &domain.CreditCard{
			Number: "4111111111111111",
			CVV:    "123",
			Expiry: "invalid",
		},
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrInvalidExpiry) {
		t.Fatalf("expected ErrInvalidExpiry, got %v", err)
	}
}

func TestProcessPayment_WalletInsufficientBalance(t *testing.T) {
	processor := newTestProcessor(func(d *testDeps) {
		d.wallet = &mockWalletService{
			balance: usdMoney(100),
		}
	})
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(500_000),
		Method:  domain.PaymentMethodWallet,
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestProcessPayment_UnsupportedMethod(t *testing.T) {
	processor := newTestProcessor()
	req := domain.PaymentRequest{
		UserID:  "user-1",
		OrderID: "order-1",
		Amount:  usdMoney(100_000),
		Method:  "bitcoin",
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if !errors.Is(err, domain.ErrUnsupportedPaymentMethod) {
		t.Fatalf("expected ErrUnsupportedPaymentMethod, got %v", err)
	}
}

func TestProcessPayment_WithPromoCode(t *testing.T) {
	processor := newTestProcessor(func(d *testDeps) {
		d.promoSvc = &mockPromoCodeService{
			discount: usdMoney(10_000),
		}
	})
	req := domain.PaymentRequest{
		UserID:     "user-1",
		OrderID:    "order-1",
		Amount:     usdMoney(100_000),
		Method:     domain.PaymentMethodCreditCard,
		PromoCode:  "SAVE10",
		CreditCard: validCreditCard(),
	}

	txnID, err := processor.ProcessPayment(context.Background(), req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if txnID == "" {
		t.Fatal("expected transaction ID, got empty string")
	}
}

func TestProcessPayment_GatewayFailure(t *testing.T) {
	processor := newTestProcessor(func(d *testDeps) {
		d.gateway = &mockPaymentGateway{err: errors.New("gateway timeout")}
	})
	req := domain.PaymentRequest{
		UserID:     "user-1",
		OrderID:    "order-1",
		Amount:     usdMoney(100_000),
		Method:     domain.PaymentMethodCreditCard,
		CreditCard: validCreditCard(),
	}

	_, err := processor.ProcessPayment(context.Background(), req)

	if err == nil {
		t.Fatal("expected error from gateway failure, got nil")
	}
}

// --- Validator unit tests ---

func TestPaymentValidator_ValidateAmount(t *testing.T) {
	v := NewPaymentValidator()

	tests := []struct {
		name    string
		amount  *domain.Money
		wantErr error
	}{
		{"valid amount", usdMoney(50_000), nil},
		{"zero amount", usdMoney(0), domain.ErrInvalidAmount},
		{"negative amount", domain.NewMoney(decimal.NewFromInt(-1), domain.CurrencyUSD), domain.ErrInvalidAmount},
		{"nil amount", nil, domain.ErrInvalidAmount},
		{"max boundary", usdMoney(MaxPaymentAmount), nil},
		{"exceeds max", usdMoney(MaxPaymentAmount + 1), domain.ErrAmountTooLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateAmount(tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateAmount() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaymentValidator_ValidateCreditCard(t *testing.T) {
	v := NewPaymentValidator()

	tests := []struct {
		name    string
		card    *domain.CreditCard
		wantErr error
	}{
		{"valid card", validCreditCard(), nil},
		{"nil card", nil, domain.ErrInvalidCardNumber},
		{"short number", &domain.CreditCard{Number: "411", CVV: "123", Expiry: "12/26"}, domain.ErrInvalidCardNumber},
		{"long number", &domain.CreditCard{Number: "41111111111111111111", CVV: "123", Expiry: "12/26"}, domain.ErrInvalidCardNumber},
		{"invalid cvv", &domain.CreditCard{Number: "4111111111111111", CVV: "12", Expiry: "12/26"}, domain.ErrInvalidCVV},
		{"4-digit cvv", &domain.CreditCard{Number: "4111111111111111", CVV: "1234", Expiry: "12/26"}, nil},
		{"invalid expiry", &domain.CreditCard{Number: "4111111111111111", CVV: "123", Expiry: "1226"}, domain.ErrInvalidExpiry},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.ValidateCreditCard(tt.card)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateCreditCard() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// --- Tax calculator tests ---

func TestTaxCalculator_Calculate(t *testing.T) {
	calc := NewTaxCalculator()

	tests := []struct {
		name     string
		amount   *domain.Money
		expected string
	}{
		{"USD 10%", domain.NewMoney(decimal.NewFromInt(100_000), domain.CurrencyUSD), "10000"},
		{"VND 8%", domain.NewMoney(decimal.NewFromInt(100_000), domain.CurrencyVND), "8000"},
		{"SGD 7%", domain.NewMoney(decimal.NewFromInt(100_000), domain.CurrencySGD), "7000"},
		{"unknown currency defaults to 10%", domain.NewMoney(decimal.NewFromInt(100_000), "BTC"), "10000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.amount)
			if result.Amount.String() != tt.expected {
				t.Errorf("Calculate() = %s, want %s", result.Amount.String(), tt.expected)
			}
		})
	}
}
