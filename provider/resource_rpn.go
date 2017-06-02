package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceRPN() *schema.Resource {
	return &schema.Resource{
		Create: resourceRPNCreate,
		Update: resourceRPNCreate,
		Read:   resourceRPNRead,
		Delete: resourceServerNone,

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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						}},
				},
			},
		},
	}
}

func resourceRPNNone(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceRPNCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	r, err := getRPN(c, d)
	if err != nil {
		return err
	}

	defer applyRPN(r, d)
	return updateRPNIfNeeded(c, r, d)
}

func getRPN(c online.Client, d *schema.ResourceData) (*online.RPNv2, error) {
	name := d.Get("name").(string)
	d.SetId(name)

	return c.RPNv2ByName(name)
}

func updateRPNIfNeeded(c online.Client, prev *online.RPNv2, d *schema.ResourceData) error {
	rpn := &online.RPNv2{
		ID:   prev.ID,
		Name: d.Get("name").(string),
		Type: online.RPNv2Type(d.Get("type").(string)),
	}

	for _, raw := range d.Get("member").([]interface{}) {
		value := raw.(map[string]interface{})

		m := &online.Member{}
		m.Linked.ID = value["id"].(int)
		m.Vlan = d.Get("vlan").(int)

		rpn.Members = append(rpn.Members, m)
	}

	return c.SetRPNv2(rpn)
}

func resourceRPNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(online.Client)
	r, err := getRPN(client, d)
	if err != nil {
		return err
	}

	applyRPN(r, d)
	return nil
}

const missingMemberStatus = "MISSING"

func applyRPN(r *online.RPNv2, d *schema.ResourceData) {
	if r == nil {
		return
	}

	d.Set("status", r.Status)

	var output []map[string]interface{}
	for _, m := range d.Get("member").([]interface{}) {
		value := m.(map[string]interface{})

		value["status"] = missingMemberStatus
		m := r.MemberByServerID(value["id"].(int))
		if m != nil {
			value["status"] = m.Status
		}

		output = append(output, value)
	}

	d.Set("member", output)
}
