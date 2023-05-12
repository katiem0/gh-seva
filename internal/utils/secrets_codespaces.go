package utils

import (
	"fmt"
	"io"
	"log"

	"go.uber.org/zap"
)

func (g *APIGetter) GetOrgCodespacesSecrets(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets", owner)

	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, err
}

func (g *APIGetter) GetRepoCodespacesSecrets(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/codespaces/secrets", owner, repo)

	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, err
}

func (g *APIGetter) GetScopedOrgCodespacesSecrets(owner string, secret string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets/%s/repositories", owner, secret)

	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData, err
}

func (g *APIGetter) GetOrgCodespacesPublicKey(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets/public-key", owner)
	zap.S().Debugf("Getting public-key for %v", url)
	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	return responseData, err
}

func (g *APIGetter) GetRepoCodespacesPublicKey(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/codespaces/secrets/public-key", owner, repo)
	zap.S().Debugf("Getting public-key for %v", url)
	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	return responseData, err
}

func (g *APIGetter) CreateOrgCodespacesSecret(owner string, secret string, data io.Reader) error {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets/%s", owner, secret)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func (g *APIGetter) CreateRepoCodespacesSecret(owner string, repo string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/codespaces/secrets/%s", owner, repo, secret)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}
