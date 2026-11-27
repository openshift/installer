package bootstrap

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/types"
	nonetypes "github.com/openshift/installer/pkg/types/none"
)

func TestMain(m *testing.M) {
	// Point data.Assets at the repo's data/data so validation can verify the
	// production install-to-disk template still passes -n.
	_, thisFile, _, ok := runtime.Caller(0)
	if ok {
		repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", ".."))
		data.Assets = http.Dir(filepath.Join(repoRoot, "data", "data"))
	}
	os.Exit(m.Run())
}

func TestValidateBootstrapInPlaceInstallationDisk(t *testing.T) {
	fld := field.NewPath("bootstrapInPlace", "installationDisk")
	cases := []struct {
		name         string
		disk         string
		expectType   field.ErrorType
		expectDetail string
	}{
		// accept
		{name: "virtio", disk: "/dev/vda"},
		{name: "scsi", disk: "/dev/sda"},
		{name: "nvme", disk: "/dev/nvme0n1"},
		{name: "dasd", disk: "/dev/dasda"},
		{name: "wwn", disk: "/dev/disk/by-id/wwn-0x64cd98f04fde100024684cf3034da5c2"},
		{name: "scsi-id", disk: "/dev/disk/by-id/scsi-3600a09803831436533244e4732444a39"},
		{name: "nvme-id", disk: "/dev/disk/by-id/nvme-eui.0102030405060708"},
		{name: "by-path colon-dot", disk: "/dev/disk/by-path/pci-0000:3b:00.0-scsi-0:0:0:0"},
		{name: "ccw", disk: "/dev/disk/by-path/ccw-0.0.1234"},
		{name: "multipath mapper", disk: "/dev/mapper/mpatha"},
		{name: "multipath by-id", disk: "/dev/disk/by-id/dm-uuid-mpath-3600a09803831436533244e4732444a39"},
		{name: "udev punctuation", disk: "/dev/disk/by-id/ata-ACME_Disk_Model#1+rev-2:=@_x"},
		{name: "encoded name", disk: `/dev/disk/by-id/scsi-SATA_Foo\x20Bar`},
		{name: "non-whitespace utf-8", disk: "/dev/disk/by-id/disk-名前"},

		// reject: exact legacy token gets the migration hint
		{
			name:         "legacy copy-network",
			disk:         "/dev/vda --copy-network",
			expectType:   field.ErrorTypeInvalid,
			expectDetail: installationDiskCopyNetworkMigrationDetail,
		},
		// reject: other extra arguments use the generic one-path message
		{
			name:         "different extra option",
			disk:         "/dev/vda --force",
			expectType:   field.ErrorTypeInvalid,
			expectDetail: installationDiskOnePathDetail,
		},
		{
			name:         "lookalike token",
			disk:         "/dev/vda --copy-networking",
			expectType:   field.ErrorTypeInvalid,
			expectDetail: installationDiskOnePathDetail,
		},
		{
			name:         "legacy token plus another option",
			disk:         "/dev/vda --copy-network --force",
			expectType:   field.ErrorTypeInvalid,
			expectDetail: installationDiskOnePathDetail,
		},
		{
			name:         "relative",
			disk:         "dev/vda",
			expectType:   field.ErrorTypeInvalid,
			expectDetail: "must be an absolute disk path",
		},
		{name: "dot relative", disk: "./dev/vda", expectType: field.ErrorTypeInvalid, expectDetail: "must be an absolute disk path"},
		{name: "root", disk: "/", expectType: field.ErrorTypeInvalid, expectDetail: "must be an absolute disk path"},
		{name: "dot root", disk: "/.", expectType: field.ErrorTypeInvalid, expectDetail: "must be an absolute disk path"},
		{name: "dotdot root", disk: "/..", expectType: field.ErrorTypeInvalid, expectDetail: "must be an absolute disk path"},
		{name: "parent component", disk: "/dev/../sda", expectType: field.ErrorTypeInvalid, expectDetail: "must be an absolute disk path"},
		{name: "tab separator other flag", disk: "/dev/vda\t--force", expectType: field.ErrorTypeInvalid, expectDetail: installationDiskOnePathDetail},
		{name: "newline", disk: "/dev/vda\n--force", expectType: field.ErrorTypeInvalid, expectDetail: installationDiskOnePathDetail},
		{name: "unicode whitespace", disk: "/dev/vda\u00a0extra", expectType: field.ErrorTypeInvalid, expectDetail: installationDiskOnePathDetail},
		{name: "semicolon", disk: "/dev/vda;reboot", expectType: field.ErrorTypeInvalid, expectDetail: "must not contain shell expression syntax"},
		{name: "command substitution", disk: "/dev/vda$(id)", expectType: field.ErrorTypeInvalid, expectDetail: "must not contain shell expression syntax"},
		{name: "backtick", disk: "/dev/vda`id`", expectType: field.ErrorTypeInvalid, expectDetail: "must not contain shell expression syntax"},
		{name: "pipe", disk: "/dev/vda|tee", expectType: field.ErrorTypeInvalid, expectDetail: "must not contain shell expression syntax"},
		{name: "ampersand", disk: "/dev/vda&", expectType: field.ErrorTypeInvalid, expectDetail: "must not contain shell expression syntax"},
		// substring lookalike on a single path token is not the legacy form
		{name: "copy-network substring in path token", disk: "/dev/disk/by-id/fake--copy-network"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateBootstrapInPlaceInstallationDisk(tc.disk, fld)
			if tc.expectType == "" {
				assert.Nil(t, err)
				return
			}
			if !assert.NotNil(t, err) {
				return
			}
			assert.Equal(t, tc.expectType, err.Type)
			assert.Equal(t, fld.String(), err.Field)
			assert.Equal(t, tc.expectDetail, err.Detail)
		})
	}
}

