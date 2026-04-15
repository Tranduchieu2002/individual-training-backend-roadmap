package promo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// CODE SMELL EXERCISE #2: Promo Code Manager
// This code focuses on different smells than Exercise #1
//
// Smells present:
// - Duplicate code everywhere
// - Comment pollution (useless comments)
// - Inconsistent naming
// - No constants for magic strings
// - Boolean parameters (flag arguments)
// - Long parameter list
// - Switch statements (type codes)
// - Feature envy (accessing too much from other objects)
// - Shotgun surgery (one change requires many edits)

// PromoCodeManager manages promo codes
type PromoCodeManager struct {
	db *sql.DB
}

// CreatePromoCode creates a new promo code
// code: the promo code string
// type: type of promo (1=percentage, 2=fixed, 3=free_shipping, 4=cashback)
// value: the value (percentage or fixed amount)
// max: maximum discount
// min: minimum order
// exp: expiry date
// limit: usage limit
// user: limit to specific user
// new: is this for new users only?
func (m *PromoCodeManager) CreatePromoCode(code string, promoType int, value float64, maxDiscount float64, minOrder float64, expiryDate string, usageLimit int, userSpecific string, newUserOnly bool) error {
	// Validate code - must be uppercase and alphanumeric
	// Code must be between 4 and 20 characters
	if len(code) < 4 || len(code) > 20 {
		return errors.New("code length invalid")
	}

	// Check if code already exists
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM promo_codes WHERE code = ?", code).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("code already exists")
	}

	// Validate promo type
	if promoType != 1 && promoType != 2 && promoType != 3 && promoType != 4 {
		return errors.New("invalid promo type")
	}

	// For percentage type, value should be between 0 and 100
	if promoType == 1 {
		if value <= 0 || value > 100 {
			return errors.New("percentage must be between 0 and 100")
		}
	}

	// For fixed type, value should be positive
	if promoType == 2 {
		if value <= 0 {
			return errors.New("fixed value must be positive")
		}
	}

	// Insert into database
	_, err = m.db.Exec(
		"INSERT INTO promo_codes (code, type, value, max_discount, min_order, expiry_date, usage_limit, user_specific, new_user_only, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		code, promoType, value, maxDiscount, minOrder, expiryDate, usageLimit, userSpecific, newUserOnly, time.Now(),
	)

	return err
}

// ValidatePromoCode validates a promo code for a user and order
func (m *PromoCodeManager) ValidatePromoCode(code string, userId string, orderAmount float64, isNewUser bool) (bool, string, error) {
	// Get promo code from database
	var promoType int
	var value float64
	var maxDiscount float64
	var minOrder float64
	var expiryDate string
	var usageLimit int
	var usageCount int
	var userSpecific string
	var newUserOnly bool

	err := m.db.QueryRow(
		"SELECT type, value, max_discount, min_order, expiry_date, usage_limit, usage_count, user_specific, new_user_only FROM promo_codes WHERE code = ?",
		code,
	).Scan(&promoType, &value, &maxDiscount, &minOrder, &expiryDate, &usageLimit, &usageCount, &userSpecific, &newUserOnly)

	if err != nil {
		return false, "Promo code not found", err
	}

	// Check if expired
	expTime, _ := time.Parse("2006-01-02", expiryDate)
	if time.Now().After(expTime) {
		return false, "Promo code expired", nil
	}

	// Check usage limit
	if usageLimit > 0 && usageCount >= usageLimit {
		return false, "Promo code usage limit reached", nil
	}

	// Check if user specific
	if userSpecific != "" && userSpecific != userId {
		return false, "Promo code not valid for this user", nil
	}

	// Check if new user only
	if newUserOnly && !isNewUser {
		return false, "Promo code valid for new users only", nil
	}

	// Check minimum order
	if orderAmount < minOrder {
		return false, fmt.Sprintf("Minimum order amount is %.2f", minOrder), nil
	}

	// Check if user already used this code
	var userUsageCount int
	m.db.QueryRow("SELECT COUNT(*) FROM promo_code_usage WHERE code = ? AND user_id = ?", code, userId).Scan(&userUsageCount)
	if userUsageCount > 0 {
		return false, "You have already used this promo code", nil
	}

	return true, "Valid", nil
}

