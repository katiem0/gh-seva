package createvars

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"

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
	sourceToken    string
	sourceOrg      string
	sourceHostname string
	fileName       string
	token          string
	hostname       string
	debug          bool
}

func NewCmdCreate() *cobra.Command {
	//var repository string
	cmdFlags := cmdFlags{}
	var authToken string

	createCmd := cobra.Command{
		Use:   "create <organization> [flags]",
		Short: "Create Organization Actions variables.",
		Long:  "Create Organization Actions variables for a specified organization or organization and repositories level variables from a file.",
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(createCmd *cobra.Command, args []string) error {
			if len(cmdFlags.fileName) == 0 && len(cmdFlags.sourceOrg) == 0 {
				return errors.New("A file or source organization must be specified where variables will be created from.")
			} else if len(cmdFlags.sourceOrg) > 0 && len(cmdFlags.sourceToken) == 0 {
				return errors.New("A Personal Access Token must be specified to access variables from the Source Organization.")
			} else if len(cmdFlags.fileName) > 0 && len(cmdFlags.sourceOrg) > 0 {
				return errors.New("Specify only one of `--source-organization` or `from-file`.")
			}
			return nil
		},
		RunE: func(createCmd *cobra.Command, args []string) error {
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

			return runCmdCreate(owner, &cmdFlags, utils.NewAPIGetter(gqlClient, restClient))
		},
	}

	// Configure flags for command
	createCmd.PersistentFlags().StringVarP(&cmdFlags.token, "token", "t", "", `GitHub personal access token for organization to write to (default "gh auth token")`)
	createCmd.PersistentFlags().StringVarP(&cmdFlags.sourceToken, "source-token", "s", "", `GitHub personal access token for Source Organization (Required for --source-organization)`)
	createCmd.PersistentFlags().StringVarP(&cmdFlags.sourceOrg, "source-organization", "o", "", `Name of the Source Organization to copy variables from (Requires --source-token)`)
	createCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	createCmd.PersistentFlags().StringVarP(&cmdFlags.sourceHostname, "source-hostname", "", "github.com", "GitHub Enterprise Server hostname where variables are copied from")
	createCmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to create variables from")
	createCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")

	return &createCmd
}

func runCmdCreate(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var variableData [][]string
	var variablesList []data.ImportedVariable

	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening variables csv file")
		}
		// remember to close the file at the end of the program
		defer f.Close()

		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		variableData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading variables from csv file")
		}
		variablesList = g.CreateVariableList(variableData)
		zap.S().Debugf("Identifying Variable list to create under %s", owner)
		zap.S().Debugf("Determining variables to create")
		for _, variable := range variablesList {

			if variable.Level == "Organization" {
				zap.S().Debugf("Gathering Organization level variable %s", variable.Name)
				importOrgVar := utils.CreateOrgVariableData(variable)
				createVariable, err := json.Marshal(importOrgVar)

				if err != nil {
					return err
				}

				reader := bytes.NewReader(createVariable)
				zap.S().Debugf("Creating Variables under %s", owner)
				err = g.CreateOrganizationVariable(owner, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating variable with %s", variable.Name)
				}
			} else if variable.Level == "Repository" {
				repoName := variable.SelectedRepos[0]
				zap.S().Debugf("Gathering Repository level variable %s", variable.Name)
				importRepoVar := utils.CreateRepoVariableData(variable)
				createVariable, err := json.Marshal(importRepoVar)

				if err != nil {
					return err
				}

				reader := bytes.NewReader(createVariable)
				zap.S().Debugf("Creating Variables under %s", repoName)
				err = g.CreateRepoVariable(owner, repoName, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating variable with %s", variable.Name)
				}
			}
		}
	} else if len(cmdFlags.sourceOrg) > 0 {
		zap.S().Debugf("Reading in variables from %s", cmdFlags.sourceOrg)
		var authToken string
		var restSourceClient api.RESTClient

		if cmdFlags.sourceToken != "" {
			authToken = cmdFlags.sourceToken
		} else {
			t, _ := auth.TokenForHost(cmdFlags.sourceHostname)
			authToken = t
		}

		restSourceClient, err := gh.RESTClient(&api.ClientOptions{
			Headers: map[string]string{
				"Accept": "application/vnd.github+json",
			},
			Host:      cmdFlags.sourceHostname,
			AuthToken: authToken,
		})
		if err != nil {
			zap.S().Errorf("Error arose retrieving source rest client")
			return err
		}

		zap.S().Debugf("Gathering variables %s", cmdFlags.sourceOrg)

		variableResponse, err := utils.GetSourceOrganizationVariables(cmdFlags.sourceOrg, utils.NewSourceAPIGetter(restSourceClient))
		if err != nil {
			return err
		}
		var response data.VariableResponse
		err = json.Unmarshal(variableResponse, &response)
		if err != nil {
			return err
		}
		for _, variable := range response.Variables {
			if variable.Visibility == "selected" {
				zap.S().Debugf("Creating Scoped Variables under %s", owner)
				var orgVariable data.CreateOrgVariable
				scoped_repo, err := utils.GetScopedSourceOrgActionVariables(cmdFlags.sourceOrg, variable.Name, utils.NewSourceAPIGetter(restSourceClient))
				if err != nil {
					zap.S().Error("Error raised in writing output", zap.Error(err))
				}

				var responseOObject data.ScopedResponse
				err = json.Unmarshal(scoped_repo, &responseOObject)
				if err != nil {
					return err
				}
				var concatRepoIds []int
				for _, scopedVar := range responseOObject.Repositories {
					concatRepoIds = append(concatRepoIds, scopedVar.ID)
				}
				orgVariable.SelectedReposIDs = concatRepoIds
				orgVariable.Name = variable.Name
				orgVariable.Value = variable.Value
				orgVariable.Visibility = variable.Visibility

				createOrgVariable, err := json.Marshal(orgVariable)

				if err != nil {
					return err
				}
				reader := bytes.NewReader(createOrgVariable)
				zap.S().Debugf("Creating Variables under %s", owner)
				err = g.CreateOrganizationVariable(owner, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating variable with %s", variable.Name)
				}
			} else {
				orgVariable := utils.CreateOrgSourceVariableData(variable)
				createOrgVariable, err := json.Marshal(orgVariable)

				if err != nil {
					return err
				}
				reader := bytes.NewReader(createOrgVariable)
				zap.S().Debugf("Creating Variable %s under %s", variable.Name, owner)
				err = g.CreateOrganizationVariable(owner, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating variable with %s", variable.Name)
				}
			}
		}
	} else {
		zap.S().Errorf("Error arose identifying variables")
	}

	fmt.Printf("Successfully created variables for: %s.", owner)
	return nil
}
