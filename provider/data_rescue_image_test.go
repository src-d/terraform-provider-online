package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestDataRescueImage(t *testing.T) {
	onlineClientMock.On("GetRescueImages", 123).Return([]string{"ubuntu-18.04-amd64", "darwin-9-armhf"}, nil)
	resource.Test(t, resource.TestCase{
		Providers:  testMockProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				data "online_rescue_image" "test" {
	 				name = "ubuntu-18.04-amd64"
					server = 123
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "name", "ubuntu-18.04-amd64"),
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "image", "ubuntu-18.04-amd64"),
				),
			},
			{
				ImportStateVerify: false,
				Config: `
				data "online_rescue_image" "test" {
	 				name_filter = "ubuntu"
					server = 123
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "name_filter", "ubuntu"),
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "image", "ubuntu-18.04-amd64"),
				),
			},
			{
				ImportStateVerify: false,
				Config: `
					data "online_rescue_image" "test" {
						 name_filter = "fedora"
						server = 123
					}
				`,
				ExpectError: regexp.MustCompile(`No image found for requirements`),
			},
		},
	})
}

func TestDataRescueImageAcceptance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testAccProviders,
		IsUnitTest: false,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				data "online_rescue_image" "test" {
	 				name = "ubuntu-18.04_amd64"
					server = 105711
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "name", "ubuntu-18.04_amd64"),
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "image", "ubuntu-18.04_amd64"),
				),
			},
			{
				ImportStateVerify: false,
				Config: `
				data "online_rescue_image" "test" {
	 				name_filter = "ubuntu-18.04"
					server = 105711
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "name_filter", "ubuntu-18.04"),
					resource.TestCheckResourceAttr("data.online_rescue_image.test", "image", "ubuntu-18.04_amd64"),
				),
			},
		},
	})
}
