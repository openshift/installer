package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceVirtualMachineScaleSetPacketCapture() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineScaleSetPacketCaptureCreate,
		Read:   resourceVirtualMachineScaleSetPacketCaptureRead,
		Delete: resourceVirtualMachineScaleSetPacketCaptureDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PacketCaptureID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"network_watcher_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.NetworkWatcherID,
			},

			"virtual_machine_scale_set_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					computeValidate.VirtualMachineScaleSetID,
				),
			},

			"maximum_bytes_per_packet": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},

			"maximum_bytes_per_session": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1073741824,
			},

			"maximum_capture_duration_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      18000,
				ValidateFunc: validation.IntBetween(1, 18000),
			},

			"storage_location": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"file_path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: networkValidate.FilePath,
							AtLeastOneOf: []string{"storage_location.0.file_path", "storage_location.0.storage_account_id"},
						},
						"storage_account_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: storageValidate.StorageAccountID,
							AtLeastOneOf: []string{"storage_location.0.file_path", "storage_location.0.storage_account_id"},
						},
						"storage_path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"filter": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"local_ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"local_port": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"protocol": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.PcProtocolAny),
								string(network.PcProtocolTCP),
								string(network.PcProtocolUDP),
							}, false),
						},
						"remote_ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"remote_port": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"machine_scope": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"exclude_instance_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"include_instance_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func resourceVirtualMachineScaleSetPacketCaptureCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCapturesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	watcherId, err := parse.NetworkWatcherID(d.Get("network_watcher_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewPacketCaptureID(subscriptionId, watcherId.ResourceGroup, watcherId.Name, d.Get("name").(string))

	targetResourceId := d.Get("virtual_machine_scale_set_id").(string)
	bytesToCapturePerPacket := d.Get("maximum_bytes_per_packet").(int)
	totalBytesPerSession := d.Get("maximum_bytes_per_session").(int)
	timeLimitInSeconds := d.Get("maximum_capture_duration_in_seconds").(int)

	existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set_packet_capture", id.ID())
	}

	storageLocation, err := expandVirtualMachineScaleSetPacketCaptureStorageLocation(d)
	if err != nil {
		return err
	}

	properties := network.PacketCapture{
		PacketCaptureParameters: &network.PacketCaptureParameters{
			Target:                  utils.String(targetResourceId),
			TargetType:              network.PacketCaptureTargetTypeAzureVMSS,
			StorageLocation:         storageLocation,
			BytesToCapturePerPacket: utils.Int64(int64(bytesToCapturePerPacket)),
			TimeLimitInSeconds:      utils.Int32(int32(timeLimitInSeconds)),
			TotalBytesPerSession:    utils.Int64(int64(totalBytesPerSession)),
			Filters:                 expandVirtualMachineScaleSetPacketCaptureFilters(d),
		},
	}

	if v, ok := d.GetOk("machine_scope"); ok {
		properties.PacketCaptureParameters.Scope = expandVirtualMachineScaleSetPacketCaptureMachineScope(v.([]interface{}))
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name, properties)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachineScaleSetPacketCaptureRead(d, meta)
}

func resourceVirtualMachineScaleSetPacketCaptureRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCapturesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PacketCaptureID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] %s not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.Name)

	networkWatcherId := parse.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroup, id.NetworkWatcherName)
	d.Set("network_watcher_id", networkWatcherId.ID())

	if props := resp.PacketCaptureResultProperties; props != nil {
		d.Set("virtual_machine_scale_set_id", props.Target)
		d.Set("maximum_bytes_per_packet", int(*props.BytesToCapturePerPacket))
		d.Set("maximum_bytes_per_session", int(*props.TotalBytesPerSession))
		d.Set("maximum_capture_duration_in_seconds", int(*props.TimeLimitInSeconds))

		location := flattenVirtualMachineScaleSetPacketCaptureStorageLocation(props.StorageLocation)
		if err := d.Set("storage_location", location); err != nil {
			return fmt.Errorf("setting `storage_location`: %+v", err)
		}

		filters := flattenVirtualMachineScaleSetPacketCaptureFilters(props.Filters)
		if err := d.Set("filter", filters); err != nil {
			return fmt.Errorf("setting `filter`: %+v", err)
		}

		scope, err := flattenVirtualMachineScaleSetPacketCaptureMachineScope(props.Scope)
		if err != nil {
			return err
		}
		if err := d.Set("machine_scope", scope); err != nil {
			return fmt.Errorf(`setting "machine_scope": %+v`, err)
		}
	}

	return nil
}

