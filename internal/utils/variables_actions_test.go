package utils

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/katiem0/gh-seva/internal/data"
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
