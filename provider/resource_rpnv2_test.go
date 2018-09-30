package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceRPNv2Acceptance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: true,
				Config: `
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2999"
					server_ids = ["${online_server.test1.server_id}"]
				}

				resource "online_server" "test1" {
					server_id = 105711
					hostname  = "stg-worker-13"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2999"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "1"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", "105711"),
				),
			},
			{
				ImportStateVerify: true,
				Config: `
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2999"
					server_ids = [
						"${online_server.test1.server_id}",
						"${online_server.test2.server_id}",
					]
				}

				resource "online_server" "test1" {
					server_id = 105711
					hostname  = "stg-worker-13"
				}

				resource "online_server" "test2" {
					server_id = 105707
					hostname  = "stg-worker-11"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2999"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "2"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", "105711"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.1", "105707"),
				),
			},
			{
				ImportStateVerify: true,
				Config: `
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2998"
					server_ids = [
						"${online_server.test1.server_id}",
						"${online_server.test2.server_id}",
					]
				}

				resource "online_server" "test1" {
					server_id = 105711
					hostname  = "stg-worker-13"
				}

				resource "online_server" "test2" {
					server_id = 105707
					hostname  = "stg-worker-11"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2998"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "2"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", "105711"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.1", "105707"),
				),
			},
			{
				ImportStateVerify: true,
				Config: `
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2998"
					server_ids = [
						"${online_server.test2.server_id}",
					]
				}

				resource "online_server" "test2" {
					server_id = 105707
					hostname  = "stg-worker-11"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2998"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "1"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", "105707"),
				),
			},
		},
	})
}
