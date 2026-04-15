# 🔍 CODE SMELL REFACTORING GUIDE

## Phase 1 Clean Code Exercises - Solutions & Learning Path

---

## 📚 EXERCISE OVERVIEW

### Exercise #1: Payment Processor (`smell_code_exercise_1_payment.go`)

**Focus:** Long methods, God class, Poor naming, Deep nesting
**Difficulty:** ⭐⭐⭐⭐ (Hard)
**Time:** 4-6 hours

### Exercise #2: Promo Code Manager (`smell_code_exercise_2_promo.go`)

**Focus:** Duplicate code, Magic numbers, Type codes
**Difficulty:** ⭐⭐⭐ (Medium)
**Time:** 3-4 hours

### Exercise #3: Notification Service (`smell_code_exercise_3_notification.go`)

**Focus:** OOP design issues, Template method pattern needed
**Difficulty:** ⭐⭐⭐⭐ (Hard)
**Time:** 4-5 hours

---

## 🎯 EXERCISE #1: PAYMENT PROCESSOR

### Code Smells Checklist (Find all 15+):

#### Method-Level Smells:

- [ ] **Long Method**: `ProcessPayment()` is 200+ lines (should be < 20 lines)
- [ ] **Deep Nesting**: 5+ levels of if-else (should be max 2-3)
- [ ] **Magic Numbers**: `10000000`, `500000`, `0.1`, `0.08`, etc.
- [ ] **Poor Variable Naming**: `x`, `y`, `z`, `o`, `os`, `p`, `pt`, `pv`, `uid`, `oid`
- [ ] **Duplicate Database Queries**: Query same table multiple times

#### Class-Level Smells:

- [ ] **God Class**: ProcessPayment does 10+ different things
- [ ] **Feature Envy**: Accesses database directly everywhere
- [ ] **Global State**: Uses global `db`, `cache`, `config` variables
- [ ] **No Error Handling**: Ignores errors from database and HTTP calls
- [ ] **Mixed Abstraction Levels**: High-level business logic mixed with low-level DB calls

#### Design Smells:

- [ ] **No Separation of Concerns**: Validation, calculation, persistence all mixed
- [ ] **Hard to Test**: Can't test without real database
- [ ] **No Dependency Injection**: Tightly coupled to global database
- [ ] **Primitive Obsession**: Using `float64` for money, `string` for IDs
- [ ] **Switch Statement**: Payment method handled with if-else chain

### Refactoring Strategy:

#### Step 1: Extract Value Objects

```go
// Create domain types
type Money struct {
    Amount   decimal.Decimal
    Currency string
}

type UserID string
type OrderID string
type PromoCode string
type TransactionID string

type CreditCard struct {
    Number string
    CVV    string
    Expiry string
}
```

#### Step 2: Extract Validation Logic

```go
type PaymentValidator struct{}

func (v *PaymentValidator) ValidateAmount(amount Money) error {
    if amount.Amount.LessThanOrEqual(decimal.Zero) {
        return ErrInvalidAmount
    }
    if amount.Amount.GreaterThan(decimal.NewFromFloat(10_000_000)) {
        return ErrAmountTooLarge
    }
    return nil
}

func (v *PaymentValidator) ValidateCreditCard(card CreditCard) error {
    if len(card.Number) < 13 || len(card.Number) > 19 {
        return ErrInvalidCardNumber
    }
    if len(card.CVV) != 3 && len(card.CVV) != 4 {
        return ErrInvalidCVV
    }
    // ... more validation
    return nil
}
```

#### Step 3: Extract Services

```go
// Separate concerns into different services

type UserService interface {
    GetUser(ctx context.Context, id UserID) (*User, error)
    IsUserBanned(ctx context.Context, id UserID) (bool, error)
}

type OrderService interface {
    GetOrder(ctx context.Context, id OrderID) (*Order, error)
    UpdateOrderStatus(ctx context.Context, id OrderID, status OrderStatus) error
}

type PromoCodeService interface {
    Validate(ctx context.Context, code PromoCode, userID UserID, amount Money) error
    CalculateDiscount(ctx context.Context, code PromoCode, amount Money) (Money, error)
    MarkAsUsed(ctx context.Context, code PromoCode) error
}

type WalletService interface {
    GetBalance(ctx context.Context, userID UserID) (Money, error)
    Deduct(ctx context.Context, userID UserID, amount Money) error
}
```

#### Step 4: Strategy Pattern for Payment Methods

