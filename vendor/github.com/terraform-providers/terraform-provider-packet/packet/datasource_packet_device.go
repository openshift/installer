package packet

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/packethost/packngo"
)

func dataSourcePacketDevice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketDeviceRead,
		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"device_id"},
			},
			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"device_id"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id", "hostname"},
			},
			"facility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_system": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_cycle": {
				Type:     schema.TypeString,
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
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hardware_reservation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage": {
				Type: schema.TypeString,
				StateFunc: func(v interface{}) string {
					s, _ := structure.NormalizeJsonString(v)
					return s
				},
				Computed: true,
			},
			"root_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"always_pxe": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ipxe_script_url": {
				Type:     schema.TypeString,
				Computed: true,
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
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bonded": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePacketDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	hostnameRaw, hostnameOK := d.GetOk("hostname")
	projectIdRaw, projectIdOK := d.GetOk("project_id")
	deviceIdRaw, deviceIdOK := d.GetOk("device_id")

	if !deviceIdOK && !hostnameOK {
		return fmt.Errorf("You must supply device_id or hostname")
	}
	var device *packngo.Device
	if hostnameOK {
		if !projectIdOK {
			return fmt.Errorf("If you lookup via hostname, you must supply project_id")
		}
		hostname := hostnameRaw.(string)
		projectId := projectIdRaw.(string)

		ds, _, err := client.Devices.List(projectId, nil)
		if err != nil {
			return err
		}

		device, err = findDeviceByHostname(ds, hostname)
		if err != nil {
			return err
		}
	} else {
		deviceId := deviceIdRaw.(string)
		var err error
		device, _, err = client.Devices.Get(deviceId, nil)
		if err != nil {
			return err
		}
	}

	d.Set("hostname", device.Hostname)
	d.Set("project_id", device.Project.ID)
	d.Set("device_id", device.ID)
	d.Set("plan", device.Plan.Slug)
	d.Set("facility", device.Facility.Code)
	d.Set("operating_system", device.OS.Slug)
	d.Set("state", device.State)
	d.Set("billing_cycle", device.BillingCycle)
	d.Set("ipxe_script_url", device.IPXEScriptURL)
	d.Set("always_pxe", device.AlwaysPXE)
	d.Set("root_password", device.RootPassword)
	if device.Storage != nil {
		rawStorageBytes, err := json.Marshal(device.Storage)
		if err != nil {
			return fmt.Errorf("[ERR] Error getting storage JSON string for device (%s): %s", d.Id(), err)
		}

		storageString, err := structure.NormalizeJsonString(string(rawStorageBytes))
		if err != nil {
			return fmt.Errorf("[ERR] Errori normalizing storage JSON string for device (%s): %s", d.Id(), err)
		}
		d.Set("storage", storageString)
	}

	if len(device.HardwareReservation.Href) > 0 {
		d.Set("hardware_reservation_id", path.Base(device.HardwareReservation.Href))
	}
	networkType, err := device.GetNetworkType()
	if err != nil {
		return err
	}

	d.Set("network_type", networkType)

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
	d.Set("access_public_ipv4", networkInfo.PublicIPv4)
	d.Set("access_private_ipv4", networkInfo.PrivateIPv4)
	d.Set("access_public_ipv6", networkInfo.PublicIPv6)

	ports := getPorts(device.NetworkPorts)
	d.Set("ports", ports)

	d.SetId(device.ID)
	return nil
}

func findDeviceByHostname(devices []packngo.Device, hostname string) (*packngo.Device, error) {
	results := make([]packngo.Device, 0)
	for _, d := range devices {
		if d.Hostname == hostname {
			results = append(results, d)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no device found with hostname %s", hostname)
	}
	return nil, fmt.Errorf("too many devices found with hostname %s (found %d, expected 1)", hostname, len(results))
}
