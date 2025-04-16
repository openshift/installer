// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceNetworkInterfaceFloatingIpAvailable = "available"
	isInstanceNetworkInterfaceFloatingIpDeleting  = "deleting"
	isInstanceNetworkInterfaceFloatingIpPending   = "pending"
	isInstanceNetworkInterfaceFloatingIpDeleted   = "deleted"
	isInstanceNetworkInterfaceFloatingIpFailed    = "failed"
	isInstanceNetworkInterface                    = "network_interface"
	isInstanceNetworkInterfaceFloatingIPID        = "floating_ip"
)

func ResourceIBMIsInstanceNetworkInterfaceFloatingIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISInstanceNetworkInterfaceFloatingIpCreate,
		ReadContext:   resourceIBMISInstanceNetworkInterfaceFloatingIpRead,
		UpdateContext: resourceIBMISInstanceNetworkInterfaceFloatingIpUpdate,
		DeleteContext: resourceIBMISInstanceNetworkInterfaceFloatingIpDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			isInstanceNetworkInterface: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance network interface identifier",
			},
			isInstanceNetworkInterfaceFloatingIPID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The floating ip identifier of the network interface associated with the Instance",
			},
			floatingIPName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the floating IP",
			},

			floatingIPAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP address",
			},

			floatingIPStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP status",
			},

			floatingIPZone: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Zone name",
			},

			floatingIPTarget: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target info",
			},

			floatingIPCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Floating IP crn",
			},
		},
	}
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := d.Get(isInstanceID).(string)
	instanceNicId := ""
	nicId := d.Get(isInstanceNetworkInterface).(string)
	if strings.Contains(nicId, "/") {
		_, instanceNicId, err = ParseNICTerraformID(nicId)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		instanceNicId = nicId
	}

	instanceNicFipId := d.Get(isInstanceNetworkInterfaceFloatingIPID).(string)

	options := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &instanceNicId,
		ID:                 &instanceNicFipId,
	}

	fip, response, err := sess.AddInstanceNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil || fip == nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Create Instance (%s) network interface (%s) floating ip (%s) err %s\n%s", instanceId, instanceNicId, instanceNicFipId, err, response))
	}
	d.SetId(MakeTerraformNICFipID(instanceId, instanceNicId, *fip.ID))
	err = instanceNICFipGet(d, fip, instanceId, instanceNicId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceId, nicID, fipId, err := ParseNICFipTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &nicID,
		ID:                 &fipId,
	}

	fip, response, err := sess.GetInstanceNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil || fip == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Instance (%s) network interface (%s): %s\n%s", instanceId, nicID, err, response))
	}
	err = instanceNICFipGet(d, fip, instanceId, nicID)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func instanceNICFipGet(d *schema.ResourceData, fip *vpcv1.FloatingIP, instanceId, nicId string) error {

	d.SetId(MakeTerraformNICFipID(instanceId, nicId, *fip.ID))
	d.Set(floatingIPName, *fip.Name)
	d.Set(floatingIPAddress, *fip.Address)
	d.Set(floatingIPStatus, fip.Status)
	d.Set(floatingIPZone, *fip.Zone.Name)

	d.Set(floatingIPCRN, *fip.CRN)

	target, ok := fip.Target.(*vpcv1.FloatingIPTarget)
	if ok {
		d.Set(floatingIPTarget, target.ID)
	}

	return nil
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange(isInstanceNetworkInterfaceFloatingIPID) {
		instanceId, nicId, _, err := ParseNICFipTerraformID(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		sess, err := vpcClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		floatingIpId := ""
		if fipOk, ok := d.GetOk(isInstanceNetworkInterfaceFloatingIPID); ok {
			floatingIpId = fipOk.(string)
		}
		options := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{
			InstanceID:         &instanceId,
			NetworkInterfaceID: &nicId,
			ID:                 &floatingIpId,
		}

		fip, response, err := sess.AddInstanceNetworkInterfaceFloatingIPWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating Instance: %s\n%s", err, response))
		}
		d.SetId(MakeTerraformNICFipID(instanceId, nicId, *fip.ID))
		return diag.FromErr(instanceNICFipGet(d, fip, instanceId, nicId))
	}
	return nil
}

