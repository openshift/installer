package baremetal

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/coreos/stream-metadata-go/arch"
	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
	"libvirt.org/go/libvirtxml"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/rhcos"
)

func newDomain(name string) libvirtxml.Domain {
	domainDef := libvirtxml.Domain{
		Name: name,
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type: "hvm",
			},
		},
		Devices: &libvirtxml.DomainDeviceList{
			Graphics: []libvirtxml.DomainGraphic{
				{
					Spice: &libvirtxml.DomainGraphicSpice{
						AutoPort: "yes",
					},
				},
			},
			Channels: []libvirtxml.DomainChannel{
				{
					Source: &libvirtxml.DomainChardevSource{
						UNIX: &libvirtxml.DomainChardevSourceUNIX{},
					},
					Target: &libvirtxml.DomainChannelTarget{
						VirtIO: &libvirtxml.DomainChannelTargetVirtIO{
							Name: "org.qemu.guest_agent.0",
						},
					},
				},
			},
		},
		Features: &libvirtxml.DomainFeatureList{
			PAE:  &libvirtxml.DomainFeature{},
			ACPI: &libvirtxml.DomainFeature{},
			APIC: &libvirtxml.DomainFeatureAPIC{},
		},

		CPU: &libvirtxml.DomainCPU{
			Mode: "host-passthrough",
		},
		Memory: &libvirtxml.DomainMemory{
			Value: 6,
			Unit:  "GiB",
		},
		VCPU: &libvirtxml.DomainVCPU{
			Value: 4,
		},
	}

	targetPort := uint(0)
	console := libvirtxml.DomainConsole{
		Target: &libvirtxml.DomainConsoleTarget{
			Port: &targetPort,
		},
	}

	domainDef.Devices.Consoles = append(domainDef.Devices.Consoles, console)

	domainDef.Devices.Graphics = []libvirtxml.DomainGraphic{
		{
			VNC: &libvirtxml.DomainGraphicVNC{
				AutoPort: "yes",
				Listeners: []libvirtxml.DomainGraphicListener{
					{
						Address: &libvirtxml.DomainGraphicListenerAddress{
							Address: "",
						},
					},
				},
			},
		},
	}

	if v := os.Getenv("TERRAFORM_LIBVIRT_TEST_DOMAIN_TYPE"); v != "" {
		domainDef.Type = v
	} else {
		domainDef.Type = "kvm"
	}

	rngDev := os.Getenv("TF_LIBVIRT_RNG_DEV")
	if rngDev == "" {
		rngDev = "/dev/urandom"
	}

	domainDef.Devices.RNGs = []libvirtxml.DomainRNG{
		{
			Model: "virtio",
			Backend: &libvirtxml.DomainRNGBackend{
				Random: &libvirtxml.DomainRNGBackendRandom{Device: rngDev},
			},
		},
	}

	return domainDef
}

func createStoragePool(virConn *libvirt.Libvirt, config baremetalConfig) (libvirt.StoragePool, error) {
	// TODO: check if unique
	bootstrapPool := libvirtxml.StoragePool{
		Type: "dir",
		Name: fmt.Sprintf("%s-bootstrap", config.ClusterID),
		Target: &libvirtxml.StoragePoolTarget{
			Path: fmt.Sprintf("/var/lib/libvirt/openshift-images/%s-bootstrap", config.ClusterID),
		},
	}

	bootstrapPoolXML, err := xml.Marshal(bootstrapPool)
	if err != nil {
		return libvirt.StoragePool{}, err
	}

	pool, err := virConn.StoragePoolDefineXML(string(bootstrapPoolXML), 0)
	if err != nil {
		return libvirt.StoragePool{}, err
	}

	err = virConn.StoragePoolBuild(pool, libvirt.StoragePoolBuildNew)
	if err != nil {
		return libvirt.StoragePool{}, err
	}

	err = virConn.StoragePoolSetAutostart(pool, 1)
	if err != nil {
		return libvirt.StoragePool{}, err
	}

	err = virConn.StoragePoolCreate(pool, libvirt.StoragePoolCreateNormal)
	if err != nil {
		return libvirt.StoragePool{}, err
	}

	err = virConn.StoragePoolRefresh(pool, 0)
	if err != nil {
		return libvirt.StoragePool{}, err
	}
	return pool, nil
}

