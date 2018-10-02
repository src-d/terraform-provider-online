package provider

import (
	"os"
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

var TestServerID string

func init() {
	TestServerID = os.Getenv("ONLINE_SERVER_ID")
}
