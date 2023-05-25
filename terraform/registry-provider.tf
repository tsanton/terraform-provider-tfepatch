resource "tfepatch_registry_provider" "this" {
  organization  = var.tfe_organization
  namespace     = var.tfe_organization
  name          = "tfepatch"
  registry_name = "private"
}
