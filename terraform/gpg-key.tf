resource "random_password" "this" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "gpg_private_key" "this" {
  name       = "gruntwork-corp"
  email      = "tobias@tsant.no"
  passphrase = random_password.this.result
  rsa_bits   = 4096
}

resource "tfepatch_gpg_key" "this" {
  organization = tfepatch_registry_provider.this.organization
  namespace    = tfepatch_registry_provider.this.namespace
  public_key   = gpg_private_key.this.public_key
}
