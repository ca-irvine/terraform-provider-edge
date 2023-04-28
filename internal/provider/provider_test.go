package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	providerConfig = `
provider "edge" {
  endpoint = "http://localhost:8018"
  api_key = "test_key"
  api_key_id = "test_key_id"
}
`
)

func protoV6ProviderFactories(cfg *config) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"edge": providerserver.NewProtocol6WithError(&EdgeProvider{
			version: "test",
			config:  cfg,
		}),
	}
}
