package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/katiem0/gh-seva/internal/data"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestGetOrgActionSecretsMock(t *testing.T) {
	// Setup - using the existing MockAPIGetter
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"secrets":[{"name":"TEST_SECRET","visibility":"all"}]}`)
	mockGetter.OrgActionSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgActionSecrets("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetRepoActionSecretsMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"secrets":[{"name":"TEST_SECRET"}]}`)
	mockGetter.RepoActionSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetRepoActionSecrets("test-org", "test-repo")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetScopedOrgActionSecretsMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"repositories":[{"id":12345,"name":"test-repo"}]}`)
	mockGetter.ScopedOrgActionSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetScopedOrgActionSecrets("test-org", "TEST_SECRET")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetOrgActionPublicKeyMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"key_id":"test-key-id","key":"base64-encoded-key"}`)
	mockGetter.PublicKeyData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgActionPublicKey("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetRepoActionPublicKeyMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"key_id":"test-key-id","key":"base64-encoded-key"}`)
	mockGetter.PublicKeyData = expectedResponse

	// Execute
	response, err := mockGetter.GetRepoActionPublicKey("test-org", "test-repo")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestCreateOrgActionSecretMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()

	// Execute
	err := mockGetter.CreateOrgActionSecret("test-org", "TEST_SECRET", nil)

	// Verify - since we're using a mock, there's not much to verify
	// other than that the function doesn't throw an error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateRepoActionSecretMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()

	// Execute
	err := mockGetter.CreateRepoActionSecret("test-org", "test-repo", "TEST_SECRET", nil)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetOrgActionSecretsError(t *testing.T) {
	// Create a custom mock that returns an error
	mockGetter := &errorMockAPIGetter{
		ErrorMessage: "network error",
	}

	// Call the overridden method that returns an error
	result, err := mockGetter.GetOrgActionSecrets("test-org")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected nil result with error, got %v", result)
	}

	if !strings.Contains(err.Error(), "network error") {
		t.Errorf("Expected 'network error' in error message, got: %v", err)
	}
}

func TestGetRepoActionSecretsSuccess(t *testing.T) {
	expectedData := data.SecretsResponse{
		TotalCount: 2,
		Secrets: []data.Secret{
			{Name: "SECRET1", Visibility: "private"},
			{Name: "SECRET2", Visibility: "private"},
		},
	}

	responseBody, _ := json.Marshal(expectedData)

	mockGetter := NewMockAPIGetter()
	mockGetter.RepoActionSecretsData = responseBody

	result, err := mockGetter.GetRepoActionSecrets("test-org", "test-repo")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	var response data.SecretsResponse
	err = json.Unmarshal(result, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.TotalCount != expectedData.TotalCount {
		t.Errorf("Expected TotalCount %d, got %d", expectedData.TotalCount, response.TotalCount)
	}

	if len(response.Secrets) != len(expectedData.Secrets) {
		t.Errorf("Expected %d secrets, got %d", len(expectedData.Secrets), len(response.Secrets))
	}
}

func TestGetScopedOrgActionSecrets(t *testing.T) {
	scopedResponse := `{
        "total_count": 2,
        "repositories": [
            {"id": 123, "name": "repo1"},
            {"id": 456, "name": "repo2"}
        ]
    }`

	mockGetter := NewMockAPIGetter()
	mockGetter.ScopedOrgActionSecretsData = []byte(scopedResponse)

	result, err := mockGetter.GetScopedOrgActionSecrets("test-org", "TEST_SECRET")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(string(result), `"total_count": 2`) {
		t.Error("Expected total_count of 2 in response")
	}

	if !strings.Contains(string(result), `{"id": 123, "name": "repo1"}`) {
		t.Error("Expected repo1 in response")
	}
}

// Test public key retrieval
func TestGetOrgActionPublicKey(t *testing.T) {
	expectedResponse := `{"key_id":"123","key":"base64key"}`

	mockGetter := NewMockAPIGetter()
	mockGetter.PublicKeyData = []byte(expectedResponse)

	result, err := mockGetter.GetOrgActionPublicKey("test-org")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(result) != expectedResponse {
		t.Errorf("Expected %s, got %s", expectedResponse, string(result))
	}
}

func TestEncryptSecret(t *testing.T) {
	expectedEncrypted := "encrypted_value_123"

	mockGetter := NewMockAPIGetter()
	mockGetter.EncryptedSecret = expectedEncrypted

	result, err := mockGetter.EncryptSecret("test_public_key", "secret_value")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != expectedEncrypted {
		t.Errorf("Expected %s, got %s", expectedEncrypted, result)
	}
}

func TestCreateOrgActionSecret(t *testing.T) {
	// First part - test success case
	mockGetter := NewMockAPIGetter()
	err := mockGetter.CreateOrgActionSecret("test-org", "TEST_SECRET", bytes.NewReader([]byte("{}")))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test with error condition - using a different mock implementation
	errorMockGetter := &secretCreateErrorMockAPIGetter{}
	err = errorMockGetter.CreateOrgActionSecret("test-org", "ERROR_SECRET", bytes.NewReader([]byte("{}")))
	if err == nil {
		t.Error("Expected error for error case, got nil")
	}
	if !strings.Contains(err.Error(), "failed to create secret") {
		t.Errorf("Expected 'failed to create secret' in error message, got: %v", err)
	}
}

func TestGetScopedOrgActionSecretsError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return nil, fmt.Errorf("forbidden")
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	t.Run("forbidden error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				// Expected to panic due to log.Fatal simulation
				if !strings.Contains(fmt.Sprintf("%v", r), "forbidden") {
					t.Errorf("Expected panic with forbidden error, got: %v", r)
				}
			} else {
				t.Error("Expected panic from log.Fatal")
			}
		}()

		_, err := getter.GetScopedOrgActionSecrets("test-org", "TEST_SECRET")
		// This line won't be reached due to panic, but linter is satisfied
		if err != nil {
			t.Logf("Got expected error: %v", err)
		}
	})
}

