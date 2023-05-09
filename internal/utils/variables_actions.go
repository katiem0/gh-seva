package utils

import (
	"fmt"
	"io"
	"log"
)

func (g *APIGetter) GetOrgActionVariables(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/variables", owner)

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

func (g *APIGetter) GetRepoActionVariables(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/variables", owner, repo)

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

func (g *APIGetter) GetScopedOrgActionVariables(owner string, secret string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/variables/%s/repositories", owner, secret)

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
