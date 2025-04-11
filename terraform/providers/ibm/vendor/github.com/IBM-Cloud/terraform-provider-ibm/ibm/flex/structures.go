// Copyright IBM Corp. 2017, 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package flex

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/api/icd/icdv4"
	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/api/schematics"
	"github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go-config/v2/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/IBM/platform-services-go-sdk/globalsearchv2"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	rc "github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	rg "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

const (
	prodBaseController  = "https://cloud.ibm.com"
	stageBaseController = "https://test.cloud.ibm.com"
	//ResourceControllerURL ...
	ResourceControllerURL = "resource_controller_url"
	//ResourceName ...
	ResourceName = "resource_name"
	//ResourceCRN ...
	ResourceCRN = "resource_crn"
	//ResourceStatus ...
	ResourceStatus = "resource_status"
	//ResourceGroupName ...
	ResourceGroupName = "resource_group_name"
	//DeletionProtection ...
	DeletionProtection = "deletion_protection"
	//RelatedCRN ...
	RelatedCRN                                = "related_crn"
	SystemIBMLabelPrefix                      = "ibm-cloud.kubernetes.io/"
	KubernetesLabelPrefix                     = "kubernetes.io/"
	K8sLabelPrefix                            = "k8s.io/"
	isLBListenerPolicyAction                  = "action"
	isLBListenerPolicyTargetID                = "target_id"
	isLBListenerPolicyTargetURL               = "target_url"
	isLBListenerPolicyTargetHTTPStatusCode    = "target_http_status_code"
	isLBListenerPolicyHTTPSRedirectStatusCode = "target_https_redirect_status_code"
	isLBListenerPolicyHTTPSRedirectURI        = "target_https_redirect_uri"
	isLBListenerPolicyHTTPSRedirectListener   = "target_https_redirect_listener"
	isLBPoolSessPersistenceType               = "session_persistence_type"
	isLBPoolSessPersistenceAppCookieName      = "session_persistence_app_cookie_name"
	isLBProfile                               = "profile"
	isLBRouteMode                             = "route_mode"
	isLBType                                  = "type"
	crnSeparator                              = ":"
	scopeSeparator                            = "/"
	crn                                       = "crn"
)

var (
	ErrMalformedCRN   = errors.New("malformed CRN")
	ErrMalformedScope = errors.New("malformed scope in CRN")
)

// HashInt ...
func HashInt(v interface{}) int { return v.(int) }

func ExpandStringList(input []interface{}) []string {
	vs := make([]string, len(input))
	for i, v := range input {
		vs[i] = v.(string)
	}
	return vs
}

func FlattenStringList(list []string) []interface{} {
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v
	}
	return vs
}

func ExpandIntList(input []interface{}) []int {
	vs := make([]int, len(input))
	for i, v := range input {
		vs[i] = v.(int)
	}
	return vs
}

func FlattenIntList(list []int) []interface{} {
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v
	}
	return vs
}

func ExpandInt64List(input []interface{}) []int64 {
	vs := make([]int64, len(input))
	for i, v := range input {
		vs[i] = v.(int64)
	}
	return vs
}

func FlattenInt64List(list []int64) []interface{} {
	vs := make([]interface{}, len(list))
	for i, v := range list {
		vs[i] = v
	}
	return vs
}

func NewStringSet(f schema.SchemaSetFunc, in []string) *schema.Set {
	var out = make([]interface{}, len(in), len(in))
	for i, v := range in {
		out[i] = v
	}
	return schema.NewSet(f, out)
}

func FlattenRoute(in []mccpv2.Route) *schema.Set {
	vs := make([]string, len(in))
	for i, v := range in {
		vs[i] = v.GUID
	}
	return NewStringSet(schema.HashString, vs)
}

func stringSliceToSet(in []string) *schema.Set {
	vs := make([]string, len(in))
	for i, v := range in {
		vs[i] = v
	}
	return NewStringSet(schema.HashString, vs)
}

func FlattenServiceBindings(in []mccpv2.ServiceBinding) *schema.Set {
	vs := make([]string, len(in))
	for i, v := range in {
		vs[i] = v.ServiceInstanceGUID
	}
	return NewStringSet(schema.HashString, vs)
}

func flattenPort(in []int) *schema.Set {
	var out = make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return schema.NewSet(HashInt, out)
}

func FlattenFileStorageID(in []datatypes.Network_Storage) *schema.Set {
	var out = []interface{}{}
	for _, v := range in {
		if *v.NasType == "NAS" {
			out = append(out, *v.Id)
		}
	}
	return schema.NewSet(HashInt, out)
}

func FlattenBlockStorageID(in []datatypes.Network_Storage) *schema.Set {
	var out = []interface{}{}
	for _, v := range in {
		if *v.NasType == "ISCSI" {
			out = append(out, *v.Id)
		}
	}
	return schema.NewSet(HashInt, out)
}

func FlattenSSHKeyIDs(in []datatypes.Security_Ssh_Key) *schema.Set {
	var out = []interface{}{}
	for _, v := range in {
		out = append(out, *v.Id)
	}
	return schema.NewSet(HashInt, out)
}

func FlattenSpaceRoleUsers(in []mccpv2.SpaceRole) *schema.Set {
	var out = []interface{}{}
	for _, v := range in {
		out = append(out, v.UserName)
	}
	return schema.NewSet(schema.HashString, out)
}

func FlattenOrgRole(in []mccpv2.OrgRole, excludeUsername string) *schema.Set {
	var out = []interface{}{}
	for _, v := range in {
		if excludeUsername == "" {
			out = append(out, v.UserName)
		} else {
			if v.UserName != excludeUsername {
				out = append(out, v.UserName)
			}
		}
	}
	return schema.NewSet(schema.HashString, out)
}

func flattenMapInterfaceVal(m map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

func flattenCredentials(creds map[string]interface{}) map[string]string {
	return flattenMapInterfaceVal(creds)
}

func flattenServiceKeyCredentials(creds map[string]interface{}) map[string]string {
	return flattenCredentials(creds)
}

func FlattenServiceInstanceCredentials(keys []mccpv2.ServiceKeyFields) []interface{} {
	var out = make([]interface{}, len(keys), len(keys))
	for i, k := range keys {
		m := make(map[string]interface{})
		m["name"] = k.Entity.Name
		m["credentials"] = Flatten(k.Entity.Credentials)
		out[i] = m
	}
	return out
}

func FlattenUsersSet(userList *schema.Set) []string {
	users := make([]string, 0)
	for _, user := range userList.List() {
		users = append(users, user.(string))
	}
	return users
}

func FlattenSet(set *schema.Set) []string {
	setList := set.List()
	elems := make([]string, 0, len(setList))
	for _, elem := range setList {
		elems = append(elems, elem.(string))
	}
	return elems
}

func ExpandMembers(configured []interface{}) []datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo {
	members := make([]datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo, 0, len(configured))
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})
		p := &datatypes.Network_LBaaS_LoadBalancerServerInstanceInfo{}
		if v, ok := data["private_ip_address"]; ok && v.(string) != "" {
			p.PrivateIpAddress = sl.String(v.(string))
		}
		if v, ok := data["weight"]; ok && v.(int) != 0 {
			p.Weight = sl.Int(v.(int))
		}

		members = append(members, *p)
	}
	return members
}

func FlattenServerInstances(list []datatypes.Network_LBaaS_Member) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"private_ip_address": *i.Address,
			"member_id":          *i.Uuid,
		}
		if i.Weight != nil {
			l["weight"] = *i.Weight
		}
		result = append(result, l)
	}
	return result
}

func FlattenProtocols(list []datatypes.Network_LBaaS_Listener) []map[string]interface{} {
	var lbIdToMethod = make(map[string]string)
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"frontend_protocol":     *i.Protocol,
			"frontend_port":         *i.ProtocolPort,
			"backend_protocol":      *i.DefaultPool.Protocol,
			"backend_port":          *i.DefaultPool.ProtocolPort,
			"load_balancing_method": lbIdToMethod[*i.DefaultPool.LoadBalancingAlgorithm],
			"protocol_id":           *i.Uuid,
		}
		if i.DefaultPool.SessionAffinity != nil && i.DefaultPool.SessionAffinity.Type != nil && *i.DefaultPool.SessionAffinity.Type != "" {
			l["session_stickiness"] = *i.DefaultPool.SessionAffinity.Type
		}
		if i.ConnectionLimit != nil && *i.ConnectionLimit != 0 {
			l["max_conn"] = *i.ConnectionLimit
		}
		if i.TlsCertificateId != nil && *i.TlsCertificateId != 0 {
			l["tls_certificate_id"] = *i.TlsCertificateId
		}
		result = append(result, l)
	}
	return result
}

func FlattenVpcWorkerPoolSecondaryDisk(secondaryDisk containerv2.DiskConfigResp) []map[string]interface{} {
	storageList := make([]map[string]interface{}, 1)
	secondary_storage := map[string]interface{}{
		"name":               secondaryDisk.Name,
		"count":              secondaryDisk.Count,
		"size":               secondaryDisk.Size,
		"device_type":        secondaryDisk.DeviceType,
		"raid_configuration": secondaryDisk.RAIDConfiguration,
		"profile":            secondaryDisk.Profile,
	}
	storageList[0] = secondary_storage
	return storageList
}
func FlattenVpcWorkerPools(list []containerv2.GetWorkerPoolResponse) []map[string]interface{} {
	workerPools := make([]map[string]interface{}, len(list))
	for i, workerPool := range list {
		l := map[string]interface{}{
			"id":               workerPool.ID,
			"name":             workerPool.PoolName,
			"flavor":           workerPool.Flavor,
			"worker_count":     workerPool.WorkerCount,
			"isolation":        workerPool.Isolation,
			"labels":           workerPool.Labels,
			"operating_system": workerPool.OperatingSystem,
			"state":            workerPool.Lifecycle.ActualState,
			"host_pool_id":     workerPool.HostPoolID,
		}
		zones := workerPool.Zones
		zonesConfig := make([]map[string]interface{}, len(zones))
		for j, zone := range zones {
			z := map[string]interface{}{
				"zone":         zone.ID,
				"worker_count": zone.WorkerCount,
			}
			subnets := zone.Subnets
			subnetConfig := make([]map[string]interface{}, len(subnets))
			for k, subnet := range subnets {
				s := map[string]interface{}{
					"id":      subnet.ID,
					"primary": subnet.Primary,
				}
				subnetConfig[k] = s
			}
			z["subnets"] = subnetConfig
			zonesConfig[j] = z
		}
		l["zones"] = zonesConfig
		if workerPool.SecondaryStorageOption != nil {
			l["secondary_storage"] = FlattenVpcWorkerPoolSecondaryDisk(*workerPool.SecondaryStorageOption)
		}
		workerPools[i] = l
	}

	return workerPools
}

func flattenVpcZones(list []containerv2.ZoneResp) []map[string]interface{} {
	zones := make([]map[string]interface{}, len(list))
	for i, zone := range list {
		l := map[string]interface{}{
			"id":           zone.ID,
			"subnet_id":    FlattenSubnets(zone.Subnets),
			"worker_count": zone.WorkerCount,
		}
		zones[i] = l
	}
	return zones
}
func FlattenConditions(list []iamaccessgroupsv2.RuleConditions) []map[string]interface{} {
	conditions := make([]map[string]interface{}, len(list))
	for i, cond := range list {
		l := map[string]interface{}{
			"claim":    cond.Claim,
			"operator": cond.Operator,
			"value":    strings.ReplaceAll(*cond.Value, "\"", ""),
		}
		conditions[i] = l
	}
	return conditions
}
func FlattenAccessGroupRules(list *iamaccessgroupsv2.RulesList) []map[string]interface{} {
	rules := make([]map[string]interface{}, len(list.Rules))
	for i, item := range list.Rules {
		l := map[string]interface{}{
			"name":              item.Name,
			"expiration":        item.Expiration,
			"identity_provider": item.RealmName,
			"conditions":        FlattenConditions(item.Conditions),
		}
		rules[i] = l
	}
	return rules
}

func FlattenSubnets(list []containerv2.Subnet) []map[string]interface{} {
	subs := make([]map[string]interface{}, len(list))
	for i, sub := range list {
		l := map[string]interface{}{
			"id":           sub.ID,
			"worker_count": sub.Primary,
		}
		subs[i] = l
	}
	return subs
}

func FlattenZones(list []containerv1.WorkerPoolZoneResponse) []map[string]interface{} {
	zones := make([]map[string]interface{}, len(list))
	for i, zone := range list {
		l := map[string]interface{}{
			"zone":         zone.WorkerPoolZone.ID,
			"private_vlan": zone.WorkerPoolZone.WorkerPoolZoneNetwork.PrivateVLAN,
			"public_vlan":  zone.WorkerPoolZone.WorkerPoolZoneNetwork.PublicVLAN,
			"worker_count": zone.WorkerCount,
		}
		zones[i] = l
	}
	return zones
}

func FlattenZonesv2(list []containerv2.ZoneResp) []map[string]interface{} {
	zones := make([]map[string]interface{}, len(list))
	for i, zone := range list {
		l := map[string]interface{}{
			"zone":         zone.ID,
			"subnets":      zone.Subnets,
			"worker_count": zone.WorkerCount,
		}
		zones[i] = l
	}
	return zones
}

func FlattenWorkerPools(list []containerv1.WorkerPoolResponse) []map[string]interface{} {
	workerPools := make([]map[string]interface{}, len(list))
	for i, workerPool := range list {
		l := map[string]interface{}{
			"id":            workerPool.ID,
			"hardware":      workerPool.Isolation,
			"name":          workerPool.Name,
			"machine_type":  workerPool.MachineType,
			"size_per_zone": workerPool.Size,
			"state":         workerPool.State,
			"labels":        workerPool.Labels,
		}
		zones := workerPool.Zones
		zonesConfig := make([]map[string]interface{}, len(zones))
		for j, zone := range zones {
			z := map[string]interface{}{
				"zone":         zone.ID,
				"private_vlan": zone.PrivateVLAN,
				"public_vlan":  zone.PublicVLAN,
				"worker_count": zone.WorkerCount,
			}
			zonesConfig[j] = z
		}
		l["zones"] = zonesConfig
		workerPools[i] = l
	}

	return workerPools
}

func FlattenAlbs(list []containerv1.ALBConfig, filterType string) []map[string]interface{} {
	albs := make([]map[string]interface{}, 0)
	for _, alb := range list {
		if alb.ALBType == filterType || filterType == "all" {
			l := map[string]interface{}{
				"id":                 alb.ALBID,
				"name":               alb.Name,
				"alb_type":           alb.ALBType,
				"enable":             alb.Enable,
				"state":              alb.State,
				"num_of_instances":   alb.NumOfInstances,
				"alb_ip":             alb.ALBIP,
				"resize":             alb.Resize,
				"disable_deployment": alb.DisableDeployment,
			}
			albs = append(albs, l)
		}
	}
	return albs
}

func FlattenVpcAlbs(list []containerv2.AlbConfig, filterType string) []map[string]interface{} {
	albs := make([]map[string]interface{}, 0)
	for _, alb := range list {
		if alb.AlbType == filterType || filterType == "all" {
			l := map[string]interface{}{
				"id":                     alb.AlbID,
				"name":                   alb.Name,
				"alb_type":               alb.AlbType,
				"enable":                 alb.Enable,
				"state":                  alb.State,
				"resize":                 alb.Resize,
				"disable_deployment":     alb.DisableDeployment,
				"load_balancer_hostname": alb.LoadBalancerHostname,
			}
			albs = append(albs, l)
		}
	}
	return albs
}

func FlattenNetworkInterfaces(list []containerv2.Network) []map[string]interface{} {
	nwInterfaces := make([]map[string]interface{}, len(list))
	for i, nw := range list {
		l := map[string]interface{}{
			"cidr":       nw.Cidr,
			"ip_address": nw.IpAddress,
			"subnet_id":  nw.SubnetID,
		}
		nwInterfaces[i] = l
	}
	return nwInterfaces
}

func FlattenVlans(list []containerv1.Vlan) []map[string]interface{} {
	vlans := make([]map[string]interface{}, len(list))
	for i, vlanR := range list {
		subnets := make([]map[string]interface{}, len(vlanR.Subnets))
		for j, subnetR := range vlanR.Subnets {
			subnet := make(map[string]interface{})
			subnet["id"] = subnetR.ID
			subnet["cidr"] = subnetR.Cidr
			subnet["is_byoip"] = subnetR.IsByOIP
			subnet["is_public"] = subnetR.IsPublic
			ips := make([]string, len(subnetR.Ips))
			for k, ip := range subnetR.Ips {
				ips[k] = ip
			}
			subnet["ips"] = ips
			subnets[j] = subnet
		}
		l := map[string]interface{}{
			"id":      vlanR.ID,
			"subnets": subnets,
		}
		vlans[i] = l
	}
	return vlans
}

func FlattenIcdGroups(groupResponse *clouddatabasesv5.ListDeploymentScalingGroupsResponse) []map[string]interface{} {
	groups := make([]map[string]interface{}, len(groupResponse.Groups))
	for i, group := range groupResponse.Groups {
		memorys := make([]map[string]interface{}, 1)
		memory := make(map[string]interface{})
		memory["units"] = group.Memory.Units
		memory["allocation_mb"] = group.Memory.AllocationMb
		memory["minimum_mb"] = group.Memory.MinimumMb
		memory["step_size_mb"] = group.Memory.StepSizeMb
		memory["is_adjustable"] = group.Memory.IsAdjustable
		memory["can_scale_down"] = group.Memory.CanScaleDown
		memorys[0] = memory

		cpus := make([]map[string]interface{}, 1)
		cpu := make(map[string]interface{})
		cpu["units"] = group.CPU.Units
		cpu["allocation_count"] = group.CPU.AllocationCount
		cpu["minimum_count"] = group.CPU.MinimumCount
		cpu["step_size_count"] = group.CPU.StepSizeCount
		cpu["is_adjustable"] = group.CPU.IsAdjustable
		cpu["can_scale_down"] = group.CPU.CanScaleDown
		cpus[0] = cpu

		disks := make([]map[string]interface{}, 1)
		disk := make(map[string]interface{})
		disk["units"] = group.Disk.Units
		disk["allocation_mb"] = group.Disk.AllocationMb
		disk["minimum_mb"] = group.Disk.MinimumMb
		disk["step_size_mb"] = group.Disk.StepSizeMb
		disk["is_adjustable"] = group.Disk.IsAdjustable
		disk["can_scale_down"] = group.Disk.CanScaleDown
		disks[0] = disk

		hostflavors := make([]map[string]interface{}, 0)
		if group.HostFlavor != nil {
			hostflavors = make([]map[string]interface{}, 1)
			hostflavor := make(map[string]interface{})
			hostflavor["id"] = group.HostFlavor.ID
			hostflavor["name"] = group.HostFlavor.Name
			hostflavor["hosting_size"] = group.HostFlavor.HostingSize
			hostflavors[0] = hostflavor
		}

		l := map[string]interface{}{
			"group_id":    group.ID,
			"count":       group.Count,
			"memory":      memorys,
			"cpu":         cpus,
			"disk":        disks,
			"host_flavor": hostflavors,
		}
		groups[i] = l
	}
	return groups
}

