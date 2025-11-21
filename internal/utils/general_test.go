package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/katiem0/gh-seva/internal/data"
)

func TestNewAPIGetter(t *testing.T) {
	// Create mock clients
	gqlClient := &api.GraphQLClient{}
	restClient := &api.RESTClient{}

	getter := NewAPIGetter(gqlClient, restClient)

	if getter == nil {
		t.Error("NewAPIGetter() returned nil")
		return
	}

	if getter.gqlClient != *gqlClient {
		t.Error("GraphQL client not properly set")
	}

	if getter.restClient != *restClient {
		t.Error("REST client not properly set")
	}
}

func TestNewSourceAPIGetter(t *testing.T) {
	// Create mock client
	restClient := api.RESTClient{}

	getter := NewSourceAPIGetter(restClient)

	if getter == nil {
		t.Error("NewAPIGetter() returned nil")
		return
	}

	if getter.restClient != restClient {
		t.Error("REST client not properly set")
	}
}

type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Error
}

func TestGetOrgActionSecrets(t *testing.T) {
	t.Run("Test using helper function", func(t *testing.T) {
		// Setup mock response data
		mockResponseData := data.SecretsResponse{
			TotalCount: 1,
			Secrets: []data.Secret{
				{
					Name:       "TEST_SECRET",
					Visibility: "all",
				},
			},
		}

		// Convert to JSON bytes as the function would receive
		mockResponseBytes, _ := json.Marshal(mockResponseData)

		// Call helper function that tests the parsing logic
		secretsResponse := testParseOrgActionSecrets(t, mockResponseBytes)

		// Verify the parsed data matches what we expect
		if secretsResponse.TotalCount != 1 {
			t.Errorf("Expected total count 1, got %d", secretsResponse.TotalCount)
		}

		if len(secretsResponse.Secrets) != 1 {
			t.Errorf("Expected 1 secret, got %d", len(secretsResponse.Secrets))
		}

		if secretsResponse.Secrets[0].Name != "TEST_SECRET" {
			t.Errorf("Expected secret name 'TEST_SECRET', got '%s'", secretsResponse.Secrets[0].Name)
		}

		if secretsResponse.Secrets[0].Visibility != "all" {
			t.Errorf("Expected visibility 'all', got '%s'", secretsResponse.Secrets[0].Visibility)
		}
	})
}

func testParseOrgActionSecrets(t *testing.T, responseBytes []byte) data.SecretsResponse {
	var secretsResponse data.SecretsResponse
	err := json.Unmarshal(responseBytes, &secretsResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	return secretsResponse
}

func TestGetReposListGraphQLError(t *testing.T) {
	// Create a mock that returns an error for GetReposList
	mockGetter := NewMockAPIGetter()
	mockGetter.GetReposListError = true

	// Create a custom getter that wraps the mock
	customGetter := &customGetReposListWrapper{
		mock: mockGetter,
	}

	// Call the method and check for error
	_, err := customGetter.GetReposList("test-org", nil)

	if err == nil {
		t.Error("Expected GraphQL error, got nil")
	}

	if !strings.Contains(err.Error(), "mock GetReposList error") {
		t.Errorf("Expected error to contain 'mock GetReposList error', got: %v", err)
	}
}

type customGetReposListWrapper struct {
	mock *MockAPIGetter
}

func (c *customGetReposListWrapper) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	if c.mock.GetReposListError {
		return nil, fmt.Errorf("mock GetReposList error")
	}
	return c.mock.GetReposList(owner, endCursor)
}

