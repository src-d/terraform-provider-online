package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceRPN(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{{
			ImportStateVerify: true,
			Config: `
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
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("online_rpn.test", "status"),
				resource.TestCheckResourceAttrSet("online_rpn.test", "member.0.status"),
				resource.TestCheckResourceAttr("online_server.test", "rpn.0", "2040"),
			),
		}, {
			ImportStateVerify: true,
			Config: `
				resource "online_rpn" "test" {
					name = "terraform"
					vlan = "2040"
					member {
						id = "105770"
					}
					member {
						id = "105771"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("online_rpn.test", "status"),
				resource.TestCheckResourceAttr("online_rpn.test", "member.#", "2"),
			),
		}},
	})
}
