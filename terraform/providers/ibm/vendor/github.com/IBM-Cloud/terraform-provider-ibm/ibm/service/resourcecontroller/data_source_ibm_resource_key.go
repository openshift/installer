// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package resourcecontroller

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"

	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
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

			// ### Modification addded onetime_credentials to Resource scehama
			"onetime_credentials": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "onetime_credentials of resource key",
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
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:%s"},
			Optional:                   true})

	ibmDataSourceKeyResourceValidator := validate.ResourceValidator{ResourceName: "ibm_resource_key", Schema: validateSchema}
	return &ibmDataSourceKeyResourceValidator
}

func dataSourceIBMResourceKeyRead(d *schema.ResourceData, meta interface{}) error {
	rsContClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	mostRecent := d.Get("most_recent").(bool)

	resourceKeys := rc.ListResourceKeysOptions{
		Name: &name,
	}

	keys, _, err := rsContClient.ListResourceKeys(&resourceKeys)
	if err != nil {
		return err
	}
	var filteredKeys []rc.ResourceKey

	if d.Get("resource_instance_id") == "" {
		filteredKeys = keys.Resources
	} else {
		crn, err := getCRN(d, meta)
		if err != nil || crn == nil {
			return err
		}
		for _, key := range keys.Resources {
			if *key.SourceCRN == *crn {
				filteredKeys = append(filteredKeys, key)
			}
		}

	}

	if len(filteredKeys) == 0 {
		return fmt.Errorf("[ERROR] No resource keys found with name [%s]", name)
	}

	var key rc.ResourceKey

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

	d.SetId(*key.ID)

	if key.Credentials != nil && key.Credentials.Redacted != nil {
		log.Printf("Credentials are redacted with code: %s.The User doesn't have the correct access to view the credentials. Refer to the API documentation for additional details.", *key.Credentials.Redacted)
	}

	if key.Credentials != nil && key.Credentials.IamRoleCRN != nil {
		roleCrn := *key.Credentials.IamRoleCRN
		iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
		if err == nil {
			var resourceCRN string
			if key.CRN != nil {
				serviceName := strings.Split(*key.CRN, ":")
				if len(serviceName) > 4 {
					resourceCRN = serviceName[4]
				}
			}
			listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
				AccountID:   key.AccountID,
				ServiceName: &resourceCRN,
			}
			roleList, resp, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
			roles := flex.MapRoleListToPolicyRoles(*roleList)
			if err == nil && len(roles) > 0 {
				for _, role := range roles {
					if *role.RoleID == roleCrn {
						RoleName := role.DisplayName
						d.Set("role", RoleName)
					}
				}
			}
			if err != nil {
				log.Printf("[ERROR] Error listing IAM Roles %s, %s", err, resp)
			}
		}
	}

	// ### Modification for onetime_credientails
	d.Set("onetime_credentials", key.OnetimeCredentials)
	var credInterface map[string]interface{}
	cred, _ := json.Marshal(key.Credentials)
	json.Unmarshal(cred, &credInterface)
	d.Set("credentials", flex.Flatten(credInterface))

	creds, err := json.Marshal(key.Credentials)
	if err != nil {
		return fmt.Errorf("[ERROR] Error marshalling resource key credentials: %s", err)
	}
	if err = d.Set("credentials_json", string(creds)); err != nil {
		return fmt.Errorf("[ERROR] Error setting the credentials json: %s", err)
	}
	d.Set("status", key.State)
	d.Set("crn", key.CRN)
	return nil
}

func getCRN(d *schema.ResourceData, meta interface{}) (*string, error) {

	rsConClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return nil, err
	}

	if insID, ok := d.GetOk("resource_instance_id"); ok {
		insIdString := insID.(string)
		resourceInstanceGet := rc.GetResourceInstanceOptions{
			ID: &insIdString,
		}
		instance, resp, err := rsConClient.GetResourceInstance(&resourceInstanceGet)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
		}
		return instance.CRN, nil

	}
	if insID, ok := d.GetOk("resource_alias_id"); ok {
		insaliasIdString := insID.(string)
		resourceInstanceAliasGet := rc.GetResourceAliasOptions{
			ID: &insaliasIdString,
		}
		instance, resp, err := rsConClient.GetResourceAlias(&resourceInstanceAliasGet)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Error retrieving resource instance: %s with resp code: %s", err, resp)
		}
		return instance.CRN, nil

	}

	return nil, nil

}

type resourceKeys []rc.ResourceKey

func (k resourceKeys) Len() int { return len(k) }

func (k resourceKeys) Swap(i, j int) { k[i], k[j] = k[j], k[i] }

func (k resourceKeys) Less(i, j int) bool {
	return (time.Time(*k[i].CreatedAt)).Before(time.Time(*k[j].CreatedAt))
}

func mostRecentResourceKey(keys resourceKeys) rc.ResourceKey {
	sortedKeys := keys
	sort.Sort(sortedKeys)
	return sortedKeys[len(sortedKeys)-1]
}
