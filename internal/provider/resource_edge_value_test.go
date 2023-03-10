package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccResourceEdgeValue_BooleanValue(t *testing.T) {
	original := client.Transport
	defer func() { client.Transport = original }()
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Create",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Get",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Update",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Delete",
		httpmock.NewStringResponder(200, "{}"),
	)
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBoolean(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "value_id", "test-bool-value"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "enabled", "true"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "description", "test bool value"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "default_variant", "off"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.#", "2"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.0.variant", "off"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.0.value", "false"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.1.variant", "on"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.1.value", "true"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.#", "2"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.0.variant", "on"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.0.spec", "cel"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.0.expr", "env == 'dev'"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.1.variant", "on"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.1.spec", "cel"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.1.expr", "userId == 'XXX'"),
				),
			},
		},
	})
}

func TestAccResourceEdgeValue_StringValue(t *testing.T) {
	original := client.Transport
	defer func() { client.Transport = original }()
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Create",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Get",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Update",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Delete",
		httpmock.NewStringResponder(200, "{}"),
	)
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceString(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edge_value.test-string-value", "value_id", "test-string-value"),
					resource.TestCheckResourceAttr("edge_value.test-string-value", "enabled", "true"),
					resource.TestCheckResourceAttr("edge_value.test-string-value", "description", "test string value"),
					resource.TestCheckResourceAttr("edge_value.test-string-value", "default_variant", "key"),
					resource.TestCheckResourceAttr("edge_value.test-string-value", "string_value.#", "1"),
					resource.TestCheckResourceAttr("edge_value.test-string-value", "string_value.0.variant", "key"),
					resource.TestCheckResourceAttr("edge_value.test-string-value", "string_value.0.value", "test value"),
					resource.TestCheckResourceAttr("edge_value.test-string-value", "targeting.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceEdgeValue_JSONValue(t *testing.T) {
	original := client.Transport
	defer func() { client.Transport = original }()
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Create",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Get",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Update",
		httpmock.NewStringResponder(200, "{}"),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Delete",
		httpmock.NewStringResponder(200, "{}"),
	)
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceJSON(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edge_value.test-json-value", "value_id", "test-json-value"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "enabled", "true"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "description", "test json value"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "default_variant", "json"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.#", "1"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.variant", "json"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.value", "{\"key1\": \"value1\"}"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "targeting.#", "0"),
				),
			},
		},
	})
}

func testAccResourceBoolean() string {
	return `
resource "edge_value" "test-bool-value" {
  value_id = "test-bool-value"
  enabled = true
  description = "test bool value"
  default_variant = "off"

  boolean_value {
	variant = "on"
	value = true
  }

  boolean_value {
	variant = "off"
	value = false
  }

  targeting {
    variant = "on"
    spec = "cel"
    expr = "env == 'dev'"
  }

  targeting {
    variant = "on"
    spec = "cel"
    expr = "userId == 'XXX'"
  }
}`
}

func testAccResourceString() string {
	return `
resource "edge_value" "test-string-value" {
  value_id = "test-string-value"
  enabled = true
  description = "test string value"
  default_variant = "key"

  string_value {
	variant = "key"
	value = "test value"
  }
}`
}

func testAccResourceJSON() string {
	return `
resource "edge_value" "test-json-value" {
  value_id = "test-json-value"
  enabled = true
  description = "test json value"
  default_variant = "json"

  json_value {
	variant = "json"
	value = "{\"key1\": \"value1\"}"
  }
}`
}
