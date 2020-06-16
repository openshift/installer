package packet

import (
	"fmt"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

func dataSourcePacketProjectSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketProjectSSHKeyRead,
		Schema: map[string]*schema.Schema{
			"search": {
				Type:         schema.TypeString,
				Description:  "The name, fingerprint, id, or public_key of the SSH Key to search for in the Equinix Metal project",
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"id": {
				Type:         schema.TypeString,
				Description:  "The id of the SSH Key",
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				Computed:     true,
			},
			"project_id": {
				Type:         schema.TypeString,
				Description:  "The Equinix Metal project id of the Equinix Metal SSH Key",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The label of the Equinix Metal SSH Key",
				Computed:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "The public SSH key that will be authorized for SSH access on Equinix Metal devices provisioned with this key",
				Computed:    true,
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePacketProjectSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	search := d.Get("search").(string)
	id := d.Get("id").(string)
	projectID := d.Get("project_id").(string)

	if id == "" && search == "" {
		return fmt.Errorf("You must supply either search or id")
	}

	var (
		key        packngo.SSHKey
		searchOpts *packngo.SearchOptions
	)

	if search != "" {
		searchOpts = &packngo.SearchOptions{Search: search}
	}
	keys, _, err := client.Projects.ListSSHKeys(projectID, searchOpts)

	if err != nil {
		err = fmt.Errorf("Error listing project ssh keys: %s", friendlyError(err))
		return err
	}

	for i := range keys {
		// use the first match for searches
		if search != "" {
			key = keys[i]
			break
		}

		// otherwise find the matching ID
		if keys[i].ID == id {
			key = keys[i]
			break
		}
	}

	if key.ID == "" {
		// Not Found
		return fmt.Errorf("Project %q SSH Key matching %q was not found", projectID, search)
	}

	ownerID := path.Base(key.Owner.Href)

	d.SetId(key.ID)
	d.Set("name", key.Label)
	d.Set("public_key", key.Key)
	d.Set("fingerprint", key.FingerPrint)
	d.Set("owner_id", ownerID)
	d.Set("created", key.Created)
	d.Set("updated", key.Updated)

	if key.Owner.Href[:10] == "/projects/" {
		d.Set("project_id", ownerID)
	}

	return nil
}
