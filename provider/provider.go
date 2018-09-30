package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/src-d/terraform-provider-online-net/online"
)

const TokenEnvVar = "ONLINE_TOKEN"

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
			"online_rpnv2":  resourceRPNv2(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"online_rescue_image": dataRescueImage(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token := d.Get("token").(string)
	return online.NewClient(token), nil
}
