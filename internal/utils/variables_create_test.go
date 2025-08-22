package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/katiem0/gh-seva/internal/data"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestCreateVariableList(t *testing.T) {
	// Setup
	g := &APIGetter{}
	filedata := [][]string{
		{"VariableLevel", "VariableName", "VariableValue", "VariableAccess", "RepositoryNames", "RepositoryIDs"},
		{"Organization", "TEST_VAR1", "value1", "all", "", ""},
		{"Organization", "TEST_VAR2", "value2", "selected", "repo1;repo2", "1234;5678"},
		{"Repository", "TEST_VAR3", "value3", "RepoOnly", "repo3", "9101"},
	}

	// Execute
	result := g.CreateVariableList(filedata)

	// Verify
	expected := []data.ImportedVariable{
		{
			Level:            "Organization",
			Name:             "TEST_VAR1",
			Value:            "value1",
			Visibility:       "all",
			SelectedRepos:    []string{""},
			SelectedReposIDs: []string{""},
		},
		{
			Level:            "Organization",
			Name:             "TEST_VAR2",
			Value:            "value2",
			Visibility:       "selected",
			SelectedRepos:    []string{"repo1", "repo2"},
			SelectedReposIDs: []string{"1234", "5678"},
		},
		{
			Level:            "Repository",
			Name:             "TEST_VAR3",
			Value:            "value3",
			Visibility:       "RepoOnly",
			SelectedRepos:    []string{"repo3"},
			SelectedReposIDs: []string{"9101"},
		},
	}

	if len(result) != len(expected) {
		t.Errorf("Expected %d variables, got %d", len(expected), len(result))
		return
	}

	for i, e := range expected {
		if result[i].Level != e.Level {
			t.Errorf("Variable %d: Expected Level %s, got %s", i, e.Level, result[i].Level)
		}
		if result[i].Name != e.Name {
			t.Errorf("Variable %d: Expected Name %s, got %s", i, e.Name, result[i].Name)
		}
		if result[i].Value != e.Value {
			t.Errorf("Variable %d: Expected Value %s, got %s", i, e.Value, result[i].Value)
		}
		if result[i].Visibility != e.Visibility {
			t.Errorf("Variable %d: Expected Visibility %s, got %s", i, e.Visibility, result[i].Visibility)
		}

		// Check arrays
		if len(result[i].SelectedRepos) != len(e.SelectedRepos) {
			t.Errorf("Variable %d: Expected %d SelectedRepos, got %d", i, len(e.SelectedRepos), len(result[i].SelectedRepos))
		} else {
			for j, repo := range e.SelectedRepos {
				if result[i].SelectedRepos[j] != repo {
					t.Errorf("Variable %d, SelectedRepo %d: Expected %s, got %s", i, j, repo, result[i].SelectedRepos[j])
				}
			}
		}

		if len(result[i].SelectedReposIDs) != len(e.SelectedReposIDs) {
			t.Errorf("Variable %d: Expected %d SelectedReposIDs, got %d", i, len(e.SelectedReposIDs), len(result[i].SelectedReposIDs))
		} else {
			for j, id := range e.SelectedReposIDs {
				if result[i].SelectedReposIDs[j] != id {
					t.Errorf("Variable %d, SelectedRepoID %d: Expected %s, got %s", i, j, id, result[i].SelectedReposIDs[j])
				}
			}
		}
	}
}

func TestCreateVariableListLogging(t *testing.T) {
	// Setup mock logger
	oldLogger := zap.L()
	defer zap.ReplaceGlobals(oldLogger)

	var logs bytes.Buffer
	testLogger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(&logs),
			zapcore.DebugLevel,
		)
	})))
	zap.ReplaceGlobals(testLogger)

	g := &APIGetter{}
	filedata := [][]string{
		{"VariableLevel", "VariableName", "VariableValue", "VariableAccess", "RepositoryNames", "RepositoryIDs"},
		{"Organization", "TEST_VAR1", "value1", "all", "", ""},
	}

	// Execute
	result := g.CreateVariableList(filedata)

	// Verify
	if len(result) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(result))
	}

	// Verify logging occurred
	logOutput := logs.String()
	if !strings.Contains(logOutput, "TEST_VAR1") {
		t.Errorf("Expected log output to contain variable name TEST_VAR1")
	}
}

