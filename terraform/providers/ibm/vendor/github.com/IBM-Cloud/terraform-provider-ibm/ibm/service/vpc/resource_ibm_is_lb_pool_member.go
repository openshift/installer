// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isLBPoolID                       = "pool"
	isLBPoolMemberPort               = "port"
	isLBPoolMemberTargetAddress      = "target_address"
	isLBPoolMemberTargetID           = "target_id"
	isLBPoolMemberWeight             = "weight"
	isLBPoolMemberProvisioningStatus = "provisioning_status"
	isLBPoolMemberHealth             = "health"
	isLBPoolMemberHref               = "href"
	isLBPoolMemberDeletePending      = "delete_pending"
	isLBPoolMemberDeleted            = "done"
	isLBPoolMemberActive             = "active"
	isLBPoolUpdating                 = "updating"
)

func ResourceIBMISLBPoolMember() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISLBPoolMemberCreate,
		Read:     resourceIBMISLBPoolMemberRead,
		Update:   resourceIBMISLBPoolMemberUpdate,
		Delete:   resourceIBMISLBPoolMemberDelete,
		Exists:   resourceIBMISLBPoolMemberExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isLBPoolID: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, o, n string, d *schema.ResourceData) bool {
					if o == "" {
						return false
					}
					// if state file entry and tf file entry matches
					if strings.Compare(n, o) == 0 {
						return true
					}

					if strings.Contains(n, "/") {
						new := strings.Split(n, "/")
						if strings.Compare(new[1], o) == 0 {
							return true
						}
					}

					return false
				},
				Description: "Loadblancer Poold ID",
			},

			isLBID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Load balancer ID",
			},

			isLBPoolMemberPort: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Load Balancer Pool port",
			},

			isLBPoolMemberTargetAddress: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isLBPoolMemberTargetAddress, isLBPoolMemberTargetID},
				Description:  "Load balancer pool member target address",
			},

			isLBPoolMemberTargetID: {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{isLBPoolMemberTargetAddress, isLBPoolMemberTargetID},
				Description:  "Load balancer pool member target id",
			},

			isLBPoolMemberWeight: {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_lb_pool_member", isLBPoolMemberWeight),
				Description:  "Load balcner pool member weight",
			},

			isLBPoolMemberProvisioningStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer Pool member provisioning status",
			},

			isLBPoolMemberHealth: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LB Pool member health",
			},

			isLBPoolMemberHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "LB pool member Href value",
			},

			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the LB resource",
			},
		},
	}
}

func ResourceIBMISLBPoolMemberValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isLBPoolMemberWeight,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "0",
			MaxValue:                   "100"})

	ibmISLBResourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_lb_pool_member", Schema: validateSchema}
	return &ibmISLBResourceValidator
}

func resourceIBMISLBPoolMemberCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("[DEBUG] LB Pool create")
	lbPoolID, err := getPoolId(d.Get(isLBPoolID).(string))
	if err != nil {
		return err
	}

	lbID := d.Get(isLBID).(string)
	port := d.Get(isLBPoolMemberPort).(int)
	port64 := int64(port)

	var weight int64

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	err = lbpMemberCreate(d, meta, lbID, lbPoolID, port64, weight)
	if err != nil {
		return err
	}

	return resourceIBMISLBPoolMemberRead(d, meta)
}

