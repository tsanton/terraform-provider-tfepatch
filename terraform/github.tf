data "github_repository" "this" {
  name = var.github_repository_name
}

resource "github_actions_variable" "organization" {
  repository    = data.github_repository.this.name
  variable_name = "TFE_ORGANIZATION"
  value         = tfepatch_registry_provider.this.organization
}

resource "github_actions_variable" "namespace" {
  repository    = data.github_repository.this.name
  variable_name = "TFE_NAMESPACE"
  value         = tfepatch_registry_provider.this.namespace
}

resource "github_actions_variable" "provider_name" {
  repository    = data.github_repository.this.name
  variable_name = "TFE_PROVIDER_NAME"
  value         = tfepatch_registry_provider.this.name
}

resource "github_actions_variable" "gpg_key_id" {
  repository    = data.github_repository.this.name
  variable_name = "TFE_GPG_KEY_ID"
  value         = tfepatch_gpg_key.this.key_id
}

resource "github_actions_secret" "passphrase" {
  repository      = data.github_repository.this.name
  secret_name     = "PASSPHRASE"
  plaintext_value = gpg_private_key.this.passphrase
}

resource "github_actions_secret" "gpg_private_key" {
  repository      = data.github_repository.this.name
  secret_name     = "GPG_PRIVATE_KEY"
  plaintext_value = gpg_private_key.this.private_key
}

resource "github_actions_secret" "tfe_token" {
  repository      = data.github_repository.this.name
  secret_name     = "TFE_TOKEN"
  plaintext_value = var.tfe_token
}
