package packet

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketProjectSSHKey() *schema.Resource {
	pkeySchema := packetSSHKeyCommonFields()
	pkeySchema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		ForceNew: true,
		Required: true,
	}
	return &schema.Resource{
		Create: resourcePacketSSHKeyCreate,
		Read:   resourcePacketProjectSSHKeyRead,
		Update: resourcePacketSSHKeyUpdate,
		Delete: resourcePacketSSHKeyDelete,

		Schema: pkeySchema,
	}
}

func resourcePacketProjectSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	projectID := d.Get("project_id").(string)
	projectKeys, _, err := client.SSHKeys.ProjectList(projectID)
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	keyFound := false
	for _, k := range projectKeys {
		if k.ID == d.Id() {
			keyFound = true
			d.Set("name", k.Label)
			d.Set("public_key", k.Key)
			d.Set("fingerprint", k.FingerPrint)
			d.Set("created", k.Created)
			d.Set("updated", k.Updated)
		}
	}
	if !keyFound {
		d.SetId("")
	}
	return nil
}
