package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/src-d/terraform-provider-online/online"
)

func TestDataSSHKeys(t *testing.T) {
	onlineClientMock.On("ListSSHKeys").Return(&online.SSHKeys{online.SSHKey{
		UUID:        "7c8157c1-367b-43c5-bb78-0c598242cfb1",
		Description: "ssh unit test",
		Fingerprint: "c1:91:5e:42:55:5c:74:65:b6:12:32:7e:1f:6d:80:3e",
	}}, nil)
	resource.Test(t, resource.TestCase{
		Providers:  testMockProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{{
			ImportStateVerify: false,
			Config: `
				data "online_ssh_keys" "test" {
				}
			`,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(
					"data.online_ssh_keys.test",
					"ssh_keys.0",
					"7c8157c1-367b-43c5-bb78-0c598242cfb1",
				),
			),
		}},
	})
}

func TestDataSSHKeysAcceptance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testAccProviders,
		IsUnitTest: false,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				data "online_ssh_keys" "test" {
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.online_ssh_keys.test",
						"ssh_keys.0",
						TestSSHUUID1,
					),
				),
			},
		},
	})
}
