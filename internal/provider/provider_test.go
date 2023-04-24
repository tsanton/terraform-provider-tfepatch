package provider_test

import (
	provider "github.com/tsanton/terraform-provider-tfepatch/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
	// The factory function will be invoked for every Terraform CLI command executed to create a provider server to which the CLI can reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"tfepatch": providerserver.NewProtocol6WithError(provider.New()()),
	}
	testAccExternalProviderFactories = map[string]resource.ExternalProvider{
		"gpg": {
			Source:            "Olivr/gpg",
			VersionConstraint: "0.2.1",
		},
	}
)
