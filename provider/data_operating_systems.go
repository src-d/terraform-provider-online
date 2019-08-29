package provider

import (
	"encoding/json"
	"fmt"

	"github.com/src-d/terraform-provider-online/online"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataOperatingSystem() *schema.Resource {
	return &schema.Resource{
		Read: dataOperatingSystemRead,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "server id",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "exact name of the desired system",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "version name of the desired system",
			},
			"os_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "identifier of the desired system",
			},
		},
	}
}

func dataOperatingSystemRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)

	serverID := d.Get("server_id").(int)
	name := d.Get("name").(string)
	version := d.Get("version").(string)

	systems, err := c.ListOperatingSystems(serverID)
	if err != nil {
		return err
	}

	for _, os := range *systems {
		if os.Version == version && os.Name == name {
			d.Set("os_id", os.ID)
			d.SetId(fmt.Sprintf("%d", os.ID))
			return nil
		}
	}

	b, err := json.Marshal(systems)
	if err != nil {
		return err
	}
	return fmt.Errorf("unable to find OS: %s %s, available ones are %s", name, version, b)
}
