package nutanix

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

const (
	// CDROM ...
	CDROM = "CDROM"
)

func expandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, utils.StringPtr(v.(string)))
		}
	}
	return vs
}

func convertMapInterfaceToMapString(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)
	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)

		mapString[strKey] = strValue
	}
	return mapString
}

func flattenNicListStatus(nics []*v3.VMNicOutputStatus) []map[string]interface{} {
	nicLists := make([]map[string]interface{}, 0)
	if nics != nil {
		nicLists = make([]map[string]interface{}, len(nics))
		for k, v := range nics {
			nic := make(map[string]interface{})
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["floating_ip"] = utils.StringValue(v.FloatingIP)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)
			var ipEndpointList []map[string]interface{}
			for _, v1 := range v.IPEndpointList {
				ipEndpoint := make(map[string]interface{})
				ipEndpoint["ip"] = utils.StringValue(v1.IP)
				ipEndpoint["type"] = utils.StringValue(v1.Type)
				ipEndpointList = append(ipEndpointList, ipEndpoint)
			}
			nic["ip_endpoint_list"] = ipEndpointList
			nic["network_function_chain_reference"] = flattenReferenceValues(v.NetworkFunctionChainReference)

			if v.SubnetReference != nil {
				nic["subnet_uuid"] = utils.StringValue(v.SubnetReference.UUID)
				nic["subnet_name"] = utils.StringValue(v.SubnetReference.Name)
			}

			if v.IsConnected != nil {
				nic["is_connected"] = strconv.FormatBool(utils.BoolValue(v.IsConnected))
			}

			nicLists[k] = nic
		}
	}

	return nicLists
}

func flattenNicList(nics []*v3.VMNic) []map[string]interface{} {
	nicLists := make([]map[string]interface{}, 0)
	if nics != nil {
		nicLists = make([]map[string]interface{}, len(nics))
		for k, v := range nics {
			nic := make(map[string]interface{})
			nic["nic_type"] = utils.StringValue(v.NicType)
			nic["uuid"] = utils.StringValue(v.UUID)
			nic["network_function_nic_type"] = utils.StringValue(v.NetworkFunctionNicType)
			nic["mac_address"] = utils.StringValue(v.MacAddress)
			nic["model"] = utils.StringValue(v.Model)
			var ipEndpointList []map[string]interface{}
			for _, v1 := range v.IPEndpointList {
				if utils.StringValue(v1.Type) != "LEARNED" {
					ipEndpoint := make(map[string]interface{})
					ipEndpoint["ip"] = utils.StringValue(v1.IP)
					ipEndpoint["type"] = utils.StringValue(v1.Type)
					ipEndpointList = append(ipEndpointList, ipEndpoint)
				}
			}
			nic["ip_endpoint_list"] = ipEndpointList
			nic["network_function_chain_reference"] = flattenReferenceValues(v.NetworkFunctionChainReference)

			if v.SubnetReference != nil {
				nic["subnet_uuid"] = utils.StringValue(v.SubnetReference.UUID)
				nic["subnet_name"] = utils.StringValue(v.SubnetReference.Name)
			}

			if v.IsConnected != nil {
				nic["is_connected"] = strconv.FormatBool(utils.BoolValue(v.IsConnected))
			}

			nicLists[k] = nic
		}
	}

	return nicLists
}

func usesGuestCustomization(d *schema.ResourceData) bool {
	keys := []string{
		"guest_customization_cloud_init_user_data",
		"guest_customization_cloud_init_meta_data",
		"guest_customization_cloud_init_custom_key_values",
		"guest_customization_is_overridable",
		"guest_customization_sysprep",
		"guest_customization_sysprep_custom_key_values"}
	for _, k := range keys {
		if _, ok := d.GetOk(k); ok {
			return true
		}
	}
	return false
}

func getDeviceIndexForDisk(disk *v3.VMDisk) (*int64, error) {
	if disk.DeviceProperties == nil {
		return nil, fmt.Errorf("deviceproperties was nil for disk")
	}
	if disk.DeviceProperties.DiskAddress == nil {
		return nil, fmt.Errorf("disk address was nil for disk")
	}
	if disk.DeviceProperties.DiskAddress.DeviceIndex == nil {
		return nil, fmt.Errorf("device index was nil for disk")
	}
	diskIndex := *disk.DeviceProperties.DiskAddress.DeviceIndex
	return &diskIndex, nil
}

