// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/helpers/location"
	"github.com/softlayer/softlayer-go/helpers/product"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

func resourceIBMIPSecVPN() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIPSecVpnCreate,
		Read:     resourceIBMIPSecVPNRead,
		Delete:   resourceIBMIPSecVPNDelete,
		Update:   resourceIBMIPSecVPNUpdate,
		Exists:   resourceIBMIPSecVPNExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Datacenter name",
			},
			"internal_peer_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phase_one": {
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "MD5",
							ValidateFunc: validateAuthProtocol,
						},
						"encryption": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "3DES",
							ValidateFunc: validateEncyptionProtocol,
						},
						"diffie_hellman_group": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2,
							ValidateFunc: validateDiffieHellmanGroup,
						},
						"keylife": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      14400,
							ValidateFunc: validatekeylife,
						},
					},
				},
			},
			"phase_two": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authentication": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "MD5",
							ValidateFunc: validateAuthProtocol,
						},
						"encryption": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "3DES",
							ValidateFunc: validateEncyptionProtocol,
						},
						"diffie_hellman_group": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2,
							ValidateFunc: validateDiffieHellmanGroup,
						},
						"keylife": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      3600,
							ValidateFunc: validatekeylife,
						},
					},
				},
			},
			"address_translation": { //Parameters for creating an adress translation
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_ip_adress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"internal_ip_adress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"preshared_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Preshared Key data",
			},
			"customer_peer_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Customer Peer IP Address",
			},
			"internal_subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Internal subnet ID value",
			},
			"remote_subnet_id": { //customer subnet id . need atleast one customer subnet id for applying the configuratons
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"remote_subnet"},
				Description:   "Remote subnet ID value",
			},
			"remote_subnet": { //parameters to be populated for creating a customer subnet. Specify only one parameter:- remote subnet/remote subnet id
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"remote_subnet_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_ip_adress": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remote_ip_cidr": {
							Type:         schema.TypeString,
							ValidateFunc: validateCIDR,
							Required:     true,
						},
						"account_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"service_subnet_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Service subnet ID value",
			},
		},
	}
}

const (
	ipsecMask = "billingItem.orderItem.order.id,serviceSubnets,staticRouteSubnets"
)

func resourceIBMIPSecVpnCreate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	datacenter := d.Get("datacenter").(string)
	dc, err := location.GetDatacenterByName(sess, datacenter, "id")
	locationid := strconv.Itoa(*dc.Id)
	packageid := 0
	if err != nil {
		return fmt.Errorf("Datacenter not found")
	}
	locationservice := services.GetLocationService(sess)
	priceidds, _ := locationservice.Id(*dc.Id).GetPriceGroups()
	var listofpriceids []int
	//store all the pricegroups a datacenter belongs to
	for _, priceidd := range priceidds {
		listofpriceids = append(listofpriceids, *priceidd.Id)
	}
	actualpriceid, err := product.GetPriceIDByPackageIdandLocationGroups(sess, listofpriceids, 0, "IPSEC - Standard")
	priceItems := []datatypes.Product_Item_Price{}
	priceItem := datatypes.Product_Item_Price{
		Id: &actualpriceid,
	}
	priceItems = append(priceItems, priceItem)
	IPSecOrder := datatypes.Container_Product_Order_Network_Tunnel_Ipsec{
		Container_Product_Order: datatypes.Container_Product_Order{
			PackageId: sl.Int(packageid),
			Prices:    priceItems,
			Quantity:  sl.Int(1),
			Location:  &locationid,
		},
	}
	//Calling verify order
	_, err = services.GetProductOrderService(sess.SetRetries(0)).
		VerifyOrder(&IPSecOrder)
	if err != nil {
		return fmt.Errorf("Error during Verify order for Creating: %s", err)
	}

	//Calling place order
	receipt, err := services.GetProductOrderService(sess.SetRetries(0)).
		PlaceOrder(&IPSecOrder, sl.Bool(false))
	if err != nil {
		return fmt.Errorf("Error during Place order for Creating: %s", err)
	}
	vpn, _ := findIPSecVpnByOrderID(sess, *receipt.OrderId, d)
	if err != nil {
		return fmt.Errorf("Error during creation of IPSec VPN: %s", err)
	}
	id := *vpn.Id
	d.SetId(fmt.Sprintf("%d", id))
	log.Printf("[INFO] IPSec VPN ID: %s", d.Id())
	return resourceIBMIPSecVPNUpdate(d, meta)
}

