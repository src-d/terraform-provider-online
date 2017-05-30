package provider

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerRead,
		Delete: resourceServerRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": { // no-effect
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_dns": { // no-effect
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			/*"rpn": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			*/
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(online.Client)
	d.SetId("foo")
	return nil
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(online.Client)

	id := d.Get("name").(string)
	s, err := client.Server(id)
	if err != nil {
		return err
	}

	public := s.InterfaceByType(online.Public)
	if public != nil {
		d.Set("public_ip", public.Address)
		d.Set("public_mac", strings.ToLower(public.MAC))
		d.Set("public_dns", public.Reverse)
	}

	private := s.InterfaceByType(online.Private)
	if private != nil {
		d.Set("private_mac", strings.ToLower(private.MAC))
	}

	d.Set("rpn", []string{"1", "2"})

	return nil
}