```go
type PaymentMethod interface {
    Process(ctx context.Context, amount Money, details interface{}) (TransactionID, error)
}

type CreditCardPayment struct {
    gateway PaymentGateway
}

func (p *CreditCardPayment) Process(ctx context.Context, amount Money, details interface{}) (TransactionID, error) {
    card := details.(CreditCard)
    // Validate card
    // Call gateway
    // Return transaction ID
}

type WalletPayment struct {
    walletService WalletService
}

func (p *WalletPayment) Process(ctx context.Context, amount Money, details interface{}) (TransactionID, error) {
    // Deduct from wallet
    // Return transaction ID
}

// Factory pattern
func GetPaymentMethod(methodType string) PaymentMethod {
    switch methodType {
    case "credit_card":
        return &CreditCardPayment{}
    case "wallet":
        return &WalletPayment{}
    // ...
    }
}
```

#### Step 5: Main Orchestrator (Clean!)

```go
type PaymentProcessor struct {
    validator         *PaymentValidator
    userService       UserService
    orderService      OrderService
    promoCodeService  PromoCodeService
    taxCalculator     *TaxCalculator
    notificationService NotificationService
}

func (p *PaymentProcessor) ProcessPayment(ctx context.Context, req PaymentRequest) (TransactionID, error) {
    // 1. Validate request
    if err := p.validator.ValidateAmount(req.Amount); err != nil {
        return "", err
    }

    // 2. Verify user and order
    user, err := p.userService.GetUser(ctx, req.UserID)
    if err != nil {
        return "", err
    }

    if banned, _ := p.userService.IsUserBanned(ctx, req.UserID); banned {
        return "", ErrUserBanned
    }

    order, err := p.orderService.GetOrder(ctx, req.OrderID)
    if err != nil {
        return "", err
    }

    // 3. Apply promo code
    discount, err := p.applyPromoCode(ctx, req.PromoCode, req.UserID, req.Amount)
    if err != nil {
        return "", err
    }

    // 4. Calculate final amount
    finalAmount := p.calculateFinalAmount(req.Amount, discount, req.Currency)

    // 5. Process payment
    paymentMethod := GetPaymentMethod(req.Method)
    txnID, err := paymentMethod.Process(ctx, finalAmount, req.PaymentDetails)
    if err != nil {
        return "", err
    }

    // 6. Update order
    if err := p.orderService.UpdateOrderStatus(ctx, req.OrderID, StatusPaid); err != nil {
        return "", err
    }

    // 7. Send notifications (async)
    go p.notificationService.SendPaymentConfirmation(user, order, txnID)

    return txnID, nil
}

// Each small method does ONE thing
func (p *PaymentProcessor) applyPromoCode(ctx context.Context, code PromoCode, userID UserID, amount Money) (Money, error) {
    if code == "" {
        return Money{}, nil
    }

    if err := p.promoCodeService.Validate(ctx, code, userID, amount); err != nil {
        return Money{}, err
    }

    discount, err := p.promoCodeService.CalculateDiscount(ctx, code, amount)
    if err != nil {
        return Money{}, err
    }

    if err := p.promoCodeService.MarkAsUsed(ctx, code); err != nil {
        return Money{}, err
    }

    return discount, nil
}

func (p *PaymentProcessor) calculateFinalAmount(base Money, discount Money, currency string) Money {
    afterDiscount := base.Amount.Sub(discount.Amount)
    tax := p.taxCalculator.Calculate(afterDiscount, currency)
    final := afterDiscount.Add(tax)

    return Money{Amount: final, Currency: currency}
}
```

### Expected Improvements:

```
Before Refactoring:
✗ ProcessPayment: 200 lines
✗ Cyclomatic Complexity: 25
✗ Test Coverage: 0%
✗ Cannot test without database

After Refactoring:
✓ ProcessPayment: 30 lines
✓ Cyclomatic Complexity: 4
✓ Test Coverage: 90%+
✓ All dependencies injected (testable with mocks)
✓ Each method has single responsibility
✓ Clear separation of concerns
```

---

## 🎯 EXERCISE #2: PROMO CODE MANAGER

### Code Smells Checklist:

- [ ] **Duplicate Code**: Same validation logic repeated 3+ times
- [ ] **Magic Numbers**: Type codes (1, 2, 3, 4), hardcoded values
- [ ] **Long Parameter List**: CreatePromoCode has 9 parameters!
- [ ] **Boolean Parameters**: `newUserOnly` flag makes code hard to read
- [ ] **Switch Statements**: Type codes (should use polymorphism)
- [ ] **Feature Envy**: Querying database directly everywhere
- [ ] **Shotgun Surgery**: Change promo type → need to edit 5 places
- [ ] **Comment Pollution**: Useless comments like "// Get promo code from database"

