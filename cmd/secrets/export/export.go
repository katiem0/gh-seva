package exportsecrets

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	gh "github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/katiem0/gh-seva/internal/data"
	"github.com/katiem0/gh-seva/internal/log"
	"github.com/katiem0/gh-seva/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
	app        string
	hostname   string
	token      string
	reportFile string
	debug      bool
}

func NewCmdExport() *cobra.Command {
	//var repository string
	cmdFlags := cmdFlags{}
	var authToken string

	exportCmd := cobra.Command{
		Use:   "export [flags] <organization> [repo ...] ",
		Short: "Generate a report of Actions, Dependabot, and Codespaces secrets for an organization and/or repositories.",
		Long:  "Generate a report of Actions, Dependabot, and Codespaces secrets for an organization and/or repositories.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(exportCmd *cobra.Command, args []string) error {
			var err error
			var gqlClient api.GQLClient
			var restClient api.RESTClient

			// Reinitialize logging if debugging was enabled
			if cmdFlags.debug {
				logger, _ := log.NewLogger(cmdFlags.debug)
				defer logger.Sync() // nolint:errcheck
				zap.ReplaceGlobals(logger)
			}

			if cmdFlags.token != "" {
				authToken = cmdFlags.token
			} else {
				t, _ := auth.TokenForHost(cmdFlags.hostname)
				authToken = t
			}

			gqlClient, err = gh.GQLClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github.hawkgirl-preview+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving graphql client")
				return err
			}

			restClient, err = gh.RESTClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github+json",
				},
				Host:      cmdFlags.hostname,
				AuthToken: authToken,
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving rest client")
				return err
			}

			owner := args[0]
			repos := args[1:]

			if _, err := os.Stat(cmdFlags.reportFile); errors.Is(err, os.ErrExist) {
				return err
			}

			reportWriter, err := os.OpenFile(cmdFlags.reportFile, os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				return err
			}

			return runCmdExport(owner, repos, &cmdFlags, utils.NewAPIGetter(gqlClient, restClient), reportWriter)
		},
	}

	// Determine default report file based on current timestamp; for more info see https://pkg.go.dev/time#pkg-constants
	reportFileDefault := fmt.Sprintf("report-secrets-%s.csv", time.Now().Format("20060102150405"))
	appDefault := "all"
	// Configure flags for command

	exportCmd.PersistentFlags().StringVarP(&cmdFlags.app, "app", "a", appDefault, "List secrets for a specific application or all: {all|actions|codespaces|dependabot}")
	exportCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub Personal Access Token (default "gh auth token")`)
	exportCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	exportCmd.Flags().StringVarP(&cmdFlags.reportFile, "output-file", "o", reportFileDefault, "Name of file to write CSV report")
	exportCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	//cmd.MarkPersistentFlagRequired("app")

	return &exportCmd
}

func runCmdExport(owner string, repos []string, cmdFlags *cmdFlags, g *utils.APIGetter, reportWriter io.Writer) error {
	var reposCursor *string
	var allRepos []data.RepoInfo

	csvWriter := csv.NewWriter(reportWriter)

	err := csvWriter.Write([]string{
		"SecretLevel",
		"SecretType",
		"SecretName",
		"SecretValue",
		"SecretAccess",
		"RepositoryNames",
		"RepositoryIDs",
	})

	if err != nil {
		return err
	}

	if len(repos) > 0 {
		zap.S().Infof("Processing repos: %s", repos)

		for _, repo := range repos {

			zap.S().Debugf("Processing %s/%s", owner, repo)

			repoQuery, err := g.GetRepo(owner, repo)
			if err != nil {
				return err
			}
			allRepos = append(allRepos, repoQuery.Repository)
		}

	} else {
		// Prepare writer for outputting report
		for {
			zap.S().Debugf("Processing list of repositories for %s", owner)
			reposQuery, err := g.GetReposList(owner, reposCursor)

			if err != nil {
				return err
			}

			allRepos = append(allRepos, reposQuery.Organization.Repositories.Nodes...)

			reposCursor = &reposQuery.Organization.Repositories.PageInfo.EndCursor

			if !reposQuery.Organization.Repositories.PageInfo.HasNextPage {
				break
			}
		}
	}

	// Writing to CSV Org level Actions secrets
	if len(repos) == 0 && (cmdFlags.app == "all" || cmdFlags.app == "actions") {
		zap.S().Debugf("Gathering Actions Secrets for %s", owner)
		orgSecrets, err := g.GetOrgActionSecrets(owner)
		if err != nil {
			return err
		}
		var oActionResponseObject data.SecretsResponse
		err = json.Unmarshal(orgSecrets, &oActionResponseObject)
		if err != nil {
			return err
		}

		for _, orgSecret := range oActionResponseObject.Secrets {
			if orgSecret.Visibility == "selected" {
				zap.S().Debugf("Gathering Actions Secrets for %s that are scoped to specific repositories", owner)
				scoped_repo, err := g.GetScopedOrgActionSecrets(owner, orgSecret.Name)
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
				var responseOObject data.ScopedResponse
				err = json.Unmarshal(scoped_repo, &responseOObject)
				if err != nil {
					return err
				}
				var concatRepos []string
				var concatRepoIds []string
				for _, scopeSecret := range responseOObject.Repositories {
					concatRepos = append(concatRepos, scopeSecret.Name)
					stringRepoId := strconv.Itoa(scopeSecret.ID)
					concatRepoIds = append(concatRepoIds, stringRepoId)
				}
				err = csvWriter.Write([]string{
					"Organization",
					"Actions",
					orgSecret.Name,
					"",
					orgSecret.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else if orgSecret.Visibility == "private" {
				zap.S().Debugf("Gathering Actions Secret %s for %s that is accessible to all internal and private repositories.", orgSecret.Name, owner)
				var concatRepos []string
				var concatRepoIds []string
				for _, repoActPrivateSecret := range allRepos {
					if repoActPrivateSecret.Visibility != "public" {
						concatRepos = append(concatRepos, repoActPrivateSecret.Name)
						stringRepoId := strconv.Itoa(repoActPrivateSecret.DatabaseId)
						concatRepoIds = append(concatRepoIds, stringRepoId)
					}
				}
				err = csvWriter.Write([]string{
					"Organization",
					"Actions",
					orgSecret.Name,
					"",
					orgSecret.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else {
				zap.S().Debugf("Gathering public Actions Secret %s for %s", orgSecret.Name, owner)
				err = csvWriter.Write([]string{
					"Organization",
					"Actions",
					orgSecret.Name,
					"",
					orgSecret.Visibility,
					"",
					"",
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			}
		}
	}

	// Writing to CSV Org level Dependabot secrets
	if len(repos) == 0 && (cmdFlags.app == "all" || cmdFlags.app == "dependabot") {
		zap.S().Debugf("Gathering Dependabot Secrets for %s", owner)

		orgDepSecrets, err := g.GetOrgDependabotSecrets(owner)
		if err != nil {
			return err
		}
		var oDepResponseObject data.SecretsResponse
		err = json.Unmarshal(orgDepSecrets, &oDepResponseObject)
		if err != nil {
			return err
		}

		for _, orgDepSecret := range oDepResponseObject.Secrets {
			if orgDepSecret.Visibility == "selected" {
				zap.S().Debugf("Gathering Dependabot Secret %s for %s that is scoped to specific repositories", orgDepSecret.Name, owner)
				scoped_repo, err := g.GetScopedOrgDependabotSecrets(owner, orgDepSecret.Name)
				if err != nil {
					return err
				}
				var rDepResponseObject data.ScopedResponse
				err = json.Unmarshal(scoped_repo, &rDepResponseObject)
				if err != nil {
					return err
				}
				var concatRepos []string
				var concatRepoIds []string
				for _, depScopeSecret := range rDepResponseObject.Repositories {
					concatRepos = append(concatRepos, depScopeSecret.Name)
					stringRepoId := strconv.Itoa(depScopeSecret.ID)
					concatRepoIds = append(concatRepoIds, stringRepoId)
				}
				err = csvWriter.Write([]string{
					"Organization",
					"Dependabot",
					orgDepSecret.Name,
					"",
					orgDepSecret.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else if orgDepSecret.Visibility == "private" {
				zap.S().Debugf("Gathering Dependabot Secret %s for %s that is accessible to all internal and private repositories.", orgDepSecret.Name, owner)
				var concatRepos []string
				var concatRepoIds []string
				for _, repoPrivateSecret := range allRepos {
					if repoPrivateSecret.Visibility != "public" {
						concatRepos = append(concatRepos, repoPrivateSecret.Name)
						stringRepoId := strconv.Itoa(repoPrivateSecret.DatabaseId)
						concatRepoIds = append(concatRepoIds, stringRepoId)
					}
				}
				err = csvWriter.Write([]string{
					"Organization",
					"Dependabot",
					orgDepSecret.Name,
					"",
					orgDepSecret.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else {
				zap.S().Debugf("Gathering public Dependabot Secret %s for %s", orgDepSecret.Name, owner)
				err = csvWriter.Write([]string{
					"Organization",
					"Dependabot",
					orgDepSecret.Name,
					"",
					orgDepSecret.Visibility,
					"",
					"",
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			}
		}
	}

	// Writing to CSV Org level Codespaces secrets
	if len(repos) == 0 && (cmdFlags.app == "all" || cmdFlags.app == "codespaces") {
		zap.S().Debugf("Gathering Codespaces Secrets for %s", owner)

		orgCodeSecrets, err := g.GetOrgCodespacesSecrets(owner)
		if err != nil {
			return err
		}
		var oCodeResponseObject data.SecretsResponse
		err = json.Unmarshal(orgCodeSecrets, &oCodeResponseObject)
		if err != nil {
			return err
		}

		for _, orgCodeSecret := range oCodeResponseObject.Secrets {
			zap.S().Debugf("Gathering Codespaces Secrets for %s that are scoped to specific repositories", owner)
			if orgCodeSecret.Visibility == "selected" {
				scoped_repo, err := g.GetScopedOrgCodespacesSecrets(owner, orgCodeSecret.Name)
				if err != nil {
					return err
				}
				var rCodeResponseObject data.ScopedResponse
				err = json.Unmarshal(scoped_repo, &rCodeResponseObject)
				if err != nil {
					return err
				}
				var concatRepos []string
				var concatRepoIds []string
				for _, codeScopeSecret := range rCodeResponseObject.Repositories {
					concatRepos = append(concatRepos, codeScopeSecret.Name)
					stringRepoId := strconv.Itoa(codeScopeSecret.ID)
					concatRepoIds = append(concatRepoIds, stringRepoId)
				}
				err = csvWriter.Write([]string{
					"Organization",
					"Codespaces",
					orgCodeSecret.Name,
					"",
					orgCodeSecret.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else if orgCodeSecret.Visibility == "private" {
				zap.S().Debugf("Gathering Codespaces Secret %s for %s that is accessible to all internal and private repositories.", orgCodeSecret.Name, owner)
				var concatRepos []string
				var concatRepoIds []string
				for _, repoCodePrivateSecret := range allRepos {
					if repoCodePrivateSecret.Visibility != "public" {
						concatRepos = append(concatRepos, repoCodePrivateSecret.Name)
						stringRepoId := strconv.Itoa(repoCodePrivateSecret.DatabaseId)
						concatRepoIds = append(concatRepoIds, stringRepoId)
					}
				}
				err = csvWriter.Write([]string{
					"Organization",
					"Codespaces",
					orgCodeSecret.Name,
					"",
					orgCodeSecret.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else {
				zap.S().Debugf("Gathering public Codespaces Secret %s for %s", orgCodeSecret.Name, owner)
				err = csvWriter.Write([]string{
					"Organization",
					"Codespaces",
					orgCodeSecret.Name,
					"",
					orgCodeSecret.Visibility,
					"",
					"",
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			}
		}
	}

	// Writing to CSV repository level Secrets
	for _, singleRepo := range allRepos {
		zap.S().Debugf("Gathering Secrets for repo %s", singleRepo.Name)

		// Writing to CSV repository level Actions secrets
		if cmdFlags.app == "all" || cmdFlags.app == "actions" {
			zap.S().Debugf("Gathering Actions Secrets for repo %s", singleRepo.Name)
			repoActionSecretsList, err := g.GetRepoActionSecrets(owner, singleRepo.Name)
			if err != nil {
				return err
			}
			var repoActionResponseObject data.SecretsResponse
			err = json.Unmarshal(repoActionSecretsList, &repoActionResponseObject)
			if err != nil {
				return err
			}
			for _, repoActionsSecret := range repoActionResponseObject.Secrets {
				err = csvWriter.Write([]string{
					"Repository",
					"Actions",
					repoActionsSecret.Name,
					"",
					"RepoOnly",
					singleRepo.Name,
					strconv.Itoa(singleRepo.DatabaseId),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			}
		}
		// Writing to CSV repository level Dependabot secrets
		if cmdFlags.app == "all" || cmdFlags.app == "dependabot" {
			zap.S().Debugf("Gathering Dependabot Secrets for repo %s", singleRepo.Name)
			repoDepSecretsList, err := g.GetRepoDependabotSecrets(owner, singleRepo.Name)
			if err != nil {
				return err
			}
			var repoDepResponseObject data.SecretsResponse
			err = json.Unmarshal(repoDepSecretsList, &repoDepResponseObject)
			if err != nil {
				return err
			}
			for _, repoDepSecret := range repoDepResponseObject.Secrets {
				err = csvWriter.Write([]string{
					"Repository",
					"Dependabot",
					repoDepSecret.Name,
					"",
					"RepoOnly",
					singleRepo.Name,
					strconv.Itoa(singleRepo.DatabaseId),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			}
		}
		// Writing to CSV repository level Codespaces secrets
		if cmdFlags.app == "all" || cmdFlags.app == "codespaces" {
			zap.S().Debugf("Gathering Codespaces Secrets for repo %s", singleRepo.Name)
			repoCodeSecretsList, err := g.GetRepoCodespacesSecrets(owner, singleRepo.Name)
			if err != nil {
				zap.S().Error("Error raised in writing output", zap.Error(err))
			}
			var repoCodeResponseObject data.SecretsResponse
			err = json.Unmarshal(repoCodeSecretsList, &repoCodeResponseObject)
			if err != nil {
				return err
			}
			for _, repoCodeSecret := range repoCodeResponseObject.Secrets {
				err = csvWriter.Write([]string{
					"Repository",
					"Codespaces",
					repoCodeSecret.Name,
					"",
					"RepoOnly",
					singleRepo.Name,
					strconv.Itoa(singleRepo.DatabaseId),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			}
		}
	}

	csvWriter.Flush()
	fmt.Printf("Successfully exported secrets for %s", owner)
	return nil

}
