package packet

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

var wgMap = map[string]*sync.WaitGroup{}
var wgMutex = sync.Mutex{}

func ifToIPCreateRequest(m interface{}) packngo.IPAddressCreateRequest {
	iacr := packngo.IPAddressCreateRequest{}
	ia := m.(map[string]interface{})
	at := ia["type"].(string)
	switch at {
	case "public_ipv4":
		iacr.AddressFamily = 4
		iacr.Public = true
	case "private_ipv4":
		iacr.AddressFamily = 4
		iacr.Public = false
	case "public_ipv6":
		iacr.AddressFamily = 6
		iacr.Public = true
	}
	iacr.CIDR = ia["cidr"].(int)
	iacr.Reservations = convertStringArr(ia["reservation_ids"].([]interface{}))
	return iacr
}

func getNewIPAddressSlice(arr []interface{}) []packngo.IPAddressCreateRequest {
	addressTypesSlice := make([]packngo.IPAddressCreateRequest, len(arr))

	for i, m := range arr {
		addressTypesSlice[i] = ifToIPCreateRequest(m)
	}
	return addressTypesSlice
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

func waitUntilReservationProvisionable(id string, meta interface{}) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"false"},
		Target:  []string{"true"},
		Refresh: func() (interface{}, string, error) {
			client := meta.(*packngo.Client)
			r, _, err := client.HardwareReservations.Get(id, nil)
			if err != nil {
				return 42, "error", friendlyError(err)
			}
			provisionableString := "false"
			if r.Provisionable {
				provisionableString = "true"
			}
			return 42, provisionableString, nil
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err := stateConf.WaitForState()
	return err
}

func getWaitForDeviceLock(deviceID string) *sync.WaitGroup {
	wgMutex.Lock()
	defer wgMutex.Unlock()
	wg, ok := wgMap[deviceID]
	if !ok {
		wg = &sync.WaitGroup{}
		wgMap[deviceID] = wg
	}
	return wg
}

func waitForDeviceAttribute(d *schema.ResourceData, targets []string, pending []string, attribute string, meta interface{}) (string, error) {

	wg := getWaitForDeviceLock(d.Id())
	wg.Wait()

	wgMutex.Lock()
	wg.Add(1)
	wgMutex.Unlock()

	defer func() {
		wgMutex.Lock()
		wg.Done()
		wgMutex.Unlock()
	}()

	if attribute != "state" && attribute != "network_type" {
		return "", fmt.Errorf("unsupported attr to wait for: %s", attribute)
	}

	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  targets,
		Refresh: func() (interface{}, string, error) {
			client := meta.(*packngo.Client)
			device, _, err := client.Devices.Get(d.Id(), &packngo.GetOptions{Includes: []string{"project"}})
			if err == nil {
				retAttrVal := device.State
				if attribute == "network_type" {
					networkType := device.GetNetworkType()
					retAttrVal = networkType
				}
				return retAttrVal, retAttrVal, nil
			}
			return "error", "error", err
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	attrValRaw, err := stateConf.WaitForState()

	if v, ok := attrValRaw.(string); ok {
		return v, err
	}

	return "", err
}

// powerOnAndWait Powers on the device and waits for it to be active.
func powerOnAndWait(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	_, err := client.Devices.PowerOn(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	_, err = waitForDeviceAttribute(d, []string{"active", "failed"}, []string{"off"}, "state", client)
	if err != nil {
		return err
	}
	state := d.Get("state").(string)
	if state != "active" {
		return friendlyError(fmt.Errorf("Device in non-active state \"%s\"", state))
	}
	return nil
}

func validateFacilityForDevice(v interface{}, k string) (ws []string, errors []error) {
	if v.(string) == "any" {
		errors = append(errors, fmt.Errorf(`Cannot use facility: "any"`))
	}
	return
}

func ipAddressSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(ipAddressTypes, false),
				Description:  fmt.Sprintf("one of %s", strings.Join(ipAddressTypes, ",")),
			},
			"cidr": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "CIDR suffix for IP block assigned to this device",
			},
			"reservation_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "IDs of reservations to pick the blocks from",
				MinItems:    1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringMatch(uuidRE, "must be a valid UUID"),
				},
			},
		},
	}
}
