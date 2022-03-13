package vmworkflow

import (
	"errors"
	"fmt"
	"net"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	cKeyPrefix        = "clone.0.customize.0"
	cLinuxKeyPrefix   = "clone.0.customize.0.linux_options.0"
	cWindowsKeyPrefix = "clone.0.customize.0.windows_options.0"
	cNetifKeyPrefix   = "clone.0.customize.0.network_interface"
)

// netifKey renders a specific network_interface key for a specific resource
// index.
func netifKey(key string, n int) string {
	return fmt.Sprintf("%s.%d.%s", cNetifKeyPrefix, n, key)
}

// matchGateway take an IP, mask, and gateway, and checks to see if the gateway
// is reachable from the IP address.
func matchGateway(a string, m int, g string) bool {
	ip := net.ParseIP(a)
	gw := net.ParseIP(g)
	var mask net.IPMask
	if ip.To4() != nil {
		mask = net.CIDRMask(m, 32)
	} else {
		mask = net.CIDRMask(m, 128)
	}
	if ip.Mask(mask).Equal(gw.Mask(mask)) {
		return true
	}
	return false
}

func v4CIDRMaskToDotted(mask int) string {
	m := net.CIDRMask(mask, 32)
	a := int(m[0])
	b := int(m[1])
	c := int(m[2])
	d := int(m[3])
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

// VirtualMachineCustomizeSchema returns the schema for VM customization.
func VirtualMachineCustomizeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// CustomizationGlobalIPSettings
		"dns_server_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The list of DNS servers for a virtual network adapter with a static IP address.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"dns_suffix_list": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A list of DNS search domains to add to the DNS configuration on the virtual machine.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},

		// CustomizationLinuxPrep
		"linux_options": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{cKeyPrefix + "." + "windows_options", cKeyPrefix + "." + "windows_sysprep_text"},
			Description:   "A list of configuration options specific to Linux virtual machines.",
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"domain": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The FQDN for this virtual machine.",
				},
				"host_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The host name for this virtual machine.",
				},
				"hw_clock_utc": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "Specifies whether or not the hardware clock should be in UTC or not.",
				},
				"time_zone": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Customize the time zone on the VM. This should be a time zone-style entry, like America/Los_Angeles.",
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile("^[-+/_a-zA-Z0-9]+$"),
						"must be similar to America/Los_Angeles or other Linux/Unix TZ format",
					),
				},
			}},
		},

		// CustomizationSysprep
		"windows_options": {
			Type:          schema.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{cKeyPrefix + "." + "linux_options", cKeyPrefix + "." + "windows_sysprep_text"},
			Description:   "A list of configuration options specific to Windows virtual machines.",
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				// CustomizationGuiRunOnce
				"run_once_command_list": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "A list of commands to run at first user logon, after guest customization.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				// CustomizationGuiUnattended
				"auto_logon": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Specifies whether or not the VM automatically logs on as Administrator.",
				},
				"auto_logon_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     1,
					Description: "Specifies how many times the VM should auto-logon the Administrator account when auto_logon is true.",
				},
				"admin_password": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "The new administrator password for this virtual machine.",
				},
				"time_zone": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     85,
					Description: "The new time zone for the virtual machine. This is a sysprep-dictated timezone code.",
				},

				// CustomizationIdentification
				"domain_admin_user": {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{cWindowsKeyPrefix + "." + "workgroup"},
					Description:   "The user account of the domain administrator used to join this virtual machine to the domain.",
				},
				"domain_admin_password": {
					Type:          schema.TypeString,
					Optional:      true,
					Sensitive:     true,
					ConflictsWith: []string{cWindowsKeyPrefix + "." + "workgroup"},
					Description:   "The password of the domain administrator used to join this virtual machine to the domain.",
				},
				"join_domain": {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{cWindowsKeyPrefix + "." + "workgroup"},
					Description:   "The domain that the virtual machine should join.",
				},
				"workgroup": {
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{cWindowsKeyPrefix + "." + "join_domain"},
					Description:   "The workgroup for this virtual machine if not joining a domain.",
				},

				// CustomizationUserData
				"computer_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The host name for this virtual machine.",
				},
				"full_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "Administrator",
					Description: "The full name of the user of this virtual machine.",
				},
				"organization_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "Managed by Terraform",
					Description: "The organization name this virtual machine is being installed for.",
				},
				"product_key": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "The product key for this virtual machine.",
				},
			}},
		},

		// CustomizationSysprepText
		"windows_sysprep_text": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{cKeyPrefix + "." + "linux_options", cKeyPrefix + "." + "windows_options"},
			Description:   "Use this option to specify a windows sysprep file directly.",
		},

		// CustomizationIPSettings
		"network_interface": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "A specification of network interface configuration options.",
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				"dns_server_list": {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "Network-interface specific DNS settings for Windows operating systems. Ignored on Linux.",
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"dns_domain": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "A DNS search domain to add to the DNS configuration on the virtual machine.",
				},
				"ipv4_address": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The IPv4 address assigned to this network adapter. If left blank, DHCP is used.",
				},
				"ipv4_netmask": {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  "The IPv4 CIDR netmask for the supplied IP address. Ignored if DHCP is selected.",
					ValidateFunc: validation.IntAtMost(32),
				},
				"ipv6_address": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The IPv6 address assigned to this network adapter. If left blank, default auto-configuration is used.",
				},
				"ipv6_netmask": {
					Type:         schema.TypeInt,
					Optional:     true,
					Description:  "The IPv6 CIDR netmask for the supplied IP address. Ignored if auto-configuration is selected.",
					ValidateFunc: validation.IntAtMost(128),
				},
			}},
		},

		// Base-level settings
		"ipv4_gateway": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The IPv4 default gateway when using network_interface customization on the virtual machine. This address must be local to a static IPv4 address configured in an interface sub-resource.",
		},
		"ipv6_gateway": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The IPv6 default gateway when using network_interface customization on the virtual machine. This address must be local to a static IPv4 address configured in an interface sub-resource.",
		},
		"timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     10,
			Description: "The amount of time, in minutes, to wait for guest OS customization to complete before returning with an error. Setting this value to 0 or a negative value skips the waiter.",
		},
	}
}