func findIPSecVpnByOrderID(sess *session.Session, orderID int, d *schema.ResourceData) (datatypes.Network_Tunnel_Module_Context, error) {
	filterPath := "networkTunnelContexts.billingItem.orderItem.order.id"
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"complete"},
		Refresh: func() (interface{}, string, error) {
			vpn, err := services.GetAccountService(sess).
				Filter(filter.Build(
					filter.Path(filterPath).
						Eq(strconv.Itoa(orderID)))).
				Mask(ipsecMask).
				GetNetworkTunnelContexts()
			if err != nil {
				return datatypes.Network_Tunnel_Module_Context{}, "", err
			}

			if len(vpn) == 1 {
				return vpn[0], "complete", nil
			} else if len(vpn) == 0 {
				return datatypes.Network_Tunnel_Module_Context{}, "pending", nil
			}
			return nil, "", fmt.Errorf("Expected one IPSec VPN: %s", err)
		},
		Timeout:        2 * time.Hour,
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 24 * 60,
	}

	pendingResult, err := stateConf.WaitForState()

	if err != nil {
		return datatypes.Network_Tunnel_Module_Context{}, err
	}
	var result, ok = pendingResult.(datatypes.Network_Tunnel_Module_Context)
	if ok {
		return result, nil
	}

	return datatypes.Network_Tunnel_Module_Context{},
		fmt.Errorf("Cannot find IPSec Vpn with order id '%d'", orderID)
}

func resourceIBMIPSecVPNRead(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	vpnID, _ := strconv.Atoi(d.Id())

	vpn, err := services.GetNetworkTunnelModuleContextService(sess).
		Id(vpnID).Mask(ipsecMask).
		GetObject()
	if err != nil {
		return fmt.Errorf("Error retrieving firewall information: %s", err)
	}
	d.Set("name", *vpn.Name)
	d.Set("internal_peer_ip_address", *vpn.InternalPeerIpAddress)
	if vpn.Datacenter != nil {
		d.Set("datacenter", *vpn.Datacenter.Name)
	}
	d.Set("phase_one", flattenPhaseOneAttributes(&vpn))
	d.Set("phase_two", flattenPhaseTwoAttributes(&vpn))
	fwID, err := strconv.Atoi(d.Id())
	if vpn.AddressTranslations != nil {
		d.Set("address_translation", flattenaddressTranslation(&vpn, fwID))
	}
	if vpn.CustomerSubnets != nil {
		d.Set("remote_subnet", flattenremoteSubnet(&vpn))
	}
	if vpn.PresharedKey != nil {
		d.Set("preshared_key", *vpn.PresharedKey)
	}
	if vpn.CustomerPeerIpAddress != nil {
		d.Set("customer_peer_ip", *vpn.CustomerPeerIpAddress)
	}
	if len(vpn.InternalSubnets) > 0 {
		d.Set("internal_subnet_id", *vpn.InternalSubnets[0].Id)
	}
	if len(vpn.CustomerSubnets) > 0 {
		d.Set("remote_subnet_id", *vpn.CustomerSubnets[0].Id)
	}
	if len(vpn.ServiceSubnets) > 0 {
		d.Set("service_subnet_id", *vpn.ServiceSubnets[0].Id)
	}
	return nil
}

func resourceIBMIPSecVPNExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess := meta.(ClientSession).SoftLayerSession()

	fwID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	_, err = services.GetNetworkTunnelModuleContextService(sess).
		Id(fwID).
		GetObject()

	if err != nil {
		if apiErr, ok := err.(sl.Error); ok && apiErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("Error retrieving vpn information: %s", err)
	}

	return true, nil
}

func resourceIBMIPSecVPNDelete(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	vpnService := services.GetNetworkTunnelModuleContextService(sess)

	vpnID, _ := strconv.Atoi(d.Id())

	// Get billing item associated with the firewall
	billingItem, err := vpnService.Id(vpnID).GetBillingItem()

	if err != nil {
		return fmt.Errorf("Error while looking up billing item associated with the ipsecvpn: %s", err)
	}

	if billingItem.Id == nil {
		return fmt.Errorf("Error while looking up billing item associated with the ipsecvpn: No billing item for ID:%d", vpnID)
	}

	success, err := services.GetBillingItemService(sess).Id(*billingItem.Id).CancelService()
	if err != nil {
		return err
	}

	if !success {
		return fmt.Errorf("SoftLayer reported an unsuccessful cancellation")
	}

	return nil
}

