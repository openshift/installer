package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cloudidentity "google.golang.org/api/cloudidentity/v1beta1"
)

func dataSourceGoogleCloudIdentityGroups() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := datasourceSchemaFromResourceSchema(resourceCloudIdentityGroup().Schema)

	return &schema.Resource{
		Read: dataSourceGoogleCloudIdentityGroupsRead,

		Schema: map[string]*schema.Schema{
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `List of Cloud Identity groups.`,
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `The resource name of the entity under which this Group resides in the
Cloud Identity resource hierarchy.

Must be of the form identitysources/{identity_source_id} for external-identity-mapped
groups or customers/{customer_id} for Google Groups.`,
			},
		},
	}
}

func dataSourceGoogleCloudIdentityGroupsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	resp, err := config.clientCloudIdentity.Groups.List().Parent(d.Get("parent").(string)).View("FULL").Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("CloudIdentityGroups %q", d.Id()))
	}

	result := []map[string]interface{}{}
	for _, group := range resp.Groups {
		result = append(result, map[string]interface{}{
			"name":         group.Name,
			"display_name": group.DisplayName,
			"labels":       group.Labels,
			"description":  group.Description,
			"group_key":    flattenCloudIdentityGroupsEntityKey(group.GroupKey),
		})
	}

	d.Set("groups", result)
	d.SetId(time.Now().UTC().String())
	return nil
}

func flattenCloudIdentityGroupsEntityKey(entityKey *cloudidentity.EntityKey) []interface{} {
	transformed := map[string]interface{}{
		"id":        entityKey.Id,
		"namespace": entityKey.Namespace,
	}
	return []interface{}{transformed}
}
