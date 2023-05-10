package utils

import (
	"io"

	"github.com/cli/go-gh/pkg/api"
	"github.com/katiem0/gh-seva/internal/data"
	"github.com/shurcooL/graphql"
)

type Getter interface {
	GetReposList(owner string, endCursor *string) ([]data.ReposQuery, error)
	GetRepo(owner string, name string) ([]data.RepoSingleQuery, error)
	GetOrgActionSecrets(owner string) ([]byte, error)
	GetRepoActionSecrets(owner string, repo string) ([]byte, error)
	GetScopedOrgActionSecrets(owner string, secret string) ([]byte, error)
	GetOrgDependabotSecrets(owner string) ([]byte, error)
	GetRepoDependabotSecrets(owner string, repo string) ([]byte, error)
	GetScopedOrgDependabotSecrets(owner string, secret string) ([]byte, error)
	GetOrgCodespacesSecrets(owner string) ([]byte, error)
	GetRepoCodespacesSecrets(owner string, repo string) ([]byte, error)
	GetScopedOrgCodespacesSecrets(owner string, secret string) ([]byte, error)
	CreateSecretsList(data [][]string) []data.ImportedSecret
	GetOrgActionPublicKey(owner string) ([]byte, error)
	GetRepoActionPublicKey(owner string, repo string) ([]byte, error)
	GetOrgCodespacesPublicKey(owner string) ([]byte, error)
	GetRepoCodespacesPublicKey(owner string, repo string) ([]byte, error)
	GetOrgDependabotPublicKey(owner string) ([]byte, error)
	GetRepoDependabotPublicKey(owner string, repo string) ([]byte, error)
	EncryptSecret(publickey string, secret string) (string, error)
	CreateOrgActionSecret(owner string, secret string, data io.Reader) error
	CreateRepoActionSecret(owner string, repo string, secret string, data io.Reader) error
	CreateOrgCodespacesSecret(owner string, secret string, data io.Reader) error
	CreateRepoCodespacesSecret(owner string, repo string, secret string, data io.Reader) error
	CreateOrgDependabotSecret(owner string, secret string, data io.Reader) error
	CreateRepoDependabotSecret(owner string, repo string, secret string, data io.Reader) error
	GetOrgActionVariables(owner string) ([]byte, error)
	GetRepoActionVariables(owner string, repo string) ([]byte, error)
	GetScopedOrgActionVariables(owner string, secret string) ([]byte, error)
}

type APIGetter struct {
	gqlClient  api.GQLClient
	restClient api.RESTClient
}

func NewAPIGetter(gqlClient api.GQLClient, restClient api.RESTClient) *APIGetter {
	return &APIGetter{
		gqlClient:  gqlClient,
		restClient: restClient,
	}
}

type sourceAPIGetter struct {
	restClient api.RESTClient
}

func NewSourceAPIGetter(restClient api.RESTClient) *sourceAPIGetter {
	return &sourceAPIGetter{
		restClient: restClient,
	}
}

func (g *APIGetter) GetReposList(owner string, endCursor *string) (*data.ReposQuery, error) {
	query := new(data.ReposQuery)
	variables := map[string]interface{}{
		"endCursor": (*graphql.String)(endCursor),
		"owner":     graphql.String(owner),
	}

	err := g.gqlClient.Query("getRepos", &query, variables)

	return query, err
}

func (g *APIGetter) GetRepo(owner string, name string) (*data.RepoSingleQuery, error) {
	query := new(data.RepoSingleQuery)
	variables := map[string]interface{}{
		"owner": graphql.String(owner),
		"name":  graphql.String(name),
	}

	err := g.gqlClient.Query("getRepo", &query, variables)
	return query, err
}
