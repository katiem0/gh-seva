package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"

	data "github.com/katiem0/gh-seva/internal/data"
	"golang.org/x/crypto/nacl/box"
)

func (g *APIGetter) CreateSecretsList(filedata [][]string) []data.ImportedSecret {
	// convert csv lines to array of structs
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

func (g *APIGetter) EncryptSecret(publicKey string, secret string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	var decodedKey [32]byte
	copy(decodedKey[:], bytes)

	encrypted, err := box.SealAnonymous(nil, []byte(secret), (*[32]byte)(bytes), rand.Reader)
	if err != nil {
		return "", err
	}
	// Encode the encrypted value in base64
	encryptedValue := base64.StdEncoding.EncodeToString(encrypted)

	return encryptedValue, nil
}

func CreateSelectedOrgSecretData(secret data.ImportedSecret, keyID string, encryptedValue string) *data.CreateOrgSecret {
	secretArray := make([]int, len(secret.RepositoryIDs))
	for i := range secretArray {
		secretArray[i], _ = strconv.Atoi(secret.RepositoryIDs[i])
	}
	s := data.CreateOrgSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     secret.Access,
		SelectedRepos:  secretArray,
	}
	return &s
}

func CreateOrgSecretData(secret data.ImportedSecret, keyID string, encryptedValue string) *data.CreateOrgSecretAll {
	s := data.CreateOrgSecretAll{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     secret.Access,
	}
	return &s
}

// Separate function to address that the Dependabot Org Secret API
// is an array of strings instead of an array of integers
func CreateOrgDependabotSecretData(secret data.ImportedSecret, keyID string, encryptedValue string) *data.CreateOrgDepSecret {
	s := data.CreateOrgDepSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
		Visibility:     secret.Access,
		SelectedRepos:  secret.RepositoryIDs,
	}
	return &s
}

func CreateRepoSecretData(keyID string, encryptedValue string) *data.CreateRepoSecret {
	s := data.CreateRepoSecret{
		EncryptedValue: encryptedValue,
		KeyID:          keyID,
	}
	return &s
}
