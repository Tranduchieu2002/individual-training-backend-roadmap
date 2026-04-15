# 💩 CODE SMELL EXERCISES - PHASE 1 TRAINING

## Luyện tập Refactoring cho Fundamental Engineers

---

## 📦 BỘ BÀI TẬP GỒM:

1. **smell_code_exercise_1_payment.go** - Payment Processor (God Class Hell)
2. **smell_code_exercise_2_promo.go** - Promo Code Manager (Duplicate Code Nightmare)
3. **smell_code_exercise_3_notification.go** - Notification Service (OOP Design Mess)
4. **refactoring_guide.md** - Hướng dẫn giải & patterns

---

## 🎯 MỤC TIÊU HỌC TẬP

Sau khi hoàn thành 3 bài tập này, bạn sẽ:

✅ **Nhận diện được 20+ code smells phổ biến** chỉ trong 5 phút đọc code
✅ **Refactor code cũ thành clean code** một cách hệ thống
✅ **Áp dụng SOLID principles** vào thiết kế thực tế
✅ **Viết code dễ test** (testable code với dependency injection)
✅ **Tư duy design patterns** (Strategy, Template Method, Factory...)

---

## 🚀 CÁCH SỬ DỤNG

### **WEEK 1: Clean Code Bootcamp**

#### **Day 1-2: Exercise #1 - Payment Processor**

**Bước 1: Đọc và phân tích** (2 hours)

```bash
# Clone and open file
cd 'WEEK 1'/excercise
code smell_code_exercise_1_payment.go

# Nhiệm vụ:
1. Đọc toàn bộ code (đừng vội refactor!)
2. List ra TẤT CẢ code smells bạn tìm thấy
3. Vẽ diagram: luồng hiện tại của ProcessPayment
4. Identify: Method nào quá dài? Class nào làm quá nhiều việc?
```

**Checklist tìm smells:**

- [ ] Method > 50 lines? → Long Method
- [ ] Nested if > 3 levels? → Arrow Anti-pattern
- [ ] Biến tên `x`, `y`, `z`? → Poor Naming
- [ ] Số `10000000`, `0.08` xuất hiện? → Magic Numbers
- [ ] Query database nhiều lần cùng bảng? → Duplicate Queries
- [ ] Code lặp lại? → Copy-Paste Programming
- [ ] Hard to test? → Tight Coupling

**Bước 2: Plan refactoring** (2 hours)

```
Không được bắt đầu code ngay!
Trước tiên, viết plan:

# Refactoring Plan - Exercise #1

## Step 1: Extract Value Objects
- [ ] Create Money type (thay vì float64)
- [ ] Create UserID, OrderID types (thay vì string)
- [ ] Create CreditCard struct

## Step 2: Extract Services
- [ ] UserService interface
- [ ] OrderService interface
- [ ] PromoCodeService interface
- [ ] WalletService interface

## Step 3: Extract Validators
- [ ] PaymentValidator
- [ ] CreditCardValidator

## Step 4: Strategy Pattern for Payment Methods
- [ ] PaymentMethod interface
- [ ] CreditCardPayment implementation
- [ ] WalletPayment implementation
- [ ] BankTransferPayment implementation

## Step 5: Refactor Main Method
- [ ] ProcessPayment giảm xuống < 30 lines
- [ ] Mỗi step có method riêng
- [ ] Clear separation of concerns

Estimate: 6-8 hours
Risk: Breaking existing functionality → Need tests first!
```

**Bước 3: Write tests FIRST** (2 hours)

```go
// Viết test cho code CŨ trước khi refactor
// Đây gọi là "Safety Net"

func TestProcessPayment_Success(t *testing.T) {
    // Setup
    db := setupTestDB()
    defer db.Close()

    // Test case: Valid credit card payment
    req := createValidPaymentRequest()

    // Execute
    txnID, err := ProcessPayment(req)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, txnID)

    // Verify database state
    assertOrderIsPaid(t, db, req.OrderID)
}

func TestProcessPayment_InsufficientBalance(t *testing.T) {
    // Test wallet payment with insufficient balance
    // ...
}

// Write 10-15 test cases covering main scenarios
```

