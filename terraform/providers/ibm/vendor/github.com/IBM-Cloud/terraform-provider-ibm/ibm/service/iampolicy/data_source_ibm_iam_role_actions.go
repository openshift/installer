// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIAMRoleAction() *schema.Resource {
	return &schema.Resource{
		Read: datasourceIBMIAMRoleActionRead,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Service Name",
				ForceNew:    true,
			},
			"reader": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Reader action ids",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"manager": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "manager action ids",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"reader_plus": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "readerplus action ids",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"writer": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "writer action ids",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"actions": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "List of actions for different services roles",
			},
		},
	}

}

func datasourceIBMIAMRoleActionRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	serviceName := d.Get("service").(string)
	d.SetId(serviceName)

	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
		ServiceName: &serviceName,
	}

	roleList, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
	if err != nil {
		return err
	}
	serviceRoles := roleList.ServiceRoles

	d.Set("reader", flex.FlattenActionbyDisplayName("Reader", serviceRoles))
	d.Set("manager", flex.FlattenActionbyDisplayName("Manager", serviceRoles))
	d.Set("reader_plus", flex.FlattenActionbyDisplayName("ReaderPlus", serviceRoles))
	d.Set("writer", flex.FlattenActionbyDisplayName("Writer", serviceRoles))
	d.Set("actions", flattenRoleActions(serviceRoles))

	return nil
}

func flattenRoleActions(object []iampolicymanagementv1.Role) map[string]string {
	actions := make(map[string]string)
	for _, item := range object {
		actions[*item.DisplayName] = strings.Join(item.Actions, ",")
	}
	return actions
}
