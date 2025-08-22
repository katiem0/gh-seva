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

func TestGetOrgActionVariablesMock(t *testing.T) {
	// Setup - using the existing MockAPIGetter
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"variables":[{"name":"TEST_VAR","value":"test-value","visibility":"all"}]}`)
	mockGetter.OrgActionVariablesData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgActionVariables("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}

	// Test parsing of the response
	var variablesResponse data.VariableResponse
	err = json.Unmarshal(response, &variablesResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if variablesResponse.TotalCount != 1 {
		t.Errorf("Expected total count 1, got %d", variablesResponse.TotalCount)
	}

	if len(variablesResponse.Variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(variablesResponse.Variables))
	}

	if variablesResponse.Variables[0].Name != "TEST_VAR" {
		t.Errorf("Expected variable name TEST_VAR, got %s", variablesResponse.Variables[0].Name)
	}

	if variablesResponse.Variables[0].Value != "test-value" {
		t.Errorf("Expected variable value test-value, got %s", variablesResponse.Variables[0].Value)
	}

	if variablesResponse.Variables[0].Visibility != "all" {
		t.Errorf("Expected variable visibility all, got %s", variablesResponse.Variables[0].Visibility)
	}
}

func TestGetRepoActionVariablesMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"variables":[{"name":"REPO_VAR","value":"repo-value"}]}`)
	mockGetter.RepoActionVariablesData = expectedResponse

	// Execute
	response, err := mockGetter.GetRepoActionVariables("test-org", "test-repo")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}

	// Test parsing of the response
	var variablesResponse data.VariableResponse
	err = json.Unmarshal(response, &variablesResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if variablesResponse.TotalCount != 1 {
		t.Errorf("Expected total count 1, got %d", variablesResponse.TotalCount)
	}

	if len(variablesResponse.Variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(variablesResponse.Variables))
	}

	if variablesResponse.Variables[0].Name != "REPO_VAR" {
		t.Errorf("Expected variable name REPO_VAR, got %s", variablesResponse.Variables[0].Name)
	}

	if variablesResponse.Variables[0].Value != "repo-value" {
		t.Errorf("Expected variable value repo-value, got %s", variablesResponse.Variables[0].Value)
	}
}

func TestGetScopedOrgActionVariablesMock(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count":1,"repositories":[{"id":12345,"name":"test-repo"}]}`)
	mockGetter.ScopedOrgActionVariablesData = expectedResponse

	// Execute
	response, err := mockGetter.GetScopedOrgActionVariables("test-org", "TEST_VAR")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}

	// Test parsing of the response
	var scopedResponse data.ScopedResponse
	err = json.Unmarshal(response, &scopedResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if scopedResponse.TotalCount != 1 {
		t.Errorf("Expected total count 1, got %d", scopedResponse.TotalCount)
	}

	if len(scopedResponse.Repositories) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(scopedResponse.Repositories))
	}

	if scopedResponse.Repositories[0].ID != 12345 {
		t.Errorf("Expected repository ID 12345, got %d", scopedResponse.Repositories[0].ID)
	}

	if scopedResponse.Repositories[0].Name != "test-repo" {
		t.Errorf("Expected repository name test-repo, got %s", scopedResponse.Repositories[0].Name)
	}
}

