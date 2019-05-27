package provider

import (
	"fmt"
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
				resource.TestCheckResourceAttr("online_server.test", "public_interface.address", "1.2.3.4"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.mac", "aa:bb:cc:dd:ee:ff"),
				resource.TestCheckResourceAttr("online_server.test", "public_interface.dns", "my.dns.address"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.address", "10.2.3.4"),
				resource.TestCheckResourceAttr("online_server.test", "private_interface.mac", "00:bb:cc:dd:ee:ff"),
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
				Config: fmt.Sprintf(`
				resource "online_server" "test" {
					server_id = %s
					hostname  = "new-stg-worker-13"
				}
			`, TestServerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_server.test", "hostname", "new-stg-worker-13"),
					resource.TestCheckResourceAttr("online_server.test", "server_id", TestServerID),
				),
			},
			{
				ImportStateVerify: true,
				Config: fmt.Sprintf(`
				resource "online_server" "test" {
					server_id = %s
					hostname  = "stg-worker-13"
				}
			`, TestServerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_server.test", "hostname", "stg-worker-13"),
					resource.TestCheckResourceAttr("online_server.test", "server_id", TestServerID),
				),
			},
		},
	})
}
