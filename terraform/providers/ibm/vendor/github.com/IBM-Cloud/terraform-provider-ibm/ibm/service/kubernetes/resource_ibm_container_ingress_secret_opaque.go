// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMContainerIngressSecretOpaque() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerIngressSecretOpaqueCreate,
		Read:     resourceIBMContainerIngressSecretOpaqueRead,
		Update:   resourceIBMContainerIngressSecretOpaqueUpdate,
		Delete:   resourceIBMContainerIngressSecretOpaqueDelete,
		Exists:   resourceIBMContainerIngressSecretOpaqueExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID or name",
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_ingress_secret_opaque",
					"cluster"),
			},
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret name",
				ForceNew:    true,
				ValidateFunc: validate.InvokeValidator(
					"ibm_container_ingress_secret_opaque",
					"secret_name"),
			},
			"secret_namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret namespace",
				ForceNew:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Opaque secret type",
			},
			"persistence": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Persistence of secret",
			},
			"user_managed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the secret was created by the user",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the secret",
			},
			"update_secret": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Updates secret from secrets manager if value is changed (increment each usage)",
			},
			"last_updated_timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp secret was last updated",
			},
			"fields": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Fields of an opaque secret",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Secret CRN corresponding to the field",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The computed field name",
						},
						"field_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The requested field name",
						},
						"expires_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field expires on date",
						},
						"last_updated_timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Field last updated timestamp",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMContainerIngressSecretOpaqueValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cluster",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			Required:                   true,
			CloudDataType:              "cluster",
			CloudDataRange:             []string{"resolved_to:id"}})

	validateSchema = append(validateSchema, validate.ValidateSchema{
		Identifier:                 "secret_name",
		ValidateFunctionIdentifier: validate.ValidateRegexpLen,
		Type:                       validate.TypeString,
		Required:                   true,
		Regexp:                     `^([a-z0-9]([-a-z0-9]*[a-z0-9])?(.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)$`,
		MinValueLength:             1,
		MaxValueLength:             63,
	})

	iBMContainerIngressInstanceValidator := validate.ResourceValidator{ResourceName: "ibm_container_ingress_secret_opaque", Schema: validateSchema}
	return &iBMContainerIngressInstanceValidator
}

func resourceIBMContainerIngressSecretOpaqueCreate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	cluster := d.Get("cluster").(string)
	secretName := d.Get("secret_name").(string)
	secretNamespace := d.Get("secret_namespace").(string)
	persistence := d.Get("persistence").(bool)

	params := v2.SecretCreateConfig{
		Cluster:     cluster,
		Name:        secretName,
		Namespace:   secretNamespace,
		Type:        "Opaque",
		Persistence: persistence,
	}

	if fields, ok := d.GetOk("fields"); ok {
		fieldList := fields.(*schema.Set).List()

		fieldsToAdd := []containerv2.FieldAdd{}
		for _, opaqueField := range fieldList {
			var fieldToAdd containerv2.FieldAdd

			opaqueFieldMap := opaqueField.(map[string]interface{})

			if fieldName, ok := opaqueFieldMap["field_name"]; ok {
				fieldToAdd.Name = fieldName.(string)
			}

			if crn, ok := opaqueFieldMap["crn"]; ok {
				fieldToAdd.CRN = crn.(string)
			}

			fieldsToAdd = append(fieldsToAdd, fieldToAdd)
		}

		params.FieldsToAdd = fieldsToAdd
	}

	ingressAPI := ingressClient.Ingresses()
	response, err := ingressAPI.CreateIngressSecret(params)

	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", cluster, response.Name, response.Namespace))

	return resourceIBMContainerIngressSecretOpaqueRead(d, meta)
}

func resourceIBMContainerIngressSecretOpaqueRead(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	ingressAPI := ingressClient.Ingresses()
	ingressSecretConfig, err := ingressAPI.GetIngressSecret(cluster, secretName, secretNamespace)
	if err != nil {
		return err
	}

	d.Set("cluster", cluster)
	d.Set("secret_name", ingressSecretConfig.Name)
	d.Set("secret_namespace", ingressSecretConfig.Namespace)
	d.Set("type", ingressSecretConfig.Type)
	d.Set("persistence", ingressSecretConfig.Persistence)
	d.Set("user_managed", ingressSecretConfig.UserManaged)
	d.Set("status", ingressSecretConfig.Status)
	d.Set("last_updated_timestamp", ingressSecretConfig.LastUpdatedTimestamp)

	if len(ingressSecretConfig.Fields) > 0 {
		d.Set("fields", flex.FlattenOpaqueSecret(ingressSecretConfig.Fields))
	}
	return nil
}

