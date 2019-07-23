package provider

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online/online"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Update: resourceServerUpdate,
		Read:   resourceServerRead,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"server_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "server id",
			},
			"hostname": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "server hostname",
			},
			"public_interface": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        resourceInterface(),
				Description: "Public interface properties",
			},
			"private_interface": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Elem:        resourceInterface(),
				Description: "Private interface properties",
			},
		},
	}
}

func resourceInterface() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mac": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hardware address of the device.",
			},
			"dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DNS server address.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Static IPv4 address.",
			},
		},
	}
}

func resourceServerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceServerRead(d, meta)
}

func resourceServerUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	s := &online.Server{
		ID:       d.Get("server_id").(int),
		Hostname: d.Get("hostname").(string),
	}

	publicDNS := d.Get("public_interface.dns").(string)
	ip := s.InterfaceByType(online.Public)
	if ip != nil && publicDNS != "" && ip.Reverse != publicDNS {
		ip.Reverse = publicDNS
	}

	if err := c.SetServer(s); err != nil {
		return err
	}

	return resourceServerRead(d, meta)
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	id := d.Get("server_id").(int)

	s, err := c.Server(id)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("hostname", s.Hostname)
	setIP(s, d)

	return nil
}

func setIP(s *online.Server, d *schema.ResourceData) {
	var public, private map[string]interface{}

	for _, iface := range s.IP {
		switch iface.Type {
		case online.Public:
			public = map[string]interface{}{
				"mac":     strings.ToLower(iface.MAC),
				"dns":     iface.Reverse,
				"address": iface.Address,
			}

		case online.Private:
			private = map[string]interface{}{
				"mac":     strings.ToLower(iface.MAC),
				"dns":     iface.Reverse,
				"address": iface.Address,
			}
		}
	}

	d.Set("public_interface", public)
	d.Set("private_interface", private)
}
