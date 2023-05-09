package utils

import (
	"fmt"
	"io"
	"log"

	"go.uber.org/zap"
)

func (g *APIGetter) GetOrgDependabotSecrets(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets", owner)

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

func (g *APIGetter) GetRepoDependabotSecrets(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/dependabot/secrets", owner, repo)

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

func (g *APIGetter) GetScopedOrgDependabotSecrets(owner string, secret string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets/%s/repositories", owner, secret)

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

func (g *APIGetter) GetOrgDependabotPublicKey(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets/public-key", owner)
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

func (g *APIGetter) GetRepoDependabotPublicKey(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/dependabot/secrets/public-key", owner, repo)
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

func (g *APIGetter) CreateOrgDependabotSecret(owner string, secret string, data io.Reader) error {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets/%s", owner, secret)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func (g *APIGetter) CreateRepoDependabotSecret(owner string, repo string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/dependabot/secrets/%s", owner, repo, secret)

	resp, err := g.restClient.Request("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}
