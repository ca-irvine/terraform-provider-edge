package provider

import (
	_ "embed"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

//go:embed testdata/boolean.json
var booleanTestdata string

func TestAccResourceEdgeValue_BooleanValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Create",
		httpmock.NewStringResponder(200, booleanTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Get",
		httpmock.NewStringResponder(200, booleanTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Update",
		httpmock.NewStringResponder(200, booleanTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Delete",
		httpmock.NewStringResponder(200, booleanTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		m:        &sync.Mutex{},
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceBoolean(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "value_id", "test-bool-value"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "enabled", "true"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "description", "test bool value"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "default_variant", "off"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.#", "2"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.0.variant", "on"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.0.value", "true"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.1.variant", "off"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "boolean_value.1.value", "false"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.#", "2"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.0.variant", "on"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.0.spec", "cel"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.0.expr", "env == 'dev'"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.1.variant", "on"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.1.spec", "cel"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "targeting.1.expr", "userId == 'XXX'"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "test.#", "1"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "test.0.variables", "{\"count\":1,\"env\":\"test\"}"),
					resource.TestCheckResourceAttr("edge_value.test-bool-value", "test.0.expected", "on"),
				),
			},
		},
	})
}

//go:embed testdata/string.json
var stringTestdata string

func TestAccResourceEdgeValue_StringValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Create",
		httpmock.NewStringResponder(200, stringTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Get",
		httpmock.NewStringResponder(200, stringTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Update",
		httpmock.NewStringResponder(200, stringTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Delete",
		httpmock.NewStringResponder(200, stringTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		m:        &sync.Mutex{},
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceString(),
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

//go:embed testdata/json.json
var jsonTestdata string

func TestAccResourceEdgeValue_JSONValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Create",
		httpmock.NewStringResponder(200, jsonTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Get",
		httpmock.NewStringResponder(200, jsonTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Update",
		httpmock.NewStringResponder(200, jsonTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Delete",
		httpmock.NewStringResponder(200, jsonTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		m:        &sync.Mutex{},
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceJSON(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edge_value.test-json-value", "value_id", "test-json-value"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "enabled", "true"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "description", "test json value"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "default_variant", "json"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.#", "1"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.variant", "json"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.value", "{\"items\":[{\"content\":\"content1\",\"viewable\":true},{\"content\":\"content2\",\"viewable\":true},{\"content\":\"content3\",\"viewable\":false}]}"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.transform.#", "2"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.transform.0.spec", "cel"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.transform.0.expr", "{\"items\":items.map(item, item.viewable ? item : item.deleteKey([\"content\"]))}"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.transform.1.spec", "cel"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "json_value.0.transform.1.expr", "{\"items\":items.map(item, item.viewable ? item.selectKey([\"content\"]) : item)}"),
					resource.TestCheckResourceAttr("edge_value.test-json-value", "targeting.#", "0"),
				),
			},
		},
	})
}

//go:embed testdata/integer.json
var integerTestdata string

func TestAccResourceEdgeValue_IntegerValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Create",
		httpmock.NewStringResponder(200, integerTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Get",
		httpmock.NewStringResponder(200, integerTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Update",
		httpmock.NewStringResponder(200, integerTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/service.Value/Delete",
		httpmock.NewStringResponder(200, integerTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		m:        &sync.Mutex{},
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceInteger(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "value_id", "test-integer-value"),
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "enabled", "true"),
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "description", "test integer value"),
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "default_variant", "one"),
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "integer_value.#", "1"),
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "integer_value.0.variant", "one"),
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "integer_value.0.value", "1"),
					resource.TestCheckResourceAttr("edge_value.test-integer-value", "targeting.#", "0"),
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
	
  test {
	variables = jsonencode({
	  env = "test"
	  count = 1
	})
	expected = "on"
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
	value = jsonencode({
	  "items": [
		{"viewable": true, "content": "content1"},
		{"viewable": true, "content": "content2"},
		{"viewable": false, "content": "content3"}
	  ]	
	})
	transform {
	  spec = "cel"
	  expr = "{\"items\":items.map(item, item.viewable ? item : item.deleteKey([\"content\"]))}"
	}
	transform {
	  spec = "cel"
	  expr = "{\"items\":items.map(item, item.viewable ? item.selectKey([\"content\"]) : item)}"
	}
  }
}`
}

func testAccResourceInteger() string {
	return `
resource "edge_value" "test-integer-value" {
  value_id = "test-integer-value"
  enabled = true
  description = "test integer value"
  default_variant = "one"

  integer_value {
	variant = "one"
	value = 1
  }
}`
}
