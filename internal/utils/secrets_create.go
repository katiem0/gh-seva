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

func CreateOrgSecretData(secret data.ImportedSecret, keyID string, encryptedValue string) *data.CreateOrgSecret {
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