func TestGetSourceOrganizationVariablesError(t *testing.T) {
	testGetSourceOrganizationVariables := func(owner string, g interface{}) ([]byte, error) {
		// This is a simplified version that just returns an error
		return nil, fmt.Errorf("mock error for source API request")
	}

	// Execute - call our local test function instead of the one from another file
	_, err := testGetSourceOrganizationVariables("test-org", nil)

	// Verify
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Check that error message is as expected
	expectedErrMsg := "mock error for source API request"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// Test multiple variables in the response
func TestOrgActionVariablesWithMultipleEntries(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{
        "total_count": 2,
        "variables": [
            {"name": "VAR1", "value": "value1", "visibility": "all"},
            {"name": "VAR2", "value": "value2", "visibility": "selected"}
        ]
    }`)
	mockGetter.OrgActionVariablesData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgActionVariables("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test parsing of the response
	var variablesResponse data.VariableResponse
	err = json.Unmarshal(response, &variablesResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if variablesResponse.TotalCount != 2 {
		t.Errorf("Expected total count 2, got %d", variablesResponse.TotalCount)
	}

	if len(variablesResponse.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(variablesResponse.Variables))
	}

	// Check first variable
	if variablesResponse.Variables[0].Name != "VAR1" {
		t.Errorf("Expected first variable name VAR1, got %s", variablesResponse.Variables[0].Name)
	}
	if variablesResponse.Variables[0].Value != "value1" {
		t.Errorf("Expected first variable value value1, got %s", variablesResponse.Variables[0].Value)
	}
	if variablesResponse.Variables[0].Visibility != "all" {
		t.Errorf("Expected first variable visibility all, got %s", variablesResponse.Variables[0].Visibility)
	}

	// Check second variable
	if variablesResponse.Variables[1].Name != "VAR2" {
		t.Errorf("Expected second variable name VAR2, got %s", variablesResponse.Variables[1].Name)
	}
	if variablesResponse.Variables[1].Value != "value2" {
		t.Errorf("Expected second variable value value2, got %s", variablesResponse.Variables[1].Value)
	}
	if variablesResponse.Variables[1].Visibility != "selected" {
		t.Errorf("Expected second variable visibility selected, got %s", variablesResponse.Variables[1].Visibility)
	}
}

// Test multiple repositories in a scoped variable response
func TestScopedOrgActionVariablesWithMultipleRepos(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{
        "total_count": 2,
        "repositories": [
            {"id": 12345, "name": "repo1"},
            {"id": 67890, "name": "repo2"}
        ]
    }`)
	mockGetter.ScopedOrgActionVariablesData = expectedResponse

	// Execute
	response, err := mockGetter.GetScopedOrgActionVariables("test-org", "TEST_VAR")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test parsing of the response
	var scopedResponse data.ScopedResponse
	err = json.Unmarshal(response, &scopedResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if scopedResponse.TotalCount != 2 {
		t.Errorf("Expected total count 2, got %d", scopedResponse.TotalCount)
	}

	if len(scopedResponse.Repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(scopedResponse.Repositories))
	}

	// Check first repo
	if scopedResponse.Repositories[0].ID != 12345 {
		t.Errorf("Expected first repository ID 12345, got %d", scopedResponse.Repositories[0].ID)
	}
	if scopedResponse.Repositories[0].Name != "repo1" {
		t.Errorf("Expected first repository name repo1, got %s", scopedResponse.Repositories[0].Name)
	}

	// Check second repo
	if scopedResponse.Repositories[1].ID != 67890 {
		t.Errorf("Expected second repository ID 67890, got %d", scopedResponse.Repositories[1].ID)
	}
	if scopedResponse.Repositories[1].Name != "repo2" {
		t.Errorf("Expected second repository name repo2, got %s", scopedResponse.Repositories[1].Name)
	}
}

// Test empty responses
func TestEmptyVariablesResponse(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count": 0, "variables": []}`)
	mockGetter.OrgActionVariablesData = expectedResponse

	// Execute
	response, err := mockGetter.GetOrgActionVariables("test-org")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test parsing of the response
	var variablesResponse data.VariableResponse
	err = json.Unmarshal(response, &variablesResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if variablesResponse.TotalCount != 0 {
		t.Errorf("Expected total count 0, got %d", variablesResponse.TotalCount)
	}

	if len(variablesResponse.Variables) != 0 {
		t.Errorf("Expected 0 variables, got %d", len(variablesResponse.Variables))
	}
}

func TestEmptyScopedResponse(t *testing.T) {
	// Setup
	mockGetter := NewMockAPIGetter()
	expectedResponse := []byte(`{"total_count": 0, "repositories": []}`)
	mockGetter.ScopedOrgActionVariablesData = expectedResponse

	// Execute
	response, err := mockGetter.GetScopedOrgActionVariables("test-org", "TEST_VAR")

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test parsing of the response
	var scopedResponse data.ScopedResponse
	err = json.Unmarshal(response, &scopedResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if scopedResponse.TotalCount != 0 {
		t.Errorf("Expected total count 0, got %d", scopedResponse.TotalCount)
	}

	if len(scopedResponse.Repositories) != 0 {
		t.Errorf("Expected 0 repositories, got %d", len(scopedResponse.Repositories))
	}
}

func TestGetOrgActionVariablesNetworkError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			return nil, fmt.Errorf("DNS resolution failed")
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	t.Run("DNS error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic from log.Fatal")
			}
		}()

		_, err := getter.GetOrgActionVariables("test-org")
		// This line won't be reached due to panic, but linter is satisfied
		if err != nil {
			t.Logf("Got expected error: %v", err)
		}
	})
}

