package packet

import (
	"fmt"
	"log"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

func dataSourcePacketProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id"},
			},
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"backend_transfer": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"payment_method_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"organization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bgp_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"local", "global"}, false),
						},
						"asn": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"md5": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_prefix": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePacketProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	nameRaw, nameOK := d.GetOk("name")
	projectIdRaw, projectIdOK := d.GetOk("project_id")

	if !projectIdOK && !nameOK {
		return fmt.Errorf("You must supply project_id or name")
	}
	var project *packngo.Project

	if nameOK {
		name := nameRaw.(string)

		os, _, err := client.Projects.List(nil)
		if err != nil {
			return err
		}

		project, err = findProjectByName(os, name)
		if err != nil {
			return err
		}
	} else {
		projectId := projectIdRaw.(string)
		log.Println(projectId)
		var err error
		project, _, err = client.Projects.Get(projectId, nil)
		if err != nil {
			return err
		}
	}

	d.SetId(project.ID)
	d.Set("payment_method_id", path.Base(project.PaymentMethod.URL))
	d.Set("name", project.Name)
	d.Set("project_id", project.ID)
	d.Set("organization_id", path.Base(project.Organization.URL))
	d.Set("created", project.Created)
	d.Set("updated", project.Updated)
	d.Set("backend_transfer", project.BackendTransfer)

	bgpConf, _, err := client.BGPConfig.Get(project.ID, nil)
	userIds := []string{}
	for _, u := range project.Users {
		userIds = append(userIds, path.Base(u.URL))
	}
	d.Set("user_ids", userIds)

	if (err == nil) && (bgpConf != nil) {
		// guard against an empty struct
		if bgpConf.ID != "" {
			err := d.Set("bgp_config", flattenBGPConfig(bgpConf))
			if err != nil {
				err = friendlyError(err)
				return err
			}
		}
	}
	return nil
}

func findProjectByName(ps []packngo.Project, name string) (*packngo.Project, error) {
	results := make([]packngo.Project, 0)
	for _, p := range ps {
		if p.Name == name {
			results = append(results, p)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no project found with name %s", name)
	}
	return nil, fmt.Errorf("too many projects found with name %s (found %d, expected 1)", name, len(results))
}
