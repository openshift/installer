// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	isFlowLogName                  = "name"
	isFlowLogActive                = "active"
	isFlowLogStorageBucket         = "storage_bucket"
	isFlowLogStorageBucketEndPoint = "endpoint"
	isFlowLogTarget                = "target"
	isFlowLogResourceGroup         = "resource_group"
	isFlowLogTargetType            = "resource_type"
	isFlowLogCreatedAt             = "created_at"
	isFlowLogCrn                   = "crn"
	isFlowLogLifecycleState        = "lifecycle_state"
	isFlowLogHref                  = "href"
	isFlowLogAutoDelete            = "auto_delete"
	isFlowLogVpc                   = "vpc"
	isFlowLogTags                  = "tags"
)

func resourceIBMISFlowLog() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISFlowLogCreate,
		Read:     resourceIBMISFlowLogRead,
		Update:   resourceIBMISFlowLogUpdate,
		Delete:   resourceIBMISFlowLogDelete,
		Exists:   resourceIBMISFlowLogExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			func(diff *schema.ResourceDiff, v interface{}) error {
				return resourceTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			isFlowLogName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				Description:  "Flow Log Collector name",
				ValidateFunc: InvokeValidator("ibm_is_flow_log", isFlowLogName),
			},

			isFlowLogStorageBucket: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Cloud Object Storage bucket name where the collected flows will be logged",
			},

			isFlowLogTarget: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target id that the flow log collector is to collect flow logs",
			},

			isFlowLogActive: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether this collector is active",
			},

			isFlowLogResourceGroup: {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "The resource group of flow log",
			},

			isFlowLogCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this flow log collector",
			},

			isFlowLogHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this flow log collector",
			},

			isFlowLogCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time flow log was created",
			},

			isFlowLogVpc: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The VPC this flow log collector is associated with",
			},

			isFlowLogAutoDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this flow log collector will be automatically deleted when the target is deleted",
			},

			isFlowLogLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the flow log collector",
			},

			isFlowLogTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: InvokeValidator("ibm_is_flow_log", "tag")},
				Set:         resourceIBMVPCHash,
				Description: "Tags for the VPC Flow logs",
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

func resourceIBMISFlowLogValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isFlowLogName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
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

	ibmISFlowLogValidator := ResourceValidator{ResourceName: "ibm_is_flow_log", Schema: validateSchema}
	return &ibmISFlowLogValidator
}

func resourceIBMISFlowLogCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	createFlowLogCollectorOptionsModel := &vpcv1.CreateFlowLogCollectorOptions{}
	name := d.Get(isFlowLogName).(string)
	createFlowLogCollectorOptionsModel.Name = &name
	if _, ok := d.GetOk(isFlowLogResourceGroup); ok {
		group := d.Get(isFlowLogResourceGroup).(string)
		resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityByID)
		resourceGroupIdentityModel.ID = &group
		createFlowLogCollectorOptionsModel.ResourceGroup = resourceGroupIdentityModel
	}

	if v, ok := d.GetOkExists(isFlowLogActive); ok {
		active := v.(bool)
		createFlowLogCollectorOptionsModel.Active = &active
	}

	target := d.Get(isFlowLogTarget).(string)
	FlowLogCollectorTargetModel := &vpcv1.FlowLogCollectorTargetPrototype{}
	FlowLogCollectorTargetModel.ID = &target
	createFlowLogCollectorOptionsModel.Target = FlowLogCollectorTargetModel

	bucketname := d.Get(isFlowLogStorageBucket).(string)
	cloudObjectStorageBucketIdentityModel := new(vpcv1.CloudObjectStorageBucketIdentityByName)
	cloudObjectStorageBucketIdentityModel.Name = &bucketname
	createFlowLogCollectorOptionsModel.StorageBucket = cloudObjectStorageBucketIdentityModel

	flowlogCollector, response, err := sess.CreateFlowLogCollector(createFlowLogCollectorOptionsModel)
	if err != nil {
		return fmt.Errorf("Create Flow Log Collector err %s\n%s", err, response)
	}
	d.SetId(*flowlogCollector.ID)

	log.Printf("Flow log collector : %s", *flowlogCollector.ID)

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isFlowLogTags); ok || v != "" {
		oldList, newList := d.GetChange(isFlowLogTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *flowlogCollector.CRN)
		if err != nil {
			log.Printf(
				"Error on create of resource vpc flow log (%s) tags: %s", d.Id(), err)
		}
	}

	return resourceIBMISFlowLogRead(d, meta)
}

func resourceIBMISFlowLogRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	ID := d.Id()

	getOptions := &vpcv1.GetFlowLogCollectorOptions{
		ID: &ID,
	}
	flowlogCollector, response, err := sess.GetFlowLogCollector(getOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Flow Log Collector: %s\n%s", err, response)
	}

	if flowlogCollector.Name != nil {
		d.Set(isFlowLogName, *flowlogCollector.Name)
	}

	if flowlogCollector.Active != nil {
		d.Set(isFlowLogActive, *flowlogCollector.Active)
	}

	if flowlogCollector.CreatedAt != nil {
		d.Set(isFlowLogCreatedAt, flowlogCollector.CreatedAt.String())
	}

	if flowlogCollector.Href != nil {
		d.Set(isFlowLogHref, *flowlogCollector.Href)
	}

	if flowlogCollector.CRN != nil {
		d.Set(isFlowLogCrn, *flowlogCollector.CRN)
	}

	if flowlogCollector.LifecycleState != nil {
		d.Set(isFlowLogLifecycleState, *flowlogCollector.LifecycleState)
	}

	if flowlogCollector.VPC != nil {
		d.Set(isFlowLogVpc, *flowlogCollector.VPC.ID)
	}

	if flowlogCollector.Target != nil {
		targetIntf := flowlogCollector.Target
		target := targetIntf.(*vpcv1.FlowLogCollectorTarget)
		d.Set(isFlowLogTarget, *target.ID)
	}

	if flowlogCollector.StorageBucket != nil {
		bucket := flowlogCollector.StorageBucket
		d.Set(isFlowLogStorageBucket, *bucket.Name)
	}

	tags, err := GetTagsUsingCRN(meta, *flowlogCollector.CRN)
	if err != nil {
		log.Printf(
			"Error on get of resource vpc flow log (%s) tags: %s", d.Id(), err)
	}
	d.Set(isFlowLogTags, tags)
	controller, err := getBaseController(meta)
	if err != nil {
		return err
	}

	d.Set(ResourceControllerURL, controller+"/vpc-ext/network/flowLogs")
	d.Set(ResourceName, *flowlogCollector.Name)
	d.Set(ResourceCRN, *flowlogCollector.CRN)
	d.Set(ResourceStatus, *flowlogCollector.LifecycleState)

	if flowlogCollector.ResourceGroup != nil {
		d.Set(isFlowLogResourceGroup, *flowlogCollector.ResourceGroup.ID)
		d.Set(ResourceGroupName, *flowlogCollector.ResourceGroup.ID)
	}

	return nil
}

func resourceIBMISFlowLogUpdate(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	ID := d.Id()

	getOptions := &vpcv1.GetFlowLogCollectorOptions{
		ID: &ID,
	}
	flowlogCollector, response, err := sess.GetFlowLogCollector(getOptions)
	if err != nil {
		return fmt.Errorf("Error Getting Flow Log Collector: %s\n%s", err, response)
	}

	if d.HasChange(isFlowLogTags) {
		oldList, newList := d.GetChange(isFlowLogTags)
		err = UpdateTagsUsingCRN(oldList, newList, meta, *flowlogCollector.CRN)
		if err != nil {
			log.Printf(
				"Error on update of resource flow log (%s) tags: %s", *flowlogCollector.ID, err)
		}
	}

	if d.HasChange(isFlowLogActive) || d.HasChange(isFlowLogName) {
		active := d.Get(isFlowLogActive).(bool)
		name := d.Get(isFlowLogName).(string)
		updoptions := &vpcv1.UpdateFlowLogCollectorOptions{
			ID: &ID,
		}
		flowLogCollectorPatchModel := &vpcv1.FlowLogCollectorPatch{
			Active: &active,
			Name:   &name,
		}
		flowLogCollectorPatch, err := flowLogCollectorPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("Error calling asPatch for FlowLogCollectorPatch: %s", err)
		}
		updoptions.FlowLogCollectorPatch = flowLogCollectorPatch
		_, response, err = sess.UpdateFlowLogCollector(updoptions)
		if err != nil {
			return fmt.Errorf("Error updating flow log collector:%s\n%s", err, response)
		}
	}

	return resourceIBMISFlowLogRead(d, meta)
}

func resourceIBMISFlowLogDelete(d *schema.ResourceData, meta interface{}) error {

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	ID := d.Id()
	delOptions := &vpcv1.DeleteFlowLogCollectorOptions{
		ID: &ID,
	}
	response, err := sess.DeleteFlowLogCollector(delOptions)

	if err != nil && response.StatusCode != 404 {
		return fmt.Errorf("Error deleting flow log collector:%s\n%s", err, response)
	}

	d.SetId("")
	return nil
}

func resourceIBMISFlowLogExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	ID := d.Id()

	getOptions := &vpcv1.GetFlowLogCollectorOptions{
		ID: &ID,
	}
	_, response, err := sess.GetFlowLogCollector(getOptions)
	if err != nil && response.StatusCode != 404 {
		return false, fmt.Errorf("Error Getting Flow Log Collector : %s\n%s", err, response)
	}
	if response.StatusCode == 404 {
		d.SetId("")
		return false, nil
	}
	return true, nil
}
