package baremetal

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
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
			Controllers: []libvirtxml.DomainController{
				{
					Type:  "scsi",
					Model: "virtio-scsi",
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

	serialLogPath := fmt.Sprintf("/var/log/libvirt/qemu/%s-serial0.log", name)
	serial := libvirtxml.DomainSerial{
		Log: &libvirtxml.DomainChardevLog{
			File:   serialLogPath,
			Append: "on",
		},
	}
	domainDef.Devices.Serials = append(domainDef.Devices.Serials, serial)

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

func getLiveISO(config baremetalConfig, arch string) (string, error) {
	fetcher := rhcos.NewBaseISOFetcher(
		rhcos.NewReleasePayload(
			rhcos.ExtractConfig{},
			config.ReleaseImagePullSpec,
			config.PullSecret,
			config.MirrorConfig,
		),
		nil)
	return fetcher.GetBaseISOFilename(context.Background(), arch)
}

func bootstrapIgnition(config baremetalConfig) (*isoeditor.IgnitionContent, error) {
	ign := &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
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
	return &isoeditor.IgnitionContent{
		Config: []byte(config.IgnitionBootstrap),
		SystemConfigs: map[string][]byte{
			"30-scratch-disk.ign": ignData,
		},
	}, nil
}

func createLiveVolume(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool) (libvirt.StorageVol, error) {
	capabilities, err := getHostCapabilities(virConn)
	if err != nil {
		return libvirt.StorageVol{}, fmt.Errorf("failed to get libvirt capabilities: %w", err)
	}

	arch := capabilities.Host.CPU.Arch
	isoFile, err := getLiveISO(config, arch)
	if err != nil {
		return libvirt.StorageVol{}, err
	}
	defer os.Remove(isoFile)

	ignition, err := bootstrapIgnition(config)
	if err != nil {
		return libvirt.StorageVol{}, err
	}
	consoleDevice := map[string]string{
		"x86_64":  "ttyS0",
		"aarch64": "ttyAMA0",
		"s390x":   "ttysclp0",
		"ppc64le": "hvc0",
	}
	var kargs string
	if dev, ok := consoleDevice[arch]; ok {
		kargs += " console=" + dev
	}
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

func configureDomainArch(dom *libvirtxml.Domain, arch string) {
	dom.OS.Type.Arch = arch

	switch arch {
	case "x86_64":
		dom.OS.Type.Machine = "q35"
		dom.OS.Firmware = "efi"
	case "aarch64":
		// reference: https://libvirt.org/formatdomain.html#bios-bootloader
		dom.OS.Firmware = "efi"
		fallthrough
	case "s390x", "ppc64le":
		dom.Devices.Graphics = nil
	}
}

func createBootstrapDomain(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool, liveCDVolume, scratchVolume libvirt.StorageVol) error {
	bootstrapDom := newDomain(fmt.Sprintf("%s-bootstrap", config.ClusterID))

	capabilities, err := getHostCapabilities(virConn)
	if err != nil {
		return fmt.Errorf("failed to get libvirt capabilities: %w", err)
	}

	arch := capabilities.Host.CPU.Arch
	configureDomainArch(&bootstrapDom, arch)

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
			Bus: "scsi",
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
	if arch == "x86_64" {
		// x86 traditionally uses IDE or SATA (only the latter is supported
		// with the q35 machine type) for cdrom devices
		liveCD.Target.Bus = "sata"
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
