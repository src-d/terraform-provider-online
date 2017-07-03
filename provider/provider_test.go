package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var providers = map[string]terraform.ResourceProvider{
	"online": Provider(),
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
