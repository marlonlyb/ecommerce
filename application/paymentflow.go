package application

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/invoice"
	"github.com/mlbautomation/ProyectoEMLB/domain/ports/purchaseorder"
)

const (
	expectedPayPalEventType = "PAYMENT.CAPTURE.COMPLETED"
	expectedPaymentStatus   = "completed"
)

type PayPalVerifier interface {
	Verify(header http.Header, body []byte) error
}

type payPalRequestData struct {
	EventType string `json:"event_type"`
	ID        string `json:"id"`
	Resource  struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		CustomID string `json:"custom_id"`
		Amount   struct {
			Value string `json:"value"`
		} `json:"amount"`
	} `json:"resource"`
}

type PaymentWebhookProcessor interface {
	Process(header http.Header, body []byte) error
}

type PaymentFlow struct {
	verifier             PayPalVerifier
	servicePurchaseOrder purchaseorder.ServicePaypal
	serviceInvoice       invoice.ServicePaypal
}

func NewPaymentFlow(verifier PayPalVerifier, purchaseOrderService purchaseorder.ServicePaypal, invoiceService invoice.ServicePaypal) *PaymentFlow {
	return &PaymentFlow{
		verifier:             verifier,
		servicePurchaseOrder: purchaseOrderService,
		serviceInvoice:       invoiceService,
	}
}

func (pf *PaymentFlow) Process(header http.Header, body []byte) error {
	requestData, err := pf.parsePayPalRequestData(body)
	if err != nil {
		return fmt.Errorf("pf.parsePayPalRequestData(): %w", err)
	}

	if pf.verifier == nil {
		return fmt.Errorf("paypal verifier is nil")
	}

	if err = pf.verifier.Verify(header, body); err != nil {
		return err
	}

	if err = pf.processPayment(&requestData); err != nil {
		return err
	}

	return nil
}

func (pf *PaymentFlow) parsePayPalRequestData(body []byte) (payPalRequestData, error) {
	data := payPalRequestData{}

	if err := json.Unmarshal(body, &data); err != nil {
		return payPalRequestData{}, fmt.Errorf("json.Unmarshal(): %w", err)
	}

	if data.EventType != expectedPayPalEventType {
		return payPalRequestData{}, fmt.Errorf("the event_type %q is not allowed", data.EventType)
	}

	return data, nil
}

func (pf *PaymentFlow) processPayment(data *payPalRequestData) error {
	if !strings.EqualFold(data.Resource.Status, expectedPaymentStatus) {
		return fmt.Errorf("el estado de la transacción: %q no es el estado esperado: %q", data.Resource.Status, expectedPaymentStatus)
	}

	ID, err := uuid.Parse(data.Resource.CustomID)
	if err != nil {
		return fmt.Errorf("uuid.Parse(): %w", err)
	}

	order, err := pf.servicePurchaseOrder.GetByID(ID)
	if err != nil {
		return fmt.Errorf("pf.servicePurchaseOrder.GetByID(ID): %w", err)
	}

	value, err := strconv.ParseFloat(data.Resource.Amount.Value, 64)
	if err != nil {
		return fmt.Errorf("strconv.ParseFloat(data.Resource.Amount.Value, 64): %w", err)
	}
	value = math.Floor(value*100) / 100

	totalAmount := pf.servicePurchaseOrder.TotalAmount(order)
	totalAmount = math.Floor(totalAmount*100) / 100

	if totalAmount != value {
		return fmt.Errorf("el valor recibido: %0.2f, es diferente al valor esperado %0.2f", value, totalAmount)
	}

	return pf.serviceInvoice.Create(&order)
}
