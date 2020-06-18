package packet

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/packethost/packngo"
)

var matchIPXEScript = regexp.MustCompile(`(?i)^#![i]?pxe`)

func resourcePacketDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketDeviceCreate,
		Read:   resourcePacketDeviceRead,
		Update: resourcePacketDeviceUpdate,
		Delete: resourcePacketDeviceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"operating_system": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"facility": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validateFacilityForDevice,
				Deprecated:    "Use the 'facilities' array instead.",
				ConflictsWith: []string{"facilities"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// ignore set of empty facility "" => "xxx1"
					if new == "" {
						return true
					}
					return false
				},
			},

			"facilities": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"facility"},
			},

			"plan": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"billing_cycle": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"root_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"access_public_ipv6": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_public_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_private_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"layer3", "layer2-bonded", "layer2-individual", "hybrid"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" {
						return true
					}
					return false
				},
			},

			"ports": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"bonded": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"family": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"cidr": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"public_ipv4_subnet_size": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"ipxe_script_url": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"always_pxe": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"hardware_reservation_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "next-available" && len(old) > 0 {
						return true
					}
					return false
				},
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"storage": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					s, _ := structure.NormalizeJsonString(v)
					return s
				},
				ValidateFunc: validation.ValidateJsonString,
			},
			"project_ssh_key_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePacketDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	var facs []string
	f, ok := d.GetOk("facility")

	if ok {
		facs = []string{f.(string)}
	} else {
		facs = convertStringArr(d.Get("facilities").([]interface{}))
		if len(facs) == 0 {
			return fmt.Errorf("You must set either 'facilities' or 'facility'")
		}
	}

	createRequest := &packngo.DeviceCreateRequest{
		Hostname:             d.Get("hostname").(string),
		Plan:                 d.Get("plan").(string),
		Facility:             facs,
		OS:                   d.Get("operating_system").(string),
		BillingCycle:         d.Get("billing_cycle").(string),
		ProjectID:            d.Get("project_id").(string),
		PublicIPv4SubnetSize: d.Get("public_ipv4_subnet_size").(int),
	}
	targetNetworkState, nTypeOk := d.GetOk("network_type")
	if attr, ok := d.GetOk("user_data"); ok {
		createRequest.UserData = attr.(string)
	}

	if attr, ok := d.GetOk("ipxe_script_url"); ok {
		createRequest.IPXEScriptURL = attr.(string)
	}

	if attr, ok := d.GetOk("hardware_reservation_id"); ok {
		createRequest.HardwareReservationID = attr.(string)
	}

	if createRequest.OS == "custom_ipxe" {
		if createRequest.IPXEScriptURL == "" && createRequest.UserData == "" {
			return friendlyError(errors.New("\"ipxe_script_url\" or \"user_data\"" +
				" must be provided when \"custom_ipxe\" OS is selected."))
		}

		// ipxe_script_url + user_data is OK, unless user_data is an ipxe script in
		// which case it's an error.
		if createRequest.IPXEScriptURL != "" {
			if matchIPXEScript.MatchString(createRequest.UserData) {
				return friendlyError(errors.New("\"user_data\" should not be an iPXE " +
					"script when \"ipxe_script_url\" is also provided."))
			}
		}
	}

	if createRequest.OS != "custom_ipxe" && createRequest.IPXEScriptURL != "" {
		return friendlyError(errors.New("\"ipxe_script_url\" argument provided, but" +
			" OS is not \"custom_ipxe\". Please verify and fix device arguments."))
	}

	if attr, ok := d.GetOk("always_pxe"); ok {
		createRequest.AlwaysPXE = attr.(bool)
	}

	projectKeys := d.Get("project_ssh_key_ids.#").(int)
	if projectKeys > 0 {
		createRequest.ProjectSSHKeys = convertStringArr(d.Get("project_ssh_key_ids").([]interface{}))
	}

	tags := d.Get("tags.#").(int)
	if tags > 0 {
		createRequest.Tags = convertStringArr(d.Get("tags").([]interface{}))
	}

	if attr, ok := d.GetOk("storage"); ok {
		s, err := structure.NormalizeJsonString(attr.(string))
		if err != nil {
			return errwrap.Wrapf("storage param contains invalid JSON: {{err}}", err)
		}
		createRequest.Storage = s
	}

	newDevice, _, err := client.Devices.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(newDevice.ID)

	// Wait for the device so we can get the networking attributes that show up after a while.
	_, err = waitForDeviceAttribute(d, "active", []string{"queued", "provisioning"}, "state", meta)
	if err != nil {
		if isForbidden(err) {
			// If the device doesn't get to the active state, we can't recover it from here.
			d.SetId("")

			return errors.New("provisioning time limit exceeded; the Packet team will investigate")
		}
		return err
	}

	if nTypeOk {
		_, err = waitForDeviceAttribute(d, "layer3", []string{"hybrid", "layer2-bonded", "layer2-individual"}, "network_type", meta)

		tns := targetNetworkState.(string)
		if tns != "layer3" {
			_, err := client.DevicePorts.DeviceToNetworkType(newDevice.ID, tns)
			if err != nil {
				return err
			}
		}
	}

	return resourcePacketDeviceRead(d, meta)
}

