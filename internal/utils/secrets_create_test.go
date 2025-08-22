package utils

import (
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/katiem0/gh-seva/internal/data"
)

func TestCreateSecretsList(t *testing.T) {
	// Setup
	g := &APIGetter{}
	filedata := [][]string{
		{"SecretLevel", "SecretType", "SecretName", "SecretValue", "SecretAccess", "RepositoryNames", "RepositoryIDs"},
		{"Organization", "Actions", "TEST_SECRET1", "secret1", "all", "", ""},
		{"Organization", "Dependabot", "TEST_SECRET2", "secret2", "selected", "repo1;repo2", "1234;5678"},
		{"Repository", "Codespaces", "TEST_SECRET3", "secret3", "RepoOnly", "repo3", "9101"},
	}

	// Execute
	result := g.CreateSecretsList(filedata)

	// Verify
	expected := []data.ImportedSecret{
		{
			Level:           "Organization",
			Type:            "Actions",
			Name:            "TEST_SECRET1",
			Value:           "secret1",
			Access:          "all",
			RepositoryNames: []string{""},
			RepositoryIDs:   []string{""},
		},
		{
			Level:           "Organization",
			Type:            "Dependabot",
			Name:            "TEST_SECRET2",
			Value:           "secret2",
			Access:          "selected",
			RepositoryNames: []string{"repo1", "repo2"},
			RepositoryIDs:   []string{"1234", "5678"},
		},
		{
			Level:           "Repository",
			Type:            "Codespaces",
			Name:            "TEST_SECRET3",
			Value:           "secret3",
			Access:          "RepoOnly",
			RepositoryNames: []string{"repo3"},
			RepositoryIDs:   []string{"9101"},
		},
	}

	if len(result) != len(expected) {
		t.Errorf("Expected %d secrets, got %d", len(expected), len(result))
		return
	}

	for i, e := range expected {
		if result[i].Level != e.Level {
			t.Errorf("Secret %d: Expected Level %s, got %s", i, e.Level, result[i].Level)
		}
		if result[i].Type != e.Type {
			t.Errorf("Secret %d: Expected Type %s, got %s", i, e.Type, result[i].Type)
		}
		if result[i].Name != e.Name {
			t.Errorf("Secret %d: Expected Name %s, got %s", i, e.Name, result[i].Name)
		}
		if result[i].Value != e.Value {
			t.Errorf("Secret %d: Expected Value %s, got %s", i, e.Value, result[i].Value)
		}
		if result[i].Access != e.Access {
			t.Errorf("Secret %d: Expected Access %s, got %s", i, e.Access, result[i].Access)
		}
	}
}

func TestEncryptSecret(t *testing.T) {
	g := &APIGetter{}

	// Create a valid base64 encoded key
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}
	publicKey := base64.StdEncoding.EncodeToString(testKey)

	secretValue := "test-secret-value"

	// Execute
	encryptedSecret, err := g.EncryptSecret(publicKey, secretValue)

	// Verify
	if err != nil {
		t.Errorf("EncryptSecret() error = %v", err)
	}

	if encryptedSecret == "" {
		t.Error("EncryptSecret() returned empty string")
	}

	// Verify it's valid base64
	_, err = base64.StdEncoding.DecodeString(encryptedSecret)
	if err != nil {
		t.Errorf("EncryptSecret() returned invalid base64: %v", err)
	}
}

func TestCreateSelectedOrgSecretData(t *testing.T) {
	// Setup
	secret := data.ImportedSecret{
		Name:          "TEST_SECRET",
		Access:        "selected",
		RepositoryIDs: []string{"1234", "5678"},
	}
	keyID := "test-key-id"
	encryptedValue := "encrypted-secret-value"

	// Execute
	result := CreateSelectedOrgSecretData(secret, keyID, encryptedValue)

	// Verify
	expected := &data.CreateOrgSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     "selected",
		SelectedRepos:  []int{1234, 5678},
	}

	if result.EncryptedValue != expected.EncryptedValue {
		t.Errorf("Expected EncryptedValue %s, got %s", expected.EncryptedValue, result.EncryptedValue)
	}

	if result.KeyID != expected.KeyID {
		t.Errorf("Expected KeyID %s, got %s", expected.KeyID, result.KeyID)
	}

	if result.Visibility != expected.Visibility {
		t.Errorf("Expected Visibility %s, got %s", expected.Visibility, result.Visibility)
	}

	if len(result.SelectedRepos) != len(expected.SelectedRepos) {
		t.Errorf("Expected %d SelectedRepos, got %d", len(expected.SelectedRepos), len(result.SelectedRepos))
		return
	}

	for i, id := range expected.SelectedRepos {
		if result.SelectedRepos[i] != id {
			t.Errorf("SelectedRepos[%d]: Expected %d, got %d", i, id, result.SelectedRepos[i])
		}
	}
}

func TestStrToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{
			name:     "Valid integer",
			input:    "123",
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "Zero",
			input:    "0",
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "Negative integer",
			input:    "-45",
			expected: -45,
			wantErr:  false,
		},
		{
			name:     "Invalid integer",
			input:    "abc",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := strToInt(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %s but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %s: %v", tt.input, err)
				}

				if result != tt.expected {
					t.Errorf("Expected %d, got %d", tt.expected, result)
				}
			}
		})
	}
}

func strToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func TestCreateOrgSecretData(t *testing.T) {
	// Setup
	secret := data.ImportedSecret{
		Name:   "TEST_SECRET",
		Access: "all",
		Value:  "secret-value",
	}
	keyID := "test-key-id"
	encryptedValue := "encrypted-secret-value"

	// Execute
	result := CreateOrgSecretData(secret, keyID, encryptedValue)

	// Verify
	expected := &data.CreateOrgSecretAll{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     "all",
	}

	if result.EncryptedValue != expected.EncryptedValue {
		t.Errorf("Expected EncryptedValue %s, got %s", expected.EncryptedValue, result.EncryptedValue)
	}

	if result.KeyID != expected.KeyID {
		t.Errorf("Expected KeyID %s, got %s", expected.KeyID, result.KeyID)
	}

	if result.Visibility != expected.Visibility {
		t.Errorf("Expected Visibility %s, got %s", expected.Visibility, result.Visibility)
	}
}

func TestCreateOrgDependabotSecretData(t *testing.T) {
	// Setup
	secret := data.ImportedSecret{
		Name:          "TEST_SECRET",
		Access:        "selected",
		RepositoryIDs: []string{"1234", "5678"},
	}
	keyID := "test-key-id"
	encryptedValue := "encrypted-secret-value"

	// Execute
	result := CreateOrgDependabotSecretData(secret, keyID, encryptedValue)

	// Verify
	expected := &data.CreateOrgDepSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     "selected",
		SelectedRepos:  []string{"1234", "5678"},
	}

	if result.EncryptedValue != expected.EncryptedValue {
		t.Errorf("Expected EncryptedValue %s, got %s", expected.EncryptedValue, result.EncryptedValue)
	}

	if result.KeyID != expected.KeyID {
		t.Errorf("Expected KeyID %s, got %s", expected.KeyID, result.KeyID)
	}

	if result.Visibility != expected.Visibility {
		t.Errorf("Expected Visibility %s, got %s", expected.Visibility, result.Visibility)
	}

	if len(result.SelectedRepos) != len(expected.SelectedRepos) {
		t.Errorf("Expected %d SelectedRepos, got %d", len(expected.SelectedRepos), len(result.SelectedRepos))
		return
	}

	for i, id := range expected.SelectedRepos {
		if result.SelectedRepos[i] != id {
			t.Errorf("SelectedRepos[%d]: Expected %s, got %s", i, id, result.SelectedRepos[i])
		}
	}
}

func TestCreateRepoSecretData(t *testing.T) {
	// Setup
	keyID := "test-key-id"
	encryptedValue := "encrypted-secret-value"

	// Execute
	result := CreateRepoSecretData(keyID, encryptedValue)

	// Verify
	expected := &data.CreateRepoSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
	}

	if result.EncryptedValue != expected.EncryptedValue {
		t.Errorf("Expected EncryptedValue %s, got %s", expected.EncryptedValue, result.EncryptedValue)
	}

	if result.KeyID != expected.KeyID {
		t.Errorf("Expected KeyID %s, got %s", expected.KeyID, result.KeyID)
	}
}
