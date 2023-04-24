terraform {
  required_providers {
    github = {
      source  = "integrations/github"
      version = "5.22.0"
    }
    gpg = {
      source  = "Olivr/gpg"
      version = "0.2.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.5.1"
    }
    /* Run the makefile in the ./internal directory to access this provider*/
    tfepatch = {
      source  = "bootstrap/automation/tfepatch"
      version = "0.0.1"
    }
    /* After release (replace '<org_name>' with your organization name */
    # tfepatch = {
    #   source  = "app.terraform.io/<org_name>/tfepatch"
    #   version = "0.1.0"
    # }
  }
}

provider "github" {
  token = var.github_token
}

provider "gpg" {}

provider "random" {}

provider "tfepatch" {
  hostname     = "https://app.terraform.io"
  token        = var.tfe_token
  organization = var.tfe_organization
}