type NetworkInfo struct {
	Networks       []map[string]interface{}
	IPv4SubnetSize int
	Host           string
	PublicIPv4     string
	PublicIPv6     string
	PrivateIPv4    string
}

func getNetworkInfo(ips []*packngo.IPAddressAssignment) NetworkInfo {
	ni := NetworkInfo{Networks: make([]map[string]interface{}, 0, 1)}
	for _, ip := range ips {
		network := map[string]interface{}{
			"address": ip.Address,
			"gateway": ip.Gateway,
			"family":  ip.AddressFamily,
			"cidr":    ip.CIDR,
			"public":  ip.Public,
		}
		ni.Networks = append(ni.Networks, network)

		// Initial device IPs are fixed and marked as "Management"
		if ip.Management {
			if ip.AddressFamily == 4 {
				if ip.Public {
					ni.Host = ip.Address
					ni.IPv4SubnetSize = ip.CIDR
					ni.PublicIPv4 = ip.Address
				} else {
					ni.PrivateIPv4 = ip.Address
				}
			} else {
				ni.PublicIPv6 = ip.Address
			}
		}
	}
	return ni
}

func resourcePacketDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	device, _, err := client.Devices.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project"}})
	if err != nil {
		err = friendlyError(err)

		// If the device somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("hostname", device.Hostname)
	d.Set("plan", device.Plan.Slug)
	d.Set("facility", device.Facility.Code)
	d.Set("operating_system", device.OS.Slug)
	d.Set("state", device.State)
	d.Set("billing_cycle", device.BillingCycle)
	d.Set("locked", device.Locked)
	d.Set("created", device.Created)
	d.Set("updated", device.Updated)
	d.Set("ipxe_script_url", device.IPXEScriptURL)
	d.Set("always_pxe", device.AlwaysPXE)
	d.Set("root_password", device.RootPassword)
	d.Set("project_id", device.Project.ID)
	storageString, err := structure.FlattenJsonToString(device.Storage)
	if err != nil {
		return fmt.Errorf("[ERR] Error getting storage JSON string for device (%s): %s", d.Id(), err)
	}
	d.Set("storage", storageString)

	if len(device.HardwareReservation.Href) > 0 {
		d.Set("hardware_reservation_id", path.Base(device.HardwareReservation.Href))
	}
	d.Set("network_type", device.NetworkType)

	d.Set("tags", device.Tags)
	keyIDs := []string{}
	for _, k := range device.SSHKeys {
		keyIDs = append(keyIDs, filepath.Base(k.URL))
	}
	d.Set("ssh_key_ids", keyIDs)
	networkInfo := getNetworkInfo(device.Network)

	sort.SliceStable(networkInfo.Networks, func(i, j int) bool {
		famI := networkInfo.Networks[i]["family"].(int)
		famJ := networkInfo.Networks[j]["family"].(int)
		pubI := networkInfo.Networks[i]["public"].(bool)
		pubJ := networkInfo.Networks[j]["public"].(bool)
		return getNetworkRank(famI, pubI) < getNetworkRank(famJ, pubJ)
	})

	d.Set("network", networkInfo.Networks)
	d.Set("public_ipv4_subnet_size", networkInfo.IPv4SubnetSize)
	d.Set("access_public_ipv4", networkInfo.PublicIPv4)
	d.Set("access_private_ipv4", networkInfo.PrivateIPv4)
	d.Set("access_public_ipv6", networkInfo.PublicIPv6)

	ports := getPorts(device.NetworkPorts)
	d.Set("ports", ports)

	if networkInfo.Host != "" {
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": networkInfo.Host,
		})
	}

	return nil
}

