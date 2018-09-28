package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestResourceServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: providers,
		Steps: []resource.TestStep{{
			ImportStateVerify: true,
			Config: `
				resource "online_server" "test" {
					server_id = "105770"
					hostname = "mvp"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.address"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.mac"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.dns"),
			),
		}},
	})
}

func TestResourceServerRPNAdd(t *testing.T) {
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
					server_id = "105770"
					hostname = "mvp"

					private_interface {
						rpn = "${online_rpn.test.name}"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.address"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.mac"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.dns"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "private_interface.0.mac"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.vlan_id", "2242"),
			),
		}, {
			ImportStateVerify: true,
			Config: `
				resource "online_rpn" "test" {
					name = "terraform-server-test"
					vlan = "2242"
				}

				resource "online_rpn" "test_alt" {
					name = "terraform-server-test_alt"
					vlan = "2243"
				}

				resource "online_server" "test" {
					server_id = "105770"
					hostname = "mvp"

					private_interface {
						rpn = "${online_rpn.test.name}"
					}

					private_interface {
						rpn = "${online_rpn.test_alt.name}"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.address"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.mac"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.dns"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "2"),
				resource.TestCheckResourceAttrSet("online_server.test", "private_interface.0.mac"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.vlan_id", "2242"),
				resource.TestCheckResourceAttrSet("online_server.test", "private_interface.1.mac"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.1.vlan_id", "2243"),
			),
		}},
	})
}

func TestResourceServerRPNDelete(t *testing.T) {
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
					server_id = "105770"
					hostname = "mvp"

					private_interface {
						rpn = "${online_rpn.test.name}"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.address"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.mac"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.dns"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "private_interface.0.mac"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.vlan_id", "2242"),
			),
		}, {
			ImportStateVerify: true,
			Config: `
				resource "online_rpn" "test" {
					name = "terraform-server-test"
					vlan = "2242"
				}

				resource "online_server" "test" {
					server_id = "105770"
					hostname = "mvp"

					private_interface {}
				}

				resource "online_server" "foo" {
					server_id = "105771"
					hostname = "mvp"

					private_interface {
						rpn = "${online_rpn.test.name}"
					}
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.address"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.mac"),
				resource.TestCheckResourceAttrSet("online_server.test", "public_interface.0.dns"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "1"),
				resource.TestCheckResourceAttrSet("online_server.test", "private_interface.0.mac"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.vlan_id", "0"),
			),
		}},
	})
}