### Refactoring Strategy:

#### Step 1: Replace Type Codes with Polymorphism

```go
// Instead of type codes (1, 2, 3, 4), use strategy pattern

type PromoCodeType interface {
    CalculateDiscount(orderAmount Money, shippingCost Money) Money
    Validate(order Order) error
}

type PercentageDiscount struct {
    Percentage  float64
    MaxDiscount Money
}

func (p *PercentageDiscount) CalculateDiscount(orderAmount Money, _ Money) Money {
    discount := orderAmount.Multiply(p.Percentage / 100)
    if discount.GreaterThan(p.MaxDiscount) {
        return p.MaxDiscount
    }
    return discount
}

type FixedDiscount struct {
    Amount Money
}

func (f *FixedDiscount) CalculateDiscount(orderAmount Money, _ Money) Money {
    if f.Amount.GreaterThan(orderAmount) {
        return orderAmount
    }
    return f.Amount
}

type FreeShipping struct{}

func (f *FreeShipping) CalculateDiscount(_ Money, shippingCost Money) Money {
    return shippingCost
}

type Cashback struct {
    Percentage float64
}

func (c *Cashback) CalculateDiscount(orderAmount Money, _ Money) Money {
    return Money{} // Cashback doesn't reduce current order
}
```

#### Step 2: Introduce Parameter Object

```go
// Instead of 9 parameters, use a struct

type CreatePromoCodeRequest struct {
    Code          string
    Type          PromoCodeType
    ExpiryDate    time.Time
    UsageLimit    int
    MinOrder      Money
    UserSpecific  *UserID  // nil means for all users
    NewUserOnly   bool
}

func (m *PromoCodeManager) CreatePromoCode(req CreatePromoCodeRequest) error {
    // Much cleaner!
}
```

#### Step 3: Extract Constants

```go
const (
    MinCodeLength = 4
    MaxCodeLength = 20
    MaxDiscount   = 500_000
)
```

#### Step 4: Remove Duplicate Validation

```go
type PromoCodeValidator struct{}

func (v *PromoCodeValidator) ValidateCode(code string) error {
    if len(code) < MinCodeLength || len(code) > MaxCodeLength {
        return ErrInvalidCodeLength
    }
    if !isValidCodeFormat(code) {
        return ErrInvalidCodeFormat
    }
    return nil
}

func (v *PromoCodeValidator) ValidateUsage(promo PromoCode, user User) error {
    if promo.IsExpired() {
        return ErrPromoExpired
    }
    if promo.IsUsageLimitReached() {
        return ErrUsageLimitReached
    }
    if !promo.IsValidForUser(user) {
        return ErrNotValidForUser
    }
    return nil
}
```

---

## 🎯 EXERCISE #3: NOTIFICATION SERVICE

### Code Smells Checklist:

- [ ] **Duplicate Code**: Same notification logic repeated for each event type
- [ ] **Template Method Pattern Needed**: Email/SMS/Push structure identical
- [ ] **Inappropriate Intimacy**: Service knows too much about User internals
- [ ] **Data Clumps**: (title, body, language) appear together
- [ ] **Divergent Change**: Adding new notification type requires changing everything
- [ ] **Lazy Class**: TemplateManager, MessageFormatter don't do enough
- [ ] **Middle Man**: MessageFormatter just delegates

### Refactoring Strategy:

#### Step 1: Template Method Pattern

```go
// Base notification interface
type Notification interface {
    Send(recipient Recipient, content NotificationContent) error
}

// Email implementation
type EmailNotification struct {
    client EmailClient
}

func (e *EmailNotification) Send(recipient Recipient, content NotificationContent) error {
    return e.client.Send(EmailMessage{
        To:      recipient.Email,
        Subject: content.Subject,
        Body:    content.Body,
    })
}

// SMS implementation
type SMSNotification struct {
    client SMSClient
}

func (s *SMSNotification) Send(recipient Recipient, content NotificationContent) error {
    return s.client.Send(SMSMessage{
        To:      recipient.Phone,
        Message: content.Body,
    })
}

// Push implementation
type PushNotification struct {
    client PushClient
}

func (p *PushNotification) Send(recipient Recipient, content NotificationContent) error {
    return p.client.Send(PushMessage{
        Token: recipient.PushToken,
        Title: content.Subject,
        Body:  content.Body,
    })
}
```

