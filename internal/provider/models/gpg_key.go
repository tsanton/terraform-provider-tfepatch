package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type GpgKey struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	Namespace    types.String `tfsdk:"namespace"`
	PublicKey    types.String `tfsdk:"public_key"`
	KeyId        types.String `tfsdk:"key_id"`
}
