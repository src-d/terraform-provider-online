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
					name = "terraform-server-test"
					vlan = "2242"
				}

				resource "online_server" "test" {
	 				name = "105770"
					hostname = "mvp"

					private_interface {
						rpn = "${online_rpn.test.name}"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.vlan_id", "2242"),
			),
		}, {
			ImportStateVerify: true,
			Config: `
				resource "online_server" "test" {
	 				name = "105770"
					hostname = "mvp"

					private_interface {}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.vlan_id", "0"),
			),
		}},
	})
}
