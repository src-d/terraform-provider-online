package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/src-d/terraform-provider-online/online"
)

func setupMock() {
	onlineClientMock.On("SetServer", &online.Server{
		ID:       123,
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
		Hostname:      "mock",
		InstallStatus: "installed",
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
	onlineClientMock.On("InstallServer", 123, &online.ServerInstall{
		Hostname:                "mock",
		OS_ID:                   "101",
		UserLogin:               "user1",
		UserPassword:            "pass1",
		RootPassword:            "rootpass",
		PartitioningTemplateRef: "81c651de-030b-41f3-8094-36f423375234",
		SSHKeys: []string{
			"81c651de-030b-41f3-8094-36f423375235",
			"81c651de-030b-41f3-8094-36f423375236",
		},
	}).Return(nil)
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
					os_id = "101"
					user_login = "user1"
					user_password = "pass1"
					root_password = "rootpass"
					partitioning_template_ref = "81c651de-030b-41f3-8094-36f423375234"
                    ssh_keys = [
                        "81c651de-030b-41f3-8094-36f423375235",
                        "81c651de-030b-41f3-8094-36f423375236",
                    ]
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
				resource.TestCheckResourceAttr("online_server.test", "status", "installed"),
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