func getLiveISO(config baremetalConfig) (string, error) {
	fetcher := rhcos.NewBaseISOFetcher(
		rhcos.NewReleasePayload(
			rhcos.ExtractConfig{},
			config.ReleaseImagePullSpec,
			config.PullSecret,
			config.MirrorConfig,
		),
		nil)
	return fetcher.GetBaseISOFilename(context.Background(), arch.RpmArch(runtime.GOARCH))
}

func bootstrapIgnition(config baremetalConfig) (*isoeditor.IgnitionContent, error) {
	ign := &igntypes.Config{}
	// TODO(zaneb): Put swap config into system ignition rather than modifying user ignition
	if err := json.Unmarshal([]byte(config.IgnitionBootstrap), &ign); err != nil {
		return nil, fmt.Errorf("failed to unmarshal bootstrap Ignition config: %w", err)
	}

	fsLabel := "var"
	partDev := fmt.Sprintf("/dev/disk/by-partlabel/%s", fsLabel)
	format := "xfs"
	path := "/var"
	ign.Storage.Disks = append(ign.Storage.Disks, igntypes.Disk{
		Device: "/dev/vda",
		Partitions: []igntypes.Partition{
			{
				Number: 1,
				Label:  &fsLabel,
			},
		},
	})
	ign.Storage.Filesystems = append(ign.Storage.Filesystems, igntypes.Filesystem{
		Device: partDev,
		Label:  &fsLabel,
		Format: &format,
		Path:   &path,
	})
	systemdPartDev := strings.ReplaceAll(strings.ReplaceAll(strings.TrimLeft(partDev, "/"), "-", "\\x2d"), "/", "-")
	enabled := true
	mountUnit := fmt.Sprintf(`[Unit]
Requires=systemd-fsck@%s.service
After=systemd-fsck@%s.service
[Mount]
Where=%s
What=%s
Type=%s

[Install]
RequiredBy=localfs.target
`,
		systemdPartDev, systemdPartDev, path, partDev, format)
	ign.Systemd.Units = append(ign.Systemd.Units, igntypes.Unit{
		Name:     fmt.Sprintf("%s.mount", fsLabel),
		Contents: &mountUnit,
		Enabled:  &enabled,
	})
	ign.Storage.Files = append(ign.Storage.Files, ignition.FileFromString("/etc/no-var-tmpfs", "root", 0o440, ""))

	ignData, err := ignition.Marshal(ign)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bootstrap Ignition config: %w", err)
	}
	return &isoeditor.IgnitionContent{Config: ignData}, nil
}

func createLiveVolume(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool) (libvirt.StorageVol, error) {
	isoFile, err := getLiveISO(config)
	if err != nil {
		return libvirt.StorageVol{}, err
	}
	defer os.Remove(isoFile)

	ignition, err := bootstrapIgnition(config)
	if err != nil {
		return libvirt.StorageVol{}, err
	}
	var kargs string
	if config.FIPS {
		kargs += " fips=1"
	}
	stream, err := isoeditor.NewRHCOSStreamReader(
		isoFile,
		ignition,
		nil,
		[]byte(kargs),
	)
	if err != nil {
		return libvirt.StorageVol{}, err
	}
	defer stream.Close()
	size, err := stream.Seek(0, io.SeekEnd)
	if err != nil {
		return libvirt.StorageVol{}, err
	}
	_, err = stream.Seek(0, io.SeekStart)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	bootstrapLiveVolume := libvirtxml.StorageVolume{
		Name: fmt.Sprintf("%s-live-provisioner", config.ClusterID),
		Type: "file",
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: "iso",
			},
			Permissions: &libvirtxml.StorageVolumeTargetPermissions{
				Mode: "644",
			},
		},
		Capacity: &libvirtxml.StorageVolumeSize{
			Value: uint64(size),
		},
	}
	bootstrapLiveVolumeXML, err := xml.Marshal(bootstrapLiveVolume)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	liveVolume, err := virConn.StorageVolCreateXML(pool, string(bootstrapLiveVolumeXML), 0)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	err = virConn.StorageVolUpload(liveVolume, stream, 0, uint64(size), 0)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	return liveVolume, nil
}