func resourceIBMContainerIngressSecretOpaqueDelete(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}

	ingressAPI := ingressClient.Ingresses()

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	params := v2.SecretDeleteConfig{
		Cluster:   cluster,
		Name:      secretName,
		Namespace: secretNamespace,
	}

	err = ingressAPI.DeleteIngressSecret(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceIBMContainerIngressSecretOpaqueUpdate(d *schema.ResourceData, meta interface{}) error {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	params := v2.SecretUpdateConfig{
		Cluster:   cluster,
		Name:      secretName,
		Namespace: secretNamespace,
	}

	ingressAPI := ingressClient.Ingresses()
	if d.HasChange("fields") {
		oldList, newList := d.GetChange("fields")

		if oldList == nil {
			oldList = new(schema.Set)
		}
		if newList == nil {
			newList = new(schema.Set)
		}
		os := oldList.(*schema.Set)
		ns := newList.(*schema.Set)

		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(remove) > 0 {
			actualSecret, err := ingressAPI.GetIngressSecret(cluster, secretName, secretNamespace)
			if err != nil {
				return err
			}

			// Loop through all actual secret fields to remove by CRN
			// As secret CRN to fields do not have 1:1 matchings
			actualFieldMap := make(map[string][]string)
			for _, actualField := range actualSecret.Fields {
				field, ok := actualFieldMap[actualField.CRN]
				if ok {
					actualFieldMap[actualField.CRN] = append(field, actualField.Name)
				} else {
					actualFieldMap[actualField.CRN] = []string{actualField.Name}
				}
			}

			for _, removeField := range remove {
				removeFieldMap := removeField.(map[string]interface{})

				if crn, ok := removeFieldMap["crn"]; ok {
					crnToRemove := crn.(string)
					secretNames, ok := actualFieldMap[crnToRemove]

					if ok {
						for _, secretName := range secretNames {
							params.FieldsToRemove = append(params.FieldsToRemove, containerv2.FieldRemove{
								Name: secretName,
							})
						}
					}
				}
			}

			_, err = ingressAPI.RemoveIngressSecretField(params)
			if err != nil {
				return err
			}
		}

		if len(add) > 0 {
			params.FieldsToAdd = nil

			for _, addField := range add {
				addFieldMap := addField.(map[string]interface{})

				var fieldToAdd containerv2.FieldAdd

				if fieldName, ok := addFieldMap["field_name"]; ok {
					fieldToAdd.Name = fieldName.(string)
				}

				if crn, ok := addFieldMap["crn"]; ok {
					fieldToAdd.CRN = crn.(string)
				}

				params.FieldsToAdd = append(params.FieldsToAdd, fieldToAdd)
			}

			_, err = ingressAPI.AddIngressSecretField(params)
			if err != nil {
				return err
			}
		}
	} else if d.HasChange("update_secret") {
		// user wants to force an upstream secret update from secrets manager onto kube cluster w/out changing crn
		_, err = ingressAPI.UpdateIngressSecret(params)
		if err != nil {
			return err
		}
	}
	return resourceIBMContainerIngressSecretOpaqueRead(d, meta)
}

func resourceIBMContainerIngressSecretOpaqueExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ingressClient, err := meta.(conns.ClientSession).VpcContainerAPI()
	if err != nil {
		return false, err
	}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	cluster := parts[0]
	secretName := parts[1]
	secretNamespace := parts[2]

	ingressAPI := ingressClient.Ingresses()
	ingressSecretConfig, err := ingressAPI.GetIngressSecret(cluster, secretName, secretNamespace)

	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("[ERROR] Error getting ingress secret: %s", err)
	}

	return ingressSecretConfig.Name == secretName && ingressSecretConfig.Namespace == secretNamespace && ingressSecretConfig.Status != "deleted", nil
}
