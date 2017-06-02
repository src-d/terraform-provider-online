package provider

import (
	"fmt"

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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  online.Standard,
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
		Name: d.Get("name").(string),
		Type: online.RPNv2Type(d.Get("type").(string)),
	}

	for _, m := range d.Get("member").([]interface{}) {
		value := m.(map[string]interface{})
		rpn.Members = append(rpn.Members, &online.Member{
			ID: value["id"].(int),
		})
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

func applyRPN(r *online.RPNv2, d *schema.ResourceData) {
	if r == nil {
		return
	}

	d.Set("status", r.Status)
	for i, m := range d.Get("member").([]interface{}) {
		value := m.(map[string]interface{})

		m := r.MemberByID(value["id"].(int))
		fmt.Println(m, r)

		d.Set(fmt.Sprintf("member.%d.status", i), m.Status)
	}
}
