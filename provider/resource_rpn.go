package provider

import (
	"time"

	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/src-d/terraform-provider-online-net/online"
)

func resourceRPN() *schema.Resource {
	return &schema.Resource{
		Create: resourceRPNCreate,
		Update: resourceRPNCreate,
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

func resourceRPNCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)
	r, err := getRPN(c, d)
	if err != nil {
		return err
	}

	if err := updateRPNIfNeeded(c, r, d); err != nil {
		return err
	}

	applyRPN(r, d)
	return nil
}

func getRPN(c online.Client, d *schema.ResourceData) (*online.RPNv2, error) {
	name := d.Get("name").(string)
	d.SetId(name)

	return c.RPNv2ByName(name)
}

func updateRPNIfNeeded(c online.Client, prev *online.RPNv2, d *schema.ResourceData) error {
	var id int
	if prev != nil {
		id = prev.ID
	}

	rpn := &online.RPNv2{
		ID:   id,
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

	return c.SetRPNv2(rpn, time.Minute)
}

func resourceRPNRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(online.Client)
	r, err := getRPN(client, d)
	if err != nil {
		return err
	}

	if r == nil {
		return nil
	}

	applyRPN(r, d)
	return nil
}

const missingMemberStatus = "MISSING"

func applyRPN(r *online.RPNv2, d *schema.ResourceData) {
	fmt.Println(r)
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

func resourceRPNDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Id() == "" {
		return nil
	}

	client := meta.(online.Client)
	rpn, err := getRPN(client, d)
	if err != nil {
		return err
	}

	if rpn == nil {
		return nil
	}

	return client.DeleteRPNv2(rpn.ID, time.Minute)
}
