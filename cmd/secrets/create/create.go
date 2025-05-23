package createsecrets

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/auth"
	"github.com/katiem0/gh-seva/internal/data"
	"github.com/katiem0/gh-seva/internal/log"
	"github.com/katiem0/gh-seva/internal/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type cmdFlags struct {
	fileName string
	token    string
	hostname string
	debug    bool
}

func NewCmdCreate() *cobra.Command {
	//var repository string
	cmdFlags := cmdFlags{}
	var authToken string

	createCmd := cobra.Command{
		Use:   "create <organization> [flags]",
		Short: "Create Actions, Dependabot, and/or Codespaces secrets from a file.",
		Long:  "Create Actions, Dependabot, and/or Codespaces secrets for an organization and/or repositories from a file.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(createCmd *cobra.Command, args []string) error {
			var err error
			var gqlClient *api.GraphQLClient
			var restClient *api.RESTClient

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

			gqlClient, err = api.NewGraphQLClient(api.ClientOptions{
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

			restClient, err = api.NewRESTClient(api.ClientOptions{
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
	createCmd.PersistentFlags().StringVarP(&cmdFlags.hostname, "hostname", "", "github.com", "GitHub Enterprise Server hostname")
	createCmd.Flags().StringVarP(&cmdFlags.fileName, "from-file", "f", "", "Path and Name of CSV file to create secrets from (required)")
	createCmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "To debug logging")
	if err := createCmd.MarkFlagRequired("from-file"); err != nil {
		zap.S().Errorf("Error marking from-file flag as required: %v", err)
		return nil
	}

	return &createCmd
}

func runCmdCreate(owner string, cmdFlags *cmdFlags, g *utils.APIGetter) error {
	var secretData [][]string
	var importSecretList []data.ImportedSecret
	if len(cmdFlags.fileName) > 0 {
		f, err := os.Open(cmdFlags.fileName)
		zap.S().Debugf("Opening up file %s", cmdFlags.fileName)
		if err != nil {
			zap.S().Errorf("Error arose opening secret csv file")
		}
		defer func() {
			if err := f.Close(); err != nil {
				zap.S().Errorf("Error closing file: %v", err)
			}
		}()
		// read csv values using csv.Reader
		csvReader := csv.NewReader(f)
		secretData, err = csvReader.ReadAll()
		zap.S().Debugf("Reading in all lines from csv file")
		if err != nil {
			zap.S().Errorf("Error arose reading secrets from csv file")
		}
		importSecretList = g.CreateSecretsList(secretData)
	} else {
		zap.S().Errorf("Error arose identifying secrets")
	}
	zap.S().Debugf("Determining secrets to create")
	for _, importSecret := range importSecretList {
		switch importSecret.Level {
		case "Organization":
			zap.S().Debugf("Gathering Organization level secret %s", importSecret.Name)
			switch importSecret.Type {
			case "Actions":
				zap.S().Debugf("Encrypting Organization level Actions secret %s", importSecret.Name)
				publicKey, err := g.GetOrgActionPublicKey(owner)
				if err != nil {
					zap.S().Errorf("Error arose reading Organization Actions secret from csv file")
				}
				var responsePublicKey data.PublicKey
				err = json.Unmarshal(publicKey, &responsePublicKey)
				if err != nil {
					return err
				}
				encryptedSecret, err := g.EncryptSecret(responsePublicKey.Key, importSecret.Value)
				if err != nil {
					return err
				}
				zap.S().Debugf("Creating Organization Actions Secret Data for %s", importSecret.Name)

				var reader io.Reader
				if importSecret.Access == "selected" {
					orgSecretObject := utils.CreateSelectedOrgSecretData(importSecret, responsePublicKey.KeyID, encryptedSecret)
					createSecret, err := json.Marshal(orgSecretObject)
					if err != nil {
						return err
					}
					reader = bytes.NewReader(createSecret)
				} else {
					orgSecretObject := utils.CreateOrgSecretData(importSecret, responsePublicKey.KeyID, encryptedSecret)
					createSecret, err := json.Marshal(orgSecretObject)
					if err != nil {
						return err
					}
					reader = bytes.NewReader(createSecret)
				}

				zap.S().Debugf("Creating Organization Actions Secret %s", importSecret.Name)
				err = g.CreateOrgActionSecret(owner, importSecret.Name, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating Actions secret %s", importSecret.Name)
				}
			case "Codespaces":
				zap.S().Debugf("Encrypting Organization level Codespaces secret %s", importSecret.Name)
				publicKey, err := g.GetOrgCodespacesPublicKey(owner)
				if err != nil {
					zap.S().Errorf("Error arose reading Organization Codespaces secret from csv file")
				}
				var responsePublicKey data.PublicKey
				err = json.Unmarshal(publicKey, &responsePublicKey)
				if err != nil {
					return err
				}
				encryptedSecret, err := g.EncryptSecret(responsePublicKey.Key, importSecret.Value)
				if err != nil {
					return err
				}
				var reader io.Reader
				if importSecret.Access == "selected" {
					orgSecretObject := utils.CreateOrgSecretData(importSecret, responsePublicKey.KeyID, encryptedSecret)
					createSecret, err := json.Marshal(orgSecretObject)
					if err != nil {
						return err
					}
					reader = bytes.NewReader(createSecret)
				} else {
					orgSecretObject := utils.CreateOrgSecretData(importSecret, responsePublicKey.KeyID, encryptedSecret)
					createSecret, err := json.Marshal(orgSecretObject)
					if err != nil {
						return err
					}
					reader = bytes.NewReader(createSecret)
				}
				zap.S().Debugf("Creating Organization Codespaces Secret %s", importSecret.Name)

				err = g.CreateOrgCodespacesSecret(owner, importSecret.Name, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating Organization Codespaces secret %s", importSecret.Name)
				}
			case "Dependabot":
				zap.S().Debugf("Encrypting Organization level Dependabot secret %s", importSecret.Name)
				publicKey, err := g.GetOrgDependabotPublicKey(owner)
				if err != nil {
					zap.S().Errorf("Error arose reading Organization Dependabot secret from csv file")
				}
				var responsePublicKey data.PublicKey
				err = json.Unmarshal(publicKey, &responsePublicKey)
				if err != nil {
					return err
				}
				encryptedSecret, err := g.EncryptSecret(responsePublicKey.Key, importSecret.Value)
				if err != nil {
					return err
				}

				var reader io.Reader
				if importSecret.Access == "selected" {
					orgSecretObject := utils.CreateOrgDependabotSecretData(importSecret, responsePublicKey.KeyID, encryptedSecret)
					createSecret, err := json.Marshal(orgSecretObject)
					if err != nil {
						return err
					}
					reader = bytes.NewReader(createSecret)
				} else {
					orgSecretObject := utils.CreateOrgSecretData(importSecret, responsePublicKey.KeyID, encryptedSecret)
					createSecret, err := json.Marshal(orgSecretObject)
					if err != nil {
						return err
					}
					reader = bytes.NewReader(createSecret)
				}

				zap.S().Debugf("Creating Organization Dependabot Secret %s", importSecret.Name)

				err = g.CreateOrgDependabotSecret(owner, importSecret.Name, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating Organization Dependabot secret %s", importSecret.Name)
				}

			default:
				zap.S().Errorf("Error arose reading secret from csv file")
			}
		case "Repository":
			repoName := importSecret.RepositoryNames[0]
			zap.S().Debugf("Gathering Repository level secret %s", importSecret.Name)
			switch importSecret.Type {
			case "Actions":
				zap.S().Debugf("Encrypting Repository %s level Actions secret %s", repoName, importSecret.Name)
				publicKey, err := g.GetRepoActionPublicKey(owner, repoName)
				if err != nil {
					zap.S().Errorf("Error arose reading Actions secret from csv file")
				}
				var responsePublicKey data.PublicKey
				err = json.Unmarshal(publicKey, &responsePublicKey)
				if err != nil {
					return err
				}
				encryptedSecret, err := g.EncryptSecret(responsePublicKey.Key, importSecret.Value)
				if err != nil {
					return err
				}
				zap.S().Debugf("Creating Repository Actions Secret Data for %s", importSecret.Name)
				repoSecretObject := utils.CreateRepoSecretData(responsePublicKey.KeyID, encryptedSecret)
				createSecret, err := json.Marshal(repoSecretObject)

				if err != nil {
					return err
				}

				reader := bytes.NewReader(createSecret)
				zap.S().Debugf("Creating Actions Secret %s", importSecret.Name)

				err = g.CreateRepoActionSecret(owner, repoName, importSecret.Name, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating Repository Actions secret %s", importSecret.Name)
				}
			case "Codespaces":
				zap.S().Debugf("Encrypting Repository level Codespaces secret %s", importSecret.Name)
				publicKey, err := g.GetRepoCodespacesPublicKey(owner, repoName)
				if err != nil {
					zap.S().Errorf("Error arose reading Repository Codespaces secret from csv file")
				}
				var responsePublicKey data.PublicKey
				err = json.Unmarshal(publicKey, &responsePublicKey)
				if err != nil {
					return err
				}
				encryptedSecret, err := g.EncryptSecret(responsePublicKey.Key, importSecret.Value)
				if err != nil {
					return err
				}
				repoSecretObject := utils.CreateRepoSecretData(responsePublicKey.KeyID, encryptedSecret)
				createSecret, err := json.Marshal(repoSecretObject)

				if err != nil {
					return err
				}

				reader := bytes.NewReader(createSecret)
				zap.S().Debugf("Creating Repository Codespaces Secret %s", importSecret.Name)

				err = g.CreateRepoCodespacesSecret(owner, repoName, importSecret.Name, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating Repository Codespaces secret %s", importSecret.Name)
				}
			case "Dependabot":
				zap.S().Debugf("Encrypting Repository level Dependabot secret %s", importSecret.Name)
				publicKey, err := g.GetRepoDependabotPublicKey(owner, repoName)
				if err != nil {
					zap.S().Errorf("Error arose reading Repository Dependabot secret from csv file")
				}
				var responsePublicKey data.PublicKey
				err = json.Unmarshal(publicKey, &responsePublicKey)
				if err != nil {
					return err
				}
				encryptedSecret, err := g.EncryptSecret(responsePublicKey.Key, importSecret.Value)
				if err != nil {
					return err
				}
				repoSecretObject := utils.CreateRepoSecretData(responsePublicKey.KeyID, encryptedSecret)
				createSecret, err := json.Marshal(repoSecretObject)

				if err != nil {
					return err
				}

				reader := bytes.NewReader(createSecret)
				zap.S().Debugf("Creating Repository Dependabot Secret %s", importSecret.Name)

				err = g.CreateRepoDependabotSecret(owner, repoName, importSecret.Name, reader)
				if err != nil {
					zap.S().Errorf("Error arose creating Repository Dependabot secret %s", importSecret.Name)
				}
			default:
				zap.S().Errorf("Error arose reading secret from csv file")
			}
		default:
			zap.S().Errorf("Error arose reading in where to create secret %s, check csv file.", importSecret.Name)
		}
	}

	fmt.Printf("Successfully created secrets for: %s.", owner)
	return nil
}