func createScratchVolume(virConn *libvirt.Libvirt, clusterID string, pool libvirt.StoragePool) (libvirt.StorageVol, error) {
	vol := libvirtxml.StorageVolume{
		Name: fmt.Sprintf("%s-scratch", clusterID),
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: "qcow2",
			},
			Permissions: &libvirtxml.StorageVolumeTargetPermissions{
				Mode: "644",
			},
		},
		Capacity: &libvirtxml.StorageVolumeSize{
			Unit:  "GiB",
			Value: 20,
		},
	}
	scratchVolumeXML, err := xml.Marshal(vol)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	scratchVolume, err := virConn.StorageVolCreateXML(pool, string(scratchVolumeXML), 0)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	return scratchVolume, nil
}

func getHostCapabilities(virConn *libvirt.Libvirt) (libvirtxml.Caps, error) {
	var caps libvirtxml.Caps

	capsBytes, err := virConn.Capabilities()
	if err != nil {
		return caps, err
	}

	err = xml.Unmarshal(capsBytes, &caps)
	if err != nil {
		return caps, err
	}

	return caps, nil
}

func createBootstrapDomain(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool, liveCDVolume, scratchVolume libvirt.StorageVol) error {
	bootstrapDom := newDomain(fmt.Sprintf("%s-bootstrap", config.ClusterID))

	capabilities, err := getHostCapabilities(virConn)
	if err != nil {
		return fmt.Errorf("failed to get libvirt capabilities: %w", err)
	}

	arch := capabilities.Host.CPU.Arch
	bootstrapDom.OS.Type.Arch = arch

	if arch == "aarch64" {
		// for aarch64 speciffying this will automatically select the firmware and NVRAM file
		// reference: https://libvirt.org/formatdomain.html#bios-bootloader
		bootstrapDom.OS.Firmware = "efi"
	}

	// For aarch64, s390x, ppc64 and ppc64le spice is not supported
	if arch == "aarch64" || arch == "s390x" || strings.HasPrefix(arch, "ppc64") {
		bootstrapDom.Devices.Graphics = nil
	}

	for _, bridge := range config.Bridges {
		netIface := libvirtxml.DomainInterface{
			Model: &libvirtxml.DomainInterfaceModel{
				Type: "virtio",
			},
			MAC: &libvirtxml.DomainInterfaceMAC{
				Address: bridge.MAC,
			},
			Source: &libvirtxml.DomainInterfaceSource{
				Bridge: &libvirtxml.DomainInterfaceSourceBridge{
					Bridge: bridge.Name,
				},
			},
		}

		bootstrapDom.Devices.Interfaces = append(bootstrapDom.Devices.Interfaces, netIface)
	}

	liveCD := libvirtxml.DomainDisk{
		Device: "cdrom",
		Target: &libvirtxml.DomainDiskTarget{
			Bus: "sata",
			Dev: "sda",
		},
		Driver: &libvirtxml.DomainDiskDriver{
			Name: "qemu",
			Type: "raw",
		},
		Source: &libvirtxml.DomainDiskSource{
			Volume: &libvirtxml.DomainDiskSourceVolume{
				Pool:   pool.Name,
				Volume: liveCDVolume.Name,
			},
		},
		Boot: &libvirtxml.DomainDeviceBoot{
			Order: 1,
		},
	}

	scratchDisk := libvirtxml.DomainDisk{
		Device: "disk",
		Target: &libvirtxml.DomainDiskTarget{
			Bus: "virtio",
			Dev: "vda",
		},
		Driver: &libvirtxml.DomainDiskDriver{
			Name: "qemu",
			Type: "qcow2",
		},
		Source: &libvirtxml.DomainDiskSource{
			Volume: &libvirtxml.DomainDiskSourceVolume{
				Pool:   pool.Name,
				Volume: scratchVolume.Name,
			},
		},
	}

	bootstrapDom.Devices.Disks = append(bootstrapDom.Devices.Disks, liveCD, scratchDisk)

	bootstrapDom.Resource = nil

	bootstrapDomXML, err := xml.Marshal(bootstrapDom)
	if err != nil {
		return err
	}

	dom, err := virConn.DomainDefineXML(string(bootstrapDomXML))
	if err != nil {
		return err
	}

	err = virConn.DomainCreate(dom)
	if err != nil {
		return err
	}

	return nil
}

