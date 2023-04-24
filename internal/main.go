package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tsanton/terraform-provider-tfepatch/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name tfepatch

func main() {
	providerHostname := "tsanton"
	namespace := "gruntwork-corp"
	providerName := "tfepatch"
	opts := providerserver.ServeOpts{
		Address: fmt.Sprintf("%s/%s/%s", providerHostname, namespace, providerName),
		Debug:   false,
	}

	err := providerserver.Serve(context.Background(), provider.New(), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
