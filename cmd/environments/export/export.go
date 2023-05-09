package exportenvs

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
		Short: "Generate a report of environments and metadata.",
		Long:  "Generate a report of environments and metadata for a single repository or all repositories in an organization.",
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
	reportFileDefault := fmt.Sprintf("report-%s.csv", time.Now().Format("20060102150405"))

	// Configure flags for command

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
		"RepositoryName",
		"RepositoryID",
		"EnvironmentName",
		"AdminBypass",
		"WaitTimer",
		"Reviewers",
		"ProtectedBranches",
		"CustomBranchPolicies",
		"SecretsTotalCount",
		"SecretsList",
		"VariablesTotalCount",
		"VariablesList",
	})

	if err != nil {
		zap.S().Error("Error raised in writing output", zap.Error(err))
	}

	if len(repos) > 0 {
		zap.S().Infof("Processing repos: %s", repos)

		for _, repo := range repos {

			zap.S().Debugf("Processing %s/%s", owner, repo)

			repoQuery, err := g.GetRepo(owner, repo)
			if err != nil {
				zap.S().Error("Error raised in gathering repo", zap.Error(err))
			}
			allRepos = append(allRepos, repoQuery.Repository)
		}

	} else {
		// Prepare writer for outputting report
		for {
			zap.S().Debugf("Processing list of repositories for %s", owner)
			reposQuery, err := g.GetReposList(owner, reposCursor)

			if err != nil {
				zap.S().Error("Error raised in gathering repos", zap.Error(err))
			}

			allRepos = append(allRepos, reposQuery.Organization.Repositories.Nodes...)

			reposCursor = &reposQuery.Organization.Repositories.PageInfo.EndCursor

			if !reposQuery.Organization.Repositories.PageInfo.HasNextPage {
				break
			}
		}
	}
	// Gathering Envs for each repository listed

	zap.S().Debug("Gathering all repository environments")
	for _, singleRepo := range allRepos {
		zap.S().Debugf("Gathering Environments for repo %s", singleRepo.Name)
		repoEnvs, err := g.GetRepoEnvironments(owner, singleRepo.Name)
		if err != nil {
			zap.S().Error("Error raised in writing output", zap.Error(err))
		}
		var responseEnvs data.EnvResponse
		err = json.Unmarshal(repoEnvs, &responseEnvs)
		if err != nil {
			return err
		}

		zap.S().Debugf("Writing data for %d environment(s) to output for repository %s", responseEnvs.TotalCount, singleRepo.Name)
		for _, env := range responseEnvs.Environments {
			var waitTimer int
			var Reviewers []string
			for _, rules := range env.ProtectionRules {
				if rules.Type == "wait_timer" {
					waitTimer = rules.WaitTimer

				} else if rules.Type == "required_reviewers" {
					for _, reviewer := range rules.Reviewers {
						var reviewList []string
						reviewList = append(reviewList, reviewer.Type)
						reviewList = append(reviewList, reviewer.Reviewer.Login)
						reviewList = append(reviewList, strconv.Itoa(reviewer.Reviewer.ID))
						reviewLists := strings.Join(reviewList, ";")
						Reviewers = append(Reviewers, reviewLists)

					}
					fmt.Println(strings.Join(Reviewers, "|"))
				} else if rules.Type == "branch_policy" {

				}
			}

			if err != nil {
				zap.S().Error("Error raised in writing output", zap.Error(err))
			}
			err = csvWriter.Write([]string{
				singleRepo.Name,
				strconv.Itoa(singleRepo.DatabaseId),
				env.Name,
				strconv.FormatBool(env.AdminByPass),
				strconv.Itoa(waitTimer),
				fmt.Sprintf(strings.Join(Reviewers, "|")),
			})

			if err != nil {
				zap.S().Error("Error raised in writing output", zap.Error(err))
			}
		}
	}
	csvWriter.Flush()
	fmt.Printf("Successfully exported environment data to csv %s", cmdFlags.reportFile)

	return nil
}
