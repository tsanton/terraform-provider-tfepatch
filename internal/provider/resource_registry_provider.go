package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	m "github.com/tsanton/terraform-provider-tfepatch/provider/models"
	api "github.com/tsanton/tfe-client/tfe"
	"github.com/tsanton/tfe-client/tfe/models/enum"
	apir "github.com/tsanton/tfe-client/tfe/models/request"
)

type ProviderRegistryResource struct {
	client *api.TerraformEnterpriseClient
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ProviderRegistryResource{}
	_ resource.ResourceWithConfigure   = &ProviderRegistryResource{}
	_ resource.ResourceWithImportState = &ProviderRegistryResource{}
)

// newResource is a helper function to simplify the provider implementation.
func newProviderRegistryResource() resource.Resource {
	return &ProviderRegistryResource{}
}

// Configure adds the provider configured client to the resource.
func (r *ProviderRegistryResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.TerraformEnterpriseClient)
}

func (r *ProviderRegistryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique id for this resource",
				MarkdownDescription: "Unique id for this resource",
			},
			// Input attributes
			"organization": schema.StringAttribute{
				Required:            true,
				Description:         "The organization name under which this provider registry will exist",
				MarkdownDescription: "The organization name under which this provider registry will exist",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"namespace": schema.StringAttribute{
				Required:            true,
				Description:         "The namespace under which the provider will exist. For provider registries with 'registy_name' == 'private' the namespace must match organization name", //TODO: test that namespace equals organization name when registry_name is private
				MarkdownDescription: "The namespace under which the provider will exist. For provider registries with 'registy_name' == 'private' the namespace must match organization name", //TODO: test that namespace equals organization name when registry_name is private
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the provider",
				MarkdownDescription: "The name of the provider",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"registry_name": schema.StringAttribute{
				Required:            true,
				Description:         "The registry type for the provider. Must be 'public' or 'private'",
				MarkdownDescription: "The registry type for the provider. Must be 'public' or 'private'",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("private", "public"),
				},
			},
		},
	}
}

// Metadata returns the resource type name.
func (r *ProviderRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry_provider"
}

func (r *ProviderRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan m.RegistryProvider
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cr, err := r.client.ProviderService.Create(ctx, plan.Organization.ValueString(), &apir.Provider{
		Data: apir.ProviderData{
			Type: "registry-providers",
			Attributes: apir.ProviderDataAttributes{
				Name:         plan.Name.ValueString(),
				Namespace:    plan.Namespace.ValueString(),
				RegistryName: enum.RegistryType(plan.RegistryName.ValueString()),
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

	plan.Id = types.StringValue(fmt.Sprintf(strings.ToLower("%s||%s||%s"), cr.Data.Attributes.Namespace, cr.Data.Attributes.Name, cr.Data.Attributes.RegistryName))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ProviderRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state m.RegistryProvider
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rr, err := r.client.ProviderService.Read(ctx, state.Organization.ValueString(), state.RegistryName.ValueString(), state.Namespace.ValueString(), state.Name.ValueString())
	//TODO: Handle deleted resource read
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource",
			"Could not read resource "+err.Error(),
		)
		return
	}

	state = m.RegistryProvider{
		Id:           types.StringValue(fmt.Sprintf(strings.ToLower("%s||%s||%s"), rr.Data.Attributes.Namespace, rr.Data.Attributes.Name, rr.Data.Attributes.RegistryName)),
		Organization: state.Organization,
		Namespace:    types.StringValue(rr.Data.Attributes.Namespace),
		Name:         types.StringValue(rr.Data.Attributes.Name),
		RegistryName: types.StringValue(string(rr.Data.Attributes.RegistryName)),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// No update as all attributes require replacement if changed
func (r *ProviderRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *ProviderRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state m.RegistryProvider
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.ProviderService.Delete(ctx, state.Organization.ValueString(), state.RegistryName.ValueString(), state.Namespace.ValueString(), state.Name.ValueString())
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
func (r *ProviderRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to the attributes that are utilized by the Read-function
	parts := strings.Split(req.ID, "||")
	resp.State.SetAttribute(ctx, path.Root("namespace"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("name"), parts[1])
	resp.State.SetAttribute(ctx, path.Root("registry_name"), parts[2])
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
