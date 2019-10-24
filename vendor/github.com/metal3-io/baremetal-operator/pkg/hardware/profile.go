package hardware

import (
	"fmt"
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
	RootDeviceHints RootDeviceHints

	// RootGB is the size of the root volume in GB
	RootGB int

	// LocalGB is the size of something(?)
	LocalGB int

	// CPUArch is the architecture of the CPU.
	CPUArch string
}

// RootDeviceHints holds the hints for specifying the storage location
// for the root filesystem for the image.
//
// NOTE(dhellmann): Valid ironic hints are: "vendor,
// wwn_vendor_extension, wwn_with_extension, by_path, serial, wwn,
// size, rotational, name, hctl, model"
type RootDeviceHints struct {
	// A device name like "/dev/vda"
	DeviceName string

	// A SCSI bus address like 0:0:0:0
	HCTL string
}

var profiles = make(map[string]Profile)

func init() {
	profiles[DefaultProfileName] = Profile{
		Name: DefaultProfileName,
		RootDeviceHints: RootDeviceHints{
			DeviceName: "/dev/sda",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["libvirt"] = Profile{
		Name: "libvirt",
		RootDeviceHints: RootDeviceHints{
			DeviceName: "/dev/vda",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["dell"] = Profile{
		Name: "dell",
		RootDeviceHints: RootDeviceHints{
			HCTL: "0:0:0:0",
		},
		RootGB:  10,
		LocalGB: 50,
		CPUArch: "x86_64",
	}

	profiles["dell-raid"] = Profile{
		Name: "dell-raid",
		RootDeviceHints: RootDeviceHints{
			HCTL: "0:2:0:0",
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