func NormalizeJSONString(jsonString interface{}) (string, error) {
	var j interface{}
	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}
	s := jsonString.(string)
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}
	bytes, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(bytes[:]), nil
}

func ExpandAnnotations(annotations string) (whisk.KeyValueArr, error) {
	var result whisk.KeyValueArr
	dc := json.NewDecoder(strings.NewReader(annotations))
	dc.UseNumber()
	err := dc.Decode(&result)
	return result, err
}

func FlattenAnnotations(in whisk.KeyValueArr) (string, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

func ExpandParameters(annotations string) (whisk.KeyValueArr, error) {
	var result whisk.KeyValueArr
	dc := json.NewDecoder(strings.NewReader(annotations))
	dc.UseNumber()
	err := dc.Decode(&result)
	return result, err
}

func FlattenParameters(in whisk.KeyValueArr) (string, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

func ExpandLimits(l []interface{}) *whisk.Limits {
	if len(l) == 0 || l[0] == nil {
		return &whisk.Limits{}
	}
	in := l[0].(map[string]interface{})
	obj := &whisk.Limits{
		Timeout: ptrToInt(in["timeout"].(int)),
		Memory:  ptrToInt(in["memory"].(int)),
		Logsize: ptrToInt(in["log_size"].(int)),
	}
	return obj
}

func FlattenActivityTrack(in *resourceconfigurationv1.ActivityTracking) []interface{} {

	att := make(map[string]interface{})
	if in != nil {
		if in.ReadDataEvents != nil {
			att["read_data_events"] = *in.ReadDataEvents
		}
		if in.WriteDataEvents != nil {
			att["write_data_events"] = *in.WriteDataEvents
		}
		if in.ManagementEvents != nil {
			att["management_events"] = *in.ManagementEvents
		}
		if in.ActivityTrackerCrn != nil {
			att["activity_tracker_crn"] = *in.ActivityTrackerCrn
		}
	}
	return []interface{}{att}
}

func FlattenMetricsMonitor(in *resourceconfigurationv1.MetricsMonitoring) []interface{} {
	att := make(map[string]interface{})
	if in != nil {
		if in.UsageMetricsEnabled != nil {
			att["usage_metrics_enabled"] = *in.UsageMetricsEnabled
		}
		if in.MetricsMonitoringCrn != nil {
			att["metrics_monitoring_crn"] = *in.MetricsMonitoringCrn
		}
		if in.RequestMetricsEnabled != nil {
			att["request_metrics_enabled"] = *in.RequestMetricsEnabled
		}
	}
	return []interface{}{att}
}

func ArchiveRuleGet(in []*s3.LifecycleRule) []interface{} {
	rules := make([]interface{}, 0, len(in))
	for _, r := range in {
		// Checking this is not an expire_rule.  LifeCycle rules are either archive or expire or non current version or abort incomplete multipart upload
		if r.Expiration == nil && r.NoncurrentVersionExpiration == nil && r.AbortIncompleteMultipartUpload == nil {
			rule := make(map[string]interface{})

			if r.Status != nil {
				if *r.Status == "Enabled" {
					rule["enable"] = true

				} else {
					rule["enable"] = false
				}

			}
			if r.ID != nil {
				rule["rule_id"] = *r.ID
			}

			for _, transition := range r.Transitions {
				if transition.Days != nil {
					rule["days"] = int(*transition.Days)
				}
				if transition.StorageClass != nil {
					rule["type"] = *transition.StorageClass
				}
			}

			rules = append(rules, rule)
		}
	}
	return rules
}

func ExpireRuleGet(in []*s3.LifecycleRule) []interface{} {
	rules := make([]interface{}, 0, len(in))
	for _, r := range in {
		if r.Expiration != nil && r.Transitions == nil {
			rule := make(map[string]interface{})

			if r.Status != nil {
				if *r.Status == "Enabled" {
					rule["enable"] = true

				} else {
					rule["enable"] = false
				}
			}
			if r.ID != nil {
				rule["rule_id"] = *r.ID
			}

			if r.Expiration != nil {
				if r.Expiration.Days != nil {
					days := int(*(r.Expiration).Days)
					if days > 0 {
						rule["days"] = days
					}
				}
				if r.Expiration.Date != nil {
					expirationTime := *(r.Expiration).Date
					d := strings.Split(expirationTime.Format(time.RFC3339), "T")
					rule["date"] = d[0]
				}

				if r.Expiration.ExpiredObjectDeleteMarker != nil {
					rule["expired_object_delete_marker"] = *(r.Expiration).ExpiredObjectDeleteMarker
				}
			}
			if r.Filter != nil && r.Filter.Prefix != nil {
				rule["prefix"] = *(r.Filter).Prefix
			}

			rules = append(rules, rule)
		}
	}

	return rules

}

func Nc_exp_RuleGet(in []*s3.LifecycleRule) []interface{} {
	rules := make([]interface{}, 0, len(in))
	for _, r := range in {
		if r.Expiration == nil && r.AbortIncompleteMultipartUpload == nil && r.Transitions == nil {
			rule := make(map[string]interface{})
			if r.Status != nil {
				if *r.Status == "Enabled" {
					rule["enable"] = true

				} else {
					rule["enable"] = false
				}

			}
			if r.ID != nil {
				rule["rule_id"] = *r.ID
			}
			if r.NoncurrentVersionExpiration != nil {
				rule["noncurrent_days"] = int(*(r.NoncurrentVersionExpiration).NoncurrentDays)
			}
			if r.Filter != nil && r.Filter.Prefix != nil {
				rule["prefix"] = *(r.Filter).Prefix
			}
			rules = append(rules, rule)
		}
	}
	return rules
}

func Abort_mpu_RuleGet(in []*s3.LifecycleRule) []interface{} {
	rules := make([]interface{}, 0, len(in))
	for _, r := range in {
		if r.Expiration == nil && r.NoncurrentVersionExpiration == nil && r.Transitions == nil {
			rule := make(map[string]interface{})
			if r.Status != nil {
				if *r.Status == "Enabled" {
					rule["enable"] = true

				} else {
					rule["enable"] = false
				}

			}
			if r.ID != nil {
				rule["rule_id"] = *r.ID
			}
			if r.AbortIncompleteMultipartUpload != nil {
				rule["days_after_initiation"] = int(*(r.AbortIncompleteMultipartUpload).DaysAfterInitiation)
			}
			if r.Filter != nil && r.Filter.Prefix != nil {
				rule["prefix"] = *(r.Filter).Prefix
			}
			rules = append(rules, rule)
		}
	}
	return rules
}

func RetentionRuleGet(in *s3.ProtectionConfiguration) []interface{} {
	rules := make([]interface{}, 0, 1)
	if in != nil && in.Status != nil && *in.Status == "COMPLIANCE" {
		protectConfig := make(map[string]interface{})
		if in.DefaultRetention != nil {
			protectConfig["default"] = int(*(in.DefaultRetention).Days)
		}
		if in.MaximumRetention != nil {
			protectConfig["maximum"] = int(*(in.MaximumRetention).Days)
		}
		if in.MinimumRetention != nil {
			protectConfig["minimum"] = int(*(in.MinimumRetention).Days)
		}
		if in.EnablePermanentRetention != nil {
			protectConfig["permanent"] = *in.EnablePermanentRetention
		}
		rules = append(rules, protectConfig)
	}
	return rules
}

func FlattenCosObejctVersioning(in *s3.GetBucketVersioningOutput) []interface{} {
	versioning := make([]interface{}, 0, 1)
	if in != nil {
		if in.Status != nil {
			att := make(map[string]interface{})
			if *in.Status == "Enabled" {
				att["enable"] = true
			} else {
				att["enable"] = false
			}
			versioning = append(versioning, att)
		}
	}
	return versioning
}

func ReplicationRuleGet(in *s3.ReplicationConfiguration) []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, 1)
	if in != nil {
		for _, replicaterule := range in.Rules {
			replicationConfig := make(map[string]interface{})
			if replicaterule.DeleteMarkerReplication != nil {
				if *(replicaterule.DeleteMarkerReplication).Status == "Enabled" {
					replicationConfig["deletemarker_replication_status"] = true
				} else {
					replicationConfig["deletemarker_replication_status"] = false
				}
			}
			if replicaterule.Destination != nil {
				replicationConfig["destination_bucket_crn"] = *(replicaterule.Destination).Bucket
			}
			if replicaterule.ID != nil {
				replicationConfig["rule_id"] = *replicaterule.ID
			}
			if replicaterule.Priority != nil {
				replicationConfig["priority"] = int(*replicaterule.Priority)
			}
			if replicaterule.Status != nil {
				if *replicaterule.Status == "Enabled" {
					replicationConfig["enable"] = true
				} else {
					replicationConfig["enable"] = false
				}
			}
			if replicaterule.Filter != nil && replicaterule.Filter.Prefix != nil {
				replicationConfig["prefix"] = *(replicaterule.Filter).Prefix
			}
			rules = append(rules, replicationConfig)
		}
	}
	return rules
}

func flattenLifecycleExpiration(expiration *s3.LifecycleExpiration) []interface{} {
	if expiration == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})
	if expiration.Date != nil {
		m["date"] = expiration.Date.Format(time.RFC3339)
	}
	if expiration.Days != nil {
		m["days"] = int(aws.Int64Value(expiration.Days))
	}
	if expiration.ExpiredObjectDeleteMarker != nil {
		m["expired_object_delete_marker"] = aws.Bool(*expiration.ExpiredObjectDeleteMarker)
	}
	return []interface{}{m}
}

func flattenNoncurrentVersionExpiration(expiration *s3.NoncurrentVersionExpiration) []interface{} {
	if expiration == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})
	if expiration.NoncurrentDays != nil {
		m["noncurrent_days"] = int(aws.Int64Value(expiration.NoncurrentDays))
	}
	return []interface{}{m}
}
func flattenTransitions(transitions []*s3.Transition) []interface{} {
	if len(transitions) == 0 {
		return []interface{}{}
	}
	var results []interface{}
	for _, transition := range transitions {
		m := make(map[string]interface{})
		if transition.StorageClass != nil {
			m["storage_class"] = transition.StorageClass
		}
		if transition.Date != nil {
			m["date"] = transition.Date.Format(time.RFC3339)
		}
		if transition.Days != nil {
			m["days"] = int(aws.Int64Value(transition.Days))
		}
		results = append(results, m)
	}
	return results
}

func flattenLifecycleRuleFilter(filter *s3.LifecycleRuleFilter) []interface{} {
	if filter == nil {
		return []interface{}{}
	}
	m := make(map[string]interface{})
	if filter.And != nil {
		m["and"] = flattenLifecycleRuleFilterMemberAnd(filter.And)

	}
	if filter.ObjectSizeGreaterThan != nil && int(aws.Int64Value(filter.ObjectSizeGreaterThan)) > 0 {

		m["object_size_greater_than"] = int(aws.Int64Value(filter.ObjectSizeGreaterThan))

	}
	if filter.ObjectSizeLessThan != nil && int(aws.Int64Value(filter.ObjectSizeLessThan)) > 0 {

		m["object_size_less_than"] = int(aws.Int64Value(filter.ObjectSizeLessThan))

	}
	if filter.Tag != nil {
		m["tag"] = flattenLifecycleRuleFilterMemberTag(filter.Tag)

	}
	if filter.Prefix != nil {
		m["prefix"] = aws.String(*filter.Prefix)
	}
	return []interface{}{m}
}
func flattenLifecycleRuleFilterMemberAnd(andOp *s3.LifecycleRuleAndOperator) []interface{} {
	if andOp == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"object_size_greater_than": andOp.ObjectSizeGreaterThan,
		"object_size_less_than":    andOp.ObjectSizeLessThan,
	}

	if v := andOp.Prefix; v != nil {
		m["prefix"] = aws.StringValue(v)
	}

	if v := andOp.Tags; v != nil {
		m["tags"] = flattenMultipleTags(v)
	}

	return []interface{}{m}
}

func flattenMultipleTags(in []*s3.Tag) []map[string]interface{} {

	tagSet := make([]map[string]interface{}, 0, len(in))
	if in != nil {
		for _, tagSetValue := range in {
			tag := make(map[string]interface{})
			if tagSetValue.Key != nil {
				tag["key"] = *tagSetValue.Key
			}
			if tagSetValue.Value != nil {
				tag["value"] = *tagSetValue.Value
			}

			tagSet = append(tagSet, tag)

		}

	}

	return tagSet
}

func flattenLifecycleRuleFilterMemberTag(op *s3.Tag) []interface{} {
	if op == nil {
		return nil
	}

	m := make(map[string]interface{})

	if v := op.Key; v != nil {
		m["key"] = aws.StringValue(v)
	}

	if v := op.Value; v != nil {
		m["value"] = aws.StringValue(v)
	}

	return []interface{}{m}
}

func flattenAbortIncompleteMultipartUpload(abortIncompleteMultipartUploadInput *s3.AbortIncompleteMultipartUpload) []interface{} {
	if abortIncompleteMultipartUploadInput == nil {
		return []interface{}{}
	}

	abortIncompleteMultipartUploadMap := make(map[string]interface{})

	if abortIncompleteMultipartUploadInput.DaysAfterInitiation != nil {
		abortIncompleteMultipartUploadMap["days_after_initiation"] = int(aws.Int64Value(abortIncompleteMultipartUploadInput.DaysAfterInitiation))
	}

	return []interface{}{abortIncompleteMultipartUploadMap}
}

func LifecylceRuleGet(lifecycleRuleInput []*s3.LifecycleRule) []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, len(lifecycleRuleInput))
	if lifecycleRuleInput != nil {
		for _, lifecyclerule := range lifecycleRuleInput {
			lifecycleRuleConfig := make(map[string]interface{})
			if lifecyclerule.Status != nil {
				if *lifecyclerule.Status == "Enabled" {
					lifecycleRuleConfig["status"] = "enable"
				} else if *lifecyclerule.Status == "Disabled" {
					lifecycleRuleConfig["status"] = "disable"

				}
			}
			if lifecyclerule.ID != nil {
				lifecycleRuleConfig["rule_id"] = *lifecyclerule.ID
			}
			if lifecyclerule.Expiration != nil {
				lifecycleRuleConfig["expiration"] = flattenLifecycleExpiration(lifecyclerule.Expiration)
			}
			if lifecyclerule.Transitions != nil {
				lifecycleRuleConfig["transition"] = flattenTransitions(lifecyclerule.Transitions)
			}
			if lifecyclerule.AbortIncompleteMultipartUpload != nil {
				lifecycleRuleConfig["abort_incomplete_multipart_upload"] = flattenAbortIncompleteMultipartUpload(lifecyclerule.AbortIncompleteMultipartUpload)
			}
			if lifecyclerule.NoncurrentVersionExpiration != nil {
				lifecycleRuleConfig["noncurrent_version_expiration"] = flattenNoncurrentVersionExpiration(lifecyclerule.NoncurrentVersionExpiration)
			}
			if lifecyclerule.Filter != nil {
				lifecycleRuleConfig["filter"] = flattenLifecycleRuleFilter(lifecyclerule.Filter)
			}
			rules = append(rules, lifecycleRuleConfig)
		}
	}
	return rules
}
func ObjectLockConfigurationGet(in *s3.ObjectLockConfiguration) []map[string]interface{} {
	configuration := make([]map[string]interface{}, 0, 1)
	if in != nil {
		objectLockConfig := make(map[string]interface{})

		if in.ObjectLockEnabled != nil {
			objectLockConfig["object_lock_enabled"] = s3.ObjectLockEnabledEnabled
		}
		if in.Rule != nil {
			objectLockConfig["object_lock_rule"] = ObjectLockRuleGet(in.Rule)
		}

		configuration = append(configuration, objectLockConfig)
	}
	return configuration
}
func ObjectLockRuleGet(in *s3.ObjectLockRule) []map[string]interface{} {
	objectLockRule := make([]map[string]interface{}, 0, 1)
	if in != nil {
		objectLockConfig := make(map[string]interface{})

		if in.DefaultRetention != nil {
			objectLockConfig["default_retention"] = ObjectLockDefaultRetentionGet(in.DefaultRetention)
		}

		objectLockRule = append(objectLockRule, objectLockConfig)
	}
	return objectLockRule
}

func ObjectLockDefaultRetentionGet(in *s3.DefaultRetention) []map[string]interface{} {
	defaultRetention := make([]map[string]interface{}, 0, 1)
	if in != nil {
		defaultRetentionMap := make(map[string]interface{})

		if in.Days != nil {
			defaultRetentionMap["days"] = int(aws.Int64Value(in.Days))
		}

		if in.Mode != nil {
			defaultRetentionMap["mode"] = aws.StringValue(in.Mode)
		}

		if in.Years != nil {
			defaultRetentionMap["years"] = int(aws.Int64Value(in.Years))
		}

		defaultRetention = append(defaultRetention, defaultRetentionMap)

	}
	return defaultRetention
}

func WebsiteConfigurationGet(in *s3.WebsiteConfiguration) []map[string]interface{} {
	configuration := make([]map[string]interface{}, 0, 1)
	if in != nil {
		websiteConfig := make(map[string]interface{})

		if in.ErrorDocument != nil {
			websiteConfig["error_document"] = GetErrorDocument(in.ErrorDocument)
		}
		if in.IndexDocument != nil {
			websiteConfig["index_document"] = GetIndexDocument(in.IndexDocument)
		}
		if in.RedirectAllRequestsTo != nil {
			websiteConfig["redirect_all_requests_to"] = RedirectAllRequestsGet(in.RedirectAllRequestsTo)
		}

		if in.RoutingRules != nil {
			websiteConfig["routing_rule"] = RoutingRulesGet(in.RoutingRules)
		}

		configuration = append(configuration, websiteConfig)
	}
	return configuration
}

