package provider

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform/terraform"
	"github.com/src-d/terraform-provider-online/online"

	"github.com/hashicorp/terraform/helper/resource"
)

func init() {
	onlineClientMock.On("EditFailoverIP", "127.0.0.1", "8.8.8.8").Return(nil)
	onlineClientMock.On("EditFailoverIP", "127.0.0.1", "").Return(nil)
	onlineClientMock.On("GenerateMACFailoverIP", "127.0.0.1", "kvm").Return("ma:ac:te:st", nil)
	onlineClientMock.On("DeleteMACFailoverIP", "127.0.0.1").Return(nil)
	onlineClientMock.On("Server", 1234).Return(&online.Server{
		IP: []*online.Interface{
			&online.Interface{
				Address: "8.8.8.8",
				Type:    online.Public,
			},
		},
	}, nil)
}

func TestResourceFailoverIP(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:  testMockProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				ImportStateVerify: false,
				Config: `
				resource "online_failover_ip" "test" {
	 				ip = "127.0.0.1"
					destination_server_ip = "8.8.8.8"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_failover_ip.test", "ip", "127.0.0.1"),
				),
			},
			{
				ImportStateVerify: false,
				Config: `
				resource "online_failover_ip" "test" {
	 				ip = "127.0.0.1"
					destination_server_id = "1234"
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_failover_ip.test", "ip", "127.0.0.1"),
				),
			},
			{
				ImportStateVerify: false,
				Config: `
				resource "online_failover_ip" "test" {
					 ip = "127.0.0.1"
					 destination_server_ip = "8.8.8.8"
					 generate_mac = true
				}
			`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_failover_ip.test", "ip", "127.0.0.1"),
					resource.TestCheckResourceAttr("online_failover_ip.test", "mac", "ma:ac:te:st"),
				),
			},
		},
	})
}

func TestResourceFailoverIPAcceptance(t *testing.T) {
	if TestFailoverIP == "" && os.Getenv("TF_ACC") == "1" {
		fmt.Println("Need ONLINE_FAILVOVER_IP to be set")
		t.Fatal("Need ONLINE_FAILVOVER_IP to be set")
		return
	}
	resource.Test(t, resource.TestCase{
		Providers:  testAccProviders,
		IsUnitTest: false,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// we are modifying routing tables here
					// online.net will error if we change these too quickly
					time.Sleep(10 * time.Second)
				},
				ImportStateVerify: false,
				Config: fmt.Sprintf(`
						resource "online_failover_ip" "test" {
			 				ip = "%s"
							destination_server_id = %s
						}
					`, TestFailoverIP, TestServerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("online_failover_ip.test", "ip", TestFailoverIP),
					func(s *terraform.State) error {
						// we are modifying routing tables here
						// online.net will error if we change these too quickly
						time.Sleep(10 * time.Second)
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					// we are modifying routing tables here
					// online.net will error if we change these too quickly
					time.Sleep(10 * time.Second)
				},
				ImportStateVerify: false,
				Config: fmt.Sprintf(`
				resource "online_failover_ip" "test" {
					 ip = "%s"
					 destination_server_id = %s
					 generate_mac = true
				}
			`, TestFailoverIP, TestServerID),
				Check: resource.ComposeAggregateTestCheckFunc(
					func(s *terraform.State) error {
						// we are modifying routing tables here
						// online.net will error if we change these too quickly
						time.Sleep(10 * time.Second)
						return nil
					},
					resource.TestCheckResourceAttr("online_failover_ip.test", "ip", TestFailoverIP),
					resource.TestCheckResourceAttrSet("online_failover_ip.test", "mac"),
				),
			},
		},
	})
}
