package libvirt

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIgnition() *schema.Resource {
	return &schema.Resource{
		Create: resourceIgnitionCreate,
		Read:   resourceIgnitionRead,
		Delete: resourceIgnitionDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pool": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIgnitionCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] creating ignition file")
	client := meta.(*Client)
	if client.libvirt == nil {
		return fmt.Errorf(LibVirtConIsNil)
	}

	ignition := newIgnitionDef()

	ignition.Name = d.Get("name").(string)
	ignition.PoolName = d.Get("pool").(string)
	ignition.Content = d.Get("content").(string)

	log.Printf("[INFO] ignition: %+v", ignition)

	key, err := ignition.CreateAndUpload(client)
	if err != nil {
		return err
	}
	d.SetId(key)

	return resourceIgnitionRead(d, meta)
}

func resourceIgnitionRead(d *schema.ResourceData, meta interface{}) error {
	virConn := meta.(*Client).libvirt
	if virConn == nil {
		return fmt.Errorf(LibVirtConIsNil)
	}

	ign, err := newIgnitionDefFromRemoteVol(virConn, d.Id())
	d.Set("pool", ign.PoolName)
	d.Set("name", ign.Name)

	if err != nil {
		return fmt.Errorf("error while retrieving remote volume: %s", err)
	}

	return nil
}

func resourceIgnitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	if client.libvirt == nil {
		return fmt.Errorf(LibVirtConIsNil)
	}

	key, err := getIgnitionVolumeKeyFromTerraformID(d.Id())
	if err != nil {
		return err
	}

	return volumeDelete(client, key)
}
