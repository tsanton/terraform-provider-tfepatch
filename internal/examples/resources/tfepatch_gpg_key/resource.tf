terraform {
  required_providers {
    gpg = {
      source  = "Olivr/gpg"
      version = "0.2.1"
    }
  }
}

provider "gpg" {}

resource "gpg_private_key" "this" {
  name       = "Gruntwork"
  email      = "donotreply@gruntwork.com"
  passphrase = "FooBarBaz"
  rsa_bits   = 4096
}

resource "tfepatch_gpg_key" "this" {
  organization = var.organization_name
  namespace    = "gruntwork-corp"
  public_key   = gpg_private_key.this.public_key
}