#### Step 2: Strategy Pattern for Content

```go
type NotificationTemplate interface {
    GenerateContent(lang string, data interface{}) NotificationContent
}

type OrderConfirmationTemplate struct{}

func (t *OrderConfirmationTemplate) GenerateContent(lang string, data interface{}) NotificationContent {
    orderData := data.(OrderData)

    switch lang {
    case "en":
        return NotificationContent{
            Subject: "Order Confirmation",
            Body:    fmt.Sprintf("Your order %s (Amount: $%.2f) has been confirmed.", orderData.OrderID, orderData.Amount),
        }
    case "vi":
        return NotificationContent{
            Subject: "Xác nhận đơn hàng",
            Body:    fmt.Sprintf("Đơn hàng %s (Số tiền: %.2f VND) đã được xác nhận.", orderData.OrderID, orderData.Amount),
        }
    default:
        return NotificationContent{}
    }
}

type OrderShippedTemplate struct{}
type PaymentFailedTemplate struct{}
// ... more templates
```

#### Step 3: Orchestrator (Clean!)

```go
type NotificationService struct {
    channels  []Notification
    templates map[string]NotificationTemplate
}

func (s *NotificationService) Notify(event NotificationEvent, recipient Recipient) error {
    // 1. Get template
    template := s.templates[event.Type]

    // 2. Generate content
    content := template.GenerateContent(recipient.Language, event.Data)

    // 3. Send via all enabled channels
    for _, channel := range s.getEnabledChannels(recipient) {
        if err := channel.Send(recipient, content); err != nil {
            // Log error but continue with other channels
            log.Error("Failed to send notification", "channel", channel, "error", err)
        }
    }

    return nil
}

func (s *NotificationService) getEnabledChannels(recipient Recipient) []Notification {
    var enabled []Notification

    for _, channel := range s.channels {
        if recipient.IsChannelEnabled(channel.Type()) {
            enabled = append(enabled, channel)
        }
    }

    return enabled
}
```

---

## 📊 SELF-ASSESSMENT RUBRIC

### After refactoring, check if your code meets these criteria:

#### Code Quality (40%)

- [ ] No method > 20 lines
- [ ] No class > 200 lines
- [ ] Cyclomatic complexity < 10 per method
- [ ] No duplicate code (DRY)
- [ ] Meaningful variable names (no abbreviations)
- [ ] No magic numbers (use constants)

#### Design (30%)

- [ ] Single Responsibility Principle (each class does ONE thing)
- [ ] Open/Closed Principle (extend without modifying)
- [ ] Dependency Inversion (depend on abstractions, not concretions)
- [ ] No global state
- [ ] Proper separation of concerns

#### Testability (20%)

- [ ] All dependencies injected
- [ ] No direct database calls in business logic
- [ ] Pure functions where possible
- [ ] Can mock external services
- [ ] Test coverage > 80%

#### Readability (10%)

- [ ] Code reads like well-written prose
- [ ] Consistent naming conventions
- [ ] Comments explain "why", not "what"
- [ ] Clear error messages

---

## 🎓 LEARNING RESOURCES

### Books (Read in this order):

1. **Clean Code** (Robert C. Martin) - Chapters 1-3, 6-7, 10
2. **Refactoring** (Martin Fowler) - Catalog of refactoring patterns
3. **Head First Design Patterns** - Strategy, Template Method, Factory

### Practice:

1. Refactor Exercise #1 → Get code review → Iterate
2. Refactor Exercise #2 → Compare with solution → Learn
3. Refactor Exercise #3 → Present to mentor → Discuss trade-offs

### Code Review Focus Points:

When reviewing your refactored code, ask:

- Can a new engineer understand this without explanation?
- Can I write tests without mocking the database?
- If requirements change (add new payment method), how many files do I touch?
- Is each method doing exactly ONE thing?

---

## ✅ COMPLETION CRITERIA

You're ready to move to Phase 2 when:

- [ ] All 3 exercises refactored to meet rubric criteria
- [ ] Unit tests written with 80%+ coverage
- [ ] Code review from mentor approved
- [ ] Can explain each design decision
- [ ] Can identify code smells in unfamiliar code within 5 minutes

**Time investment:** ~15-20 hours total
**Expected outcome:** Code quality improves 10x, ready for production-level work

---

Good luck! Remember: **Refactoring is not about being clever, it's about being kind to the next person who reads your code (which is often future you).** 🚀
