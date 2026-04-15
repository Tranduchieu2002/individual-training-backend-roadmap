package payment

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CODE SMELL EXERCISE #1: Payment Processor
// This code has 15+ code smells. Find and fix them all!
//
// Smells present:
// - Long method (>100 lines)
// - God class (does too many things)
// - Magic numbers everywhere
// - Poor naming (a, b, c, x, y, z)
// - Deep nesting (if inside if inside if...)
// - No error handling
// - Mixed abstraction levels
// - Duplicate code
// - Hard to test (global state, database calls everywhere)
// - No separation of concerns
// - Feature envy
// - Primitive obsession

// Global variables (BAD!)
var db *sql.DB
var cache map[string]interface{}
var config map[string]string

// ProcessPayment handles payment processing
// TODO: This function is too long, need to refactor someday
func ProcessPayment(r *http.Request) (string, error) {
	// Get data from request
	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	// uid = user id, oid = order id, amt = amount, pm = payment method,
	// pc = promo code, cur = currency, cc = credit card
	uid := data["user_id"].(string)
	oid := data["order_id"].(string)
	amt := data["amount"].(float64)
	pm := data["payment_method"].(string)
	pc := data["promo_code"].(string)
	cur := data["currency"].(string)
	cc := data["credit_card"].(map[string]interface{})

	// Check if user exists
	var x int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", uid).Scan(&x)
	if x == 0 {
		return "", fmt.Errorf("user not found")
	}

	// Check if user is banned
	var y string
	db.QueryRow("SELECT status FROM users WHERE id = ?", uid).Scan(&y)
	if y == "banned" {
		return "", fmt.Errorf("user banned")
	}

	// Check if user has enough balance for wallet payment
	if pm == "wallet" {
		var z float64
		db.QueryRow("SELECT balance FROM wallets WHERE user_id = ?", uid).Scan(&z)
		if z < amt {
			return "", fmt.Errorf("insufficient balance")
		}
	}

	// Validate amount
	if amt <= 0 {
		return "", fmt.Errorf("invalid amount")
	}
	if amt > 10000000 {
		return "", fmt.Errorf("amount too large")
	}

	// Check order exists
	var o int
	db.QueryRow("SELECT COUNT(*) FROM orders WHERE id = ?", oid).Scan(&o)
	if o == 0 {
		return "", fmt.Errorf("order not found")
	}

	// Check order status
	var os string
	db.QueryRow("SELECT status FROM orders WHERE id = ?", oid).Scan(&os)
	if os != "pending" {
		return "", fmt.Errorf("order already processed")
	}

	// Process promo code
	var discount float64 = 0
	if pc != "" {
		// Check if promo code exists
		var p int
		db.QueryRow("SELECT COUNT(*) FROM promo_codes WHERE code = ?", pc).Scan(&p)
		if p > 0 {
			// Check if promo code is expired
			var exp string
			db.QueryRow("SELECT expires_at FROM promo_codes WHERE code = ?", pc).Scan(&exp)
			t, _ := time.Parse("2006-01-02", exp)
			if t.After(time.Now()) {
				// Check if promo code is used
				var used int
				db.QueryRow("SELECT used FROM promo_codes WHERE code = ?", pc).Scan(&used)
				if used == 0 {
					// Check promo type
					var pt string
					db.QueryRow("SELECT promo_type FROM promo_codes WHERE code = ?", pc).Scan(&pt)
					if pt == "percentage" {
						var pv float64
						db.QueryRow("SELECT value FROM promo_codes WHERE code = ?", pc).Scan(&pv)
						discount = amt * (pv / 100)
						// Max discount 500000
						if discount > 500000 {
							discount = 500000
						}
					} else if pt == "fixed" {
						var pv float64
						db.QueryRow("SELECT value FROM promo_codes WHERE code = ?", pc).Scan(&pv)
						discount = pv
						if discount > amt {
							discount = amt
						}
					} else if pt == "free_shipping" {
						// Get shipping cost
						var sc float64
						db.QueryRow("SELECT shipping_cost FROM orders WHERE id = ?", oid).Scan(&sc)
						discount = sc
					}

					// Mark promo as used
					db.Exec("UPDATE promo_codes SET used = 1 WHERE code = ?", pc)
				}
			}
		}
	}

	// Calculate final amount
	finalAmount := amt - discount

	// Add tax
	var tax float64
	if cur == "USD" {
		tax = finalAmount * 0.1
	} else if cur == "VND" {
		tax = finalAmount * 0.08
	} else if cur == "SGD" {
		tax = finalAmount * 0.07
	} else {
		tax = finalAmount * 0.1
	}
	finalAmount = finalAmount + tax

	// Process payment based on method
	var transactionId string
	if pm == "credit_card" {
		// Validate credit card
		cardNum := cc["number"].(string)
		cvv := cc["cvv"].(string)
		exp := cc["expiry"].(string)

		// Check card number length
		if len(cardNum) < 13 || len(cardNum) > 19 {
			return "", fmt.Errorf("invalid card")
		}

		// Check CVV
		if len(cvv) != 3 && len(cvv) != 4 {
			return "", fmt.Errorf("invalid cvv")
		}

		// Check expiry
		parts := strings.Split(exp, "/")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid expiry")
		}

		// Call external payment gateway
		resp, _ := http.Post("https://payment-gateway.com/charge", "application/json", nil)
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		transactionId = result["transaction_id"].(string)

		// Save to database
		db.Exec("INSERT INTO transactions (id, user_id, order_id, amount, method, status) VALUES (?, ?, ?, ?, ?, ?)",
			transactionId, uid, oid, finalAmount, pm, "success")

	} else if pm == "wallet" {
		// Deduct from wallet
		db.Exec("UPDATE wallets SET balance = balance - ? WHERE user_id = ?", finalAmount, uid)

		// Generate transaction id
		transactionId = fmt.Sprintf("TXN_%d", time.Now().Unix())

		// Save to database
		db.Exec("INSERT INTO transactions (id, user_id, order_id, amount, method, status) VALUES (?, ?, ?, ?, ?, ?)",
			transactionId, uid, oid, finalAmount, pm, "success")

	} else if pm == "bank_transfer" {
		// Generate virtual account
		va := fmt.Sprintf("VA_%s_%d", uid, time.Now().Unix())

		// Save pending transaction
		transactionId = fmt.Sprintf("TXN_%d", time.Now().Unix())
		db.Exec("INSERT INTO transactions (id, user_id, order_id, amount, method, status, virtual_account) VALUES (?, ?, ?, ?, ?, ?, ?)",
			transactionId, uid, oid, finalAmount, pm, "pending", va)

	} else if pm == "cash" {
		// Generate transaction id
		transactionId = fmt.Sprintf("TXN_%d", time.Now().Unix())

		// Save to database
		db.Exec("INSERT INTO transactions (id, user_id, order_id, amount, method, status) VALUES (?, ?, ?, ?, ?, ?)",
			transactionId, uid, oid, finalAmount, pm, "pending")
	}

	// Update order status
	db.Exec("UPDATE orders SET status = ?, transaction_id = ? WHERE id = ?", "paid", transactionId, oid)

	// Send notification
	if pm == "credit_card" || pm == "wallet" {
		// Get user email
		var email string
		db.QueryRow("SELECT email FROM users WHERE id = ?", uid).Scan(&email)

		// Send email
		http.Post("https://email-service.com/send", "application/json", nil)

		// Send SMS
		var phone string
		db.QueryRow("SELECT phone FROM users WHERE id = ?", uid).Scan(&phone)
		http.Post("https://sms-service.com/send", "application/json", nil)

		// Send push notification
		http.Post("https://push-service.com/send", "application/json", nil)
	}

	// Update user stats
	db.Exec("UPDATE user_stats SET total_spent = total_spent + ?, total_orders = total_orders + 1 WHERE user_id = ?", finalAmount, uid)

	// Add loyalty points
	points := int(finalAmount / 1000)
	db.Exec("UPDATE users SET loyalty_points = loyalty_points + ? WHERE id = ?", points, uid)

	// Clear cache
	delete(cache, "user_"+uid)
	delete(cache, "order_"+oid)

	return transactionId, nil
}

// ValidatePaymentMethod checks if payment method is valid
func ValidatePaymentMethod(pm string) bool {
	if pm == "credit_card" || pm == "wallet" || pm == "bank_transfer" || pm == "cash" || pm == "ewallet" || pm == "cod" {
		return true
	}
	return false
}

// GetTransactionStatus returns transaction status
func GetTransactionStatus(tid string) string {
	var s string
	db.QueryRow("SELECT status FROM transactions WHERE id = ?", tid).Scan(&s)
	return s
}
