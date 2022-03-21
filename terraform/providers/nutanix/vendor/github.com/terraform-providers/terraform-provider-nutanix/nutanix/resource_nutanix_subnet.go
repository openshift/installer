package nutanix

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

var (
	subnetTimeout    = 10 * time.Minute
	subnetDelay      = 10 * time.Second
	subnetMinTimeout = 3 * time.Second
)

func resourceNutanixSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceNutanixSubnetCreate,
		Read:   resourceNutanixSubnetRead,
		Update: resourceNutanixSubnetUpdate,
		Delete: resourceNutanixSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNutanixSubnetInstanceResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceSubnetInstanceStateUpgradeV0,
				Version: 0,
			},
		},
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"categories": categoriesSchema(),
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prefix_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"subnet_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_server_address": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dhcp_server_address_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ip_config_pool_list_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(
							"^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)["+
								" ](?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"),
						"please see https://developer.nutanix.com/reference/prism_central/v3/#definitions-ip_pool"),
				},
			},
			"dhcp_options": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot_file_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tftp_server_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dhcp_domain_name_server_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dhcp_domain_search_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"network_function_chain_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNutanixSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	request := &v3.SubnetIntentInput{}
	spec := &v3.Subnet{}
	metadata := &v3.Metadata{}
	subnet := &v3.SubnetResources{}

	n, nok := d.GetOk("name")
	azr, azrok := d.GetOk("availability_zone_reference")
	clusterUUID, crok := d.GetOk("cluster_uuid")
	_, stok := d.GetOk("subnet_type")

	if !stok && !nok {
		return fmt.Errorf("please provide the required attributes name, subnet_type")
	}

	if !nok {
		return fmt.Errorf("please provide the required name attribute")
	}
	if err := getMetadataAttributes(d, metadata, "subnet"); err != nil {
		return err
	}

	if azrok {
		a := azr.(map[string]interface{})
		spec.AvailabilityZoneReference = validateRef(a)
	}
	if crok {
		spec.ClusterReference = buildReference(clusterUUID.(string), "cluster")
	}

	getSubnetResources(d, subnet)

	spec.Description = utils.StringPtr(d.Get("description").(string))

	spec.Name = utils.StringPtr(n.(string))
	spec.Resources = subnet
	request.Metadata = metadata
	request.Spec = spec

	resp, err := conn.V3.CreateSubnet(request)
	if err != nil {
		return fmt.Errorf("error creating Nutanix Subnet %s: %+v", utils.StringValue(spec.Name), err)
	}

	d.SetId(*resp.Metadata.UUID)

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the Subnet to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		id := d.Id()
		d.SetId("")
		return fmt.Errorf("error waiting for subnet id (%s) to create: %+v", id, err)
	}

	// Setting Description because in Get request is not present.
	d.Set("description", utils.StringValue(resp.Spec.Description))

	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	id := d.Id()
	resp, err := conn.V3.GetSubnet(id)
	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		errDel := resourceNutanixSubnetDelete(d, meta)
		if errDel != nil {
			return fmt.Errorf("error deleting subnet (%s) after read error: %+v", id, errDel)
		}
		d.SetId("")
		return fmt.Errorf("error reading subnet id (%s): %+v", id, err)
	}

	m, c := setRSEntityMetadata(resp.Metadata)

	if err := d.Set("metadata", m); err != nil {
		return err
	}
	if err := d.Set("categories", c); err != nil {
		return err
	}
	if err := d.Set("project_reference", flattenReferenceValues(resp.Metadata.ProjectReference)); err != nil {
		return err
	}
	if err := d.Set("owner_reference", flattenReferenceValues(resp.Metadata.OwnerReference)); err != nil {
		return err
	}
	if err := d.Set("availability_zone_reference", flattenReferenceValues(resp.Status.AvailabilityZoneReference)); err != nil {
		return err
	}
	if err := flattenClusterReference(resp.Status.ClusterReference, d); err != nil {
		return err
	}

	dgIP := ""
	sIP := ""
	pl := int64(0)
	port := int64(0)
	dhcpSA := make(map[string]interface{})
	dOptions := make(map[string]interface{})
	ipcpl := make([]string, 0)
	dnsList := make([]string, 0)
	dsList := make([]string, 0)

	if resp.Status.Resources.IPConfig != nil {
		dgIP = utils.StringValue(resp.Status.Resources.IPConfig.DefaultGatewayIP)
		pl = utils.Int64Value(resp.Status.Resources.IPConfig.PrefixLength)
		sIP = utils.StringValue(resp.Status.Resources.IPConfig.SubnetIP)

		if dhcpServerAdd := resp.Status.Resources.IPConfig.DHCPServerAddress; dhcpServerAdd != nil {
			if dhcpServerAdd.IP != nil {
				dhcpSA["ip"] = utils.StringValue(dhcpServerAdd.IP)
			}
			if dhcpServerAdd.FQDN != nil {
				dhcpSA["fqdn"] = utils.StringValue(dhcpServerAdd.FQDN)
			}
			if dhcpServerAdd.IPV6 != nil {
				dhcpSA["ipv6"] = utils.StringValue(dhcpServerAdd.IPV6)
			}
			port = utils.Int64Value(dhcpServerAdd.Port)
		}

		if resp.Status.Resources.IPConfig.PoolList != nil {
			pl := resp.Status.Resources.IPConfig.PoolList
			poolList := make([]string, len(pl))
			for k, v := range pl {
				poolList[k] = utils.StringValue(v.Range)
			}
			ipcpl = poolList
		}
		if resp.Status.Resources.IPConfig.DHCPOptions != nil {
			dOptions["boot_file_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.BootFileName)
			dOptions["domain_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.DomainName)
			dOptions["tftp_server_name"] = utils.StringValue(resp.Status.Resources.IPConfig.DHCPOptions.TFTPServerName)

			if resp.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList != nil {
				dnsList = utils.StringValueSlice(resp.Status.Resources.IPConfig.DHCPOptions.DomainNameServerList)
			}
			if resp.Status.Resources.IPConfig.DHCPOptions.DomainSearchList != nil {
				dsList = utils.StringValueSlice(resp.Status.Resources.IPConfig.DHCPOptions.DomainSearchList)
			}
		}
	}

	if err := d.Set("dhcp_server_address", dhcpSA); err != nil {
		return fmt.Errorf("error setting attribute for subnet id (%s) dhcp_server_address: %s", d.Id(), err)
	}
	if err := d.Set("ip_config_pool_list_ranges", ipcpl); err != nil {
		return fmt.Errorf("error setting attribute for subnet id (%s) ip_config_pool_list_ranges: %s", d.Id(), err)
	}
	if err := d.Set("dhcp_options", dOptions); err != nil {
		return fmt.Errorf("error setting attribute for subnet id (%s) dhcp_options: %s", d.Id(), err)
	}
	if err := d.Set("dhcp_domain_name_server_list", dnsList); err != nil {
		return fmt.Errorf("error setting attribute for subnet id (%s) dhcp_domain_name_server_list: %s", d.Id(), err)
	}
	if err := d.Set("dhcp_domain_search_list", dsList); err != nil {
		return fmt.Errorf("error setting attribute for subnet id (%s) dhcp_domain_search_list: %s", d.Id(), err)
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))

	nfcr := make(map[string]interface{})

	if status := resp.Status; status != nil {
		d.Set("name", utils.StringValue(status.Name))
		d.Set("state", utils.StringValue(status.State))

		if res := status.Resources; res != nil {
			d.Set("vswitch_name", utils.StringValue(res.VswitchName))
			d.Set("subnet_type", utils.StringValue(res.SubnetType))

			d.Set("vlan_id", utils.Int64Value(res.VlanID))

			if res.NetworkFunctionChainReference != nil {
				nfcr = flattenReferenceValues(res.NetworkFunctionChainReference)
			}
		}
	}

	d.Set("network_function_chain_reference", nfcr)
	d.Set("default_gateway_ip", dgIP)
	d.Set("prefix_length", pl)
	d.Set("subnet_ip", sIP)
	d.Set("dhcp_server_address_port", port)

	d.SetId(*resp.Metadata.UUID)

	return nil
}

func resourceNutanixSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API
	request := &v3.SubnetIntentInput{}
	metadata := &v3.Metadata{}
	res := &v3.SubnetResources{}
	ipcfg := &v3.IPConfig{}
	dhcpO := &v3.DHCPOptions{}
	spec := &v3.Subnet{}

	id := d.Id()
	response, err := conn.V3.GetSubnet(id)

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
		}
		return fmt.Errorf("error retrieving for subnet id (%s) :%+v", id, err)
	}

	if response.Metadata != nil {
		metadata = response.Metadata
	}

	if response.Spec != nil {
		spec = response.Spec

		if response.Spec.Resources != nil {
			res = response.Spec.Resources

			if res.IPConfig != nil {
				ipcfg = res.IPConfig
			}
			if ipcfg != nil {
				dhcpO = ipcfg.DHCPOptions
			}
		}
	}

	if d.HasChange("categories") {
		metadata.Categories = expandCategories(d.Get("categories"))
	}
	if d.HasChange("owner_reference") {
		or := d.Get("owner_reference").(map[string]interface{})
		metadata.OwnerReference = validateRef(or)
	}
	if d.HasChange("project_reference") {
		pr := d.Get("project_reference").(map[string]interface{})
		metadata.ProjectReference = validateRef(pr)
	}
	if d.HasChange("name") {
		spec.Name = utils.StringPtr(d.Get("name").(string))
	}
	if d.HasChange("description") {
		spec.Description = utils.StringPtr(d.Get("description").(string))
	}
	if d.HasChange("availability_zone_reference") {
		a := d.Get("availability_zone_reference").(map[string]interface{})
		spec.AvailabilityZoneReference = validateRef(a)
	}
	if d.HasChange("cluster_uuid") {
		uuid := d.Get("cluster_uuid").(string)
		spec.ClusterReference = buildReference(uuid, "cluster")
	}
	if d.HasChange("dhcp_domain_name_server_list") {
		dhcpO.DomainNameServerList = expandStringList(d.Get("dhcp_domain_name_server_list").([]interface{}))
	}
	if d.HasChange("dhcp_domain_search_list") {
		dhcpO.DomainSearchList = expandStringList(d.Get("dhcp_domain_search_list").([]interface{}))
	}
	if d.HasChange("ip_config_pool_list_ranges") {
		dd := d.Get("ip_config_pool_list_ranges").([]interface{})
		ddn := make([]*v3.IPPool, len(dd))
		for k, v := range dd {
			i := &v3.IPPool{}
			i.Range = utils.StringPtr(v.(string))
			ddn[k] = i
		}
		ipcfg.PoolList = ddn
	}
	if d.HasChange("dhcp_options") {
		dOptions := d.Get("dhcp_options").(map[string]interface{})

		dhcpO.BootFileName = validateMapStringValue(dOptions, "boot_file_name")
		dhcpO.DomainName = validateMapStringValue(dOptions, "domain_name")
		dhcpO.TFTPServerName = validateMapStringValue(dOptions, "tftp_server_name")
	}
	if d.HasChange("network_function_chain_reference") {
		res.NetworkFunctionChainReference =
			validateRef(d.Get("network_function_chain_reference").(map[string]interface{}))
	}
	if d.HasChange("vswitch_name") {
		res.VswitchName = utils.StringPtr(d.Get("vswitch_name").(string))
	}
	if d.HasChange("subnet_type") {
		res.SubnetType = utils.StringPtr(d.Get("subnet_type").(string))
	}
	if d.HasChange("default_gateway_ip") {
		ipcfg.DefaultGatewayIP = utils.StringPtr(d.Get("default_gateway_ip").(string))
	}
	if d.HasChange("prefix_length") {
		ipcfg.PrefixLength = utils.Int64Ptr(int64(d.Get("prefix_length").(int)))
	}
	if d.HasChange("subnet_ip") {
		ipcfg.SubnetIP = utils.StringPtr(d.Get("subnet_ip").(string))
	}
	if d.HasChange("dhcp_server_address") {
		dh := d.Get("dhcp_server_address").(map[string]interface{})

		ipcfg.DHCPServerAddress = &v3.Address{
			IP:   validateMapStringValue(dh, "ip"),
			IPV6: validateMapStringValue(dh, "ipv6"),
			FQDN: validateMapStringValue(dh, "fqdn"),
		}
	}
	if d.HasChange("dhcp_server_address_port") {
		ipcfg.DHCPServerAddress.Port = utils.Int64Ptr(int64(d.Get("dhcp_server_address_port").(int)))
	}
	if d.HasChange("vlan_id") {
		res.VlanID = utils.Int64Ptr(int64(d.Get("vlan_id").(int)))
	}

	ipcfg.DHCPOptions = dhcpO
	res.IPConfig = ipcfg
	spec.Resources = res
	request.Metadata = metadata
	request.Spec = spec

	resp, errUpdate := conn.V3.UpdateSubnet(d.Id(), request)
	if errUpdate != nil {
		return fmt.Errorf("error updating subnet id %s): %s", d.Id(), errUpdate)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for subnet (%s) to update: %s", d.Id(), err)
	}
	// Setting Description because in Get request is not present.
	d.Set("description", utils.StringValue(resp.Spec.Description))

	return resourceNutanixSubnetRead(d, meta)
}

func resourceNutanixSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).API

	resp, err := conn.V3.DeleteSubnet(d.Id())

	if err != nil {
		return fmt.Errorf("error deleting subnet id %s): %s", d.Id(), err)
	}

	taskUUID := resp.Status.ExecutionContext.TaskUUID.(string)

	// Wait for the VM to be available
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"QUEUED", "RUNNING"},
		Target:     []string{"SUCCEEDED"},
		Refresh:    taskStateRefreshFunc(conn, taskUUID),
		Timeout:    subnetTimeout,
		Delay:      subnetDelay,
		MinTimeout: subnetMinTimeout,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for subnet (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceNutanixSubnetExists(conn *v3.Client, name string) (*string, error) {
	var subnetUUID *string

	filter := fmt.Sprintf("name==%s", name)
	subnetList, err := conn.V3.ListAllSubnet(filter)

	if err != nil {
		return nil, err
	}

	for _, subnet := range subnetList.Entities {
		if utils.StringValue(subnet.Status.Name) == name {
			subnetUUID = subnet.Metadata.UUID
		}
	}
	return subnetUUID, nil
}

func getSubnetResources(d *schema.ResourceData, subnet *v3.SubnetResources) {
	ip := &v3.IPConfig{}
	dhcpo := &v3.DHCPOptions{}

	if v, ok := d.GetOk("vswitch_name"); ok {
		subnet.VswitchName = utils.StringPtr(v.(string))
	}
	if st, ok := d.GetOk("subnet_type"); ok {
		subnet.SubnetType = utils.StringPtr(st.(string))
	}
	if v, ok := d.GetOk("default_gateway_ip"); ok {
		ip.DefaultGatewayIP = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("prefix_length"); ok {
		ip.PrefixLength = utils.Int64Ptr(int64(v.(int)))
	}
	if v, ok := d.GetOk("subnet_ip"); ok {
		ip.SubnetIP = utils.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("dhcp_server_address"); ok {
		dhcpa := v.(map[string]interface{})
		address := &v3.Address{}

		if ip, ok := dhcpa["ip"]; ok {
			address.IP = utils.StringPtr(ip.(string))
		}
		if fqdn, ok := dhcpa["fqdn"]; ok {
			address.FQDN = utils.StringPtr(fqdn.(string))
		}
		if v, ok := d.GetOk("dhcp_server_address_port"); ok {
			address.Port = utils.Int64Ptr(int64(v.(int)))
		}
		if ipv6, ok := dhcpa["ipv6"]; ok {
			address.IPV6 = utils.StringPtr(ipv6.(string))
		}

		ip.DHCPServerAddress = address
	}
	if v, ok := d.GetOk("ip_config_pool_list_ranges"); ok {
		p := v.([]interface{})
		pool := make([]*v3.IPPool, len(p))

		for k, v := range p {
			pItem := &v3.IPPool{}
			pItem.Range = utils.StringPtr(v.(string))
			pool[k] = pItem
		}

		ip.PoolList = pool
	}
	if v, ok := d.GetOk("dhcp_options"); ok {
		dop := v.(map[string]interface{})

		if boot, ok := dop["boot_file_name"]; ok {
			dhcpo.BootFileName = utils.StringPtr(boot.(string))
		}

		if dn, ok := dop["domain_name"]; ok {
			dhcpo.DomainName = utils.StringPtr(dn.(string))
		}

		if tsn, ok := dop["tftp_server_name"]; ok {
			dhcpo.TFTPServerName = utils.StringPtr(tsn.(string))
		}
	}

	if v, ok := d.GetOk("dhcp_domain_name_server_list"); ok {
		dhcpo.DomainNameServerList = expandStringList(v.([]interface{}))
	}
	if v, ok := d.GetOk("dhcp_domain_search_list"); ok {
		dhcpo.DomainSearchList = expandStringList(v.([]interface{}))
	}

	v, ok := d.GetOk("vlan_id")
	if v.(int) == 0 || ok {
		subnet.VlanID = utils.Int64Ptr(int64(v.(int)))
	}

	if v, ok := d.GetOk("network_function_chain_reference"); ok {
		subnet.NetworkFunctionChainReference = validateRef(v.(map[string]interface{}))
	}

	ip.DHCPOptions = dhcpo

	subnet.IPConfig = ip
}

func resourceSubnetInstanceStateUpgradeV0(is map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Entering resourceSubnetInstanceStateUpgradeV0")
	return resourceNutanixCategoriesMigrateState(is, meta)
}

func resourceNutanixSubnetInstanceResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_update_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"spec_hash": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"categories": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"owner_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cluster_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vswitch_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_gateway_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prefix_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"subnet_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_server_address": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dhcp_server_address_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ip_config_pool_list_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(
							"^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)["+
								" ](?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"),
						"please see https://developer.nutanix.com/reference/prism_central/v3/#definitions-ip_pool"),
				},
			},
			"dhcp_options": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boot_file_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"tftp_server_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dhcp_domain_name_server_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dhcp_domain_search_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"network_function_chain_reference": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