func GetErrorDocument(in *s3.ErrorDocument) []map[string]interface{} {
	errorDocumentMap := make([]map[string]interface{}, 0, 1)
	if in != nil {
		errorDocValue := make(map[string]interface{})

		if in.Key != nil {
			errorDocValue["key"] = aws.StringValue(in.Key)
		}
		errorDocumentMap = append(errorDocumentMap, errorDocValue)
	}
	return errorDocumentMap
}

func GetIndexDocument(in *s3.IndexDocument) []map[string]interface{} {
	indexDocumentMap := make([]map[string]interface{}, 0, 1)
	if in != nil {
		indexDocumentValue := make(map[string]interface{})

		if in.Suffix != nil {
			indexDocumentValue["suffix"] = aws.StringValue(in.Suffix)
		}
		indexDocumentMap = append(indexDocumentMap, indexDocumentValue)
	}
	return indexDocumentMap
}

func RedirectAllRequestsGet(in *s3.RedirectAllRequestsTo) []map[string]interface{} {
	redirectRequests := make([]map[string]interface{}, 0, 1)
	if in != nil {
		redirectRequestConfig := make(map[string]interface{})

		if in.HostName != nil {
			redirectRequestConfig["host_name"] = aws.StringValue(in.HostName)
		}
		if in.Protocol != nil {
			redirectRequestConfig["protocol"] = aws.StringValue(in.Protocol)
		}

		redirectRequests = append(redirectRequests, redirectRequestConfig)
	}
	return redirectRequests
}

func RoutingRulesGet(in []*s3.RoutingRule) []map[string]interface{} {
	routingRules := make([]map[string]interface{}, 0, len(in))
	if in != nil {
		for _, routingRuleValue := range in {
			rule := make(map[string]interface{})

			if routingRuleValue.Condition != nil {
				rule["condition"] = RoutingRuleConditionGet(routingRuleValue.Condition)
			}

			if routingRuleValue.Redirect != nil {
				rule["redirect"] = RoutingRuleRedirectGet(routingRuleValue.Redirect)
			}

			routingRules = append(routingRules, rule)
		}
	}
	return routingRules
}

func RoutingRuleConditionGet(in *s3.Condition) []map[string]interface{} {
	condition := make([]map[string]interface{}, 0, 1)
	if in != nil {
		conditionConfig := make(map[string]interface{})

		if in.HttpErrorCodeReturnedEquals != nil {
			conditionConfig["http_error_code_returned_equals"] = aws.StringValue(in.HttpErrorCodeReturnedEquals)
		}
		if in.KeyPrefixEquals != nil {
			conditionConfig["key_prefix_equals"] = aws.StringValue(in.KeyPrefixEquals)
		}

		condition = append(condition, conditionConfig)
	}
	return condition
}

func RoutingRuleRedirectGet(in *s3.Redirect) []map[string]interface{} {
	redirect := make([]map[string]interface{}, 0, 1)
	if in != nil {
		redirectConfig := make(map[string]interface{})

		if in.HostName != nil {
			redirectConfig["host_name"] = aws.StringValue(in.HostName)
		}
		if in.HttpRedirectCode != nil {
			redirectConfig["http_redirect_code"] = aws.StringValue(in.HttpRedirectCode)
		}
		if in.Protocol != nil {
			redirectConfig["protocol"] = aws.StringValue(in.Protocol)
		}
		if in.ReplaceKeyPrefixWith != nil {
			redirectConfig["replace_key_prefix_with"] = aws.StringValue(in.ReplaceKeyPrefixWith)
		}
		if in.ReplaceKeyWith != nil {
			redirectConfig["replace_key_with"] = aws.StringValue(in.ReplaceKeyWith)
		}

		redirect = append(redirect, redirectConfig)
	}
	return redirect
}

func FlattenLimits(in *whisk.Limits) []interface{} {
	att := make(map[string]interface{})
	if in.Timeout != nil {
		att["timeout"] = *in.Timeout
	}
	if in.Memory != nil {
		att["memory"] = *in.Memory
	}
	if in.Memory != nil {
		att["log_size"] = *in.Logsize
	}
	return []interface{}{att}
}

func ExpandExec(execs []interface{}) *whisk.Exec {
	var code string
	var document []byte
	for _, exec := range execs {
		e, _ := exec.(map[string]interface{})
		code_path := e["code_path"].(string)
		if code_path != "" {
			ext := path.Ext(code_path)
			if strings.ToLower(ext) == ".zip" {
				data, err := ioutil.ReadFile(code_path)
				if err != nil {
					log.Println("Error reading file", err)
					return &whisk.Exec{}
				}
				sEnc := b64.StdEncoding.EncodeToString([]byte(data))
				code = sEnc
			} else {
				data, err := ioutil.ReadFile(code_path)
				if err != nil {
					log.Println("Error reading file", err)
					return &whisk.Exec{}
				}
				document = data
				code = string(document)
			}
		} else {
			code = e["code"].(string)
		}
		obj := &whisk.Exec{
			Image:      e["image"].(string),
			Init:       e["init"].(string),
			Code:       PtrToString(code),
			Kind:       e["kind"].(string),
			Main:       e["main"].(string),
			Components: ExpandStringList(e["components"].([]interface{})),
		}
		return obj
	}

	return &whisk.Exec{}
}

func FlattenExec(in *whisk.Exec, d *schema.ResourceData) []interface{} {
	code_data := 4194304 // length of 'code' parameter should be always <= 4MB data
	att := make(map[string]interface{})
	// open-whisk SDK will not return the value for code_path
	// Hence using d.GetOk method to setback the code_path value.
	if cPath, ok := d.GetOk("exec.0.code_path"); ok {
		att["code_path"] = cPath.(string)
	}
	if in.Image != "" {
		att["image"] = in.Image
	}
	if in.Init != "" {
		att["init"] = in.Init
	}
	if in != nil && in.Code != nil && len(*in.Code) <= code_data {
		att["code"] = *in.Code
	}
	if in.Kind != "" {
		att["kind"] = in.Kind
	}
	if in.Main != "" {
		att["main"] = in.Main
	}

	if len(in.Components) > 0 {
		att["components"] = FlattenStringList(in.Components)
	}

	return []interface{}{att}
}

func ptrToInt(i int) *int {
	return &i
}

func PtrToString(s string) *string {
	return &s
}

func PtrToBool(b bool) *bool {
	return &b
}

func IntValue(i64 *int64) (i int) {
	if i64 != nil {
		i = int(*i64)
	}
	return
}

func Float64Value(f32 *float32) (f float64) {
	if f32 != nil {
		f = float64(*f32)
	}
	return
}

func StringValue(strPtr *string) (_ string) {
	if strPtr != nil {
		return *strPtr
	}
	return
}

func DateToString(d *strfmt.Date) (s string) {
	if d != nil {
		s = d.String()
	}
	return
}

func DateTimeToString(dt *strfmt.DateTime) (s string) {
	if dt != nil {
		s = dt.String()
	}
	return
}

func FilterActionAnnotations(in whisk.KeyValueArr) (string, error) {
	noExec := make(whisk.KeyValueArr, 0, len(in))
	for _, v := range in {
		if v.Key == "exec" {
			continue
		}
		noExec = append(noExec, v)
	}

	return FlattenAnnotations(noExec)
}

func FilterActionParameters(in whisk.KeyValueArr) (string, error) {
	noAction := make(whisk.KeyValueArr, 0, len(in))
	for _, v := range in {
		if v.Key == "_actions" {
			continue
		}
		noAction = append(noAction, v)
	}
	return FlattenParameters(noAction)
}

func FilterInheritedAnnotations(inheritedAnnotations, annotations whisk.KeyValueArr) whisk.KeyValueArr {
	userDefinedAnnotations := make(whisk.KeyValueArr, 0)
	for _, a := range annotations {
		insert := false
		if a.Key == "binding" || a.Key == "exec" {
			insert = false
			break
		}
		for _, b := range inheritedAnnotations {
			if a.Key == b.Key && reflect.DeepEqual(a.Value, b.Value) {
				insert = false
				break
			}
			insert = true
		}
		if insert {
			userDefinedAnnotations = append(userDefinedAnnotations, a)
		}
	}
	return userDefinedAnnotations
}

func FilterInheritedParameters(inheritedParameters, parameters whisk.KeyValueArr) whisk.KeyValueArr {
	userDefinedParameters := make(whisk.KeyValueArr, 0)
	for _, p := range parameters {
		insert := false
		if p.Key == "_actions" {
			insert = false
			break
		}
		for _, b := range inheritedParameters {
			if p.Key == b.Key && reflect.DeepEqual(p.Value, b.Value) {
				insert = false
				break
			}
			insert = true
		}
		if insert {
			userDefinedParameters = append(userDefinedParameters, p)
		}

	}
	return userDefinedParameters
}

func IsEmpty(object interface{}) bool {
	//First check normal definitions of empty
	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == false {
		return true
	}

	//Then see if it's a struct
	if reflect.ValueOf(object).Kind() == reflect.Struct {
		// and create an empty copy of the struct object to compare against
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}
	return false
}

func FilterTriggerAnnotations(in whisk.KeyValueArr) (string, error) {
	noFeed := make(whisk.KeyValueArr, 0, len(in))
	for _, v := range in {
		if v.Key == "feed" {
			continue
		}
		noFeed = append(noFeed, v)
	}
	return FlattenParameters(noFeed)
}

func FlattenFeed(feedName string) []interface{} {
	att := make(map[string]interface{})
	att["name"] = feedName
	att["parameters"] = "[]"
	return []interface{}{att}
}

func FlattenGatewayVlans(list []datatypes.Network_Gateway_Vlan) []map[string]interface{} {
	vlans := make([]map[string]interface{}, len(list))
	for i, ele := range list {
		vlan := make(map[string]interface{})
		vlan["bypass"] = *ele.BypassFlag
		vlan["network_vlan_id"] = *ele.NetworkVlanId
		vlan["vlan_id"] = *ele.Id
		vlans[i] = vlan
	}
	return vlans
}

func FlattenGatewayMembers(d *schema.ResourceData, list []datatypes.Network_Gateway_Member) []map[string]interface{} {
	members := make([]map[string]interface{}, len(list))
	for i, ele := range list {
		hardware := *ele.Hardware
		member := make(map[string]interface{})
		member["member_id"] = *ele.HardwareId
		member["hostname"] = *hardware.Hostname
		member["domain"] = *hardware.Domain
		if hardware.Notes != nil {
			member["notes"] = *hardware.Notes
		}
		if hardware.Datacenter != nil {
			member["datacenter"] = *hardware.Datacenter.Name
		}
		if hardware.PrimaryNetworkComponent.MaxSpeed != nil {
			member["network_speed"] = *hardware.PrimaryNetworkComponent.MaxSpeed
		}
		member["redundant_network"] = false
		member["unbonded_network"] = false
		backendNetworkComponent := ele.Hardware.BackendNetworkComponents

		if len(backendNetworkComponent) > 2 && ele.Hardware.PrimaryBackendNetworkComponent != nil {
			if *hardware.PrimaryBackendNetworkComponent.RedundancyEnabledFlag {
				member["redundant_network"] = true
			} else {
				member["unbonded_network"] = true
			}
		}
		tagReferences := ele.Hardware.TagReferences
		tagReferencesLen := len(tagReferences)
		if tagReferencesLen > 0 {
			tags := make([]interface{}, 0, tagReferencesLen)
			for _, tagRef := range tagReferences {
				tags = append(tags, *tagRef.Tag.Name)
			}
			member["tags"] = schema.NewSet(schema.HashString, tags)
		}

		member["redundant_power_supply"] = false

		if *hardware.PowerSupplyCount == 2 {
			member["redundant_power_supply"] = true
		}
		member["memory"] = *hardware.MemoryCapacity
		if !(*hardware.PrivateNetworkOnlyFlag) {
			member["public_vlan_id"] = *hardware.NetworkVlans[1].Id
		}
		member["private_vlan_id"] = *hardware.NetworkVlans[0].Id

		if hardware.PrimaryIpAddress != nil {
			member["public_ipv4_address"] = *hardware.PrimaryIpAddress
		}
		if hardware.PrimaryBackendIpAddress != nil {
			member["private_ipv4_address"] = *hardware.PrimaryBackendIpAddress
		}
		member["ipv6_enabled"] = false
		if ele.Hardware.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord != nil {
			member["ipv6_enabled"] = true
			member["ipv6_address"] = *hardware.PrimaryNetworkComponent.PrimaryVersion6IpAddressRecord.IpAddress
		}

		member["private_network_only"] = *hardware.PrivateNetworkOnlyFlag
		userData := hardware.UserData
		if len(userData) > 0 && userData[0].Value != nil {
			member["user_metadata"] = *userData[0].Value
		}
		members[i] = member
	}
	return members
}

func FlattenDisks(result datatypes.Virtual_Guest) []int {
	var out = make([]int, 0)

	for _, v := range result.BlockDevices {
		// skip 1,7 which is reserved for the swap disk and metadata
		_, ok := sl.GrabOk(result, "BillingItem.OrderItem.Preset")
		if ok {
			if *v.Device != "1" && *v.Device != "7" && *v.Device != "0" {
				capacity, ok := sl.GrabOk(v, "DiskImage.Capacity")

				if ok {
					out = append(out, capacity.(int))
				}

			}
		} else {
			if *v.Device != "1" && *v.Device != "7" {
				capacity, ok := sl.GrabOk(v, "DiskImage.Capacity")

				if ok {
					out = append(out, capacity.(int))
				}
			}
		}
	}

	return out
}

func FlattenDisksForWindows(result datatypes.Virtual_Guest) []int {
	var out = make([]int, 0)

	for _, v := range result.BlockDevices {
		// skip 1,7 which is reserved for the swap disk and metadata
		_, ok := sl.GrabOk(result, "BillingItem.OrderItem.Preset")
		if ok {
			if *v.Device != "1" && *v.Device != "7" && *v.Device != "0" && *v.Device != "3" {
				capacity, ok := sl.GrabOk(v, "DiskImage.Capacity")

				if ok {
					out = append(out, capacity.(int))
				}
			}
		} else {
			if *v.Device != "1" && *v.Device != "7" && *v.Device != "3" {
				capacity, ok := sl.GrabOk(v, "DiskImage.Capacity")

				if ok {
					out = append(out, capacity.(int))
				}
			}
		}
	}

	return out
}

func filterResourceKeyParameters(params map[string]interface{}) map[string]interface{} {
	delete(params, "role_crn")
	return params
}

func IdParts(id string) ([]string, error) {
	if strings.Contains(id, "/") {
		parts := strings.Split(id, "/")
		return parts, nil
	}
	return []string{}, fmt.Errorf("The given id %s does not contain / please check documentation on how to provider id during import command", id)
}

func SepIdParts(id string, separator string) ([]string, error) {
	if strings.Contains(id, separator) {
		parts := strings.Split(id, separator)
		return parts, nil
	}
	return []string{}, fmt.Errorf("The given id %s does not contain %s please check documentation on how to provider id during import command", id, separator)
}

func VmIdParts(id string) ([]string, error) {
	parts := strings.Split(id, "/")
	return parts, nil
}

func CfIdParts(id string) ([]string, error) {
	parts := strings.Split(id, ":")
	return parts, nil
}

// getCustomAttributes will return all attributes which are not system defined
func getCustomAttributes(r iampolicymanagementv1.PolicyResource) []iampolicymanagementv1.ResourceAttribute {
	attributes := []iampolicymanagementv1.ResourceAttribute{}
	for _, a := range r.Attributes {
		switch *a.Name {
		case "accesGroupId":
		case "accountId":
		case "organizationId":
		case "spaceId":
		case "region":
		case "resource":
		case "resourceType":
		case "resourceGroupId":
		case "serviceType":
		case "serviceName":
		case "serviceInstance":
		default:
			attributes = append(attributes, a)
		}
	}
	return attributes
}

func GetV2PolicyCustomAttributes(r iampolicymanagementv1.V2PolicyResource) []iampolicymanagementv1.V2PolicyResourceAttribute {
	attributes := []iampolicymanagementv1.V2PolicyResourceAttribute{}
	for _, a := range r.Attributes {

		switch *a.Key {
		case "accesGroupId":
		case "accountId":
		case "organizationId":
		case "spaceId":
		case "region":
		case "resource":
		case "resourceType":
		case "resourceGroupId":
		case "serviceType":
		case "serviceName":
		case "serviceInstance":
		case "service_group_id":
		default:
			attributes = append(attributes, a)
		}
	}
	return attributes
}

func FlattenPolicyResource(list []iampolicymanagementv1.PolicyResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"service":              GetResourceAttribute("serviceName", i),
			"resource_instance_id": GetResourceAttribute("serviceInstance", i),
			"region":               GetResourceAttribute("region", i),
			"resource_type":        GetResourceAttribute("resourceType", i),
			"resource":             GetResourceAttribute("resource", i),
			"resource_group_id":    GetResourceAttribute("resourceGroupId", i),
			"service_type":         GetResourceAttribute("serviceType", i),
		}
		customAttributes := getCustomAttributes(i)
		if len(customAttributes) > 0 {
			out := make(map[string]string)
			for _, a := range customAttributes {
				out[*a.Name] = *a.Value
			}
			l["attributes"] = out
		}

		result = append(result, l)
	}
	return result
}
func FlattenPolicyResourceAttributes(list []iampolicymanagementv1.PolicyResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, i := range list {
		for _, a := range i.Attributes {
			if *a.Name != "accountId" {
				l := map[string]interface{}{
					"name":     a.Name,
					"value":    a.Value,
					"operator": a.Operator,
				}
				result = append(result, l)
			}
		}
	}
	return result
}

func FlattenPolicyResourceTags(resources []iampolicymanagementv1.PolicyResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	for _, resource := range resources {
		if resource.Tags != nil {
			for _, tags := range resource.Tags {
				tag := map[string]interface{}{
					"name":     tags.Name,
					"value":    tags.Value,
					"operator": tags.Operator,
				}
				result = append(result, tag)
			}
		}
	}
	return result
}