func flattenDiskListFilterCloudInit(d *schema.ResourceData, disks []*v3.VMDisk) ([]map[string]interface{}, error) {
	//todo check if guestcust is passed -> if it is not passed, just continue without searching for cloud-init uuid
	// reason: no device_index or disk id will result in crash
	cloudInitCdromUUID := ""
	if cloudInitCdromUUIDInput, cliOk := d.GetOk("cloud_init_cdrom_uuid"); cliOk {
		cloudInitCdromUUID = cloudInitCdromUUIDInput.(string)
	}
	filteredDiskList := disks
	potentialCloudInitIDs := make([]string, 0)
	if cloudInitCdromUUID == "" && usesGuestCustomization(d) {
		filteredDiskList = make([]*v3.VMDisk, 0)
		//expand the user inputted list of disks
		expandedOrgDiskList := expandDiskList(d)
		//extract the CD-rom drives
		userCdromDiskList := GetCdromDiskList(expandedOrgDiskList)
		for _, eDisk := range disks {
			//if existing disk is not CD-rom, append it to the list and continue
			if !isCdromDisk(eDisk) {
				filteredDiskList = append(filteredDiskList, eDisk)
				continue
			} else {
				//Get existing CDrom device Index
				eDiskIndexP, err := getDeviceIndexForDisk(eDisk) //*eDisk.DeviceProperties.DiskAddress.DeviceIndex
				if err != nil {
					return nil, err
				}
				eDiskIndex := *eDiskIndexP
				match := false
				// Loop over the user defined cdrom drives
				for _, uDisk := range userCdromDiskList {
					//extract the device index of the user defined cdrom
					uDiskIndexP, err := getDeviceIndexForDisk(uDisk)
					if err != nil {
						return nil, err
					}
					uDiskIndex := *uDiskIndexP
					// if there is a matching device index for a userdefined and an existing cdrom, it is not the cloud-init one
					if eDiskIndex == uDiskIndex {
						filteredDiskList = append(filteredDiskList, eDisk)
						match = true
						break
					}
				}
				if !match {
					potentialCloudInitIDs = append(potentialCloudInitIDs, *eDisk.UUID)
				}
			}
		}
		if len(potentialCloudInitIDs) == 1 {
			cloudInitCdromUUID = potentialCloudInitIDs[0]
			d.Set("cloud_init_cdrom_uuid", cloudInitCdromUUID)
		}
		if len(potentialCloudInitIDs) > 1 {
			return nil, fmt.Errorf("more than 1 unknown cd-rom device: %v", potentialCloudInitIDs)
		}
	}
	fDiskList := flattenDiskListHelper(filteredDiskList, cloudInitCdromUUID)
	return fDiskList, nil
}
func flattenDiskList(disks []*v3.VMDisk) []map[string]interface{} {
	return flattenDiskListHelper(disks, "")
}
func flattenDiskListHelper(disks []*v3.VMDisk, cloudInitCdromUUID string) []map[string]interface{} {
	diskList := make([]map[string]interface{}, 0)
	for _, v := range disks {
		var deviceProps []map[string]interface{}
		var storageConfig []map[string]interface{}

		if v.DeviceProperties != nil {
			deviceProps = make([]map[string]interface{}, 1)
			index := fmt.Sprintf("%d", utils.Int64Value(v.DeviceProperties.DiskAddress.DeviceIndex))
			adapter := v.DeviceProperties.DiskAddress.AdapterType

			deviceProps[0] = map[string]interface{}{
				"device_type": v.DeviceProperties.DeviceType,
				"disk_address": map[string]interface{}{
					"device_index": index,
					"adapter_type": adapter,
				},
			}
		}

		if v.StorageConfig != nil {
			storageConfig = append(storageConfig, map[string]interface{}{
				"flash_mode": cast.ToString(v.StorageConfig.FlashMode),
				"storage_container_reference": []map[string]interface{}{
					{
						"url":  cast.ToString(v.StorageConfig.StorageContainerReference.URL),
						"kind": cast.ToString(v.StorageConfig.StorageContainerReference.Kind),
						"name": cast.ToString(v.StorageConfig.StorageContainerReference.Name),
						"uuid": cast.ToString(v.StorageConfig.StorageContainerReference.UUID),
					},
				},
			})
		}

		diskUUID := utils.StringValue(v.UUID)
		if cloudInitCdromUUID == diskUUID {
			continue
		}
		diskList = append(diskList, map[string]interface{}{
			"uuid":                   diskUUID,
			"disk_size_bytes":        utils.Int64Value(v.DiskSizeBytes),
			"disk_size_mib":          utils.Int64Value(v.DiskSizeMib),
			"device_properties":      deviceProps,
			"storage_config":         storageConfig,
			"data_source_reference":  flattenReferenceValues(v.DataSourceReference),
			"volume_group_reference": flattenReferenceValues(v.VolumeGroupReference),
		})
	}
	return diskList
}

