package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/ca-irvine/terraform-provider-edge/internal/model"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	client = &http.Client{}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_key_id": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("EDGE_API_KEY_ID", nil),
				},
				"api_key": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("EDGE_API_KEY", nil),
				},
				"endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("EDGE_API_ENDPOINT", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"cairvine_edge_value": resourceEdgeValue(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, data *schema.ResourceData) (any, diag.Diagnostics) {
		api := &config{
			m:        &sync.Mutex{},
			ua:       p.UserAgent("terraform-provider-edge", version),
			keyID:    data.Get("api_key_id").(string),
			key:      data.Get("api_key").(string),
			endpoint: data.Get("endpoint").(string),
		}
		log.Println("[INFO] Initializing edge client")
		return api, nil
	}
}

type config struct {
	m        *sync.Mutex
	ua       string
	keyID    string
	key      string
	endpoint string
}

var client *http.Client

func (v *config) GetValue(ctx context.Context, id string) (*model.Value, error) {
	const path = "/service.Value/Get"

	u, err := url.JoinPath(v.endpoint, path)
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

	b, err := v.do(ctx, req)

	tflog.Debug(ctx, "response", map[string]interface{}{"body": string(b)})
	value := new(model.Value)
	if err = json.Unmarshal(b, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (v *config) CreateValue(ctx context.Context, value *model.Value) error {
	const path = "/service.Value/Create"

	u, err := url.JoinPath(v.endpoint, path)
	if err != nil {
		return err
	}

	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	tflog.Debug(ctx, "CreateValue", map[string]interface{}{"json": string(j)})

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(j))
	if err != nil {
		return err
	}

	_, err = v.do(ctx, req)

	return err
}

func (v *config) UpdateValue(ctx context.Context, value *model.Value) error {
	const path = "/service.Value/Update"
	getVal, err := v.GetValue(ctx, value.ID)
	if err != nil {
		return err
	}
	value.CreateTime = getVal.CreateTime
	value.UpdateTime = getVal.UpdateTime
	u, err := url.JoinPath(v.endpoint, path)
	if err != nil {
		return err
	}

	j, err := json.Marshal(value)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(j))
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "UpdateValue", map[string]interface{}{"request": string(j)})

	_, err = v.do(ctx, req)

	return err
}

func (v *config) DeleteValue(ctx context.Context, id string) error {
	const path = "/service.Value/Delete"
	u, err := url.JoinPath(v.endpoint, path)
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

	_, err = v.do(ctx, req)

	return err
}

func (v *config) do(ctx context.Context, req *http.Request) ([]byte, error) {
	req.Header.Set(headerKeyID, v.keyID)
	req.Header.Set(headerKey, v.key)
	req.Header.Set(headerUA, v.ua)
	req.Header.Set(headerContentType, applicationJSON)
	req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body := res.Body
	statusCd := res.StatusCode
	defer func() { _ = body.Close() }()
	b, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if statusCd >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected status code %d, %s", statusCd, string(b))
	}
	return b, nil
}
