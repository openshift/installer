// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isNetworkInterfacePending   = "pending"
	isNetworkInterfaceAvailable = "available"
	isNetworkInterfaceFailed    = "failed"
	isNetworkInterfaceDeleting  = "deleting"
	isNetworkInterfaceDeleted   = "deleted"
)

func resourceIBMIsInstanceNetworkInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsInstanceNetworkInterfaceCreate,
		ReadContext:   resourceIBMIsInstanceNetworkInterfaceRead,
		UpdateContext: resourceIBMIsInstanceNetworkInterfaceUpdate,
		DeleteContext: resourceIBMIsInstanceNetworkInterfaceDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the instance.",
			},
			isInstanceNicSubnet: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the subnet.",
			},
			isInstanceNicAllowIPSpoofing: &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether source IP spoofing is allowed on this interface. If false, source IP spoofing is prevented on this interface. If true, source IP spoofing is allowed on this interface.",
			},
			isInstanceNicName: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_network_interface", isInstanceNicName),
				Description:  "The user-defined name for this network interface. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			isInstanceNicPrimaryIpv4Address: &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: InvokeValidator("ibm_is_instance_network_interface", isInstanceNicPrimaryIpv4Address),
				Description:  "The primary IPv4 address. If specified, it must be an available address on the network interface's subnet. If unspecified, an available address on the subnet will be automatically selected.",
			},
			"network_interface": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique ID of this network interface",
			},
			isInstanceNicSecurityGroups: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the network interface was created.",
			},
			isInstanceNicFloatingIP: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the floating IP to attach to this network interface",
			},
			isInstanceNicFloatingIPs: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The floating IPs associated with this network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique IP address.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this floating IP.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this floating IP.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unique identifier for this floating IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this floating IP.",
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this network interface.",
			},
			"port_speed": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The network interface port speed in Mbps.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the network interface.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of this network interface as it relates to an instance.",
			},
		},
	}
}

func resourceIBMIsInstanceNetworkInterfaceValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 isInstanceNicName,
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		ValidateSchema{
			Identifier:                 isInstanceNicPrimaryIpv4Address,
			ValidateFunctionIdentifier: ValidateRegexp,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_instance_network_interface", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsInstanceNetworkInterfaceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	instance_id := d.Get("instance").(string)

	createInstanceNetworkInterfaceOptions := &vpcv1.CreateInstanceNetworkInterfaceOptions{}

	createInstanceNetworkInterfaceOptions.SetInstanceID(instance_id)

	subnetId := d.Get(isInstanceNicSubnet).(string)
	subnetIdentity := vpcv1.SubnetIdentity{
		ID: &subnetId,
	}
	createInstanceNetworkInterfaceOptions.SetSubnet(&subnetIdentity)
	if allow_ip_spoofing, ok := d.GetOk(isInstanceNicAllowIPSpoofing); ok {
		createInstanceNetworkInterfaceOptions.SetAllowIPSpoofing(allow_ip_spoofing.(bool))
	}
	if name, ok := d.GetOk(isInstanceNicName); ok {
		createInstanceNetworkInterfaceOptions.SetName(name.(string))
	}
	if primary_ipv4, ok := d.GetOk(isInstanceNicPrimaryIpv4Address); ok {
		createInstanceNetworkInterfaceOptions.SetPrimaryIpv4Address(primary_ipv4.(string))
	}

	if secgrpintf, ok := d.GetOk(isInstanceNicSecurityGroups); ok {
		secgrpSet := secgrpintf.(*schema.Set)
		if secgrpSet.Len() != 0 {
			var secgrpobjs = make([]vpcv1.SecurityGroupIdentityIntf, secgrpSet.Len())
			for i, secgrpIntf := range secgrpSet.List() {
				secgrpIntfstr := secgrpIntf.(string)
				secgrpobjs[i] = &vpcv1.SecurityGroupIdentity{
					ID: &secgrpIntfstr,
				}
			}
			createInstanceNetworkInterfaceOptions.SecurityGroups = secgrpobjs
		}
	}

	isNICKey := "instance_network_interface_key_" + instance_id
	ibmMutexKV.Lock(isNICKey)
	defer ibmMutexKV.Unlock(isNICKey)

	networkInterface, response, err := vpcClient.CreateInstanceNetworkInterfaceWithContext(context, createInstanceNetworkInterfaceOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateInstanceNetworkInterfaceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateInstanceNetworkInterfaceWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createInstanceNetworkInterfaceOptions.InstanceID, *networkInterface.ID))
	d.Set("network_interface", *networkInterface.ID)

	_, err = isWaitForNetworkInterfaceAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error occured while waiting for network interface %s", err))
	}

	if floating_ip_Intf, ok := d.GetOk(isInstanceNicFloatingIP); ok && floating_ip_Intf.(string) != "" {
		floating_ip_id := floating_ip_Intf.(string)

		addInstanceNetworkInterfaceFloatingIPOptions := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{
			InstanceID:         createInstanceNetworkInterfaceOptions.InstanceID,
			NetworkInterfaceID: networkInterface.ID,
			ID:                 &floating_ip_id,
		}

		_, response, err := vpcClient.AddInstanceNetworkInterfaceFloatingIP(addInstanceNetworkInterfaceFloatingIPOptions)

		if err != nil {
			d.Set(isInstanceNicFloatingIP, "")
			return diag.FromErr(fmt.Errorf("[DEBUG] Error adding Floating IP to network interface %s\n%s", err, response))
		}
		_, err = isWaitForNetworkInterfaceAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error occured while waiting for network interface %s", err))
		}

	}

	_, err = isWaitForInstanceAvailable(vpcClient, instance_id, d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceIBMIsInstanceNetworkInterfaceRead(context, d, meta)
}

