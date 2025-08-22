package utils

import (
	"fmt"
	"io"
	"strings"

	"github.com/katiem0/gh-seva/internal/data"
)

// MockAPIGetter implements the Getter interface for testing
type MockAPIGetter struct {
	// Mocked responses
	ReposResponse                  *data.ReposQuery
	RepoResponse                   *data.RepoSingleQuery
	OrgActionSecretsData           []byte
	RepoActionSecretsData          []byte
	ScopedOrgActionSecretsData     []byte
	OrgDependabotSecretsData       []byte
	RepoDependabotSecretsData      []byte
	ScopedOrgDependabotSecretsData []byte
	OrgCodespacesSecretsData       []byte
	RepoCodespacesSecretsData      []byte
	ScopedOrgCodespacesSecretsData []byte
	OrgActionVariablesData         []byte
	RepoActionVariablesData        []byte
	ScopedOrgActionVariablesData   []byte
	PublicKeyData                  []byte
	EncryptedSecret                string
	ImportedSecrets                []data.ImportedSecret
	GetRepoError                   bool
	GetReposListError              bool
	GetOrgActionSecretsError       bool
	ShouldReturnError              bool
}

// NewMockAPIGetter creates a new instance of MockAPIGetter
func NewMockAPIGetter() *MockAPIGetter {
	return &MockAPIGetter{}
}

// GetReposList mocks the retrieval of repositories list
func (m *MockAPIGetter) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	return m.ReposResponse, nil
}

// GetRepo mocks the retrieval of a single repository
func (m *MockAPIGetter) GetRepo(owner string, name string) (*data.RepoSingleQuery, error) {
	return m.RepoResponse, nil
}

// GetOrgActionSecrets mocks retrieving organization action secrets
func (m *MockAPIGetter) GetOrgActionSecrets(owner string) ([]byte, error) {
	return m.OrgActionSecretsData, nil
}

// GetRepoActionSecrets mocks retrieving repository action secrets
func (m *MockAPIGetter) GetRepoActionSecrets(owner string, repo string) ([]byte, error) {
	return m.RepoActionSecretsData, nil
}

// GetScopedOrgActionSecrets mocks retrieving scoped organization action secrets
func (m *MockAPIGetter) GetScopedOrgActionSecrets(owner string, secret string) ([]byte, error) {
	return m.ScopedOrgActionSecretsData, nil
}

// GetOrgDependabotSecrets mocks retrieving organization dependabot secrets
func (m *MockAPIGetter) GetOrgDependabotSecrets(owner string) ([]byte, error) {
	return m.OrgDependabotSecretsData, nil
}

// GetRepoDependabotSecrets mocks retrieving repository dependabot secrets
func (m *MockAPIGetter) GetRepoDependabotSecrets(owner string, repo string) ([]byte, error) {
	return m.RepoDependabotSecretsData, nil
}

// GetScopedOrgDependabotSecrets mocks retrieving scoped organization dependabot secrets
func (m *MockAPIGetter) GetScopedOrgDependabotSecrets(owner string, secret string) ([]byte, error) {
	return m.ScopedOrgDependabotSecretsData, nil
}

// GetOrgCodespacesSecrets mocks retrieving organization codespaces secrets
func (m *MockAPIGetter) GetOrgCodespacesSecrets(owner string) ([]byte, error) {
	return m.OrgCodespacesSecretsData, nil
}

// GetRepoCodespacesSecrets mocks retrieving repository codespaces secrets
func (m *MockAPIGetter) GetRepoCodespacesSecrets(owner string, repo string) ([]byte, error) {
	return m.RepoCodespacesSecretsData, nil
}

// GetScopedOrgCodespacesSecrets mocks retrieving scoped organization codespaces secrets
func (m *MockAPIGetter) GetScopedOrgCodespacesSecrets(owner string, secret string) ([]byte, error) {
	return m.ScopedOrgCodespacesSecretsData, nil
}

// CreateSecretsList mocks creating a list of secrets from CSV data
func (m *MockAPIGetter) CreateSecretsList(filedata [][]string) []data.ImportedSecret {
	if m.ImportedSecrets != nil {
		return m.ImportedSecrets
	}

	// Provide a default implementation if no mock data is set
	var importSecretList []data.ImportedSecret
	var secret data.ImportedSecret
	for _, each := range filedata[1:] {
		secret.Level = each[0]
		secret.Type = each[1]
		secret.Name = each[2]
		secret.Value = each[3]
		secret.Access = each[4]
		secret.RepositoryNames = strings.Split(each[5], ";")
		secret.RepositoryIDs = strings.Split(each[6], ";")
		importSecretList = append(importSecretList, secret)
	}
	return importSecretList
}

