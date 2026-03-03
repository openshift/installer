package baremetal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"libvirt.org/go/libvirtxml"
)

func TestNewDomain(t *testing.T) {
	dom := newDomain("test-bootstrap")

	assert.Equal(t, "test-bootstrap", dom.Name)
	assert.Equal(t, "kvm", dom.Type)
	assert.Equal(t, "hvm", dom.OS.Type.Type)
	assert.Empty(t, dom.OS.Type.Arch, "arch should not be set by newDomain")
	assert.Empty(t, dom.OS.Type.Machine, "machine should not be set by newDomain")
	assert.Equal(t, uint(6), dom.Memory.Value)
	assert.Equal(t, "GiB", dom.Memory.Unit)
	assert.Equal(t, uint(4), dom.VCPU.Value)
	assert.Equal(t, "host-passthrough", dom.CPU.Mode)
	assert.NotEmpty(t, dom.Devices.RNGs, "RNG device should be configured")
	assert.NotEmpty(t, dom.Devices.Graphics, "graphics should be configured")
	assert.NotEmpty(t, dom.Devices.Consoles, "console should be configured")
}

func TestConfigureDomainArch(t *testing.T) {
	tests := []struct {
		name           string
		arch           string
		expectMachine  string
		expectFirmware string
		expectGraphics bool
	}{
		{
			name:           "x86_64 sets q35 and efi",
			arch:           "x86_64",
			expectMachine:  "q35",
			expectFirmware: "efi",
			expectGraphics: true,
		},
		{
			name:           "aarch64 sets efi and removes graphics",
			arch:           "aarch64",
			expectMachine:  "",
			expectFirmware: "efi",
			expectGraphics: false,
		},
		{
			name:           "s390x removes graphics",
			arch:           "s390x",
			expectMachine:  "",
			expectFirmware: "",
			expectGraphics: false,
		},
		{
			name:           "ppc64le removes graphics",
			arch:           "ppc64le",
			expectMachine:  "",
			expectFirmware: "",
			expectGraphics: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dom := newDomain("test-bootstrap")
			configureDomainArch(&dom, tc.arch)

			assert.Equal(t, tc.arch, dom.OS.Type.Arch)
			assert.Equal(t, tc.expectMachine, dom.OS.Type.Machine)
			assert.Equal(t, tc.expectFirmware, dom.OS.Firmware)

			if tc.expectGraphics {
				assert.NotNil(t, dom.Devices.Graphics, "graphics should be present for %s", tc.arch)
			} else {
				assert.Nil(t, dom.Devices.Graphics, "graphics should be removed for %s", tc.arch)
			}
		})
	}
}

func TestConfigureDomainArchPreservesOtherFields(t *testing.T) {
	dom := newDomain("test-bootstrap")

	dom.Devices.Interfaces = []libvirtxml.DomainInterface{
		{
			Model: &libvirtxml.DomainInterfaceModel{Type: "virtio"},
		},
	}

	configureDomainArch(&dom, "x86_64")

	assert.Len(t, dom.Devices.Interfaces, 1, "interfaces should not be modified")
	assert.Equal(t, "hvm", dom.OS.Type.Type, "OS type should remain hvm")
	assert.Equal(t, "kvm", dom.Type, "domain type should remain kvm")
}
