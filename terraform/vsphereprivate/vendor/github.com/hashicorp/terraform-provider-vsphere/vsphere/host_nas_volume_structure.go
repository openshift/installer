package vsphere

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	hostNasVolumeAccessModeReadOnly  = "readOnly"
	hostNasVolumeAccessModeReadWrite = "readWrite"

	hostNasVolumeSecurityTypeAuthSys  = "AUTH_SYS"
	hostNasVolumeSecurityTypeSecKrb5  = "SEC_KRB5"
	hostNasVolumeSecurityTypeSecKrb5i = "SEC_KRB5I"
)

// schemaHostNasVolumeSpec returns schema items for resources that need to work
// with a HostNasVolumeSpec.
func schemaHostNasVolumeSpec() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// HostNasVolumeSpec
		// Skipped attributes: localPath (this is the name attribute)
		// All CIFS attributes (we currently do not support CIFS as it's not
		// available in the vSphere client and there is not much data about how to
		// get it working)
		"access_mode": {
			Type:        schema.TypeString,
			Default:     hostNasVolumeAccessModeReadWrite,
			Description: "Access mode for the mount point. Can be one of readOnly or readWrite.",
			ForceNew:    true,
			Optional:    true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					hostNasVolumeAccessModeReadOnly,
					hostNasVolumeAccessModeReadWrite,
				},
				false,
			),
		},
		"remote_hosts": {
			Type:        schema.TypeList,
			Description: "The hostnames or IP addresses of the remote server or servers. Only one element should be present for NFS v3 but multiple can be present for NFS v4.1.",
			Elem:        &schema.Schema{Type: schema.TypeString},
			ForceNew:    true,
			MinItems:    1,
			Required:    true,
		},
		"remote_path": {
			Type:        schema.TypeString,
			Description: "The remote path of the mount point.",
			ForceNew:    true,
			Required:    true,
		},
		"security_type": {
			Type:        schema.TypeString,
			Description: "The security type to use.",
			ForceNew:    true,
			Optional:    true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					hostNasVolumeSecurityTypeAuthSys,
					hostNasVolumeSecurityTypeSecKrb5,
					hostNasVolumeSecurityTypeSecKrb5i,
				},
				false,
			),
		},
		"type": {
			Type:        schema.TypeString,
			Default:     "NFS",
			Description: "The type of NAS volume. Can be one of NFS (to denote v3) or NFS41 (to denote NFS v4.1).",
			ForceNew:    true,
			Optional:    true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					string(types.HostFileSystemVolumeFileSystemTypeNFS),
					string(types.HostFileSystemVolumeFileSystemTypeNFS41),
				},
				false,
			),
		},
		"protocol_endpoint": {
			Type:        schema.TypeString,
			Description: "Indicates that this NAS volume is a protocol endpoint. This field is only populated if the host supports virtual datastores.",
			Computed:    true,
		},
	}
}

// expandHostNasVolumeSpec reads certain ResourceData keys and returns a
// HostNasVolumeSpec.
func expandHostNasVolumeSpec(d *schema.ResourceData) (*types.HostNasVolumeSpec, error) {
	remoteHosts := d.Get("remote_hosts").([]interface{})
	var remoteHost string
	if len(remoteHosts) == 0 || (len(remoteHosts) > 0 && remoteHosts[0] == nil) {
		return nil, fmt.Errorf("remote hosts cannot be empty")
	} else {
		remoteHost = remoteHosts[0].(string)
	}
	obj := &types.HostNasVolumeSpec{
		AccessMode:      d.Get("access_mode").(string),
		LocalPath:       d.Get("name").(string),
		RemoteHost:      remoteHost,
		RemoteHostNames: structure.SliceInterfacesToStrings(d.Get("remote_hosts").([]interface{})),
		RemotePath:      d.Get("remote_path").(string),
		SecurityType:    d.Get("security_type").(string),
		Type:            d.Get("type").(string),
	}

	return obj, nil
}

// flattenHostNasVolume reads various fields from a HostNasVolume into the
// passed in ResourceData.
//
// Note the name attribute is not set here, bur rather set in
// flattenDatastoreSummary and sourced from there.
func flattenHostNasVolume(d *schema.ResourceData, obj *types.HostNasVolume) error {
	d.Set("remote_path", obj.RemotePath)
	d.Set("security_type", obj.SecurityType)
	d.Set("protocol_endpoint", obj.ProtocolEndpoint)

	if err := d.Set("remote_hosts", obj.RemoteHostNames); err != nil {
		return err
	}
	return nil
}

// isNasVolume returns true if the HostFileSystemVolumeFileSystemType matches
// one of the possible filesystem types that a NAS datastore supports.
func isNasVolume(t types.HostFileSystemVolumeFileSystemType) bool {
	switch t {
	case types.HostFileSystemVolumeFileSystemTypeNFS, types.HostFileSystemVolumeFileSystemTypeNFS41:
		return true
	}
	return false
}