// expandCustomizationGlobalIPSettings reads certain ResourceData keys and
// returns a CustomizationGlobalIPSettings.
func expandCustomizationGlobalIPSettings(d *schema.ResourceData) types.CustomizationGlobalIPSettings {
	obj := types.CustomizationGlobalIPSettings{
		DnsSuffixList: structure.SliceInterfacesToStrings(d.Get(cKeyPrefix + "." + "dns_suffix_list").([]interface{})),
		DnsServerList: structure.SliceInterfacesToStrings(d.Get(cKeyPrefix + "." + "dns_server_list").([]interface{})),
	}
	return obj
}

// expandCustomizationLinuxPrep reads certain ResourceData keys and
// returns a CustomizationLinuxPrep.
func expandCustomizationLinuxPrep(d *schema.ResourceData) *types.CustomizationLinuxPrep {
	obj := &types.CustomizationLinuxPrep{
		HostName: &types.CustomizationFixedName{
			Name: d.Get(cLinuxKeyPrefix + "." + "host_name").(string),
		},
		Domain:     d.Get(cLinuxKeyPrefix + "." + "domain").(string),
		TimeZone:   d.Get(cLinuxKeyPrefix + "." + "time_zone").(string),
		HwClockUTC: structure.GetBoolPtr(d, cLinuxKeyPrefix+"."+"hw_clock_utc"),
	}
	return obj
}

// expandCustomizationGuiRunOnce reads certain ResourceData keys and
// returns a CustomizationGuiRunOnce.
func expandCustomizationGuiRunOnce(d *schema.ResourceData) *types.CustomizationGuiRunOnce {
	obj := &types.CustomizationGuiRunOnce{
		CommandList: structure.SliceInterfacesToStrings(d.Get(cWindowsKeyPrefix + "." + "run_once_command_list").([]interface{})),
	}
	if len(obj.CommandList) < 1 {
		return nil
	}
	return obj
}

// expandCustomizationGuiUnattended reads certain ResourceData keys and
// returns a CustomizationGuiUnattended.
func expandCustomizationGuiUnattended(d *schema.ResourceData) types.CustomizationGuiUnattended {
	obj := types.CustomizationGuiUnattended{
		TimeZone:       int32(d.Get(cWindowsKeyPrefix + "." + "time_zone").(int)),
		AutoLogon:      d.Get(cWindowsKeyPrefix + "." + "auto_logon").(bool),
		AutoLogonCount: int32(d.Get(cWindowsKeyPrefix + "." + "auto_logon_count").(int)),
	}
	if v, ok := d.GetOk(cWindowsKeyPrefix + "." + "admin_password"); ok {
		obj.Password = &types.CustomizationPassword{
			Value:     v.(string),
			PlainText: true,
		}
	}

	return obj
}

