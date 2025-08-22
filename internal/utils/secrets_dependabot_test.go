package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestGetOrgDependabotSecretsMock(t *testing.T) {
	// Setup - using the existing MockAPIGetter
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"secrets":[{"name":"TEST_SECRET","visibility":"all"}]}`)
	mockGetter.OrgDependabotSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgDependabotSecrets("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetRepoDependabotSecretsMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"secrets":[{"name":"TEST_SECRET"}]}`)
	mockGetter.RepoDependabotSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetRepoDependabotSecrets("test-org", "test-repo")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetScopedOrgDependabotSecretsMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"repositories":[{"id":12345,"name":"test-repo"}]}`)
	mockGetter.ScopedOrgDependabotSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetScopedOrgDependabotSecrets("test-org", "TEST_SECRET")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetOrgDependabotPublicKeyMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"key_id":"test-key-id","key":"base64-encoded-key"}`)
	mockGetter.PublicKeyData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgDependabotPublicKey("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetRepoDependabotPublicKeyMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"key_id":"test-key-id","key":"base64-encoded-key"}`)
	mockGetter.PublicKeyData = expectedResponse

	// Execute
	response, err := mockGetter.GetRepoDependabotPublicKey("test-org", "test-repo")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestCreateOrgDependabotSecretMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()

	// Execute
	err := mockGetter.CreateOrgDependabotSecret("test-org", "TEST_SECRET", nil)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateRepoDependabotSecretMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()

	// Execute
	err := mockGetter.CreateRepoDependabotSecret("test-org", "test-repo", "TEST_SECRET", nil)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetOrgDependabotSecretsNetworkError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return nil, fmt.Errorf("connection timeout")
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	t.Run("network timeout", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic from log.Fatal")
			} else if !strings.Contains(fmt.Sprintf("%v", r), "connection timeout") {
				t.Errorf("Expected connection timeout error, got: %v", r)
			}
		}()

		_, err := getter.GetOrgDependabotSecrets("test-org")
		// This line won't be reached due to panic, but linter is satisfied
		if err != nil {
			t.Logf("Got expected error: %v", err)
		}
	})
}

// Test GetRepoDependabotSecrets with malformed response
func TestGetRepoDependabotSecretsMalformedResponse(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Return malformed JSON
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"invalid json`))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	result, err := getter.GetRepoDependabotSecrets("test-org", "test-repo")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// The function returns the raw response even if it's malformed
	if !strings.Contains(string(result), `{"invalid json`) {
		t.Errorf("Expected malformed JSON in response, got: %s", string(result))
	}
}

// Test GetOrgDependabotPublicKey with different key formats
func TestGetOrgDependabotPublicKeyFormats(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	testCases := []struct {
		name         string
		responseBody string
	}{
		{"Standard format", `{"key_id":"12345","key":"base64encodedkey=="}`},
		{"With extra fields", `{"key_id":"67890","key":"anotherkey==","created_at":"2023-01-01"}`},
		{"Minimal format", `{"key_id":"1","key":"k"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockRESTClient{
				RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader([]byte(tc.responseBody))),
					}, nil
				},
			}

			getter := newAPIGetterWithMockREST(mockClient)

			result, err := getter.GetOrgDependabotPublicKey("test-org")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if string(result) != tc.responseBody {
				t.Errorf("Expected %s, got %s", tc.responseBody, string(result))
			}
		})
	}
}

// Test CreateOrgDependabotSecret with large payload
func TestCreateOrgDependabotSecretLargePayload(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	// Create a large payload
	largePayload := bytes.Repeat([]byte("a"), 64*1024) // 64KB

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Read the body to ensure it's the expected size
			bodyBytes, _ := io.ReadAll(body)
			if len(bodyBytes) != len(largePayload) {
				t.Errorf("Expected payload size %d, got %d", len(largePayload), len(bodyBytes))
			}

			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	err := getter.CreateOrgDependabotSecret("test-org", "LARGE_SECRET", bytes.NewReader(largePayload))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Test GetScopedOrgDependabotSecrets with empty response
func TestGetScopedOrgDependabotSecretsEmpty(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"total_count":0,"repositories":[]}`))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	result, err := getter.GetScopedOrgDependabotSecrets("test-org", "EMPTY_SECRET")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !strings.Contains(string(result), `"total_count":0`) {
		t.Error("Expected total_count of 0")
	}

	if !strings.Contains(string(result), `"repositories":[]`) {
		t.Error("Expected empty repositories array")
	}
}

// Test CreateRepoDependabotSecret error scenarios
func TestCreateRepoDependabotSecretErrors(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	testCases := []struct {
		name          string
		returnError   error
		expectedPanic string
	}{
		{"Network error", fmt.Errorf("network unreachable"), "network unreachable"},
		{"Auth error", fmt.Errorf("authentication failed"), "authentication failed"},
		{"Rate limit", fmt.Errorf("rate limit exceeded"), "rate limit exceeded"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockRESTClient{
				RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
					return nil, tc.returnError
				},
			}

			getter := newAPIGetterWithMockREST(mockClient)

			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic from log.Fatal")
				} else if !strings.Contains(fmt.Sprintf("%v", r), tc.expectedPanic) {
					t.Errorf("Expected panic containing '%s', got: %v", tc.expectedPanic, r)
				}
			}()

			err := getter.CreateRepoDependabotSecret("test-org", "test-repo", "SECRET", bytes.NewReader([]byte("{}")))
			// This line won't be reached due to panic, but linter is satisfied
			if err != nil {
				t.Logf("Got expected error: %v", err)
			}
		})
	}
}
