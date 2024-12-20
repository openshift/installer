package ignition

import (
	"fmt"
	"path/filepath"

	"github.com/clarketm/json"
	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
)

// Marshal is a helper function to use the marshaler function from "github.com/clarketm/json".
// It supports zero values of structs with the omittempty annotation.
// In effect this excludes empty pointer struct fields from the marshaled data,
// instead of inserting nil values into them.
// This is necessary for ignition configs to pass openAPI validation on fields
// that are not supposed to contain nil pointers, but e.g. strings.
// It can be used as a dropin replacement for "encoding/json".Marshal
func Marshal(input interface{}) ([]byte, error) {
	return json.Marshal(input)
}

// FilesFromAsset creates ignition files for each of the files in the specified
// asset.
func FilesFromAsset(pathPrefix string, username string, mode int, asset asset.WritableAsset) []igntypes.File {
	var files []igntypes.File
	for _, f := range asset.Files() {
		files = append(files, FileFromBytes(filepath.Join(pathPrefix, f.Filename), username, mode, f.Data))
	}
	return files
}

// FileFromString creates an ignition-config file with the given contents.
func FileFromString(path string, username string, mode int, contents string) igntypes.File {
	return FileFromBytes(path, username, mode, []byte(contents))
}

// FileFromBytes creates an ignition-config file with the given contents.
func FileFromBytes(path string, username string, mode int, contents []byte) igntypes.File {
	return igntypes.File{
		Node: igntypes.Node{
			Path: path,
			User: igntypes.NodeUser{
				Name: &username,
			},
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &mode,
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(dataurl.EncodeBytes(contents)),
			},
		},
	}
}

// ConvertToRawExtension converts and ignition config to a RawExtension containing the ignition as raw bytes
func ConvertToRawExtension(config igntypes.Config) (runtime.RawExtension, error) {
	rawIgnConfig, err := json.Marshal(config)
	if err != nil {
		return runtime.RawExtension{}, fmt.Errorf("failed to marshal Ignition config: %v", err)
	}

	return runtime.RawExtension{
		Raw: rawIgnConfig,
	}, nil
}

// ConvertToAppendix converts the contents of an ignition file to an appendix.
// In ignition config spec v2 the `Append` boolean value was used to denote whether
// the `Contents` field was an appendix or not. It was also permitted to define
// multiple file configs (appendix or not) that would be merged/overwritten
// sequentially in the order of the json data, which made them non-deterministic.
// In spec v3 this has changed with only one config allowed for each file config,
// and `Append` now being a list of objects that are being appended to `Contents`,
// with the `Contents` field itself never being an appendix.
// This function moves an ignition file's `Contents` object into the `Append` list.
// Since the resulting ignition file of this function has an empty `Contents` field,
// `Overwrite` must be set to false, per the spec.
// The output is an ignition file config that will write a new file with only the
// appendix contents in the case of a file not already existing on disk,
// or append the appendix contents to a file already existing.
func ConvertToAppendix(file *igntypes.File) {
	file.Append = []igntypes.Resource{
		file.Contents,
	}
	file.Contents = igntypes.Resource{}
	file.Overwrite = ignutil.BoolToPtr(false)
}

// AppendVarPartition appends a /var partition to the ignition configuration to avoid growfs.
func AppendVarPartition(config *igntypes.Config) {
	// https://docs.openshift.com/container-platform/4.17/installing/installing_platform_agnostic/installing-platform-agnostic.html#installation-user-infra-machines-advanced_vardisk_installing-platform-agnostic
	config.Storage.Disks = append(config.Storage.Disks, igntypes.Disk{
		Device: "/dev/disk/by-id/coreos-boot-disk",
		Partitions: []igntypes.Partition{
			{
				Label:    ptr.To("var"),
				Number:   5,
				StartMiB: ptr.To(50000),
				SizeMiB:  ptr.To(0),
			},
		},
	})
	config.Storage.Filesystems = append(config.Storage.Filesystems, igntypes.Filesystem{
		Path:         ptr.To("/var"),
		Device:       "/dev/disk/by-partlabel/var",
		Format:       ptr.To("xfs"),
		MountOptions: []igntypes.MountOption{"defaults", "prjquota"},
	})
	// generate a mount unit so that this filesystem gets mounted in the real root.
	config.Systemd.Units = append(config.Systemd.Units, igntypes.Unit{
		Name:    "var.mount",
		Enabled: ptr.To(true),
		Contents: ptr.To(`
[Unit]
Requires=systemd-fsck@dev-disk-by\x2dpartlabel-var.service
After=systemd-fsck@dev-disk-by\x2dpartlabel-var.service

[Mount]
Where=/var
What=/dev/disk/by-partlabel/var
Type=xfs
Options=defaults,prjquota

[Install]
RequiredBy=local-fs.target
`),
	})
}
