// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMISInstanceGroupManagerPolicy() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISInstanceGroupManagerPolicyCreate,
		Read:     resourceIBMISInstanceGroupManagerPolicyRead,
		Update:   resourceIBMISInstanceGroupManagerPolicyUpdate,
		Delete:   resourceIBMISInstanceGroupManagerPolicyDelete,
		Exists:   resourceIBMISInstanceGroupManagerPolicyExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_manager_policy", "name"),
				Description:  "instance group manager policy name",
			},

			"instance_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance group ID",
			},

			"instance_group_manager": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance group manager ID",
			},

			"metric_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_manager_policy", "metric_type"),
				Description:  "The type of metric to be evaluated",
			},

			"metric_value": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The metric value to be evaluated",
			},

			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_group_manager_policy", "policy_type"),
				Description:  "The type of Policy for the Instance Group",
			},

			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Policy ID",
			},
		},
	}
}

func resourceIBMISInstanceGroupManagerPolicyValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	metricTypes := "cpu,memory,network_in,network_out"
	policyType := "target"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "metric_type",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              metricTypes})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "policy_type",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              policyType})

	ibmISInstanceGroupManagerPolicyResourceValidator := ResourceValidator{ResourceName: "ibm_is_instance_group_manager_policy", Schema: validateSchema}
	return &ibmISInstanceGroupManagerPolicyResourceValidator
}

func resourceIBMISInstanceGroupManagerPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	instanceGroupID := d.Get("instance_group").(string)
	instanceGroupManagerID := d.Get("instance_group_manager").(string)

	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceGroupManagerPolicyPrototype := vpcv1.InstanceGroupManagerPolicyPrototype{}

	name := d.Get("name").(string)
	metricType := d.Get("metric_type").(string)
	metricValue := int64(d.Get("metric_value").(int))
	policyType := d.Get("policy_type").(string)

	instanceGroupManagerPolicyPrototype.Name = &name
	instanceGroupManagerPolicyPrototype.MetricType = &metricType
	instanceGroupManagerPolicyPrototype.MetricValue = &metricValue
	instanceGroupManagerPolicyPrototype.PolicyType = &policyType

	createInstanceGroupManagerPolicyOptions := vpcv1.CreateInstanceGroupManagerPolicyOptions{
		InstanceGroupID:                     &instanceGroupID,
		InstanceGroupManagerID:              &instanceGroupManagerID,
		InstanceGroupManagerPolicyPrototype: &instanceGroupManagerPolicyPrototype,
	}

	isInsGrpKey := "Instance_Group_Key_" + instanceGroupID
	ibmMutexKV.Lock(isInsGrpKey)
	defer ibmMutexKV.Unlock(isInsGrpKey)

	_, healthError := waitForHealthyInstanceGroup(instanceGroupID, meta, d.Timeout(schema.TimeoutCreate))
	if healthError != nil {
		return healthError
	}

	data, response, err := sess.CreateInstanceGroupManagerPolicy(&createInstanceGroupManagerPolicyOptions)
	if err != nil || data == nil {
		return fmt.Errorf("Error Creating InstanceGroup Manager Policy: %s\n%s", err, response)
	}
	instanceGroupManagerPolicy := data.(*vpcv1.InstanceGroupManagerPolicy)

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceGroupID, instanceGroupManagerID, *instanceGroupManagerPolicy.ID))

	return resourceIBMISInstanceGroupManagerPolicyRead(d, meta)

}

func resourceIBMISInstanceGroupManagerPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	var changed bool
	updateInstanceGroupManagerPolicyOptions := vpcv1.UpdateInstanceGroupManagerPolicyOptions{}
	instanceGroupManagerPolicyPatchModel := &vpcv1.InstanceGroupManagerPolicyPatch{}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		instanceGroupManagerPolicyPatchModel.Name = &name
		changed = true
	}

	if d.HasChange("metric_type") {
		metricType := d.Get("metric_type").(string)
		instanceGroupManagerPolicyPatchModel.MetricType = &metricType
		changed = true
	}

	if d.HasChange("metric_value") {
		metricValue := int64(d.Get("metric_value").(int))
		instanceGroupManagerPolicyPatchModel.MetricValue = &metricValue
		changed = true
	}

	if changed {
		parts, err := idParts(d.Id())
		if err != nil {
			return err
		}
		instanceGroupID := parts[0]
		instanceGroupManagerID := parts[1]
		instanceGroupManagerPolicyID := parts[2]

		updateInstanceGroupManagerPolicyOptions.ID = &instanceGroupManagerPolicyID
		updateInstanceGroupManagerPolicyOptions.InstanceGroupID = &instanceGroupID
		updateInstanceGroupManagerPolicyOptions.InstanceGroupManagerID = &instanceGroupManagerID

		isInsGrpKey := "Instance_Group_Key_" + instanceGroupID
		ibmMutexKV.Lock(isInsGrpKey)
		defer ibmMutexKV.Unlock(isInsGrpKey)

		_, healthError := waitForHealthyInstanceGroup(instanceGroupID, meta, d.Timeout(schema.TimeoutUpdate))
		if healthError != nil {
			return healthError
		}

		_, response, err := sess.UpdateInstanceGroupManagerPolicy(&updateInstanceGroupManagerPolicyOptions)
		if err != nil {
			return fmt.Errorf("Error Updating InstanceGroup Manager Policy: %s\n%s", err, response)
		}
	}
	return resourceIBMISInstanceGroupManagerPolicyRead(d, meta)
}

func resourceIBMISInstanceGroupManagerPolicyRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	instanceGroupID := parts[0]
	instanceGroupManagerID := parts[1]
	instanceGroupManagerPolicyID := parts[2]

	getInstanceGroupManagerPolicyOptions := vpcv1.GetInstanceGroupManagerPolicyOptions{
		ID:                     &instanceGroupManagerPolicyID,
		InstanceGroupID:        &instanceGroupID,
		InstanceGroupManagerID: &instanceGroupManagerID,
	}
	data, response, err := sess.GetInstanceGroupManagerPolicy(&getInstanceGroupManagerPolicyOptions)
	if err != nil || data == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Getting InstanceGroup Manager Policy: %s\n%s", err, response)
	}
	instanceGroupManagerPolicy := data.(*vpcv1.InstanceGroupManagerPolicy)
	d.Set("name", *instanceGroupManagerPolicy.Name)
	d.Set("metric_value", instanceGroupManagerPolicy.MetricValue)
	d.Set("metric_type", instanceGroupManagerPolicy.MetricType)
	d.Set("policy_type", instanceGroupManagerPolicy.PolicyType)
	d.Set("policy_id", instanceGroupManagerPolicyID)
	d.Set("instance_group", instanceGroupID)
	d.Set("instance_group_manager", instanceGroupManagerID)

	return nil
}

func resourceIBMISInstanceGroupManagerPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	instanceGroupID := parts[0]
	instanceGroupManagerID := parts[1]
	instanceGroupManagerPolicyID := parts[2]

	deleteInstanceGroupManagerPolicyOptions := vpcv1.DeleteInstanceGroupManagerPolicyOptions{
		ID:                     &instanceGroupManagerPolicyID,
		InstanceGroupManagerID: &instanceGroupManagerID,
		InstanceGroupID:        &instanceGroupID,
	}

	isInsGrpKey := "Instance_Group_Key_" + instanceGroupID
	ibmMutexKV.Lock(isInsGrpKey)
	defer ibmMutexKV.Unlock(isInsGrpKey)

	_, healthError := waitForHealthyInstanceGroup(instanceGroupID, meta, d.Timeout(schema.TimeoutDelete))
	if healthError != nil {
		return healthError
	}

	response, err := sess.DeleteInstanceGroupManagerPolicy(&deleteInstanceGroupManagerPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Deleting the InstanceGroup Manager Policy: %s\n%s", err, response)
	}
	return nil
}

func resourceIBMISInstanceGroupManagerPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	instanceGroupID := parts[0]
	instanceGroupManagerID := parts[1]
	instanceGroupManagerPolicyID := parts[2]

	getInstanceGroupManagerPolicyOptions := vpcv1.GetInstanceGroupManagerPolicyOptions{
		ID:                     &instanceGroupManagerPolicyID,
		InstanceGroupManagerID: &instanceGroupManagerID,
		InstanceGroupID:        &instanceGroupID,
	}

	_, response, err := sess.GetInstanceGroupManagerPolicy(&getInstanceGroupManagerPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error Getting InstanceGroup Manager Policy: %s\n%s", err, response)
	}
	return true, nil
}
