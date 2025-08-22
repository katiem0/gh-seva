package utils

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/katiem0/gh-seva/internal/data"
	"go.uber.org/zap"
)

func (g *APIGetter) CreateVariableList(filedata [][]string) []data.ImportedVariable {
	var variableList []data.ImportedVariable

	if len(filedata) <= 1 {
		zap.S().Warn("Empty variable data provided")
		return variableList
	}

	for _, each := range filedata[1:] {
		// Skip if we don't have at least name, level, and value
		if len(each) < 3 {
			zap.S().Warn("Skipping row with insufficient fields")
			continue
		}

		variable := data.ImportedVariable{
			Level: each[0],
			Name:  each[1],
			Value: each[2],
		}

		// Validate required fields
		if variable.Name == "" {
			zap.S().Warn("Skipping variable with empty name")
			continue
		}

		// Handle optional fields
		if len(each) > 3 {
			variable.Visibility = each[3]
		}

		// Handle repository data
		if len(each) > 4 {
			variable.SelectedRepos = strings.Split(each[4], ";")
		}

		if len(each) > 5 {
			variable.SelectedReposIDs = strings.Split(each[5], ";")
		}

		zap.S().Debugf("Processed variable: %s/%s", variable.Level, variable.Name)
		variableList = append(variableList, variable)
	}

	return variableList
}

func GetSourceOrganizationVariables(owner string, g *sourceAPIGetter) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/variables", owner)
	zap.S().Debugf("Reading in variables from %v", url)
	resp, err := g.restClient.Request("GET", url, nil)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			zap.S().Errorf("Error closing response body: %v", err)
		}
	}()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Body read error, %v", err)
	}
	return responseData, err
}

func GetScopedSourceOrgActionVariables(owner string, secret string, g *sourceAPIGetter) ([]byte, error) {
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

func (g *APIGetter) CreateOrganizationVariable(owner string, data io.Reader) error {
	url := fmt.Sprintf("orgs/%s/actions/variables", owner)

	resp, err := g.restClient.Request("POST", url, data)
	if err != nil {
		zap.S().Errorf("Error making request to create organization variable: %v", err)
		return fmt.Errorf("failed to create organization variable: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			zap.S().Errorf("Error closing response body: %v", err)
		}
	}()

	// Check response status code
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	return nil
}

func (g *APIGetter) CreateRepoVariable(owner string, repo string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/actions/variables", owner, repo)

	resp, err := g.restClient.Request("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			zap.S().Errorf("Error closing response body: %v", err)
		}
	}()
	return err
}

func CreateSelectedOrgVariableData(variable data.ImportedVariable) *data.CreateOrgVariable {
	var validIDs []int

	for _, idStr := range variable.SelectedReposIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			zap.S().Warnf("Invalid repository ID '%s' will be skipped", idStr)
			continue
		}
		validIDs = append(validIDs, id)
	}

	s := data.CreateOrgVariable{
		Name:             variable.Name,
		Value:            variable.Value,
		Visibility:       variable.Visibility,
		SelectedReposIDs: validIDs,
	}
	return &s
}

func CreateOrgVariableData(variable data.ImportedVariable) *data.CreateVariableAll {
	s := data.CreateVariableAll{
		Name:       variable.Name,
		Value:      variable.Value,
		Visibility: variable.Visibility,
	}
	return &s
}

func CreateRepoVariableData(variable data.ImportedVariable) *data.CreateRepoVariable {
	s := data.CreateRepoVariable{
		Name:  variable.Name,
		Value: variable.Value,
	}
	return &s
}

func CreateOrgSourceVariableData(variable data.Variable) *data.CreateVariableAll {
	s := data.CreateVariableAll{
		Name:       variable.Name,
		Value:      variable.Value,
		Visibility: variable.Visibility,
	}
	return &s
}
