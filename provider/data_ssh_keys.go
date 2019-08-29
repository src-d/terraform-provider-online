package provider

import (
	"time"

	"github.com/src-d/terraform-provider-online/online"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSSHKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSSHKeysRead,

		Schema: map[string]*schema.Schema{
			"ssh_keys": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "User's SSH keys.",
			},
		},
	}
}

func dataSSHKeysRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)

	keys, err := c.ListSSHKeys()
	if err != nil {
		return err
	}

	userKeys := make([]interface{}, len(*keys))
	for i, v := range *keys {
		userKeys[i] = v.UUID
	}

	d.SetId(time.Now().UTC().String())
	if err = d.Set("ssh_keys", userKeys); err != nil {
		return err
	}

	return nil
}
