package provider

import (
	"strings"
	"time"

	"fmt"

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
			"public_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     resourcePublicInterface(),
			},
			"private_interface": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: false,
				Elem:     resourcePrivateInterface(),
			},
		},
	}
}

func resourcePrivateInterface() *schema.Resource {
	r := resourcePublicInterface()
	r.Schema["rpn"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "ID of the RPNv2 to be used by this interface.",
	}

	r.Schema["vlan_id"] = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Computed:    true,
		Description: "VLAN ID from the RPN assigned to this interface.",
	}

	return r
}

func resourcePublicInterface() *schema.Resource {
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

func resourceServerNone(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceServerCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	s, err := getServer(c, d)
	if err != nil {
		return err
	}

	if err := updateServerIfNeeded(c, s, d); err != nil {
		return err
	}

	if err := updateRPNsIfNeeded(c, s, d); err != nil {
		return err
	}

	return resourceServerRead(d, meta)
}

func updateServerIfNeeded(c online.Client, s *online.Server, d *schema.ResourceData) error {
	hostname := d.Get("hostname").(string)

	var changed bool
	if s.Hostname != hostname {
		changed = true
		s.Hostname = hostname
	}

	publicDNS := d.Get("public_interface.0.dns").(string)
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

func updateRPNsIfNeeded(c online.Client, s *online.Server, d *schema.ResourceData) error {
	rpns := make(map[string]bool)
	for _, p := range d.Get("private_interface").([]interface{}) {
		m := p.(map[string]interface{})
		rpn := m["rpn"].(string)
		if rpn == "" {
			continue
		}

		rpns[rpn] = true

		if err := addServerToRPN(c, s, rpn); err != nil {
			return err
		}
	}

	return removeServerFromRPNsIfNeeded(c, s, rpns)
}

func addServerToRPN(c online.Client, s *online.Server, rpn string) error {
	cache, ok := globalCache.rpn[rpn]
	if !ok {
		return fmt.Errorf("unknown rpn %q", rpn)
	}

	m := &online.Member{}
	m.Linked.ID = s.ID
	m.VLAN = cache.VLAN

	cache.Members = append(cache.Members, m)

	return c.SetRPNv2(&cache.RPNv2, time.Minute)
}

func removeServerFromRPNsIfNeeded(c online.Client, s *online.Server, rpns map[string]bool) error {
	list, err := c.ListRPNv2()
	if err != nil {
		return err
	}

	for _, rpn := range list {
		if _, ok := rpns[rpn.Name]; ok {
			continue
		}

		var keep []*online.Member
		for _, m := range rpn.Members {
			if m.Linked.ID != s.ID {
				keep = append(keep, m)
			}
		}

		if len(keep) == 0 {
			fmt.Printf("unused rpn %q please remove it from configuration\n", rpn.Name)
		} else if len(keep) < len(rpn.Members) {
			rpn.Members = keep
			if err := c.SetRPNv2(rpn, time.Minute); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(online.Client)
	s, err := getServer(client, d)
	if err != nil {
		return err
	}

	id := d.Get("name").(int)
	rpns, err := getRPNbyServer(client, id)
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

func getRPNbyServer(c online.Client, serverID int) (map[string]*online.RPNv2, error) {
	r, err := c.ListRPNv2()
	if err != nil {
		return nil, err
	}

	m := map[string]*online.RPNv2{}
	for _, rpn := range r {
		if rpn.MemberByServerID(serverID) != nil {
			m[rpn.Name] = rpn
		}
	}

	return m, nil
}

func applyServer(s *online.Server, rpns map[string]*online.RPNv2, d *schema.ResourceData) {
	var public, private []map[string]interface{}

	for _, iface := range s.IP {
		switch iface.Type {
		case online.Public:
			public = append(public, map[string]interface{}{
				"mac":     strings.ToLower(iface.MAC),
				"dns":     iface.Reverse,
				"address": iface.Address,
			})

		case online.Private:
			for _, p := range d.Get("private_interface").([]interface{}) {
				m := p.(map[string]interface{})
				m["mac"] = strings.ToLower(iface.MAC)
				m["vlan_id"] = 0

				rpn, ok := rpns[m["rpn"].(string)]
				if ok {
					m["vlan_id"] = rpn.MemberByServerID(s.ID).VLAN
				}

				private = append(private, m)
			}

		}
	}

	d.Set("public_interface", public)
	d.Set("private_interface", private)
}
