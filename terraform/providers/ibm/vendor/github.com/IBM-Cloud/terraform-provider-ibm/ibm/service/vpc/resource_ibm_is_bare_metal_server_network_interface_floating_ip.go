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
	isBareMetalServerNetworkInterfaceFloatingIpAvailable  = "available"
	isBareMetalServerNetworkInterfaceFloatingIpDeleting   = "deleting"
	isBareMetalServerNetworkInterfaceFloatingIpPending    = "pending"
	isBareMetalServerNetworkInterfacePCIFloatingIpPending = "pci_pending"
	isBareMetalServerNetworkInterfaceFloatingIpDeleted    = "deleted"
	isBareMetalServerNetworkInterfaceFloatingIpFailed     = "failed"
	isBareMetalServerFloatingIpHardStop                   = "hard_stop"
)

func ResourceIBMIsBareMetalServerNetworkInterfaceFloatingIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISBareMetalServerNetworkInterfaceFloatingIpCreate,
		ReadContext:   resourceIBMISBareMetalServerNetworkInterfaceFloatingIpRead,
		UpdateContext: resourceIBMISBareMetalServerNetworkInterfaceFloatingIpUpdate,
		DeleteContext: resourceIBMISBareMetalServerNetworkInterfaceFloatingIpDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bare metal server identifier",
			},
			isBareMetalServerNetworkInterface: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bare metal server network interface identifier",
			},
			isBareMetalServerNetworkInterfaceFloatingIPID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The floating ip identifier of the network interface associated with the bare metal server",
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

func resourceIBMISBareMetalServerNetworkInterfaceFloatingIpCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bareMetalServerId := ""
	if bmsId, ok := d.GetOk(isBareMetalServerID); ok {
		bareMetalServerId = bmsId.(string)
	}
	bareMetalServerNicId := ""
	if nicId, ok := d.GetOk(isBareMetalServerNetworkInterface); ok {
		if strings.Contains(nicId.(string), "/") {
			_, bareMetalServerNicId, err = ParseNICTerraformID(nicId.(string))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			bareMetalServerNicId = nicId.(string)
		}

	}
	bareMetalServerNicFipId := ""
	if fipId, ok := d.GetOk(isBareMetalServerNetworkInterfaceFloatingIPID); ok {
		bareMetalServerNicFipId = fipId.(string)
	}

	options := &vpcv1.AddBareMetalServerNetworkInterfaceFloatingIPOptions{
		BareMetalServerID:  &bareMetalServerId,
		NetworkInterfaceID: &bareMetalServerNicId,
		ID:                 &bareMetalServerNicFipId,
	}

	fip, response, err := sess.AddBareMetalServerNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil || fip == nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] Create bare metal server (%s) network interface (%s) floating ip (%s) err %s\n%s", bareMetalServerId, bareMetalServerNicId, bareMetalServerNicFipId, err, response))
	}
	d.SetId(MakeTerraformNICFipID(bareMetalServerId, bareMetalServerNicId, *fip.ID))
	err = bareMetalServerNICFipGet(d, fip, bareMetalServerId, bareMetalServerNicId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIBMISBareMetalServerNetworkInterfaceFloatingIpRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerId, nicID, fipId, err := ParseNICFipTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetBareMetalServerNetworkInterfaceFloatingIPOptions{
		BareMetalServerID:  &bareMetalServerId,
		NetworkInterfaceID: &nicID,
		ID:                 &fipId,
	}

	fip, response, err := sess.GetBareMetalServerNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil || fip == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) network interface (%s): %s\n%s", bareMetalServerId, nicID, err, response))
	}
	err = bareMetalServerNICFipGet(d, fip, bareMetalServerId, nicID)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func bareMetalServerNICFipGet(d *schema.ResourceData, fip *vpcv1.FloatingIP, bareMetalServerId, nicId string) error {

	d.SetId(MakeTerraformNICFipID(bareMetalServerId, nicId, *fip.ID))
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

func resourceIBMISBareMetalServerNetworkInterfaceFloatingIpUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange(isBareMetalServerNetworkInterfaceFloatingIPID) {
		bareMetalServerId, nicId, _, err := ParseNICFipTerraformID(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		sess, err := vpcClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		floatingIpId := ""
		if fipOk, ok := d.GetOk(isBareMetalServerNetworkInterfaceFloatingIPID); ok {
			floatingIpId = fipOk.(string)
		}
		options := &vpcv1.AddBareMetalServerNetworkInterfaceFloatingIPOptions{
			BareMetalServerID:  &bareMetalServerId,
			NetworkInterfaceID: &nicId,
			ID:                 &floatingIpId,
		}

		fip, response, err := sess.AddBareMetalServerNetworkInterfaceFloatingIPWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating Bare Metal Server: %s\n%s", err, response))
		}
		d.SetId(MakeTerraformNICFipID(bareMetalServerId, nicId, *fip.ID))
		return diag.FromErr(bareMetalServerNICFipGet(d, fip, bareMetalServerId, nicId))
	}
	return nil
}