func lbpMemberCreate(d *schema.ResourceData, meta interface{}, lbID, lbPoolID string, port, weight int64) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	options := &vpcv1.CreateLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		Port:           &port,
	}

	if _, ok := d.GetOk(isLBPoolMemberTargetAddress); ok {
		targetAddress := d.Get(isLBPoolMemberTargetAddress).(string)
		target := &vpcv1.LoadBalancerPoolMemberTargetPrototype{
			Address: &targetAddress,
		}
		options.Target = target
	} else {
		targetID := d.Get(isLBPoolMemberTargetID).(string)
		target := &vpcv1.LoadBalancerPoolMemberTargetPrototype{
			ID: &targetID,
		}
		options.Target = target
	}
	if w, ok := d.GetOkExists(isLBPoolMemberWeight); ok {
		weight = int64(w.(int))
		options.Weight = &weight
	}

	lbPoolMember, response, err := sess.CreateLoadBalancerPoolMember(options)
	if err != nil {
		return fmt.Errorf("[DEBUG] lbpool member create err: %s\n%s", err, response)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", lbID, lbPoolID, *lbPoolMember.ID))
	log.Printf("[INFO] lbpool member : %s", *lbPoolMember.ID)

	_, err = isWaitForLBPoolMemberAvailable(sess, lbID, lbPoolID, *lbPoolMember.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	return nil
}

func isWaitForLBPoolMemberAvailable(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for load balancer pool member(%s) to be available.", lbPoolMemID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"create_pending", "update_pending", "maintenance_pending"},
		Target:     []string{isLBPoolMemberActive, ""},
		Refresh:    isLBPoolMemberRefreshFunc(lbc, lbID, lbPoolID, lbPoolMemID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isLBPoolMemberRefreshFunc(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
			LoadBalancerID: &lbID,
			PoolID:         &lbPoolID,
			ID:             &lbPoolMemID,
		}
		lbPoolMem, response, err := lbc.GetLoadBalancerPoolMember(getlbpmoptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error Getting Load Balancer Pool Member: %s\n%s", err, response)
		}

		if *lbPoolMem.ProvisioningStatus == isLBPoolMemberActive {
			return lbPoolMem, *lbPoolMem.ProvisioningStatus, nil
		}

		return lbPoolMem, *lbPoolMem.ProvisioningStatus, nil
	}
}

func resourceIBMISLBPoolMemberRead(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	if len(parts) < 3 {
		return fmt.Errorf(
			"The id should contain loadbalancer Id, loadbalancer pool Id and loadbalancer poolmemebr Id")
	}

	lbID := parts[0]
	lbPoolID := parts[1]
	lbPoolMemID := parts[2]

	err = lbpmemberGet(d, meta, lbID, lbPoolID, lbPoolMemID)
	if err != nil {
		return err
	}

	return nil
}

func lbpmemberGet(d *schema.ResourceData, meta interface{}, lbID, lbPoolID, lbPoolMemID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		ID:             &lbPoolMemID,
	}
	lbPoolMem, response, err := sess.GetLoadBalancerPoolMember(getlbpmoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Load Balancer Pool Member: %s\n%s", err, response)
	}
	d.Set(isLBPoolID, lbPoolID)
	d.Set(isLBID, lbID)
	d.Set(isLBPoolMemberPort, *lbPoolMem.Port)

	target := lbPoolMem.Target.(*vpcv1.LoadBalancerPoolMemberTarget)
	if target.Address != nil {
		d.Set(isLBPoolMemberTargetAddress, *target.Address)
	}
	if target.ID != nil {
		d.Set(isLBPoolMemberTargetID, *target.ID)
	}
	d.Set(isLBPoolMemberWeight, *lbPoolMem.Weight)
	d.Set(isLBPoolMemberProvisioningStatus, *lbPoolMem.ProvisioningStatus)
	d.Set(isLBPoolMemberHealth, *lbPoolMem.Health)
	d.Set(isLBPoolMemberHref, *lbPoolMem.Href)
	getLoadBalancerOptions := &vpcv1.GetLoadBalancerOptions{
		ID: &lbID,
	}
	lb, response, err := sess.GetLoadBalancer(getLoadBalancerOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting Load Balancer : %s\n%s", err, response)
	}
	d.Set(flex.RelatedCRN, *lb.CRN)
	return nil
}

func resourceIBMISLBPoolMemberUpdate(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]
	lbPoolMemID := parts[2]

	err = lbpmemberUpdate(d, meta, lbID, lbPoolID, lbPoolMemID)
	if err != nil {
		return err
	}

	return resourceIBMISLBPoolMemberRead(d, meta)
}

func lbpmemberUpdate(d *schema.ResourceData, meta interface{}, lbID, lbPoolID, lbPoolMemID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	if d.HasChange(isLBPoolMemberTargetID) || d.HasChange(isLBPoolMemberTargetAddress) || d.HasChange(isLBPoolMemberPort) || d.HasChange(isLBPoolMemberWeight) {

		port := int64(d.Get(isLBPoolMemberPort).(int))
		weight := int64(d.Get(isLBPoolMemberWeight).(int))

		isLBKey := "load_balancer_key_" + lbID
		conns.IbmMutexKV.Lock(isLBKey)
		defer conns.IbmMutexKV.Unlock(isLBKey)

		_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		}

		_, err = isWaitForLBPoolMemberAvailable(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}

		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}

		updatelbpmoptions := &vpcv1.UpdateLoadBalancerPoolMemberOptions{
			LoadBalancerID: &lbID,
			PoolID:         &lbPoolID,
			ID:             &lbPoolMemID,
		}

		loadBalancerPoolMemberPatchModel := &vpcv1.LoadBalancerPoolMemberPatch{
			Port:   &port,
			Weight: &weight,
		}

		if _, ok := d.GetOk(isLBPoolMemberTargetAddress); ok {
			targetAddress := d.Get(isLBPoolMemberTargetAddress).(string)
			target := &vpcv1.LoadBalancerPoolMemberTargetPrototype{
				Address: &targetAddress,
			}
			loadBalancerPoolMemberPatchModel.Target = target
		} else {
			targetID := d.Get(isLBPoolMemberTargetID).(string)
			target := &vpcv1.LoadBalancerPoolMemberTargetPrototype{
				ID: &targetID,
			}
			loadBalancerPoolMemberPatchModel.Target = target
		}

		loadBalancerPoolMemberPatch, err := loadBalancerPoolMemberPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for LoadBalancerPoolMemberPatch: %s", err)
		}
		updatelbpmoptions.LoadBalancerPoolMemberPatch = loadBalancerPoolMemberPatch

		_, response, err := sess.UpdateLoadBalancerPoolMember(updatelbpmoptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Updating Load Balancer Pool Member: %s\n%s", err, response)
		}
		_, err = isWaitForLBPoolMemberAvailable(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}

		_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
		}

		_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf(
				"Error checking for load balancer (%s) is active: %s", lbID, err)
		}
	}
	return nil
}

