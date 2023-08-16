// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package hpcs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-hpcs-uko-sdk/ukov4"
)

func ResourceIbmVault() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIbmVaultCreate,
		ReadContext:   ResourceIbmVaultRead,
		UpdateContext: ResourceIbmVaultUpdate,
		DeleteContext: ResourceIbmVaultDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the UKO instance this resource exists in.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region of the UKO instance this resource exists in.",
			},
			"vault_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the vault.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_hpcs_vault", "name"),
				Description:  "A human-readable name to assign to your vault. To protect your privacy, do not use personal data, such as your name or location.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_hpcs_vault", "description"),
				Description:  "Description of the vault.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the vault was created.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and time when the vault was last updated.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that created the vault.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the user that last updated the vault.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL that uniquely identifies your cloud resource.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmVaultValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[A-Za-z0-9#@!$%'_-][A-Za-z0-9#@!$% '_-]*$`,
			MinValueLength:             1,
			MaxValueLength:             100,
		},
		validate.ValidateSchema{
			Identifier:                 "region",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "au-syd, in-che, jp-osa, jp-tok, kr-seo, eu-de, eu-gb, ca-tor, us-south, us-south-test, us-east, br-sao",
		},
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `(.|\n)*`,
			// TODO: Old regex, in case there are problems
			// Regexp:                     `(.|\\n)*`,
			MinValueLength: 0,
			MaxValueLength: 200,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_hpcs_vault", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIbmVaultCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	instance_id := d.Get("instance_id").(string)
	region := d.Get("region").(string)

	createVaultOptions := &ukov4.CreateVaultOptions{}

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	createVaultOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("description"); ok {
		createVaultOptions.SetDescription(d.Get("description").(string))
	}

	vault, response, err := ukoClient.CreateVaultWithContext(context, createVaultOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateVaultWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateVaultWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instance_id, *vault.ID))

	return ResourceIbmVaultRead(context, d, meta)
}

func ResourceIbmVaultRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	getVaultOptions := &ukov4.GetVaultOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	getVaultOptions.SetID(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	vault, response, err := ukoClient.GetVaultWithContext(context, getVaultOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVaultWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVaultWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("vault_id", vault.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vault ID: %s", err))
	}
	if err = d.Set("name", vault.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", vault.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(vault.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(vault.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("created_by", vault.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("updated_by", vault.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}
	if err = d.Set("href", vault.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func ResourceIbmVaultUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	updateVaultOptions := &ukov4.UpdateVaultOptions{}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	updateVaultOptions.SetID(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	hasChange := false

	if d.HasChange("name") {
		updateVaultOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateVaultOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}

	// Etag support
	updateVaultOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		_, response, err := ukoClient.UpdateVaultWithContext(context, updateVaultOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateVaultWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateVaultWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIbmVaultRead(context, d, meta)
}

func ResourceIbmVaultDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ukoClient, err := meta.(conns.ClientSession).UkoV4()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteVaultOptions := &ukov4.DeleteVaultOptions{}

	// Etag support
	deleteVaultOptions.SetIfMatch(d.Get("etag").(string))

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instance_id := id[1]
	vault_id := id[2]
	deleteVaultOptions.SetID(vault_id)

	url, err := getUkoUrl(context, region, instance_id, ukoClient)
	if err != nil {
		return diag.FromErr(err)
	}
	ukoClient.SetServiceURL(url)

	response, err := ukoClient.DeleteVaultWithContext(context, deleteVaultOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteVaultWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteVaultWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func getUkoUrl(context context.Context, region string, instance_id string, ukoClient *ukov4.UkoV4) (string, error) {

	// Get url
	pathParamsMap := map[string]string{
		"id": instance_id,
	}

	var brokerUrl string
	if v := os.Getenv("IBMCLOUD_HPCS_UKO_URL"); v != "" {
		if strings.Contains(v, "https://") {
			return v, nil
		} else {
			return fmt.Sprintf("https://%s/", v), nil
		}
	} else {
		if os.Getenv("IBMCLOUD_IAM_API_ENDPOINT") == "https://iam.test.cloud.ibm.com" {
			if region == "us-south" {
				brokerUrl = "https://broker.us-south.hs-crypto.test.cloud.ibm.com"
			} else if region == "us-east-dev" {
				brokerUrl = "https://broker.us-east.hs-crypto.test.cloud.ibm.com"
			} else if region == "us-south-svt" {
				brokerUrl = "https://broker.svt.us-south.hs-crypto.test.cloud.ibm.com"
			} else if region == "us-south-test" {
				brokerUrl = "https://broker.vpc.us-south.hs-crypto.test.cloud.ibm.com"
			}
		} else {
			brokerUrl = fmt.Sprintf("https://%s.broker.hs-crypto.cloud.ibm.com", region)
		}
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(context)
	builder.EnableGzipCompression = ukoClient.GetEnableGzipCompression()
	_, err := builder.ResolveRequestURL(brokerUrl,
		`/crypto_v2/instances/{id}`, pathParamsMap)
	if err != nil {
		return "", err
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return "", err
	}

	var rawResponse map[string]json.RawMessage
	_, err = ukoClient.Service.Request(request, &rawResponse)
	if err != nil {
		return "", err
	}

	var uko *map[string]string
	if rawResponse != nil {
		err = core.UnmarshalPrimitive(rawResponse, "uko", &uko)
		if err != nil {
			return "", err
		}
		url := (*uko)["public"]
		log.Printf(url)
		return "https://" + url, nil
	}

	return "", fmt.Errorf("Could not get response")
}
