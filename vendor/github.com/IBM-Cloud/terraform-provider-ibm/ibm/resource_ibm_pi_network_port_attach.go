// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMPINetworkPortAttach() *schema.Resource {
	return &schema.Resource{

		Create: resourceIBMPINetworkPortAttachCreate,
		Read:   resourceIBMPINetworkPortAttachRead,
		Update: resourceIBMPINetworkPortAttachUpdate,
		Delete: resourceIBMPINetworkPortAttachDelete,
		//Exists:   resourceIBMPINetworkExists,
		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
		},
		Schema: map[string]*schema.Schema{

			"port_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			helpers.PICloudInstanceId: {
				Type:     schema.TypeString,
				Required: true,
			},

			helpers.PIInstanceName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance name to attach the network port to",
			},

			helpers.PINetworkName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network Name - This is the subnet name  in the Cloud instance",
			},

			helpers.PINetworkPortDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human readable description for this network Port",
				Default:     "Port Created via Terraform",
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func resourceIBMPINetworkPortAttachCreate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return err
	}
	powerinstanceid := d.Get(helpers.PICloudInstanceId).(string)
	networkname := d.Get(helpers.PINetworkName).(string)
	portid := d.Get("port_id").(string)
	instancename := d.Get(helpers.PIInstanceName).(string)
	description := d.Get(helpers.PINetworkPortDescription).(string)
	client := st.NewIBMPINetworkClient(sess, powerinstanceid)

	log.Printf("Printing the input to the resource powerinstance [%s] and network name [%s] and the portid [%s]", powerinstanceid, networkname, portid)
	networkPortResponse, err := client.AttachPort(powerinstanceid, networkname, portid, description, instancename, postTimeOut)

	if err != nil {
		return err
	}

	log.Printf("Printing the networkresponse %+v", &networkPortResponse)

	IBMPINetworkPortID := *networkPortResponse.PortID

	d.SetId(fmt.Sprintf("%s/%s", powerinstanceid, IBMPINetworkPortID))
	if err != nil {
		log.Printf("[DEBUG]  err %s", err)
		return err
	}
	_, err = isWaitForIBMPINetworkPortAttachAvailable(client, IBMPINetworkPortID, d.Timeout(schema.TimeoutCreate), powerinstanceid, networkname)
	if err != nil {
		return err
	}

	return resourceIBMPINetworkPortAttachRead(d, meta)
}

func resourceIBMPINetworkPortAttachRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Calling ther Network Port Attach Read code")
	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		fmt.Printf("failed to get  a session from the IBM Cloud Service %v", err)
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	powerinstanceid := parts[0]
	powernetworkname := d.Get(helpers.PINetworkName).(string)
	networkC := st.NewIBMPINetworkClient(sess, powerinstanceid)
	networkdata, err := networkC.GetPort(powernetworkname, powerinstanceid, parts[1], getTimeOut)

	d.Set("ipaddress", networkdata.IPAddress)
	d.Set("macaddress", networkdata.MacAddress)
	d.Set("status", networkdata.Status)
	d.Set("portid", networkdata.PortID)
	d.Set("pvminstance", networkdata.PvmInstance.Href)
	d.Set("public_ip", networkdata.ExternalIP)

	return nil
}

func resourceIBMPINetworkPortAttachUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Calling the attach update ")
	return nil
}

func resourceIBMPINetworkPortAttachDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Detaching the network port from the Instance ")

	sess, err := meta.(ClientSession).IBMPISession()

	if err != nil {
		fmt.Printf("failed to get  a session from the IBM Cloud Service %v", err)

	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	powerinstanceid := parts[0]
	powernetworkname := d.Get(helpers.PINetworkName).(string)
	portid := d.Get("port_id").(string)

	client := st.NewIBMPINetworkClient(sess, powerinstanceid)
	log.Printf("Calling the network delete functions. ")
	network, err := client.DetachPort(powerinstanceid, powernetworkname, portid, deleteTimeOut)
	if err != nil {
		return err
	}

	log.Printf("Printing the networkresponse %+v", &network)

	//log.Printf("Printing the networkresponse %s", network.Status)

	d.SetId("")
	return nil

}

func isWaitForIBMPINetworkPortAttachAvailable(client *st.IBMPINetworkClient, id string, timeout time.Duration, powerinstanceid, networkname string) (interface{}, error) {
	log.Printf("Waiting for Power Network (%s) that was created for Network Zone (%s) to be available.", id, networkname)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"retry", helpers.PINetworkProvisioning},
		Target:     []string{"ACTIVE"},
		Refresh:    isIBMPINetworkPortAttachRefreshFunc(client, id, powerinstanceid, networkname),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Minute,
	}

	return stateConf.WaitForState()
}

func isIBMPINetworkPortAttachRefreshFunc(client *st.IBMPINetworkClient, id, powerinstanceid, networkname string) resource.StateRefreshFunc {

	log.Printf("Calling the IsIBMPINetwork Refresh Function....with the following id (%s) for network port and following id (%s) for network name and waiting for network to be READY", id, networkname)
	return func() (interface{}, string, error) {
		network, err := client.GetPort(networkname, powerinstanceid, id, getTimeOut)
		if err != nil {
			return nil, "", err
		}

		if &network.PortID != nil && &network.PvmInstance.PvmInstanceID != nil {
			//if network.State == "available" {
			log.Printf(" The port has been created with the following ip address and attached to an instance ")
			return network, "ACTIVE", nil
		}

		return network, helpers.PINetworkProvisioning, nil
	}
}