func resourceIBMISBareMetalServerNetworkInterfaceFloatingIpDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerId, nicId, fipId, err := ParseNICFipTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = bareMetalServerNetworkInterfaceFipDelete(context, d, meta, bareMetalServerId, nicId, fipId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func bareMetalServerNetworkInterfaceFipDelete(context context.Context, d *schema.ResourceData, meta interface{}, bareMetalServerId, nicId, fipId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getBmsNicFipOptions := &vpcv1.GetBareMetalServerNetworkInterfaceFloatingIPOptions{
		BareMetalServerID:  &bareMetalServerId,
		NetworkInterfaceID: &nicId,
		ID:                 &fipId,
	}
	fip, response, err := sess.GetBareMetalServerNetworkInterfaceFloatingIPWithContext(context, getBmsNicFipOptions)
	if err != nil || fip == nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) network interface(%s) Floating Ip(%s) : %s\n%s", bareMetalServerId, nicId, fipId, err, response)
	}

	options := &vpcv1.RemoveBareMetalServerNetworkInterfaceFloatingIPOptions{
		BareMetalServerID:  &bareMetalServerId,
		NetworkInterfaceID: &nicId,
		ID:                 &fipId,
	}
	response, err = sess.RemoveBareMetalServerNetworkInterfaceFloatingIPWithContext(context, options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Bare Metal Server (%s) network interface (%s) Floating Ip(%s) : %s\n%s", bareMetalServerId, nicId, fipId, err, response)
	}
	_, err = isWaitForBareMetalServerNetworkInterfaceFloatingIpDeleted(sess, bareMetalServerId, nicId, fipId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForBareMetalServerNetworkInterfaceFloatingIpDeleted(bmsC *vpcv1.VpcV1, bareMetalServerId, nicId, fipId string, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for (%s) / (%s) / (%s) to be deleted.", bareMetalServerId, nicId, fipId)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isBareMetalServerNetworkInterfaceFloatingIpAvailable, isBareMetalServerNetworkInterfaceFloatingIpDeleting, isBareMetalServerNetworkInterfaceFloatingIpPending},
		Target:     []string{isBareMetalServerNetworkInterfaceFloatingIpDeleted, isBareMetalServerNetworkInterfaceFailed, ""},
		Refresh:    isBareMetalServerNetworkInterfaceFloatingIpDeleteRefreshFunc(bmsC, bareMetalServerId, nicId, fipId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isBareMetalServerNetworkInterfaceFloatingIpDeleteRefreshFunc(bmsC *vpcv1.VpcV1, bareMetalServerId, nicId, fipId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getBmsNicFloatingIpOptions := &vpcv1.GetBareMetalServerNetworkInterfaceFloatingIPOptions{
			BareMetalServerID:  &bareMetalServerId,
			NetworkInterfaceID: &nicId,
			ID:                 &fipId,
		}
		fip, response, err := bmsC.GetBareMetalServerNetworkInterfaceFloatingIP(getBmsNicFloatingIpOptions)

		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return fip, isBareMetalServerNetworkInterfaceFloatingIpDeleted, nil
			}
			return fip, isBareMetalServerNetworkInterfaceFloatingIpFailed, fmt.Errorf("[ERROR] Error getting Bare Metal Server(%s) Network Interface (%s) FloatingIp(%s) : %s\n%s", bareMetalServerId, nicId, fipId, err, response)
		}
		return fip, isBareMetalServerNetworkInterfaceFloatingIpDeleting, err
	}
}

func isWaitForBareMetalServerNetworkInterfaceFloatingIpAvailable(client *vpcv1.VpcV1, bareMetalServerId, nicId, fipId string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for Bare Metal Server (%s) Network Interface (%s) to be available.", bareMetalServerId, nicId)
	communicator := make(chan interface{})
	stateConf := &resource.StateChangeConf{
		Pending:    []string{isBareMetalServerNetworkInterfaceFloatingIpPending},
		Target:     []string{isBareMetalServerNetworkInterfaceFloatingIpAvailable, isBareMetalServerNetworkInterfaceFloatingIpFailed},
		Refresh:    isBareMetalServerNetworkInterfaceFloatingIpRefreshFunc(client, bareMetalServerId, nicId, fipId, d, communicator),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isBareMetalServerNetworkInterfaceFloatingIpRefreshFunc(client *vpcv1.VpcV1, bareMetalServerId, nicId, fipId string, d *schema.ResourceData, communicator chan interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getBmsNicFloatingIpOptions := &vpcv1.GetBareMetalServerNetworkInterfaceFloatingIPOptions{
			BareMetalServerID:  &bareMetalServerId,
			NetworkInterfaceID: &nicId,
			ID:                 &fipId,
		}
		fip, response, err := client.GetBareMetalServerNetworkInterfaceFloatingIP(getBmsNicFloatingIpOptions)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) Network Interface (%s) FloatingIp(%s) : %s\n%s", bareMetalServerId, nicId, fipId, err, response)
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

func MakeTerraformNICFipID(id1, id2, id3 string) string {
	// Include bare metal sever id, network interface id, floating ip id to create a unique Terraform id.  As a bonus,
	// we can extract the bare metal sever id as needed for API calls such as READ.
	return fmt.Sprintf("%s/%s/%s", id1, id2, id3)
}

func ParseNICFipTerraformID(s string) (string, string, string, error) {
	segments := strings.Split(s, "/")
	if len(segments) != 3 {
		return "", "", "", fmt.Errorf("invalid terraform Id %s (incorrect number of segments)", s)
	}
	if segments[0] == "" || segments[1] == "" || segments[2] == "" {
		return "", "", "", fmt.Errorf("invalid terraform Id %s (one or more empty segments)", s)
	}
	return segments[0], segments[1], segments[2], nil
}
