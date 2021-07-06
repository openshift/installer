// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
)

func resourceIBMPINetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMPINetworkCreate,
		Read:   resourceIBMPINetworkRead,
		Update: resourceIBMPINetworkUpdate,
		Delete: resourceIBMPINetworkDelete,
		//Exists:   resourceIBMPINetworkExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PINetworkType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"vlan", "pub-vlan"}),
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
				Description: "PI network gateway",
			},

			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
			},

			//Computed Attributes

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

func resourceIBMPINetworkCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	networkname := d.Get(helpers.PINetworkName).(string)
	networktype := d.Get(helpers.PINetworkType).(string)
	networkcidr := d.Get(helpers.PINetworkCidr).(string)
	networkdns := expandStringList((d.Get(helpers.PINetworkDNS).(*schema.Set)).List())

	client := st.NewIBMPINetworkClient(sess, powerinstanceid)
	var networkgateway, firstip, lastip string
	if networktype == "vlan" {
		networkgateway, firstip, lastip = generateIPData(networkcidr)
	}
	networkResponse, _, err := client.Create(networkname, networktype, networkcidr, networkdns, networkgateway, firstip, lastip, powerinstanceid, postTimeOut)
	if err != nil {
		return err
	}

	IBMPINetworkID := *networkResponse.NetworkID

	d.SetId(fmt.Sprintf("%s/%s", powerinstanceid, IBMPINetworkID))

	_, err = isWaitForIBMPINetworkAvailable(client, IBMPINetworkID, d.Timeout(schema.TimeoutCreate), powerinstanceid)
	if err != nil {
		return err
	}

	return resourceIBMPINetworkRead(d, meta)
}

func resourceIBMPINetworkRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	powerinstanceid := parts[0]
	networkC := st.NewIBMPINetworkClient(sess, powerinstanceid)
	networkdata, err := networkC.Get(parts[1], powerinstanceid, getTimeOut)
	if err != nil {
		return err
	}

	d.Set("network_id", networkdata.NetworkID)
	d.Set(helpers.PINetworkCidr, networkdata.Cidr)
	d.Set(helpers.PINetworkDNS, networkdata.DNSServers)
	d.Set("vlan_id", networkdata.VlanID)
	d.Set(helpers.PINetworkName, networkdata.Name)
	d.Set(helpers.PINetworkType, networkdata.Type)
	d.Set(helpers.PICloudInstanceId, powerinstanceid)

	return nil

}

func resourceIBMPINetworkUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMPINetworkDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("Calling the network delete functions. ")
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	powerinstanceid := parts[0]
	networkC := st.NewIBMPINetworkClient(sess, powerinstanceid)
	err = networkC.Delete(parts[1], powerinstanceid, deleteTimeOut)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceIBMPINetworkExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return false, err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	powerinstanceid := parts[0]
	client := st.NewIBMPINetworkClient(sess, powerinstanceid)

	network, err := client.Get(parts[0], powerinstanceid, getTimeOut)
	if err != nil {

		return false, err
	}
	return *network.NetworkID == parts[1], nil
}

func isWaitForIBMPINetworkAvailable(client *st.IBMPINetworkClient, id string, timeout time.Duration, powerinstanceid string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PINetworkProvisioning},
		Target:     []string{"NETWORK_READY"},
		Refresh:    isIBMPINetworkRefreshFunc(client, id, powerinstanceid),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isIBMPINetworkRefreshFunc(client *st.IBMPINetworkClient, id, powerinstanceid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		network, err := client.Get(id, powerinstanceid, getTimeOut)
		if err != nil {
			return nil, "", err
		}

		if &network.VlanID != nil {
			return network, "NETWORK_READY", nil
		}

		return network, helpers.PINetworkProvisioning, nil
	}
}

func generateIPData(cdir string) (gway, firstip, lastip string) {
	_, ipv4Net, err := net.ParseCIDR(cdir)

	if err != nil {
		log.Fatal(err)
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
		log.Printf("Failed to get the gateway for this cdir passed in %s", cdir)
		log.Fatal(err)
	}
	ad := cidr.AddressCount(ipv4Net)

	convertedad := strconv.FormatUint(ad, 10)
	// Powervc in wdc04 has to reserve 3 ip address hence we start from the 4th. This will be the default behaviour
	firstusable, err := cidr.Host(ipv4Net, 4)
	if err != nil {
		log.Fatal(err)
	}
	lastusable, err := cidr.Host(ipv4Net, subnetToSize[convertedad]-2)
	if err != nil {
		log.Fatal(err)
	}
	return gateway.String(), firstusable.String(), lastusable.String()

}
