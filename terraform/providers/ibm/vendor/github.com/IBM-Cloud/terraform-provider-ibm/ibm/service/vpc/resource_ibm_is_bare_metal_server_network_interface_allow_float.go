// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIsBareMetalServerNetworkInterfaceAllowFloat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISBareMetalServerNetworkInterfaceAllowFloatCreate,
		ReadContext:   resourceIBMISBareMetalServerNetworkInterfaceAllowFloatRead,
		UpdateContext: resourceIBMISBareMetalServerNetworkInterfaceAllowFloatUpdate,
		DeleteContext: resourceIBMISBareMetalServerNetworkInterfaceAllowFloatDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			isBareMetalServerID: {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: flex.ApplyOnce,
				Description:      "Bare metal server identifier",
			},
			isBareMetalServerNicID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bare metal server network interface identifier",
			},
			isBareMetalServerNicAllowIPSpoofing: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether source IP spoofing is allowed on this interface. If false, source IP spoofing is prevented on this interface. If true, source IP spoofing is allowed on this interface.",
			},
			isBareMetalServerNicEnableInfraNAT: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If true, the VPC infrastructure performs any needed NAT operations. If false, the packet is passed unmodified to/from the network interface, allowing the workload to perform any needed NAT operations.",
			},
			isBareMetalServerNicFloatingIPs: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The floating IPs associated with this network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerNicIpAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique IP address",
						},
					},
				},
			},
			isBareMetalServerNicHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this network interface",
			},
			isBareMetalServerNicInterfaceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network interface type: [ pci, vlan ]",
			},
			isBareMetalServerNicMacAddress: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The MAC address of the interface. If absent, the value is not known.",
			},
			isBareMetalServerNicName: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user-defined name for this network interface",
			},
			isBareMetalServerNicPortSpeed: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The network interface port speed in Mbps",
			},
			isBareMetalServerNicPrimaryIP: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "title: IPv4, The IP address. ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerNicIpAddress: {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The globally unique IP address",
						},
					},
				},
			},
			isBareMetalServerNicResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type : [ subnet_reserved_ip ]",
			},

			isBareMetalServerNicSecurityGroups: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Collection of security groups ids",
			},

			isBareMetalServerNicStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the network interface : [ available, deleting, failed, pending ]",
			},

			isBareMetalServerNicSubnet: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the associated subnet",
			},

			isBareMetalServerNicType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of this bare metal server network interface : [ primary, secondary ]",
			},
			isBareMetalServerNicVlan: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this interface",
			},
			isBareMetalServerNicAllowInterfaceToFloat: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates if the interface can float to any other server within the same resource_group. The interface will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to vlan type interfaces.",
			},
		},
	}
}

