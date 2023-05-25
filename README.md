# **Terraform cloud private provider**

The goal of this project is to set up a custom terraform cloud provider registry, build our provider, and release it to terraform cloud where it can be distributed to our authenticated consumers.

Inspirations for this guide:

- [Private module and provider registries](https://developer.hashicorp.com/terraform/cloud-docs/registry)
- [Publish private providers](https://developer.hashicorp.com/terraform/cloud-docs/registry/publish-providers)
- [Registry Providers API](https://developer.hashicorp.com/terraform/cloud-docs/api-docs/private-registry/providers#create-a-provider)
- [API Authentication](https://developer.hashicorp.com/terraform/cloud-docs/api-docs#authentication)

## **TL;DR**

We're going to complete the following steps:

- clone down the repository
- Move directory into [./internal/](./internal/) and run our `Makefile` (build and move our bootstrap providers into our .terraform.d cache)
- Switch directory into [./terraform/](./terraform/) where we are to create `GPG key` and `provider registry`.
- Push the repository to your Git (enable XXXX) and trigger our pipeline
- Authenticate to Terraform Cloud and init our new CI/CD provider platform version

## **Bootstrap**

### **Terraform Cloud API Token**

Firstly you must sign into your [terraform cloud organization](https://app.terraform.io/) and create API keys.

There are three different types of keys, of whom all three types can be used:

- User tokens — each Terraform Cloud user can have any number of API tokens, which can make requests on their behalf.
- Team tokens — each team can have one API token at a time. This is intended for performing plans and applies via a CI/CD pipeline.
- Organization tokens — each organization can have one API token at a time. This is intended for automating the management of teams, team membership, and workspaces

For the CI/CD of our provider, we chose to utilize the `team token` from the owners team.

FYI: You must be a member of the owners team or a team with Manage Private Registry permissions to publish and delete private providers from the private registry.

### **Terraform Cloud GPG key**

First run this [makefile](./internal/Makefile) to create a local copy (0.0.1) of the provider. We will utilise this provider to generate a GPG key, upload that key to the Terraform Cloud Registry, and to create our registry provider platform.

Navigate into the [./terraform](./terraform/) directory and supply the following key value pairs in a *providers.auto.tfvars* file:

- tfe_token="very.long.token"
- tfe_organization="my-org-name"

As you'll see from the resource configuration, for private providers the namespace must equal the organization name. Completing the same logic: the GPG key must match with the namespace of the provider.

Do not `terrform init` just yet. You must first complete the next step.

### **GitHub Action**

For [this](./.github/workflows/release.yml) release pipeline to run you must first upload the following secret and variables to your GitHub respository:

- secrets.TFE_TOKEN
- secrets.GPG_PRIVATE_KEY
- secrets.PASSPHRASE
- vars.TFE_ORGANIZATION
- vars.TFE_NAMESPACE
- vars.TFE_PROVIDER_NAME
- vars.TFE_GPG_KEY_ID

All these secrets and variables will be uploaded to your repository as long as you utilise the included [terraform boostrap](./terraform/) code.

Now again, navigate into the [./terraform](./terraform/) directory and supply the following additional key value pairs in a *providers.auto.tfvars* file:

- github_repository_name="name-of-forked-repo"
- github_token="my-personal-access-token-with-read-write-on-repository-and-workflows"

This is the point in time where we go `terraform apply`

**NB! Store the your terraform state away in a secure remote backend**.

## **Test and build the provider**

This provider comes with integration tests, but they're not run in the release pipeline. Feel free to add that step for yourself.

## **Release to Terraform cloud**

As a last step prior to fireing off our release we must configure our repo to `choose whether GitHub Actions can create pull requests or submit approving pull requests reviews`; naturally we want this!

When this is in order, commit with a conventional commit message along the line of `feat: init release` and push your code. This should trigger the release of your provider.

The release entails a few steps:

- Use [release-please-actions](https://github.com/google-github-actions/release-please-action) to generate semantic versioned tag based off [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)
- Use [goreleaser](https://github.com/goreleaser/goreleaser-action) to create and upload the release assets to your tag
- Use [this](https://github.com/Tsanton/tfe-provider-release-action) composite action to create a new provider version and to upload your provider binaries.

While you read up on what each of those actions do individually, rest assured that your provider is being released as we speak.

## **Authentication for consumption**

In order to use a remote published artifacts, we must authenticate to our Terraform Cloud Organization. \
To do so, we can create a .terraformrc holding the terraform API token: See [this](https://developer.hashicorp.com/terraform/cli/config/config-file) doc for where to place your `.terraformrc` file.

```hcl
credentials "app.terraform.io" {
 token = "dz8xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxK0Y"
}
```

As a sidenote I usually "cheat" mount/symlink the `.terraformrc` file into the users home (`~/`) directory.
That our you can init with the following command:

```sh
TF_CLI_CONFIG_FILE="/path/to/.terraformrc" terraform init
```

## **Extra: Generate Documentation for Wiki release?**

See [this](https://developer.hashicorp.com/terraform/tutorials/providers/provider-release-publish#generate-provider-documentation) for how to automatically generate docs.
In short it boils down to completing the following steps:

- run ```go get -d github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs```
- Create a *tools* folder with a *tools.go* in your go.mod working directory
- Modify your main.go to include *generate*
- Ensure the correct 'GOOS' and 'GOARCH' is set as environment variables
- From your project root, run ```go generate ./...```

This created a *./docs* output with merged information from your examples folder. \
See [main.go](./internal/main.go) for generate config.
