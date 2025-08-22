package utils

import (
	"testing"
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