func createBootstrap(config baremetalConfig) error {
	logrus.Debug("libvirt: Creating bootstrap")
	uri, err := url.Parse(config.LibvirtURI)
	if err != nil {
		return err
	}

	virConn, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return err
	}

	logrus.Debug("  Creating storage pool")
	pool, err := createStoragePool(virConn, config)
	if err != nil {
		return err
	}

	logrus.Debug("  Creating live volume")
	liveVolume, err := createLiveVolume(virConn, config, pool)
	if err != nil {
		return err
	}

	logrus.Debug("  Creating scratch volume")
	scratchVolume, err := createScratchVolume(virConn, config.ClusterID, pool)
	if err != nil {
		return err
	}

	logrus.Debug("  Creating bootstrap domain")
	err = createBootstrapDomain(virConn, config, pool, liveVolume, scratchVolume)
	if err != nil {
		return err
	}

	return nil
}

func destroyBootstrap(config baremetalConfig) error {
	logrus.Debug("libvirt: Destroying bootstrap")

	uri, err := url.Parse(config.LibvirtURI)
	if err != nil {
		return err
	}

	virConn, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("%s-bootstrap", config.ClusterID)

	dom, err := virConn.DomainLookupByName(name)
	if err != nil {
		return err
	}

	logrus.Debug("  Destroying domain")
	err = virConn.DomainDestroy(dom)
	if err != nil {
		return err
	}

	logrus.Debug("  Undefining domain")

	if err := virConn.DomainUndefineFlags(dom, libvirt.DomainUndefineNvram|libvirt.DomainUndefineSnapshotsMetadata|libvirt.DomainUndefineManagedSave|libvirt.DomainUndefineCheckpointsMetadata); err != nil {
		var libvirtErr *libvirt.Error

		if !errors.As(err, &libvirtErr) {
			return fmt.Errorf("failed to cast to libvirt.Error: %w", err)
		}

		if libvirtErr.Code == uint32(libvirt.ErrNoSupport) || libvirtErr.Code == uint32(libvirt.ErrInvalidArg) {
			logrus.Printf("libvirt does not support undefine flags: will try again without flags")
			if err := virConn.DomainUndefine(dom); err != nil {
				return fmt.Errorf("couldn't undefine libvirt domain: %w", err)
			}
		} else {
			return fmt.Errorf("couldn't undefine libvirt domain with flags: %w", err)
		}
	}

	pool, err := virConn.StoragePoolLookupByName(name)
	if err != nil {
		return err
	}

	vol, err := virConn.StorageVolLookupByName(pool, fmt.Sprintf("%s-live-provisioner", config.ClusterID))
	if err != nil {
		return err
	}

	logrus.Debug("  Deleting live volume")
	err = virConn.StorageVolDelete(vol, libvirt.StorageVolDeleteNormal)
	if err != nil {
		return err
	}

	vol, err = virConn.StorageVolLookupByName(pool, fmt.Sprintf("%s-scratch", config.ClusterID))
	if err != nil {
		return err
	}

	logrus.Debug("  Deleting scratch volume")
	err = virConn.StorageVolDelete(vol, libvirt.StorageVolDeleteNormal)
	if err != nil {
		return err
	}

	logrus.Debug("  Destroying pool")
	err = virConn.StoragePoolDestroy(pool)
	if err != nil {
		return err
	}

	logrus.Debug("  Deleting pool")
	err = virConn.StoragePoolDelete(pool, libvirt.StoragePoolDeleteNormal)
	if err != nil {
		return err
	}

	return nil
}
