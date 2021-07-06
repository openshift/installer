// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	tgGateways      = "transit_gateways"
	tgResourceGroup = "resource_group"
	tgID            = "id"
	tgCrn           = "crn"
	tgName          = "name"
	tgLocation      = "location"
	tgCreatedAt     = "created_at"
	tgGlobal        = "global"
	tgStatus        = "status"
	tgUpdatedAt     = "updated_at"
	tgGatewayTags   = "tags"

	isTransitGatewayProvisioning     = "provisioning"
	isTransitGatewayProvisioningDone = "done"
	isTransitGatewayDeleting         = "deleting"
	isTransitGatewayDeleted          = "done"
)

func resourceIBMTransitGateway() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMTransitGatewayCreate,
		Read:     resourceIBMTransitGatewayRead,
		Delete:   resourceIBMTransitGatewayDelete,
		Exists:   resourceIBMTransitGatewayExists,
		Update:   resourceIBMTransitGatewayUpdate,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			tgLocation: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Location of Transit Gateway Services",
			},

			tgName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				Description:  "Name Transit Gateway Services",
				ValidateFunc: InvokeValidator("ibm_tg_gateway", tgName),
			},

			tgGlobal: {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     false,
				Description: "Allow global routing for a Transit Gateway. If unspecified, the default value is false",
			},

			tgGatewayTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_tg_gateway", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "Tags for the transit gateway instance",
			},

			tgResourceGroup: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			tgCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},
			tgCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the resource",
			},
			tgUpdatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The updation time of the resource",
			},
			tgStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Status of the resource",
			},

			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the resource",
			},

			ResourceStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func resourceIBMTGValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 tgName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              `^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$}`,
			MinValueLength:             1,
			MaxValueLength:             63})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128})

	ibmTGResourceValidator := ResourceValidator{ResourceName: "ibm_tg_gateway", Schema: validateSchema}
	return &ibmTGResourceValidator
}

func transitgatewayClient(meta interface{}) (*transitgatewayapisv1.TransitGatewayApisV1, error) {
	sess, err := meta.(ClientSession).TransitGatewayV1API()
	return sess, err
}

func resourceIBMTransitGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	location := d.Get(tgLocation).(string)
	name := d.Get(tgName).(string)
	global := d.Get(tgGlobal).(bool)

	createTransitGatewayOptions := &transitgatewayapisv1.CreateTransitGatewayOptions{}

	createTransitGatewayOptions.Name = &name
	createTransitGatewayOptions.Location = &location
	createTransitGatewayOptions.Global = &global

	if rsg, ok := d.GetOk(tgResourceGroup); ok {
		resourceGroup := rsg.(string)
		createTransitGatewayOptions.ResourceGroup = &transitgatewayapisv1.ResourceGroupIdentity{ID: &resourceGroup}
	}

	//log.Println("going to create tgw now with options", *createTransitGatewayOptions.ResourceGroup)
	tgw, response, err := client.CreateTransitGateway(createTransitGatewayOptions)

	if err != nil {
		log.Printf("[DEBUG] Create Transit Gateway err %s\n%s", err, response)
		return err
	}
	d.SetId(*tgw.ID)

	_, err = isWaitForTransitGatewayAvailable(client, d.Id(), d.Timeout(schema.TimeoutCreate))

	if err != nil {
		return err
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(tgGatewayTags); ok || v != "" {
		oldList, newList := d.GetChange(tgGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *tgw.Crn)
		if err != nil {
			log.Printf(
				"Error on create of transit gateway (%s) tags: %s", d.Id(), err)
		}
	}
	return resourceIBMTransitGatewayRead(d, meta)
}

func isWaitForTransitGatewayAvailable(client *transitgatewayapisv1.TransitGatewayApisV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for transit gateway (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isTransitGatewayProvisioning},
		Target:     []string{isTransitGatewayProvisioningDone, ""},
		Refresh:    isTransitGatewayRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isTransitGatewayRefreshFunc(client *transitgatewayapisv1.TransitGatewayApisV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		gettgwoptions := &transitgatewayapisv1.GetTransitGatewayOptions{
			ID: &id,
		}
		transitGateway, response, err := client.GetTransitGateway(gettgwoptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting Transit Gateway: %s\n%s", err, response)
		}

		if *transitGateway.Status == "available" || *transitGateway.Status == "failed" {
			return transitGateway, isTransitGatewayProvisioningDone, nil
		}

		return transitGateway, isTransitGatewayProvisioning, nil
	}
}

func resourceIBMTransitGatewayRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}
	tgOptions := &transitgatewayapisv1.GetTransitGatewayOptions{}
	if id != "" {
		tgOptions.ID = &id
	}

	tgw, response, err := client.GetTransitGateway(tgOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(*tgw.ID)
	d.Set(tgCrn, tgw.Crn)
	d.Set(tgName, tgw.Name)
	d.Set(tgLocation, tgw.Location)
	d.Set(tgCreatedAt, tgw.CreatedAt.String())

	if tgw.UpdatedAt != nil {
		d.Set(tgUpdatedAt, tgw.UpdatedAt.String())
	}
	d.Set(tgGlobal, tgw.Global)
	d.Set(tgStatus, tgw.Status)

	tags, err := GetTagsUsingCRN(meta, *tgw.Crn)
	if err != nil {
		log.Printf(
			"Error on get of transit gateway (%s) tags: %s", d.Id(), err)
	}
	d.Set(tgGatewayTags, tags)

	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}

	d.Set(ResourceControllerURL, controller+"/interconnectivity/transit")
	d.Set(ResourceName, *tgw.Name)
	d.Set(ResourceCRN, *tgw.Crn)
	d.Set(ResourceStatus, *tgw.Status)
	if tgw.ResourceGroup != nil {
		rg := tgw.ResourceGroup
		d.Set(tgResourceGroup, *rg.ID)
		d.Set(ResourceGroupName, *rg.ID)
	}
	return nil
}

func resourceIBMTransitGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := transitgatewayClient(meta)

	if err != nil {
		return err
	}

	ID := d.Id()
	tgOptions := &transitgatewayapisv1.GetTransitGatewayOptions{
		ID: &ID,
	}
	tgw, resp, err := client.GetTransitGateway(tgOptions)

	if err != nil {
		log.Printf("Error fetching Tranisit  Gateway: %s", resp)
		return err
	}

	updateTransitGatewayOptions := &transitgatewayapisv1.UpdateTransitGatewayOptions{}
	updateTransitGatewayOptions.ID = &ID
	if d.HasChange(tgName) {
		if tgwname, ok := d.GetOk(tgName); ok {
			name := tgwname.(string)
			updateTransitGatewayOptions.Name = &name
		}
	}
	if d.HasChange(tgGlobal) {
		if tgwglobal, ok := d.GetOk(tgGlobal); ok {
			global := tgwglobal.(bool)
			updateTransitGatewayOptions.Global = &global
		}
	}
	if d.HasChange(tgGatewayTags) {
		oldList, newList := d.GetChange(tgGatewayTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *tgw.Crn)
		if err != nil {
			log.Printf(
				"Error on update of transit gateway (%s) tags: %s", ID, err)
		}
	}

	_, response, err := client.UpdateTransitGateway(updateTransitGatewayOptions)
	if err != nil {
		log.Printf("[DEBUG] Update Transit Gateway err %s\n%s", err, response)
		return err
	}

	return resourceIBMTransitGatewayRead(d, meta)
}

func resourceIBMTransitGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return err
	}

	ID := d.Id()
	delOptions := &transitgatewayapisv1.DeleteTransitGatewayOptions{
		ID: &ID,
	}
	response, err := client.DeleteTransitGateway(delOptions)

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("Error deleting Transit Gateway (%s): %s\n%s", ID, err, response)
	}
	_, err = isWaitForTransitGatewayDeleted(client, ID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForTransitGatewayDeleted(client *transitgatewayapisv1.TransitGatewayApisV1, id string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for transit gateway (%s) to be deleted.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", isTransitGatewayDeleting},
		Target:     []string{"", isTransitGatewayDeleted},
		Refresh:    isTransitGatewayDeleteRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isTransitGatewayDeleteRefreshFunc(client *transitgatewayapisv1.TransitGatewayApisV1, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] delete function here")
		gettgwoptions := &transitgatewayapisv1.GetTransitGatewayOptions{
			ID: &id,
		}
		transitGateway, response, err := client.GetTransitGateway(gettgwoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return transitGateway, isTransitGatewayDeleted, nil
			}
			return nil, "", fmt.Errorf("Error Getting Transit Gateway: %s\n%s", err, response)
		}
		return transitGateway, isTransitGatewayDeleting, err
	}
}

func resourceIBMTransitGatewayExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	client, err := transitgatewayClient(meta)
	if err != nil {
		return false, err
	}

	ID := d.Id()

	tgOptions := &transitgatewayapisv1.GetTransitGatewayOptions{}
	if ID != "" {
		tgOptions.ID = &ID
	}
	_, response, err := client.GetTransitGateway(tgOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return false, nil
		}
		return false, fmt.Errorf("Error Getting Transit Gateway: %s\n%s", err, response)
	}

	return true, nil
}
