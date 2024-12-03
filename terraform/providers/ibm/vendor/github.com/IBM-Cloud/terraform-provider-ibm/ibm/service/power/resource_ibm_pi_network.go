// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
)

const (
	piEndingIPAaddress   = "pi_ending_ip_address"
	piStartingIPAaddress = "pi_starting_ip_address"
)

func ResourceIBMPINetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPINetworkCreate,
		ReadContext:   resourceIBMPINetworkRead,
		UpdateContext: resourceIBMPINetworkUpdate,
		DeleteContext: resourceIBMPINetworkDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			helpers.PINetworkType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{DHCPVlan, PubVlan, Vlan}),
				Description:  "PI network type",
			},
			helpers.PINetworkName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI network name",
			},
			helpers.PINetworkDNS: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of PI network DNS name",
			},
			helpers.PINetworkCidr: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "PI network CIDR",
			},
			helpers.PINetworkGateway: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "PI network gateway",
			},
			helpers.PINetworkJumbo: {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				Deprecated:    "This field is deprecated, use pi_network_mtu instead.",
				ConflictsWith: []string{helpers.PINetworkMtu},
				Description:   "PI network enable MTU Jumbo option",
			},
			helpers.PINetworkMtu: {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{helpers.PINetworkJumbo},
				Description:   "PI Maximum Transmission Unit",
			},
			helpers.PINetworkAccessConfig: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"internal-only", "outbound-only", "bidirectional-static-route", "bidirectional-bgp", "bidirectional-l2out"}),
				Description:  "PI network communication configuration",
			},
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},
			helpers.PINetworkIPAddressRange: {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of one or more ip address range(s)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						piEndingIPAaddress: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ending ip address",
						},
						piStartingIPAaddress: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Starting ip address",
						},
					},
				},
			},
			Arg_UserTags: {
				Description: "The user tags attached to this resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},

			//Computed Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The CRN of this resource.",
				Type:        schema.TypeString,
			},
			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "PI network ID",
			},
			"vlan_id": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "VLAN Id value",
			},
		},
	}
}

func resourceIBMPINetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}
	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	networkname := d.Get(helpers.PINetworkName).(string)
	networktype := d.Get(helpers.PINetworkType).(string)

	client := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	var body = &models.NetworkCreate{
		Type: &networktype,
		Name: networkname,
	}
	if v, ok := d.GetOk(helpers.PINetworkDNS); ok {
		networkdns := flex.ExpandStringList((v.(*schema.Set)).List())
		if len(networkdns) > 0 {
			body.DNSServers = networkdns
		}
	}
	if tags, ok := d.GetOk(Arg_UserTags); ok {
		body.UserTags = flex.FlattenSet(tags.(*schema.Set))
	}
	if v, ok := d.GetOk(helpers.PINetworkJumbo); ok {
		body.Jumbo = v.(bool)
	}
	if v, ok := d.GetOk(helpers.PINetworkMtu); ok {
		var mtu int64 = int64(v.(int))
		body.Mtu = &mtu
	}
	if v, ok := d.GetOk(helpers.PINetworkAccessConfig); ok {
		body.AccessConfig = models.AccessConfig(v.(string))
	}

	if networktype == DHCPVlan || networktype == Vlan {
		var networkcidr string
		var ipBodyRanges []*models.IPAddressRange
		if v, ok := d.GetOk(helpers.PINetworkCidr); ok {
			networkcidr = v.(string)
		} else {
			return diag.Errorf("%s is required when %s is vlan", helpers.PINetworkCidr, helpers.PINetworkType)
		}

		gateway, firstip, lastip, err := generateIPData(networkcidr)
		if err != nil {
			return diag.FromErr(err)
		}

		ipBodyRanges = []*models.IPAddressRange{{EndingIPAddress: &lastip, StartingIPAddress: &firstip}}

		if g, ok := d.GetOk(helpers.PINetworkGateway); ok {
			gateway = g.(string)
		}

		if ips, ok := d.GetOk(helpers.PINetworkIPAddressRange); ok {
			ipBodyRanges = getIPAddressRanges(ips.([]interface{}))
		}

		body.IPAddressRanges = ipBodyRanges
		body.Gateway = gateway
		body.Cidr = networkcidr
	}
	wsclient := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	wsData, err := wsclient.Get(cloudInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	if wsData.Capabilities[PER] {
		_, err = waitForPERWorkspaceActive(ctx, wsclient, cloudInstanceID, d.Timeout(schema.TimeoutRead))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	networkResponse, err := client.Create(body)
	if err != nil {
		return diag.FromErr(err)
	}

	networkID := *networkResponse.NetworkID

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, networkID))

	_, err = isWaitForIBMPINetworkAvailable(ctx, client, networkID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk(Arg_UserTags); ok {
		if networkResponse.Crn != "" {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, string(networkResponse.Crn), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi snapshot (%s) pi_user_tags during creation: %s", networkID, err)
			}
		}
	}

	return resourceIBMPINetworkRead(ctx, d, meta)
}

func resourceIBMPINetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, networkID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	networkC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	networkdata, err := networkC.Get(networkID)
	if err != nil {
		return diag.FromErr(err)
	}
	if networkdata.Crn != "" {
		d.Set(Attr_CRN, networkdata.Crn)
		tags, err := flex.GetGlobalTagsUsingCRN(meta, string(networkdata.Crn), "", UserTagType)
		if err != nil {
			log.Printf("Error on get of pi network (%s) pi_user_tags: %s", *networkdata.NetworkID, err)
		}
		d.Set(Arg_UserTags, tags)
	}
	d.Set("network_id", networkdata.NetworkID)
	d.Set(helpers.PINetworkCidr, networkdata.Cidr)
	d.Set(helpers.PINetworkDNS, networkdata.DNSServers)
	d.Set("vlan_id", networkdata.VlanID)
	d.Set(helpers.PINetworkName, networkdata.Name)
	d.Set(helpers.PINetworkType, networkdata.Type)
	d.Set(helpers.PINetworkJumbo, networkdata.Jumbo)
	d.Set(helpers.PINetworkMtu, networkdata.Mtu)
	d.Set(helpers.PINetworkAccessConfig, networkdata.AccessConfig)
	d.Set(helpers.PINetworkGateway, networkdata.Gateway)
	ipRangesMap := []map[string]interface{}{}
	if networkdata.IPAddressRanges != nil {
		for _, n := range networkdata.IPAddressRanges {
			if n != nil {
				v := map[string]interface{}{
					piEndingIPAaddress:   n.EndingIPAddress,
					piStartingIPAaddress: n.StartingIPAddress,
				}
				ipRangesMap = append(ipRangesMap, v)
			}
		}
	}
	d.Set(helpers.PINetworkIPAddressRange, ipRangesMap)

	return nil
}

func resourceIBMPINetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, networkID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(helpers.PINetworkName, helpers.PINetworkDNS, helpers.PINetworkGateway, helpers.PINetworkIPAddressRange) {
		networkC := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
		body := &models.NetworkUpdate{
			DNSServers: flex.ExpandStringList((d.Get(helpers.PINetworkDNS).(*schema.Set)).List()),
		}
		networkType := d.Get(helpers.PINetworkType).(string)
		if d.HasChange(helpers.PINetworkIPAddressRange) || d.HasChange(helpers.PINetworkGateway) {
			if networkType == Vlan {
				if d.HasChange(helpers.PINetworkIPAddressRange) {
					body.IPAddressRanges = getIPAddressRanges(d.Get(helpers.PINetworkIPAddressRange).([]interface{}))
				}
				if d.HasChange(helpers.PINetworkGateway) {
					body.Gateway = flex.PtrToString(d.Get(helpers.PINetworkGateway).(string))
				}
			} else {
				return diag.Errorf("%v type does not allow ip-address range or gateway update", networkType)
			}
		}

		if d.HasChange(helpers.PINetworkName) {
			body.Name = flex.PtrToString(d.Get(helpers.PINetworkName).(string))
		}

		_, err = networkC.Update(networkID, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of pi network (%s) pi_user_tags: %s", networkID, err)
			}
		}
	}

	return resourceIBMPINetworkRead(ctx, d, meta)
}

func resourceIBMPINetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	log.Printf("Calling the network delete functions. ")
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, networkID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := instance.NewIBMPINetworkClient(ctx, sess, cloudInstanceID)
	err = client.Delete(networkID)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = isWaitForIBMPINetworkDeleted(ctx, client, networkID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func isWaitForIBMPINetworkAvailable(ctx context.Context, client *instance.IBMPINetworkClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PINetworkProvisioning},
		Target:     []string{"NETWORK_READY"},
		Refresh:    isIBMPINetworkRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPINetworkRefreshFunc(client *instance.IBMPINetworkClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		network, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}

		if network.VlanID != nil {
			return network, "NETWORK_READY", nil
		}

		return network, helpers.PINetworkProvisioning, nil
	}
}

func isWaitForIBMPINetworkDeleted(ctx context.Context, client *instance.IBMPINetworkClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Found},
		Target:     []string{State_NotFound},
		Refresh:    isIBMPINetworkRefreshDeleteFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isIBMPINetworkRefreshDeleteFunc(client *instance.IBMPINetworkClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		network, err := client.Get(id)
		if err != nil {
			return network, State_NotFound, nil
		}
		return network, State_Found, nil
	}
}

func generateIPData(cdir string) (gway, firstip, lastip string, err error) {
	_, ipv4Net, err := net.ParseCIDR(cdir)

	if err != nil {
		return "", "", "", err
	}

	var subnetToSize = map[string]int{
		"21": 2048,
		"22": 1024,
		"23": 512,
		"24": 256,
		"25": 128,
		"26": 64,
		"27": 32,
		"28": 16,
		"29": 8,
		"30": 4,
		"31": 2,
	}

	//subnetsize, _ := ipv4Net.Mask.Size()

	gateway, err := cidr.Host(ipv4Net, 1)
	if err != nil {
		log.Printf("Failed to get the gateway for this cidr passed in %s", cdir)
		return "", "", "", err
	}
	ad := cidr.AddressCount(ipv4Net)

	convertedad := strconv.FormatUint(ad, 10)
	// Powervc in wdc04 has to reserve 3 ip address hence we start from the 4th. This will be the default behaviour
	firstusable, err := cidr.Host(ipv4Net, 4)
	if err != nil {
		log.Print(err)
		return "", "", "", err
	}
	lastusable, err := cidr.Host(ipv4Net, subnetToSize[convertedad]-2)
	if err != nil {
		log.Print(err)
		return "", "", "", err
	}
	return gateway.String(), firstusable.String(), lastusable.String(), nil

}

func getIPAddressRanges(ipAddressRanges []interface{}) []*models.IPAddressRange {
	ipRanges := make([]*models.IPAddressRange, 0, len(ipAddressRanges))
	for _, v := range ipAddressRanges {
		if v != nil {
			ipAddressRange := v.(map[string]interface{})
			ipRange := &models.IPAddressRange{
				EndingIPAddress:   flex.PtrToString(ipAddressRange[piEndingIPAaddress].(string)),
				StartingIPAddress: flex.PtrToString(ipAddressRange[piStartingIPAaddress].(string)),
			}
			ipRanges = append(ipRanges, ipRange)
		}
	}
	return ipRanges
}
func waitForPERWorkspaceActive(ctx context.Context, client *instance.IBMPIWorkspacesClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_Inactive, State_Configuring},
		Target:     []string{State_Active},
		Refresh:    isPERWorkspaceRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func isPERWorkspaceRefreshFunc(client *instance.IBMPIWorkspacesClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ws, err := client.Get(id)
		if err != nil {
			return nil, "", err
		}
		// Check for backward compatibility for legacy workspace.
		if ws.Details.PowerEdgeRouter == nil {
			return ws, State_Active, nil
		}
		if *(ws.Details.PowerEdgeRouter.State) == State_Active {
			return ws, State_Active, nil
		}
		if *(ws.Details.PowerEdgeRouter.State) == State_Inactive {
			return ws, State_Inactive, nil
		}
		if *(ws.Details.PowerEdgeRouter.State) == State_Error {
			return ws, *ws.Details.PowerEdgeRouter.State, fmt.Errorf("[ERROR] workspace PER configuration failed to initialize. Please try again later")
		}

		return ws, State_Configuring, nil
	}
}