func FlattenV2PolicyResource(resource iampolicymanagementv1.V2PolicyResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	l := map[string]interface{}{
		"service":              GetV2PolicyResourceAttribute("serviceName", resource),
		"resource_instance_id": GetV2PolicyResourceAttribute("serviceInstance", resource),
		"region":               GetV2PolicyResourceAttribute("region", resource),
		"resource_type":        GetV2PolicyResourceAttribute("resourceType", resource),
		"resource":             GetV2PolicyResourceAttribute("resource", resource),
		"resource_group_id":    GetV2PolicyResourceAttribute("resourceGroupId", resource),
		"service_type":         GetV2PolicyResourceAttribute("serviceType", resource),
		"service_group_id":     GetV2PolicyResourceAttribute("service_group_id", resource),
	}
	customAttributes := GetV2PolicyCustomAttributes(resource)

	if len(customAttributes) > 0 {
		out := make(map[string]string)
		for _, a := range customAttributes {
			if *a.Operator == "stringExists" && a.Value == true {
				out[*a.Key] = fmt.Sprint("*")
			} else if *a.Operator == "stringMatch" || *a.Operator == "stringEquals" {
				out[*a.Key] = fmt.Sprint(a.Value)
			}
		}
		l["attributes"] = out
	}
	result = append(result, l)

	return result
}

func FlattenV2PolicyResourceAttributes(attributes []iampolicymanagementv1.V2PolicyResourceAttribute) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, a := range attributes {
		if *a.Key != "accountId" {
			l := map[string]interface{}{
				"name":     a.Key,
				"value":    a.Value,
				"operator": a.Operator,
			}
			result = append(result, l)
		}
	}
	return result
}

func FlattenV2PolicyResourceTags(resource iampolicymanagementv1.V2PolicyResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if resource.Tags != nil {
		for _, tags := range resource.Tags {
			tag := map[string]interface{}{
				"name":     tags.Key,
				"value":    tags.Value,
				"operator": tags.Operator,
			}
			result = append(result, tag)
		}
	}
	return result
}

func getConditionValues(v interface{}) []string {
	var values []string
	switch value := v.(type) {
	case string:
		values = append(values, value)
	case []interface{}:
		for _, v := range value {
			values = append(values, fmt.Sprint(v))
		}
	case nil:
	default:
		values = append(values, fmt.Sprintf("%v", value))
	}
	return values
}

func FlattenRuleConditions(rule iampolicymanagementv1.V2PolicyRule) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if len(rule.Conditions) > 0 {
		for _, cIntf := range rule.Conditions {
			c := cIntf.(*iampolicymanagementv1.NestedCondition)
			if len(c.Conditions) > 0 {
				nestedConditions := make([]map[string]interface{}, 0)
				for _, nc := range c.Conditions {
					values := getConditionValues(nc.Value)
					nestedCondition := map[string]interface{}{
						"key":      nc.Key,
						"value":    values,
						"operator": nc.Operator,
					}
					nestedConditions = append(nestedConditions, nestedCondition)
				}
				condition := map[string]interface{}{
					"operator":   c.Operator,
					"conditions": nestedConditions,
				}
				result = append(result, condition)
			} else {
				values := getConditionValues(c.Value)
				condition := map[string]interface{}{
					"key":      c.Key,
					"value":    values,
					"operator": c.Operator,
				}
				result = append(result, condition)
			}
		}
	} else {
		values := getConditionValues(rule.Value)
		condition := map[string]interface{}{
			"key":      rule.Key,
			"value":    values,
			"operator": rule.Operator,
		}
		result = append(result, condition)
	}
	return result
}

func FlattenAMSettingsExternalIdentityInteraction(amAccountSettings *iampolicymanagementv1.AccountSettingsAccessManagement) []map[string]interface{} {
	identityTypes := make([]map[string]interface{}, 0)
	if amAccountSettings.ExternalAccountIdentityInteraction != nil && amAccountSettings.ExternalAccountIdentityInteraction.IdentityTypes != nil {
		iTypes := amAccountSettings.ExternalAccountIdentityInteraction.IdentityTypes
		user := make([]map[string]interface{}, 0)
		if iTypes.User != nil {
			u := map[string]interface{}{
				"state":                     iTypes.User.State,
				"external_allowed_accounts": iTypes.User.ExternalAllowedAccounts,
			}
			user = append(user, u)
		}
		serviceId := make([]map[string]interface{}, 0)
		if iTypes.ServiceID != nil {
			sid := map[string]interface{}{
				"state":                     iTypes.ServiceID.State,
				"external_allowed_accounts": iTypes.ServiceID.ExternalAllowedAccounts,
			}
			serviceId = append(serviceId, sid)
		}
		service := make([]map[string]interface{}, 0)
		if iTypes.Service != nil {
			s := map[string]interface{}{
				"state":                     iTypes.Service.State,
				"external_allowed_accounts": iTypes.Service.ExternalAllowedAccounts,
			}
			service = append(service, s)
		}
		identityTypes = append(identityTypes, map[string]interface{}{
			"user":       user,
			"service_id": serviceId,
			"service":    service,
		})
	}
	externalIdentityInteraction := make([]map[string]interface{}, 0)
	externalIdentityInteraction = append(externalIdentityInteraction, map[string]interface{}{
		"identity_types": identityTypes,
	})
	return externalIdentityInteraction
}

func GenerateExternalAccountIdentityInteraction(d *schema.ResourceData) iampolicymanagementv1.ExternalAccountIdentityInteractionPatch {
	interaction := getElementFromResource(d, "external_account_identity_interaction")
	if interaction != nil {
		identityTypes := getFirstElementFromSet(interaction["identity_types"])
		if identityTypes != nil {
			identityTypesPatch := iampolicymanagementv1.IdentityTypesPatch{}
			user := getFirstElementFromSet(identityTypes["user"])
			if user != nil {
				state := user["state"].(string)
				accounts := user["external_allowed_accounts"].([]interface{})
				userPatch := iampolicymanagementv1.IdentityTypesBase{
					State:                   &state,
					ExternalAllowedAccounts: interfaceSliceToStringSlice(accounts),
				}
				identityTypesPatch.User = &userPatch
			}

			serviceId := getFirstElementFromSet(identityTypes["service_id"])
			if serviceId != nil {
				state := serviceId["state"].(string)
				accounts := serviceId["external_allowed_accounts"].([]interface{})
				serviceIdPatch := iampolicymanagementv1.IdentityTypesBase{
					State:                   &state,
					ExternalAllowedAccounts: interfaceSliceToStringSlice(accounts),
				}
				identityTypesPatch.ServiceID = &serviceIdPatch
			}

			service := getFirstElementFromSet(identityTypes["service"])
			if service != nil {
				state := service["state"].(string)
				accounts := service["external_allowed_accounts"].([]interface{})
				servicePatch := iampolicymanagementv1.IdentityTypesBase{
					State:                   &state,
					ExternalAllowedAccounts: interfaceSliceToStringSlice(accounts),
				}
				identityTypesPatch.Service = &servicePatch
			}
			return iampolicymanagementv1.ExternalAccountIdentityInteractionPatch{
				IdentityTypes: &identityTypesPatch,
			}
		}
	}
	return iampolicymanagementv1.ExternalAccountIdentityInteractionPatch{}
}

func getFirstElementFromSet(elem interface{}) map[string]interface{} {
	set, ok := elem.(*schema.Set)
	if ok && len(set.List()) != 0 {
		return set.List()[0].(map[string]interface{})
	}
	return nil
}

func getElementFromResource(d *schema.ResourceData, key string) map[string]interface{} {
	if elem, ok := d.GetOk(key); ok {
		return getFirstElementFromSet(elem)
	}
	return nil
}

func interfaceSliceToStringSlice(interfaceSlice []interface{}) []string {
	stringSlice := make([]string, 0, len(interfaceSlice))
	for _, element := range interfaceSlice {
		if str, ok := element.(string); ok {
			stringSlice = append(stringSlice, str)
		} else {
			stringSlice = append(stringSlice, fmt.Sprintf("%v", element))
		}
	}
	return stringSlice
}

// Cloud Internet Services
func FlattenHealthMonitors(list []datatypes.Network_LBaaS_Listener) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	ports := make([]int, 0, 0)
	for _, i := range list {
		l := map[string]interface{}{
			"protocol":    *i.DefaultPool.Protocol,
			"port":        *i.DefaultPool.ProtocolPort,
			"interval":    *i.DefaultPool.HealthMonitor.Interval,
			"max_retries": *i.DefaultPool.HealthMonitor.MaxRetries,
			"timeout":     *i.DefaultPool.HealthMonitor.Timeout,
			"monitor_id":  *i.DefaultPool.HealthMonitor.Uuid,
		}

		if i.DefaultPool.HealthMonitor.UrlPath != nil {
			l["url_path"] = *i.DefaultPool.HealthMonitor.UrlPath
		}

		if !contains(ports, *i.DefaultPool.ProtocolPort) {
			result = append(result, l)
		}

		ports = append(ports, *i.DefaultPool.ProtocolPort)
	}
	return result
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func StringContains(s []string, str string) bool {
	for _, a := range s {
		if a == str {
			return true
		}
	}
	return false
}

func FlattenMembersData(list []iamaccessgroupsv2.ListGroupMembersResponseMember, users []usermanagementv2.UserInfo, serviceids []iamidentityv1.ServiceID, profileids []iamidentityv1.TrustedProfile) ([]string, []string, []string) {
	var ibmid []string
	var serviceid []string
	var profileid []string
	for _, m := range list {
		if *m.Type == "user" {
			for _, user := range users {
				if user.IamID == *m.IamID {
					ibmid = append(ibmid, user.Email)
					break
				}
			}
		} else if *m.Type == "profile" {
			for _, prid := range profileids {
				if *prid.IamID == *m.IamID {
					profileid = append(profileid, *prid.ID)
					break
				}
			}
		} else {
			for _, srid := range serviceids {
				if *srid.IamID == *m.IamID {
					serviceid = append(serviceid, *srid.ID)
					break
				}
			}
		}

	}
	return ibmid, serviceid, profileid
}

func FlattenAccessGroupMembers(list []iamaccessgroupsv2.ListGroupMembersResponseMember, users []usermanagementv2.UserInfo, serviceids []iamidentityv1.ServiceID) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, m := range list {
		var value, vtype string
		vtype = *m.Type
		if *m.Type == "user" {
			for _, user := range users {
				if user.IamID == *m.IamID {
					value = user.Email
					break
				}
			}
		} else {
			for _, srid := range serviceids {
				if *srid.IamID == *m.IamID {
					value = *srid.ID
					break
				}
			}

		}
		l := map[string]interface{}{
			"iam_id": value,
			"type":   vtype,
		}
		result = append(result, l)
	}
	return result
}

func FlattenUserIds(accountID string, users []string, meta interface{}) ([]string, error) {
	userids := make([]string, len(users))
	for i, name := range users {
		iamID, err := GetIBMUniqueId(accountID, name, meta)
		if err != nil {
			return nil, err
		}
		userids[i] = iamID
	}
	return userids, nil
}

func ExpandUsers(userList *schema.Set) (users []icdv4.User) {
	for _, iface := range userList.List() {
		userEl := iface.(map[string]interface{})
		user := icdv4.User{
			UserName: userEl["name"].(string),
			Password: userEl["password"].(string),
		}
		users = append(users, user)
	}
	return
}

type CsEntry struct {
	Name       string
	Password   string
	String     string
	Composed   string
	CertName   string
	CertBase64 string
	Hosts      []struct {
		HostName string `json:"hostname"`
		Port     int    `json:"port"`
	}
	Scheme       string
	QueryOptions map[string]interface{}
	Path         string
	Database     string
	BundleName   string
	BundleBase64 string
}

// IBM Cloud Databases
func FlattenConnectionStrings(cs []CsEntry) []map[string]interface{} {
	entries := make([]map[string]interface{}, len(cs), len(cs))
	for i, csEntry := range cs {
		l := map[string]interface{}{
			"name":         csEntry.Name,
			"password":     csEntry.Password,
			"composed":     csEntry.Composed,
			"certname":     csEntry.CertName,
			"certbase64":   csEntry.CertBase64,
			"queryoptions": csEntry.QueryOptions,
			"scheme":       csEntry.Scheme,
			"path":         csEntry.Path,
			"database":     csEntry.Database,
			"bundlename":   csEntry.BundleName,
			"bundlebase64": csEntry.BundleBase64,
		}
		hosts := csEntry.Hosts
		hostsList := make([]map[string]interface{}, len(hosts), len(hosts))
		for j, host := range hosts {
			z := map[string]interface{}{
				"hostname": host.HostName,
				"port":     strconv.Itoa(host.Port),
			}
			hostsList[j] = z
		}
		l["hosts"] = hostsList
		var queryOpts string
		if len(csEntry.QueryOptions) != 0 {
			queryOpts = "?"
			count := 0
			for k, v := range csEntry.QueryOptions {
				if count >= 1 {
					queryOpts = queryOpts + "&"
				}
				queryOpts = queryOpts + fmt.Sprintf("%v", k) + "=" + fmt.Sprintf("%v", v)
				count++
			}
		} else {
			queryOpts = ""
		}
		l["queryoptions"] = queryOpts
		entries[i] = l
	}

	return entries
}

func FlattenPhaseOneAttributes(vpn *datatypes.Network_Tunnel_Module_Context) []map[string]interface{} {
	phaseoneAttributesMap := make([]map[string]interface{}, 0, 1)
	phaseoneAttributes := make(map[string]interface{})
	phaseoneAttributes["authentication"] = *vpn.PhaseOneAuthentication
	phaseoneAttributes["encryption"] = *vpn.PhaseOneEncryption
	phaseoneAttributes["diffie_hellman_group"] = *vpn.PhaseOneDiffieHellmanGroup
	phaseoneAttributes["keylife"] = *vpn.PhaseOneKeylife
	phaseoneAttributesMap = append(phaseoneAttributesMap, phaseoneAttributes)
	return phaseoneAttributesMap
}

func FlattenPhaseTwoAttributes(vpn *datatypes.Network_Tunnel_Module_Context) []map[string]interface{} {
	phasetwoAttributesMap := make([]map[string]interface{}, 0, 1)
	phasetwoAttributes := make(map[string]interface{})
	phasetwoAttributes["authentication"] = *vpn.PhaseTwoAuthentication
	phasetwoAttributes["encryption"] = *vpn.PhaseTwoEncryption
	phasetwoAttributes["diffie_hellman_group"] = *vpn.PhaseTwoDiffieHellmanGroup
	phasetwoAttributes["keylife"] = *vpn.PhaseTwoKeylife
	phasetwoAttributesMap = append(phasetwoAttributesMap, phasetwoAttributes)
	return phasetwoAttributesMap
}

func FlattenaddressTranslation(vpn *datatypes.Network_Tunnel_Module_Context, fwID int) []map[string]interface{} {
	addressTranslationMap := make([]map[string]interface{}, 0, 1)
	addressTranslationAttributes := make(map[string]interface{})
	for _, networkAddressTranslation := range vpn.AddressTranslations {
		if *networkAddressTranslation.NetworkTunnelContext.Id == fwID {
			addressTranslationAttributes["remote_ip_adress"] = *networkAddressTranslation.CustomerIpAddress
			addressTranslationAttributes["internal_ip_adress"] = *networkAddressTranslation.InternalIpAddress
			addressTranslationAttributes["notes"] = *networkAddressTranslation.Notes
		}
	}
	addressTranslationMap = append(addressTranslationMap, addressTranslationAttributes)
	return addressTranslationMap
}

func FlattenremoteSubnet(vpn *datatypes.Network_Tunnel_Module_Context) []map[string]interface{} {
	remoteSubnetMap := make([]map[string]interface{}, 0, 1)
	remoteSubnetAttributes := make(map[string]interface{})
	for _, customerSubnet := range vpn.CustomerSubnets {
		remoteSubnetAttributes["remote_ip_adress"] = customerSubnet.NetworkIdentifier
		remoteSubnetAttributes["remote_ip_cidr"] = customerSubnet.Cidr
		remoteSubnetAttributes["account_id"] = customerSubnet.AccountId
	}
	remoteSubnetMap = append(remoteSubnetMap, remoteSubnetAttributes)
	return remoteSubnetMap
}

// IBM Cloud Databases
func ExpandAllowlist(allowList *schema.Set) (entries []clouddatabasesv5.AllowlistEntry) {
	entries = make([]clouddatabasesv5.AllowlistEntry, 0, len(allowList.List()))
	for _, iface := range allowList.List() {
		alItem := iface.(map[string]interface{})
		alEntry := &clouddatabasesv5.AllowlistEntry{
			Address:     core.StringPtr(alItem["address"].(string)),
			Description: core.StringPtr(alItem["description"].(string)),
		}
		entries = append(entries, *alEntry)
	}
	return
}

// IBM Cloud Databases
func FlattenAllowlist(allowlist []clouddatabasesv5.AllowlistEntry) []map[string]interface{} {
	entries := make([]map[string]interface{}, 0, len(allowlist))
	for _, ip := range allowlist {
		entry := map[string]interface{}{
			"address":     ip.Address,
			"description": ip.Description,
		}
		entries = append(entries, entry)
	}
	return entries
}

func ExpandPlatformOptions(deployment clouddatabasesv5.Deployment) []map[string]interface{} {
	pltOptions := make([]map[string]interface{}, 0, 1)
	pltOption := make(map[string]interface{})
	pltOption["disk_encryption_key_crn"] = deployment.PlatformOptions["disk_encryption_key_crn"]
	pltOption["backup_encryption_key_crn"] = deployment.PlatformOptions["backup_encryption_key_crn"]
	pltOptions = append(pltOptions, pltOption)
	return pltOptions
}

func expandStringMap(inVal interface{}) map[string]string {
	outVal := make(map[string]string)
	if inVal == nil {
		return outVal
	}
	for k, v := range inVal.(map[string]interface{}) {
		strValue := fmt.Sprintf("%v", v)
		outVal[k] = strValue
	}
	return outVal
}

