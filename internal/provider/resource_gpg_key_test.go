package provider_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	u "github.com/tsanton/terraform-provider-tfepatch/utilities"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

func Test_provider_gpg_key(t *testing.T) {
	/* Arrange */
	t.Log("Arranging")
	orgName := u.GetEnv("TFE_ORG_NAME", "")
	namespace := orgName
	entity, err := openpgp.NewEntity("Gruntwork", "Integration test GPG key", "donotreply@gruntwork.com", &packet.Config{RSABits: 4096})
	if err != nil {
		t.Log("Unable to generate GPG key entity")
		t.FailNow()
	}
	publicKey, err := generateGpgKey(entity)
	if err != nil {
		t.Log("Unable to generate GPG key")
		t.FailNow()
	}

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
				resource "tfepatch_gpg_key" "this" {
					organization = "%s"
					namespace    = "%s"
					public_key   = trimspace(<<EOF
%s
EOF
)
				  }
				`, orgName, namespace, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfepatch_gpg_key.this", "organization", orgName),
					resource.TestCheckResourceAttr("tfepatch_gpg_key.this", "namespace", namespace),
					resource.TestCheckResourceAttr("tfepatch_gpg_key.this", "public_key", publicKey),
				),
			},
			{
				RefreshState: true,
				// Destroy: true,
				PreConfig: func() {
					keys, err := cli.GpgService.List(context.Background(), []string{namespace})
					assert.Nil(t, err)
					assert.Equal(t, 1, len(keys.Data))
				},
			},
		},
	})
}

func Test_provider_gpg_key_with_gpg_provider(t *testing.T) {
	/* Arrange */
	t.Log("Arranging")
	orgName := u.GetEnv("TFE_ORG_NAME", "")
	namespace := orgName

	/* Act */
	log.Println("Invoking tests")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders:        testAccExternalProviderFactories,
		Steps: []resource.TestStep{
			//--------------------------------------------------------------------------
			//--- Create and Read testing
			//--------------------------------------------------------------------------
			{
				Config: providerConfig + fmt.Sprintf(`
				provider "gpg" {}

				resource "gpg_private_key" "this" {
					name       = "Gruntwork"
					email      = "donotreply@gruntwork.com"
					passphrase = "FooBarBaz"
					rsa_bits   = 4096
				  }

				resource "tfepatch_gpg_key" "this" {
					organization = "%s"
					namespace    = "%s"
					public_key   = gpg_private_key.this.public_key
				  }
				`, orgName, namespace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tfepatch_gpg_key.this", "organization", orgName),
					resource.TestCheckResourceAttr("tfepatch_gpg_key.this", "namespace", namespace),
					// resource.TestCheckResourceAttrPair("tfepatch_gpg_key.this", "public_key", "gpg_private_key.this", "public_key"),
				),
			},
			{
				RefreshState: true,
				// Destroy: true,
				PreConfig: func() {
					keys, err := cli.GpgService.List(context.Background(), []string{namespace})
					assert.Nil(t, err)
					assert.Equal(t, 1, len(keys.Data))
				},
			},
		},
	})
}

func generateGpgKey(entity *openpgp.Entity) (string, error) {
	var publicKeyBuf bytes.Buffer
	err := entity.Serialize(&publicKeyBuf)
	if err != nil {
		fmt.Println("Error serializing public key:", err)
		return "", err
	}

	// Convert the public key to an armored string
	publicKeyArmorBuf := bytes.Buffer{}
	w, err := armor.Encode(&publicKeyArmorBuf, "PGP PUBLIC KEY BLOCK", nil)
	if err != nil {
		fmt.Println("Error encoding public key:", err)
		return "", err
	}
	_, err = w.Write(publicKeyBuf.Bytes())
	if err != nil {
		fmt.Println("Error writing public key to armored buffer:", err)
		return "", err
	}
	w.Close()

	return publicKeyArmorBuf.String(), nil
}