func resourceIBMISBareMetalServerNetworkInterfaceAllowFloatCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	bareMetalServerId := ""
	if bmsId, ok := d.GetOk(isBareMetalServerID); ok {
		bareMetalServerId = bmsId.(string)
	}

	err := createVlanTypeNetworkInterfaceAllowFloat(context, d, meta, bareMetalServerId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func createVlanTypeNetworkInterfaceAllowFloat(context context.Context, d *schema.ResourceData, meta interface{}, bareMetalServerId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	options := &vpcv1.CreateBareMetalServerNetworkInterfaceOptions{}
	interfaceType := "vlan"
	nicOptions := &vpcv1.BareMetalServerNetworkInterfacePrototypeBareMetalServerNetworkInterfaceByVlanPrototype{}
	allowInterfaceToFloat := true
	nicOptions.AllowInterfaceToFloat = &allowInterfaceToFloat
	if vlan, ok := d.GetOk(isBareMetalServerNicVlan); ok {
		vlanInt := int64(vlan.(int))
		nicOptions.Vlan = &vlanInt
	}

	if name, ok := d.GetOk(isBareMetalServerNicName); ok {
		nameStr := name.(string)
		nicOptions.Name = &nameStr
	}
	nicOptions.InterfaceType = &interfaceType

	if ais, ok := d.GetOk(isBareMetalServerNicAllowIPSpoofing); ok {
		allowIPSpoofing := ais.(bool)
		if allowIPSpoofing {
			nicOptions.AllowIPSpoofing = &allowIPSpoofing
		}
	}
	if ein, ok := d.GetOk(isBareMetalServerNicEnableInfraNAT); ok {
		enableInfrastructureNat := ein.(bool)
		nicOptions.EnableInfrastructureNat = &enableInfrastructureNat
	}
	if subnetOk, ok := d.GetOk(isBareMetalServerNicSubnet); ok {
		subnet := subnetOk.(string)
		nicOptions.Subnet = &vpcv1.SubnetIdentity{
			ID: &subnet,
		}
	}

	if primaryIpIntf, ok := d.GetOk(isBareMetalServerNicPrimaryIP); ok && len(primaryIpIntf.([]interface{})) > 0 {
		primaryIp := primaryIpIntf.([]interface{})[0].(map[string]interface{})

		reservedIpAddressOk, ok := primaryIp[isBareMetalServerNicIpAddress]
		if ok && reservedIpAddressOk.(string) != "" {
			reservedIpAddress := reservedIpAddressOk.(string)
			nicOptions.PrimaryIpv4Address = &reservedIpAddress
		}
	}

	sGroups := d.Get(isBareMetalServerNicSecurityGroups).(*schema.Set).List()
	var sGroupList []vpcv1.SecurityGroupIdentityIntf
	// Add new allowed_subnets
	for _, sGroup := range sGroups {
		sGroupStr := sGroup.(string)
		sgModel := &vpcv1.SecurityGroupIdentity{
			ID: &sGroupStr,
		}
		sGroupList = append(sGroupList, sgModel)
	}
	nicOptions.SecurityGroups = sGroupList
	options.BareMetalServerID = &bareMetalServerId
	options.BareMetalServerNetworkInterfacePrototype = nicOptions
	nic, response, err := sess.CreateBareMetalServerNetworkInterfaceWithContext(context, options)
	if err != nil || nic == nil {
		return fmt.Errorf("[DEBUG] Create bare metal server (%s) network interface err %s\n%s", bareMetalServerId, err, response)
	}
	err = bareMetalServerNICGet(d, meta, nic, bareMetalServerId)
	if err != nil {
		return err
	}
	_, nicId, err := ParseNICTerraformID(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Bare Metal Server Network Interface : %s", d.Id())
	_, err = isWaitForBareMetalServerNetworkInterfaceAvailable(sess, bareMetalServerId, nicId, d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return err
	}
	return nil
}

func resourceIBMISBareMetalServerNetworkInterfaceAllowFloatRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerId, nicID, err := ParseNICTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
		BareMetalServerID: &bareMetalServerId,
		ID:                &nicID,
	}
	var nicIntf vpcv1.BareMetalServerNetworkInterfaceIntf
	// try to fetch original nic
	nicIntf, response, err := sess.GetBareMetalServerNetworkInterfaceWithContext(context, options)
	if (err != nil || nicIntf == nil) && response != nil {
		//if original nic is not present, try fetching nic without server id
		nicIntf, response, err = findNicsWithoutBMS(context, sess, nicID)
		// response here can be either nil or not nil and if it returns 404 means nic is deleted
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		// if response returns an error
		if err != nil || nicIntf == nil {
			if response != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) network interface (%s): %s\n%s", bareMetalServerId, nicID, err, response))
			} else {
				return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) network interface (%s): %s", bareMetalServerId, nicID, err))
			}
		}
	}
	err = bareMetalServerNICAllowFloatGet(d, meta, nicIntf, bareMetalServerId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func findNicsWithoutBMS(context context.Context, sess *vpcv1.VpcV1, nicId string) (result vpcv1.BareMetalServerNetworkInterfaceIntf, response *core.DetailedResponse, err error) {
	// listing all servers
	start := ""
	allrecs := []vpcv1.BareMetalServer{}
	for {
		listBareMetalServersOptions := &vpcv1.ListBareMetalServersOptions{}
		if start != "" {
			listBareMetalServersOptions.Start = &start
		}
		availableServers, response, err := sess.ListBareMetalServersWithContext(context, listBareMetalServersOptions)
		if err != nil {
			return nil, nil, fmt.Errorf("[ERROR] Error fetching Bare Metal Servers %s\n%s", err, response)
		}
		start = flex.GetNext(availableServers.Next)
		allrecs = append(allrecs, availableServers.BareMetalServers...)
		if start == "" {
			break
		}
	}
	// finding nic id each server
	for _, server := range allrecs {
		nics := server.NetworkInterfaces
		for _, nic := range nics {
			if *nic.ID == nicId {
				options := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
					BareMetalServerID: server.ID,
					ID:                &nicId,
				}
				//return response of the server nic matches
				return sess.GetBareMetalServerNetworkInterfaceWithContext(context, options)
			}
		}
	}
	// if not found return nil response and error
	return nil, nil, fmt.Errorf("[ERROR] Error Network interface not found")
}