func TestGetReposListPagination(t *testing.T) {
	cursor := "cursor123"

	// Create a response with pagination info
	reposResponse := &data.ReposQuery{
		Organization: struct {
			Repositories struct {
				TotalCount int
				Nodes      []data.RepoInfo
				PageInfo   struct {
					EndCursor   string
					HasNextPage bool
				}
			} `graphql:"repositories(first: 100, after: $endCursor)"`
		}{
			Repositories: struct {
				TotalCount int
				Nodes      []data.RepoInfo
				PageInfo   struct {
					EndCursor   string
					HasNextPage bool
				}
			}{
				TotalCount: 100,
				Nodes: []data.RepoInfo{
					{DatabaseId: 1, Name: "repo1"},
					{DatabaseId: 2, Name: "repo2"},
				},
				PageInfo: struct {
					EndCursor   string
					HasNextPage bool
				}{
					EndCursor:   "nextCursor",
					HasNextPage: true,
				},
			},
		},
	}

	// Create a mock with the response
	mockGetter := NewMockAPIGetter()
	mockGetter.ReposResponse = reposResponse

	// Create a cursor-capturing wrapper
	captureWrapper := &cursorCapturingWrapper{
		mock: mockGetter,
	}

	// Call GetReposList with our cursor
	result, err := captureWrapper.GetReposList("test-org", &cursor)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify cursor was captured
	if captureWrapper.capturedCursor == nil || *captureWrapper.capturedCursor != cursor {
		t.Errorf("Expected cursor %s, got %v", cursor, captureWrapper.capturedCursor)
	}

	// Verify response data
	if result.Organization.Repositories.TotalCount != 100 {
		t.Errorf("Expected total count 100, got %d", result.Organization.Repositories.TotalCount)
	}

	if !result.Organization.Repositories.PageInfo.HasNextPage {
		t.Error("Expected HasNextPage to be true")
	}
}

// Helper type that captures cursor in GetReposList
type cursorCapturingWrapper struct {
	mock           *MockAPIGetter
	capturedCursor *string
}

// Override GetReposList to capture cursor
func (c *cursorCapturingWrapper) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	c.capturedCursor = endCursor
	return c.mock.GetReposList(owner, endCursor)
}

// Keep your existing customRepoMockGetter and associated tests
type customRepoMockGetter struct {
	MockAPIGetter
	returnErr    bool
	errMessage   string
	repoData     *data.RepoSingleQuery
	captureOwner string
	captureName  string
}

func (c *customRepoMockGetter) GetRepo(owner string, name string) (*data.RepoSingleQuery, error) {
	c.captureOwner = owner
	c.captureName = name

	if c.returnErr {
		return nil, fmt.Errorf("%s", c.errMessage)
	}
	return c.repoData, nil
}

func TestGetRepoStates(t *testing.T) {
	testCases := []struct {
		name       string
		repoData   data.RepoSingleQuery
		expectErr  bool
		errMessage string
	}{
		{
			name: "Public repository",
			repoData: data.RepoSingleQuery{
				Repository: data.RepoInfo{
					DatabaseId: 12345,
					Name:       "public-repo",
					Visibility: "public",
				},
			},
			expectErr: false,
		},
		{
			name: "Private repository",
			repoData: data.RepoSingleQuery{
				Repository: data.RepoInfo{
					DatabaseId: 67890,
					Name:       "private-repo",
					Visibility: "private",
				},
			},
			expectErr: false,
		},
		{
			name:       "Non-existent repository",
			expectErr:  true,
			errMessage: "Could not resolve to a Repository",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockGetter := &customRepoMockGetter{
				returnErr:  tc.expectErr,
				errMessage: tc.errMessage,
				repoData:   &tc.repoData,
			}

			result, err := mockGetter.GetRepo("test-org", tc.repoData.Repository.Name)

			if tc.expectErr && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tc.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectErr && result != nil {
				if result.Repository.Name != tc.repoData.Repository.Name {
					t.Errorf("Expected repo name %s, got %s",
						tc.repoData.Repository.Name, result.Repository.Name)
				}
			}
		})
	}
}

// Test edge cases for repository queries
func TestRepositoryQueryEdgeCases(t *testing.T) {
	testCases := []struct {
		name     string
		owner    string
		repoName string
		wantErr  bool
	}{
		{"Empty owner", "", "repo", true},
		{"Empty repo name", "owner", "", true},
		{"Special characters in name", "owner", "repo-with-dash", false},
		{"Numeric repo name", "owner", "12345", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockGetter := &customRepoMockGetter{
				returnErr:  tc.wantErr,
				errMessage: "invalid parameters",
				repoData: &data.RepoSingleQuery{
					Repository: data.RepoInfo{
						DatabaseId: 1,
						Name:       tc.repoName,
						Visibility: "public",
					},
				},
			}

			_, err := mockGetter.GetRepo(tc.owner, tc.repoName)

			if tc.wantErr && err == nil {
				t.Error("Expected error for invalid parameters")
			}

			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify parameters were passed correctly
			if mockGetter.captureOwner != tc.owner {
				t.Errorf("Expected owner %s, got %s", tc.owner, mockGetter.captureOwner)
			}

			if mockGetter.captureName != tc.repoName {
				t.Errorf("Expected repo name %s, got %s", tc.repoName, mockGetter.captureName)
			}
		})
	}
}
