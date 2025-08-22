package data

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSecret(t *testing.T) {
	// Test Secret struct
	now := time.Now()
	secret := Secret{
		Name:          "TEST_SECRET",
		CreatedAt:     now,
		UpdatedAt:     now,
		Visibility:    "all",
		SelectedRepos: "https://api.github.com/orgs/testorg/actions/secrets/TEST_SECRET/repositories",
	}

	if secret.Name != "TEST_SECRET" {
		t.Errorf("Expected Name to be 'TEST_SECRET', got %s", secret.Name)
	}

	if secret.Visibility != "all" {
		t.Errorf("Expected Visibility to be 'all', got %s", secret.Visibility)
	}

	if secret.SelectedRepos != "https://api.github.com/orgs/testorg/actions/secrets/TEST_SECRET/repositories" {
		t.Errorf("Expected SelectedRepos URL to match, got %s", secret.SelectedRepos)
	}
}

func TestSecretsResponse(t *testing.T) {
	// Test SecretsResponse struct
	response := SecretsResponse{
		TotalCount: 2,
		Secrets: []Secret{
			{
				Name:       "SECRET1",
				Visibility: "all",
			},
			{
				Name:       "SECRET2",
				Visibility: "selected",
			},
		},
	}

	if response.TotalCount != 2 {
		t.Errorf("Expected TotalCount to be 2, got %d", response.TotalCount)
	}

	if len(response.Secrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(response.Secrets))
	}

	if response.Secrets[0].Name != "SECRET1" {
		t.Errorf("Expected first secret name to be 'SECRET1', got %s", response.Secrets[0].Name)
	}

	if response.Secrets[1].Visibility != "selected" {
		t.Errorf("Expected second secret visibility to be 'selected', got %s", response.Secrets[1].Visibility)
	}
}

func TestImportedSecret(t *testing.T) {
	// Test ImportedSecret struct
	secret := ImportedSecret{
		Level:           "Organization",
		Type:            "Actions",
		Name:            "TEST_SECRET",
		Value:           "secret-value",
		Access:          "selected",
		RepositoryNames: []string{"repo1", "repo2"},
		RepositoryIDs:   []string{"1", "2"},
	}

	if secret.Level != "Organization" {
		t.Errorf("Expected Level to be 'Organization', got %s", secret.Level)
	}

	if secret.Type != "Actions" {
		t.Errorf("Expected Type to be 'Actions', got %s", secret.Type)
	}

	if secret.Name != "TEST_SECRET" {
		t.Errorf("Expected Name to be 'TEST_SECRET', got %s", secret.Name)
	}

	if secret.Value != "secret-value" {
		t.Errorf("Expected Value to be 'secret-value', got %s", secret.Value)
	}

	if secret.Access != "selected" {
		t.Errorf("Expected Access to be 'selected', got %s", secret.Access)
	}

	if len(secret.RepositoryNames) != 2 {
		t.Errorf("Expected 2 repository names, got %d", len(secret.RepositoryNames))
	}

	if len(secret.RepositoryIDs) != 2 {
		t.Errorf("Expected 2 repository IDs, got %d", len(secret.RepositoryIDs))
	}
}

func TestPublicKey(t *testing.T) {
	// Test PublicKey struct
	key := PublicKey{
		KeyID: "key-id-123",
		Key:   "base64-encoded-key",
	}

	if key.KeyID != "key-id-123" {
		t.Errorf("Expected KeyID to be 'key-id-123', got %s", key.KeyID)
	}

	if key.Key != "base64-encoded-key" {
		t.Errorf("Expected Key to be 'base64-encoded-key', got %s", key.Key)
	}
}

func TestSecretStructsJSON(t *testing.T) {
	// Test JSON marshaling/unmarshaling

	// Test CreateOrgSecret
	orgSecret := CreateOrgSecret{
		EncryptedValue: "encrypted-value",
		KeyID:          "key-id",
		Visibility:     "selected",
		SelectedRepos:  []int{1, 2, 3},
	}

	jsonBytes, err := json.Marshal(orgSecret)
	if err != nil {
		t.Fatalf("Failed to marshal CreateOrgSecret: %v", err)
	}

	var unmarshaledOrgSecret CreateOrgSecret
	err = json.Unmarshal(jsonBytes, &unmarshaledOrgSecret)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateOrgSecret: %v", err)
	}

	if unmarshaledOrgSecret.KeyID != "key-id" {
		t.Errorf("Expected KeyID to be 'key-id', got %s", unmarshaledOrgSecret.KeyID)
	}

	if len(unmarshaledOrgSecret.SelectedRepos) != 3 {
		t.Errorf("Expected 3 SelectedRepos, got %d", len(unmarshaledOrgSecret.SelectedRepos))
	}

	// Test CreateOrgSecretAll
	orgSecretAll := CreateOrgSecretAll{
		EncryptedValue: "encrypted-value",
		KeyID:          "key-id",
		Visibility:     "all",
	}

	jsonBytes, err = json.Marshal(orgSecretAll)
	if err != nil {
		t.Fatalf("Failed to marshal CreateOrgSecretAll: %v", err)
	}

	var unmarshaledOrgSecretAll CreateOrgSecretAll
	err = json.Unmarshal(jsonBytes, &unmarshaledOrgSecretAll)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateOrgSecretAll: %v", err)
	}

	if unmarshaledOrgSecretAll.Visibility != "all" {
		t.Errorf("Expected Visibility to be 'all', got %s", unmarshaledOrgSecretAll.Visibility)
	}
}
