package provider

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/src-d/terraform-provider-online-net/online/mock"
)

var testAccProviders = map[string]terraform.ResourceProvider{}
var testMockProviders = map[string]terraform.ResourceProvider{}

var onlineClientMock = new(mock.OnlineClientMock)

func init() {
	if os.Getenv("TF_ACC") == "1" {
		testAccProviders["online"] = Provider()
	}

	os.Setenv(TokenEnvVar, "test-token")

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
		os.Setenv(TokenEnvVar, "test-token")
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

var TestServerID string

func init() {
	TestServerID = os.Getenv("ONLINE_SERVER_ID")
}
