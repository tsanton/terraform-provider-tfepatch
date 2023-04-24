package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type RegistryProvider struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Namespace    types.String `tfsdk:"namespace"`
	Name         types.String `tfsdk:"name"`
	RegistryName types.String `tfsdk:"registry_name"`
}
