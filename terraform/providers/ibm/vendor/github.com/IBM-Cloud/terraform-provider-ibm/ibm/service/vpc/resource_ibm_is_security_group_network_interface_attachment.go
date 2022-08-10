// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isSGNICAGroupId               = "security_group"
	isSGNICANicId                 = "network_interface"
	isSGNICAInstanceNwInterfaceID = "instance_network_interface"
	isSGNICAName                  = "name"
	isSGNICAPortSpeed             = "port_speed"
	isSGNICAPrimaryIPV4Address    = "primary_ipv4_address"
	isSGNICASecondaryAddresses    = "secondary_address"
	isSGNICASecurityGroups        = "security_groups"
	isSGNICASecurityGroupCRN      = "crn"
	isSGNICASecurityGroupID       = "id"
	isSGNICASecurityGroupName     = "name"
	isSGNICAStatus                = "status"
	isSGNICASubnet                = "subnet"
	isSGNICAType                  = "type"
	isSGNICAFloatingIps           = "floating_ips"
	isSGNICAFloatingIpID          = "id"
	isSGNICAFloatingIpName        = "name"
	isSGNICAFloatingIpCRN         = "crn"
)

func ResourceIBMISSecurityGroupNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMISSecurityGroupNetworkInterfaceAttachmentCreate,
		Read:     resourceIBMISSecurityGroupNetworkInterfaceAttachmentRead,
		Delete:   resourceIBMISSecurityGroupNetworkInterfaceAttachmentDelete,
		Exists:   resourceIBMISSecurityGroupNetworkInterfaceAttachmentExists,
		Importer: &schema.ResourceImporter{},

		DeprecationMessage: "Resource ibm_is_security_group_network_interface_attachment is deprecated. Use ibm_is_security_group_target to attach a network interface to a security group",

		Schema: map[string]*schema.Schema{
			isSGNICAGroupId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "security group network interface attachment group ID",
			},
			isSGNICANicId: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "security group network interface attachment NIC ID",
			},
			isSGNICAName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group network interface attachment name",
			},
			isSGNICAInstanceNwInterfaceID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group network interface attachment network interface ID",
			},
			isSGNICAPortSpeed: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "security group network interface attachment port speed",
			},
			isSGNICAPrimaryIPV4Address: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group network interface attachment Primary IPV4 address",
			},
			isSGNICASecondaryAddresses: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "security group network interface attachment secondary address",
			},
			isSGNICAStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group network interface attachment status",
			},
			isSGNICASubnet: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group network interface attachment subnet",
			},
			isSGNICAType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "security group network interface attachment type",
			},
			isSGNICAFloatingIps: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isSGNICAFloatingIpID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group network interface attachment floating IP ID",
						},
						isSGNICAFloatingIpName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group network interface attachment floating IP name",
						},
						isSGNICAFloatingIpCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group network interface attachment floating IP CRN",
						},
					},
				},
			},
			isSGNICASecurityGroups: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isSGNICASecurityGroupID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group network interface attachment security group ID",
						},
						isSGNICASecurityGroupCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group network interface attachment security group CRN",
						},
						isSGNICASecurityGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "security group network interface attachment security group name",
						},
					},
				},
			},

			flex.RelatedCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The crn of the Security Group",
			},
		},
	}
}

func resourceIBMISSecurityGroupNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	sgID := d.Get(isSGNICAGroupId).(string)
	nicID := d.Get(isSGNICANicId).(string)

	options := &vpcv1.CreateSecurityGroupTargetBindingOptions{
		SecurityGroupID: &sgID,
		ID:              &nicID,
	}
	_, response, err := sess.CreateSecurityGroupTargetBinding(options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error while creating SecurityGroup NetworkInterface Binding %s\n%s", err, response)
	}
	d.SetId(fmt.Sprintf("%s/%s", sgID, nicID))
	return resourceIBMISSecurityGroupNetworkInterfaceAttachmentRead(d, meta)

}

func resourceIBMISSecurityGroupNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}
	sgID := parts[0]
	nicID := parts[1]

	getSecurityGroupNetworkInterfaceOptions := &vpcv1.GetSecurityGroupTargetOptions{
		SecurityGroupID: &sgID,
		ID:              &nicID,
	}
	secGroupTarget, response, err := sess.GetSecurityGroupTarget(getSecurityGroupNetworkInterfaceOptions)
	if err != nil || secGroupTarget == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting target(%s) for the SecurityGroup (%s) : %s\n%s", nicID, sgID, err, response)
	}
	instance_id := strings.Split(*secGroupTarget.(*vpcv1.SecurityGroupTargetReference).Href, "/")[5]
	net_interf_id := *secGroupTarget.(*vpcv1.SecurityGroupTargetReference).ID
	getnicoptions := &vpcv1.GetInstanceNetworkInterfaceOptions{
		InstanceID: &instance_id,
		ID:         &net_interf_id,
	}
	instanceNic, response, err := sess.GetInstanceNetworkInterface(getnicoptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting network interfaces attached to the instance %s %s\n%s", instance_id, err, response)
	}
	d.Set(isSGNICAGroupId, sgID)
	d.Set(isSGNICANicId, nicID)
	d.Set(isSGNICAInstanceNwInterfaceID, *instanceNic.ID)
	d.Set(isSGNICAName, *instanceNic.Name)
	d.Set(isSGNICAPortSpeed, *instanceNic.PortSpeed)
	if instanceNic.PrimaryIP != nil && instanceNic.PrimaryIP.Address != nil {
		d.Set(isSGNICAPrimaryIPV4Address, *instanceNic.PrimaryIP.Address)
	}
	d.Set(isSGNICAStatus, *instanceNic.Status)
	d.Set(isSGNICAType, *instanceNic.Type)
	if instanceNic.Subnet != nil {
		d.Set(isSGNICASubnet, *instanceNic.Subnet.ID)
	}
	sgs := make([]map[string]interface{}, len(instanceNic.SecurityGroups))
	for i, sgObj := range instanceNic.SecurityGroups {
		sg := make(map[string]interface{})
		sg[isSGNICASecurityGroupCRN] = *sgObj.CRN
		sg[isSGNICASecurityGroupID] = *sgObj.ID
		sg[isSGNICASecurityGroupName] = *sgObj.Name
		sgs[i] = sg
	}
	d.Set(isSGNICASecurityGroups, sgs)

	fps := make([]map[string]interface{}, len(instanceNic.FloatingIps))
	for i, fpObj := range instanceNic.FloatingIps {
		fp := make(map[string]interface{})
		fp[isSGNICAFloatingIpCRN] = fpObj.CRN
		fp[isSGNICAFloatingIpID] = *fpObj.ID
		fp[isSGNICAFloatingIpName] = *fpObj.Name
		fps[i] = fp
	}
	d.Set(isSGNICAFloatingIps, fps)

	// d.Set(isSGNICASecondaryAddresses, *instanceNic.SecondaryAddresses)
	getSecurityGroupOptions := &vpcv1.GetSecurityGroupOptions{
		ID: &sgID,
	}
	sg, response, err := sess.GetSecurityGroup(getSecurityGroupOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Getting Security Group : %s\n%s", err, response)
	}
	d.Set(flex.RelatedCRN, *sg.CRN)
	return nil
}

func resourceIBMISSecurityGroupNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return err
	}

	sgID := parts[0]
	nicID := parts[1]

	getSecurityGroupNetworkInterfaceOptions := &vpcv1.GetSecurityGroupTargetOptions{
		SecurityGroupID: &sgID,
		ID:              &nicID,
	}
	_, response, err := sess.GetSecurityGroupTarget(getSecurityGroupNetworkInterfaceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting NetworkInterface(%s) for the SecurityGroup (%s) : %s\n%s", nicID, sgID, err, response)
	}

	removeSecurityGroupNetworkInterfaceOptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
		SecurityGroupID: &sgID,
		ID:              &nicID,
	}
	response, err = sess.DeleteSecurityGroupTargetBinding(removeSecurityGroupNetworkInterfaceOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting NetworkInterface(%s) for the SecurityGroup (%s) : %s\n%s", nicID, sgID, err, response)
	}
	d.SetId("")
	return nil
}

func resourceIBMISSecurityGroupNetworkInterfaceAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		return false, err
	}
	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) != 2 {
		return false, fmt.Errorf("[ERROR] Incorrect ID %s: ID should be a combination of sgID/nicID", d.Id())
	}
	sgID := parts[0]
	nicID := parts[1]
	getSecurityGroupNetworkInterfaceOptions := &vpcv1.GetSecurityGroupTargetOptions{
		SecurityGroupID: &sgID,
		ID:              &nicID,
	}
	_, response, err := sess.GetSecurityGroupTarget(getSecurityGroupNetworkInterfaceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error getting NetworkInterface(%s) for the SecurityGroup (%s) : %s\n%s", nicID, sgID, err, response)
	}
	return true, nil
}
