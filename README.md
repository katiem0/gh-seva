# gh-seva

A GitHub `gh` [CLI](https://cli.github.com/) extension to list and create Secrets and Variables defined at an Organization level and/or Repository level.

## Installation

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation) instructions.

2. Install the extension:
   ```sh
   gh extension install katiem0/gh-seva
   ```

For more information: [`gh extension install`](https://cli.github.com/manual/gh_extension_install).

## Usage

This extension supports listing and creating secrets and variables between `GitHub.com` and GitHub Enterprise Server, through the use of `--hostname` and `--source-hostname`.

```sh
$ gh seva -h
Export and Create secrets and variables for an organization and/or repositories.

Usage:
  seva [command]

Available Commands:
  secrets     Export and Create secrets for an organization and/or repositories.
  variables   Export and Create variables for an organization and/or repositories.

Flags:
      --help   Show help for command

Use "seva [command] --help" for more information about a command.
```

### Secrets

The `gh seva secrets` command comprises of two subcommands, `export` and `create`, to access and create Organization level and repository level secrets.

```sh
$ gh seva secrets -h
Export and Create Actions, Dependabot, and Codespaces secrets for an organization and/or repositories.

Usage:
  seva secrets [command]

Available Commands:
  create      Create Actions, Dependabot, and/or Codespaces secrets from a file.
  export      Generate a report of Actions, Dependabot, and Codespaces secrets for an organization and/or repositories.

Flags:
      --help   Show help for command

Use "seva secrets [command] --help" for more information about a command.
```

#### Create Secrets

The `gh seva secrets create` command will create secrets from a `csv` file that contains the following information:

- `SecretLevel`: If the secret was created at the organization or repository level
- `SecretType`: If the secret was created for `Actions`, `Dependabot` or `Codespaces`
- `SecretName`: The name of the secret
-	`SecretValue`: The value of the secret that will be [encrypted using the associated `public key`](https://docs.github.com/en/actions/security-guides/encrypted-secrets)
- `SecretAccess`: If an organization level secret, the visibility of the secret (i.e. `all`, `private`, or `scoped`)
- `RepositoryNames`: The name of the repositories that the secret can be accessed from (delimited with `;`)
- `RepositoryIDs`: The `id` of the repositories that the secret can be accessed from (delimited with `;`)

This extension supports `GitHub.com` and GHES, through the use of `--hostname` and `--token`.

```sh
$ gh seva secrets create -h
Create Actions, Dependabot, and/or Codespaces secrets for an organization and/or repositories from a file.

Usage:
  seva secrets create <organization> [flags]

Flags:
  -d, --debug              To debug logging
  -f, --from-file string   Path and Name of CSV file to create webhooks from (required)
      --hostname string    GitHub Enterprise Server hostname (default "github.com")
  -t, --token string       GitHub personal access token for organization to write to (default "gh auth token")

Global Flags:
      --help   Show help for command
```

#### Export Secrets

The `gh seva secrets export` command exports secrets for the specified `<organization>` or `[repo ..]` list. If `<organization>` is selected, **both organization level and repository level secrets will be exported**. The report will contain secrets produces a `csv` report with the following:

- `SecretLevel`: If the secret was created at the organization or repository level
- `SecretType`: If the secret was created for `Actions`, `Dependabot` or `Codespaces`
- `SecretName`: The name of the secret
-	`SecretValue`: This field **will be blank**, we cannot export secret values.
- `SecretAccess`: If an organization level secret, this is the visibility of the secret (i.e. `all`, `private`, or `scoped`)
- `RepositoryNames`: The name of the repositories that the secret can be accessed from (delimited with `;`)
- `RepositoryIDs`: The `id` of the repositories that the secret can be accessed from (delimited with `;`)

This extension supports `GitHub.com` and GHES, through the use of `--hostname` and `--token`.

```sh
$ gh seva secrets export -h
Generate a report of Actions, Dependabot, and Codespaces secrets for an organization and/or repositories.

Usage:
  seva secrets export [flags] <organization> [repo ...] 

Flags:
  -a, --app string           List secrets for a specific application or all: {all|actions|codespaces|dependabot} (default "all")
  -d, --debug                To debug logging
      --hostname string      GitHub Enterprise Server hostname (default "github.com")
  -o, --output-file string   Name of file to write CSV report (default "report-20230505162601.csv")
  -t, --token string         GitHub Personal Access Token (default "gh auth token")

Global Flags:
      --help   Show help for command
```

### Variables

Organization level Actions variables can be created and exported, relying on the `csv` file syntax:

- `VariableLevel`: If the variable was created at the organization or repository level
- `VariableName`: The name of the Actions variable
- `VariableValue`: The value of the Actions variable
- `VariableAccess`: If an organization level variable, this is the visibility of the variable (i.e. `all`, `private`, or `scoped`)
- `RepositoryNames`: The name of the repositories that the variable can be accessed from (delimited with `;`)
- `RepositoryIDs`: The `id` of the repositories that the variable can be accessed from (delimited with `;`)


```sh
$ gh seva variables -h
Export and Create Actions variables for an organization and/or repositories.

Usage:
  seva variables [command]

Available Commands:
  create      Create Organization Actions variables.
  export      Generate a report of Actions variables for an organization and/or repositories.

Flags:
      --help   Show help for command

Use "seva variables [command] --help" for more information about a command.
```

#### Create Variables

Organization level variables can be created from a `csv` file using `--from-file` following the format outlined in [`gh seva variables`](#variables).

* If specifying a Source Organization (`--source-organization`) to retrieve secrets and create under a new Org, the `--source-token` is required.

```sh
$ gh seva variables create -h

Create Organization Actions variables for a specified organization or organization and repositories level variables from a file.

Usage:
  seva variables create <organization> [flags]

Flags:
  -d, --debug                        To debug logging
  -f, --from-file string             Path and Name of CSV file to create variables from
      --hostname string              GitHub Enterprise Server hostname (default "github.com")
      --source-hostname string       GitHub Enterprise Server hostname where variables are copied from (default "github.com")
  -o, --source-organization string   Name of the Source Organization to copy variables from (Requires --source-token)
  -s, --source-token string          GitHub personal access token for Source Organization (Required for --source-organization)
  -t, --token string                 GitHub personal access token for organization to write to (default "gh auth token")

Global Flags:
      --help   Show help for command
```

#### Export Variables

The `gh seva variables export` command exports variables for the specified `<organization>` or `[repo ..]` list. If `<organization>` is selected, **both organization level and repository level variables will be exported**. The report will contain secrets produces a `csv` report with the following:

- `VariableLevel`: If the variable was created at the organization or repository level
- `VariableName`: The name of the Actions variable
- `VariableValue`: The value of the Actions variable
- `VariableAccess`: If an organization level variable, this is the visibility of the variable (i.e. `all`, `private`, or `scoped`)
- `RepositoryNames`: The name of the repositories that the variable can be accessed from (delimited with `;`)
- `RepositoryIDs`: The `id` of the repositories that the variable can be accessed from (delimited with `;`)

This extension supports `GitHub.com` and GHES, through the use of `--hostname` and `--token`.

```sh
$ gh seva variables export -h
Generate a report of Actions variables for an organization and/or repositories.

Usage:
  seva variables export [flags] <organization> [repo ...] 

Flags:
  -d, --debug                To debug logging
      --hostname string      GitHub Enterprise Server hostname (default "github.com")
  -o, --output-file string   Name of file to write CSV report (default "report-20230505163210.csv")
  -t, --token string         GitHub Personal Access Token (default "gh auth token")

Global Flags:
      --help   Show help for command
```
