package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func init() {
	onlineClientMock.On("EditReverseHostname", "127.0.0.1", "google.com").Return(nil)
	onlineClientMock.On("EditReverseHostname", "127.0.0.1", "false").Return(nil)
}

func TestResourceReverseHostname(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testMockProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				resource "online_reverse_hostname" "test" {
	 				"ip" = "127.0.0.1"
					"hostname" = "google.com"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_reverse_hostname.test", "ip", "127.0.0.1"),
				),
			},
		},
	})
}

func TestResourceReverseHostnameAcceptance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testAccProviders,
		IsUnitTest: false,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				resource "online_reverse_hostname" "test" {
	 				"ip" = "51.158.20.2"
					"hostname" = "niceip.eyskens.me"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_reverse_hostname.test", "ip", "51.158.20.2"),
				),
			},
		},
	})
}
