package provider

import (
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/src-d/terraform-provider-online-net/online"
)

const TokenEnvVar = "ONLINE_TOKEN"

// globalCache keeps the the internal implate types generated  by the different
// data resources with the goal to be reused by the other resources. The key of
// the maps are the name of resource.
var globalCache = newCache()

// Provider returns the provider schema to Terraform.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc(TokenEnvVar, ""),
				Required:    true,
				Sensitive:   true,
				Description: "Online.net private API token, by default the ONLINE_TOKEN environment variable is used.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"online_server": resourceServer(),
			"online_rpn":    resourceRPN(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("token").(string)
	return online.NewClient(token), nil
}

type cache struct {
	rpn map[string]*rpnCache

	sync.Mutex
}

func newCache() *cache {
	return &cache{
		rpn: make(map[string]*rpnCache, 0),
	}
}

type rpnCache struct {
	online.RPNv2
	VLAN int
}

func (c *cache) addRPN(r *online.RPNv2, vlan int) {
	c.Lock()
	defer c.Unlock()

	c.rpn[r.Name] = &rpnCache{*r, vlan}
}