func TestCreateSelectedOrgVariableData(t *testing.T) {
	// Setup
	variable := data.ImportedVariable{
		Name:             "TEST_VAR",
		Value:            "test-value",
		Visibility:       "selected",
		SelectedReposIDs: []string{"1234", "5678"},
	}

	// Execute
	result := CreateSelectedOrgVariableData(variable)

	// Verify
	expected := &data.CreateOrgVariable{
		Name:             "TEST_VAR",
		Value:            "test-value",
		Visibility:       "selected",
		SelectedReposIDs: []int{1234, 5678},
	}

	if result.Name != expected.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, result.Name)
	}

	if result.Value != expected.Value {
		t.Errorf("Expected Value %s, got %s", expected.Value, result.Value)
	}

	if result.Visibility != expected.Visibility {
		t.Errorf("Expected Visibility %s, got %s", expected.Visibility, result.Visibility)
	}

	if len(result.SelectedReposIDs) != len(expected.SelectedReposIDs) {
		t.Errorf("Expected %d SelectedReposIDs, got %d", len(expected.SelectedReposIDs), len(result.SelectedReposIDs))
		return
	}

	for i, id := range expected.SelectedReposIDs {
		if result.SelectedReposIDs[i] != id {
			t.Errorf("SelectedReposIDs[%d]: Expected %d, got %d", i, id, result.SelectedReposIDs[i])
		}
	}
}

func TestCreateSelectedOrgVariableDataWithInvalidIDs(t *testing.T) {
	// Setup
	variable := data.ImportedVariable{
		Name:             "TEST_VAR",
		Value:            "test-value",
		Visibility:       "selected",
		SelectedReposIDs: []string{"1234", "invalid", "5678"},
	}

	// Execute
	result := CreateSelectedOrgVariableData(variable)

	// Verify - only valid IDs should be included
	if len(result.SelectedReposIDs) != 2 {
		t.Errorf("Expected 2 valid IDs, got %d", len(result.SelectedReposIDs))
	}

	// Check specific IDs
	validIDs := map[int]bool{1234: true, 5678: true}
	for _, id := range result.SelectedReposIDs {
		if !validIDs[id] {
			t.Errorf("Unexpected ID in result: %d", id)
		}
	}
}

func TestCreateOrgVariableData(t *testing.T) {
	// Setup
	variable := data.ImportedVariable{
		Name:       "TEST_VAR",
		Value:      "test-value",
		Visibility: "all",
	}

	// Execute
	result := CreateOrgVariableData(variable)

	// Verify
	expected := &data.CreateVariableAll{
		Name:       "TEST_VAR",
		Value:      "test-value",
		Visibility: "all",
	}

	if result.Name != expected.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, result.Name)
	}

	if result.Value != expected.Value {
		t.Errorf("Expected Value %s, got %s", expected.Value, result.Value)
	}

	if result.Visibility != expected.Visibility {
		t.Errorf("Expected Visibility %s, got %s", expected.Visibility, result.Visibility)
	}
}

func TestCreateRepoVariableData(t *testing.T) {
	// Setup
	variable := data.ImportedVariable{
		Name:  "TEST_VAR",
		Value: "test-value",
	}

	// Execute
	result := CreateRepoVariableData(variable)

	// Verify
	expected := &data.CreateRepoVariable{
		Name:  "TEST_VAR",
		Value: "test-value",
	}

	if result.Name != expected.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, result.Name)
	}

	if result.Value != expected.Value {
		t.Errorf("Expected Value %s, got %s", expected.Value, result.Value)
	}
}

func TestCreateOrgSourceVariableData(t *testing.T) {
	// Setup
	variable := data.Variable{
		Name:       "TEST_VAR",
		Value:      "test-value",
		Visibility: "all",
	}

	// Execute
	result := CreateOrgSourceVariableData(variable)

	// Verify
	expected := &data.CreateVariableAll{
		Name:       "TEST_VAR",
		Value:      "test-value",
		Visibility: "all",
	}

	if result.Name != expected.Name {
		t.Errorf("Expected Name %s, got %s", expected.Name, result.Name)
	}

	if result.Value != expected.Value {
		t.Errorf("Expected Value %s, got %s", expected.Value, result.Value)
	}

	if result.Visibility != expected.Visibility {
		t.Errorf("Expected Visibility %s, got %s", expected.Visibility, result.Visibility)
	}
}

