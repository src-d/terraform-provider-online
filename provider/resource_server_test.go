package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceServer(t *testing.T) {
	hcl := `
		resource "online_server" "test" {
 			name = "105770"
			hostname = "mvp"
		}
	`

	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{{
			ImportStateVerify: true,
			Config:            hcl,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("online_server.test", "private_mac"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_ip"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_mac"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_dns"),
			),
		}},
	})

}