func resourceIBMIsInstanceNetworkInterfaceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceNetworkInterfaceOptions := &vpcv1.GetInstanceNetworkInterfaceOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceNetworkInterfaceOptions.SetInstanceID(parts[0])
	getInstanceNetworkInterfaceOptions.SetID(parts[1])

	networkInterface, response, err := vpcClient.GetInstanceNetworkInterfaceWithContext(context, getInstanceNetworkInterfaceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetInstanceNetworkInterfaceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetInstanceNetworkInterfaceWithContext failed %s\n%s", err, response))
	}
	d.SetId(fmt.Sprintf("%s/%s", parts[0], *networkInterface.ID))
	d.Set("network_interface", *networkInterface.ID)
	d.Set("instance", parts[0])
	if err = d.Set(isInstanceNicSubnet, *networkInterface.Subnet.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting subnet: %s", err))
	}
	if err = d.Set(isInstanceNicAllowIPSpoofing, *networkInterface.AllowIPSpoofing); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting allow_ip_spoofing: %s", err))
	}
	if err = d.Set(isInstanceNicName, *networkInterface.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set(isInstanceNicPrimaryIpv4Address, *networkInterface.PrimaryIpv4Address); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting primary_ipv4_address: %s", err))
	}
	if networkInterface.SecurityGroups != nil && len(networkInterface.SecurityGroups) != 0 {
		secgrpList := []string{}
		for _, secGrp := range networkInterface.SecurityGroups {
			secgrpList = append(secgrpList, string(*(secGrp.ID)))
		}
		d.Set("security_groups", newStringSet(schema.HashString, secgrpList))
	}

	if err = d.Set("created_at", dateTimeToString(networkInterface.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	floatingIps := []map[string]interface{}{}
	if networkInterface.FloatingIps != nil {

		for _, floatingIpsItem := range networkInterface.FloatingIps {
			floatingIpsItemMap := resourceIBMIsInstanceNetworkInterfaceFloatingIPReferenceToMap(floatingIpsItem)
			floatingIps = append(floatingIps, floatingIpsItemMap)
		}
	}
	if err = d.Set(isInstanceNicFloatingIPs, floatingIps); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting floating_ips: %s", err))
	}
	if err = d.Set("href", networkInterface.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("port_speed", intValue(networkInterface.PortSpeed)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting port_speed: %s", err))
	}
	if err = d.Set("resource_type", networkInterface.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("status", networkInterface.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	if err = d.Set("type", networkInterface.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	return nil
}

func resourceIBMIsInstanceNetworkInterfaceFloatingIPReferenceToMap(floatingIPReference vpcv1.FloatingIPReference) map[string]interface{} {
	floatingIPReferenceMap := map[string]interface{}{}

	floatingIPReferenceMap["address"] = floatingIPReference.Address
	floatingIPReferenceMap["crn"] = floatingIPReference.CRN
	if floatingIPReference.Deleted != nil {
		DeletedMap := resourceIBMIsInstanceNetworkInterfaceFloatingIPReferenceDeletedToMap(*floatingIPReference.Deleted)
		floatingIPReferenceMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	floatingIPReferenceMap["href"] = floatingIPReference.Href
	floatingIPReferenceMap["id"] = floatingIPReference.ID
	floatingIPReferenceMap["name"] = floatingIPReference.Name

	return floatingIPReferenceMap
}

func resourceIBMIsInstanceNetworkInterfaceFloatingIPReferenceDeletedToMap(floatingIPReferenceDeleted vpcv1.FloatingIPReferenceDeleted) map[string]interface{} {
	floatingIPReferenceDeletedMap := map[string]interface{}{}

	floatingIPReferenceDeletedMap["more_info"] = floatingIPReferenceDeleted.MoreInfo

	return floatingIPReferenceDeletedMap
}

func resourceIBMIsInstanceNetworkInterfaceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateInstanceNetworkInterfaceOptions := &vpcv1.UpdateInstanceNetworkInterfaceOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	instance_id := parts[0]
	network_interface_id := parts[1]
	updateInstanceNetworkInterfaceOptions.SetInstanceID(instance_id)
	updateInstanceNetworkInterfaceOptions.SetID(network_interface_id)

	hasChange := false

	patchVals := &vpcv1.NetworkInterfacePatch{}

	if d.HasChange(isInstanceNicAllowIPSpoofing) {
		patchVals.AllowIPSpoofing = core.BoolPtr(d.Get(isInstanceNicAllowIPSpoofing).(bool))
		hasChange = true
	}
	if d.HasChange(isInstanceNicName) {
		patchVals.Name = core.StringPtr(d.Get(isInstanceNicName).(string))
		hasChange = true
	}
	if d.HasChange(isInstanceNicSecurityGroups) && !d.IsNewResource() {

		ovs, nvs := d.GetChange(isInstanceNicSecurityGroups)
		ov := ovs.(*schema.Set)
		nv := nvs.(*schema.Set)
		remove := expandStringList(ov.Difference(nv).List())
		add := expandStringList(nv.Difference(ov).List())
		if len(add) > 0 {
			for i := range add {
				createsgnicoptions := &vpcv1.CreateSecurityGroupTargetBindingOptions{
					SecurityGroupID: &add[i],
					ID:              &network_interface_id,
				}
				_, response, err := vpcClient.CreateSecurityGroupTargetBinding(createsgnicoptions)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error while creating security group %q for network interface of instance %s\n%s: %q", add[i], d.Id(), err, response))
				}
				_, err = isWaitForInstanceAvailable(vpcClient, instance_id, d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return diag.FromErr(err)
				}
			}

		}
		if len(remove) > 0 {
			for i := range remove {
				deletesgnicoptions := &vpcv1.DeleteSecurityGroupTargetBindingOptions{
					SecurityGroupID: &remove[i],
					ID:              &network_interface_id,
				}
				response, err := vpcClient.DeleteSecurityGroupTargetBinding(deletesgnicoptions)
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error while removing security group %q for network interface of instance %s\n%s: %q", remove[i], d.Id(), err, response))
				}
				_, err = isWaitForInstanceAvailable(vpcClient, instance_id, d.Timeout(schema.TimeoutUpdate), d)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
		hasChange = true
	}
	if hasChange {
		isNICKey := "instance_network_interface_key_" + instance_id
		ibmMutexKV.Lock(isNICKey)
		defer ibmMutexKV.Unlock(isNICKey)
		updateInstanceNetworkInterfaceOptions.NetworkInterfacePatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateInstanceNetworkInterfaceWithContext(context, updateInstanceNetworkInterfaceOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateInstanceNetworkInterfaceWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateInstanceNetworkInterfaceWithContext failed %s\n%s", err, response))
		}
	}

	if d.HasChange(isInstanceNicFloatingIP) {
		floating_ip_id_old, floating_ip_id_new := d.GetChange(isInstanceNicFloatingIP)
		instance_id := parts[0]
		network_interface_id := parts[1]
		if floating_ip_id_new == nil || floating_ip_id_new.(string) == "" {
			removeInstanceNetworkInterfaceFloatingIPOptions := vpcClient.NewRemoveInstanceNetworkInterfaceFloatingIPOptions(instance_id, network_interface_id, floating_ip_id_old.(string))
			response, err := vpcClient.RemoveInstanceNetworkInterfaceFloatingIP(removeInstanceNetworkInterfaceFloatingIPOptions)
			if err != nil {
				if response.StatusCode == 404 {
					log.Println("[DEBUG] The specified floating IP address is not associated with the network interface with the specified identifier. ", err.Error())
				} else {
					return diag.FromErr(fmt.Errorf("Error de-associating the floating ip %s in network interface %s of instance %s, %s\n%s", floating_ip_id_old.(string), network_interface_id, instance_id, err, response))
				}
			}
		} else {
			floating_ip_id := floating_ip_id_new.(string)
			getFloatingIPOptions := &vpcv1.GetFloatingIPOptions{
				ID: &floating_ip_id,
			}
			floatingip, response, err := vpcClient.GetFloatingIP(getFloatingIPOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					d.SetId("")
					return nil
				}
				return diag.FromErr(fmt.Errorf("Error Getting Floating IP (%s): %s\n%s", floating_ip_id, err, response))

			}

			if floatingip != nil && floatingip.Target != nil {
				floatingIpTarget := floatingip.Target.(*vpcv1.FloatingIPTarget)
				if *floatingIpTarget.ID != network_interface_id {
					d.Set(isInstanceNicFloatingIP, "")
					return diag.FromErr(fmt.Errorf("[Error] Provided floating ip is already bound to another resource"))
				}
			}

			addInstanceNetworkInterfaceFloatingIPOptions := &vpcv1.AddInstanceNetworkInterfaceFloatingIPOptions{
				InstanceID:         &instance_id,
				NetworkInterfaceID: &network_interface_id,
				ID:                 &floating_ip_id,
			}

			_, response, err = vpcClient.AddInstanceNetworkInterfaceFloatingIP(addInstanceNetworkInterfaceFloatingIPOptions)

			if err != nil {
				d.Set(isInstanceNicFloatingIP, "")
				return diag.FromErr(fmt.Errorf("[DEBUG] Error adding Floating IP to network interface %s\n%s", err, response))
			}
		}

	}

	_, err = isWaitForNetworkInterfaceAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error occured while waiting for network interface %s", err))
	}
	_, err = isWaitForInstanceAvailable(vpcClient, instance_id, d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceIBMIsInstanceNetworkInterfaceRead(context, d, meta)
}

func resourceIBMIsInstanceNetworkInterfaceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteInstanceNetworkInterfaceOptions := &vpcv1.DeleteInstanceNetworkInterfaceOptions{}

	parts, err := sepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	instance_id := parts[0]
	network_intf_id := parts[1]
	isNICKey := "instance_network_interface_key_" + instance_id
	ibmMutexKV.Lock(isNICKey)
	defer ibmMutexKV.Unlock(isNICKey)

	deleteInstanceNetworkInterfaceOptions.SetInstanceID(instance_id)
	deleteInstanceNetworkInterfaceOptions.SetID(network_intf_id)

	response, err := vpcClient.DeleteInstanceNetworkInterfaceWithContext(context, deleteInstanceNetworkInterfaceOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteInstanceNetworkInterfaceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteInstanceNetworkInterfaceWithContext failed %s\n%s", err, response))
	}

	_, err = isWaitForNetworkInterfaceDelete(vpcClient, d.Id(), d.Timeout(schema.TimeoutUpdate), d)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error occured while waiting for network interface %s", err))
	}

	_, err = isWaitForInstanceAvailable(vpcClient, instance_id, d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func isWaitForNetworkInterfaceAvailable(vpcClient *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for dedicated host (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isNetworkInterfacePending},
		Target:     []string{isNetworkInterfaceAvailable, isNetworkInterfaceFailed},
		Refresh:    isNetworkInterfaceRefreshFunc(vpcClient, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isNetworkInterfaceRefreshFunc(vpcClient *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getInstanceNetworkInterfaceOptions := &vpcv1.GetInstanceNetworkInterfaceOptions{}
		parts, err := sepIdParts(d.Id(), "/")
		if err != nil {
			return nil, "", fmt.Errorf("Error splitting ID in parts %s", err)
		}

		getInstanceNetworkInterfaceOptions.SetInstanceID(parts[0])
		getInstanceNetworkInterfaceOptions.SetID(parts[1])

		networkInterface, response, err := vpcClient.GetInstanceNetworkInterface(getInstanceNetworkInterfaceOptions)
		if err != nil {
			return nil, "", fmt.Errorf("GetInstanceNetworkInterface failed %s\n%s", err, response)
		}
		d.Set("status", *networkInterface.Status)

		if *networkInterface.Status == isNetworkInterfaceFailed {
			return networkInterface, *networkInterface.Status, fmt.Errorf("Network Interface creationg failed with status %s ", *networkInterface.Status)
		}
		return networkInterface, *networkInterface.Status, nil
	}
}

func isWaitForNetworkInterfaceDelete(vpcClient *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for dedicated host (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isNetworkInterfacePending, isNetworkInterfaceDeleting, isNetworkInterfaceAvailable},
		Target:     []string{isNetworkInterfaceDeleted},
		Refresh:    isNetworkInterfaceRefreshDeleteFunc(vpcClient, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isNetworkInterfaceRefreshDeleteFunc(vpcClient *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		getInstanceNetworkInterfaceOptions := &vpcv1.GetInstanceNetworkInterfaceOptions{}
		parts, err := sepIdParts(d.Id(), "/")
		if err != nil {
			return nil, "", fmt.Errorf("Error splitting ID in parts %s", err)
		}

		getInstanceNetworkInterfaceOptions.SetInstanceID(parts[0])
		getInstanceNetworkInterfaceOptions.SetID(parts[1])

		networkInterface, response, err := vpcClient.GetInstanceNetworkInterface(getInstanceNetworkInterfaceOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return networkInterface, isNetworkInterfaceDeleted, nil
			}
			return nil, "", fmt.Errorf("GetInstanceNetworkInterface failed %s\n%s", err, response)
		}
		d.Set("status", *networkInterface.Status)

		if *networkInterface.Status == isNetworkInterfaceFailed {
			return networkInterface, *networkInterface.Status, fmt.Errorf("Network Interface creationg failed with status %s ", *networkInterface.Status)
		}
		return networkInterface, *networkInterface.Status, nil
	}
}
