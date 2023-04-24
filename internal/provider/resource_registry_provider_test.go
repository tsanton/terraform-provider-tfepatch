package provider_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	u "github.com/tsanton/terraform-provider-tfepatch/utilities"

	me "github.com/tsanton/tfe-client/tfe/models/enum"
)

func Test_provider_public_registry_provider(t *testing.T) {
	/* Arrange */
	log.Println("Arranging")
	orgName := u.GetEnv("TFE_ORG_NAME", "")
	namespace := "hashicorp"
	name := "aws"
	registryType := me.RegistryTypePublic

	/* Act */
	log.Println("Invoking tests")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//--------------------------------------------------------------------------
			//--- Create and Read testing
			//--------------------------------------------------------------------------
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tfepatch_registry_provider" "this" {
					organization  = "%s"
					namespace     = "%s"
					name          = "%s"
					registry_name = "%s"
				  }
				`, orgName, namespace, name, registryType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "organization", orgName),
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "namespace", namespace),
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "name", name),
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "registry_name", string(registryType)),
				),
			},
			{
				RefreshState: true,
				// Destroy: true,
				PreConfig: func() {
					registryProvider, err := cli.ProviderService.Read(context.Background(), orgName, string(registryType), namespace, name)
					assert.Nil(t, err)
					assert.Equal(t, name, registryProvider.Data.Attributes.Name)
					assert.Equal(t, namespace, registryProvider.Data.Attributes.Namespace)
					assert.Equal(t, registryType, registryProvider.Data.Attributes.RegistryName)
				},
			},
		},
	})
}

func Test_provider_private_registry_provider(t *testing.T) {
	/* Arrange */
	log.Println("Arranging")
	orgName := u.GetEnv("TFE_ORG_NAME", "")
	namespace := orgName
	name := "demo-provider"
	registryType := me.RegistryTypePrivate

	/* Check func */

	/* Act */
	log.Println("Invoking tests")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			//--------------------------------------------------------------------------
			//--- Create and Read testing
			//--------------------------------------------------------------------------
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "tfepatch_registry_provider" "this" {
					organization  = "%s"
					namespace     = "%s"
					name          = "%s"
					registry_name = "%s"
				  }
				`, orgName, namespace, name, registryType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "organization", orgName),
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "namespace", namespace),
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "name", name),
					resource.TestCheckResourceAttr("tfepatch_registry_provider.this", "registry_name", string(registryType)),
				),
			},
			{
				RefreshState: true,
				// Destroy: true,
				PreConfig: func() {
					registryProvider, err := cli.ProviderService.Read(context.Background(), orgName, string(registryType), namespace, name)
					assert.Nil(t, err)
					assert.Equal(t, name, registryProvider.Data.Attributes.Name)
					assert.Equal(t, namespace, registryProvider.Data.Attributes.Namespace)
					assert.Equal(t, registryType, registryProvider.Data.Attributes.RegistryName)
				},
			},
		},
	})
}
