package data

import "time"

type CreateOrgSecret struct {
	EncryptedValue string `json:"encrypted_value"`
	KeyID          string `json:"key_id"`
	Visibility     string `json:"visibility"`
	SelectedRepos  []int  `json:"selected_repository_ids"`
}

type CreateOrgSecretAll struct {
	EncryptedValue string `json:"encrypted_value"`
	KeyID          string `json:"key_id"`
	Visibility     string `json:"visibility"`
}

// Address Dependabot API differences
type CreateOrgDepSecret struct {
	EncryptedValue string   `json:"encrypted_value"`
	KeyID          string   `json:"key_id"`
	Visibility     string   `json:"visibility"`
	SelectedRepos  []string `json:"selected_repository_ids"`
}

type CreateRepoSecret struct {
	EncryptedValue string `json:"encrypted_value"`
	KeyID          string `json:"key_id"`
}

type ImportedSecret struct {
	Level           string   `json:"level"`
	Type            string   `json:"type"`
	Name            string   `json:"name"`
	Value           string   `json:"value"`
	Access          string   `json:"visibility"`
	RepositoryNames []string `json:"selected_repositories"`
	RepositoryIDs   []string `json:"selected_repository_ids"`
}

type PublicKey struct {
	KeyID string `json:"key_id"`
	Key   string `json:"key"`
}

type Secret struct {
	Name          string    `json:"name"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Visibility    string    `json:"visibility"`
	SelectedRepos string    `json:"selected_repositories_url"`
}

type SecretExport struct {
	SecretLevel    string
	SecretType     string
	SecretName     string
	SecretAccess   string
	RepositoryName string
	RepositoryID   int
}

type SecretsResponse struct {
	TotalCount int      `json:"total_count"`
	Secrets    []Secret `json:"secrets"`
}
