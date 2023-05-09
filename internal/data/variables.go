package data

import "time"

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

type ImportedVariable struct {
	Level            string
	Name             string `json:"name"`
	Value            string `json:"value"`
	Visibility       string `json:"visibility"`
	SelectedRepos    []string
	SelectedReposIDs []string `json:"selected_repository_ids"`
}

type Variable struct {
	Name          string    `json:"name"`
	Value         string    `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Visibility    string    `json:"visibility"`
	SelectedRepos string    `json:"selected_repositories_url"`
}

type VariableResponse struct {
	TotalCount int        `json:"total_count"`
	Variables  []Variable `json:"variables"`
}
