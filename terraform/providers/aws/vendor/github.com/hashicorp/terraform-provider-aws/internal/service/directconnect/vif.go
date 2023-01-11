package directconnect

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directconnect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func virtualInterfaceRead(id string, conn *directconnect.DirectConnect) (*directconnect.VirtualInterface, error) {
	resp, state, err := virtualInterfaceStateRefresh(conn, id)()
	if err != nil {
		return nil, fmt.Errorf("error reading Direct Connect virtual interface (%s): %s", id, err)
	}
	if state == directconnect.VirtualInterfaceStateDeleted {
		return nil, nil
	}

	return resp.(*directconnect.VirtualInterface), nil
}

func virtualInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn

	if d.HasChange("mtu") {
		req := &directconnect.UpdateVirtualInterfaceAttributesInput{
			Mtu:                aws.Int64(int64(d.Get("mtu").(int))),
			VirtualInterfaceId: aws.String(d.Id()),
		}
		log.Printf("[DEBUG] Modifying Direct Connect virtual interface attributes: %s", req)
		_, err := conn.UpdateVirtualInterfaceAttributes(req)
		if err != nil {
			return fmt.Errorf("error modifying Direct Connect virtual interface (%s) attributes: %s", d.Id(), err)
		}
	}
	if d.HasChange("sitelink_enabled") {
		req := &directconnect.UpdateVirtualInterfaceAttributesInput{
			EnableSiteLink:     aws.Bool(d.Get("sitelink_enabled").(bool)),
			VirtualInterfaceId: aws.String(d.Id()),
		}
		log.Printf("[DEBUG] Modifying Direct Connect virtual interface attributes: %s", req)
		_, err := conn.UpdateVirtualInterfaceAttributes(req)
		if err != nil {
			return fmt.Errorf("error modifying Direct Connect virtual interface (%s) attributes: %s", d.Id(), err)
		}
	}

	arn := d.Get("arn").(string)
	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		if err := UpdateTags(conn, arn, o, n); err != nil {
			return fmt.Errorf("error updating Direct Connect virtual interface (%s) tags: %s", arn, err)
		}
	}

	return nil
}

func virtualInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).DirectConnectConn

	log.Printf("[DEBUG] Deleting Direct Connect virtual interface: %s", d.Id())
	_, err := conn.DeleteVirtualInterface(&directconnect.DeleteVirtualInterfaceInput{
		VirtualInterfaceId: aws.String(d.Id()),
	})
	if err != nil {
		if tfawserr.ErrMessageContains(err, directconnect.ErrCodeClientException, "does not exist") {
			return nil
		}
		return fmt.Errorf("error deleting Direct Connect virtual interface (%s): %s", d.Id(), err)
	}

	deleteStateConf := &resource.StateChangeConf{
		Pending: []string{
			directconnect.VirtualInterfaceStateAvailable,
			directconnect.VirtualInterfaceStateConfirming,
			directconnect.VirtualInterfaceStateDeleting,
			directconnect.VirtualInterfaceStateDown,
			directconnect.VirtualInterfaceStatePending,
			directconnect.VirtualInterfaceStateRejected,
			directconnect.VirtualInterfaceStateVerifying,
		},
		Target: []string{
			directconnect.VirtualInterfaceStateDeleted,
		},
		Refresh:    virtualInterfaceStateRefresh(conn, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = deleteStateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for Direct Connect virtual interface (%s) to be deleted: %s", d.Id(), err)
	}

	return nil
}

func virtualInterfaceStateRefresh(conn *directconnect.DirectConnect, vifId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeVirtualInterfaces(&directconnect.DescribeVirtualInterfacesInput{
			VirtualInterfaceId: aws.String(vifId),
		})
		if err != nil {
			return nil, "", err
		}

		n := len(resp.VirtualInterfaces)
		switch n {
		case 0:
			return "", directconnect.VirtualInterfaceStateDeleted, nil

		case 1:
			vif := resp.VirtualInterfaces[0]
			return vif, aws.StringValue(vif.VirtualInterfaceState), nil

		default:
			return nil, "", fmt.Errorf("Found %d Direct Connect virtual interfaces for %s, expected 1", n, vifId)
		}
	}
}

func virtualInterfaceWaitUntilAvailable(conn *directconnect.DirectConnect, vifId string, timeout time.Duration, pending, target []string) error {
	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    virtualInterfaceStateRefresh(conn, vifId),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for Direct Connect virtual interface (%s) to become available: %s", vifId, err)
	}

	return nil
}