func getNetworkRank(family int, public bool) int {
	switch {
	case family == 4 && public:
		return 0
	case family == 6:
		return 1
	case family == 4 && public:
		return 2
	}
	return 3
}

func getPorts(ps []packngo.Port) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0, 1)
	for _, p := range ps {
		port := map[string]interface{}{
			"name":   p.Name,
			"id":     p.ID,
			"type":   p.Type,
			"mac":    p.Data.MAC,
			"bonded": p.Data.Bonded,
		}
		ret = append(ret, port)
	}
	return ret
}

func resourcePacketDeviceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if d.HasChange("locked") {
		var action func(string) (*packngo.Response, error)
		if d.Get("locked").(bool) {
			action = client.Devices.Lock
		} else {
			action = client.Devices.Unlock
		}
		if _, err := action(d.Id()); err != nil {
			return friendlyError(err)
		}
	}
	ur := packngo.DeviceUpdateRequest{}

	if d.HasChange("description") {
		dDesc := d.Get("description").(string)
		ur.Description = &dDesc
	}
	if d.HasChange("hostname") {
		dHostname := d.Get("hostname").(string)
		ur.Hostname = &dHostname
	}
	if d.HasChange("tags") {
		ts := d.Get("tags")
		sts := []string{}

		switch ts.(type) {
		case []interface{}:
			for _, v := range ts.([]interface{}) {
				sts = append(sts, v.(string))
			}
			ur.Tags = &sts
		default:
			return friendlyError(fmt.Errorf("garbage in tags: %s", ts))
		}
	}
	if d.HasChange("ipxe_script_url") {
		dUrl := d.Get("ipxe_script_url").(string)
		ur.IPXEScriptURL = &dUrl
	}
	if d.HasChange("always_pxe") {
		dPXE := d.Get("always_pxe").(bool)
		ur.AlwaysPXE = &dPXE
	}
	if !reflect.DeepEqual(ur, packngo.DeviceUpdateRequest{}) {
		if _, _, err := client.Devices.Update(d.Id(), &ur); err != nil {
			return friendlyError(err)
		}

	}
	if d.HasChange("network_type") {
		target, ok := d.GetOk("network_type")
		if ok {
			targetType := target.(string)
			_, err := client.DevicePorts.DeviceToNetworkType(d.Id(), targetType)
			if err != nil {
				return err
			}
		}
	}
	return resourcePacketDeviceRead(d, meta)
}

func resourcePacketDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	if _, err := client.Devices.Delete(d.Id()); err != nil {
		return friendlyError(err)
	}

	return nil
}

func waitForDeviceAttribute(d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newDeviceStateRefreshFunc(d, attribute, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newDeviceStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*packngo.Client)

	return func() (interface{}, string, error) {
		if err := resourcePacketDeviceRead(d, meta); err != nil {
			return nil, "", err
		}

		if attr, ok := d.GetOk(attribute); ok {
			device, _, err := client.Devices.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project"}})
			if err != nil {
				return nil, "", friendlyError(err)
			}
			return &device, attr.(string), nil
		}

		return nil, "", nil
	}
}

// powerOnAndWait Powers on the device and waits for it to be active.
func powerOnAndWait(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	_, err := client.Devices.PowerOn(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	_, err = waitForDeviceAttribute(d, "active", []string{"off"}, "state", client)
	return err
}

func validateFacilityForDevice(v interface{}, k string) (ws []string, errors []error) {
	if v.(string) == "any" {
		errors = append(errors, fmt.Errorf(`Cannot use facility: "any"`))
	}
	return
}
