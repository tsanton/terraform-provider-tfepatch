package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	m "github.com/tsanton/terraform-provider-tfepatch/provider/models"
	api "github.com/tsanton/tfe-client/tfe"
	apir "github.com/tsanton/tfe-client/tfe/models/request"
)

type GpgKeyResource struct {
	client *api.TerraformEnterpriseClient
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &GpgKeyResource{}
	_ resource.ResourceWithConfigure   = &GpgKeyResource{}
	_ resource.ResourceWithImportState = &GpgKeyResource{}
)

// newResource is a helper function to simplify the provider implementation.
func newGpgKeyResource() resource.Resource {
	return &GpgKeyResource{}
}

// Configure adds the provider configured client to the resource.
func (r *GpgKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.TerraformEnterpriseClient)
}

func (r *GpgKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique id for this resource",
				MarkdownDescription: "Unique id for this resource",
			},
			"key_id": schema.StringAttribute{
				Computed:            true,
				Description:         "The identity of the generated key",
				MarkdownDescription: "The identity of the generated key",
			},
			// Input attributes
			"organization": schema.StringAttribute{
				Required:            true,
				Description:         "The organization name under which this GPG key will exist",
				MarkdownDescription: "The organization name under which this GPG key will exist",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				Required:            true,
				Description:         "The provider, by namespace, that this GPG key is affiliated with.",
				MarkdownDescription: "The provider, by namespace, that this GPG key is affiliated with.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Required:            true,
				Description:         "The ASCII-armored public GPG key ",
				MarkdownDescription: "The ASCII-armored public GPG key ",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *GpgKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gpg_key"
}

func (r *GpgKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan m.GpgKey
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr, err := r.client.GpgService.Create(ctx, &apir.Gpg{
		Data: apir.GpgData{
			Type: "gpg-keys",
			Attributes: apir.GpgDataAttributes{
				AsciiArmor: plan.PublicKey.ValueString(),
				Namespace:  plan.Namespace.ValueString(),
			},
		},
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource",
			"Could not read resource "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf(strings.ToLower("%s||%s"), cr.Data.Attributes.Namespace, cr.Data.Attributes.KeyId))
	plan.KeyId = types.StringValue(cr.Data.Attributes.KeyId)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *GpgKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state m.GpgKey
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rr, err := r.client.GpgService.Read(ctx, state.Namespace.ValueString(), state.KeyId.ValueString())

	//TODO: Handle deleted resource read
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource",
			"Could not read resource "+err.Error(),
		)
		return
	}

	state = m.GpgKey{
		Id:           types.StringValue(fmt.Sprintf(strings.ToLower("%s||%s"), rr.Data.Attributes.Namespace, rr.Data.Attributes.KeyId)),
		Organization: state.Organization,
		Namespace:    types.StringValue(rr.Data.Attributes.Namespace),
		PublicKey:    types.StringValue(rr.Data.Attributes.AsciiArmor),
		KeyId:        types.StringValue(rr.Data.Attributes.KeyId),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// No update as all attributes require replacement if changed
func (r *GpgKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *GpgKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state m.GpgKey
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.GpgService.Delete(ctx, state.Namespace.ValueString(), state.KeyId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource",
			"Could not read resource "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// TODO: test import
func (r *GpgKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to the attributes that are utilized by the Read-function
	parts := strings.Split(req.ID, "||")
	resp.State.SetAttribute(ctx, path.Root("namespace"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("key_id"), parts[1])
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
