package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmUserAssignedIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmUserAssignedIdentityRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 24),
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmUserAssignedIdentityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).userAssignedIdentitiesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on User Assigned Identity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.IdentityProperties; props != nil {
		if principalId := props.PrincipalID; principalId != nil {
			d.Set("principal_id", principalId.String())
		}

		if clientId := props.ClientID; clientId != nil {
			d.Set("client_id", clientId.String())
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
