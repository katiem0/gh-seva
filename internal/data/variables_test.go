package data

import (
	"encoding/json"
	"testing"
	"time"
)

func TestVariable(t *testing.T) {
	// Test Variable struct
	now := time.Now()
	variable := Variable{
		Name:          "TEST_VAR",
		Value:         "test-value",
		CreatedAt:     now,
		UpdatedAt:     now,
		Visibility:    "all",
		SelectedRepos: "https://api.github.com/orgs/testorg/actions/variables/TEST_VAR/repositories",
	}

	if variable.Name != "TEST_VAR" {
		t.Errorf("Expected Name to be 'TEST_VAR', got %s", variable.Name)
	}

	if variable.Value != "test-value" {
		t.Errorf("Expected Value to be 'test-value', got %s", variable.Value)
	}

	if variable.Visibility != "all" {
		t.Errorf("Expected Visibility to be 'all', got %s", variable.Visibility)
	}

	if variable.SelectedRepos != "https://api.github.com/orgs/testorg/actions/variables/TEST_VAR/repositories" {
		t.Errorf("Expected SelectedRepos URL to match, got %s", variable.SelectedRepos)
	}
}

func TestVariableResponse(t *testing.T) {
	// Test VariableResponse struct
	response := VariableResponse{
		TotalCount: 2,
		Variables: []Variable{
			{
				Name:       "VAR1",
				Value:      "value1",
				Visibility: "all",
			},
			{
				Name:       "VAR2",
				Value:      "value2",
				Visibility: "selected",
			},
		},
	}

	if response.TotalCount != 2 {
		t.Errorf("Expected TotalCount to be 2, got %d", response.TotalCount)
	}

	if len(response.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(response.Variables))
	}

	if response.Variables[0].Name != "VAR1" {
		t.Errorf("Expected first variable name to be 'VAR1', got %s", response.Variables[0].Name)
	}

	if response.Variables[1].Value != "value2" {
		t.Errorf("Expected second variable value to be 'value2', got %s", response.Variables[1].Value)
	}
}

func TestImportedVariable(t *testing.T) {
	// Test ImportedVariable struct
	variable := ImportedVariable{
		Level:            "Organization",
		Name:             "TEST_VAR",
		Value:            "var-value",
		Visibility:       "selected",
		SelectedRepos:    []string{"repo1", "repo2"},
		SelectedReposIDs: []string{"1", "2"},
	}

	if variable.Level != "Organization" {
		t.Errorf("Expected Level to be 'Organization', got %s", variable.Level)
	}

	if variable.Name != "TEST_VAR" {
		t.Errorf("Expected Name to be 'TEST_VAR', got %s", variable.Name)
	}

	if variable.Value != "var-value" {
		t.Errorf("Expected Value to be 'var-value', got %s", variable.Value)
	}

	if variable.Visibility != "selected" {
		t.Errorf("Expected Visibility to be 'selected', got %s", variable.Visibility)
	}

	if len(variable.SelectedRepos) != 2 {
		t.Errorf("Expected 2 selected repos, got %d", len(variable.SelectedRepos))
	}

	if len(variable.SelectedReposIDs) != 2 {
		t.Errorf("Expected 2 selected repo IDs, got %d", len(variable.SelectedReposIDs))
	}
}

func TestVariableStructsJSON(t *testing.T) {
	// Test JSON marshaling/unmarshaling

	// Test CreateOrgVariable
	orgVar := CreateOrgVariable{
		Name:             "TEST_VAR",
		Value:            "test-value",
		Visibility:       "selected",
		SelectedReposIDs: []int{1, 2, 3},
	}

	jsonBytes, err := json.Marshal(orgVar)
	if err != nil {
		t.Fatalf("Failed to marshal CreateOrgVariable: %v", err)
	}

	var unmarshaledOrgVar CreateOrgVariable
	err = json.Unmarshal(jsonBytes, &unmarshaledOrgVar)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateOrgVariable: %v", err)
	}

	if unmarshaledOrgVar.Name != "TEST_VAR" {
		t.Errorf("Expected Name to be 'TEST_VAR', got %s", unmarshaledOrgVar.Name)
	}

	if len(unmarshaledOrgVar.SelectedReposIDs) != 3 {
		t.Errorf("Expected 3 SelectedReposIDs, got %d", len(unmarshaledOrgVar.SelectedReposIDs))
	}

	// Test CreateVariableAll
	varAll := CreateVariableAll{
		Name:       "TEST_VAR",
		Value:      "test-value",
		Visibility: "all",
	}

	jsonBytes, err = json.Marshal(varAll)
	if err != nil {
		t.Fatalf("Failed to marshal CreateVariableAll: %v", err)
	}

	var unmarshaledVarAll CreateVariableAll
	err = json.Unmarshal(jsonBytes, &unmarshaledVarAll)
	if err != nil {
		t.Fatalf("Failed to unmarshal CreateVariableAll: %v", err)
	}

	if unmarshaledVarAll.Visibility != "all" {
		t.Errorf("Expected Visibility to be 'all', got %s", unmarshaledVarAll.Visibility)
	}
}
