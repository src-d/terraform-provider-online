package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/src-d/terraform-provider-online/online"
)

func setupMock() {
	onlineClientMock.On("SetServer", &online.Server{
		Hostname: "mock",
		IP: []*online.Interface{
			&online.Interface{
				Address: "1.2.3.4",
				MAC:     "aa:bb:cc:dd:ee:ff",
				Reverse: "my.dns.address",
				Type:    online.Public,
			},
			&online.Interface{
				Address: "10.2.3.4",
				MAC:     "00:bb:cc:dd:ee:ff",
				Type:    online.Private,
			},
		},
	}).Return(nil)
	onlineClientMock.On("Server", 123).Return(&online.Server{
		IP: []*online.Interface{
			&online.Interface{
				Address: "1.2.3.4",
				MAC:     "aa:bb:cc:dd:ee:ff",
				Reverse: "my.dns.address",
				Type:    online.Public,
			},
			&online.Interface{
				Address: "10.2.3.4",
				MAC:     "00:bb:cc:dd:ee:ff",
				Type:    online.Private,
			},
		},
	}, nil)
}

func TestResourceServerUnit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testMockProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{{
			ImportStateVerify: true,
			PreConfig:         setupMock,
			Config: `
				resource "online_server" "test" {
					server_id = 123
					hostname = "mock"
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("online_server.test", "hostname", "mock"),
				resource.TestCheckResourceAttr("online_server.test", "server_id", "123"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.0.address", "1.2.3.4"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.0.mac", "aa:bb:cc:dd:ee:ff"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.0.dns", "my.dns.address"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "1"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.address", "10.2.3.4"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.0.mac", "00:bb:cc:dd:ee:ff"),
			),
		}},
	})
}

func TestResourceServerAcceptance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: true,
				Config: `
				resource "online_server" "test" {
					server_id = 105711
					hostname  = "new-stg-worker-13"
					public_interface {
						dns = "terraform-provider-online-test-01.srcd.run."
					}
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_server.test", "hostname", "new-stg-worker-13"),
					resource.TestCheckResourceAttr("online_server.test", "server_id", "105711"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.0.address", "163.172.75.177"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.0.mac", "14:18:77:51:90:00"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.0.dns", "terraform-provider-online-test-01.srcd.run."),
					resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "1"),
					resource.TestCheckResourceAttr("online_server.test", "private_interface.0.mac", "a0:36:9f:b3:e9:ec"),
				),
			},
			{
				ImportStateVerify: true,
				Config: `
				resource "online_server" "test" {
					server_id = 105711
					hostname  = "stg-worker-13"
					public_interface {
						dns = "worker-13.infra.pipeline.staging.srcd.host."
					}
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_server.test", "hostname", "stg-worker-13"),
					resource.TestCheckResourceAttr("online_server.test", "server_id", "105711"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.#", "1"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.0.address", "163.172.75.177"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.0.mac", "14:18:77:51:90:00"),
					resource.TestCheckResourceAttr("online_server.test", "public_interface.0.dns", "worker-13.infra.pipeline.staging.srcd.host."),
					resource.TestCheckResourceAttr("online_server.test", "private_interface.#", "1"),
					resource.TestCheckResourceAttr("online_server.test", "private_interface.0.mac", "a0:36:9f:b3:e9:ec"),
				),
			},
		},
	})
}