func bareMetalServerNICAllowFloatGet(d *schema.ResourceData, meta interface{}, nicIntf interface{}, bareMetalServerId string) error {
	switch reflect.TypeOf(nicIntf).String() {
	case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
		{
			nic := nicIntf.(*vpcv1.BareMetalServerNetworkInterfaceByPci)
			d.SetId(MakeTerraformNICID(bareMetalServerId, *nic.ID))
			d.Set(isBareMetalServerNicAllowIPSpoofing, *nic.AllowIPSpoofing)
			d.Set(isBareMetalServerNicEnableInfraNAT, *nic.EnableInfrastructureNat)
			d.Set(isBareMetalServerNicStatus, *nic.Status)

			floatingIPList := make([]map[string]interface{}, 0)
			if nic.FloatingIps != nil {
				for _, ip := range nic.FloatingIps {
					currentIP := map[string]interface{}{
						isBareMetalServerNicIpID:      *ip.ID,
						isBareMetalServerNicIpAddress: *ip.Address,
					}
					floatingIPList = append(floatingIPList, currentIP)
				}
			}
			d.Set(isBareMetalServerNicFloatingIPs, floatingIPList)

			d.Set(isBareMetalServerNicHref, *nic.Href)
			d.Set(isBareMetalServerNicID, *nic.ID)
			d.Set(isBareMetalServerNicInterfaceType, *nic.InterfaceType)

			d.Set(isBareMetalServerNicMacAddress, *nic.MacAddress)
			d.Set(isBareMetalServerNicName, *nic.Name)
			if nic.PortSpeed != nil {
				d.Set(isBareMetalServerNicPortSpeed, *nic.PortSpeed)
			}
			primaryIpList := make([]map[string]interface{}, 0)
			currentIP := map[string]interface{}{

				isBareMetalServerNicIpAddress: *nic.PrimaryIpv4Address,
			}
			primaryIpList = append(primaryIpList, currentIP)
			d.Set(isBareMetalServerNicPrimaryIP, primaryIpList)

			d.Set(isBareMetalServerNicResourceType, *nic.ResourceType)

			if nic.SecurityGroups != nil && len(nic.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(nic.SecurityGroups); i++ {
					secgrpList = append(secgrpList, string(*(nic.SecurityGroups[i].ID)))
				}
				d.Set(isBareMetalServerNicSecurityGroups, flex.NewStringSet(schema.HashString, secgrpList))
			}

			d.Set(isBareMetalServerNicStatus, *nic.Status)
			d.Set(isBareMetalServerNicSubnet, *nic.Subnet.ID)
			d.Set(isBareMetalServerNicType, *nic.Type)
		}
	case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
		{
			nic := nicIntf.(*vpcv1.BareMetalServerNetworkInterfaceByVlan)
			d.SetId(MakeTerraformNICID(bareMetalServerId, *nic.ID))
			d.Set(isBareMetalServerNicAllowIPSpoofing, *nic.AllowIPSpoofing)
			d.Set(isBareMetalServerNicEnableInfraNAT, *nic.EnableInfrastructureNat)
			d.Set(isBareMetalServerNicStatus, *nic.Status)
			floatingIPList := make([]map[string]interface{}, 0)
			if nic.FloatingIps != nil {
				for _, ip := range nic.FloatingIps {
					currentIP := map[string]interface{}{
						isBareMetalServerNicIpID:      *ip.ID,
						isBareMetalServerNicIpAddress: *ip.Address,
					}
					floatingIPList = append(floatingIPList, currentIP)
				}
			}
			d.Set(isBareMetalServerNicFloatingIPs, floatingIPList)

			d.Set(isBareMetalServerNicHref, nic.Href)
			d.Set(isBareMetalServerNicID, *nic.ID)
			d.Set(isBareMetalServerNicInterfaceType, *nic.InterfaceType)

			d.Set(isBareMetalServerNicMacAddress, *nic.MacAddress)
			d.Set(isBareMetalServerNicName, *nic.Name)
			d.Set(isBareMetalServerNicPortSpeed, nic.PortSpeed)

			primaryIpList := make([]map[string]interface{}, 0)
			currentIP := map[string]interface{}{

				isBareMetalServerNicIpAddress: *nic.PrimaryIpv4Address,
			}
			primaryIpList = append(primaryIpList, currentIP)
			d.Set(isBareMetalServerNicPrimaryIP, primaryIpList)

			d.Set(isBareMetalServerNicResourceType, nic.ResourceType)

			if len(nic.SecurityGroups) != 0 {
				secgrpList := []string{}
				for i := 0; i < len(nic.SecurityGroups); i++ {
					secgrpList = append(secgrpList, string(*(nic.SecurityGroups[i].ID)))
				}
				d.Set(isBareMetalServerNicSecurityGroups, flex.NewStringSet(schema.HashString, secgrpList))
			}

			d.Set(isBareMetalServerNicStatus, *nic.Status)
			d.Set(isBareMetalServerNicSubnet, *nic.Subnet.ID)
			d.Set(isBareMetalServerNicType, *nic.Type)
			d.Set(isBareMetalServerNicAllowInterfaceToFloat, *nic.AllowInterfaceToFloat)
			d.Set(isBareMetalServerNicVlan, *nic.Vlan)
		}
	}
	return nil
}

