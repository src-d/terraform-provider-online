package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceRPNv2Acceptance(t *testing.T) {
	if TestServerID2 == "" && os.Getenv("TF_ACC") == "1" {
		t.Fatal("Need ONLINE_SERVER_ID_2 to be set")
		return
	}
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: true,
				Config: fmt.Sprintf(`
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2999"
					server_ids = ["${online_server.test1.server_id}"]
				}

				resource "online_server" "test1" {
					server_id = %s
					hostname  = "stg-worker-13"
				}
			`, TestServerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2999"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "1"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", TestServerID),
				),
			},
			{
				ImportStateVerify: true,
				Config: fmt.Sprintf(`
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2999"
					server_ids = [
						"${online_server.test1.server_id}",
						"${online_server.test2.server_id}",
					]
				}

				resource "online_server" "test1" {
					server_id = %s
					hostname  = "stg-worker-13"
				}

				resource "online_server" "test2" {
					server_id = %s
					hostname  = "stg-worker-11"
				}
			`, TestServerID, TestServerID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2999"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "2"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", TestServerID),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.1", TestServerID2),
				),
			},
			{
				ImportStateVerify: true,
				Config: fmt.Sprintf(`
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2998"
					server_ids = [
						"${online_server.test1.server_id}",
						"${online_server.test2.server_id}",
					]
				}

				resource "online_server" "test1" {
					server_id = %s
					hostname  = "stg-worker-13"
				}

				resource "online_server" "test2" {
					server_id = %s
					hostname  = "stg-worker-11"
				}
			`, TestServerID, TestServerID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2998"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "2"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", TestServerID),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.1", TestServerID2),
				),
			},
			{
				ImportStateVerify: true,
				Config: fmt.Sprintf(`
				resource "online_rpnv2" "test" {
					name = "terraform-provider-online-acceptance"
					vlan = "2998"
					server_ids = [
						"${online_server.test2.server_id}",
					]
				}

				resource "online_server" "test2" {
					server_id = %s
					hostname  = "stg-worker-11"
				}
			`, TestServerID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_rpnv2.test", "name", "terraform-provider-online-acceptance"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "vlan", "2998"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.#", "1"),
					resource.TestCheckResourceAttr("online_rpnv2.test", "server_ids.0", TestServerID2),
				),
			},
		},
	})
}