type mockSourceAPIGetter struct {
	VariablesData      []byte
	ScopedVariableData []byte
	ShouldReturnError  bool
}

func GetSourceOrganizationVariablesTest(owner string, g interface{}) ([]byte, error) {
	// Cast the interface to our mock type
	mockGetter, ok := g.(*mockSourceAPIGetter)
	if !ok {
		return nil, fmt.Errorf("invalid getter type")
	}

	// Return the mock data
	return mockGetter.VariablesData, nil
}

func GetScopedSourceOrgActionVariablesTest(owner string, secret string, g interface{}) ([]byte, error) {
	// Cast the interface to our mock type
	mockGetter, ok := g.(*mockSourceAPIGetter)
	if !ok {
		return nil, fmt.Errorf("invalid getter type")
	}

	// Return the mock data
	return mockGetter.ScopedVariableData, nil
}

func (m *mockSourceAPIGetter) Request(method string, url string, body io.Reader) (*mockResponse, error) {
	if m.ShouldReturnError {
		return nil, fmt.Errorf("mock error for source API request")
	}

	if strings.Contains(url, "repositories") {
		return &mockResponse{Body: io.NopCloser(bytes.NewReader(m.ScopedVariableData))}, nil
	}
	return &mockResponse{Body: io.NopCloser(bytes.NewReader(m.VariablesData))}, nil
}

type MockRESTClient struct {
	StatusCode int
	Response   []byte
}

func (m MockRESTClient) Request(method string, url string, body io.Reader) (*http.Response, error) {
	// Create a mock response with the configured status code
	resp := &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(bytes.NewReader(m.Response)),
	}

	return resp, nil
}

func (m MockRESTClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(bytes.NewReader(m.Response)),
	}

	return resp, nil
}

func (m MockRESTClient) BuildRequestURL(path string) string {
	return path
}

func (m MockRESTClient) GraphQL(query string, variables map[string]interface{}, result interface{}) error {
	return nil
}

type mockResponse struct {
	Body io.ReadCloser
}

func TestGetSourceOrganizationVariables(t *testing.T) {
	// Setup
	expectedResponse := []byte(`{"total_count":1,"variables":[{"name":"TEST_VAR","value":"test-value","visibility":"all"}]}`)
	mockGetter := &mockSourceAPIGetter{
		VariablesData: expectedResponse,
	}

	// Execute
	response, err := GetSourceOrganizationVariablesTest("test-org", mockGetter)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}

	// Test parsing of the response
	var variablesResponse data.VariableResponse
	err = json.Unmarshal(response, &variablesResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if variablesResponse.TotalCount != 1 {
		t.Errorf("Expected total count 1, got %d", variablesResponse.TotalCount)
	}

	if len(variablesResponse.Variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(variablesResponse.Variables))
	}

	if variablesResponse.Variables[0].Name != "TEST_VAR" {
		t.Errorf("Expected variable name TEST_VAR, got %s", variablesResponse.Variables[0].Name)
	}

	if variablesResponse.Variables[0].Value != "test-value" {
		t.Errorf("Expected variable value test-value, got %s", variablesResponse.Variables[0].Value)
	}

	if variablesResponse.Variables[0].Visibility != "all" {
		t.Errorf("Expected variable visibility all, got %s", variablesResponse.Variables[0].Visibility)
	}
}

func TestGetScopedSourceOrgActionVariables(t *testing.T) {
	// Setup
	expectedResponse := []byte(`{"total_count":1,"repositories":[{"id":12345,"name":"test-repo"}]}`)
	mockGetter := &mockSourceAPIGetter{
		ScopedVariableData: expectedResponse,
	}

	// Execute
	response, err := GetScopedSourceOrgActionVariablesTest("test-org", "TEST_VAR", mockGetter)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if string(response) != string(expectedResponse) {
		t.Errorf("Expected response %s, got %s", string(expectedResponse), string(response))
	}

	// Test parsing of the response
	var scopedResponse data.ScopedResponse
	err = json.Unmarshal(response, &scopedResponse)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if scopedResponse.TotalCount != 1 {
		t.Errorf("Expected total count 1, got %d", scopedResponse.TotalCount)
	}

	if len(scopedResponse.Repositories) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(scopedResponse.Repositories))
	}

	if scopedResponse.Repositories[0].ID != 12345 {
		t.Errorf("Expected repository ID 12345, got %d", scopedResponse.Repositories[0].ID)
	}

	if scopedResponse.Repositories[0].Name != "test-repo" {
		t.Errorf("Expected repository name test-repo, got %s", scopedResponse.Repositories[0].Name)
	}
}

