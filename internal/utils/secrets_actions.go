package utils

import (
	"fmt"
	"io"
	"log"

	"go.uber.org/zap"
)

func (g *APIGetter) GetOrgActionSecrets(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/secrets", owner)

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

func (g *APIGetter) GetRepoActionSecrets(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/secrets", owner, repo)

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

func (g *APIGetter) GetScopedOrgActionSecrets(owner string, secret string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/secrets/%s/repositories", owner, secret)

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

func (g *APIGetter) GetOrgActionPublicKey(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/secrets/public-key", owner)
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

func (g *APIGetter) GetRepoActionPublicKey(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/secrets/public-key", owner, repo)
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

func (g *APIGetter) CreateOrgActionSecret(owner string, secret string, data io.Reader) error {
	url := fmt.Sprintf("orgs/%s/actions/secrets/%s", owner, secret)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func (g *APIGetter) CreateRepoActionSecret(owner string, repo string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/actions/secrets/%s", owner, repo, secret)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}