// Cloud Internet Services
func ConvertTfToCisThreeVar(glbTfId string) (glbId string, zoneId string, cisId string, err error) {
	g := strings.SplitN(glbTfId, ":", 3)
	glbId = g[0]
	if len(g) > 2 {
		zoneId = g[1]
		cisId = g[2]
	} else {
		err = errors.New("cis_id or zone_id not passed")
		return
	}
	return
}
func ConvertCisToTfFourVar(firewallType string, ID string, ID2 string, cisID string) (buildID string) {
	if ID != "" {
		buildID = firewallType + ":" + ID + ":" + ID2 + ":" + cisID
	} else {
		buildID = ""
	}
	return
}
func ConvertTfToCisFourVar(TfID string) (firewallType string, ID string, zoneID string, cisID string, err error) {
	g := strings.SplitN(TfID, ":", 4)
	firewallType = g[0]
	if len(g) > 3 {
		ID = g[1]
		zoneID = g[2]
		cisID = g[3]
	} else {
		err = errors.New("Id or cis_id or zone_id not passed")
		return
	}
	return
}

// Cloud Internet Services
func ConvertCisToTfThreeVar(Id string, Id2 string, cisId string) (buildId string) {
	if Id != "" {
		buildId = Id + ":" + Id2 + ":" + cisId
	} else {
		buildId = ""
	}
	return
}

// Cloud Internet Services
func ConvertTfToCisTwoVarSlice(tfIds []string) (Ids []string, cisId string, err error) {
	for _, item := range tfIds {
		Id := strings.SplitN(item, ":", 2)
		if len(Id) < 2 {
			err = errors.New("cis_id not passed")
			return
		}
		Ids = append(Ids, Id[0])
		cisId = Id[1]
	}
	return
}

// Cloud Internet Services
func ConvertCisToTfTwoVarSlice(Ids []string, cisId string) (buildIds []string) {
	for _, Id := range Ids {
		buildIds = append(buildIds, Id+":"+cisId)
	}
	return
}

// Cloud Internet Services
func ConvertCisToTfTwoVar(Id string, cisId string) (buildId string) {
	if Id != "" {
		buildId = Id + ":" + cisId
	} else {
		buildId = ""
	}
	return
}

// Cloud Internet Services
func ConvertTftoCisTwoVar(tfId string) (Id string, cisId string, err error) {
	g := strings.SplitN(tfId, ":", 2)
	Id = g[0]
	if len(g) > 1 {
		cisId = g[1]
	} else {
		err = errors.New(" cis_id or zone_id not passed")
		return
	}
	return
}
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

var dnsTypeIntFields = []string{
	"algorithm",
	"key_tag",
	"type",
	"usage",
	"selector",
	"matching_type",
	"weight",
	"priority",
	"port",
	"long_degrees",
	"lat_degrees",
	"long_minutes",
	"lat_minutes",
	"protocol",
	"digest_type",
	"order",
	"preference",
}

var dnsTypeFloatFields = []string{
	"size",
	"altitude",
	"precision_horz",
	"precision_vert",
	"long_seconds",
	"lat_seconds",
}

// Cloud Internet Services
func TransformToIBMCISDnsData(recordType string, id string, value interface{}) (newValue interface{}, err error) {
	switch {
	case id == "flags":
		switch {
		case strings.ToUpper(recordType) == "SRV",
			strings.ToUpper(recordType) == "CAA",
			strings.ToUpper(recordType) == "DNSKEY":
			newValue, err = strconv.Atoi(value.(string))
		case strings.ToUpper(recordType) == "NAPTR":
			newValue, err = value.(string), nil
		}
	case stringInSlice(id, dnsTypeIntFields):
		newValue, err = strconv.Atoi(value.(string))
	case stringInSlice(id, dnsTypeFloatFields):
		newValue, err = strconv.ParseFloat(value.(string), 32)
	default:
		newValue, err = value.(string), nil
	}

	return
}

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func rcInstanceExists(resourceId string, resourceType string, meta interface{}) (bool, error) {
	// Check to see if Resource Manager instance exists
	rsConClient, err := meta.(conns.ClientSession).ResourceControllerAPI()
	if err != nil {
		return true, nil
	}
	exists := true
	instance, err := rsConClient.ResourceServiceInstance().GetInstance(resourceId)
	if err != nil {
		if strings.Contains(err.Error(), "Object not found") ||
			strings.Contains(err.Error(), "status code: 404") {
			exists = false
		} else {
			return true, fmt.Errorf("[ERROR] Error checking resource instance exists: %s", err)
		}
	} else {
		if strings.Contains(instance.State, "removed") {
			exists = false
		}
	}
	if exists {
		return true, nil
	}
	// Implement when pointer to terraform.State available
	// If rcInstance is now in removed state, set TF state to removed
	// s := *terraform.State
	// for _, r := range s.RootModule().Resources {
	//  if r.Type != resourceType {
	//      continue
	//  }
	//  if r.Primary.ID == resourceId {
	//      r.Primary.Set("status", "removed")
	//  }
	// }
	return false, nil
}

// Implement when pointer to terraform.State available
// func resourceInstanceExistsTf(resourceId string, resourceType string) bool {
//  // Check TF state to see if Cloud resource instance has already been removed
//  s := *terraform.State
//  for _, r := range s.RootModule().Resources {
//      if r.Type != resourceType {
//          continue
//      }
//      if r.Primary.ID == resourceId {
//          if strings.Contains(r.Primary.Attributes["status"], "removed") {
//              return false
//          }
//      }
//  }
//  return true
// }

// convert CRN to be url safe
func EscapeUrlParm(urlParm string) string {
	if strings.Contains(urlParm, "/") {
		newUrlParm := url.PathEscape(urlParm)
		return newUrlParm
	}
	return urlParm
}
func GetLocation(instance models.ServiceInstanceV2) string {
	region := instance.Crn.Region
	cName := instance.Crn.CName
	if cName == "bluemix" || cName == "staging" {
		return region
	} else {
		return cName + "-" + region
	}
}

type CRN struct {
	Scheme          string
	Version         string
	CName           string
	CType           string
	ServiceName     string
	Region          string
	ScopeType       string
	Scope           string
	ServiceInstance string
	ResourceType    string
	Resource        string
}

func Parse(s string) (CRN, error) {
	if s == "" {
		return CRN{}, nil
	}

	segments := strings.Split(s, crnSeparator)
	if len(segments) != 10 || segments[0] != crn {
		return CRN{}, ErrMalformedCRN
	}

	crn := CRN{
		Scheme:          segments[0],
		Version:         segments[1],
		CName:           segments[2],
		CType:           segments[3],
		ServiceName:     segments[4],
		Region:          segments[5],
		ServiceInstance: segments[7],
		ResourceType:    segments[8],
		Resource:        segments[9],
	}

	scopeSegments := segments[6]
	if scopeSegments != "" {
		if scopeSegments == "global" {
			crn.Scope = "global"
		} else {
			scopeParts := strings.Split(scopeSegments, scopeSeparator)
			if len(scopeParts) == 2 {
				crn.ScopeType, crn.Scope = scopeParts[0], scopeParts[1]
			} else {
				return CRN{}, ErrMalformedScope
			}
		}
	}

	return crn, nil
}
func GetLocationV2(instance rc.ResourceInstance) string {
	crn, err := Parse(*instance.CRN)
	if err != nil {
		log.Fatal(err)
	}
	region := crn.Region
	cName := crn.CName
	if cName == "bluemix" || cName == "staging" {
		return region
	} else {
		return cName + "-" + region
	}
}

func GetTags(d *schema.ResourceData, meta interface{}) error {
	resourceID := d.Id()
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPI()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}
	taggingResult, err := gtClient.Tags().GetTags(resourceID)
	if err != nil {
		return err
	}
	var taglist []string
	for _, item := range taggingResult.Items {
		taglist = append(taglist, item.Name)
	}
	d.Set("tags", FlattenStringList(taglist))
	return nil
}

// func UpdateTags(d *schema.ResourceData, meta interface{}) error {
// 	resourceID := d.Id()
// 	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPI()
// 	if err != nil {
// 		return fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
// 	}
// 	oldList, newList := d.GetChange("tags")
// 	if oldList == nil {
// 		oldList = new(schema.Set)
// 	}
// 	if newList == nil {
// 		newList = new(schema.Set)
// 	}
// 	olds := oldList.(*schema.Set)
// 	news := newList.(*schema.Set)
// 	removeInt := olds.Difference(news).List()
// 	addInt := news.Difference(olds).List()
// 	add := make([]string, len(addInt))
// 	for i, v := range addInt {
// 		add[i] = fmt.Sprint(v)
// 	}
// 	remove := make([]string, len(removeInt))
// 	for i, v := range removeInt {
// 		remove[i] = fmt.Sprint(v)
// 	}

// 	if len(add) > 0 {
// 		_, err := gtClient.Tags().AttachTags(resourceID, add)
// 		if err != nil {
// 			return fmt.Errorf("[ERROR] Error updating database tags %v : %s", add, err)
// 		}
// 	}
// 	if len(remove) > 0 {
// 		_, err := gtClient.Tags().DetachTags(resourceID, remove)
// 		if err != nil {
// 			return fmt.Errorf("[ERROR] Error detaching database tags %v: %s", remove, err)
// 		}
// 		for _, v := range remove {
// 			_, err := gtClient.Tags().DeleteTag(v)
// 			if err != nil {
// 				return fmt.Errorf("[ERROR] Error deleting database tag %v: %s", v, err)
// 			}
// 		}
// 	}
// 	return nil
// }

func GetGlobalTagsUsingCRN(meta interface{}, resourceID, resourceType, tagType string) (*schema.Set, error) {
	taggingResult, err := GetGlobalTagsUsingSearchAPI(meta, resourceID, resourceType, tagType)
	if err != nil {
		return nil, err
	}
	return taggingResult, nil
}

func GetTagsUsingResourceCRNFromTaggingApi(meta interface{}, resourceID, resourceType, tagType string) (*schema.Set, error) {
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return nil, err
	}
	accountID := userDetails.UserAccount
	ListTagsOptions := &globaltaggingv1.ListTagsOptions{}
	ListTagsOptions.AttachedTo = &resourceID
	if strings.HasPrefix(resourceType, "Softlayer_") {
		ListTagsOptions.Providers = []string{"ims"}
	}
	if len(tagType) > 0 {
		ListTagsOptions.TagType = PtrToString(tagType)
		if tagType == "service" {
			ListTagsOptions.AccountID = PtrToString(accountID)
		}
	}
	taggingResult, _, err := gtClient.ListTags(ListTagsOptions)
	if err != nil {
		return nil, err
	}
	var taglist []string
	for _, item := range taggingResult.Items {
		taglist = append(taglist, *item.Name)
	}
	return NewStringSet(ResourceIBMVPCHash, taglist), nil
}

func GetGlobalTagsUsingSearchAPI(meta interface{}, resourceID, resourceType, tagType string) (*schema.Set, error) {
	gsClient, err := meta.(conns.ClientSession).GlobalSearchAPIV2()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting global search client settings: %s", err)
	}
	options := globalsearchv2.SearchOptions{}
	var query string
	if strings.Contains(resourceType, "SoftLayer_") {
		query = fmt.Sprintf("doc.id:%s AND family:ims", resourceID)
		options.SetQuery(query)
	} else {
		query = fmt.Sprintf("crn:\"%s\"", resourceID)
		options.SetQuery(query)
	}
	if tagType == "service" {
		userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
		if err != nil {
			return nil, err
		}
		options.SetAccountID(userDetails.UserAccount)
	}
	options.SetFields([]string{"access_tags", "tags", "service_tags"})
	result, resp, err := gsClient.Search(&options)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error to query the tags for the resource: %s %s", err, resp)
	}
	var taglist []string
	var t interface{}
	if len(result.Items) > 0 {
		if tagType == "access" {
			t = result.Items[0].GetProperty("access_tags")
		} else if tagType == "service" {
			t = result.Items[0].GetProperty("service_tags")
		} else {
			t = result.Items[0].GetProperty("tags")
		}
		switch reflect.TypeOf(t).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(t)

			for i := 0; i < s.Len(); i++ {
				t := fmt.Sprintf("%s", (s.Index(i)))
				taglist = append(taglist, t)
			}
		}
	}
	return NewStringSet(ResourceIBMVPCHash, taglist), nil
}

func UpdateGlobalTagsUsingCRN(oldList, newList interface{}, meta interface{}, resourceID, resourceType, tagType string) error {
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	acctID := userDetails.UserAccount

	resources := []globaltaggingv1.Resource{}
	r := globaltaggingv1.Resource{ResourceID: PtrToString(resourceID), ResourceType: PtrToString(resourceType)}
	resources = append(resources, r)

	if oldList == nil {
		oldList = new(schema.Set)
	}
	if newList == nil {
		newList = new(schema.Set)
	}
	olds := oldList.(*schema.Set)
	news := newList.(*schema.Set)
	removeInt := olds.Difference(news).List()
	addInt := news.Difference(olds).List()
	add := make([]string, len(addInt))
	for i, v := range addInt {
		add[i] = fmt.Sprint(v)
	}
	remove := make([]string, len(removeInt))
	for i, v := range removeInt {
		remove[i] = fmt.Sprint(v)
	}

	if strings.TrimSpace(tagType) == "" || tagType == "user" {
		schematicTags := os.Getenv("IC_ENV_TAGS")
		var envTags []string
		if schematicTags != "" {
			envTags = strings.Split(schematicTags, ",")
			add = append(add, envTags...)
		}
	}

	if len(remove) > 0 {
		detachTagOptions := &globaltaggingv1.DetachTagOptions{}
		detachTagOptions.Resources = resources
		detachTagOptions.TagNames = remove
		if len(tagType) > 0 {
			detachTagOptions.TagType = PtrToString(tagType)
			if tagType == "service" {
				detachTagOptions.AccountID = PtrToString(acctID)
			}
		}
		results, fullResponse, err := gtClient.DetachTag(detachTagOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error detaching tags calling api %v: %s\n%s", remove, err, fullResponse)
		}
		if results != nil {
			errMap := make([]globaltaggingv1.TagResultsItem, 0)
			for _, res := range results.Results {
				if res.IsError != nil && *res.IsError {
					errMap = append(errMap, res)
				}
			}
			if len(errMap) > 0 {
				output, _ := json.MarshalIndent(errMap, "", "    ")
				return fmt.Errorf("[ERROR] Error detaching tag in results %v: %s\n%s", remove, string(output), fullResponse)
			}
		}
	}

	if len(add) > 0 {
		AttachTagOptions := &globaltaggingv1.AttachTagOptions{}
		AttachTagOptions.Resources = resources
		AttachTagOptions.TagNames = add
		if len(tagType) > 0 {
			AttachTagOptions.TagType = PtrToString(tagType)
			if tagType == "service" {
				AttachTagOptions.AccountID = PtrToString(acctID)
			}
		}

		_, resp, err := gtClient.AttachTag(AttachTagOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating database tags %v : %s\n%s", add, err, resp)
		}
		response, errored := WaitForTagsAvailable(meta, resourceID, resourceType, tagType, news, 30*time.Second)
		if errored != nil {
			log.Printf(`[ERROR] Error waiting for resource tags %s : %v
%v`, resourceID, errored, response)
		}
	}

	return nil
}

func WaitForTagsAvailable(meta interface{}, resourceID, resourceType, tagType string, desired *schema.Set, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for tag attachment (%s) to be successful.", resourceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"success", "error"},
		Refresh:    tagsRefreshFunc(meta, resourceID, resourceType, tagType, desired),
		Timeout:    timeout,
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func tagsRefreshFunc(meta interface{}, resourceID, resourceType, tagType string, desired *schema.Set) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		tags, err := GetGlobalTagsUsingCRN(meta, resourceID, resourceType, tagType)
		if err != nil {
			return tags, "error", fmt.Errorf("[ERROR] Error on get of resource tags (%s) tags: %s", resourceID, err)
		}
		if tags.Equal(desired) {
			return tags, "success", nil
		} else {
			return tags, "pending", nil
		}
	}
}

func ResourceIBMVPCHash(v interface{}) int {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s",
		strings.ToLower(v.(string))))
	return conns.String(buf.String())
}

// Use this function for attributes which only should be applied in resource creation time.
func ApplyOnce(k, o, n string, d *schema.ResourceData) bool {
	if len(d.Id()) == 0 {
		return false
	}
	return true
}

func ApplyOnlyOnce(k, o, n string, d *schema.ResourceData) bool {
	// For new resources, allow the first value to be set
	if len(d.Id()) == 0 {
		return false
	}

	// For existing resources, don't allow changes (keep the original value)
	if o == "" {
		return false
	}
	return true
}
func GetTagsUsingCRN(meta interface{}, resourceCRN string) (*schema.Set, error) {
	// Move the API to use globalsearch API instead of globalTags API due to rate limit
	taggingResult, err := GetGlobalTagsUsingSearchAPI(meta, resourceCRN, "", "user")
	if err != nil {
		return nil, err
	}
	return taggingResult, nil
}