func TestVerifyBootstrapInPlaceInstallationDisk(t *testing.T) {
	base := func(disk string) *types.InstallConfig {
		return &types.InstallConfig{
			ControlPlane: &types.MachinePool{Replicas: ptr.To[int64](1)},
			BootstrapInPlace: &types.BootstrapInPlace{
				InstallationDisk: disk,
			},
			Platform: types.Platform{None: &nonetypes.Platform{}},
		}
	}

	t.Run("valid disk", func(t *testing.T) {
		assert.NoError(t, verifyBootstrapInPlace(base("/dev/vda")))
	})

	t.Run("empty disk is Required", func(t *testing.T) {
		err := verifyBootstrapInPlace(base(""))
		if !assert.Error(t, err) {
			return
		}
		assert.Contains(t, err.Error(), "bootstrapInPlace.installationDisk")
		assert.Contains(t, err.Error(), "Required")
		assert.NotContains(t, err.Error(), "Invalid")
	})

	t.Run("legacy copy-network is Invalid with migration hint", func(t *testing.T) {
		err := verifyBootstrapInPlace(base("/dev/vda --copy-network"))
		if !assert.Error(t, err) {
			return
		}
		assert.Contains(t, err.Error(), "bootstrapInPlace.installationDisk")
		assert.Contains(t, err.Error(), "Invalid")
		assert.Contains(t, err.Error(), installationDiskCopyNetworkMigrationDetail)
	})

	t.Run("different extra option is Invalid without migration hint", func(t *testing.T) {
		err := verifyBootstrapInPlace(base("/dev/vda --force"))
		if !assert.Error(t, err) {
			return
		}
		assert.Contains(t, err.Error(), "bootstrapInPlace.installationDisk")
		assert.Contains(t, err.Error(), "Invalid")
		assert.Contains(t, err.Error(), installationDiskOnePathDetail)
		assert.NotContains(t, err.Error(), installationDiskCopyNetworkMigrationDetail)
	})

	t.Run("missing bootstrapInPlace", func(t *testing.T) {
		ic := base("/dev/vda")
		ic.BootstrapInPlace = nil
		err := verifyBootstrapInPlace(ic)
		if !assert.Error(t, err) {
			return
		}
		assert.Contains(t, err.Error(), "bootstrapInPlace")
		assert.Contains(t, err.Error(), "Required")
	})

	t.Run("replicas not one", func(t *testing.T) {
		ic := base("/dev/vda")
		ic.ControlPlane.Replicas = ptr.To[int64](3)
		err := verifyBootstrapInPlace(ic)
		if !assert.Error(t, err) {
			return
		}
		assert.Contains(t, err.Error(), "controlPlane.replicas")
	})
}

