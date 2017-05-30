package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerNone,
		Delete: resourceServerNone,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeInt,
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
				Optional: true,
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

func resourceServerNone(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	id := d.Get("name").(int)
	d.SetId(string(id))

	s := &online.Server{
		ID:       id,
		Hostname: d.Get("hostname").(string),
	}

	client := meta.(online.Client)
	if err := client.SetServer(s); err != nil {
		return err
	}

	fmt.Println(s)

	return resourceServerRead(d, meta)
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(online.Client)

	id := d.Get("name").(int)
	fmt.Println(id)
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
