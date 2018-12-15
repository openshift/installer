package openstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceComputeInterfaceAttachV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInterfaceAttachV2Create,
		Read:   resourceComputeInterfaceAttachV2Read,
		Delete: resourceComputeInterfaceAttachV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"network_id"},
			},

			"network_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
			},

			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fixed_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeInterfaceAttachV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

	var portId string
	if v, ok := d.GetOk("port_id"); ok {
		portId = v.(string)
	}

	var networkId string
	if v, ok := d.GetOk("network_id"); ok {
		networkId = v.(string)
	}

	if networkId == "" && portId == "" {
		return fmt.Errorf("Must set one of network_id and port_id")
	}

	// For some odd reason the API takes an array of IPs, but you can only have one element in the array.
	var fixedIPs []attachinterfaces.FixedIP
	if v, ok := d.GetOk("fixed_ip"); ok {
		fixedIPs = append(fixedIPs, attachinterfaces.FixedIP{IPAddress: v.(string)})
	}

	attachOpts := attachinterfaces.CreateOpts{
		PortID:    portId,
		NetworkID: networkId,
		FixedIPs:  fixedIPs,
	}

	log.Printf("[DEBUG] Creating interface attachment: %#v", attachOpts)

	attachment, err := attachinterfaces.Create(computeClient, instanceId, attachOpts).Extract()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    resourceComputeInterfaceAttachV2AttachFunc(computeClient, instanceId, attachment.PortID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error attaching interface: %s", err)
	}

	log.Printf("[DEBUG] Created interface attachment: %#v", attachment)

	// Use the instance ID and attachment ID as the resource ID.
	id := fmt.Sprintf("%s/%s", instanceId, attachment.PortID)

	d.SetId(id)

	return resourceComputeInterfaceAttachV2Read(d, meta)
}

func resourceComputeInterfaceAttachV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceId, attachmentId, err := parseComputeInterfaceAttachId(d.Id())
	if err != nil {
		return err
	}

	attachment, err := attachinterfaces.Get(computeClient, instanceId, attachmentId).Extract()
	if err != nil {
		return CheckDeleted(d, err, "compute_interface_attach")
	}

	log.Printf("[DEBUG] Retrieved interface attachment: %#v", attachment)

	d.Set("port_id", attachment.PortID)
	d.Set("network_id", attachment.NetID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeInterfaceAttachV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	instanceId, attachmentId, err := parseComputeInterfaceAttachId(d.Id())
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    resourceComputeInterfaceAttachV2DetachFunc(computeClient, instanceId, attachmentId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error detaching interface: %s", err)
	}

	return nil
}

func resourceComputeInterfaceAttachV2AttachFunc(
	computeClient *gophercloud.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := attachinterfaces.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		return va, "ATTACHED", nil
	}
}

func resourceComputeInterfaceAttachV2DetachFunc(
	computeClient *gophercloud.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to detach OpenStack interface %s from instance %s",
			attachmentId, instanceId)

		va, err := attachinterfaces.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = attachinterfaces.Delete(computeClient, instanceId, attachmentId).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(gophercloud.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] Interface Attachment (%s) is still active.", attachmentId)
		return nil, "", nil
	}
}

func parseComputeInterfaceAttachId(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine interface attachment ID")
	}

	instanceId := idParts[0]
	attachmentId := idParts[1]

	return instanceId, attachmentId, nil
}