func resourceIBMISInstanceNetworkInterfaceFloatingIpDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	instanceId, nicId, fipId, err := ParseNICFipTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = instanceNetworkInterfaceFipDelete(context, d, meta, instanceId, nicId, fipId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func instanceNetworkInterfaceFipDelete(context context.Context, d *schema.ResourceData, meta interface{}, instanceId, nicId, fipId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getBmsNicFipOptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &nicId,
		ID:                 &fipId,
	}
	fip, response, err := sess.GetInstanceNetworkInterfaceFloatingIPWithContext(context, getBmsNicFipOptions)
	if err != nil || fip == nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Instance (%s) network interface(%s) Floating Ip(%s) : %s\n%s", instanceId, nicId, fipId, err, response)
	}

	options := &vpcv1.RemoveInstanceNetworkInterfaceFloatingIPOptions{
		InstanceID:         &instanceId,
		NetworkInterfaceID: &nicId,
		ID:                 &fipId,
	}
	response, err = sess.RemoveInstanceNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Instance (%s) network interface (%s) Floating Ip(%s) : %s\n%s", instanceId, nicId, fipId, err, response)
	}
	_, err = isWaitForInstanceNetworkInterfaceFloatingIpDeleted(sess, instanceId, nicId, fipId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForInstanceNetworkInterfaceFloatingIpDeleted(instanceC *vpcv1.VpcV1, instanceId, nicId, fipId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for (%s) / (%s) / (%s) to be deleted.", instanceId, nicId, fipId)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isInstanceNetworkInterfaceFloatingIpAvailable, isInstanceNetworkInterfaceFloatingIpDeleting, isInstanceNetworkInterfaceFloatingIpPending},
		Target:     []string{isInstanceNetworkInterfaceFloatingIpDeleted, isInstanceNetworkInterfaceFloatingIpFailed, ""},
		Refresh:    isInstanceNetworkInterfaceFloatingIpDeleteRefreshFunc(instanceC, instanceId, nicId, fipId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isInstanceNetworkInterfaceFloatingIpDeleteRefreshFunc(instanceC *vpcv1.VpcV1, instanceId, nicId, fipId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getBmsNicFloatingIpOptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
			InstanceID:         &instanceId,
			NetworkInterfaceID: &nicId,
			ID:                 &fipId,
		}
		fip, response, err := instanceC.GetInstanceNetworkInterfaceFloatingIP(getBmsNicFloatingIpOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return fip, isInstanceNetworkInterfaceFloatingIpDeleted, nil
			}
			return fip, isInstanceNetworkInterfaceFloatingIpFailed, fmt.Errorf("[ERROR] Error getting Instance(%s) Network Interface (%s) FloatingIp(%s) : %s\n%s", instanceId, nicId, fipId, err, response)
		}
		return fip, isInstanceNetworkInterfaceFloatingIpDeleting, err
	}
}

func isWaitForInstanceNetworkInterfaceFloatingIpAvailable(client *vpcv1.VpcV1, instanceId, nicId, fipId string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for Instance (%s) Network Interface (%s) to be available.", instanceId, nicId)
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isInstanceNetworkInterfaceFloatingIpPending},
		Target:     []string{isInstanceNetworkInterfaceFloatingIpAvailable, isInstanceNetworkInterfaceFloatingIpFailed},
		Refresh:    isInstanceNetworkInterfaceFloatingIpRefreshFunc(client, instanceId, nicId, fipId, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isInstanceNetworkInterfaceFloatingIpRefreshFunc(client *vpcv1.VpcV1, instanceId, nicId, fipId string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getBmsNicFloatingIpOptions := &vpcv1.GetInstanceNetworkInterfaceFloatingIPOptions{
			InstanceID:         &instanceId,
			NetworkInterfaceID: &nicId,
			ID:                 &fipId,
		}
		fip, response, err := client.GetInstanceNetworkInterfaceFloatingIP(getBmsNicFloatingIpOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting Instance (%s) Network Interface (%s) FloatingIp(%s) : %s\n%s", instanceId, nicId, fipId, err, response)
		}
		status := ""

		status = *fip.Status
		d.Set(floatingIPStatus, *fip.Status)

		select {
		case data := <-communicator:
			return nil, "", data.(error)
		default:
			fmt.Println("no message sent")
		}

		if status == "available" || status == "failed" {
			close(communicator)
			return fip, status, nil

		}

		return fip, "pending", nil
	}
}