func resourceIBMISBareMetalServerNetworkInterfaceAllowFloatUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	bareMetalServerId, nicId, err := ParseNICTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	options := &vpcv1.UpdateBareMetalServerNetworkInterfaceOptions{
		BareMetalServerID: &bareMetalServerId,
		ID:                &nicId,
	}
	nicPatchModel := &vpcv1.BareMetalServerNetworkInterfacePatch{}
	flag := false
	if d.HasChange(isBareMetalServerNicAllowIPSpoofing) {
		flag = true
		aisBool := false
		if ais, ok := d.GetOk(isBareMetalServerNicEnableInfraNAT); ok {
			aisBool = ais.(bool)
		}
		nicPatchModel.AllowIPSpoofing = &aisBool
	}
	if d.HasChange(isBareMetalServerNicEnableInfraNAT) {
		flag = true
		einBool := false
		if ein, ok := d.GetOk(isBareMetalServerNicEnableInfraNAT); ok {
			einBool = ein.(bool)
		}
		nicPatchModel.EnableInfrastructureNat = &einBool
	}
	if d.HasChange(isBareMetalServerNicName) {
		flag = true
		nameStr := ""
		if name, ok := d.GetOk(isBareMetalServerNicName); ok {
			nameStr = name.(string)
		}
		nicPatchModel.Name = &nameStr
	}

	if flag {
		nicPatchModelAsPatch, err := nicPatchModel.AsPatch()
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error calling asPatch for BareMetalServerNetworkInterfacePatch %s", err))
		}
		options.BareMetalServerNetworkInterfacePatch = nicPatchModelAsPatch

		nicIntf, response, err := sess.UpdateBareMetalServerNetworkInterfaceWithContext(context, options)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error updating Bare Metal Server: %s\n%s", err, response))
		}
		return diag.FromErr(bareMetalServerNICGet(d, meta, nicIntf, bareMetalServerId))
	}

	return nil
}