func resourceVirtualMachineScaleSetPacketCaptureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCapturesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PacketCaptureID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkWatcherName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandVirtualMachineScaleSetPacketCaptureStorageLocation(d *pluginsdk.ResourceData) (*network.PacketCaptureStorageLocation, error) {
	locations := d.Get("storage_location").([]interface{})
	if len(locations) == 0 {
		return nil, fmt.Errorf("expandng `storage_location`: not found")
	}

	location := locations[0].(map[string]interface{})

	storageLocation := network.PacketCaptureStorageLocation{}

	if v := location["file_path"]; v != "" {
		storageLocation.FilePath = utils.String(v.(string))
	}
	if v := location["storage_account_id"]; v != "" {
		storageLocation.StorageID = utils.String(v.(string))
	}

	return &storageLocation, nil
}

func flattenVirtualMachineScaleSetPacketCaptureStorageLocation(input *network.PacketCaptureStorageLocation) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make(map[string]interface{})

	if path := input.FilePath; path != nil {
		output["file_path"] = *path
	}

	if account := input.StorageID; account != nil {
		output["storage_account_id"] = *account
	}

	if path := input.StoragePath; path != nil {
		output["storage_path"] = *path
	}

	return []interface{}{output}
}

func expandVirtualMachineScaleSetPacketCaptureFilters(d *pluginsdk.ResourceData) *[]network.PacketCaptureFilter {
	inputFilters := d.Get("filter").([]interface{})
	if len(inputFilters) == 0 {
		return nil
	}

	filters := make([]network.PacketCaptureFilter, 0)

	for _, v := range inputFilters {
		inputFilter := v.(map[string]interface{})

		localIPAddress := inputFilter["local_ip_address"].(string)
		localPort := inputFilter["local_port"].(string) // TODO: should this be an int?
		protocol := inputFilter["protocol"].(string)
		remoteIPAddress := inputFilter["remote_ip_address"].(string)
		remotePort := inputFilter["remote_port"].(string)

		filter := network.PacketCaptureFilter{
			LocalIPAddress:  utils.String(localIPAddress),
			LocalPort:       utils.String(localPort),
			Protocol:        network.PcProtocol(protocol),
			RemoteIPAddress: utils.String(remoteIPAddress),
			RemotePort:      utils.String(remotePort),
		}
		filters = append(filters, filter)
	}

	return &filters
}

func flattenVirtualMachineScaleSetPacketCaptureFilters(input *[]network.PacketCaptureFilter) []interface{} {
	filters := make([]interface{}, 0)

	if inFilter := input; inFilter != nil {
		for _, v := range *inFilter {
			filter := make(map[string]interface{})

			if address := v.LocalIPAddress; address != nil {
				filter["local_ip_address"] = *address
			}

			if port := v.LocalPort; port != nil {
				filter["local_port"] = *port
			}

			filter["protocol"] = string(v.Protocol)

			if address := v.RemoteIPAddress; address != nil {
				filter["remote_ip_address"] = *address
			}

			if port := v.RemotePort; port != nil {
				filter["remote_port"] = *port
			}

			filters = append(filters, filter)
		}
	}

	return filters
}

func expandVirtualMachineScaleSetPacketCaptureMachineScope(input []interface{}) *network.PacketCaptureMachineScope {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &network.PacketCaptureMachineScope{}

	if exclude := raw["exclude_instance_ids"].([]interface{}); len(exclude) > 0 {
		output.Exclude = utils.ExpandStringSlice(exclude)
	}

	if include := raw["include_instance_ids"].([]interface{}); len(include) > 0 {
		output.Include = utils.ExpandStringSlice(include)
	}

	return output
}

func flattenVirtualMachineScaleSetPacketCaptureMachineScope(input *network.PacketCaptureMachineScope) ([]interface{}, error) {
	outputs := make([]interface{}, 0)
	if input == nil || (input.Exclude == nil && input.Include == nil) || (len(*input.Exclude) == 0 && len(*input.Include) == 0) {
		return outputs, nil
	}

	output := make(map[string]interface{}, 0)

	excludedInstanceIds, err := flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input.Exclude)
	if err != nil {
		return nil, err
	}
	output["exclude_instance_ids"] = excludedInstanceIds

	includedInstanceIds, err := flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input.Include)
	if err != nil {
		return nil, err
	}
	output["include_instance_ids"] = includedInstanceIds

	outputs = append(outputs, output)

	return outputs, nil
}

func flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input *[]string) ([]string, error) {
	instances := make([]string, 0)
	if input == nil {
		return instances, nil
	}

	for _, instance := range *input {
		instance, err := computeParse.VMSSInstanceID(instance)
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance.VirtualMachineName)
	}

	return instances, nil
}
