package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIdentityAuthScopeV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityAuthScopeV3Read,

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

			"service_catalog": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"interface": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityAuthScopeV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	d.SetId(d.Get("name").(string))

	tokenDetails, err := getTokenDetails(identityClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("user_name", tokenDetails.user.Name)
	d.Set("user_id", tokenDetails.user.ID)
	d.Set("user_domain_name", tokenDetails.user.Domain.Name)
	d.Set("user_domain_id", tokenDetails.user.Domain.ID)

	if tokenDetails.domain != nil {
		d.Set("domain_name", tokenDetails.domain.Name)
		d.Set("domain_id", tokenDetails.domain.ID)
	} else {
		d.Set("domain_name", "")
		d.Set("domain_id", "")
	}

	if tokenDetails.project != nil {
		d.Set("project_name", tokenDetails.project.Name)
		d.Set("project_id", tokenDetails.project.ID)
		d.Set("project_domain_name", tokenDetails.project.Domain.Name)
		d.Set("project_domain_id", tokenDetails.project.Domain.ID)
	} else {
		d.Set("project_name", "")
		d.Set("project_id", "")
		d.Set("project_domain_name", "")
		d.Set("project_domain_id", "")
	}

	allRoles := flattenIdentityAuthScopeV3Roles(tokenDetails.roles)
	if err := d.Set("roles", allRoles); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_identity_auth_scope_v3 roles: %s", err)
	}

	if tokenDetails.catalog != nil {
		flatCatalog := flattenIdentityAuthScopeV3ServiceCatalog(tokenDetails.catalog)
		if err := d.Set("service_catalog", flatCatalog); err != nil {
			log.Printf("[DEBUG] Unable to set openstack_identity_auth_scope_v3 service_catalog: %s", err)
		}
	}

	d.Set("region", GetRegion(d, config))

	return nil
}
