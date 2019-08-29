package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/src-d/terraform-provider-online/online"
)

func TestDataOperatingSystem(t *testing.T) {
	onlineClientMock.On("ListOperatingSystems", 123).Return(&online.OperatingSystems{online.OperatingSystem{
		OS: online.OS{
			Name:    "centos",
			Version: "CentOS 7.6",
		},
		ID:      305,
		Type:    "server",
		Release: "2014-07-06T22:00:00.000Z",
		Arch:    "64 bits",
	}}, nil)
	resource.Test(t, resource.TestCase{
		Providers:  testMockProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				data "online_operating_system" "test" {
	 				name = "centos"
	 				version = "CentOS 7.6"
					server_id = 123
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.online_operating_system.test", "os_id", "305"),
				),
			},
			{
				ImportStateVerify: false,
				Config: `
					data "online_operating_system" "test" {
						name = "centos"
	 				    version = "CentOS 666"
						server_id = 123
					}
				`,
				ExpectError: regexp.MustCompile("unable to find OS"),
			},
		},
	})
}

func TestDataOperatingSystemAcceptance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testAccProviders,
		IsUnitTest: false,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: fmt.Sprintf(`
				data "online_operating_system" "test" {
	 				name = "centos"
	 				version = "CentOS 7.6"
					server_id = %s
				}
			`, TestServerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.online_operating_system.test", "os_id", "305"),
				),
			},
		},
	})
}