// expandCustomizationIdentification reads certain ResourceData keys and
// returns a CustomizationIdentification.
func expandCustomizationIdentification(d *schema.ResourceData) types.CustomizationIdentification {
	obj := types.CustomizationIdentification{
		JoinWorkgroup: d.Get(cWindowsKeyPrefix + "." + "workgroup").(string),
		JoinDomain:    d.Get(cWindowsKeyPrefix + "." + "join_domain").(string),
		DomainAdmin:   d.Get(cWindowsKeyPrefix + "." + "domain_admin_user").(string),
	}
	if v, ok := d.GetOk(cWindowsKeyPrefix + "." + "domain_admin_password"); ok {
		obj.DomainAdminPassword = &types.CustomizationPassword{
			Value:     v.(string),
			PlainText: true,
		}
	}
	return obj
}

// expandCustomizationUserData reads certain ResourceData keys and
// returns a CustomizationUserData.
func expandCustomizationUserData(d *schema.ResourceData) types.CustomizationUserData {
	obj := types.CustomizationUserData{
		FullName: d.Get(cWindowsKeyPrefix + "." + "full_name").(string),
		OrgName:  d.Get(cWindowsKeyPrefix + "." + "organization_name").(string),
		ComputerName: &types.CustomizationFixedName{
			Name: d.Get(cWindowsKeyPrefix + "." + "computer_name").(string),
		},
		ProductId: d.Get(cWindowsKeyPrefix + "." + "product_key").(string),
	}
	return obj
}

// expandCustomizationSysprep reads certain ResourceData keys and
// returns a CustomizationSysprep.
func expandCustomizationSysprep(d *schema.ResourceData) *types.CustomizationSysprep {
	obj := &types.CustomizationSysprep{
		GuiUnattended:  expandCustomizationGuiUnattended(d),
		UserData:       expandCustomizationUserData(d),
		GuiRunOnce:     expandCustomizationGuiRunOnce(d),
		Identification: expandCustomizationIdentification(d),
	}
	return obj
}

// expandCustomizationSysprepText reads certain ResourceData keys and
// returns a CustomizationSysprepText.
func expandCustomizationSysprepText(d *schema.ResourceData) *types.CustomizationSysprepText {
	obj := &types.CustomizationSysprepText{
		Value: d.Get(cKeyPrefix + "." + "windows_sysprep_text").(string),
	}
	return obj
}

// expandBaseCustomizationIdentitySettings returns a
// BaseCustomizationIdentitySettings, depending on what is defined.
//
// Only one of the three types of identity settings can be specified: Linux
// settings (from linux_options), Windows settings (from windows_options), and
// the raw Windows sysprep file (via windows_sysprep_text).
func expandBaseCustomizationIdentitySettings(d *schema.ResourceData, family string) types.BaseCustomizationIdentitySettings {
	var obj types.BaseCustomizationIdentitySettings
	_, windowsExists := d.GetOkExists(cKeyPrefix + "." + "windows_options")
	_, sysprepExists := d.GetOkExists(cKeyPrefix + "." + "windows_sysprep_text")
	switch {
	case family == string(types.VirtualMachineGuestOsFamilyLinuxGuest):
		obj = expandCustomizationLinuxPrep(d)
	case family == string(types.VirtualMachineGuestOsFamilyWindowsGuest) && windowsExists:
		obj = expandCustomizationSysprep(d)
	case family == string(types.VirtualMachineGuestOsFamilyWindowsGuest) && sysprepExists:
		obj = expandCustomizationSysprepText(d)
	default:
		obj = &types.CustomizationIdentitySettings{}
	}
	return obj
}

// expandCustomizationIPSettingsIPV6AddressSpec reads certain ResourceData keys and
// returns a CustomizationIPSettingsIpV6AddressSpec.
func expandCustomizationIPSettingsIPV6AddressSpec(d *schema.ResourceData, n int, gwAdd bool) (*types.CustomizationIPSettingsIpV6AddressSpec, bool) {
	v, ok := d.GetOk(netifKey("ipv6_address", n))
	var gwFound bool
	if !ok {
		return nil, gwFound
	}
	addr := v.(string)
	mask := d.Get(netifKey("ipv6_netmask", n)).(int)
	gw, gwOk := d.Get(cKeyPrefix + "." + "ipv6_gateway").(string)
	obj := &types.CustomizationIPSettingsIpV6AddressSpec{
		Ip: []types.BaseCustomizationIpV6Generator{
			&types.CustomizationFixedIpV6{
				IpAddress:  addr,
				SubnetMask: int32(mask),
			},
		},
	}
	if gwAdd && gwOk && matchGateway(addr, mask, gw) {
		obj.Gateway = []string{gw}
		gwFound = true
	}
	return obj, gwFound
}

