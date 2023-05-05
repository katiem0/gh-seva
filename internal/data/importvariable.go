package data

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ImportedVariable struct {
	Level            string
	Name             string `json:"name"`
	Value            string `json:"value"`
	Visibility       string `json:"visibility"`
	SelectedRepos    []string
	SelectedReposIDs []string `json:"selected_repository_ids"`
}

type CreateOrgVariable struct {
	Name             string `json:"name"`
	Value            string `json:"value"`
	Visibility       string `json:"visibility"`
	SelectedReposIDs []int  `json:"selected_repository_ids"`
}

type CreateVariableAll struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Visibility string `json:"visibility"`
}

type CreateRepoVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (g *APIGetter) CreateVariableList(data [][]string) []ImportedVariable {
	// convert csv lines to array of structs
	var variableList []ImportedVariable
	var vars ImportedVariable
	for _, each := range data[1:] {
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

func CreateOrgVariableData(variable ImportedVariable) *CreateOrgVariable {
	variableArray := make([]int, len(variable.SelectedReposIDs))
	for i := range variableArray {
		variableArray[i], _ = strconv.Atoi(variable.SelectedReposIDs[i])
	}
	fmt.Println(variableArray)
	s := CreateOrgVariable{
		Name:             variable.Name,
		Value:            variable.Value,
		Visibility:       variable.Visibility,
		SelectedReposIDs: variableArray,
	}
	return &s
}

func CreateRepoVariableData(variable ImportedVariable) *CreateRepoVariable {
	s := CreateRepoVariable{
		Name:  variable.Name,
		Value: variable.Value,
	}
	return &s
}

func CreateOrgSourceVariableData(variable Variable) *CreateVariableAll {
	s := CreateVariableAll{
		Name:       variable.Name,
		Value:      variable.Value,
		Visibility: variable.Visibility,
	}
	return &s
}
