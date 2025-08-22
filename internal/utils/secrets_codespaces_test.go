package utils

import (
	"testing"
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