func UpdateTagsUsingCRN(oldList, newList interface{}, meta interface{}, resourceCRN string) error {
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}
	if oldList == nil {
		oldList = new(schema.Set)
	}
	if newList == nil {
		newList = new(schema.Set)
	}
	olds := oldList.(*schema.Set)
	news := newList.(*schema.Set)
	removeInt := olds.Difference(news).List()
	addInt := news.Difference(olds).List()
	add := make([]string, len(addInt))
	for i, v := range addInt {
		add[i] = fmt.Sprint(v)
	}
	remove := make([]string, len(removeInt))
	for i, v := range removeInt {
		remove[i] = fmt.Sprint(v)
	}

	schematicTags := os.Getenv("IC_ENV_TAGS")
	var envTags []string
	if schematicTags != "" {
		envTags = strings.Split(schematicTags, ",")
		add = append(add, envTags...)
	}

	resources := []globaltaggingv1.Resource{}
	r := globaltaggingv1.Resource{ResourceID: &resourceCRN}
	resources = append(resources, r)

	if len(remove) > 0 {
		detachTagOptions := &globaltaggingv1.DetachTagOptions{}
		detachTagOptions.Resources = resources
		detachTagOptions.TagNames = remove

		results, fullResponse, err := gtClient.DetachTag(detachTagOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error detaching tags %v: %s", remove, err)
		}
		if results != nil {
			errMap := make([]globaltaggingv1.TagResultsItem, 0)
			for _, res := range results.Results {
				if res.IsError != nil && *res.IsError {
					errMap = append(errMap, res)
				}
			}
			if len(errMap) > 0 {
				output, _ := json.MarshalIndent(errMap, "", "    ")
				return fmt.Errorf("[ERROR] Error detaching tag %v: %s\n%s", remove, string(output), fullResponse)
			}
		}
		for _, v := range remove {
			delTagOptions := &globaltaggingv1.DeleteTagOptions{
				TagName: PtrToString(v),
			}
			results, fullResponse, err := gtClient.DeleteTag(delTagOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error deleting tag %v: %s\n%s", v, err, fullResponse)
			}

			if results != nil {
				errMap := make([]globaltaggingv1.DeleteTagResultsItem, 0)
				for _, res := range results.Results {
					if res.IsError != nil && *res.IsError {
						errMap = append(errMap, res)
					}
				}
				if len(errMap) > 0 {
					output, _ := json.MarshalIndent(errMap, "", "    ")
					return fmt.Errorf("[ERROR] Error deleting tag %s: %s\n%s", v, string(output), fullResponse)
				}
			}
		}
	}

	if len(add) > 0 {
		AttachTagOptions := &globaltaggingv1.AttachTagOptions{}
		AttachTagOptions.Resources = resources
		AttachTagOptions.TagNames = add
		results, fullResponse, err := gtClient.AttachTag(AttachTagOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating tags %v : %s", add, err)
		}
		if results != nil {
			errMap := make([]globaltaggingv1.TagResultsItem, 0)
			for _, res := range results.Results {
				if res.IsError != nil && *res.IsError {
					errMap = append(errMap, res)
				}
			}
			if len(errMap) > 0 {
				output, _ := json.MarshalIndent(errMap, "", "    ")
				return fmt.Errorf("Error while updating tag: %s - Full response: %s", string(output), fullResponse)
			}
		}
	}

	return nil
}

func GetBaseController(meta interface{}) (string, error) {
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return "", err
	}
	if userDetails != nil && userDetails.CloudName == "staging" {
		return stageBaseController, nil
	}
	return prodBaseController, nil
}

func FlattenSSLCiphers(ciphers []datatypes.Network_LBaaS_SSLCipher) *schema.Set {
	c := make([]string, len(ciphers))
	for i, v := range ciphers {
		c[i] = *v.Name
	}
	return NewStringSet(schema.HashString, c)
}

func ResourceTagsCustomizeDiff(diff *schema.ResourceDiff) error {

	if diff.Id() != "" && diff.HasChange("tags") {
		o, n := diff.GetChange("tags")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)
		removeInt := oldSet.Difference(newSet).List()
		addInt := newSet.Difference(oldSet).List()
		if v := os.Getenv("IC_ENV_TAGS"); v != "" {
			s := strings.Split(v, ",")
			if len(removeInt) == len(s) && len(addInt) == 0 {
				fmt.Println("Suppresing the TAG diff ")
				return diff.Clear("tags")
			}
		}
	}
	return nil
}
func ResourcePowerUserTagsCustomizeDiff(diff *schema.ResourceDiff) error {

	if diff.Id() != "" && diff.HasChange("pi_user_tags") {
		// power tags
		o, n := diff.GetChange("pi_user_tags")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)
		removeInt := oldSet.Difference(newSet).List()
		addInt := newSet.Difference(oldSet).List()
		if v := os.Getenv("IC_ENV_TAGS"); v != "" {
			s := strings.Split(v, ",")
			if len(removeInt) == len(s) && len(addInt) == 0 {
				fmt.Println("Suppresing the TAG diff ")
				return diff.Clear("pi_user_tags")
			}
		}
	}
	return nil
}
func OnlyInUpdateDiff(resources []string, diff *schema.ResourceDiff) error {
	for _, r := range resources {
		if diff.HasChange(r) && diff.Id() == "" {
			return fmt.Errorf("the %s can't be used at create time", r)
		}
	}
	return nil
}

func ResourceValidateAccessTags(diff *schema.ResourceDiff, meta interface{}) error {

	if value, ok := diff.GetOkExists("access_tags"); ok {
		tagSet := value.(*schema.Set)
		newTagList := tagSet.List()
		tagType := "access"
		gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
		if err != nil {
			return fmt.Errorf("Error getting global tagging client settings: %s", err)
		}

		listTagsOptions := &globaltaggingv1.ListTagsOptions{
			TagType: &tagType,
		}
		taggingResult, _, err := gtClient.ListTags(listTagsOptions)
		if err != nil {
			return err
		}
		var taglist []string
		for _, item := range taggingResult.Items {
			taglist = append(taglist, *item.Name)
		}
		existingAccessTags := NewStringSet(ResourceIBMVPCHash, taglist)
		errStatement := ""
		for _, tag := range newTagList {
			if !existingAccessTags.Contains(tag) {
				errStatement = errStatement + " " + tag.(string)
			}
		}
		if errStatement != "" {
			return fmt.Errorf("[ERROR] Error : Access tag(s) %s does not exist", errStatement)
		}
	}
	return nil
}

func ResourceIBMISLBPoolCookieValidate(diff *schema.ResourceDiff) error {
	_, sessionPersistenceTypeIntf := diff.GetChange(isLBPoolSessPersistenceType)
	_, sessionPersistenceCookieNameIntf := diff.GetChange(isLBPoolSessPersistenceAppCookieName)
	sessionPersistenceType := sessionPersistenceTypeIntf.(string)
	sessionPersistenceCookieName := sessionPersistenceCookieNameIntf.(string)

	if sessionPersistenceType == "app_cookie" {
		if sessionPersistenceCookieName == "" {
			return fmt.Errorf("Load Balancer Pool: %s is required for %s 'app_cookie'", isLBPoolSessPersistenceAppCookieName, isLBPoolSessPersistenceType)
		}
		if strings.HasPrefix(sessionPersistenceCookieName, "IBM") {
			return fmt.Errorf("Load Balancer Pool: %s starting with IBM are not allowed", isLBPoolSessPersistenceAppCookieName)
		}
	}

	if sessionPersistenceCookieName != "" && sessionPersistenceType != "app_cookie" {
		return fmt.Errorf("Load Balancer Pool: %s is only applicable for %s 'app_cookie'.", isLBPoolSessPersistenceAppCookieName, isLBPoolSessPersistenceType)
	}
	return nil
}

func ResourceVolumeAttachmentValidate(diff *schema.ResourceDiff) error {

	if volsintf, ok := diff.GetOk("volume_attachments"); ok {
		vols := volsintf.([]interface{})
		for volAttIdx := range vols {
			volumeid := "volume_attachments." + strconv.Itoa(volAttIdx) + "." + "volume"
			volumePrototype := "volume_attachments." + strconv.Itoa(volAttIdx) + "." + "volume_prototype"
			var volIdnterpolated = false
			var volumeIdFound = false
			if _, volumeIdFound = diff.GetOk(volumeid); !volumeIdFound {
				if !diff.NewValueKnown(volumeid) {
					volIdnterpolated = true
				}
			}
			_, volPrototypeFound := diff.GetOk(volumePrototype)

			if volPrototypeFound && (volumeIdFound || volIdnterpolated) {
				return fmt.Errorf("InstanceTemplate - volume_attachments[%d]: Cannot provide both 'volume' and 'volume_prototype' together.", volAttIdx)
			}
			if !volPrototypeFound && !volumeIdFound && !volIdnterpolated {
				return fmt.Errorf("InstanceTemplate - volume_attachments[%d]: Volume details missing. Provide either 'volume' or 'volume_prototype'.", volAttIdx)
			}
		}
	}

	return nil
}

func ResourceIPSecPolicyValidate(diff *schema.ResourceDiff) error {

	newEncAlgo := diff.Get("encryption_algorithm").(string)
	newAuthAlgo := diff.Get("authentication_algorithm").(string)
	if (newEncAlgo == "aes128gcm16" || newEncAlgo == "aes192gcm16" || newEncAlgo == "aes256gcm16") && newAuthAlgo != "disabled" {
		return fmt.Errorf("authentication_algorithm must be set to 'disabled' when the encryption_algorithm is either one of 'aes128gcm16', 'aes192gcm16', 'aes256gcm16'")
	}

	return nil
}

