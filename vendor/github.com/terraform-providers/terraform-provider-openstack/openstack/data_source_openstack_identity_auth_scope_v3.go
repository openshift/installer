package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func dataSourceIdentityAuthScopeV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityAuthScopeV3Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// computed attributes
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityAuthScopeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}
	tokenID := config.OsClient.TokenID

	d.SetId(d.Get("name").(string))

	result := tokens.Get(identityClient, tokenID)
	if result.Err != nil {
		return result.Err
	}

	user, err := result.ExtractUser()
	if err != nil {
		return err
	}

	d.Set("user_name", user.Name)
	d.Set("user_id", user.ID)
	d.Set("user_domain_name", user.Domain.Name)
	d.Set("user_domain_id", user.Domain.ID)

	domain, err := result.ExtractDomain()
	if err != nil {
		return err
	}
	if domain != nil {
		d.Set("domain_name", domain.Name)
		d.Set("domain_id", domain.ID)
	} else {
		d.Set("domain_name", "")
		d.Set("domain_id", "")
	}

	project, err := result.ExtractProject()
	if err != nil {
		return err
	}
	if project != nil {
		d.Set("project_name", project.Name)
		d.Set("project_id", project.ID)
		d.Set("project_domain_name", project.Domain.Name)
		d.Set("project_domain_id", project.Domain.ID)
	} else {
		d.Set("project_name", "")
		d.Set("project_id", "")
		d.Set("project_domain_name", "")
		d.Set("project_domain_id", "")
	}

	roles, err := result.ExtractRoles()
	if err != nil {
		return err
	}

	allRoles := flattenIdentityAuthScopeV3Roles(roles)
	if err := d.Set("roles", allRoles); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_identity_auth_scope_v3 roles: %s", err)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}
