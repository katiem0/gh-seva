package data

import (
	"testing"
	"time"
)

func TestRepoInfo(t *testing.T) {
	// Create a valid RepoInfo struct with all required fields
	repo := RepoInfo{
		DatabaseId: 12345,
		Name:       "test-repo",
		UpdatedAt:  time.Now(),
		Visibility: "private",
	}

	// Verify the values
	if repo.DatabaseId != 12345 {
		t.Errorf("Expected DatabaseId to be 12345, got %d", repo.DatabaseId)
	}
	if repo.Name != "test-repo" {
		t.Errorf("Expected Name to be 'test-repo', got '%s'", repo.Name)
	}
	if repo.Visibility != "private" {
		t.Errorf("Expected Visibility to be 'private', got '%s'", repo.Visibility)
	}
}

func TestReposQuery(t *testing.T) {
	// Create a valid ReposQuery struct
	query := ReposQuery{
		Organization: struct {
			Repositories struct {
				TotalCount int
				Nodes      []RepoInfo
				PageInfo   struct {
					EndCursor   string
					HasNextPage bool
				}
			} `graphql:"repositories(first: 100, after: $endCursor)"`
		}{
			Repositories: struct {
				TotalCount int
				Nodes      []RepoInfo
				PageInfo   struct {
					EndCursor   string
					HasNextPage bool
				}
			}{
				TotalCount: 1,
				Nodes: []RepoInfo{
					{
						DatabaseId: 12345,
						Name:       "test-repo",
						UpdatedAt:  time.Now(),
						Visibility: "private",
					},
				},
				PageInfo: struct {
					EndCursor   string
					HasNextPage bool
				}{
					EndCursor:   "cursor",
					HasNextPage: false,
				},
			},
		},
	}

	// Verify the values
	if query.Organization.Repositories.TotalCount != 1 {
		t.Errorf("Expected TotalCount to be 1, got %d", query.Organization.Repositories.TotalCount)
	}
	if len(query.Organization.Repositories.Nodes) != 1 {
		t.Errorf("Expected 1 repo, got %d", len(query.Organization.Repositories.Nodes))
	}
	if query.Organization.Repositories.Nodes[0].DatabaseId != 12345 {
		t.Errorf("Expected DatabaseId to be 12345, got %d", query.Organization.Repositories.Nodes[0].DatabaseId)
	}
}

func TestRepoSingleQuery(t *testing.T) {
	// Create a valid RepoSingleQuery struct
	query := RepoSingleQuery{
		Repository: RepoInfo{
			DatabaseId: 12345,
			Name:       "test-repo",
			UpdatedAt:  time.Now(),
			Visibility: "private",
		},
	}

	// Verify the values
	if query.Repository.DatabaseId != 12345 {
		t.Errorf("Expected DatabaseId to be 12345, got %d", query.Repository.DatabaseId)
	}
	if query.Repository.Name != "test-repo" {
		t.Errorf("Expected Name to be 'test-repo', got '%s'", query.Repository.Name)
	}
}

func TestScopedRepository(t *testing.T) {
	// Create a valid ScopedRepository struct
	repo := ScopedRepository{
		ID:   12345,
		Name: "test-repo",
	}

	// Verify the values
	if repo.ID != 12345 {
		t.Errorf("Expected ID to be 12345, got %d", repo.ID)
	}
	if repo.Name != "test-repo" {
		t.Errorf("Expected Name to be 'test-repo', got '%s'", repo.Name)
	}
}

func TestScopedResponse(t *testing.T) {
	// Create a valid ScopedResponse struct
	response := ScopedResponse{
		TotalCount: 1,
		Repositories: []ScopedRepository{
			{
				ID:   12345,
				Name: "test-repo",
			},
		},
	}

	// Verify the values
	if response.TotalCount != 1 {
		t.Errorf("Expected TotalCount to be 1, got %d", response.TotalCount)
	}
	if len(response.Repositories) != 1 {
		t.Errorf("Expected 1 repo, got %d", len(response.Repositories))
	}
	if response.Repositories[0].ID != 12345 {
		t.Errorf("Expected ID to be 12345, got %d", response.Repositories[0].ID)
	}
}
