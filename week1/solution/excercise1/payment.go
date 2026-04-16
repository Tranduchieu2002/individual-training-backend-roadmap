package payment

import (
	"context"
	"fmt"
	"log"

	"individual-training-backend-roadmap/week1/solution/excercise1/domain"
)

// PaymentProcessor orchestrates the payment flow with injected dependencies.
type PaymentProcessor struct {
	validator     *PaymentValidator
	taxCalculator *TaxCalculator
	userService   UserService
	orderService  OrderService
	promoService  PromoCodeService
	notifyService NotificationService
	methodFactory *PaymentMethodFactory
	cacheService  CacheService
}

func NewPaymentProcessor(
	userService UserService,
	orderService OrderService,
	promoService PromoCodeService,
	notifyService NotificationService,
	methodFactory *PaymentMethodFactory,
	cacheService CacheService,
) *PaymentProcessor {
	return &PaymentProcessor{
		validator:     NewPaymentValidator(),
		taxCalculator: NewTaxCalculator(),
		userService:   userService,
		orderService:  orderService,
		promoService:  promoService,
		notifyService: notifyService,
		methodFactory: methodFactory,
		cacheService:  cacheService,
	}
}

func (p *PaymentProcessor) ProcessPayment(ctx context.Context, req domain.PaymentRequest) (domain.TransactionID, error) {
	if err := p.validator.ValidateAmount(req.Amount); err != nil {
		return "", err
	}

	if req.Method == domain.PaymentMethodCreditCard {
		if err := p.validator.ValidateCreditCard(req.CreditCard); err != nil {
			return "", err
		}
	}

	user, err := p.verifyUser(ctx, req.UserID)
	if err != nil {
		return "", err
	}

	order, err := p.verifyOrder(ctx, req.OrderID)
	if err != nil {
		return "", err
	}

	finalAmount, err := p.calculateFinalAmount(ctx, req)
	if err != nil {
		return "", err
	}

	txnID, err := p.executePayment(ctx, req, finalAmount)
	if err != nil {
		return "", err
	}

	if err := p.orderService.UpdateOrderStatus(ctx, req.OrderID, domain.OrderStatusPaid, txnID); err != nil {
		return "", fmt.Errorf("failed to update order status: %w", err)
	}

	p.handlePostPayment(ctx, user, order, txnID, finalAmount)

	return txnID, nil
}

func (p *PaymentProcessor) verifyUser(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	user, err := p.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	banned, err := p.userService.IsUserBanned(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user ban status: %w", err)
	}
	if banned {
		return nil, domain.ErrUserBanned
	}

	return user, nil
}

func (p *PaymentProcessor) verifyOrder(ctx context.Context, orderID domain.OrderID) (*domain.Order, error) {
	order, err := p.orderService.GetOrder(ctx, orderID)
	if err != nil {
		return nil, domain.ErrOrderNotFound
	}

	if order.Status != domain.OrderStatusPending {
		return nil, domain.ErrOrderAlreadyProcessed
	}

	return order, nil
}

func (p *PaymentProcessor) calculateFinalAmount(ctx context.Context, req domain.PaymentRequest) (*domain.Money, error) {
	discount, err := p.applyPromoCode(ctx, req.PromoCode, req.Amount, req.OrderID)
	if err != nil {
		return nil, err
	}

	afterDiscount := req.Amount.Subtract(discount)
	tax := p.taxCalculator.Calculate(afterDiscount)
	finalAmount := afterDiscount.Add(tax)

	return finalAmount, nil
}

func (p *PaymentProcessor) applyPromoCode(ctx context.Context, code domain.PromoCode, amount *domain.Money, orderID domain.OrderID) (*domain.Money, error) {
	if code == "" {
		return domain.ZeroMoney(amount.Currency), nil
	}

	if err := p.promoService.Validate(ctx, code); err != nil {
		return domain.ZeroMoney(amount.Currency), nil
	}

	discount, err := p.promoService.CalculateDiscount(ctx, code, amount, orderID)
	if err != nil {
		return domain.ZeroMoney(amount.Currency), nil
	}

	if err := p.promoService.MarkAsUsed(ctx, code); err != nil {
		return nil, fmt.Errorf("failed to mark promo code as used: %w", err)
	}

	return discount, nil
}

func (p *PaymentProcessor) executePayment(ctx context.Context, req domain.PaymentRequest, amount *domain.Money) (domain.TransactionID, error) {
	method, err := p.methodFactory.Create(req.Method)
	if err != nil {
		return "", err
	}

	return method.Process(ctx, PaymentMethodRequest{
		UserID:  req.UserID,
		OrderID: req.OrderID,
		Amount:  amount,
		Card:    req.CreditCard,
	})
}

func (p *PaymentProcessor) handlePostPayment(ctx context.Context, user *domain.User, order *domain.Order, txnID domain.TransactionID, amount *domain.Money) {
	go func() {
		if err := p.notifyService.SendPaymentConfirmation(ctx, user, order, txnID); err != nil {
			log.Printf("failed to send payment confirmation: %v", err)
		}
	}()

	loyaltyPoints := int(amount.Amount.IntPart()) / LoyaltyPointsPerUnit
	if err := p.userService.UpdateLoyaltyPoints(ctx, user.ID, loyaltyPoints); err != nil {
		log.Printf("failed to update loyalty points: %v", err)
	}

	if err := p.userService.UpdateStats(ctx, user.ID, amount); err != nil {
		log.Printf("failed to update user stats: %v", err)
	}

	if err := p.cacheService.Invalidate(ctx, "user_"+string(user.ID), "order_"+string(order.ID)); err != nil {
		log.Printf("failed to invalidate cache: %v", err)
	}
}