func TestGetOrgActionPublicKeyResponseCodes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	testCases := []struct {
		name         string
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{"Success", 200, `{"key_id":"123","key":"base64key"}`, false},
		{"Not Found", 404, "Not Found", false}, // Function returns body regardless
		{"Server Error", 500, "Internal Server Error", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockRESTClient{
				RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       io.NopCloser(bytes.NewReader([]byte(tc.responseBody))),
					}, nil
				},
			}

			getter := newAPIGetterWithMockREST(mockClient)

			result, err := getter.GetOrgActionPublicKey("test-org")

			if tc.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if string(result) != tc.responseBody {
				t.Errorf("Expected response body %s, got %s", tc.responseBody, string(result))
			}
		})
	}
}

func TestCreateOrgActionSecretScenarios(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	testCases := []struct {
		name        string
		secretName  string
		secretData  string
		statusCode  int
		shouldPanic bool
	}{
		{"Success", "TEST_SECRET", `{"encrypted_value":"encrypted","key_id":"123"}`, 201, false},
		{"Invalid Data", "BAD_SECRET", `invalid`, 400, true},
		{"Forbidden", "FORBIDDEN_SECRET", `{"encrypted_value":"encrypted","key_id":"123"}`, 403, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockRESTClient{
				RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
					if method != "PUT" {
						t.Errorf("Expected PUT method, got %s", method)
					}

					if !strings.Contains(path, fmt.Sprintf("orgs/test-org/actions/secrets/%s", tc.secretName)) {
						t.Errorf("Unexpected path: %s", path)
					}

					if tc.statusCode >= 400 {
						return nil, fmt.Errorf("HTTP %d", tc.statusCode)
					}

					return &http.Response{
						StatusCode: tc.statusCode,
						Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
					}, nil
				},
			}

			getter := newAPIGetterWithMockREST(mockClient)

			if tc.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic from log.Fatal")
					}
				}()
			}

			err := getter.CreateOrgActionSecret("test-org", tc.secretName, bytes.NewReader([]byte(tc.secretData)))

			if !tc.shouldPanic && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCreateRepoActionSecretPath(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	var capturedPath string
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			capturedPath = path
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	err := getter.CreateRepoActionSecret("test-org", "test-repo", "REPO_SECRET", bytes.NewReader([]byte("{}")))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedPath := "repos/test-org/test-repo/actions/secrets/REPO_SECRET"
	if capturedPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, capturedPath)
	}
}

func TestGetOrgActionSecretsWithRestError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	t.Run("network error handling", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				// Expected to panic due to log.Fatal simulation
				if !strings.Contains(fmt.Sprintf("%v", r), "network error") {
					t.Errorf("Expected panic with network error, got: %v", r)
				}
			} else {
				t.Error("Expected panic from log.Fatal, but code continued")
			}
		}()

		_, err := getter.GetOrgActionSecrets("test-org")
		// This line won't be reached due to panic, but linter is satisfied
		if err != nil {
			t.Logf("Got expected error: %v", err)
		}
	})
}