// CalculateDiscount calculates discount for a promo code
func (m *PromoCodeManager) CalculateDiscount(code string, orderAmount float64, shippingCost float64) (float64, error) {
	// Get promo code from database
	var promoType int
	var value float64
	var maxDiscount float64

	err := m.db.QueryRow(
		"SELECT type, value, max_discount FROM promo_codes WHERE code = ?",
		code,
	).Scan(&promoType, &value, &maxDiscount)

	if err != nil {
		return 0, err
	}

	var discount float64

	// Calculate based on type
	// Type 1: Percentage
	if promoType == 1 {
		discount = orderAmount * (value / 100)
		if maxDiscount > 0 && discount > maxDiscount {
			discount = maxDiscount
		}
	}

	// Type 2: Fixed amount
	if promoType == 2 {
		discount = value
		if discount > orderAmount {
			discount = orderAmount
		}
	}

	// Type 3: Free shipping
	if promoType == 3 {
		discount = shippingCost
	}

	// Type 4: Cashback (no immediate discount)
	if promoType == 4 {
		discount = 0
	}

	return discount, nil
}

// ApplyPromoCode applies promo code to an order
func (m *PromoCodeManager) ApplyPromoCode(code string, userId string, orderId string, orderAmount float64, shippingCost float64) error {
	// Get promo code from database
	var promoType int
	var value float64

	err := m.db.QueryRow(
		"SELECT type, value FROM promo_codes WHERE code = ?",
		code,
	).Scan(&promoType, &value)

	if err != nil {
		return err
	}

	// Increment usage count
	_, err = m.db.Exec("UPDATE promo_codes SET usage_count = usage_count + 1 WHERE code = ?", code)
	if err != nil {
		return err
	}

	// Record usage
	_, err = m.db.Exec(
		"INSERT INTO promo_code_usage (code, user_id, order_id, used_at) VALUES (?, ?, ?, ?)",
		code, userId, orderId, time.Now(),
	)
	if err != nil {
		return err
	}

	// If cashback, add to user wallet
	if promoType == 4 {
		cashbackAmount := orderAmount * (value / 100)
		_, err = m.db.Exec("UPDATE wallets SET balance = balance + ? WHERE user_id = ?", cashbackAmount, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPromoCodesByUser returns all promo codes available for a user
func (m *PromoCodeManager) GetPromoCodesByUser(userId string, isNewUser bool) ([]string, error) {
	var codes []string

	// Get all non-expired codes
	rows, err := m.db.Query(
		"SELECT code FROM promo_codes WHERE expiry_date > ? AND (user_specific = '' OR user_specific = ?)",
		time.Now().Format("2006-01-02"), userId,
	)
	if err != nil {
		return codes, err
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		rows.Scan(&code)
		codes = append(codes, code)
	}

	return codes, nil
}

// DeactivatePromoCode deactivates a promo code
func (m *PromoCodeManager) DeactivatePromoCode(code string, reason string) error {
	// Update promo code status
	_, err := m.db.Exec(
		"UPDATE promo_codes SET is_active = 0, deactivated_at = ?, deactivation_reason = ? WHERE code = ?",
		time.Now(), reason, code,
	)
	return err
}

// Helper function to check if code is valid format
func isValidCodeFormat(code string) bool {
	// Code must be uppercase
	if strings.ToUpper(code) != code {
		return false
	}

	// Code must be alphanumeric
	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')) {
			return false
		}
	}

	return true
}

// Helper function to get promo type name
func getPromoTypeName(promoType int) string {
	// Return promo type name based on type code
	if promoType == 1 {
		return "Percentage Discount"
	} else if promoType == 2 {
		return "Fixed Amount Discount"
	} else if promoType == 3 {
		return "Free Shipping"
	} else if promoType == 4 {
		return "Cashback"
	}
	return "Unknown"
}
