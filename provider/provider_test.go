package provider

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/src-d/terraform-provider-online/online/mock"
)

var testAccProviders = map[string]terraform.ResourceProvider{}
var testMockProviders = map[string]terraform.ResourceProvider{}

var onlineClientMock = new(mock.OnlineClientMock)

var TestServerID string
var TestServerID2 string
var TestToken = "test-token"
var TestFailoverIP string
var TestSSHUUID1 string

func init() {
	if os.Getenv("TF_ACC") == "1" {
		testAccProviders["online"] = Provider()
		TestServerID = os.Getenv("ONLINE_SERVER_ID")
		TestServerID2 = os.Getenv("ONLINE_SERVER_ID_2")
		TestFailoverIP = os.Getenv("ONLINE_FAILOVER_IP")
		TestSSHUUID1 = os.Getenv("ONLINE_SSH_UUID_1")
		TestToken = os.Getenv(TokenEnvVar)

		if TestToken == "" {
			fmt.Println("Need ONLINE_TOKEN to be set")
			os.Exit(1)
		}
		if TestServerID == "" {
			fmt.Println("Need ONLINE_SERVER_ID to be set")
			os.Exit(1)
		}
	}

	os.Setenv(TokenEnvVar, TestToken)

	// creating the provider with a mocked online.net api client
	provider := Provider().(*schema.Provider)
	provider.ConfigureFunc = providerConfigureMock
	testMockProviders["online"] = provider
}

func providerConfigureMock(d *schema.ResourceData) (interface{}, error) {
	return onlineClientMock, nil
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderMissingToken(t *testing.T) {
	os.Setenv(TokenEnvVar, "")

	defer func() {
		os.Setenv(TokenEnvVar, TestToken)
	}()

	_, fails := Provider().Validate(&terraform.ResourceConfig{})
	expectedErr := `"token": required field is not set`
	var err error

	for _, e := range fails {
		if strings.Contains(e.Error(), expectedErr) {
			err = e
			break
		}
	}

	if err == nil {
		t.Fatalf("no error received, but expected: %s", expectedErr)
	}
}
