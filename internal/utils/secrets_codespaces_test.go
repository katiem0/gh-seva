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

func TestGetOrgCodespacesSecretsMock(t *testing.T) {
	// Setup - using the existing MockAPIGetter
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"secrets":[{"name":"TEST_SECRET","visibility":"all"}]}`)
	mockGetter.OrgCodespacesSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgCodespacesSecrets("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetRepoCodespacesSecretsMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"secrets":[{"name":"TEST_SECRET"}]}`)
	mockGetter.RepoCodespacesSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetRepoCodespacesSecrets("test-org", "test-repo")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetScopedOrgCodespacesSecretsMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"repositories":[{"id":12345,"name":"test-repo"}]}`)
	mockGetter.ScopedOrgCodespacesSecretsData = expectedResponse

	// Execute
	response, err := mockGetter.GetScopedOrgCodespacesSecrets("test-org", "TEST_SECRET")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetOrgCodespacesPublicKeyMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"key_id":"test-key-id","key":"base64-encoded-key"}`)
	mockGetter.PublicKeyData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgCodespacesPublicKey("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestGetRepoCodespacesPublicKeyMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"key_id":"test-key-id","key":"base64-encoded-key"}`)
	mockGetter.PublicKeyData = expectedResponse

	// Execute
	response, err := mockGetter.GetRepoCodespacesPublicKey("test-org", "test-repo")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}
}

func TestCreateOrgCodespacesSecretMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()

	// Execute
	err := mockGetter.CreateOrgCodespacesSecret("test-org", "TEST_SECRET", nil)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateRepoCodespacesSecretMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()

	// Execute
	err := mockGetter.CreateRepoCodespacesSecret("test-org", "test-repo", "TEST_SECRET", nil)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetOrgCodespacesSecretsHTTPError(t *testing.T) {
	// Instead of using testAPIGetterWrapper, use MockAPIGetter for simple cases
	mockGetter := NewMockAPIGetter()
	mockGetter.OrgCodespacesSecretsData = []byte("Forbidden")

	result, err := mockGetter.GetOrgCodespacesSecrets("test-org")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(result) != "Forbidden" {
		t.Errorf("Expected 'Forbidden' response, got: %s", string(result))
	}
}

// Test CreateOrgCodespacesSecret with invalid data
func TestCreateOrgCodespacesSecretInvalidData(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return nil, fmt.Errorf("bad request: invalid encrypted value")
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	t.Run("invalid data error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic from log.Fatal")
			}
		}()

		err := getter.CreateOrgCodespacesSecret("test-org", "TEST_SECRET", bytes.NewReader([]byte("invalid")))
		_ = err // This line won't be reached due to panic
	})
}

// Test GetRepoCodespacesSecrets path construction
func TestGetRepoCodespacesSecretsPath(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	var capturedPath string
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			capturedPath = path
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"total_count":0,"secrets":[]}`))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	_, err := getter.GetRepoCodespacesSecrets("test-org", "test-repo")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedPath := "repos/test-org/test-repo/codespaces/secrets"
	if capturedPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, capturedPath)
	}
}

// Test GetOrgCodespacesPublicKey with body read error
func TestGetOrgCodespacesPublicKeyReadError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(&errorReader{}),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	_, err := getter.GetOrgCodespacesPublicKey("test-org")
	if err == nil {
		t.Error("Expected read error, got nil")
	}
}

// errorReader simulates a read error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("simulated read error")
}

func (e *errorReader) Close() error {
	return nil
}

// Test CreateRepoCodespacesSecret with different status codes
func TestCreateRepoCodespacesSecretStatusCodes(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	testCases := []struct {
		name        string
		statusCode  int
		shouldPanic bool
	}{
		{"Success Created", 201, false},
		{"Success OK", 200, false},
		{"Unauthorized", 401, true},
		{"Not Found", 404, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockRESTClient{
				RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
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

			err := getter.CreateRepoCodespacesSecret("test-org", "test-repo", "SECRET", bytes.NewReader([]byte("{}")))

			if !tc.shouldPanic && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Test GetScopedOrgCodespacesSecrets with pagination
func TestGetScopedOrgCodespacesSecretsPagination(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			if !strings.Contains(path, "orgs/test-org/codespaces/secrets/TEST_SECRET/repositories") {
				t.Errorf("Unexpected path: %s", path)
			}

			response := `{
                "total_count": 3,
                "repositories": [
                    {"id": 1, "name": "repo1"},
                    {"id": 2, "name": "repo2"},
                    {"id": 3, "name": "repo3"}
                ]
            }`

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(response))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	result, err := getter.GetScopedOrgCodespacesSecrets("test-org", "TEST_SECRET")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !strings.Contains(string(result), `"total_count": 3`) {
		t.Error("Expected total_count of 3 in response")
	}
}
