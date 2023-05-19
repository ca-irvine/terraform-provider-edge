package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/ca-irvine/terraform-provider-edge/internal/model"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	headerKeyID       = "X-API-KEY-ID"
	headerKey         = "X-API-KEY"
	headerUA          = "User-Agent"
	headerContentType = "Content-Type"
)

const (
	applicationJSON = "application/json"
)

var _ provider.Provider = &EdgeProvider{}

type EdgeProvider struct {
	version string
	config  *config
}

type edgeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKeyID types.String `tfsdk:"api_key_id"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *EdgeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "edge"
}

func (p *EdgeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Edge.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
			},
			"api_key_id": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"api_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *EdgeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg edgeProviderModel
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Edge API Endpoint",
			"The provider cannot create the Edge API client as there is an unknown configuration value for the Edge API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDGE_ENDPOINT environment variable.",
		)
	}

	if cfg.APIKeyID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key_id"),
			"Unknown Edge API Key ID",
			"The provider cannot create the Edge API client as there is an unknown configuration value for the Edge API Key ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDGE_API_KEY_ID environment variable.",
		)
	}

	if cfg.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Edge API Key",
			"The provider cannot create the Edge API client as there is an unknown configuration value for the Edge API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the EDGE_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("EDGE_ENDPOINT")
	apiKeyID := os.Getenv("EDGE_API_KEY_ID")
	apiKey := os.Getenv("EDGE_API_KEY")

	if !cfg.Endpoint.IsNull() {
		endpoint = cfg.Endpoint.ValueString()
	}

	if !cfg.APIKeyID.IsNull() {
		apiKeyID = cfg.APIKeyID.ValueString()
	}

	if !cfg.APIKey.IsNull() {
		apiKey = cfg.APIKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Edge API Endpoint",
			"The provider cannot create the Edge API client as there is a missing or empty value for the Edge API endpoint. "+
				"Set the endpoint value in the configuration or use the EDGE_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKeyID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key_id"),
			"Missing Edge API Key",
			"The provider cannot create the Edge API client as there is a missing or empty value for the Edge API Key. "+
				"Set the api_key value in the configuration or use the EDGE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Edge API Key",
			"The provider cannot create the Edge API client as there is a missing or empty value for the Edge API Key. "+
				"Set the api_key value in the configuration or use the EDGE_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "edge_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "edge_api_key_id", apiKeyID)
	ctx = tflog.SetField(ctx, "edge_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "edge_api_key")

	tflog.Debug(ctx, "Creating Edge client")

	if p.config == nil {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 5
		rc := retryClient.StandardClient()
		p.config = &config{
			m:        &sync.Mutex{},
			ua:       "terraform-provider-edge",
			keyID:    apiKeyID,
			key:      apiKey,
			endpoint: endpoint,
			client:   rc,
		}
	}

	resp.DataSourceData = p.config
	resp.ResourceData = p.config

	tflog.Info(ctx, "Configured Edge client", map[string]any{"success": true})
}

func (p *EdgeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *EdgeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewValueResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &EdgeProvider{
			version: version,
		}
	}
}

type config struct {
	m        *sync.Mutex
	ua       string
	keyID    string
	key      string
	endpoint string
	client   *http.Client
}

func (c *config) GetValue(ctx context.Context, id string) (*model.Value, error) {
	const path = "/service.Value/Get"
	u, err := url.JoinPath(c.endpoint, path)
	if err != nil {
		return nil, err
	}

	m := &model.GetValueRequest{ID: id}
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req, false)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= http.StatusBadRequest || resp.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	value := new(model.Value)
	err = json.NewDecoder(resp.Body).Decode(value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *config) CreateValue(ctx context.Context, value *model.Value) (*model.Value, error) {
	const path = "/service.Value/Create"
	u, err := url.JoinPath(c.endpoint, path)
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	v := new(model.Value)
	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		return nil, err
	}
	return v, err
}

func (c *config) UpdateValue(ctx context.Context, value *model.Value) (*model.Value, error) {
	const path = "/service.Value/Update"
	u, err := url.JoinPath(c.endpoint, path)
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	v := new(model.Value)
	err = json.NewDecoder(resp.Body).Decode(v)
	return v, err
}

func (c *config) DeleteValue(ctx context.Context, id string) error {
	const path = "/service.Value/Delete"
	u, err := url.JoinPath(c.endpoint, path)
	if err != nil {
		return err
	}

	m := &model.DeleteValueRequest{ID: id}
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(j))
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, true)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest || resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *config) do(ctx context.Context, req *http.Request, useMutex bool) (*http.Response, error) {
	if useMutex {
		c.m.Lock()
		defer c.m.Unlock()
	}
	req.Header.Set(headerKeyID, c.keyID)
	req.Header.Set(headerKey, c.key)
	req.Header.Set(headerUA, c.ua)
	req.Header.Set(headerContentType, applicationJSON)
	req.WithContext(ctx)
	res, err := c.client.Do(req)
	return res, err
}
