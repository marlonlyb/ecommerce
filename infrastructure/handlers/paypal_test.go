package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mlbautomation/ProyectoEMLB/infrastructure/handlers/response"
)

type mockPaymentWebhookProcessor struct {
	err error
}

func (m *mockPaymentWebhookProcessor) Process(header http.Header, body []byte) error {
	return m.err
}

func TestPaypal_Webhook(t *testing.T) {
	tests := []struct {
		name         string
		reqBody      string
		processorErr error
		expectedCode int
	}{
		{
			name:         "success",
			reqBody:      `{"event_type":"PAYMENT.CAPTURE.COMPLETED"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "processor error",
			reqBody:      `{"event_type":"PAYMENT.CAPTURE.COMPLETED"}`,
			processorErr: errors.New("processor failed"),
			expectedCode: http.StatusOK, // The endpoint returns 200 OK anyway for paypal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/public/paypal", bytes.NewBufferString(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockProcessor := &mockPaymentWebhookProcessor{err: tt.processorErr}
			h := &Paypal{
				processor: mockProcessor,
				responser: response.API{},
			}

			err := h.Webhook(c)

			if err != nil {
				t.Errorf("Webhook() unexpected error: %v", err)
			}

			if rec.Code != tt.expectedCode {
				t.Errorf("Webhook() status code = %v, want %v", rec.Code, tt.expectedCode)
			}
		})
	}
}
