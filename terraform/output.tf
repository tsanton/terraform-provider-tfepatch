output "tfe_config" {
  value = {
    "tfe_organization"  = tfepatch_registry_provider.this.organization
    "tfe_namespace"     = tfepatch_registry_provider.this.namespace
    "tfe_provider_name" = tfepatch_registry_provider.this.name
    "tfe_gpg_key_id"    = tfepatch_gpg_key.this.key_id
  }
}
