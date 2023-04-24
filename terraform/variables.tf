variable "github_repository_name" {
  type        = string
  description = "The github "
  sensitive   = false
}

variable "github_token" {
  type        = string
  description = "A github PAT that has access to the 'var.github_repository_name' repo"
  sensitive   = true
}

variable "tfe_token" {
  type        = string
  description = "The team or personal access token to your Terraform Cloud organization. NB: You must be a member of the owners team or a team with Manage Private Registry permissions to publish and delete private providers from the private registry"
  sensitive   = true
}

variable "tfe_organization" {
  type        = string
  description = "The name of your Terraform Cloud organization"
  sensitive   = false
}
