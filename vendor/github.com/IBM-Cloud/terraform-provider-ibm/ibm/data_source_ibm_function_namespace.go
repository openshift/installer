// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/api/functions"
)

func dataSourceIBMFunctionNamespace() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMFunctionNamespaceRead,
		Schema: map[string]*schema.Schema{
			funcNamespaceName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of namespace.",
				ValidateFunc: InvokeValidator("ibm_function_namespace", funcNamespaceName),
			},
			funcNamespaceDesc: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Namespace Description.",
			},
			funcNamespaceResGrpId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource Group ID.",
			},
			funcNamespaceLoc: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Namespace Location.",
			},
		},
	}
}

func dataSourceIBMFunctionNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	functionNamespaceAPI, err := meta.(ClientSession).FunctionIAMNamespaceAPI()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	nsList, err := functionNamespaceAPI.Namespaces().GetNamespaces()
	if err != nil {
		return err
	}
	for _, n := range nsList.Namespaces {
		if n.Name != nil && *n.Name == name {
			getOptions := functions.GetNamespaceOptions{
				ID: n.ID,
			}

			instance, err := functionNamespaceAPI.Namespaces().GetNamespace(getOptions)
			if err != nil {
				d.SetId("")
				return nil
			}

			if instance.ID != nil {
				d.SetId(*instance.ID)
			}

			if instance.Name != nil {
				d.Set(funcNamespaceName, *instance.Name)
			}

			if instance.ResourceGroupID != nil {
				d.Set(funcNamespaceResGrpId, *instance.ResourceGroupID)
			}

			if instance.Location != nil {
				d.Set(funcNamespaceLoc, *instance.Location)
			}

			if instance.Description != nil {
				d.Set(funcNamespaceDesc, *instance.Description)
			}

			return nil
		}
	}

	return fmt.Errorf("No cloud function namespace found with name [%s]", name)
}
