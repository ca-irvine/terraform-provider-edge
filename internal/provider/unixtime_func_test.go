package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func Test_UnixTimeConverterFunc(t *testing.T) {
	t.Parallel()

	cfg := &config{}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: testUnixTimeConverterFuncConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckOutput("test", "1710325173"),
				),
			},
		},
	})
}

func testUnixTimeConverterFuncConfig() string {
	return `
output "test" {
  value = provider::edge::unixtime("2024-03-13T19:19:33+09:00")
}`
}