func TestInstallToDiskTemplateEnablesNetworkCopy(t *testing.T) {
	assert.True(t, installToDiskTemplateEnablesNetworkCopy(),
		"production install-to-disk.sh.template must retain -n before claiming default network copy in migration messages")
}

func TestShellSingleQuote(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{`/dev/vda`, `'/dev/vda'`},
		{`/dev/disk/by-path/pci-0000:3b:00.0`, `'/dev/disk/by-path/pci-0000:3b:00.0'`},
		{`/dev/disk/by-id/scsi-SATA_Foo\x20Bar`, `'/dev/disk/by-id/scsi-SATA_Foo\x20Bar'`},
		{`foo'bar`, `'foo'"'"'bar'`},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, shellSingleQuote(tc.in))
	}
}

func productionInstallToDiskTemplatePath(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	assert.True(t, ok)
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", ".."))
	return filepath.Join(repoRoot, "data", "data", "bootstrap", "bootstrap-in-place", "files", "usr", "local", "bin", "install-to-disk.sh.template")
}

func TestInstallToDiskTemplateRender(t *testing.T) {
	f, err := os.Open(productionInstallToDiskTemplatePath(t))
	assert.NoError(t, err)
	defer f.Close()

	td := &bootstrapTemplateData{
		BootstrapInPlace: &types.BootstrapInPlace{
			InstallationDisk: `/dev/disk/by-id/scsi-SATA_Foo\x20Bar`,
		},
	}
	name, rendered, err := readFile(filepath.Base(productionInstallToDiskTemplatePath(t)), f, td)
	assert.NoError(t, err)
	assert.Equal(t, "install-to-disk.sh", name)

	got := string(rendered)
	assert.Contains(t, got, `disk='/dev/disk/by-id/scsi-SATA_Foo\x20Bar'`)
	assert.Contains(t, got, `coreos-installer install -n -i /opt/openshift/master.ign "$disk"`)
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "coreos-installer install") {
			assert.NotContains(t, line, "--copy-network")
			assert.Contains(t, line, " -n ")
			assert.Contains(t, line, " -i /opt/openshift/master.ign ")
		}
	}
}

func TestInstallToDiskArgumentBoundary(t *testing.T) {
	// Render the production template with benign accepted paths and assert the
	// assignment / installer lines as text. Malformed disks belong in validator
	// tests. Avoid executing a shell recorder (gosec G204/G306).
	accepted := []string{
		"/dev/vda",
		"/dev/disk/by-path/pci-0000:3b:00.0-scsi-0:0:0:0",
		`/dev/disk/by-id/scsi-SATA_Foo\x20Bar`,
		"/dev/mapper/mpatha",
	}

	for _, disk := range accepted {
		t.Run(disk, func(t *testing.T) {
			f, err := os.Open(productionInstallToDiskTemplatePath(t))
			assert.NoError(t, err)
			defer f.Close()

			name, rendered, err := readFile(filepath.Base(productionInstallToDiskTemplatePath(t)), f, &bootstrapTemplateData{
				BootstrapInPlace: &types.BootstrapInPlace{InstallationDisk: disk},
			})
			assert.NoError(t, err)
			assert.Equal(t, "install-to-disk.sh", name)

			got := string(rendered)
			assert.Contains(t, got, "disk="+shellSingleQuote(disk))
			assert.Contains(t, got, `coreos-installer install -n -i /opt/openshift/master.ign "$disk"`)

			var diskLine, installLine string
			for _, line := range strings.Split(got, "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "disk=") {
					diskLine = trimmed
				}
				if strings.Contains(trimmed, "coreos-installer install") {
					installLine = trimmed
				}
			}
			assert.Equal(t, "disk="+shellSingleQuote(disk), diskLine)
			assert.Equal(t, `coreos-installer install -n -i /opt/openshift/master.ign "$disk"`, installLine)
			assert.NotContains(t, installLine, "--copy-network")
		})
	}
}
