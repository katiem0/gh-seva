package data

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"

	"golang.org/x/crypto/nacl/box"
)

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

type CreateOrgSecret struct {
	EncryptedValue string `json:"encrypted_value"`
	KeyID          string `json:"key_id"`
	Visibility     string `json:"visibility"`
	SelectedRepos  []int  `json:"selected_repository_ids"`
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

func (g *APIGetter) CreateSecretsList(data [][]string) []ImportedSecret {
	// convert csv lines to array of structs
	var importSecretList []ImportedSecret
	var secret ImportedSecret
	for _, each := range data[1:] {
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

func (g *APIGetter) EncryptSecret(publickey string, secret string) (string, error) {
	var pkBytes [32]byte
	copy(pkBytes[:], publickey)
	secretBytes := secret

	out := make([]byte, 0,
		len(secretBytes)+
			box.Overhead+
			len(pkBytes))

	enc, err := box.SealAnonymous(
		out, []byte(secretBytes), &pkBytes, rand.Reader,
	)
	if err != nil {
		return "", err
	}

	encEnc := base64.StdEncoding.EncodeToString(enc)

	return encEnc, nil
}

func CreateOrgSecretData(secret ImportedSecret, keyID string, encryptedValue string) *CreateOrgSecret {
	secretArray := make([]int, len(secret.RepositoryIDs))
	for i := range secretArray {
		secretArray[i], _ = strconv.Atoi(secret.RepositoryIDs[i])
	}
	s := CreateOrgSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     secret.Access,
		SelectedRepos:  secretArray,
	}
	return &s
}

// Separate function to address that the Dependabot Org Secret API
// is an array of strings instead of an array of integers
func CreateOrgDependabotSecretData(secret ImportedSecret, keyID string, encryptedValue string) *CreateOrgDepSecret {
	s := CreateOrgDepSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     secret.Access,
		SelectedRepos:  secret.RepositoryIDs,
	}
	return &s
}

func CreateRepoSecretData(keyID string, encryptedValue string) *CreateRepoSecret {
	s := CreateRepoSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
	}
	return &s
}
