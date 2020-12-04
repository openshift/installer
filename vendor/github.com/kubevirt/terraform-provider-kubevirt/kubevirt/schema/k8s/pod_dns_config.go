package k8s

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	v1 "k8s.io/api/core/v1"
)

func podDnsConfigFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"nameservers": {
			Type:        schema.TypeList,
			Description: "A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed.",
			Optional:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.IsIPAddress,
			},
		},
		"option": {
			Type:        schema.TypeList,
			Description: "A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Description: "Name of the option.",
						Required:    true,
					},
					"value": {
						Type:        schema.TypeString,
						Description: "Value of the option. Optional: Defaults to empty.",
						Optional:    true,
					},
				},
			},
		},
		"searches": {
			Type:        schema.TypeList,
			Description: "A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed.",
			Optional:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: utils.ValidateName,
			},
		},
	}
}

func PodDnsConfigSchema() *schema.Schema {
	fields := podDnsConfigFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Specifies the DNS parameters of a pod. Parameters specified here will be merged to the generated DNS configuration based on DNSPolicy. Optional: Defaults to empty",
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func ExpandPodDNSConfig(l []interface{}) (*v1.PodDNSConfig, error) {
	if len(l) == 0 || l[0] == nil {
		return &v1.PodDNSConfig{}, nil
	}
	in := l[0].(map[string]interface{})
	obj := &v1.PodDNSConfig{}
	if v, ok := in["nameservers"].([]interface{}); ok {
		obj.Nameservers = utils.ExpandStringSlice(v)
	}
	if v, ok := in["searches"].([]interface{}); ok {
		obj.Searches = utils.ExpandStringSlice(v)
	}
	if v, ok := in["option"].([]interface{}); ok {
		opts, err := expandDNSConfigOptions(v)
		if err != nil {
			return obj, err
		}
		obj.Options = opts
	}
	return obj, nil
}

func expandDNSConfigOptions(options []interface{}) ([]v1.PodDNSConfigOption, error) {
	if len(options) == 0 {
		return []v1.PodDNSConfigOption{}, nil
	}
	opts := make([]v1.PodDNSConfigOption, len(options))
	for i, c := range options {
		in := c.(map[string]interface{})
		opt := v1.PodDNSConfigOption{}
		if v, ok := in["name"].(string); ok {
			opt.Name = v
		}
		if v, ok := in["value"].(string); ok {
			opt.Value = utils.PtrToString(v)
		}
		opts[i] = opt
	}

	return opts, nil
}

func FlattenPodDNSConfig(in *v1.PodDNSConfig) []interface{} {
	att := make(map[string]interface{})

	if len(in.Nameservers) > 0 {
		att["nameservers"] = in.Nameservers
	}
	if len(in.Searches) > 0 {
		att["searches"] = in.Searches
	}
	if len(in.Options) > 0 {
		att["option"] = flattenPodDNSConfigOptions(in.Options)
	}

	if len(att) > 0 {
		return []interface{}{att}
	}
	return []interface{}{}
}

func flattenPodDNSConfigOptions(options []v1.PodDNSConfigOption) []interface{} {
	att := make([]interface{}, len(options))
	for i, v := range options {
		obj := map[string]interface{}{}

		if v.Name != "" {
			obj["name"] = v.Name
		}
		if v.Value != nil {
			obj["value"] = *v.Value
		}
		att[i] = obj
	}
	return att
}
