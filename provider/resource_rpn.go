package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceRPN() *schema.Resource {
	return &schema.Resource{
		Create: resourceRPNCreate,
		Update: resourceRPNUpdate,
		Read:   resourceRPNRead,
		Delete: resourceRPNDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  online.Standard,
			},
			"vlan": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceRPNCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	c := meta.(online.Client)
	rpn, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpn == nil {
		rpn = &online.RPNv2{Name: name}
	}

	rpn.Type = online.RPNv2Type(d.Get("type").(string))

	d.SetId(rpn.Name)
	globalCache.addRPN(rpn, d.Get("vlan").(int))

	return nil
}

func resourceRPNUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceRPNRead(d, meta); err != nil {
		return err
	}

	name := d.Get("name").(string)
	cache := globalCache.rpn[name]

	for _, m := range cache.Members {
		m.VLAN = cache.VLAN
	}

	cache.Type = online.RPNv2Type(d.Get("type").(string))

	c := meta.(online.Client)
	return c.SetRPNv2(&cache.RPNv2, time.Minute)
}

func resourceRPNRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	c := meta.(online.Client)
	rpn, err := c.RPNv2ByName(name)
	if err != nil {
		return err
	}

	if rpn == nil {
		return fmt.Errorf("missing RPNv2 group: %q", name)
	}

	d.SetId(rpn.Name)
	globalCache.addRPN(rpn, d.Get("vlan").(int))

	return nil
}

func resourceRPNDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	c := meta.(online.Client)
	rpn, err := c.RPNv2ByName(d.Id())
	if err != nil {
		return err
	}

	if rpn == nil {
		return nil
	}

	return c.DeleteRPNv2(rpn.ID, time.Minute)
}
