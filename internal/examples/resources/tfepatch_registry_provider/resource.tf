# Manage registry provider.
resource "tfepatch_registry_provider" "this" {
  organization  = var.organization_name
  namespace     = "hashicorp"
  name          = "aws"
  registry_name = "public"
}
