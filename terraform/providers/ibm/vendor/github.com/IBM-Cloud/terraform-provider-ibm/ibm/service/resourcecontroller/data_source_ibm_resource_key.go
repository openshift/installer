// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package resourcecontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMResourceKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMResourceKeyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the resource key",
				Type:        schema.TypeString,
				Required:    true,
			},

			"resource_instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The id of the resource instance",
				ConflictsWith: []string{"resource_alias_id"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_resource_key",
					"resource_instance_id"),
			},

			"resource_alias_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The id of the resource alias",
				ConflictsWith: []string{"resource_instance_id"},
			},

			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User role",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of resource key",
			},

			"credentials": {
				Description: "Credentials asociated with the key",
				Sensitive:   true,
				Type:        schema.TypeMap,
				Computed:    true,
			},

			"credentials_json": {
				Description: "Credentials asociated with the key in json string",
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
			},

			"most_recent": {
				Description: "If true and multiple entries are found, the most recently created resource key is used. " +
					"If false, an error is returned",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "crn of resource key",
			},
		},
	}
}
func DataSourceIBMResourceKeyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "resource_instance_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "ResourceInstance",
			CloudDataRange:             []string{"service:%s"},
			Optional:                   true})

	ibmDataSourceKeyResourceValidator := validate.ResourceValidator{ResourceName: "ibm_resource_key", Schema: validateSchema}
	return &ibmDataSourceKeyResourceValidator
}

func dataSourceIBMResourceKeyRead(d *schema.ResourceData, meta interface{}) error {
	rsContClient, err := meta.(conns.ClientSession).ResourceControllerAPI()
	if err != nil {
		return err
	}
	rkAPI := rsContClient.ResourceServiceKey()
	name := d.Get("name").(string)
	mostRecent := d.Get("most_recent").(bool)

	keys, err := rkAPI.GetKeys(name)
	if err != nil {
		return err
	}
	var filteredKeys []models.ServiceKey

	if d.Get("resource_instance_id") == "" {
		filteredKeys = keys
	} else {
		crn, err := getCRN(d, meta)
		if err != nil {
			return err
		}
		for _, key := range keys {
			if key.SourceCrn == *crn {
				filteredKeys = append(filteredKeys, key)
			}
		}

	}

	if len(filteredKeys) == 0 {
		return fmt.Errorf("[ERROR] No resource keys found with name [%s]", name)
	}

	var key models.ServiceKey

	if len(filteredKeys) > 1 {
		if mostRecent {
			key = mostRecentResourceKey(filteredKeys)
		} else {
			return fmt.Errorf("[ERROR] More than one resource key found with name matching [%s]. "+
				"Set 'most_recent' to true in your configuration to force the most recent resource key "+
				"to be used", name)
		}
	} else {
		key = filteredKeys[0]
	}

	d.SetId(key.ID)

	if redacted, ok := key.Credentials["redacted"].(string); ok {
		log.Printf("Credentials are redacted with code: %s.The User doesn't have the correct access to view the credentials. Refer to the API documentation for additional details.", redacted)
	}
	if roleCrn, ok := key.Parameters["role_crn"].(string); ok {
		d.Set("role", roleCrn[strings.LastIndex(roleCrn, ":")+1:])
	} else if roleCrn, ok := key.Credentials["iam_role_crn"].(string); ok {
		d.Set("role", roleCrn[strings.LastIndex(roleCrn, ":")+1:])
	}

	d.Set("credentials", flex.Flatten(key.Credentials))
	creds, err := json.Marshal(key.Credentials)
	if err != nil {
		return fmt.Errorf("[ERROR] Error marshalling resource key credentials: %s", err)
	}
	if err = d.Set("credentials_json", string(creds)); err != nil {
		return fmt.Errorf("[ERROR] Error setting the credentials json: %s", err)
	}
	d.Set("status", key.State)
	d.Set("crn", key.Crn.String())
	return nil
}

func getCRN(d *schema.ResourceData, meta interface{}) (*crn.CRN, error) {

	rsContClient, err := meta.(conns.ClientSession).ResourceControllerAPI()
	if err != nil {
		return nil, err
	}

	if insID, ok := d.GetOk("resource_instance_id"); ok {
		instance, err := rsContClient.ResourceServiceInstance().GetInstance(insID.(string))
		if err != nil {
			return nil, err
		}
		return &(instance.Crn), nil

	}

	alias, err := rsContClient.ResourceServiceAlias().Alias(d.Get("resource_alias_id").(string))
	if err != nil {
		return nil, err
	}
	return &(alias.CRN), nil

}

type resourceKeys []models.ServiceKey

func (k resourceKeys) Len() int { return len(k) }

func (k resourceKeys) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

func (k resourceKeys) Less(i, j int) bool {
	return (*k[i].CreatedAt).Before(*k[j].CreatedAt)
}

func mostRecentResourceKey(keys resourceKeys) models.ServiceKey {
	sortedKeys := keys
	sort.Sort(sortedKeys)
	return sortedKeys[len(sortedKeys)-1]
}