// GetOrgActionPublicKey mocks retrieving organization action public key
func (m *MockAPIGetter) GetOrgActionPublicKey(owner string) ([]byte, error) {
	return m.PublicKeyData, nil
}

// GetRepoActionPublicKey mocks retrieving repository action public key
func (m *MockAPIGetter) GetRepoActionPublicKey(owner string, repo string) ([]byte, error) {
	return m.PublicKeyData, nil
}

// GetOrgCodespacesPublicKey mocks retrieving organization codespaces public key
func (m *MockAPIGetter) GetOrgCodespacesPublicKey(owner string) ([]byte, error) {
	return m.PublicKeyData, nil
}

// GetRepoCodespacesPublicKey mocks retrieving repository codespaces public key
func (m *MockAPIGetter) GetRepoCodespacesPublicKey(owner string, repo string) ([]byte, error) {
	return m.PublicKeyData, nil
}

// GetOrgDependabotPublicKey mocks retrieving organization dependabot public key
func (m *MockAPIGetter) GetOrgDependabotPublicKey(owner string) ([]byte, error) {
	return m.PublicKeyData, nil
}

// GetRepoDependabotPublicKey mocks retrieving repository dependabot public key
func (m *MockAPIGetter) GetRepoDependabotPublicKey(owner string, repo string) ([]byte, error) {
	return m.PublicKeyData, nil
}

// EncryptSecret mocks encrypting a secret
func (m *MockAPIGetter) EncryptSecret(publickey string, secret string) (string, error) {
	return m.EncryptedSecret, nil
}

// CreateOrgActionSecret mocks creating an organization action secret
func (m *MockAPIGetter) CreateOrgActionSecret(owner string, secret string, data io.Reader) error {
	return nil
}

// CreateRepoActionSecret mocks creating a repository action secret
func (m *MockAPIGetter) CreateRepoActionSecret(owner string, repo string, secret string, data io.Reader) error {
	return nil
}

// CreateOrgCodespacesSecret mocks creating an organization codespaces secret
func (m *MockAPIGetter) CreateOrgCodespacesSecret(owner string, secret string, data io.Reader) error {
	return nil
}

// CreateRepoCodespacesSecret mocks creating a repository codespaces secret
func (m *MockAPIGetter) CreateRepoCodespacesSecret(owner string, repo string, secret string, data io.Reader) error {
	return nil
}

// CreateOrgDependabotSecret mocks creating an organization dependabot secret
func (m *MockAPIGetter) CreateOrgDependabotSecret(owner string, secret string, data io.Reader) error {
	return nil
}

// CreateRepoDependabotSecret mocks creating a repository dependabot secret
func (m *MockAPIGetter) CreateRepoDependabotSecret(owner string, repo string, secret string, data io.Reader) error {
	return nil
}

// GetOrgActionVariables mocks retrieving organization action variables
func (m *MockAPIGetter) GetOrgActionVariables(owner string) ([]byte, error) {
	return m.OrgActionVariablesData, nil
}

// GetRepoActionVariables mocks retrieving repository action variables
func (m *MockAPIGetter) GetRepoActionVariables(owner string, repo string) ([]byte, error) {
	return m.RepoActionVariablesData, nil
}

// GetScopedOrgActionVariables mocks retrieving scoped organization action variables
func (m *MockAPIGetter) GetScopedOrgActionVariables(owner string, variable string) ([]byte, error) {
	return m.ScopedOrgActionVariablesData, nil
}

// CreateVariableList mocks creating a list of variables from CSV data
func (m *MockAPIGetter) CreateVariableList(filedata [][]string) []data.ImportedVariable {
	var variableList []data.ImportedVariable
	var vars data.ImportedVariable
	for _, each := range filedata[1:] {
		vars.Level = each[0]
		vars.Name = each[1]
		vars.Value = each[2]
		vars.Visibility = each[3]
		vars.SelectedRepos = strings.Split(each[4], ";")
		vars.SelectedReposIDs = strings.Split(each[5], ";")
		variableList = append(variableList, vars)
	}
	return variableList
}

func (m *MockAPIGetter) CreateOrganizationVariable(owner string, data io.Reader) error {
	if m.ShouldReturnError {
		return fmt.Errorf("mock error for CreateOrganizationVariable")
	}
	return nil
}

func (m *MockAPIGetter) CreateRepoVariable(owner string, repo string, data io.Reader) error {
	if m.ShouldReturnError {
		return fmt.Errorf("mock error for CreateRepoVariable")
	}
	return nil
}