func resourceIBMISBareMetalServerNetworkInterfaceAllowFloatDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerId, nicId, err := ParseNICTerraformID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = bareMetalServerNetworkInterfaceAllowFloatDelete(context, d, meta, bareMetalServerId, nicId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func bareMetalServerNetworkInterfaceAllowFloatDelete(context context.Context, d *schema.ResourceData, meta interface{}, bareMetalServerId, nicId string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	getBmsNicOptions := &vpcv1.GetBareMetalServerNetworkInterfaceOptions{
		BareMetalServerID: &bareMetalServerId,
		ID:                &nicId,
	}
	nicIntf, response, err := sess.GetBareMetalServerNetworkInterfaceWithContext(context, getBmsNicOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) network interface(%s) : %s\n%s", bareMetalServerId, nicId, err, response)
	}
	nicType := ""
	switch reflect.TypeOf(nicIntf).String() {
	case "*vpcv1.BareMetalServerNetworkInterfaceByPci":
		{
			nicType = "pci"
			log.Printf("[DEBUG] PCI type network interface needs the server in stopped state")
			log.Printf("[DEBUG] Stopping the bare metal server %s", bareMetalServerId)
			// to delete pci, server needs to be in stopped state

			getbmsoptions := &vpcv1.GetBareMetalServerOptions{
				ID: &bareMetalServerId,
			}

			bms, response, err := sess.GetBareMetalServerWithContext(context, getbmsoptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error fetching bare metal server (%s) err %s\n%s", bareMetalServerId, err, response)
			}
			// failed, pending, restarting, running, starting, stopped, stopping, maintenance
			if *bms.Status == "failed" {
				return fmt.Errorf("[ERROR] Error cannot detach network interface from a failed bare metal server")
			} else if *bms.Status == "running" {
				log.Printf("[DEBUG] Stopping bare metal server (%s) to create a PCI network interface", bareMetalServerId)
				stopType := "soft"
				if d.Get(isBareMetalServerHardStop).(bool) {
					stopType = "hard"
				}
				createstopaction := &vpcv1.StopBareMetalServerOptions{
					ID:   &bareMetalServerId,
					Type: &stopType,
				}
				res, err := sess.StopBareMetalServerWithContext(context, createstopaction)
				if err != nil || res.StatusCode != 204 {
					return fmt.Errorf("[ERROR] Error stopping bare metal server (%s) err %s\n%s", bareMetalServerId, err, response)
				}
				_, err = isWaitForBareMetalServerStoppedForNIC(sess, bareMetalServerId, d.Timeout(schema.TimeoutCreate), d)
				if err != nil || res.StatusCode != 204 {
					return err
				}
			} else if *bms.Status != "stopped" {
				return fmt.Errorf("[ERROR] Error bare metal server in %s state, please try after some time", *bms.Status)
			}
		}
	case "*vpcv1.BareMetalServerNetworkInterfaceByVlan":
		{
			nicType = "vlan"
		}
	}

	options := &vpcv1.DeleteBareMetalServerNetworkInterfaceOptions{
		BareMetalServerID: &bareMetalServerId,
		ID:                &nicId,
	}
	response, err = sess.DeleteBareMetalServerNetworkInterfaceWithContext(context, options)
	if err != nil {
		return fmt.Errorf("[ERROR] Error Deleting Bare Metal Server (%s) network interface (%s) : %s\n%s", bareMetalServerId, nicId, err, response)
	}
	_, err = isWaitForBareMetalServerNetworkInterfaceDeleted(sess, bareMetalServerId, nicId, nicType, nicIntf, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
