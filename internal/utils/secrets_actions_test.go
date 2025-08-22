package utils

import (
	"testing"
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
