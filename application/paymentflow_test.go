package application

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

// Mock objects
type mockVerifier struct {
	err error
}

func (m *mockVerifier) Verify(header http.Header, body []byte) error {
	return m.err
}

type mockPurchaseOrderService struct {
	order       model.PurchaseOrder
	totalAmount float64
	err         error
}

func (m *mockPurchaseOrderService) GetByID(ID uuid.UUID) (model.PurchaseOrder, error) {
	return m.order, m.err
}

func (m *mockPurchaseOrderService) TotalAmount(order model.PurchaseOrder) float64 {
	return m.totalAmount
}

type mockInvoiceService struct {
	err error
}

func (m *mockInvoiceService) Create(order *model.PurchaseOrder) error {
	return m.err
}

func TestPaymentFlow_Process(t *testing.T) {
	validUUID := uuid.New().String()

	tests := []struct {
		name          string
		body          string
		verifierErr   error
		poGetErr      error
		poTotalAmount float64
		invCreateErr  error
		wantErr       bool
		errContains   string
	}{
		{
			name:          "success",
			body:          `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"completed","custom_id":"` + validUUID + `","amount":{"value":"100.50"}}}`,
			poTotalAmount: 100.50,
			wantErr:       false,
		},
		{
			name:        "invalid json",
			body:        `{invalid}`,
			wantErr:     true,
			errContains: "json.Unmarshal()",
		},
		{
			name:        "invalid event type",
			body:        `{"event_type":"OTHER_EVENT"}`,
			wantErr:     true,
			errContains: "not allowed",
		},
		{
			name:        "verifier error",
			body:        `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"completed","custom_id":"` + validUUID + `","amount":{"value":"100.50"}}}`,
			verifierErr: errors.New("verification failed"),
			wantErr:     true,
			errContains: "verification failed",
		},
		{
			name:        "payment not completed",
			body:        `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"pending","custom_id":"` + validUUID + `","amount":{"value":"100.50"}}}`,
			wantErr:     true,
			errContains: "no es el estado esperado",
		},
		{
			name:        "invalid custom id",
			body:        `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"completed","custom_id":"invalid-uuid","amount":{"value":"100.50"}}}`,
			wantErr:     true,
			errContains: "uuid.Parse()",
		},
		{
			name:        "purchase order not found",
			body:        `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"completed","custom_id":"` + validUUID + `","amount":{"value":"100.50"}}}`,
			poGetErr:    errors.New("order not found"),
			wantErr:     true,
			errContains: "pf.servicePurchaseOrder.GetByID(ID)",
		},
		{
			name:        "invalid amount format",
			body:        `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"completed","custom_id":"` + validUUID + `","amount":{"value":"invalid-float"}}}`,
			wantErr:     true,
			errContains: "strconv.ParseFloat",
		},
		{
			name:          "amount mismatch",
			body:          `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"completed","custom_id":"` + validUUID + `","amount":{"value":"100.50"}}}`,
			poTotalAmount: 50.00, // Mismatch
			wantErr:       true,
			errContains:   "diferente al valor esperado",
		},
		{
			name:          "invoice creation error",
			body:          `{"event_type":"PAYMENT.CAPTURE.COMPLETED","resource":{"status":"completed","custom_id":"` + validUUID + `","amount":{"value":"100.50"}}}`,
			poTotalAmount: 100.50,
			invCreateErr:  errors.New("could not create invoice"),
			wantErr:       true,
			errContains:   "could not create invoice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := NewPaymentFlow(
				&mockVerifier{err: tt.verifierErr},
				&mockPurchaseOrderService{
					order:       model.PurchaseOrder{ID: uuid.MustParse(validUUID)},
					err:         tt.poGetErr,
					totalAmount: tt.poTotalAmount,
				},
				&mockInvoiceService{err: tt.invCreateErr},
			)

			err := pf.Process(http.Header{}, []byte(tt.body))

			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Process() error = %v, expected it to contain %q", err, tt.errContains)
				}
			}
		})
	}
}
