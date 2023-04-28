package main

import (
	"context"
	"flag"

	"github.com/ca-irvine/terraform-provider-edge/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	version = "dev"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Debug:   debugMode,
		Address: "registry.terraform.io/ca-irvine/edge",
	}

	providerserver.Serve(context.Background(), provider.New(version), opts)
}
