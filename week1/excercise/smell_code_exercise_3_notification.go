package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CODE SMELL EXERCISE #3: Notification Service
// This code focuses on OOP design smells
//
// Smells present:
// - Refused bequest (inheritance used incorrectly)
// - Inappropriate intimacy (classes know too much about each other)
// - Data clumps (same group of data appears together)
// - Lazy class (class doesn't do enough)
// - Speculative generality (over-engineered)
// - Message chains (a.b().c().d())
// - Middle man (class just delegates)
// - Divergent change (one class changed for many reasons)

// NotificationService handles all notifications
type NotificationService struct {
	emailConfig EmailConfig
	smsConfig   SMSConfig
	pushConfig  PushConfig
}

// EmailConfig stores email configuration
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// SMSConfig stores SMS configuration
type SMSConfig struct {
	Provider string
	ApiKey   string
	Sender   string
}

// PushConfig stores push notification configuration
type PushConfig struct {
	Provider string
	ApiKey   string
}

// User represents a user
type User struct {
	ID            string
	Name          string
	Email         string
	Phone         string
	PushToken     string
	NotifyByEmail bool
	NotifyBySMS   bool
	NotifyByPush  bool
	Language      string
	Timezone      string
	EmailVerified bool
	PhoneVerified bool
}

// SendOrderConfirmation sends order confirmation
func (ns *NotificationService) SendOrderConfirmation(user User, orderId string, amount float64, items []string) error {
	// Prepare email
	if user.NotifyByEmail && user.EmailVerified {
		subject := ""
		body := ""

		if user.Language == "en" {
			subject = "Order Confirmation"
			body = fmt.Sprintf("Dear %s,\n\nYour order %s has been confirmed.\n\nAmount: $%.2f\n\nItems:\n", user.Name, orderId, amount)
			for _, item := range items {
				body += "- " + item + "\n"
			}
		} else if user.Language == "vi" {
			subject = "Xác nhận đơn hàng"
			body = fmt.Sprintf("Kính gửi %s,\n\nĐơn hàng %s đã được xác nhận.\n\nSố tiền: %.2f VND\n\nSản phẩm:\n", user.Name, orderId, amount)
			for _, item := range items {
				body += "- " + item + "\n"
			}
		}

		// Send email
		err := ns.sendEmail(user.Email, subject, body)
		if err != nil {
			return err
		}
	}

	// Prepare SMS
	if user.NotifyBySMS && user.PhoneVerified {
		message := ""

		if user.Language == "en" {
			message = fmt.Sprintf("Your order %s (Amount: $%.2f) has been confirmed.", orderId, amount)
		} else if user.Language == "vi" {
			message = fmt.Sprintf("Đơn hàng %s (Số tiền: %.2f VND) đã được xác nhận.", orderId, amount)
		}

		// Send SMS
		err := ns.sendSMS(user.Phone, message)
		if err != nil {
			return err
		}
	}

	// Prepare push notification
	if user.NotifyByPush && user.PushToken != "" {
		title := ""
		body := ""

		if user.Language == "en" {
			title = "Order Confirmed"
			body = fmt.Sprintf("Your order %s has been confirmed. Amount: $%.2f", orderId, amount)
		} else if user.Language == "vi" {
			title = "Đơn hàng đã xác nhận"
			body = fmt.Sprintf("Đơn hàng %s đã được xác nhận. Số tiền: %.2f VND", orderId, amount)
		}

		// Send push
		err := ns.sendPush(user.PushToken, title, body)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendOrderShipped sends order shipped notification
func (ns *NotificationService) SendOrderShipped(user User, orderId string, trackingNumber string) error {
	// Prepare email
	if user.NotifyByEmail && user.EmailVerified {
		subject := ""
		body := ""

		if user.Language == "en" {
			subject = "Order Shipped"
			body = fmt.Sprintf("Dear %s,\n\nYour order %s has been shipped.\n\nTracking Number: %s\n", user.Name, orderId, trackingNumber)
		} else if user.Language == "vi" {
			subject = "Đơn hàng đã gửi"
			body = fmt.Sprintf("Kính gửi %s,\n\nĐơn hàng %s đã được gửi đi.\n\nMã vận đơn: %s\n", user.Name, orderId, trackingNumber)
		}

		// Send email
		err := ns.sendEmail(user.Email, subject, body)
		if err != nil {
			return err
		}
	}

	// Prepare SMS
	if user.NotifyBySMS && user.PhoneVerified {
		message := ""

		if user.Language == "en" {
			message = fmt.Sprintf("Your order %s has been shipped. Tracking: %s", orderId, trackingNumber)
		} else if user.Language == "vi" {
			message = fmt.Sprintf("Đơn hàng %s đã được gửi. Mã vận đơn: %s", orderId, trackingNumber)
		}

		// Send SMS
		err := ns.sendSMS(user.Phone, message)
		if err != nil {
			return err
		}
	}

	// Prepare push notification
	if user.NotifyByPush && user.PushToken != "" {
		title := ""
		body := ""

		if user.Language == "en" {
			title = "Order Shipped"
			body = fmt.Sprintf("Your order %s has been shipped. Track it with %s", orderId, trackingNumber)
		} else if user.Language == "vi" {
			title = "Đơn hàng đã gửi"
			body = fmt.Sprintf("Đơn hàng %s đã được gửi. Theo dõi với mã %s", orderId, trackingNumber)
		}

		// Send push
		err := ns.sendPush(user.PushToken, title, body)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendPaymentFailed sends payment failed notification
func (ns *NotificationService) SendPaymentFailed(user User, orderId string, reason string) error {
	// Prepare email
	if user.NotifyByEmail && user.EmailVerified {
		subject := ""
		body := ""

		if user.Language == "en" {
			subject = "Payment Failed"
			body = fmt.Sprintf("Dear %s,\n\nPayment for order %s has failed.\n\nReason: %s\n\nPlease try again.\n", user.Name, orderId, reason)
		} else if user.Language == "vi" {
			subject = "Thanh toán thất bại"
			body = fmt.Sprintf("Kính gửi %s,\n\nThanh toán cho đơn hàng %s đã thất bại.\n\nLý do: %s\n\nVui lòng thử lại.\n", user.Name, orderId, reason)
		}

		// Send email
		err := ns.sendEmail(user.Email, subject, body)
		if err != nil {
			return err
		}
	}

	// Prepare SMS
	if user.NotifyBySMS && user.PhoneVerified {
		message := ""

		if user.Language == "en" {
			message = fmt.Sprintf("Payment failed for order %s. Reason: %s. Please try again.", orderId, reason)
		} else if user.Language == "vi" {
			message = fmt.Sprintf("Thanh toán thất bại cho đơn hàng %s. Lý do: %s. Vui lòng thử lại.", orderId, reason)
		}

		// Send SMS
		err := ns.sendSMS(user.Phone, message)
		if err != nil {
			return err
		}
	}

	// Prepare push notification
	if user.NotifyByPush && user.PushToken != "" {
		title := ""
		body := ""

		if user.Language == "en" {
			title = "Payment Failed"
			body = fmt.Sprintf("Payment for order %s failed. %s. Please try again.", orderId, reason)
		} else if user.Language == "vi" {
			title = "Thanh toán thất bại"
			body = fmt.Sprintf("Thanh toán cho đơn hàng %s thất bại. %s. Vui lòng thử lại.", orderId, reason)
		}

		// Send push
		err := ns.sendPush(user.PushToken, title, body)
		if err != nil {
			return err
		}
	}

	return nil
}

// sendEmail sends an email
func (ns *NotificationService) sendEmail(to string, subject string, body string) error {
	// Create email payload
	payload := map[string]interface{}{
		"from":    ns.emailConfig.From,
		"to":      to,
		"subject": subject,
		"body":    body,
	}

	jsonData, _ := json.Marshal(payload)

	// Call email API
	req, _ := http.NewRequest("POST", "https://email-api.com/send", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(ns.emailConfig.Username, ns.emailConfig.Password)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("email send failed with status %d", resp.StatusCode)
	}

	return nil
}

// sendSMS sends an SMS
func (ns *NotificationService) sendSMS(to string, message string) error {
	// Create SMS payload
	payload := map[string]interface{}{
		"from":    ns.smsConfig.Sender,
		"to":      to,
		"message": message,
	}

	jsonData, _ := json.Marshal(payload)

	// Call SMS API
	req, _ := http.NewRequest("POST", "https://sms-api.com/send", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ns.smsConfig.ApiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("sms send failed with status %d", resp.StatusCode)
	}

	return nil
}

// sendPush sends a push notification
func (ns *NotificationService) sendPush(token string, title string, body string) error {
	// Create push payload
	payload := map[string]interface{}{
		"token": token,
		"title": title,
		"body":  body,
	}

	jsonData, _ := json.Marshal(payload)

	// Call Push API
	req, _ := http.NewRequest("POST", "https://push-api.com/send", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ns.pushConfig.ApiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("push send failed with status %d", resp.StatusCode)
	}

	return nil
}

// TemplateManager manages notification templates (lazy class - doesn't do much)
type TemplateManager struct {
	templates map[string]string
}

func (tm *TemplateManager) GetTemplate(name string) string {
	return tm.templates[name]
}

// MessageFormatter formats messages (middle man - just delegates)
type MessageFormatter struct {
	templateManager *TemplateManager
}

func (mf *MessageFormatter) Format(template string, data map[string]interface{}) string {
	message := mf.templateManager.GetTemplate(template)
	// Simple string replacement (over-simplified)
	for key, value := range data {
		message = fmt.Sprintf(message, value)
	}
	return message
}
