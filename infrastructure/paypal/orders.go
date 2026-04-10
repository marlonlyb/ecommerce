package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	orderport "github.com/mlbautomation/ProyectoEMLB/domain/ports/order"
	"github.com/mlbautomation/ProyectoEMLB/model"
)

type OrdersClient struct {
	client   HTTPClient
	baseURL  string
	clientID string
	secretID string
}

func NewOrdersClient(client HTTPClient) *OrdersClient {
	if client == nil {
		client = &http.Client{}
	}

	return &OrdersClient{
		client:   client,
		baseURL:  resolvePayPalBaseURL(os.Getenv("VALIDATION_URL")),
		clientID: os.Getenv("CLIENT_ID"),
		secretID: os.Getenv("SECRET_ID"),
	}
}

func (c *OrdersClient) CreateOrder(order model.Order) (string, error) {
	token, err := c.authToken()
	if err != nil {
		return "", err
	}

	payload := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"custom_id": order.ID.String(),
				"amount": map[string]string{
					"currency_code": order.Currency,
					"value":         fmt.Sprintf("%.2f", order.Total),
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/v2/checkout/orders", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("paypal create order status=%d body=%s", resp.StatusCode, string(responseBody))
	}

	responseData := struct {
		ID string `json:"id"`
	}{}
	if err = json.Unmarshal(responseBody, &responseData); err != nil {
		return "", err
	}

	if responseData.ID == "" {
		return "", fmt.Errorf("paypal create order without id")
	}

	return responseData.ID, nil
}

func (c *OrdersClient) CaptureOrder(payPalOrderID string) (orderport.CaptureResult, error) {
	token, err := c.authToken()
	if err != nil {
		return orderport.CaptureResult{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/v2/checkout/orders/"+payPalOrderID+"/capture", nil)
	if err != nil {
		return orderport.CaptureResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(req)
	if err != nil {
		return orderport.CaptureResult{}, err
	}
	defer func() { _ = resp.Body.Close() }()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return orderport.CaptureResult{}, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return orderport.CaptureResult{}, fmt.Errorf("paypal capture status=%d body=%s", resp.StatusCode, string(responseBody))
	}

	responseData := struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		PurchaseUnits []struct {
			Amount struct {
				CurrencyCode string `json:"currency_code"`
				Value        string `json:"value"`
			} `json:"amount"`
			Payments struct {
				Captures []struct {
					ID     string `json:"id"`
					Status string `json:"status"`
					Amount struct {
						CurrencyCode string `json:"currency_code"`
						Value        string `json:"value"`
					} `json:"amount"`
				} `json:"captures"`
			} `json:"payments"`
		} `json:"purchase_units"`
	}{}
	if err = json.Unmarshal(responseBody, &responseData); err != nil {
		return orderport.CaptureResult{}, err
	}

	result := orderport.CaptureResult{OrderID: responseData.ID, Status: responseData.Status, RawPayload: map[string]interface{}{}}
	if len(responseData.PurchaseUnits) > 0 {
		unit := responseData.PurchaseUnits[0]
		result.Currency = unit.Amount.CurrencyCode
		fmt.Sscanf(unit.Amount.Value, "%f", &result.Total)
		if len(unit.Payments.Captures) > 0 {
			capture := unit.Payments.Captures[0]
			result.CaptureID = capture.ID
			result.Status = capture.Status
			if capture.Amount.CurrencyCode != "" {
				result.Currency = capture.Amount.CurrencyCode
			}
			if capture.Amount.Value != "" {
				fmt.Sscanf(capture.Amount.Value, "%f", &result.Total)
			}
		}
	}

	return result, nil
}

func (c *OrdersClient) authToken() (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/v1/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.secretID)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("paypal auth status=%d body=%s", resp.StatusCode, string(body))
	}

	authResponse := struct {
		AccessToken string `json:"access_token"`
	}{}
	if err = json.Unmarshal(body, &authResponse); err != nil {
		return "", err
	}

	if authResponse.AccessToken == "" {
		return "", fmt.Errorf("paypal auth without access token")
	}

	return authResponse.AccessToken, nil
}

func resolvePayPalBaseURL(validationURL string) string {
	if idx := strings.Index(validationURL, "/v1/"); idx > 0 {
		return strings.TrimRight(validationURL[:idx], "/")
	}
	return strings.TrimRight(validationURL, "/")
}
