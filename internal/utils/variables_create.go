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
	// convert csv lines to array of structs
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

func GetSourceOrganizationVariables(owner string, g *sourceAPIGetter) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/variables", owner)
	zap.S().Debugf("Reading in variables from %v", url)
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
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func (g *APIGetter) CreateRepoVariable(owner string, repo string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/actions/variables", owner, repo)

	resp, err := g.restClient.Request("POST", url, data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	return err
}

func CreateSelectedOrgVariableData(variable data.ImportedVariable) *data.CreateOrgVariable {
	variableArray := make([]int, len(variable.SelectedReposIDs))
	for i := range variableArray {
		variableArray[i], _ = strconv.Atoi(variable.SelectedReposIDs[i])
	}
	s := data.CreateOrgVariable{
		Name:             variable.Name,
		Value:            variable.Value,
		Visibility:       variable.Visibility,
		SelectedReposIDs: variableArray,
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
