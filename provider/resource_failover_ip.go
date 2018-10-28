package provider

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online/online"
)

func resourceFailoverIP() *schema.Resource {
	return &schema.Resource{
		Read:   resourceFailoverIPRead,
		Create: resourceFailoverIPCreate,
		Delete: resourceFailoverIPDelete,
		Update: resourceFailoverIPUpdate,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the failover IP to modify",
			},
			"destination_server_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "The ID of the server to route the failover IP to",
				ConflictsWith: []string{"destination_server_ip"},
			},
			"destination_server_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The ID of the server to route the failover IP to",
				ConflictsWith: []string{"destination_server_id"},
			},
			"generate_mac": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "should a virtual mac be created for the IP",
			},
			"generate_mac_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "should a virtual mac be created for the IP",
				Default:     "kvm",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "vmware" && v != "xen" && v != "kvm" && v != "" {
						errs = append(errs, errors.New("generate_mac_type must be either vmware, xen or kvm"))
					}
					return
				},
			},
			"mac": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "the generated mac",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the reverse hostname",
			},
		},
	}
}

func resourceFailoverIPRead(d *schema.ResourceData, meta interface{}) error {
	ip := d.Get("ip").(string)
	d.SetId(ip)

	return nil
}

func resourceFailoverIPCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)

	ip := d.Get("ip").(string)
	serverIDInterface, hasServerID := d.GetOk("destination_server_id")
	serverIPInterface, hasServerIP := d.GetOk("destination_server_ip")
	generateMac := d.Get("generate_mac").(bool)
	macType := d.Get("generate_mac_type").(string)
	hostnameInterface, hasHostname := d.GetOk("hostname")

	dstIP := ""

	if hasServerIP {
		dstIP = serverIPInterface.(string)
	}

	if hasServerID {
		server, err := c.Server(serverIDInterface.(int))
		if err != nil {
			return err
		}
		dstIP = server.InterfaceByType(online.Public).Address
	}

	if hasHostname {
		hostname := hostnameInterface.(string)
		err := c.SetReverseFailoverIP(ip, hostname)
		if err != nil {
			return err
		}
	}

	err := c.EditFailoverIP(ip, dstIP)
	if err != nil && !strings.Contains(err.Error(), "Address already provisioned") {
		return err
	}

	if generateMac {
		mac, err := c.GenerateMACFailoverIP(ip, macType)
		if err != nil {
			return err
		}
		d.Set("mac", mac)
	}

	d.SetId(ip)
	return nil
}

func resourceFailoverIPDelete(d *schema.ResourceData, meta interface{}) error {
	ip := d.Get("ip").(string)
	_, macExists := d.GetOkExists("mac")
	c := meta.(online.Client)

	if macExists {
		err := c.DeleteMACFailoverIP(ip)
		if err != nil {
			return err
		}
		d.Set("mac", "")
	}

	err := c.EditFailoverIP(ip, "")
	if err != nil {
		return err
	}

	return err
}

func resourceFailoverIPUpdate(d *schema.ResourceData, meta interface{}) error {
	ip := d.Get("ip").(string)
	c := meta.(online.Client)

	hasNewMACRequest := d.HasChange("generate_mac")

	hasNewServerID := d.HasChange("destination_server_id")
	hasNewServerIP := d.HasChange("destination_server_ip")

	hasNewHostname := d.HasChange("hostname")

	if hasNewServerID && hasNewServerIP {
		// we switched here!
		_, hasdID := d.GetOkExists("destination_server_id")
		_, hasIP := d.GetOkExists("destination_server_ip")
		serverID := d.Get("destination_server_id").(int)

		// make sure the next steps will use the correct key
		if hasdID && serverID != 0 {
			hasNewServerIP = false
		} else if hasIP {
			hasNewServerID = false
		}
	}

	if hasNewServerID {
		serverID := d.Get("destination_server_id").(int)
		server, err := c.Server(serverID)
		if err != nil {
			return err
		}

		dstIP := server.InterfaceByType(online.Public).Address
		err = c.EditFailoverIP(ip, dstIP)
		if err != nil {
			return err
		}
	} else if hasNewServerIP {
		dstIP := d.Get("destination_server_ip").(string)
		err := c.EditFailoverIP(ip, dstIP)
		if err != nil {
			return err
		}
	}

	if hasNewHostname {
		hostname := d.Get("hostname").(string)
		err := c.SetReverseFailoverIP(ip, hostname)
		if err != nil {
			return err
		}
	}

	if hasNewMACRequest {
		// we need to enable or disable the generated MAC
		_, newValue := d.GetChange("generate_mac")
		if newValue.(bool) {
			macType := d.Get("generate_mac_type").(string)
			mac, err := c.GenerateMACFailoverIP(ip, macType)
			if err != nil {
				return err
			}

			d.Set("mac", mac)
		} else {
			err := c.DeleteMACFailoverIP(ip)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