func ResourceVolumeValidate(diff *schema.ResourceDiff) error {

	if diff.Id() != "" && diff.HasChange("capacity") {
		o, n := diff.GetChange("capacity")
		old := int64(o.(int))
		new := int64(n.(int))
		if new < old {
			return fmt.Errorf("'%s' attribute has a constraint, it supports only expansion and can't be changed from %d to %d.", "capacity", old, new)
		}
	}

	profile := ""
	var capacity, iops int64
	if profileOk, ok := diff.GetOk("profile"); ok {
		profile = profileOk.(string)
	}
	if capacityOk, ok := diff.GetOk("capacity"); ok {
		capacity = int64(capacityOk.(int))
	}

	if capacity == int64(0) {
		capacity = int64(100)
	}
	if profile == "5iops-tier" && capacity > 9600 {
		return fmt.Errorf("'%s' storage block supports capacity up to %d.", profile, 9600)
	} else if profile == "10iops-tier" && capacity > 4800 {
		return fmt.Errorf("'%s' storage block supports capacity up to %d.", profile, 4800)
	}

	if iopsOk, ok := diff.GetOk("iops"); ok {
		iops = int64(iopsOk.(int))
	}

	if diff.HasChange("profile") {
		oldProfile, newProfile := diff.GetChange("profile")
		if oldProfile.(string) == "custom" || newProfile.(string) == "custom" {
			diff.ForceNew("profile")
		}
	}

	if profile != "custom" && profile != "sdp" {
		if iops != 0 && diff.NewValueKnown("iops") && diff.HasChange("iops") {
			return fmt.Errorf("VolumeError : iops is applicable for only custom/sdp volume profiles")
		}
	} else if profile != "sdp" {
		if capacity == 0 {
			capacity = int64(100)
		}
		if capacity >= 10 && capacity <= 39 {
			min := int64(100)
			max := int64(1000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 40 && capacity <= 79 {
			min := int64(100)
			max := int64(2000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 80 && capacity <= 99 {
			min := int64(100)
			max := int64(4000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 100 && capacity <= 499 {
			min := int64(100)
			max := int64(6000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 500 && capacity <= 999 {
			min := int64(100)
			max := int64(10000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 1000 && capacity <= 1999 {
			min := int64(100)
			max := int64(20000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 2000 && capacity <= 3999 {
			min := int64(200)
			max := int64(40000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 4000 && capacity <= 7999 {
			min := int64(300)
			max := int64(40000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 8000 && capacity <= 9999 {
			min := int64(500)
			max := int64(48000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
		if capacity >= 10000 && capacity <= 16000 {
			min := int64(1000)
			max := int64(48000)
			if !(iops >= min && iops <= max) {
				return fmt.Errorf("VolumeError : allowed iops value for capacity(%d) is [%d-%d] ", capacity, min, max)
			}
		}
	}
	return nil
}

func ResourceRouteModeValidate(diff *schema.ResourceDiff) error {

	var lbtype, lbprofile string
	if typeOk, ok := diff.GetOk(isLBType); ok {
		lbtype = typeOk.(string)
	}
	if profileOk, ok := diff.GetOk(isLBProfile); ok {
		lbprofile = profileOk.(string)
	}
	if rmOk, ok := diff.GetOk(isLBRouteMode); ok {
		routeMode := rmOk.(bool)

		if routeMode && lbtype != "private" {
			return fmt.Errorf("'type' must be 'private', at present public load balancers are not supported with route mode enabled.")
		}
		if routeMode && lbprofile != "network-fixed" {
			return fmt.Errorf("'profile' must be 'network-fixed', route mode is supported by private network load balancer.")
		}
	}

	return nil
}

func FlattenRoleData(object []iampolicymanagementv1.Role, roleType string) []map[string]string {
	var roles []map[string]string

	for _, item := range object {
		role := make(map[string]string)
		role["name"] = *item.DisplayName
		role["type"] = roleType
		role["description"] = *item.Description
		roles = append(roles, role)
	}
	return roles
}

func FlattenCustomRoleData(object []iampolicymanagementv1.CustomRole, roleType string) []map[string]string {
	var roles []map[string]string

	for _, item := range object {
		role := make(map[string]string)
		role["name"] = *item.DisplayName
		role["type"] = roleType
		role["description"] = *item.Description
		roles = append(roles, role)
	}
	return roles
}

func flattenActions(object []iampolicymanagementv1.Role) map[string]interface{} {
	actions := map[string]interface{}{
		"reader":      FlattenActionbyDisplayName("Reader", object),
		"manager":     FlattenActionbyDisplayName("Manager", object),
		"reader_plus": FlattenActionbyDisplayName("ReaderPlus", object),
		"writer":      FlattenActionbyDisplayName("Writer", object),
	}
	return actions
}

func FlattenActionbyDisplayName(displayName string, object []iampolicymanagementv1.Role) []string {
	var actionIDs []string
	for _, role := range object {
		if *role.DisplayName == displayName {
			actionIDs = role.Actions
		}
	}
	return actionIDs
}

func flattenCatalogRef(object schematics.CatalogInfo) map[string]interface{} {
	catalogRef := map[string]interface{}{
		"item_id":          object.ItemID,
		"item_name":        object.ItemName,
		"item_url":         object.ItemURL,
		"offering_version": object.OfferingVersion,
	}
	return catalogRef
}

// GetNext ...
func GetNext(next interface{}) string {
	if reflect.ValueOf(next).IsNil() {
		return ""
	}

	u, err := url.Parse(reflect.ValueOf(next).Elem().FieldByName("Href").Elem().String())
	if err != nil {
		return ""
	}

	q := u.Query()
	return q.Get("start")
}

// GetNextIAM ...
func GetNextIAM(next interface{}) string {
	if reflect.ValueOf(next).IsNil() {
		return ""
	}

	u, err := url.Parse(reflect.ValueOf(next).Elem().String())
	if err != nil {
		return ""
	}
	q := u.Query()
	return q.Get("pagetoken")
}

/* Return the default resource group */
func DefaultResourceGroup(meta interface{}) (string, error) {

	rMgtClient, err := meta.(conns.ClientSession).ResourceManagerV2API()
	if err != nil {
		return "", err
	}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return "", err
	}
	accountID := userDetails.UserAccount
	defaultGrp := true
	resourceGroupList := rg.ListResourceGroupsOptions{
		Default: &defaultGrp,
	}
	if accountID != "" {
		resourceGroupList.AccountID = &accountID
	}
	grpList, resp, err := rMgtClient.ListResourceGroups(&resourceGroupList)
	if err != nil || grpList == nil || grpList.Resources == nil {
		return "", fmt.Errorf("[ERROR] Error retrieving resource group: %s %s", err, resp)
	}
	if len(grpList.Resources) <= 0 {
		return "", fmt.Errorf("[ERROR] The default resource group could not be found. Make sure you have required permissions to access the resource group")
	}
	return *grpList.Resources[0].ID, nil
}

func FlattenKeyPolicies(policies []kp.Policy) []map[string]interface{} {
	policyMap := make([]map[string]interface{}, 0, 1)
	rotationMap := make([]map[string]interface{}, 0, 1)
	dualAuthMap := make([]map[string]interface{}, 0, 1)
	for _, policy := range policies {
		log.Println("Policy CRN Data =============>", policy.CRN)
		policyCRNData := strings.Split(policy.CRN, ":")
		policyInstance := map[string]interface{}{
			"id":               policyCRNData[9],
			"crn":              policy.CRN,
			"created_by":       policy.CreatedBy,
			"creation_date":    (*(policy.CreatedAt)).String(),
			"updated_by":       policy.UpdatedBy,
			"last_update_date": (*(policy.UpdatedAt)).String(),
		}
		if policy.Rotation != nil {
			policyInstance["interval_month"] = policy.Rotation.Interval
			policyInstance["enabled"] = *policy.Rotation.Enabled
			rotationMap = append(rotationMap, policyInstance)
		} else if policy.DualAuth != nil {
			policyInstance["enabled"] = *(policy.DualAuth.Enabled)
			dualAuthMap = append(dualAuthMap, policyInstance)
		}
	}
	tempMap := map[string]interface{}{
		"rotation":         rotationMap,
		"dual_auth_delete": dualAuthMap,
	}
	policyMap = append(policyMap, tempMap)
	return policyMap
}

func FlattenKeyIndividualPolicy(policy string, policies []kp.Policy) []map[string]interface{} {
	rotationMap := make([]map[string]interface{}, 0, 1)
	dualAuthMap := make([]map[string]interface{}, 0, 1)
	for _, policy := range policies {
		policyCRNData := strings.Split(policy.CRN, ":")
		policyInstance := map[string]interface{}{
			"id":               policyCRNData[9],
			"crn":              policy.CRN,
			"created_by":       policy.CreatedBy,
			"creation_date":    (*(policy.CreatedAt)).String(),
			"updated_by":       policy.UpdatedBy,
			"last_update_date": (*(policy.UpdatedAt)).String(),
		}
		if policy.Rotation != nil {
			policyInstance["interval_month"] = policy.Rotation.Interval
			policyInstance["enabled"] = *policy.Rotation.Enabled
			rotationMap = append(rotationMap, policyInstance)
		} else if policy.DualAuth != nil {
			policyInstance["enabled"] = *(policy.DualAuth.Enabled)
			dualAuthMap = append(dualAuthMap, policyInstance)
		}
	}
	if policy == "rotation" {
		return rotationMap
	} else if policy == "dual_auth_delete" {
		return dualAuthMap
	}
	return nil
}

func FlattenInstancePolicy(policyType string, policies []kp.InstancePolicy) []map[string]interface{} {
	dualAuthMap := make([]map[string]interface{}, 0, 1)
	rotationMap := make([]map[string]interface{}, 0, 1)
	metricsMap := make([]map[string]interface{}, 0, 1)
	keyCreateImportAccessMap := make([]map[string]interface{}, 0, 1)
	for _, policy := range policies {
		policyInstance := map[string]interface{}{
			"created_by":    policy.CreatedBy,
			"creation_date": (*policy.CreatedAt).String(),
			"updated_by":    policy.UpdatedBy,
			"last_updated":  (*policy.UpdatedAt).String(),
		}
		if policy.PolicyType == "dualAuthDelete" {
			policyInstance["enabled"] = policy.PolicyData.Enabled
			dualAuthMap = append(dualAuthMap, policyInstance)
		}
		if policy.PolicyType == "rotation" {
			policyInstance["enabled"] = policy.PolicyData.Enabled
			if policy.PolicyData.Attributes != nil {
				policyInstance["interval_month"] = policy.PolicyData.Attributes.IntervalMonth
			}
			rotationMap = append(rotationMap, policyInstance)
		}
		if policy.PolicyType == "metrics" {
			policyInstance["enabled"] = policy.PolicyData.Enabled
			metricsMap = append(metricsMap, policyInstance)
		}
		if policy.PolicyType == "keyCreateImportAccess" {
			if policy.PolicyData.Enabled != nil {
				policyInstance["enabled"] = *policy.PolicyData.Enabled
			}
			if policy.PolicyData.Attributes.CreateRootKey != nil {
				policyInstance["create_root_key"] = *policy.PolicyData.Attributes.CreateRootKey
			}
			if policy.PolicyData.Attributes.CreateStandardKey != nil {
				policyInstance["create_standard_key"] = *policy.PolicyData.Attributes.CreateStandardKey
			}
			if policy.PolicyData.Attributes.ImportRootKey != nil {
				policyInstance["import_root_key"] = *policy.PolicyData.Attributes.ImportRootKey
			}
			if policy.PolicyData.Attributes.ImportStandardKey != nil {
				policyInstance["import_standard_key"] = *policy.PolicyData.Attributes.ImportStandardKey
			}
			if policy.PolicyData.Attributes.EnforceToken != nil {
				policyInstance["enforce_token"] = *policy.PolicyData.Attributes.EnforceToken
			}

			keyCreateImportAccessMap = append(keyCreateImportAccessMap, policyInstance)
		}
	}
	if policyType == "rotation" {
		return rotationMap
	} else if policyType == "dual_auth_delete" {
		return dualAuthMap
	} else if policyType == "metrics" {
		return metricsMap
	} else if policyType == "key_create_import_access" {
		return keyCreateImportAccessMap
	}
	return nil
}
func FlattenKeyPoliciesKey(policies []kp.Policy) []map[string]interface{} {
	policyMap := make([]map[string]interface{}, 0, 1)
	rotationMap := make([]map[string]interface{}, 0, 1)
	dualAuthMap := make([]map[string]interface{}, 0, 1)
	for _, policy := range policies {
		policyInstance := map[string]interface{}{}
		if policy.Rotation != nil {
			policyInstance["interval_month"] = policy.Rotation.Interval
			policyInstance["enabled"] = *(policy.Rotation.Enabled)
			rotationMap = append(rotationMap, policyInstance)
		} else if policy.DualAuth != nil {
			policyInstance["enabled"] = *(policy.DualAuth.Enabled)
			dualAuthMap = append(dualAuthMap, policyInstance)
		}
	}
	tempMap := map[string]interface{}{
		"rotation":         rotationMap,
		"dual_auth_delete": dualAuthMap,
	}
	policyMap = append(policyMap, tempMap)
	return policyMap
}

// IgnoreSystemLabels returns non-IBM tag keys.
func IgnoreSystemLabels(labels map[string]string) map[string]string {
	result := make(map[string]string)

	for k, v := range labels {
		if (strings.HasPrefix(k, SystemIBMLabelPrefix) ||
			strings.HasPrefix(k, KubernetesLabelPrefix) ||
			strings.HasPrefix(k, K8sLabelPrefix)) &&
			!strings.Contains(k, "node-local-dns-enabled") {
			continue
		}

		result[k] = v
	}

	return result
}

// ExpandCosConfig ..
func ExpandCosConfig(cos []interface{}) *kubernetesserviceapiv1.COSBucket {
	if len(cos) == 0 || cos[0] == nil {
		return &kubernetesserviceapiv1.COSBucket{}
	}
	in := cos[0].(map[string]interface{})
	obj := &kubernetesserviceapiv1.COSBucket{
		Bucket:   PtrToString(in["bucket"].(string)),
		Endpoint: PtrToString(in["endpoint"].(string)),
		Region:   PtrToString(in["region"].(string)),
	}
	return obj
}

// expandCosCredentials ..
func ExpandCosCredentials(cos []interface{}) *kubernetesserviceapiv1.COSAuthorization {
	if len(cos) == 0 || cos[0] == nil {
		return &kubernetesserviceapiv1.COSAuthorization{}
	}
	in := cos[0].(map[string]interface{})
	obj := &kubernetesserviceapiv1.COSAuthorization{
		AccessKeyID:     PtrToString(in["access_key-id"].(string)),
		SecretAccessKey: PtrToString(in["secret_access_key"].(string)),
	}
	return obj
}
func FlattenNlbConfigs(nlbData []containerv2.NlbVPCListConfig) []map[string]interface{} {
	nlbConfigList := make([]map[string]interface{}, 0)
	for _, n := range nlbData {
		nlbConfig := make(map[string]interface{})
		nlbConfig["secret_name"] = n.SecretName
		nlbConfig["secret_status"] = n.SecretStatus
		c := n.Nlb
		nlbConfig["cluster"] = c.Cluster
		nlbConfig["dns_type"] = c.DnsType
		nlbConfig["lb_hostname"] = c.LbHostname
		nlbConfig["nlb_ips"] = c.NlbIPArray
		nlbConfig["nlb_sub_domain"] = c.NlbSubdomain
		nlbConfig["secret_namespace"] = c.SecretNamespace
		nlbConfig["type"] = c.Type
		nlbConfigList = append(nlbConfigList, nlbConfig)
	}

	return nlbConfigList
}

func FlattenOpaqueSecret(fields containerv2.Fields) []map[string]interface{} {
	flattenedOpaqueSecret := make([]map[string]interface{}, 0)

	for _, field := range fields {
		opaqueSecretField := map[string]interface{}{
			"name":                   field.Name,
			"crn":                    field.CRN,
			"expires_on":             field.ExpiresOn,
			"last_updated_timestamp": field.LastUpdatedTimestamp,
		}
		flattenedOpaqueSecret = append(flattenedOpaqueSecret, opaqueSecretField)
	}

	return flattenedOpaqueSecret
}

// flatten the provided key-value pairs
func FlattenKeyValues(keyValues []interface{}) map[string]string {
	labels := make(map[string]string)
	for _, v := range keyValues {
		parts := strings.Split(v.(string), ":")
		if len(parts) != 2 {
			log.Fatal("Entered key-value " + v.(string) + "is in incorrect format.")
		}
		labels[parts[0]] = parts[1]
	}

	return labels
}

func FlattenSatelliteZones(zones *schema.Set) []string {
	zoneList := make([]string, zones.Len())
	for i, v := range zones.List() {
		zoneList[i] = fmt.Sprint(v)
	}

	return zoneList
}

// error object
type ServiceErrorResponse struct {
	Message    string
	StatusCode int
	Result     interface{}
	Error      error
}

func BeautifyError(err error, response *core.DetailedResponse) *ServiceErrorResponse {
	var (
		statusCode int
		result     interface{}
	)
	if response != nil {
		statusCode = response.StatusCode
		result = response.Result
	}
	return &ServiceErrorResponse{
		Message:    err.Error(),
		StatusCode: statusCode,
		Result:     result,
		Error:      err,
	}
}

func (response *ServiceErrorResponse) String() string {
	output, err := json.MarshalIndent(response, "", "    ")
	if err == nil {
		return fmt.Sprintf("%+v\n", string(output))
	}
	return fmt.Sprintf("Error : %#v", response)
}

// IAM Policy Management
func GetResourceAttribute(name string, r iampolicymanagementv1.PolicyResource) *string {
	for _, a := range r.Attributes {
		if *a.Name == name {
			return a.Value
		}
	}
	return core.StringPtr("")
}

func GetV2PolicyResourceAttribute(key string, r iampolicymanagementv1.V2PolicyResource) string {
	for _, a := range r.Attributes {
		if *a.Key == key {
			if *a.Operator == "stringExists" && a.Value == true {
				return fmt.Sprint("*")
			} else if *a.Operator == "stringMatch" || *a.Operator == "stringEquals" {
				return a.Value.(string)
			}
		}
	}
	return *core.StringPtr("")
}

func GetSubjectAttribute(name string, s iampolicymanagementv1.PolicySubject) *string {
	for _, a := range s.Attributes {
		if *a.Name == name {
			return a.Value
		}
	}
	return core.StringPtr("")
}

func GetV2PolicySubjectAttribute(key string, s iampolicymanagementv1.V2PolicySubject) interface{} {
	for _, a := range s.Attributes {
		if *a.Key == key &&
			(*a.Operator == "stringMatch" ||
				*a.Operator == "stringEquals") {
			return a.Value
		}
	}
	return interface{}(core.StringPtr(""))
}

func SetResourceAttribute(name *string, value *string, r []iampolicymanagementv1.ResourceAttribute) []iampolicymanagementv1.ResourceAttribute {
	for _, a := range r {
		if *a.Name == *name {
			a.Value = value
			return r
		}
	}
	r = append(r, iampolicymanagementv1.ResourceAttribute{
		Name:     name,
		Value:    value,
		Operator: core.StringPtr("stringEquals"),
	})
	return r
}

func SetV2PolicyResourceAttribute(key *string, value *string, r []iampolicymanagementv1.V2PolicyResourceAttribute) []iampolicymanagementv1.V2PolicyResourceAttribute {
	for _, a := range r {
		if *a.Key == *key {
			a.Value = value
			return r
		}
	}
	r = append(r, iampolicymanagementv1.V2PolicyResourceAttribute{
		Key:      key,
		Value:    value,
		Operator: core.StringPtr("stringEquals"),
	})
	return r
}

func FindRoleByName(supported []iampolicymanagementv1.PolicyRole, name string) (iampolicymanagementv1.PolicyRole, error) {
	for _, role := range supported {
		if role.DisplayName != nil {
			if *role.DisplayName == name {
				role.DisplayName = nil
				return role, nil
			}
		}
	}
	if name == "NONE" {
		name := "NONE"
		r := iampolicymanagementv1.PolicyRole{
			DisplayName: &name,
			RoleID:      &name,
		}
		return r, nil
	}
	supportedRoles := getSupportedRolesStr(supported)
	return iampolicymanagementv1.PolicyRole{}, bmxerror.New("RoleDoesnotExist",
		fmt.Sprintf("%s was not found. Valid roles are %s", name, supportedRoles))

}

func FindRoleByCRN(supported []iampolicymanagementv1.PolicyRole, crn string) (iampolicymanagementv1.PolicyRole, error) {
	for _, role := range supported {
		if role.RoleID != nil {
			if *role.RoleID == crn {
				role.RoleID = nil
				return role, nil
			}
		}
	}
	supportedRoles := getSupportedRolesStr(supported)
	return iampolicymanagementv1.PolicyRole{}, bmxerror.New("RoleDoesnotExist",
		fmt.Sprintf("%s was not found. Valid roles are %s", crn, supportedRoles))

}

func getSupportedRolesStr(supported []iampolicymanagementv1.PolicyRole) string {
	rolesStr := "NONE, "
	for index, role := range supported {
		if index != 0 {
			rolesStr += ", "
		}
		if role.DisplayName != nil {
			rolesStr += *role.DisplayName
		}
	}
	return rolesStr
}

func GetRolesFromRoleNames(roleNames []string, roles []iampolicymanagementv1.PolicyRole) ([]iampolicymanagementv1.PolicyRole, error) {

	filteredRoles := []iampolicymanagementv1.PolicyRole{}
	for _, roleName := range roleNames {
		role, err := FindRoleByName(roles, roleName)
		if err != nil {
			return []iampolicymanagementv1.PolicyRole{}, err
		}
		role.DisplayName = nil
		filteredRoles = append(filteredRoles, role)
	}
	return filteredRoles, nil
}

func MapRoleListToPolicyRoles(roleList iampolicymanagementv1.RoleCollection) []iampolicymanagementv1.PolicyRole {
	var policyRoles []iampolicymanagementv1.PolicyRole
	for _, customRole := range roleList.CustomRoles {
		newPolicyRole := iampolicymanagementv1.PolicyRole{
			DisplayName: customRole.DisplayName,
			RoleID:      customRole.CRN,
		}
		policyRoles = append(policyRoles, newPolicyRole)
	}
	for _, serviceRole := range roleList.ServiceRoles {
		newPolicyRole := iampolicymanagementv1.PolicyRole{
			DisplayName: serviceRole.DisplayName,
			RoleID:      serviceRole.CRN,
		}
		policyRoles = append(policyRoles, newPolicyRole)
	}
	for _, systemRole := range roleList.SystemRoles {
		newPolicyRole := iampolicymanagementv1.PolicyRole{
			DisplayName: systemRole.DisplayName,
			RoleID:      systemRole.CRN,
		}
		policyRoles = append(policyRoles, newPolicyRole)
	}
	return policyRoles
}

func MapPolicyRolesToRoles(policyRoles []iampolicymanagementv1.PolicyRole) []iampolicymanagementv1.Roles {
	roles := make([]iampolicymanagementv1.Roles, len(policyRoles))
	for i, policyRole := range policyRoles {
		roles[i].RoleID = policyRole.RoleID
	}
	return roles
}

func MapRolesToPolicyRoles(roles []iampolicymanagementv1.Roles) []iampolicymanagementv1.PolicyRole {
	policyRoles := make([]iampolicymanagementv1.PolicyRole, len(roles))
	for i, role := range roles {
		policyRoles[i].RoleID = role.RoleID
	}
	return policyRoles
}

func GetRoleNamesFromPolicyResponse(policy iampolicymanagementv1.V2PolicyTemplateMetaData, d *schema.ResourceData, meta interface{}) ([]string, error) {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return []string{}, err
	}

	controlResponse := policy.Control.(*iampolicymanagementv1.ControlResponse)
	policyRoles := MapRolesToPolicyRoles(controlResponse.Grant.Roles)
	resourceAttributes := policy.Resource.Attributes
	subjectAttributes := policy.Subject.Attributes

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return []string{}, err
	}

	var (
		serviceName       string
		sourceServiceName string
		resourceType      string
		serviceGroupID    string
	)

	for _, a := range subjectAttributes {
		if *a.Key == "serviceName" &&
			(*a.Operator == "stringMatch" ||
				*a.Operator == "stringEquals") {
			sourceServiceName = a.Value.(string)
		}
	}

	for _, a := range resourceAttributes {
		if *a.Key == "serviceName" &&
			(*a.Operator == "stringMatch" ||
				*a.Operator == "stringEquals") {
			serviceName = a.Value.(string)
		}
		if *a.Key == "resourceType" &&
			(*a.Operator == "stringMatch" ||
				*a.Operator == "stringEquals") {
			resourceType = a.Value.(string)
		}
		if *a.Key == "service_group_id" &&
			(*a.Operator == "stringMatch" ||
				*a.Operator == "stringEquals") {
			serviceGroupID = a.Value.(string)
		}
	}

	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
		AccountID: &userDetails.UserAccount,
	}

	var isAccountManagementPolicy bool
	if accountManagement, ok := d.GetOk("account_management"); ok {
		isAccountManagementPolicy = accountManagement.(bool)
	}

	if serviceName == "" && resourceType == "resource-group" {
		serviceName = "resource-controller"
	}

	if serviceName == "" && // no specific service specified
		!isAccountManagementPolicy && // not all account management services
		resourceType != "resource-group" && // not to a resource group
		serviceGroupID == "" {
		listRoleOptions.ServiceName = core.StringPtr("alliamserviceroles")
	}

	if serviceName != "" {
		listRoleOptions.ServiceName = &serviceName
	}

	if serviceGroupID != "" {
		listRoleOptions.ServiceGroupID = &serviceGroupID
	}

	if sourceServiceName != "" {
		listRoleOptions.SourceServiceName = &sourceServiceName
	}

	if *policy.Type != "" {
		listRoleOptions.PolicyType = policy.Type
	}

	roleList, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)

	if err != nil {
		return []string{}, err
	}
	roles := MapRoleListToPolicyRoles(*roleList)
	roleNames := []string{}
	for _, role := range policyRoles {
		role, err := FindRoleByCRN(roles, *role.RoleID)
		if err != nil {
			return []string{}, err
		}
		roleNames = append(roleNames, *role.DisplayName)
	}

	return roleNames, nil
}

func GeneratePolicyOptions(d *schema.ResourceData, meta interface{}) (iampolicymanagementv1.CreatePolicyOptions, error) {

	var serviceName string
	var resourceType string
	var serviceGroupID string
	resourceAttributes := []iampolicymanagementv1.ResourceAttribute{}

	if res, ok := d.GetOk("resources"); ok {
		resources := res.([]interface{})
		for _, resource := range resources {
			r, _ := resource.(map[string]interface{})

			if r, ok := r["service"]; ok && r != nil {
				serviceName = r.(string)
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("serviceName"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["service_group_id"]; ok && r != nil {
				serviceGroupID = r.(string)
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("service_group_id"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource_instance_id"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("serviceInstance"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["region"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("region"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource_type"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("resourceType"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("resource"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource_group_id"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("resourceGroupId"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["service_type"]; ok && r != nil {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.ResourceAttribute{
						Name:     core.StringPtr("serviceType"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["attributes"]; ok {
				for k, v := range r.(map[string]interface{}) {
					resourceAttributes = SetResourceAttribute(core.StringPtr(k), core.StringPtr(v.(string)), resourceAttributes)
				}
			}
		}
	}
	if r, ok := d.GetOk("resource_attributes"); ok {
		for _, attribute := range r.(*schema.Set).List() {
			a := attribute.(map[string]interface{})
			name := a["name"].(string)
			value := a["value"].(string)
			operator := a["operator"].(string)
			if name == "serviceName" {
				serviceName = value
			}
			if name == "service_group_id" {
				serviceGroupID = value
			}
			at := iampolicymanagementv1.ResourceAttribute{
				Name:     &name,
				Value:    &value,
				Operator: &operator,
			}
			resourceAttributes = append(resourceAttributes, at)
		}
	}

	var serviceTypeResourceAttribute iampolicymanagementv1.ResourceAttribute

	if d.Get("account_management").(bool) {
		serviceTypeResourceAttribute = iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("serviceType"),
			Value:    core.StringPtr("platform_service"),
			Operator: core.StringPtr("stringEquals"),
		}
		resourceAttributes = append(resourceAttributes, serviceTypeResourceAttribute)
	}

	if len(resourceAttributes) == 0 {
		serviceTypeResourceAttribute = iampolicymanagementv1.ResourceAttribute{
			Name:     core.StringPtr("serviceType"),
			Value:    core.StringPtr("service"),
			Operator: core.StringPtr("stringEquals"),
		}
		resourceAttributes = append(resourceAttributes, serviceTypeResourceAttribute)
	}

	policyResources := iampolicymanagementv1.PolicyResource{
		Attributes: resourceAttributes,
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return iampolicymanagementv1.CreatePolicyOptions{}, err
	}

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()

	if err != nil {
		return iampolicymanagementv1.CreatePolicyOptions{}, err
	}

	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
		AccountID: &userDetails.UserAccount,
	}
	if serviceName == "" && // no specific service specified
		!d.Get("account_management").(bool) && // not all account management services
		resourceType != "resource-group" && // not to a resource group
		serviceGroupID == "" { // service_group_id and service is mutually exclusive
		listRoleOptions.ServiceName = core.StringPtr("alliamserviceroles")
	}
	if serviceName != "" {
		listRoleOptions.ServiceName = &serviceName
	}
	if serviceGroupID != "" {
		listRoleOptions.ServiceGroupID = &serviceGroupID
	}

	roleList, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
	if err != nil {
		return iampolicymanagementv1.CreatePolicyOptions{}, err
	}

	roles := MapRoleListToPolicyRoles(*roleList)
	policyRoles, err := GetRolesFromRoleNames(ExpandStringList(d.Get("roles").([]interface{})), roles)
	if err != nil {
		return iampolicymanagementv1.CreatePolicyOptions{}, err
	}

	return iampolicymanagementv1.CreatePolicyOptions{Roles: policyRoles, Resources: []iampolicymanagementv1.PolicyResource{policyResources}}, nil
}

func GenerateV2PolicyOptions(d *schema.ResourceData, meta interface{}) (iampolicymanagementv1.CreateV2PolicyOptions, error) {

	var serviceName string
	var resourceType string
	var serviceGroupID string
	resourceAttributes := []iampolicymanagementv1.V2PolicyResourceAttribute{}

	if res, ok := d.GetOk("resources"); ok {
		resources := res.([]interface{})
		for _, resource := range resources {
			r, _ := resource.(map[string]interface{})

			if r, ok := r["service"]; ok && r != nil {
				serviceName = r.(string)
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("serviceName"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["service_group_id"]; ok && r != nil {
				serviceGroupID = r.(string)
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("service_group_id"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource_instance_id"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("serviceInstance"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["region"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("region"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource_type"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("resourceType"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("resource"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["resource_group_id"]; ok {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("resourceGroupId"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["service_type"]; ok && r != nil {
				if r.(string) != "" {
					resourceAttr := iampolicymanagementv1.V2PolicyResourceAttribute{
						Key:      core.StringPtr("serviceType"),
						Value:    core.StringPtr(r.(string)),
						Operator: core.StringPtr("stringEquals"),
					}
					resourceAttributes = append(resourceAttributes, resourceAttr)
				}
			}

			if r, ok := r["attributes"]; ok {
				for k, v := range r.(map[string]interface{}) {
					resourceAttributes = SetV2PolicyResourceAttribute(core.StringPtr(k), core.StringPtr(v.(string)), resourceAttributes)
				}
			}
		}
	}
	if r, ok := d.GetOk("resource_attributes"); ok {
		for _, attribute := range r.(*schema.Set).List() {
			a := attribute.(map[string]interface{})
			name := a["name"].(string)
			value := a["value"].(string)
			operator := a["operator"].(string)
			if name == "serviceName" {
				serviceName = value
			}
			if name == "service_group_id" {
				serviceGroupID = value
			}
			at := iampolicymanagementv1.V2PolicyResourceAttribute{
				Key:      &name,
				Value:    &value,
				Operator: &operator,
			}
			resourceAttributes = append(resourceAttributes, at)
		}
	}

	var serviceTypeResourceAttribute iampolicymanagementv1.V2PolicyResourceAttribute

	if d.Get("account_management").(bool) {
		serviceTypeResourceAttribute = iampolicymanagementv1.V2PolicyResourceAttribute{
			Key:      core.StringPtr("serviceType"),
			Value:    core.StringPtr("platform_service"),
			Operator: core.StringPtr("stringEquals"),
		}
		resourceAttributes = append(resourceAttributes, serviceTypeResourceAttribute)
	}

	if len(resourceAttributes) == 0 {
		serviceTypeResourceAttribute = iampolicymanagementv1.V2PolicyResourceAttribute{
			Key:      core.StringPtr("serviceType"),
			Value:    core.StringPtr("service"),
			Operator: core.StringPtr("stringEquals"),
		}
		resourceAttributes = append(resourceAttributes, serviceTypeResourceAttribute)
	}

	policyResource := iampolicymanagementv1.V2PolicyResource{
		Attributes: resourceAttributes,
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return iampolicymanagementv1.CreateV2PolicyOptions{}, err
	}

	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()

	if err != nil {
		return iampolicymanagementv1.CreateV2PolicyOptions{}, err
	}

	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
		AccountID: &userDetails.UserAccount,
	}

	if serviceName == "" && // no specific service specified
		!d.Get("account_management").(bool) && // not all account management services
		resourceType != "resource-group" && // not to a resource group
		serviceGroupID == "" {
		listRoleOptions.ServiceName = core.StringPtr("alliamserviceroles")
	}

	if serviceName != "" {
		listRoleOptions.ServiceName = &serviceName
	}

	if serviceGroupID != "" {
		listRoleOptions.ServiceGroupID = &serviceGroupID
	}

	roleList, _, err := iamPolicyManagementClient.ListRoles(listRoleOptions)
	if err != nil {
		return iampolicymanagementv1.CreateV2PolicyOptions{}, err
	}

	roles := MapRoleListToPolicyRoles(*roleList)
	policyRoles, err := GetRolesFromRoleNames(ExpandStringList(d.Get("roles").([]interface{})), roles)
	if err != nil {
		return iampolicymanagementv1.CreateV2PolicyOptions{}, err
	}
	policyGrant := &iampolicymanagementv1.Grant{
		Roles: MapPolicyRolesToRoles(policyRoles),
	}
	policyControl := &iampolicymanagementv1.Control{
		Grant: policyGrant,
	}

	return iampolicymanagementv1.CreateV2PolicyOptions{Control: policyControl, Resource: &policyResource}, nil
}

func generatePolicyRuleCondition(c map[string]interface{}) iampolicymanagementv1.RuleAttribute {
	key := c["key"].(string)
	operator := c["operator"].(string)
	r := iampolicymanagementv1.RuleAttribute{
		Key:      &key,
		Operator: &operator,
	}

	interfaceValues := c["value"].([]interface{})
	values := make([]string, len(interfaceValues))
	for i, v := range interfaceValues {
		values[i] = fmt.Sprint(v)
	}

	if len(values) > 1 {
		r.Value = &values
	} else if operator == "stringExists" && values[0] == "true" {
		r.Value = true
	} else if operator == "stringExists" && values[0] == "false" {
		r.Value = false
	} else {
		r.Value = &values[0]
	}
	return r
}

func GeneratePolicyRule(d *schema.ResourceData, ruleConditions interface{}) *iampolicymanagementv1.V2PolicyRule {
	conditions := []iampolicymanagementv1.NestedConditionIntf{}

	for _, ruleCondition := range ruleConditions.(*schema.Set).List() {
		rc := ruleCondition.(map[string]interface{})
		con := rc["conditions"].([]interface{})
		if len(con) > 0 {
			nestedConditions := []iampolicymanagementv1.RuleAttribute{}
			for _, nc := range con {
				nestedConditions = append(nestedConditions, generatePolicyRuleCondition(nc.(map[string]interface{})))
			}
			nestedCondition := &iampolicymanagementv1.NestedCondition{}
			nestedConditionsOperator := rc["operator"].(string)
			nestedCondition.Operator = &nestedConditionsOperator
			nestedCondition.Conditions = nestedConditions
			conditions = append(conditions, nestedCondition)
		} else {
			ruleAttribute := generatePolicyRuleCondition(rc)
			nestedCondition := &iampolicymanagementv1.NestedCondition{
				Key:      ruleAttribute.Key,
				Operator: ruleAttribute.Operator,
				Value:    ruleAttribute.Value,
			}
			conditions = append(conditions, nestedCondition)
		}
	}
	rule := new(iampolicymanagementv1.V2PolicyRule)
	if len(conditions) == 1 {
		ruleCondition := conditions[0].(*iampolicymanagementv1.NestedCondition)
		rule.Key = ruleCondition.Key
		rule.Operator = ruleCondition.Operator
		rule.Value = ruleCondition.Value
	} else {
		ruleOperator := d.Get("rule_operator").(string)
		rule.Operator = &ruleOperator
		rule.Conditions = conditions
	}

	return rule
}

func SetTags(d *schema.ResourceData) []iampolicymanagementv1.ResourceTag {
	resourceAttributes := []iampolicymanagementv1.ResourceTag{}
	if r, ok := d.GetOk("resource_tags"); ok {
		for _, attribute := range r.(*schema.Set).List() {
			a := attribute.(map[string]interface{})
			name := a["name"].(string)
			value := a["value"].(string)
			operator := a["operator"].(string)
			tag := iampolicymanagementv1.ResourceTag{
				Name:     &name,
				Value:    &value,
				Operator: &operator,
			}
			resourceAttributes = append(resourceAttributes, tag)
		}
		return resourceAttributes
	}
	return []iampolicymanagementv1.ResourceTag{}
}

func SetV2PolicyTags(d *schema.ResourceData) []iampolicymanagementv1.V2PolicyResourceTag {
	resourceAttributes := []iampolicymanagementv1.V2PolicyResourceTag{}
	if r, ok := d.GetOk("resource_tags"); ok {
		for _, attribute := range r.(*schema.Set).List() {
			a := attribute.(map[string]interface{})
			name := a["name"].(string)
			value := a["value"].(string)
			operator := a["operator"].(string)
			tag := iampolicymanagementv1.V2PolicyResourceTag{
				Key:      &name,
				Value:    &value,
				Operator: &operator,
			}
			resourceAttributes = append(resourceAttributes, tag)
		}
		return resourceAttributes
	}
	return []iampolicymanagementv1.V2PolicyResourceTag{}
}

func GetIBMUniqueId(accountID, userEmail string, meta interface{}) (string, error) {
	userManagement, err := meta.(conns.ClientSession).UserManagementAPI()
	if err != nil {
		return "", err
	}
	client := userManagement.UserInvite()
	res, err := client.ListUsers(accountID)
	if err != nil {
		return "", err
	}
	for _, userInfo := range res {
		//handling case-sensitivity in userEmail
		if strings.ToLower(userInfo.Email) == strings.ToLower(userEmail) {
			return userInfo.IamID, nil
		}
	}
	return "", fmt.Errorf("User %s is not found under account %s", userEmail, accountID)
}

func ImmutableResourceCustomizeDiff(resourceList []string, diff *schema.ResourceDiff) error {

	sateLocZone := "managed_from"
	for _, rName := range resourceList {
		if diff.Id() != "" && diff.HasChange(rName) && rName != sateLocZone {
			return fmt.Errorf("'%s' attribute is immutable and can't be changed", rName)
		} else if diff.Id() != "" && diff.HasChange(rName) && rName == sateLocZone {
			o, n := diff.GetChange(rName)
			old := o.(string)
			new := n.(string)
			if len(old) >= 3 && len(new) >= 3 {
				if old[0:3] != new[0:3] {
					return fmt.Errorf("'%s' attribute is immutable and can't be changed from %s to %s", rName, old, new)
				}
			}
		}
	}
	return nil
}

func FlattenSatelliteWorkerPoolZones(zones *schema.Set) []kubernetesserviceapiv1.SatelliteCreateWorkerPoolZone {
	zoneList := make([]kubernetesserviceapiv1.SatelliteCreateWorkerPoolZone, zones.Len())
	for i, v := range zones.List() {
		data := v.(map[string]interface{})
		if v, ok := data["id"]; ok && v.(string) != "" {
			zoneList[i].ID = sl.String(v.(string))
		}
	}

	return zoneList
}

func FlattenSatelliteWorkerPools(list []kubernetesserviceapiv1.GetWorkerPoolResponse) []map[string]interface{} {
	workerPools := make([]map[string]interface{}, len(list))
	for i, workerPool := range list {
		l := map[string]interface{}{
			"id":                         *workerPool.ID,
			"name":                       *workerPool.PoolName,
			"isolation":                  *workerPool.Isolation,
			"flavour":                    *workerPool.Flavor,
			"size_per_zone":              *workerPool.WorkerCount,
			"state":                      *workerPool.Lifecycle.ActualState,
			"default_worker_pool_labels": workerPool.Labels,
			"host_labels":                workerPool.HostLabels,
		}
		zones := workerPool.Zones
		zonesConfig := make([]map[string]interface{}, len(zones))
		for j, zone := range zones {
			z := map[string]interface{}{
				"zone":         *zone.ID,
				"worker_count": int(*zone.WorkerCount),
			}
			zonesConfig[j] = z
		}
		l["zones"] = zonesConfig
		workerPools[i] = l
	}

	return workerPools
}

func FlattenSatelliteHosts(hostList []kubernetesserviceapiv1.MultishiftQueueNode) []map[string]interface{} {
	hosts := make([]map[string]interface{}, len(hostList))
	for i, host := range hostList {
		l := map[string]interface{}{
			"host_id":      *host.ID,
			"host_name":    *host.Name,
			"status":       *host.Health.Status,
			"ip_address":   *host.Assignment.IpAddress,
			"cluster_name": *host.Assignment.ClusterName,
			"zone":         *host.Assignment.Zone,
			"host_labels":  *&host.Labels,
		}
		hosts[i] = l
	}

	return hosts
}

func FlattenSatelliteCapabilities(capabilities *schema.Set) []kubernetesserviceapiv1.CapabilityManagedBySatellite {
	result := make([]kubernetesserviceapiv1.CapabilityManagedBySatellite, capabilities.Len())
	for i, v := range capabilities.List() {
		result[i] = kubernetesserviceapiv1.CapabilityManagedBySatellite(v.(string))
	}

	return result
}

func FlattenWorkerPoolHostLabels(hostLabels map[string]string) *schema.Set {
	mapped := make([]string, 0)
	for k, v := range hostLabels {
		if strings.HasPrefix(k, "os") {
			continue
		}
		mapped = append(mapped, fmt.Sprintf("%s:%v", k, v))
	}

	return NewStringSet(schema.HashString, mapped)
}

// KMS Private Endpoint
func updatePrivateURL(kpURL string) (string, error) {
	var kmsEndpointURL string
	if !strings.Contains(kpURL, "private") {
		kmsEndpURL := strings.SplitAfter(kpURL, "https://")
		if len(kmsEndpURL) == 2 {
			kmsEndpointURL = kmsEndpURL[0] + "private." + kmsEndpURL[1] + "/api/v2/"

		} else {
			return "", fmt.Errorf("[ERROR] Error in Kms EndPoint URL ")
		}
	}
	return kmsEndpointURL, nil
}

func FlattenSatelliteClusterZones(list []string) []map[string]interface{} {
	zones := make([]map[string]interface{}, len(list))
	for i, zone := range list {
		l := map[string]interface{}{
			"id": zone,
		}
		zones[i] = l
	}
	return zones
}

func FetchResourceInstanceDetails(d *schema.ResourceData, meta interface{}, instanceID string) error {
	// Get ResourceController from ClientSession
	resourceControllerClient, err := meta.(conns.ClientSession).ResourceControllerV2API()
	if err != nil {
		return err
	}

	getResourceOpts := resourcecontrollerv2.GetResourceInstanceOptions{
		ID: &instanceID,
	}

	instance, response, err := resourceControllerClient.GetResourceInstance(&getResourceOpts)
	if err != nil {
		log.Printf("[DEBUG] Error retrieving resource instance: %s\n%s", err, response)
		return fmt.Errorf("Error retrieving resource instance: %s\n%s", err, response)
	}
	if strings.Contains(*instance.State, "removed") {
		log.Printf("[DEBUG] Error retrieving resource instance details: Resource has been removed")
		return fmt.Errorf("Error retrieving resource instance details: Resource has been removed")
	}

	extensionsMap := Flatten(instance.Extensions)
	if extensionsMap == nil {
		log.Printf("[DEBUG] Error parsing resource instance: Endpoints are missing in instance Extensions map")
		return fmt.Errorf("Error parsing resource instance: Endpoints are missing in instance Extensions map")
	}
	d.Set("extensions", extensionsMap)

	return nil
}

func GetResourceInstanceURL(d *schema.ResourceData, meta interface{}) (*string, error) {

	var endpoint string
	extensions := d.Get("extensions").(map[string]interface{})

	if url, ok := extensions["endpoints.public"]; ok {
		endpoint = "https://" + url.(string)
	}

	if endpoint == "" {
		return nil, fmt.Errorf("[ERROR] Missing endpoints.public in extensions")
	}

	return &endpoint, nil
}

// Converts a struct to a map while maintaining the json alias as keys
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}

// This function takes two lists and returns the difference between the two lists
// Listdifference([1,2] [2,3]) = [1]
func Listdifference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// Stringify returns the stringified form of value "v".
// If "v" is a string-based type (string, strfmt.Date, strfmt.DateTime, strfmt.UUID, etc.),
// then it is returned unchanged (e.g. `this is a string`, `foo`, `2025-06-03`).
// Otherwise, json.Marshal() is used to serialze "v" and the resulting string is returned
// (e.g. `32`, `true`, `[true, false, true]`, `{"foo": "bar"}`).
// Note: the backticks in the comments above are not part of the returned strings.
func Stringify(v interface{}) string {
	if !core.IsNil(v) {
		if s, ok := v.(string); ok {
			return s
		} else if s, ok := v.(interface{ String() string }); ok {
			return s.String()
		} else {
			bytes, err := json.Marshal(v)
			if err != nil {
				log.Printf("[ERROR] Error marshaling 'any type' value as string: %s", err.Error())
				return ""
			}
			return string(bytes)
		}
	}
	return ""
}
