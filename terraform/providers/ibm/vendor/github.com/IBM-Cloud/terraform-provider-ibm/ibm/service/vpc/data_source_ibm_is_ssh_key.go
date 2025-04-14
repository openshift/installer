// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISSSHKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISSSHKeyRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group ID",
			},

			isKeyName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the ssh key",
			},

			isKeyType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ssh key type",
			},

			isKeyFingerprint: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ssh key Fingerprint",
			},

			isKeyPublicKey: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SSH Public key data",
			},

			isKeyLength: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ssh key length",
			},

			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			IsKeyCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},

			isKeyAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},
		},
	}
}

func dataSourceIBMISSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get(isKeyName).(string)

	err := keyGetByName(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func keyGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listKeysOptions := &vpcv1.ListKeysOptions{}

	start := ""
	allrecs := []vpcv1.Key{}
	for {
		if start != "" {
			listKeysOptions.Start = &start
		}

		keys, response, err := sess.ListKeys(listKeysOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error fetching Keys %s\n%s", err, response)
		}
		start = flex.GetNext(keys.Next)
		allrecs = append(allrecs, keys.Keys...)
		if start == "" {
			break
		}
	}

	for _, key := range allrecs {
		if *key.Name == name {
			d.SetId(*key.ID)
			d.Set("name", *key.Name)
			d.Set(isKeyType, *key.Type)
			d.Set(isKeyFingerprint, *key.Fingerprint)
			d.Set(isKeyLength, *key.Length)
			controller, err := flex.GetBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(flex.ResourceControllerURL, controller+"/vpc/compute/sshKeys")
			d.Set(flex.ResourceName, *key.Name)
			d.Set(flex.ResourceCRN, *key.CRN)
			d.Set(IsKeyCRN, *key.CRN)
			if key.ResourceGroup != nil {
				d.Set(flex.ResourceGroupName, *key.ResourceGroup.ID)
			}
			if key.PublicKey != nil {
				d.Set(isKeyPublicKey, *key.PublicKey)
			}
			accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *key.CRN, "", isKeyAccessTagType)
			if err != nil {
				log.Printf(
					"Error on get of resource SSH Key (%s) access tags: %s", d.Id(), err)
			}
			d.Set(isKeyAccessTags, accesstags)
			return nil
		}
	}
	return fmt.Errorf("[ERROR] No SSH Key found with name %s", name)
}