**Bước 4: Refactor từng bước nhỏ** (4-6 hours)

```bash
# Quy tắc vàng: Mỗi commit là một refactoring nhỏ

git commit -m "refactor: extract Money value object"
git commit -m "refactor: extract UserService interface"
git commit -m "refactor: apply strategy pattern for payment methods"
git commit -m "refactor: simplify ProcessPayment orchestration"

# Sau mỗi commit: Run tests!
go test ./... -v

# Nếu test fail → Revert và refactor lại cẩn thận hơn
```

**Bước 5: Code review** (1 hour)

```
Tự review code của mình:

Questions to ask:
□ Có method nào > 20 lines không?
□ Có class nào > 200 lines không?
□ ProcessPayment còn làm quá nhiều việc không?
□ Tất cả dependencies đã inject chưa?
□ Code có dễ test hơn không?
□ Tên biến/method có rõ ràng không?

Nếu answer là "No" cho bất kỳ câu nào → Keep refactoring!
```

---

#### **Day 3-4: Exercise #2 - Promo Code Manager**

**Focus:** Duplicate code, Magic numbers, Type codes

**Challenge:**

```
Bài tập này test khả năng DRY (Don't Repeat Yourself).

Key smells:
1. Logic validation giống nhau lặp lại 3 lần
2. Type codes (1, 2, 3, 4) thay vì polymorphism
3. CreatePromoCode có 9 parameters (quá nhiều!)
4. Copy-paste code ở CalculateDiscount và ApplyPromoCode

Your task:
→ Reduce duplication by 80%
→ Replace type codes với Strategy pattern
→ Introduce Parameter Object
→ Extract constants for all magic numbers
```

**Deliverable:**

- Refactored code với 0% duplication
- All type codes replaced với interfaces
- Test coverage > 85%

---

#### **Day 5: Exercise #3 - Notification Service**

**Focus:** OOP design, Template Method pattern

**Challenge:**

```
Đây là bài khó nhất vì test OOP design skills.

Current problems:
1. SendOrderConfirmation, SendOrderShipped, SendPaymentFailed
   → 90% duplicate code, chỉ khác content
2. Thêm notification type mới → Phải copy-paste toàn bộ
3. TemplateManager, MessageFormatter là "lazy classes"

Your task:
→ Apply Template Method pattern
→ Separate concerns: Content generation vs Delivery
→ Make adding new notification type trivial (< 10 lines)
```

**Deliverable:**

- Refactored với Template Method pattern
- Adding new notification type chỉ cần thêm 1 class
- Test coverage > 80%

---

## 📊 ĐÁNH GIÁ & TIÊU CHÍ ĐẠT

### Rubric for each exercise:

| Criteria         | Weight | Pass Criteria                                                                                     |
| ---------------- | ------ | ------------------------------------------------------------------------------------------------- |
| **Code Quality** | 40%    | ✓ No method > 20 lines<br>✓ No duplicate code<br>✓ Meaningful names<br>✓ No magic numbers         |
| **Design**       | 30%    | ✓ SOLID principles applied<br>✓ Proper separation of concerns<br>✓ Design patterns used correctly |
| **Testability**  | 20%    | ✓ Dependencies injected<br>✓ Pure functions<br>✓ Test coverage > 80%                              |
| **Readability**  | 10%    | ✓ Code self-documenting<br>✓ Consistent style<br>✓ Clear error messages                           |

**Passing score:** 80% overall, no category < 70%

---

## 🎓 TIPS & BEST PRACTICES

### ⚠️ Common Mistakes to Avoid:

**1. Refactor quá nhanh, bỏ qua tests**

```
❌ Wrong approach:
1. Đọc code xấu
2. Xóa hết, viết lại từ đầu
3. Hope it works

✅ Correct approach:
1. Đọc code xấu
2. Viết tests cho behavior hiện tại
3. Refactor từng bước nhỏ
4. Run tests sau mỗi step
5. Commit frequently
```

**2. Over-engineering**

