// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
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

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/bluemix-go/models"
	"github.com/IBM-Cloud/container-services-go-sdk/kubernetesserviceapiv1"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go-config/resourceconfigurationv1"
	"github.com/IBM/ibm-cos-sdk-go/service/s3"
	kp "github.com/IBM/keyprotect-go-client"
	"github.com/IBM/platform-services-go-sdk/globaltaggingv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	rg "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"

	"github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/IBM-Cloud/bluemix-go/api/icd/icdv4"
	"github.com/IBM-Cloud/bluemix-go/api/mccp/mccpv2"
	"github.com/IBM-Cloud/bluemix-go/api/schematics"
	"github.com/IBM-Cloud/bluemix-go/api/usermanagement/usermanagementv2"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
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
)

//HashInt ...
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

func ExpandProtocols(configured []interface{}) ([]datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration, error) {
	protocols := make([]datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration, 0, len(configured))
	var lbMethodToId = make(map[string]string)
	for _, lRaw := range configured {
		data := lRaw.(map[string]interface{})
		p := &datatypes.Network_LBaaS_LoadBalancerProtocolConfiguration{
			FrontendProtocol: sl.String(data["frontend_protocol"].(string)),
			BackendProtocol:  sl.String(data["backend_protocol"].(string)),
			FrontendPort:     sl.Int(data["frontend_port"].(int)),
			BackendPort:      sl.Int(data["backend_port"].(int)),
		}
		if v, ok := data["session_stickiness"]; ok && v.(string) != "" {
			p.SessionType = sl.String(v.(string))
		}
		if v, ok := data["max_conn"]; ok && v.(int) != 0 {
			p.MaxConn = sl.Int(v.(int))
		}
		if v, ok := data["tls_certificate_id"]; ok && v.(int) != 0 {
			p.TlsCertificateId = sl.Int(v.(int))
		}
		if v, ok := data["load_balancing_method"]; ok {
			p.LoadBalancingMethod = sl.String(lbMethodToId[v.(string)])
		}
		if v, ok := data["protocol_id"]; ok && v.(string) != "" {
			p.ListenerUuid = sl.String(v.(string))
		}

		var isValid bool
		if p.TlsCertificateId != nil && *p.TlsCertificateId != 0 {
			// validate the protocol is correct
			if *p.FrontendProtocol == "HTTPS" {
				isValid = true
			}
		} else {
			isValid = true
		}

		if isValid {
			protocols = append(protocols, *p)
		} else {
			return protocols, fmt.Errorf("tls_certificate_id may be set only when frontend protocol is 'HTTPS'")
		}

	}
	return protocols, nil
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

func FlattenVpcWorkerPools(list []containerv2.GetWorkerPoolResponse) []map[string]interface{} {
	workerPools := make([]map[string]interface{}, len(list))
	for i, workerPool := range list {
		l := map[string]interface{}{
			"id":           workerPool.ID,
			"name":         workerPool.PoolName,
			"flavor":       workerPool.Flavor,
			"worker_count": workerPool.WorkerCount,
			"isolation":    workerPool.Isolation,
			"labels":       workerPool.Labels,
			"state":        workerPool.Lifecycle.ActualState,
			"host_pool_id": workerPool.HostPoolID,
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

func FlattenIcdGroups(grouplist icdv4.GroupList) []map[string]interface{} {
	groups := make([]map[string]interface{}, len(grouplist.Groups))
	for i, group := range grouplist.Groups {
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
		cpu["units"] = group.Cpu.Units
		cpu["allocation_count"] = group.Cpu.AllocationCount
		cpu["minimum_count"] = group.Cpu.MinimumCount
		cpu["step_size_count"] = group.Cpu.StepSizeCount
		cpu["is_adjustable"] = group.Cpu.IsAdjustable
		cpu["can_scale_down"] = group.Cpu.CanScaleDown
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

		l := map[string]interface{}{
			"group_id": group.Id,
			"count":    group.Count,
			"memory":   memorys,
			"cpu":      cpus,
			"disk":     disks,
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

func IntValue(i64 *int64) (i int) {
	if i64 != nil {
		i = int(*i64)
	}
	return
}

func float64Value(f32 *float32) (f float64) {
	if f32 != nil {
		f = float64(*f32)
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
func ExpandWhitelist(whiteList *schema.Set) (whitelist []icdv4.WhitelistEntry) {
	for _, iface := range whiteList.List() {
		wlItem := iface.(map[string]interface{})
		wlEntry := icdv4.WhitelistEntry{
			Address:     wlItem["address"].(string),
			Description: wlItem["description"].(string),
		}
		whitelist = append(whitelist, wlEntry)
	}
	return
}

// Cloud Internet Services
func FlattenWhitelist(whitelist icdv4.Whitelist) []map[string]interface{} {
	entries := make([]map[string]interface{}, len(whitelist.WhitelistEntrys), len(whitelist.WhitelistEntrys))
	for i, whitelistEntry := range whitelist.WhitelistEntrys {
		l := map[string]interface{}{
			"address":     whitelistEntry.Address,
			"description": whitelistEntry.Description,
		}
		entries[i] = l
	}
	return entries
}

func ExpandPlatformOptions(platformOptions icdv4.PlatformOptions) []map[string]interface{} {
	pltOptions := make([]map[string]interface{}, 0, 1)
	pltOption := make(map[string]interface{})
	pltOption["key_protect_key_id"] = platformOptions.KeyProtectKey
	pltOption["disk_encryption_key_crn"] = platformOptions.DiskENcryptionKeyCrn
	pltOption["backup_encryption_key_crn"] = platformOptions.BackUpEncryptionKeyCrn
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

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPIv1()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return nil, err
	}
	accountID := userDetails.UserAccount

	var providers []string
	if strings.Contains(resourceType, "SoftLayer_") {
		providers = []string{"ims"}
	}

	ListTagsOptions := &globaltaggingv1.ListTagsOptions{}
	if resourceID != "" {
		ListTagsOptions.AttachedTo = &resourceID
	}
	ListTagsOptions.Providers = providers
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
	log.Println("tagList: ", taglist)
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

		_, resp, err := gtClient.DetachTag(detachTagOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error detaching database tags %v: %s\n%s", remove, err, resp)
		}
		for _, v := range remove {
			delTagOptions := &globaltaggingv1.DeleteTagOptions{
				TagName: PtrToString(v),
			}
			_, resp, err := gtClient.DeleteTag(delTagOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error deleting database tag %v: %s\n%s", v, err, resp)
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
	}

	return nil
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
func GetTagsUsingCRN(meta interface{}, resourceCRN string) (*schema.Set, error) {

	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPI()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error getting global tagging client settings: %s", err)
	}
	taggingResult, err := gtClient.Tags().GetTags(resourceCRN)
	if err != nil {
		return nil, err
	}
	var taglist []string
	for _, item := range taggingResult.Items {
		taglist = append(taglist, item.Name)
	}
	log.Println("tagList: ", taglist)
	return NewStringSet(ResourceIBMVPCHash, taglist), nil
}

func UpdateTagsUsingCRN(oldList, newList interface{}, meta interface{}, resourceCRN string) error {
	gtClient, err := meta.(conns.ClientSession).GlobalTaggingAPI()
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

	if len(remove) > 0 {
		_, err := gtClient.Tags().DetachTags(resourceCRN, remove)
		if err != nil {
			return fmt.Errorf("[ERROR] Error detaching database tags %v: %s", remove, err)
		}
		for _, v := range remove {
			_, err := gtClient.Tags().DeleteTag(v)
			if err != nil {
				return fmt.Errorf("[ERROR] Error deleting database tag %v: %s", v, err)
			}
		}
	}

	if len(add) > 0 {
		_, err := gtClient.Tags().AttachTags(resourceCRN, add)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating database tags %v : %s", add, err)
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

func ResourceLBListenerPolicyCustomizeDiff(diff *schema.ResourceDiff) error {
	policyActionIntf, _ := diff.GetOk(isLBListenerPolicyAction)
	policyAction := policyActionIntf.(string)

	if policyAction == "forward" {
		_, policyTargetIDSet := diff.GetOk(isLBListenerPolicyTargetID)

		if !policyTargetIDSet && diff.NewValueKnown(isLBListenerPolicyTargetID) {
			return fmt.Errorf("Load balancer listener policy: When action is forward please specify target_id")
		}
	} else if policyAction == "redirect" {
		_, httpsStatusCodeSet := diff.GetOk(isLBListenerPolicyTargetHTTPStatusCode)
		_, targetURLSet := diff.GetOk(isLBListenerPolicyTargetURL)

		if !httpsStatusCodeSet && diff.NewValueKnown(isLBListenerPolicyTargetHTTPStatusCode) {
			return fmt.Errorf("Load balancer listener policy: When action is redirect please specify target_http_status_code")
		}

		if !targetURLSet && diff.NewValueKnown(isLBListenerPolicyTargetURL) {
			return fmt.Errorf("Load balancer listener policy: When action is redirect please specify target_url")
		}
	} else if policyAction == "https_redirect" {
		_, listenerSet := diff.GetOk(isLBListenerPolicyHTTPSRedirectListener)
		_, httpsStatusSet := diff.GetOk(isLBListenerPolicyHTTPSRedirectStatusCode)

		if !listenerSet && diff.NewValueKnown(isLBListenerPolicyHTTPSRedirectListener) {
			return fmt.Errorf("Load balancer listener policy: When action is https_redirect please specify target_https_redirect_listener")
		}

		if !httpsStatusSet && diff.NewValueKnown(isLBListenerPolicyHTTPSRedirectStatusCode) {
			return fmt.Errorf("When action is https_redirect please specify target_https_redirect_status_code")
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

func InstanceProfileValidate(diff *schema.ResourceDiff) error {
	if diff.Id() != "" && diff.HasChange("profile") {
		o, n := diff.GetChange("profile")
		old := o.(string)
		new := n.(string)
		log.Println("old profile : ", old)
		log.Println("new profile : ", new)
		if !strings.Contains(old, "d") && strings.Contains(new, "d") {
			diff.ForceNew("profile")
		} else if strings.Contains(old, "d") && !strings.Contains(new, "d") {
			diff.ForceNew("profile")
		}
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

	if profile != "custom" {
		if iops != 0 && diff.NewValueKnown("iops") && diff.HasChange("iops") {
			return fmt.Errorf("VolumeError : iops is applicable for only custom volume profiles")
		}
	} else {
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

// flattenHostLabels ..
func FlattenHostLabels(hostLabels []interface{}) map[string]string {
	labels := make(map[string]string)
	for _, v := range hostLabels {
		parts := strings.Split(v.(string), ":")
		if parts != nil {
			labels[parts[0]] = parts[1]
		}
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

func GetSubjectAttribute(name string, s iampolicymanagementv1.PolicySubject) *string {
	for _, a := range s.Attributes {
		if *a.Name == name {
			return a.Value
		}
	}
	return core.StringPtr("")
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
func FindRoleByName(supported []iampolicymanagementv1.PolicyRole, name string) (iampolicymanagementv1.PolicyRole, error) {
	for _, role := range supported {
		if role.DisplayName != nil {
			if *role.DisplayName == name {
				role.DisplayName = nil
				return role, nil
			}
		}
	}
	supportedRoles := getSupportedRolesStr(supported)
	return iampolicymanagementv1.PolicyRole{}, bmxerror.New("RoleDoesnotExist",
		fmt.Sprintf("%s was not found. Valid roles are %s", name, supportedRoles))

}

func getSupportedRolesStr(supported []iampolicymanagementv1.PolicyRole) string {
	rolesStr := ""
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

func MapRoleListToPolicyRoles(roleList iampolicymanagementv1.RoleList) []iampolicymanagementv1.PolicyRole {
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

func GeneratePolicyOptions(d *schema.ResourceData, meta interface{}) (iampolicymanagementv1.CreatePolicyOptions, error) {

	var serviceName string
	var resourceType string
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

	serviceToQuery := serviceName

	if serviceName == "" && // no specific service specified
		!d.Get("account_management").(bool) && // not all account management services
		resourceType != "resource-group" { // not to a resource group
		serviceToQuery = "alliamserviceroles"
	}

	listRoleOptions := &iampolicymanagementv1.ListRolesOptions{
		AccountID:   &userDetails.UserAccount,
		ServiceName: &serviceToQuery,
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

func FlattenWorkerPoolHostLabels(hostLabels map[string]string) *schema.Set {
	mapped := make([]string, len(hostLabels))
	idx := 0
	for k, v := range hostLabels {
		mapped[idx] = fmt.Sprintf("%s:%v", k, v)
		idx++
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
