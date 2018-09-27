package provider

import (
	"errors"
	"fmt"
	"strings"

	"github.com/src-d/terraform-provider-online-net/online"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataRescueImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceResourceImageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "exact name of the desired image",
				ConflictsWith: []string{"name_filter"},
			},
			"name_filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "partial name of the desired image to filter with, in case multiple get found the last one will be used",
				ConflictsWith: []string{"name"},
			},
			"server": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "server for the desired image",
			},
			"image": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the requested image",
			},
		},
	}
}

func dataSourceResourceImageRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(online.Client)

	server := d.Get("server").(int)
	nameInterface, hasName := d.GetOk("name")
	nameFilterInterface, hasNameFilter := d.GetOk("name_filter")

	if !hasName && !hasNameFilter {
		return errors.New("Need either a name or a name_filter")
	}

	images, err := c.GetRescueImages(server)
	if err != nil {
		return err
	}

	selectedImage := ""

	if hasName {
		want := nameInterface.(string)
		for _, image := range images {
			if image == want {
				selectedImage = image
				break
			}
		}
	} else if hasNameFilter {
		want := nameFilterInterface.(string)
		for _, image := range images {
			if strings.Contains(image, want) {
				selectedImage = image
			}
		}
	}

	if selectedImage == "" {
		return fmt.Errorf("No image found for requirements, options are: %s", strings.Join(images, ","))
	}

	d.Set("image", selectedImage)
	d.SetId(fmt.Sprintf("%d", server))

	return nil
}
