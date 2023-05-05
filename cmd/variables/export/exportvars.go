package exportvars

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

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/auth"
	"github.com/katiem0/gh-seva/internal/data"
	"github.com/katiem0/gh-seva/internal/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
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
		Short: "Generate a report of Actions variables for an organization and/or repositories.",
		Long:  "Generate a report of Actions variables for an organization and/or repositories.",
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

			return runCmdExport(owner, repos, &cmdFlags, data.NewAPIGetter(gqlClient, restClient), reportWriter)
		},
	}

	// Determine default report file based on current timestamp; for more info see https://pkg.go.dev/time#pkg-constants
	reportFileDefault := fmt.Sprintf("report-%s.csv", time.Now().Format("20060102150405"))
	// Configure flags for command
	exportCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub Personal Access Token (default "gh auth token")`)
	exportCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	exportCmd.Flags().StringVarP(&cmdFlags.reportFile, "output-file", "o", reportFileDefault, "Name of file to write CSV report")
	exportCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	//cmd.MarkPersistentFlagRequired("app")

	return &exportCmd
}

func runCmdExport(owner string, repos []string, cmdFlags *cmdFlags, g *data.APIGetter, reportWriter io.Writer) error {
	var reposCursor *string
	var allRepos []data.RepoInfo

	csvWriter := csv.NewWriter(reportWriter)

	err := csvWriter.Write([]string{
		"VariableLevel",
		"VariableName",
		"VariableValue",
		"VariableAccess",
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
	// Writing to CSV Org level Actions Variables
	if len(repos) == 0 {
		zap.S().Debugf("Gathering ORganization level Actions Variables for %s", owner)
		orgVariables, err := g.GetOrgActionVariables(owner)
		if err != nil {
			return err
		}
		var oActionResponseObject data.VariableResponse
		err = json.Unmarshal(orgVariables, &oActionResponseObject)
		if err != nil {
			return err
		}

		for _, orgVariable := range oActionResponseObject.Variables {
			if orgVariable.Visibility == "selected" {
				zap.S().Debugf("Gathering Actions Variables for %s that are scoped to specific repositories", owner)
				scoped_repo, err := g.GetScopedOrgActionVariables(owner, orgVariable.Name)
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
				var responseOObject data.ScopedVariableResponse
				err = json.Unmarshal(scoped_repo, &responseOObject)
				if err != nil {
					return err
				}
				var concatRepos []string
				var concatRepoIds []string
				for _, scopeVariable := range responseOObject.Repositories {
					concatRepos = append(concatRepos, scopeVariable.Name)
					stringRepoId := strconv.Itoa(scopeVariable.ID)
					concatRepoIds = append(concatRepoIds, stringRepoId)
				}
				err = csvWriter.Write([]string{
					"Organization",
					orgVariable.Name,
					orgVariable.Value,
					orgVariable.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else if orgVariable.Visibility == "private" {
				zap.S().Debugf("Gathering Actions Variables %s for %s that is accessible to all internal and private repositories.", orgVariable.Name, owner)
				var concatRepos []string
				var concatRepoIds []string
				for _, repoActPrivateVars := range allRepos {
					if repoActPrivateVars.Visibility != "public" {
						concatRepos = append(concatRepos, repoActPrivateVars.Name)
						stringRepoId := strconv.Itoa(repoActPrivateVars.DatabaseId)
						concatRepoIds = append(concatRepoIds, stringRepoId)
					}
				}
				err = csvWriter.Write([]string{
					"Organization",
					orgVariable.Name,
					orgVariable.Value,
					orgVariable.Visibility,
					strings.Join(concatRepos, ";"),
					strings.Join(concatRepoIds, ";"),
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			} else {
				zap.S().Debugf("Gathering public Actions Secret %s for %s", orgVariable.Name, owner)
				err = csvWriter.Write([]string{
					"Organization",
					orgVariable.Name,
					orgVariable.Value,
					orgVariable.Visibility,
					"",
					"",
				})
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}
			}
		}
	}

	// Writing to CSV repository level Variables
	for _, singleRepo := range allRepos {
		// Writing to CSV repository level Actions Variables
		repoActionVariablesList, err := g.GetRepoActionVariables(owner, singleRepo.Name)
		if err != nil {
			return err
		}
		var repoActionResponseObject data.VariableResponse
		err = json.Unmarshal(repoActionVariablesList, &repoActionResponseObject)
		if err != nil {
			return err
		}
		for _, repoActionsVars := range repoActionResponseObject.Variables {
			err = csvWriter.Write([]string{
				"Repository",
				repoActionsVars.Name,
				repoActionsVars.Value,
				"RepoOnly",
				singleRepo.Name,
				strconv.Itoa(singleRepo.DatabaseId),
			})
			if err != nil {
				zap.S().Error("Error raised in writing output", zap.Error(err))
			}
		}
	}

	csvWriter.Flush()
	fmt.Printf("Successfully exported variables for %s", owner)
	return nil
}