func resourceIBMISLBPoolMemberDelete(d *schema.ResourceData, meta interface{}) error {

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	lbID := parts[0]
	lbPoolID := parts[1]
	lbPoolMemID := parts[2]

	isLBKey := "load_balancer_key_" + lbID
	conns.IbmMutexKV.Lock(isLBKey)
	defer conns.IbmMutexKV.Unlock(isLBKey)

	err = lbpmemberDelete(d, meta, lbID, lbPoolID, lbPoolMemID)
	if err != nil {
		return err
	}

	return nil
}

func lbpmemberDelete(d *schema.ResourceData, meta interface{}, lbID, lbPoolID, lbPoolMemID string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		ID:             &lbPoolMemID,
	}
	_, response, err := sess.GetLoadBalancerPoolMember(getlbpmoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error Getting Load Balancer Pool Member: %s\n%s", err, response)
	}
	_, err = isWaitForLBPoolMemberAvailable(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	dellbpmoptions := &vpcv1.DeleteLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		ID:             &lbPoolMemID,
	}
	response, err = sess.DeleteLoadBalancerPoolMember(dellbpmoptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Load Balancer Pool Member: %s\n%s", err, response)
	}

	_, err = isWaitForLBPoolMemberDeleted(sess, lbID, lbPoolID, lbPoolMemID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	_, err = isWaitForLBPoolActive(sess, lbID, lbPoolID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer pool (%s) is active: %s", lbPoolID, err)
	}

	_, err = isWaitForLBAvailable(sess, lbID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("[ERROR] Error checking for load balancer (%s) is active: %s", lbID, err)
	}

	d.SetId("")
	return nil
}

func isWaitForLBPoolMemberDeleted(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for  (%s) to be deleted.", lbPoolMemID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isLBPoolMemberDeletePending},
		Target:     []string{isLBPoolMemberDeleted, ""},
		Refresh:    isDeleteLBPoolMemberRefreshFunc(lbc, lbID, lbPoolID, lbPoolMemID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isDeleteLBPoolMemberRefreshFunc(lbc *vpcv1.VpcV1, lbID, lbPoolID, lbPoolMemID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
			LoadBalancerID: &lbID,
			PoolID:         &lbPoolID,
			ID:             &lbPoolMemID,
		}
		lbPoolMem, response, err := lbc.GetLoadBalancerPoolMember(getlbpmoptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return lbPoolMem, isLBPoolMemberDeleted, nil
			}
			return nil, "", fmt.Errorf("[ERROR] Error Deleting Load balancer pool member: %s\n%s", err, response)
		}
		return lbPoolMem, isLBPoolMemberDeletePending, nil
	}
}

func resourceIBMISLBPoolMemberExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 3 {
		return false, fmt.Errorf(
			"The id should contain loadbalancer Id, loadbalancer pool Id and loadbalancer poolmemebr Id")
	}

	lbID := parts[0]
	lbPoolID := parts[1]
	lbPoolMemID := parts[2]

	exists, err := lbpmemberExists(d, meta, lbID, lbPoolID, lbPoolMemID)
	return exists, err

}

func lbpmemberExists(d *schema.ResourceData, meta interface{}, lbID, lbPoolID, lbPoolMemID string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}

	getlbpmoptions := &vpcv1.GetLoadBalancerPoolMemberOptions{
		LoadBalancerID: &lbID,
		PoolID:         &lbPoolID,
		ID:             &lbPoolMemID,
	}
	_, response, err := sess.GetLoadBalancerPoolMember(getlbpmoptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting Load balancer pool member: %s\n%s", err, response)
	}
	return true, nil
}

func getPoolId(id string) (string, error) {
	if strings.Contains(id, "/") {
		parts, err := flex.IdParts(id)
		if err != nil {
			return "", err
		}

		return parts[1], nil
	} else {
		return id, nil
	}
}
