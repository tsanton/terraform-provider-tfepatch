provider "tfepatch" {
  hostname        = "https://app.terraform.io"
  token           = "abcde12345"
  organization    = "gruntwork-corp"
  ssl_skip_verify = false
}
