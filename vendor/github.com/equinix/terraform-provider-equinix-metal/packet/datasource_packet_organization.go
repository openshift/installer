package packet

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourcePacketOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketOrganizationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"organization_id"},
			},
			"organization_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"website": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"twitter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func findOrgByName(os []packngo.Organization, name string) (*packngo.Organization, error) {
	results := make([]packngo.Organization, 0)
	for _, o := range os {
		if o.Name == name {
			results = append(results, o)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no organization found with name %s", name)
	}
	return nil, fmt.Errorf("too many organizations found with name %s (found %d, expected 1)", name, len(results))
}

func dataSourcePacketOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	nameRaw, nameOK := d.GetOk("name")
	orgIdRaw, orgIdOK := d.GetOk("organization_id")

	if !orgIdOK && !nameOK {
		return fmt.Errorf("You must supply organization_id or name")
	}
	var org *packngo.Organization

	if nameOK {
		name := nameRaw.(string)

		os, _, err := client.Organizations.List(nil)
		if err != nil {
			return err
		}

		org, err = findOrgByName(os, name)
		if err != nil {
			return err
		}
	} else {
		orgId := orgIdRaw.(string)
		log.Println(orgId)
		var err error
		org, _, err = client.Organizations.Get(orgId, nil)
		if err != nil {
			return err
		}
	}
	projectIds := []string{}

	for _, p := range org.Projects {
		projectIds = append(projectIds, filepath.Base(p.URL))
	}

	d.Set("organization_id", org.ID)
	d.Set("name", org.Name)
	d.Set("description", org.Description)
	d.Set("website", org.Website)
	d.Set("twitter", org.Twitter)
	d.Set("logo", org.Logo)
	d.Set("project_ids", projectIds)
	d.SetId(org.ID)

	return nil
}