// Test the API methods that need a REST client
func TestCreateOrganizationVariable(t *testing.T) {
	// Setup - using mock
	mockGetter := NewMockAPIGetter()

	// Prepare test data
	variableData := data.CreateVariableAll{
		Name:       "TEST_VAR",
		Value:      "test-value",
		Visibility: "all",
	}
	jsonData, _ := json.Marshal(variableData)
	reader := bytes.NewReader(jsonData)

	// Execute
	err := mockGetter.CreateOrganizationVariable("test-org", reader)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCreateRepoVariable(t *testing.T) {
	// Setup - using mock
	mockGetter := NewMockAPIGetter()

	// Prepare test data
	variableData := data.CreateRepoVariable{
		Name:  "TEST_VAR",
		Value: "test-value",
	}
	jsonData, _ := json.Marshal(variableData)
	reader := bytes.NewReader(jsonData)

	// Execute
	err := mockGetter.CreateRepoVariable("test-org", "test-repo", reader)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Helper function tests
func TestStrToIntVariable(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{"Valid integer", "123", 123, false},
		{"Zero", "0", 0, false},
		{"Negative integer", "-123", -123, false},
		{"Empty string", "", 0, true},
		{"Invalid input", "abc", 0, true},
		{"Mixed input", "123abc", 0, true},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := strconv.Atoi(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %s but got none", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %s: %v", tc.input, err)
				}

				if result != tc.expected {
					t.Errorf("Expected %d, got %d", tc.expected, result)
				}
			}
		})
	}
}

func TestCreateOrganizationVariableError(t *testing.T) {
	// Setup - mock that returns an error
	mockGetter := &MockAPIGetter{
		ShouldReturnError: true,
	}

	// Prepare test data
	variableData := data.CreateVariableAll{
		Name:       "TEST_VAR",
		Value:      "test-value",
		Visibility: "all",
	}
	jsonData, _ := json.Marshal(variableData)
	reader := bytes.NewReader(jsonData)

	// Execute
	err := mockGetter.CreateOrganizationVariable("test-org", reader)

	// Verify
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreateRepoVariableError(t *testing.T) {
	// Setup - mock that returns an error
	mockGetter := &MockAPIGetter{
		ShouldReturnError: true,
	}

	// Prepare test data
	variableData := data.CreateRepoVariable{
		Name:  "TEST_VAR",
		Value: "test-value",
	}
	jsonData, _ := json.Marshal(variableData)
	reader := bytes.NewReader(jsonData)

	// Execute
	err := mockGetter.CreateRepoVariable("test-org", "test-repo", reader)

	// Verify
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreateVariableListEmptyInput(t *testing.T) {
	// Setup
	g := &APIGetter{}

	// Test with empty input
	emptyData := [][]string{
		{"VariableLevel", "VariableName", "VariableValue", "VariableAccess", "RepositoryNames", "RepositoryIDs"},
	}

	// Execute
	result := g.CreateVariableList(emptyData)

	// Verify
	if len(result) != 0 {
		t.Errorf("Expected empty result for empty input, got %d items", len(result))
	}
}

func TestCreateVariableListInvalidData(t *testing.T) {
	// Setup
	g := &APIGetter{}

	// Test with invalid data (missing required fields)
	invalidData := [][]string{
		{"VariableLevel", "VariableName", "VariableValue", "VariableAccess", "RepositoryNames", "RepositoryIDs"},
		{"Organization", "", "value1", "all", "", ""}, // Missing name
	}

	// Execute
	result := g.CreateVariableList(invalidData)

	// Verify
	if len(result) != 0 {
		t.Errorf("Expected empty result for invalid input, got %d items", len(result))
	}
}
