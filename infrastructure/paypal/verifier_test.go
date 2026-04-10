package paypal

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestVerifier_Verify(t *testing.T) {
	tests := []struct {
		name        string
		clientErr   error
		statusCode  int
		body        string
		wantErr     bool
		errContains string
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			body:       `{"verification_status":"SUCCESS"}`,
			wantErr:    false,
		},
		{
			name:        "http client error",
			clientErr:   errors.New("network error"),
			wantErr:     true,
			errContains: "network error",
		},
		{
			name:        "non-200 status code",
			statusCode:  http.StatusBadRequest,
			body:        `{"error":"bad request"}`,
			wantErr:     true,
			errContains: "PayPal response with status code 400",
		},
		{
			name:        "invalid json response",
			statusCode:  http.StatusOK,
			body:        `{invalid json}`,
			wantErr:     true,
			errContains: "invalid character",
		},
		{
			name:        "verification failed",
			statusCode:  http.StatusOK,
			body:        `{"verification_status":"FAILURE"}`,
			wantErr:     true,
			errContains: "verification status is FAILURE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			if tt.clientErr == nil {
				resp = &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
				}
			}

			v := &Verifier{
				client:        &mockHTTPClient{response: resp, err: tt.clientErr},
				validationURL: "http://example.com/validate",
				clientID:      "client",
				secretID:      "secret",
				webhookID:     "webhook",
			}

			err := v.Verify(http.Header{}, []byte(`{"some":"event"}`))

			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Verify() error = %v, expected it to contain %q", err, tt.errContains)
				}
			}
		})
	}
}
