package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceReverseHostname() *schema.Resource {
	return &schema.Resource{
		Read:   resourceReverseHostnameRead,
		Create: resourceReverseHostnameCreate,
		Delete: resourceReverseHostnameDelete,
		Update: resourceReverseHostnameUpdate,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the IP to set the reverse hostname on",
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The reverse hostname",
			},
		},
	}
}

func resourceReverseHostnameRead(d *schema.ResourceData, meta interface{}) error {
	ip := d.Get("ip").(string)
	d.SetId(ip)

	return nil
}

func resourceReverseHostnameCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	ip := d.Get("ip").(string)
	hostname := d.Get("hostname").(string)

	d.SetId(ip)

	err := c.EditReverseHostname(ip, hostname)
	if err != nil {
		return err
	}

	return nil
}

func resourceReverseHostnameDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	ip := d.Get("ip").(string)

	return c.EditReverseHostname(ip, "false")
}

func resourceReverseHostnameUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)

	if d.HasChange("hostname") {
		ip := d.Get("ip").(string)
		hostname, _ := d.GetChange("hostname")

		err := c.EditReverseHostname(ip, hostname.(string))
		if err != nil {
			return err
		}
	}

	return nil
}
