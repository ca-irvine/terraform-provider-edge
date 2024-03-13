package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = unixTimeConverterFunc{}

func NewUnixTimeConverterFunc() function.Function {
	return &unixTimeConverterFunc{}
}

type unixTimeConverterFunc struct{}

func (u unixTimeConverterFunc) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "unixtime"
}

func (u unixTimeConverterFunc) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "The utility function for converting format to unixtime.",
		MarkdownDescription: "The utility function for converting format to unixtime.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "format",
				MarkdownDescription: "format to convert to unixtime",
			},
		},
		Return: function.Int64Return{},
	}
}

func (u unixTimeConverterFunc) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var format string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &format))
	if resp.Error != nil {
		return
	}

	t, err := time.Parse(time.RFC3339, format)
	if err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, t.Unix()))
}
