// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/posturemanagementv2"
)

func resourceIBMSccPostureCollectors() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccPostureCollectorsCreate,
		ReadContext:   resourceIBMSccPostureCollectorsRead,
		UpdateContext: resourceIBMSccPostureCollectorsUpdate,
		DeleteContext: resourceIBMSccPostureCollectorsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_scc_posture_collector", "name"),
				Description:  "A unique name for your collector.",
			},
			"is_public": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Determines whether the collector endpoint is accessible on a public network. If set to `true`, the collector connects to resources in your account over a public network. If set to `false`, the collector connects to resources by using a private IP that is accessible only through the IBM Cloud private network.",
			},
			"managed_by": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_scc_posture_collector", "managed_by"),
				Description:  "Determines whether the collector is an IBM or customer-managed virtual machine. Use `ibm` to allow Security and Compliance Center to create, install, and manage the collector on your behalf. The collector is installed in an OpenShift cluster and approved automatically for use. Use `customer` if you would like to install the collector by using your own virtual machine. For more information, check out the [docs](https://cloud.ibm.com/docs/security-compliance?topic=security-compliance-collector).",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				ValidateFunc: InvokeValidator("ibm_scc_posture_collector", "description"),
				Description:  "A detailed description of the collector.",
			},
			"passphrase": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_scc_posture_collector", "passphrase"),
				Description:  "To protect the credentials that you add to the service, a passphrase is used to generate a data encryption key. The key is used to securely store your credentials and prevent anyone from accessing them.",
			},
			"is_ubi_image": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Determines whether the collector has a Ubi image.",
			},
		},
	}
}

func resourceIBMSccPostureCollectorsValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9-\\.,_\\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             32,
		},
		ValidateSchema{
			Identifier:                 "managed_by",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "customer, ibm",
		},
		ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9-\\._,\\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             1000,
		},
		ValidateSchema{
			Identifier:                 "passphrase",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9-\\._,\\s]*$`,
			MinValueLength:             1,
			MaxValueLength:             200,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_collectors", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSccPostureCollectorsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createCollectorOptions := &posturemanagementv2.CreateCollectorOptions{}
	createCollectorOptions.SetAccountID(os.Getenv("SCC_POSTURE_ACCOUNT_ID"))

	createCollectorOptions.SetName(d.Get("name").(string))
	createCollectorOptions.SetIsPublic(d.Get("is_public").(bool))
	createCollectorOptions.SetManagedBy(d.Get("managed_by").(string))
	if _, ok := d.GetOk("description"); ok {
		createCollectorOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("passphrase"); ok {
		createCollectorOptions.SetPassphrase(d.Get("passphrase").(string))
	}
	if _, ok := d.GetOk("is_ubi_image"); ok {
		createCollectorOptions.SetIsUbiImage(d.Get("is_ubi_image").(bool))
	}

	collector, response, err := postureManagementClient.CreateCollectorWithContext(context, createCollectorOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateCollectorWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateCollectorWithContext failed %s\n%s", err, response))
	}

	d.SetId(*collector.ID)

	return resourceIBMSccPostureCollectorsRead(context, d, meta)
}

func resourceIBMSccPostureCollectorsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	listCollectorsOptions := &posturemanagementv2.ListCollectorsOptions{}
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting userDetails %s", err))
	}

	accountID := userDetails.userAccount
	listCollectorsOptions.SetAccountID(accountID)

	collectorList, response, err := postureManagementClient.ListCollectorsWithContext(context, listCollectorsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] ListCollectorsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListCollectorsWithContext failed %s\n%s", err, response))
	}
	d.SetId(*(collectorList.Collectors[0].ID))
	return nil
}

func resourceIBMSccPostureCollectorsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateCollectorOptions := &posturemanagementv2.UpdateCollectorOptions{}
	updateCollectorOptions.SetAccountID(os.Getenv("SCC_POSTURE_ACCOUNT_ID"))

	updateCollectorOptions.SetID(d.Id())

	hasChange := false

	if hasChange {
		//updateCollectorOptions.CollectorUpdatePatch, _ = patchVals.AsPatch()
		_, response, err := postureManagementClient.UpdateCollectorWithContext(context, updateCollectorOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateCollectorWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateCollectorWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMSccPostureCollectorsRead(context, d, meta)
}

func resourceIBMSccPostureCollectorsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteCollectorOptions := &posturemanagementv2.DeleteCollectorOptions{}
	deleteCollectorOptions.SetAccountID(os.Getenv("SCC_POSTURE_ACCOUNT_ID"))

	deleteCollectorOptions.SetID(d.Id())

	response, err := postureManagementClient.DeleteCollectorWithContext(context, deleteCollectorOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteCollectorWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteCollectorWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
