// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func datasourceIBMIAMRole() *schema.Resource {
	return &schema.Resource{
		Read: datasourceIBMIAMRoleRead,

		Schema: map[string]*schema.Schema{
			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Service Name",
				ForceNew:    true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

}

func datasourceIBMIAMRoleRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyManagementClient, err := meta.(ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return err
	}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}

	var serviceName string
	var customRoles []iampolicymanagementv1.CustomRole
	var serviceRoles, systemRoles []iampolicymanagementv1.Role

	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
		AccountID: &userDetails.userAccount,
	}

	if service, ok := d.GetOk("service"); ok {
		serviceName = service.(string)
		listRoleOptions.ServiceName = &serviceName
	}
	roleList, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
	if err != nil {
		return err
	}
	customRoles = roleList.CustomRoles
	serviceRoles = roleList.ServiceRoles
	systemRoles = roleList.SystemRoles

	d.SetId(userDetails.userAccount)

	var roles []map[string]string

	roles = append(flattenRoleData(systemRoles, "platform"), append(flattenRoleData(serviceRoles, "service"), flattenCustomRoleData(customRoles, "custom")...)...)

	d.Set("roles", roles)

	return nil
}
