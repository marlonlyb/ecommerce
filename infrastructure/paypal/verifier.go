package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const expectedVerification = "SUCCESS"

type payPalRequestValidator struct {
	// Headers
	AuthAlgo         string `json:"auth_algo"`
	CertURL          string `json:"cert_url"`
	TransmissionID   string `json:"transmission_id"`
	TransmissionSig  string `json:"transmission_sig"`
	TransmissionTime string `json:"transmission_time"`

	// Body
	WebhookID    string          `json:"webhook_id"`
	WebhookEvent json.RawMessage `json:"webhook_event"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Verifier struct {
	client        HTTPClient
	validationURL string
	clientID      string
	secretID      string
	webhookID     string
}

func NewVerifier(client HTTPClient) *Verifier {
	if client == nil {
		client = &http.Client{}
	}

	return &Verifier{
		client:        client,
		validationURL: os.Getenv("VALIDATION_URL"),
		clientID:      os.Getenv("CLIENT_ID"),
		secretID:      os.Getenv("SECRET_ID"),
		webhookID:     os.Getenv("WEBHOOK_ID"),
	}
}

func (v *Verifier) Verify(headers http.Header, body []byte) error {
	payload, err := json.Marshal(payPalRequestValidator{
		AuthAlgo:         headers.Get("Paypal-Auth-Algo"),
		CertURL:          headers.Get("Paypal-Cert-Url"),
		TransmissionID:   headers.Get("Paypal-Transmission-Id"),
		TransmissionSig:  headers.Get("Paypal-Transmission-Sig"),
		TransmissionTime: headers.Get("Paypal-Transmission-Time"),
		WebhookEvent:     body,
		WebhookID:        v.webhookID,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, v.validationURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(v.clientID, v.secretID)

	response, err := v.client.Do(request)
	if err != nil {
		return err
	}
	defer func(r *http.Response) {
		_ = r.Body.Close()
	}(response)

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("PayPal response with status code %d, body: %s", response.StatusCode, string(responseBody))
	}

	bodyMap := make(map[string]string)
	if err = json.Unmarshal(responseBody, &bodyMap); err != nil {
		return err
	}

	if bodyMap["verification_status"] != expectedVerification {
		return fmt.Errorf("verification status is %s", bodyMap["verification_status"])
	}

	return nil
}