// expandCustomizationIPSettings reads certain ResourceData keys and
// returns a CustomizationIPSettings.
func expandCustomizationIPSettings(d *schema.ResourceData, n int, v4gwAdd, v6gwAdd bool) (types.CustomizationIPSettings, bool, bool) {
	var v4gwFound, v6gwFound bool
	v4addr, v4addrOk := d.GetOk(netifKey("ipv4_address", n))
	v4mask := d.Get(netifKey("ipv4_netmask", n)).(int)
	v4gw, v4gwOk := d.Get(cKeyPrefix + "." + "ipv4_gateway").(string)
	var obj types.CustomizationIPSettings
	switch {
	case v4addrOk:
		obj.Ip = &types.CustomizationFixedIp{
			IpAddress: v4addr.(string),
		}
		obj.SubnetMask = v4CIDRMaskToDotted(v4mask)
		// Check for the gateway
		if v4gwAdd && v4gwOk && matchGateway(v4addr.(string), v4mask, v4gw) {
			obj.Gateway = []string{v4gw}
			v4gwFound = true
		}
	default:
		obj.Ip = &types.CustomizationDhcpIpGenerator{}
	}
	obj.DnsServerList = structure.SliceInterfacesToStrings(d.Get(netifKey("dns_server_list", n)).([]interface{}))
	obj.DnsDomain = d.Get(netifKey("dns_domain", n)).(string)
	obj.IpV6Spec, v6gwFound = expandCustomizationIPSettingsIPV6AddressSpec(d, n, v6gwAdd)
	return obj, v4gwFound, v6gwFound
}

// expandSliceOfCustomizationAdapterMapping reads certain ResourceData keys and
// returns a CustomizationAdapterMapping slice.
func expandSliceOfCustomizationAdapterMapping(d *schema.ResourceData) []types.CustomizationAdapterMapping {
	s := d.Get(cKeyPrefix + "." + "network_interface").([]interface{})
	if len(s) < 1 {
		return nil
	}
	result := make([]types.CustomizationAdapterMapping, len(s))
	var v4gwFound, v6gwFound bool
	for i := range s {
		var adapter types.CustomizationIPSettings
		adapter, v4gwFound, v6gwFound = expandCustomizationIPSettings(d, i, !v4gwFound, !v6gwFound)
		obj := types.CustomizationAdapterMapping{
			Adapter: adapter,
		}
		result[i] = obj
	}
	return result
}

// ExpandCustomizationSpec reads certain ResourceData keys and
// returns a CustomizationSpec.
func ExpandCustomizationSpec(d *schema.ResourceData, family string) types.CustomizationSpec {
	obj := types.CustomizationSpec{
		Identity:         expandBaseCustomizationIdentitySettings(d, family),
		GlobalIPSettings: expandCustomizationGlobalIPSettings(d),
		NicSettingMap:    expandSliceOfCustomizationAdapterMapping(d),
	}
	return obj
}

// ValidateCustomizationSpec checks the validity of the supplied customization
// spec. It should be called during diff customization to veto invalid configs.
func ValidateCustomizationSpec(d *schema.ResourceDiff, family string) error {
	// Validate that the proper section exists for OS family suboptions.
	linuxExists := len(d.Get(cKeyPrefix+"."+"linux_options").([]interface{})) > 0 || !structure.ValuesAvailable(cKeyPrefix+"."+"linux_options.", []string{"host_name", "domain"}, d)
	windowsExists := len(d.Get(cKeyPrefix+"."+"windows_options").([]interface{})) > 0 || !structure.ValuesAvailable(cKeyPrefix+"."+"windows_options.", []string{"computer_name"}, d)
	sysprepExists := d.Get(cKeyPrefix+"."+"windows_sysprep_text").(string) != "" || !structure.ValuesAvailable(cKeyPrefix+".", []string{"windows_sysprep_text"}, d)
	switch {
	case family == string(types.VirtualMachineGuestOsFamilyLinuxGuest) && !linuxExists:
		return errors.New("linux_options must exist in VM customization options for Linux operating systems")
	case family == string(types.VirtualMachineGuestOsFamilyWindowsGuest) && !windowsExists && !sysprepExists:
		return errors.New("one of windows_options or windows_sysprep_text must exist in VM customization options for Windows operating systems")
	}
	return nil
}
