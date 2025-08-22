//go:build !cover
// +build !cover

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// Shared mockRESTClient for all test files
type mockRESTClient struct {
	RequestFunc func(method string, path string, body io.Reader) (*http.Response, error)
}

func (m *mockRESTClient) Request(method string, path string, body io.Reader) (*http.Response, error) {
	if m.RequestFunc != nil {
		return m.RequestFunc(method, path, body)
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
	}, nil
}

func (m *mockRESTClient) RequestWithContext(ctx context.Context, method string, path string, body io.Reader) (*http.Response, error) {
	return m.Request(method, path, body)
}

// Shared testAPIGetterWrapper with all methods
type testAPIGetterWrapper struct {
	mockClient *mockRESTClient
}

func newAPIGetterWithMockREST(client *mockRESTClient) *testAPIGetterWrapper {
	return &testAPIGetterWrapper{
		mockClient: client,
	}
}

// Actions methods
func (t *testAPIGetterWrapper) GetOrgActionSecrets(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/secrets", owner)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err) // Simulate log.Fatal behavior
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetScopedOrgActionSecrets(owner string, secret string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/secrets/%s/repositories", owner, secret)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetOrgActionPublicKey(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/secrets/public-key", owner)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) CreateOrgActionSecret(owner string, secret string, data io.Reader) error {
	url := fmt.Sprintf("orgs/%s/actions/secrets/%s", owner, secret)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

func (t *testAPIGetterWrapper) CreateRepoActionSecret(owner string, repo string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/actions/secrets/%s", owner, repo, secret)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

// Codespaces methods
func (t *testAPIGetterWrapper) GetOrgCodespacesSecrets(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets", owner)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetRepoCodespacesSecrets(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/codespaces/secrets", owner, repo)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetScopedOrgCodespacesSecrets(owner string, secret string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets/%s/repositories", owner, secret)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetOrgCodespacesPublicKey(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets/public-key", owner)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) CreateOrgCodespacesSecret(owner string, secret string, data io.Reader) error {
	url := fmt.Sprintf("orgs/%s/codespaces/secrets/%s", owner, secret)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

func (t *testAPIGetterWrapper) CreateRepoCodespacesSecret(owner string, repo string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/codespaces/secrets/%s", owner, repo, secret)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

// Dependabot methods - Add these missing methods
func (t *testAPIGetterWrapper) GetOrgDependabotSecrets(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets", owner)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetRepoDependabotSecrets(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/dependabot/secrets", owner, repo)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetScopedOrgDependabotSecrets(owner string, secret string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets/%s/repositories", owner, secret)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetOrgDependabotPublicKey(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets/public-key", owner)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) CreateOrgDependabotSecret(owner string, secret string, data io.Reader) error {
	url := fmt.Sprintf("orgs/%s/dependabot/secrets/%s", owner, secret)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

func (t *testAPIGetterWrapper) CreateRepoDependabotSecret(owner string, repo string, secret string, data io.Reader) error {
	url := fmt.Sprintf("repos/%s/%s/dependabot/secrets/%s", owner, repo, secret)
	resp, err := t.mockClient.Request("PUT", url, data)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()
	return nil
}

// Variables methods
func (t *testAPIGetterWrapper) GetOrgActionVariables(owner string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/variables", owner)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetRepoActionVariables(owner string, repo string) ([]byte, error) {
	url := fmt.Sprintf("repos/%s/%s/actions/variables", owner, repo)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

func (t *testAPIGetterWrapper) GetScopedOrgActionVariables(owner string, variable string) ([]byte, error) {
	url := fmt.Sprintf("orgs/%s/actions/variables/%s/repositories", owner, variable)
	resp, err := t.mockClient.Request("GET", url, nil)
	if err != nil {
		panic(err)
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return responseData, err
}

// Add other error mock types here
type errorMockAPIGetter struct {
	MockAPIGetter
	ErrorMessage string
}

func (e *errorMockAPIGetter) GetOrgActionSecrets(owner string) ([]byte, error) {
	return nil, fmt.Errorf(e.ErrorMessage)
}

type secretCreateErrorMockAPIGetter struct {
	MockAPIGetter
}

func (s *secretCreateErrorMockAPIGetter) CreateOrgActionSecret(owner string, secret string, data io.Reader) error {
	return fmt.Errorf("failed to create secret")
}
