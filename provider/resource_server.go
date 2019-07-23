package provider

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
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
			"os_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "OS identifier",
			},
			"user_login": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "User login",
			},
			"user_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "User password",
			},
			"root_password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
				Description: "Root password",
			},
			"partitioning_template_ref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "UUID of the partitioning template created from " +
					"https://console.online.net/en/template/partition",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Install status",
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
	c := meta.(online.Client)
	s := &online.ServerInstall{
		Hostname:                d.Get("hostname").(string),
		OS_ID:                   d.Get("os_id").(string),
		UserLogin:               d.Get("user_login").(string),
		UserPassword:            d.Get("user_password").(string),
		RootPassword:            d.Get("root_password").(string),
		PartitioningTemplateRef: d.Get("partitioning_template_ref").(string),
	}
	id := d.Get("server_id").(int)

	if err := c.InstallServer(id, s); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"installing"},
		Target:     []string{"installed"},
		Refresh:    waitForServerInstall(c, id),
		Timeout:    60 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceServerRead(d, meta)
}

func waitForServerInstall(c online.Client, id int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := c.Server(id)
		if err != nil {
			return s, "", err
		}
		return s, s.InstallStatus, nil
	}
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
	d.Set("status", s.InstallStatus)

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