func resourceIBMIPSecVPNUpdate(d *schema.ResourceData, meta interface{}) error {
	sess := meta.(ClientSession).SoftLayerSession()
	vpnID, err := strconv.Atoi(d.Id())
	var addresstranslation datatypes.Network_Tunnel_Module_Context_Address_Translation
	if err != nil {
		return fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	vpn, err := services.GetNetworkTunnelModuleContextService(sess).
		Id(vpnID).Mask(ipsecMask).
		GetObject()

	if err != nil {
		return fmt.Errorf("Error updating storage information: %s", err)
	}
	if d.HasChange("phase_one") {
		for _, e := range d.Get("phase_one").([]interface{}) {
			value := e.(map[string]interface{})
			auth := value["authentication"].(string)
			vpn.PhaseOneAuthentication = &auth
			encryption := value["encryption"].(string)
			vpn.PhaseOneEncryption = &encryption
			diffiehellman := value["diffie_hellman_group"].(int)
			vpn.PhaseOneDiffieHellmanGroup = &diffiehellman
			keylife := value["keylife"].(int)
			vpn.PhaseOneKeylife = &keylife
		}
	}
	if d.HasChange("phase_two") {
		for _, e := range d.Get("phase_two").([]interface{}) {
			value := e.(map[string]interface{})
			auth := value["authentication"].(string)
			vpn.PhaseTwoAuthentication = &auth
			encryption := value["encryption"].(string)
			vpn.PhaseTwoEncryption = &encryption
			diffiehellman := value["diffie_hellman_group"].(int)
			vpn.PhaseTwoDiffieHellmanGroup = &diffiehellman
			keylife := value["keylife"].(int)
			vpn.PhaseTwoKeylife = &keylife
		}
	}
	if d.HasChange("preshared_key") {
		presharedkey := d.Get("preshared_key").(string)
		vpn.PresharedKey = &presharedkey
	}
	if _, ok := d.GetOk("customer_peer_ip"); ok {
		if d.HasChange("customer_peer_ip") {
			customeripaddr := d.Get("customer_peer_ip").(string)
			vpn.CustomerPeerIpAddress = &customeripaddr
		}
		_, err = services.GetNetworkTunnelModuleContextService(sess).Id(vpnID).EditObject(&vpn)
		if err != nil {
			return fmt.Errorf("SoftLayer reported an unsuccessful edit")
		}
	}
	if d.HasChange("internal_subnet_id") {
		subnetid := d.Get("internal_subnet_id").(int)
		_, err = services.GetNetworkTunnelModuleContextService(sess).AddPrivateSubnetToNetworkTunnel(&subnetid)
		if err != nil {
			return fmt.Errorf("Unable to find object with id of: %s", err)
		}
	}
	if d.HasChange("remote_subnet_id") {
		subnetid := d.Get("remote_subnet_id").(int)
		_, err = services.GetNetworkTunnelModuleContextService(sess).AddCustomerSubnetToNetworkTunnel(&subnetid)
		if err != nil {
			return fmt.Errorf("Unable to find object with id of: %s", err)
		}
	}
	if d.HasChange("service_subnet_id") {
		subnetid := d.Get("service_subnet_id").(int)
		_, err = services.GetNetworkTunnelModuleContextService(sess).AddServiceSubnetToNetworkTunnel(&subnetid)
		if err != nil {
			return fmt.Errorf("Unable to find object with id of: %s", err)
		}
	}
	if d.HasChange("address_translation") {
		for _, e := range d.Get("address_translation").([]interface{}) {
			value := e.(map[string]interface{})
			customerIP := value["remote_ip_adress"].(string)
			addresstranslation.CustomerIpAddress = &customerIP
			internalIP := value["internal_ip_adress"].(string)
			addresstranslation.InternalIpAddress = &internalIP
			notes := value["notes"].(string)
			addresstranslation.Notes = &notes
		}
		_, err = services.GetNetworkTunnelModuleContextService(sess).Id(vpnID).CreateAddressTranslation(&addresstranslation)
		if err != nil {
			return fmt.Errorf("Unable to create the address translation: %s", err)
		}
	}
	if d.HasChange("remote_subnet") {
		for _, e := range d.Get("remote_subnet").([]interface{}) {
			remoteSubnet := datatypes.Network_Customer_Subnet{}
			value := e.(map[string]interface{})
			customerIP := value["remote_ip_adress"].(string)
			s := strings.Split(customerIP, "/")
			ip, cidr := s[0], s[1]
			actualcidr, _ := strconv.Atoi(cidr)
			accountID := value["accountid"].(int)
			remoteSubnet.NetworkIdentifier = &ip
			remoteSubnet.Cidr = &actualcidr
			remoteSubnet.AccountId = &accountID
			subnet, err := services.GetNetworkCustomerSubnetService(sess).Id(vpnID).CreateObject(&remoteSubnet)
			if err != nil {
				return fmt.Errorf("Some error occured creating the customer subnet resource %s", err)
			}
			_, err = services.GetNetworkTunnelModuleContextService(sess).Id(vpnID).AddCustomerSubnetToNetworkTunnel(subnet.Id)
			if err != nil {
				return fmt.Errorf("Some error occured adding the customer subnet to the network tunnel module %s", err)
			}

		}
	}
	if _, ok := d.GetOk("remote_subnet_id"); ok {

		_, err = services.GetNetworkTunnelModuleContextService(sess).Id(vpnID).ApplyConfigurationsToDevice()
		if err != nil {
			return fmt.Errorf("There is some erorr applying the configuration %s", err)
		}
	} else if _, ok := d.GetOk("remote_subnet"); ok {
		_, err = services.GetNetworkTunnelModuleContextService(sess).Id(vpnID).ApplyConfigurationsToDevice()
		if err != nil {
			return fmt.Errorf("There is some erorr applying the configuration %s", err)
		}
	}

	return resourceIBMIPSecVPNRead(d, meta)
}