```
❌ Bad: Tạo 20 abstractions cho vấn đề đơn giản
✅ Good: Refactor đủ để code clean, không làm phức tạp thêm
```

**3. Không commit đủ thường xuyên**

```
❌ Bad: 1 commit với 2000 lines changed
✅ Good: 10 commits, mỗi commit 1 refactoring nhỏ

Example good commit history:
- refactor: extract constants for magic numbers
- refactor: rename variables (x → userCount, y → status)
- refactor: extract validateUser method
- refactor: extract calculateDiscount method
- refactor: apply dependency injection for database
```

**4. Quên viết commit message rõ ràng**

```
❌ Bad commits:
- "fix"
- "update code"
- "wip"

✅ Good commits:
- "refactor: extract PromoCodeValidator class"
- "refactor: replace type codes with strategy pattern"
- "refactor: reduce ProcessPayment from 200 to 25 lines"
```

---

## 🛠️ SETUP & TOOLS

### Required:

```bash
# Install Go (nếu chưa có)
go version  # should be 1.21+

# Install testing tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install coverage tool
go install github.com/axw/gocov/gocov@latest
go install github.com/AlekSi/gocov-xml@latest
```

### Run tests:

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run linter
golangci-lint run
```

---

## 📚 RECOMMENDED READING ORDER

Đọc song song với làm bài tập:

### Day 1-2 (Exercise #1):

- Clean Code Chapter 3: Functions (về Long Method)
- Clean Code Chapter 10: Classes (về Single Responsibility)

### Day 3-4 (Exercise #2):

- Refactoring Chapter 3: Bad Smells in Code
- Refactoring Chapter 10: Simplifying Conditional Logic

### Day 5 (Exercise #3):

- Head First Design Patterns: Template Method Pattern
- Head First Design Patterns: Strategy Pattern

---

## 🎯 NEXT STEPS (After completing)

Khi đã pass cả 3 exercises:

1. **Code Review Session với Mentor**
   - Present refactoring decisions
   - Discuss trade-offs
   - Get feedback

2. **Write Blog Post**
   - Document learnings
   - Share before/after code
   - Explain design patterns used

3. **Ready for Phase 2!**
   - Move to System Design exercises
   - Apply clean code principles to larger systems

---

## 💬 SUPPORT & COMMUNITY

**Stuck? Need help?**

- Đọc `refactoring_guide.md` để xem solutions
- Post questions in `#training-help` Slack channel
- Book 1-on-1 with mentor
- Join Friday code review sessions

**Share your progress:**

- Post screenshots in `#training-wins`
- Celebrate small wins (mỗi PR merged!)
- Help others in `#code-review`

---

## 📈 EXPECTED TIMELINE

**Realistic timeline:**

- Exercise #1: 8-10 hours (phức tạp nhất)
- Exercise #2: 6-8 hours (medium)
- Exercise #3: 6-8 hours (cần hiểu design patterns)

**Total:** 20-26 hours (~ 1 week full-time hoặc 2 weeks part-time)

---

## ✨ SUCCESS STORIES

> "Tuần đầu mình nhìn Exercise #1 mà muốn khóc. Code 200 lines không biết bắt đầu từ đâu. Nhưng sau khi học cách refactor từng bước nhỏ, giờ mình tự tin refactor bất kỳ legacy code nào!"
> — **Minh Nguyen**, Junior Engineer → Mid-level trong 6 tháng

> "Ban đầu mình cứ nghĩ viết tests là tốn thời gian. Nhưng khi refactor Exercise #1, tests đã cứu mình 5 lần khỏi break production code. Bây giờ không dám code mà không có tests!"
> — **Thu Tran**, Backend Engineer

> "Exercise #3 khó nhất vì buộc mình phải suy nghĩ về design. Nhưng sau khi hiểu Template Method pattern, mình áp dụng được vào real project và giảm được 70% duplicate code!"
> — **Huy Le**, Full-Stack Engineer

---

**Ready to start? Pick Exercise #1 và bắt đầu hành trình refactoring! 🚀**

**Remember:**

> "Any fool can write code that a computer can understand. Good programmers write code that humans can understand."
> — Martin Fowler

Good luck! 💪
