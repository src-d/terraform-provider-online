package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceRPN(t *testing.T) {
	hcl := `
		resource "online_rpn" "test" {
			name = "terraform"
			vlan = "2040"
			member {
				id = "105770"
			}
		}

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
				resource.TestCheckResourceAttrSet("online_rpn.test", "status"),
				resource.TestCheckResourceAttrSet("online_rpn.test", "member.0.status"),
				resource.TestCheckResourceAttr("online_server.test", "rpn.0", "2040"),
			),
		}},
	})
}
