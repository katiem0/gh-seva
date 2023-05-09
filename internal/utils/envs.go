package utils

import (
	"fmt"
	"io"
	"log"
)

func (g *APIGetter) GetRepoEnvironments(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/environments", owner, repo)
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