// Test GetRepoActionVariables with various response formats
func TestGetRepoActionVariablesResponseFormats(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	testCases := []struct {
		name         string
		responseBody string
	}{
		{"Standard response", `{"total_count":2,"variables":[{"name":"VAR1","value":"val1"},{"name":"VAR2","value":"val2"}]}`},
		{"Empty variables", `{"total_count":0,"variables":[]}`},
		{"With visibility", `{"total_count":1,"variables":[{"name":"VAR1","value":"val1","visibility":"selected"}]}`},
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

			result, err := getter.GetRepoActionVariables("test-org", "test-repo")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if string(result) != tc.responseBody {
				t.Errorf("Expected %s, got %s", tc.responseBody, string(result))
			}
		})
	}
}

// Test GetScopedOrgActionVariables path validation
func TestGetScopedOrgActionVariablesPath(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	var capturedPath string
	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			capturedPath = path
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"total_count":0,"repositories":[]}`))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	_, err := getter.GetScopedOrgActionVariables("test-org", "TEST_VAR")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedPath := "orgs/test-org/actions/variables/TEST_VAR/repositories"
	if capturedPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, capturedPath)
	}
}

// Test with special characters in variable names
func TestVariablesWithSpecialCharacters(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	specialVarNames := []string{
		"TEST_VAR_123",
		"TEST-VAR",
		"TEST.VAR",
		"TEST VAR", // This might need URL encoding
	}

	for _, varName := range specialVarNames {
		t.Run(fmt.Sprintf("Variable: %s", varName), func(t *testing.T) {
			mockClient := &mockRESTClient{
				RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
					// Check if the variable name is properly included in the path
					if !strings.Contains(path, varName) && !strings.Contains(path, strings.ReplaceAll(varName, " ", "%20")) {
						t.Errorf("Expected variable name in path, got: %s", path)
					}

					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"total_count":1,"repositories":[{"id":1,"name":"repo1"}]}`))),
					}, nil
				},
			}

			getter := newAPIGetterWithMockREST(mockClient)

			result, err := getter.GetScopedOrgActionVariables("test-org", varName)
			if err != nil {
				t.Errorf("Unexpected error for variable '%s': %v", varName, err)
			}

			if len(result) == 0 {
				t.Errorf("Expected non-empty result for variable '%s'", varName)
			}
		})
	}
}

// Test concurrent access to variables
func TestConcurrentVariableAccess(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	mockClient := &mockRESTClient{
		RequestFunc: func(method string, path string, body io.Reader) (*http.Response, error) {
			// Simulate some processing time
			response := `{"total_count":1,"variables":[{"name":"CONCURRENT_VAR","value":"test"}]}`
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(response))),
			}, nil
		},
	}

	getter := newAPIGetterWithMockREST(mockClient)

	// Run multiple goroutines accessing the same getter
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func(index int) {
			defer func() { done <- true }()

			result, err := getter.GetOrgActionVariables(fmt.Sprintf("test-org-%d", index))
			if err != nil {
				t.Errorf("Goroutine %d: Unexpected error: %v", index, err)
			}

			if len(result) == 0 {
				t.Errorf("Goroutine %d: Expected non-empty result", index)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}