func flattenSerialPortList(serialPorts []*v3.VMSerialPort) []map[string]interface{} {
	serialPortList := make([]map[string]interface{}, 0)
	if serialPorts != nil {
		serialPortList = make([]map[string]interface{}, len(serialPorts))
		for k, v := range serialPorts {
			serialPort := make(map[string]interface{})
			serialPort["index"] = utils.Int64Value(v.Index)
			serialPort["is_connected"] = utils.BoolValue(v.IsConnected)
			serialPortList[k] = serialPort
		}
	}
	return serialPortList
}

func flattenGPUList(gpu []*v3.VMGpuOutputStatus) []map[string]interface{} {
	gpuList := make([]map[string]interface{}, 0)
	if gpu != nil {
		gpuList = make([]map[string]interface{}, len(gpu))
		for k, v := range gpu {
			gpu := make(map[string]interface{})
			gpu["frame_buffer_size_mib"] = utils.Int64Value(v.FrameBufferSizeMib)
			gpu["vendor"] = utils.StringValue(v.Vendor)
			gpu["uuid"] = utils.StringValue(v.UUID)
			gpu["name"] = utils.StringValue(v.Name)
			gpu["pci_address"] = utils.StringValue(v.PCIAddress)
			gpu["fraction"] = utils.Int64Value(v.Fraction)
			gpu["mode"] = utils.StringValue(v.Mode)
			gpu["num_virtual_display_heads"] = utils.Int64Value(v.NumVirtualDisplayHeads)
			gpu["guest_driver_version"] = utils.StringValue(v.GuestDriverVersion)
			gpu["device_id"] = utils.Int64Value(v.DeviceID)
			gpuList[k] = gpu
		}
	}
	return gpuList
}

func flattenNutanixGuestTools(d *schema.ResourceData, guest *v3.GuestToolsStatus) error {
	nutanixGuestTools := make(map[string]interface{})
	ngtCredentials := make(map[string]string)
	ngtEnabledCapabilityList := make([]string, 0)

	if guest != nil && guest.NutanixGuestTools != nil {
		tools := guest.NutanixGuestTools
		ngtCredentials = tools.Credentials
		ngtEnabledCapabilityList = utils.StringValueSlice(tools.EnabledCapabilityList)

		nutanixGuestTools["available_version"] = utils.StringValue(tools.AvailableVersion)
		nutanixGuestTools["iso_mount_state"] = utils.StringValue(tools.IsoMountState)
		nutanixGuestTools["ngt_state"] = utils.StringValue(tools.NgtState)
		nutanixGuestTools["state"] = utils.StringValue(tools.State)
		nutanixGuestTools["version"] = utils.StringValue(tools.Version)
		nutanixGuestTools["guest_os_version"] = utils.StringValue(tools.GuestOsVersion)
		nutanixGuestTools["vss_snapshot_capable"] = strconv.FormatBool(utils.BoolValue(tools.VSSSnapshotCapable))
		nutanixGuestTools["is_reachable"] = strconv.FormatBool(utils.BoolValue(tools.IsReachable))
		nutanixGuestTools["vm_mobility_drivers_installed"] = strconv.FormatBool(utils.BoolValue(tools.VMMobilityDriversInstalled))
	}

	if err := d.Set("ngt_enabled_capability_list", ngtEnabledCapabilityList); err != nil {
		return err
	}

	if err := d.Set("ngt_credentials", ngtCredentials); err != nil {
		return err
	}

	if err := d.Set("nutanix_guest_tools", nutanixGuestTools); err != nil {
		return err
	}
	return nil
}
