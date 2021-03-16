package virtualmachineinstance

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/schema/k8s"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	k8sv1 "k8s.io/api/core/v1"
	kubevirtapiv1 "kubevirt.io/client-go/api/v1"
)

func virtualMachineInstanceSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"priority_class_name": {
			Type:        schema.TypeString,
			Description: "If specified, indicates the pod's priority. If not specified, the pod priority will be default or zero if there is no default.",
			Optional:    true,
		},
		"domain": domainSpecSchema(),
		"node_selector": {
			Type:        schema.TypeMap,
			Description: "NodeSelector is a selector which must be true for the vmi to fit on a node. Selector which must match a node's labels for the vmi to be scheduled on that node.",
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"affinity": k8s.AffinitySchema(),
		"scheduler_name": {
			Type:        schema.TypeString,
			Description: "If specified, the VMI will be dispatched by specified scheduler. If not specified, the VMI will be dispatched by default scheduler.",
			Optional:    true,
		},
		"tolerations": k8s.TolerationSchema(),
		"eviction_strategy": {
			Type:        schema.TypeString,
			Description: "EvictionStrategy can be set to \"LiveMigrate\" if the VirtualMachineInstance should be migrated instead of shut-off in case of a node drain.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"LiveMigrate",
			}, false),
		},
		"termination_grace_period_seconds": {
			Type:        schema.TypeInt,
			Description: "Grace period observed after signalling a VirtualMachineInstance to stop after which the VirtualMachineInstance is force terminated.",
			Optional:    true,
		},
		"volume":          volumesSchema(),
		"liveness_probe":  probeSchema(),
		"readiness_probe": probeSchema(),
		"hostname": {
			Type:        schema.TypeString,
			Description: "Specifies the hostname of the vmi.",
			Optional:    true,
		},
		"subdomain": {
			Type:        schema.TypeString,
			Description: "If specified, the fully qualified vmi hostname will be \"<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>\".",
			Optional:    true,
		},
		"network": networksSchema(),
		"dns_policy": {
			Type:        schema.TypeString,
			Description: "DNSPolicy defines how a pod's DNS will be configured.",
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"ClusterFirstWithHostNet",
				"ClusterFirst",
				"Default",
				"None",
			}, false),
		},
		"pod_dns_config": k8s.PodDnsConfigSchema(),
	}
}

func virtualMachineInstanceSpecSchema() *schema.Schema {
	fields := virtualMachineInstanceSpecFields()

	return &schema.Schema{
		Type: schema.TypeList,

		Description: fmt.Sprintf("Template is the direct specification of VirtualMachineInstance."),
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func expandVirtualMachineInstanceSpec(virtualMachineInstanceSpec []interface{}) (kubevirtapiv1.VirtualMachineInstanceSpec, error) {
	result := kubevirtapiv1.VirtualMachineInstanceSpec{}

	if len(virtualMachineInstanceSpec) == 0 || virtualMachineInstanceSpec[0] == nil {
		return result, nil
	}

	in := virtualMachineInstanceSpec[0].(map[string]interface{})

	if v, ok := in["priority_class_name"].(string); ok {
		result.PriorityClassName = v
	}
	if v, ok := in["domain"].([]interface{}); ok {
		domain, err := expandDomainSpec(v)
		if err != nil {
			return result, err
		}
		result.Domain = domain
	}
	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		result.NodeSelector = utils.ExpandStringMap(v)
	}
	if v, ok := in["affinity"].([]interface{}); ok {
		result.Affinity = k8s.ExpandAffinity(v)
	}
	if v, ok := in["scheduler_name"].(string); ok {
		result.SchedulerName = v
	}
	if v, ok := in["tolerations"].([]interface{}); ok {
		tolerations, err := k8s.ExpandTolerations(v)
		if err != nil {
			return result, err
		}
		result.Tolerations = tolerations
	}
	if v, ok := in["eviction_strategy"].(string); ok {
		if v != "" {
			evictionStrategy := kubevirtapiv1.EvictionStrategy(v)
			result.EvictionStrategy = &evictionStrategy
		}
	}
	if v, ok := in["termination_grace_period_seconds"].(int); ok {
		seconds := int64(v)
		result.TerminationGracePeriodSeconds = &seconds
	}
	if v, ok := in["volume"].([]interface{}); ok {
		result.Volumes = expandVolumes(v)
	}
	if v, ok := in["liveness_probe"].([]interface{}); ok {
		result.LivenessProbe = expandProbe(v)
	}
	if v, ok := in["readiness_probe"].([]interface{}); ok {
		result.ReadinessProbe = expandProbe(v)
	}
	if v, ok := in["hostname"].(string); ok {
		result.Hostname = v
	}
	if v, ok := in["subdomain"].(string); ok {
		result.Subdomain = v
	}
	if v, ok := in["network"].([]interface{}); ok {
		result.Networks = expandNetworks(v)
	}
	if v, ok := in["dns_policy"].(string); ok {
		result.DNSPolicy = k8sv1.DNSPolicy(v)
	}
	if v, ok := in["pod_dns_config"].([]interface{}); ok {
		dnsConfig, err := k8s.ExpandPodDNSConfig(v)
		if err != nil {
			return result, err
		}
		result.DNSConfig = dnsConfig
	}

	return result, nil
}

func flattenVirtualMachineInstanceSpec(in kubevirtapiv1.VirtualMachineInstanceSpec) []interface{} {
	att := make(map[string]interface{})

	att["priority_class_name"] = in.PriorityClassName
	att["domain"] = flattenDomainSpec(in.Domain)
	att["node_selector"] = utils.FlattenStringMap(in.NodeSelector)
	att["affinity"] = k8s.FlattenAffinity(in.Affinity)
	att["scheduler_name"] = in.SchedulerName
	att["tolerations"] = k8s.FlattenTolerations(in.Tolerations)
	if in.EvictionStrategy != nil {
		att["eviction_strategy"] = string(*in.EvictionStrategy)
	}
	if in.TerminationGracePeriodSeconds != nil {
		att["termination_grace_period_seconds"] = *in.TerminationGracePeriodSeconds
	}
	att["volume"] = flattenVolumes(in.Volumes)
	if in.LivenessProbe != nil {
		att["liveness_probe"] = flattenProbe(*in.LivenessProbe)
	}
	if in.ReadinessProbe != nil {
		att["readiness_probe"] = flattenProbe(*in.ReadinessProbe)
	}
	att["hostname"] = in.Hostname
	att["subdomain"] = in.Subdomain
	att["network"] = flattenNetworks(in.Networks)
	att["dns_policy"] = string(in.DNSPolicy)
	if in.DNSConfig != nil {
		att["pod_dns_config"] = k8s.FlattenPodDNSConfig(in.DNSConfig)
	}

	return []interface{}{att}
}
