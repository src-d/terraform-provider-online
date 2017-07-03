package provider

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Update: resourceServerCreate,
		Read:   resourceServerRead,
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
			"rpn": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
		},
	}
}

func resourceServerNone(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	s, err := getServer(c, d)
	if err != nil {
		return err
	}

	rpns, err := getRPNbyServer(c, d)
	if err != nil {
		return err
	}

	defer applyServer(s, rpns, d)
	return updateServerIfNeeded(c, s, d)
}

func updateServerIfNeeded(c online.Client, s *online.Server, d *schema.ResourceData) error {
	hostname := d.Get("hostname").(string)
	publicDNS := d.Get("public_dns").(string)

	var changed bool
	if s.Hostname != hostname {
		changed = true
		s.Hostname = hostname
	}

	ip := s.InterfaceByType(online.Public)
	if ip != nil && publicDNS != "" && ip.Reverse != publicDNS {
		changed = true
		ip.Reverse = publicDNS
	}

	if !changed {
		return nil
	}

	return c.SetServer(s)
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(online.Client)
	s, err := getServer(client, d)
	if err != nil {
		return err
	}

	rpns, err := getRPNbyServer(client, d)
	if err != nil {
		return err
	}

	applyServer(s, rpns, d)
	return nil
}

func getServer(c online.Client, d *schema.ResourceData) (*online.Server, error) {
	id := d.Get("name").(int)
	d.SetId(string(id))

	return c.Server(id)
}

func getRPNbyServer(c online.Client, d *schema.ResourceData) ([]int, error) {
	id := d.Get("name").(int)

	r, err := c.ListRPNv2()
	if err != nil {
		return nil, err
	}

	var list []int
	for _, rpn := range r {
		m := rpn.MemberByServerID(id)
		if m != nil {
			list = append(list, m.Vlan)
		}
	}

	return list, nil
}

func applyServer(s *online.Server, rpn []int, d *schema.ResourceData) {
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

	d.Set("rpn", rpn)
}
