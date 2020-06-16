package hardware

import (
	"fmt"

	metal3v1alpha1 "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

const (
	// DefaultProfileName is the default hardware profile to use when
	// no other profile matches.
	DefaultProfileName string = "unknown"
)

// Profile holds the settings for a class of hardware.
type Profile struct {
	// Name holds the profile name
	Name string

	// RootDeviceHints holds the suggestions for placing the storage
	// for the root filesystem.
	RootDeviceHints metal3v1alpha1.RootDeviceHints

	// RootGB is the size of the root volume in GB
	RootGB int

	// LocalGB is the size of something(?)
	LocalGB int

	// CPUArch is the architecture of the CPU.
	CPUArch string
}

var profiles = make(map[string]Profile)

func init() {
	profiles[DefaultProfileName] = Profile{
		Name: DefaultProfileName,
		RootDeviceHints: metal3v1alpha1.RootDeviceHints{
			DeviceName: "/dev/sda",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["libvirt"] = Profile{
		Name: "libvirt",
		RootDeviceHints: metal3v1alpha1.RootDeviceHints{
			DeviceName: "/dev/vda",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["dell"] = Profile{
		Name: "dell",
		RootDeviceHints: metal3v1alpha1.RootDeviceHints{
			HCTL: "0:0:0:0",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["dell-raid"] = Profile{
		Name: "dell-raid",
		RootDeviceHints: metal3v1alpha1.RootDeviceHints{
			HCTL: "0:2:0:0",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["openstack"] = Profile{
		Name: "openstack",
		RootDeviceHints: metal3v1alpha1.RootDeviceHints{
			DeviceName: "/dev/vdb",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["openstack"] = Profile{
		Name: "openstack",
		RootDeviceHints: metal3v1alpha1.RootDeviceHints{
			DeviceName: "/dev/vdb",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}
}

// GetProfile returns the named profile
func GetProfile(name string) (Profile, error) {
	profile, ok := profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("No hardware profile named %q", name)
	}
	return profile, nil
}
