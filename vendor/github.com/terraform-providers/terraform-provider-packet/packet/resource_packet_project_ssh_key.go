package packet

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		Read:   resourcePacketSSHKeyRead,
		Update: resourcePacketSSHKeyUpdate,
		Delete: resourcePacketSSHKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: pkeySchema,
	}
}
