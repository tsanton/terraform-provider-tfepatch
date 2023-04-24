package provider

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	tfeclient "github.com/tsanton/tfe-client/tfe"

	log "github.com/sirupsen/logrus"

	m "github.com/tsanton/tfe-client/tfe/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &TfeProvider{}
)

type TfeProvider struct{}

// First, the provider calls the New function to create an instance of the provider
func New() func() provider.Provider {
	return func() provider.Provider {
		return &TfeProvider{}
	}
}

func (p *TfeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tfepatch"
}

func (p *TfeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Terraform provider to patch/work alongside the TFE for resources/data sources that are yet to be implemented",
		MarkdownDescription: "Terraform provider to patch/work alongside the TFE for resources/data sources that are yet to be implemented",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Required:            true,
				Sensitive:           false,
				Description:         "The Terraform Enterprise hostname to connect to. Defaults to app.terraform.io.",
				MarkdownDescription: "The Terraform Enterprise hostname to connect to. Defaults to app.terraform.io.",
			},
			"token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				Description:         "The token used to authenticate with Terraform Enterprise. We recommend omitting the token which can be set as credentials in the CLI config file.",
				MarkdownDescription: "The token used to authenticate with Terraform Enterprise. We recommend omitting the token which can be set as credentials in the CLI config file.",
			},
			"organization": schema.StringAttribute{
				Required:            true,
				Sensitive:           false,
				Description:         "The organization to apply to a resource if one is not defined on the resource itself.",
				MarkdownDescription: "The organization to apply to a resource if one is not defined on the resource itself.",
			},
			"ssl_skip_verify": schema.BoolAttribute{
				Optional:            true,
				Sensitive:           false,
				Description:         "Whether or not to skip certificate verifications.",
				MarkdownDescription: "Whether or not to skip certificate verifications.",
			},
		},
	}
}

type providerConfig struct {
	Hostname     types.String `tfsdk:"hostname"`
	Token        types.String `tfsdk:"token"`
	Organization types.String `tfsdk:"organization"`
	VerifyTls    types.Bool   `tfsdk:"ssl_skip_verify"`
}

func (p *TfeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	logger := log.New()
	logger.Info("Configuring TFE client")

	// Retrieve provider data from configuration
	var config providerConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//TODO: verify ENV variable function
	if config.Hostname.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hostname"),
			"Unknown Terraform Enterprice API Host",
			"The provider cannot create the Terraform Enterprice API client as there is an unknown configuration value for the API host. "+
				"Either set the value statically in the configuration, or use the provider_HOST environment variable.",
		)
	}

	//TODO: verify ENV variable function
	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Terraform Enterprice API Token",
			"The provider cannot create the Terraform Enterprice API client as there is an unknown configuration value for the API token. "+
				"Either set the value statically in the configuration, or use the provider_TOKEN environment variable.",
		)
	}

	//TODO: verify ENV variable function
	if config.Organization.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Terraform Enterprice API Token",
			"The provider cannot create the Terraform Enterprice API client as there is an unknown configuration value for the API token. "+
				"Either set the value statically in the configuration, or use the provider_TOKEN environment variable.",
		)
	}

	transport := &http.Transport{}
	if transport.TLSClientConfig == nil {
		transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	if !config.VerifyTls.ValueBool() {
		tflog.Warn(ctx, "Client configured to skip certificate verifications")
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
	httpClient := cleanhttp.DefaultPooledClient()
	httpClient.Transport = transport

	// Create a new TFE client config
	cfg := m.ClientConfig{
		Address: config.Hostname.ValueString(),
		Token:   config.Token.ValueString(),
	}

	// Create a new TFE client.
	client, err := tfeclient.NewClient(logger, &cfg)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to configure up a new TFE API Client",
			"Unable to configure up a new TFE API Client",
		)
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

// GetDataSources satisfies the provider.Provider interface.
func (p *TfeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Provider specific implementation
	}
}

// GetResources satisfies the provider.Provider interface.
func (p *TfeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newProviderRegistryResource,
		newGpgKeyResource,
	}
}
