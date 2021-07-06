// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_networks"
	"github.com/IBM-Cloud/power-go-client/power/models"
)

func resourceIBMPINetworkPort() *schema.Resource {
	return &schema.Resource{
		Create: resourceIBMPINetworkPortCreate,
		Read:   resourceIBMPINetworkPortRead,
		Update: resourceIBMPINetworkPortUpdate,
		Delete: resourceIBMPINetworkPortDelete,
		//Exists:   resourceIBMPINetworkExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			helpers.PINetworkName: {
				Type:     schema.TypeString,
				Required: true,
			},

			helpers.PICloudInstanceId: {
				Type:     schema.TypeString,
				Required: true,
			},

			helpers.PINetworkPortDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},

			//Computed Attributes

			helpers.PINetworkPortIPAddress: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"macaddress": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"portid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMPINetworkPortCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	networkname := d.Get(helpers.PINetworkName).(string)
	description := d.Get(helpers.PINetworkPortDescription).(string)

	ipaddress := d.Get(helpers.PINetworkPortIPAddress).(string)

	nwportBody := &models.NetworkPortCreate{Description: description}

	if ipaddress != "" {
		log.Printf("IP address provided. ")
		nwportBody.IPAddress = ipaddress
	}

	client := st.NewIBMPINetworkClient(sess, powerinstanceid)

	networkPortResponse, err := client.CreatePort(networkname, powerinstanceid, &p_cloud_networks.PcloudNetworksPortsPostParams{Body: nwportBody}, postTimeOut)

	if err != nil {
		return err
	}

	log.Printf("Printing the networkresponse %+v", &networkPortResponse)

	IBMPINetworkPortID := *networkPortResponse.PortID

	d.SetId(fmt.Sprintf("%s/%s/%s", powerinstanceid, IBMPINetworkPortID, networkname))
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return err
	}
	_, err = isWaitForIBMPINetworkPortAvailable(client, IBMPINetworkPortID, d.Timeout(schema.TimeoutCreate), powerinstanceid, networkname)
	if err != nil {
		return err
	}

	return resourceIBMPINetworkPortRead(d, meta)
}

func resourceIBMPINetworkPortRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}
	var powernetworkname string
	if len(parts) > 2 {
		powernetworkname = parts[2]
	} else {
		powernetworkname = d.Get(helpers.PINetworkName).(string)
		d.SetId(fmt.Sprintf("%s/%s", d.Id(), powernetworkname))
	}

	powerinstanceid := parts[0]
	networkC := st.NewIBMPINetworkClient(sess, powerinstanceid)
	networkdata, err := networkC.GetPort(powernetworkname, powerinstanceid, parts[1], getTimeOut)

	if err != nil {
		return err
	}

	d.Set(helpers.PINetworkPortIPAddress, networkdata.IPAddress)
	d.Set("macaddress", networkdata.MacAddress)
	d.Set("status", networkdata.Status)
	d.Set("portid", networkdata.PortID)
	d.Set("public_ip", networkdata.ExternalIP)

	return nil

}

func resourceIBMPINetworkPortUpdate(data *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMPINetworkPortDelete(d *schema.ResourceData, meta interface{}) error {

	log.Printf("Calling the network delete functions. ")
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())

	if err != nil {
		return err
	}
	var powernetworkname string
	if len(parts) > 2 {
		powernetworkname = parts[2]
	} else {
		powernetworkname = d.Get(helpers.PINetworkName).(string)
	}
	powerinstanceid := parts[0]
	client := st.NewIBMPINetworkClient(sess, powerinstanceid)
	log.Printf("Calling the client %v", client)

	log.Printf("Calling the delete with the following params delete with cloudinstance -> (%s) and networkid -->  (%s) and portid --> (%s) ", powerinstanceid, powernetworkname, parts[1])
	networkdata, err := client.DeletePort(powernetworkname, powerinstanceid, parts[1], deleteTimeOut)

	log.Printf("Response from the deleteport call %v", networkdata)

	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func isWaitForIBMPINetworkPortAvailable(client *st.IBMPINetworkClient, id string, timeout time.Duration, powerinstanceid, networkname string) (interface{}, error) {
	log.Printf("Waiting for Power Network (%s) that was created for Network Zone (%s) to be available.", id, networkname)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PINetworkProvisioning},
		Target:     []string{"DOWN"},
		Refresh:    isIBMPINetworkPortRefreshFunc(client, id, powerinstanceid, networkname),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Minute,
	}

	return stateConf.WaitForState()
}

func isIBMPINetworkPortRefreshFunc(client *st.IBMPINetworkClient, id, powerinstanceid, networkname string) resource.StateRefreshFunc {

	log.Printf("Calling the IsIBMPINetwork Refresh Function....with the following id (%s) for network port and following id (%s) for network name and waiting for network to be READY", id, networkname)
	return func() (interface{}, string, error) {
		network, err := client.GetPort(networkname, powerinstanceid, id, getTimeOut)
		if err != nil {
			return nil, "", err
		}

		if &network.PortID != nil {
			//if network.State == "available" {
			log.Printf(" The port has been created with the following ip address and attached to an instance ")
			return network, "DOWN", nil
		}

		return network, helpers.PINetworkProvisioning, nil
	}
}
