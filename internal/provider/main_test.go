package provider_test

import (
	"fmt"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	u "github.com/tsanton/terraform-provider-tfepatch/utilities"
	api "github.com/tsanton/tfe-client/tfe"
	apim "github.com/tsanton/tfe-client/tfe/models"
)

var providerConfig string
var cli *api.TerraformEnterpriseClient
var logger u.ILogger

func TestMain(m *testing.M) {
	logger = log.New()

	logger.Info("Test suite setup")

	tfeToken := u.GetEnv("TFE_TOKEN", "")
	providerConfig = fmt.Sprintf(`
provider "tfepatch" {
	hostname  		= "https://app.terraform.io"
	token 			= "%s"
	organization	= "%s"
	ssl_skip_verify = false

}`, tfeToken, u.GetEnv("TFE_ORG_NAME", ""))

	var err error
	cli, err = api.NewClient(logger, &apim.ClientConfig{
		Address: "https://app.terraform.io",
		Token:   tfeToken,
	})
	if err != nil {
		panic("unable to create test client")
	}

	logger.Info("Invoking tests")
	exitVal := m.Run()

	logger.Info("Test suite teardown")

	os.Exit(exitVal)
}
