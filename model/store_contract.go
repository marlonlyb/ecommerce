package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderStatusPendingPayment = "pending_payment"
	OrderStatusPaid           = "paid"
	OrderStatusPaymentFailed  = "payment_failed"
	OrderStatusCancelled      = "cancelled"
	OrderStatusRefunded       = "refunded"

	PaymentProviderPayPal = "paypal"

	PaymentStatusPending  = "pending"
	PaymentStatusApproved = "approved"
	PaymentStatusCaptured = "captured"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

type APIErrorDetail struct {
	Field string `json:"field,omitempty"`
	Issue string `json:"issue"`
}

type APIErrorPayload struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	Details   []APIErrorDetail `json:"details,omitempty"`
	RequestID string           `json:"request_id,omitempty"`
}

type APIErrorResponse struct {
	Error APIErrorPayload `json:"error"`
}

type StoreUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type StoreProductVariant struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	SKU       string    `json:"sku"`
	Color     string    `json:"color"`
	Size      string    `json:"size"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	ImageURL  string    `json:"image_url,omitempty"`
}

type StoreProduct struct {
	ID              uuid.UUID             `json:"id"`
	Name            string                `json:"name"`
	Slug            string                `json:"slug"`
	Description     string                `json:"description"`
	Category        string                `json:"category"`
	Brand           string                `json:"brand,omitempty"`
	Images          []string              `json:"images"`
	Active          bool                  `json:"active"`
	PriceFrom       float64               `json:"price_from,omitempty"`
	AvailableColors []string              `json:"available_colors,omitempty"`
	AvailableSizes  []string              `json:"available_sizes,omitempty"`
	Variants        []StoreProductVariant `json:"variants,omitempty"`
}

type OrderItem struct {
	ID          uuid.UUID `json:"id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	VariantID   uuid.UUID `json:"variant_id"`
	VariantSKU  string    `json:"variant_sku"`
	Color       string    `json:"color"`
	Size        string    `json:"size"`
	UnitPrice   float64   `json:"unit_price"`
	Quantity    int       `json:"quantity"`
	LineTotal   float64   `json:"line_total"`
}

type Order struct {
	ID              uuid.UUID   `json:"id"`
	UserID          uuid.UUID   `json:"user_id"`
	Status          string      `json:"status"`
	PaymentProvider string      `json:"payment_provider"`
	PaymentStatus   string      `json:"payment_status"`
	Currency        string      `json:"currency"`
	Subtotal        float64     `json:"subtotal"`
	Total           float64     `json:"total"`
	PayPalOrderID   string      `json:"paypal_order_id,omitempty"`
	PayPalCaptureID string      `json:"paypal_capture_id,omitempty"`
	PaidAt          *time.Time  `json:"paid_at,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	Items           []OrderItem `json:"items,omitempty"`
}

type CheckoutPayPalItem struct {
	VariantID uuid.UUID `json:"variant_id"`
	Quantity  int       `json:"quantity"`
}

type CheckoutPayPalRequest struct {
	Items []CheckoutPayPalItem `json:"items"`
}

type CapturePayPalOrderRequest struct {
	PayPalOrderID string `json:"paypal_order_id"`
}
