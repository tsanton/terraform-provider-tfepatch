# Terraform cloud private provider

The goal of this project is to set up a custom terraform cloud provider registry, build our provider, and release it to terraform cloud where it can be distributed to our authenticated consumers.

Inspirations for this guide:

- [Private module and provider registries](https://developer.hashicorp.com/terraform/cloud-docs/registry)
- [Publish private providers](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers)
- [Registry Providers API](https://developer.hashicorp.com/terraform/cloud-docs/api-docs/private-registry/providers#create-a-provider)
- [API Authentication](https://developer.hashicorp.com/terraform/cloud-docs/api-docs#authentication)

## TL;DR

We're going to complete the following steps:

- clone down the repository
- Move directory into [./internal/](./internal/) and run our `Makefile` (build and move our bootstrap providers into our .terraform.d cache)
- Switch directory into [./terraform/](./terraform/) where we are to create `GPG key` and `provider registry`.
- Push the repository to your Git (enable XXXX) and trigger our pipeline
- Authenticate to Terraform Cloud and init our new CI/CD provider platform version

## Bootstrap

### Terraform Cloud API Token

Firstly you must sign into your [terraform cloud organization](https://app.terraform.io/) and create API keys.

There are three different types of keys, of whom all three types can be used:

- User tokens — each Terraform Cloud user can have any number of API tokens, which can make requests on their behalf.
- Team tokens — each team can have one API token at a time. This is intended for performing plans and applies via a CI/CD pipeline.
- Organization tokens — each organization can have one API token at a time. This is intended for automating the management of teams, team membership, and workspaces

For the CI/CD of our provider, we chose to utilize the `team token` from the owners team.

FYI: You must be a member of the owners team or a team with Manage Private Registry permissions to publish and delete private providers from the private registry.

### Terraform Cloud GPG key

First run this [makefile](./internal/Makefile) to create a local copy (0.0.1) of the provider. \
We will utilise this provider to upload the public GPG key to the Terraform Cloud Registry, and to create our registry provider platform.

Navigate into the [./terraform](./terraform/) directory and supply the following key value pairs in a *providers.auto.tfvars* file:

- tfe_token="very.long.token"
- tfe_organization="my-org-name"

As you'll see from the resource configuration, for private providers the namespace must equal the organization name. Completing the same logic: the GPG key must match with the namespace of the provider.

Do not `terrform init` just yet. You must first complete the next step.

### GitHub Action

Creating the GPG-key and the provider platform gives you a base to create releases out from. \
To generate versioned releases we will use [this](./.github/workflows/release.yml) release pipeline. \
For it to run we must first upload the following secret and variables to your GitHub respository:

- secrets.TFE_TOKEN
- secrets.GPG_PRIVATE_KEY
- secrets.PASSPHRASE
- vars.TFE_ORGANIZATION
- vars.TFE_NAMESPACE
- vars.TFE_PROVIDER_NAME
- vars.TFE_GPG_KEY_ID

These secrets and variables will be uploaded to your repository if you utilise the included [terraform boostrap](./terraform/) code. \
If you don't have access to create a GitHub personal access token (PAT) or have access to a GitHub Application, feel free to remote the [github.tf](./terraform/github.tf) script, remove the [provider](./terraform/providers.tf) configuration for the `GitHub provider`, remove the [variable](./terraform/variables.tf) references and manually upload your secrets.

If you opt to upload the secrets and variables to GitHub through terraform, navigate to the [./terraform](./terraform/) directory and supply the following additional key value pairs in a *providers.auto.tfvars* file:

- github_repository_name="name-of-forked-repo"
- github_token="my-personal-access-token-with-read-write-on-repository-and-workflows"

This is the point in time where we go `terraform apply`

## Test and build the provider

This provider comes with integration tests, but they're not run in the release pipeline. Feel free to add that step for yourself.

## Release to Terraform cloud

As a last step prior to fireing off our release we must configure our repo to allow `[...] whether GitHub Actions can create pull requests or submit approving pull requests reviews`: we want this!

We're almost ready to commit and generate our first release. Note that the release entails a few steps:

- We use [release-please-actions](https://github.com/google-github-actions/release-please-action) to generate semantic versioned tag based off [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)
- Use [goreleaser](https://github.com/goreleaser/goreleaser-action) to create and upload the release assets to your tag
- Use [this](https://github.com/Tsanton/tfe-provider-release-action) composite action to create a new provider version and to upload your provider binaries.

Create a commit with a conventional commit message along the line of `feat: init release` and push your code. This should trigger the release of your provider.
While you read up on what each of those actions do individually, rest assured that your provider is being released as we speak.

## Authentication for consumption

In order to use a remote published artifacts, we must authenticate to our Terraform Cloud Organization. \
To do so, we can create a .terraformrc holding the terraform API token: See [this](https://developer.hashicorp.com/terraform/cli/config/config-file) doc for where to place your `.terraformrc` file.

```hcl
credentials "app.terraform.io" {
 token = "dz8xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxK0Y"
}
```

I find that mounting or using symlink to place the `.terraformrc` file into the users home (`~/`) directory is the easiest way to ensure access to the token for authentication purposes.

## Extra: Generate Documentation for Wiki release?

First you must create examples of how to use your resources and data sources. This must be done in the [<package-root>./examples](./internal/examples/) directory. \
The expected folder structure and file names are as follows:

```txt
└── examples/
    ├── resources/
    │ ├── <provider-name>_<resource-name>/
    │ │ ├── resource.tf
    │ │ └── import.sh
    │ └── <provider-name>_another_resource/
    │ ├── resource.tf
    │ └── import.sh
    ├── data-sources/
    │ ├── <provider-name>_<data-source-name>/
    │ │ ├── resource.tf
    │ │ └── import.sh
    │ └── <provider-name>_another_data_source/
    │ ├── resource.tf
    │ └── import.sh
    └── provider/
    └── provider.tf
```

In order to automate the docs generation we must complete the following steps:

- Setup your ./examples folder according to the structure above and fill the files (.tf and .sh) with resource usage definition and example import statements.
- Run ```go get -d github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs```
- Create a *tools* folder with a *tools.go* in your go.mod working directory
  - Add the following to [tools.go](./internal/tools/tools.go):

    ```go
    //go:build tools
    package tools

    import (
        // Documentation generation
        _ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
    )
    ```

- Modify your [main.go](./internal/main.go) to include the following comments in your main.go:

    ```go
    //Generate the Terraform provider documentation using `tfplugindocs`:
    //go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name <**insert-your-provider-name**>
    ```

- Ensure the correct 'GOOS' and 'GOARCH' is set as environment variables
- From your project root, run ```GOOS=linux GOARCH=amd64 go generate ./...```

This created a *./docs* output with merged information from your examples folder and the resource/attribute descriptions. \
This *docs* folder can then be released in order to provide proper documentation of the usage of your custom provider.

See [this](https://developer.hashicorp.com/terraform/tutorials/providers/provider-release-publish#generate-provider-documentation) for more info about how to generate docs.
