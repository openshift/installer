package baremetal

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/digitalocean/go-libvirt"
	"github.com/sirupsen/logrus"
	"libvirt.org/go/libvirtxml"
)

func newCopier(virConn *libvirt.Libvirt, volume libvirt.StorageVol, size uint64) func(src io.Reader) error {
	copier := func(src io.Reader) error {
		return virConn.StorageVolUpload(volume, src, 0, size, 0)
	}
	return copier
}

func newVolumeFromXML(s string) (libvirtxml.StorageVolume, error) {
	var volumeDef libvirtxml.StorageVolume
	err := xml.Unmarshal([]byte(s), &volumeDef)
	if err != nil {
		return libvirtxml.StorageVolume{}, err
	}
	return volumeDef, nil
}

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
			Value: 20,
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

func newVolume(name string) libvirtxml.StorageVolume {
	return libvirtxml.StorageVolume{
		Name: name,
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: "qcow2",
			},
			Permissions: &libvirtxml.StorageVolumeTargetPermissions{
				Mode: "644",
			},
		},
		Capacity: &libvirtxml.StorageVolumeSize{
			Unit:  "bytes",
			Value: 1,
		},
	}
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

func createBaseVolume(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool) (libvirt.StorageVol, error) {
	bootstrapBaseVolume := newVolume(fmt.Sprintf("%s-bootstrap-base", config.ClusterID))
	image, err := newImage(config.BootstrapOSImage)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	isQCOW2, err := image.IsQCOW2()
	if err != nil {
		return libvirt.StorageVol{}, err
	}
	if isQCOW2 {
		bootstrapBaseVolume.Target.Format.Type = "qcow2"
	}

	size, err := image.Size()
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	bootstrapBaseVolume.Capacity.Unit = "B"
	bootstrapBaseVolume.Capacity.Value = size

	bootstrapBaseVolumeXML, err := xml.Marshal(bootstrapBaseVolume)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	baseVolume, err := virConn.StorageVolCreateXML(pool, string(bootstrapBaseVolumeXML), 0)

	if err != nil {
		return libvirt.StorageVol{}, err
	}

	err = image.Import(newCopier(virConn, baseVolume, bootstrapBaseVolume.Capacity.Value), bootstrapBaseVolume)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	return baseVolume, nil
}

func createMainVolume(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool, baseVolume libvirt.StorageVol) (libvirt.StorageVol, error) {
	bootstrapVolume := newVolume(fmt.Sprintf("%s-bootstrap", config.ClusterID))
	bootstrapVolume.Capacity.Value = 34359738368

	volPath, err := virConn.StorageVolGetPath(baseVolume)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	baseVolumeXMLDesc, err := virConn.StorageVolGetXMLDesc(baseVolume, 0)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	baseVolFromLibvirt, err := newVolumeFromXML(baseVolumeXMLDesc)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	backingStoreVolumeDef := libvirtxml.StorageVolumeBackingStore{
		Path:   volPath,
		Format: baseVolFromLibvirt.Target.Format,
	}

	bootstrapVolume.BackingStore = &backingStoreVolumeDef

	bootstrapVolumeXML, err := xml.Marshal(bootstrapVolume)
	if err != nil {
		return libvirt.StorageVol{}, err
	}

	return virConn.StorageVolCreateXML(pool, string(bootstrapVolumeXML), 0)
}
func createIgnition(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool) error {
	bootstrapIgnition := defIgnition{
		Name:     fmt.Sprintf("%s-bootstrap.ign", config.ClusterID),
		PoolName: pool.Name,
		Content:  config.IgnitionBootstrap,
	}

	_, err := bootstrapIgnition.CreateAndUpload(virConn)
	if err != nil {
		return err
	}

	return nil
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

func createBootstrapDomain(virConn *libvirt.Libvirt, config baremetalConfig, pool libvirt.StoragePool, volume libvirt.StorageVol) error {
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

	disk := libvirtxml.DomainDisk{
		Device: "disk",
		Target: &libvirtxml.DomainDiskTarget{
			Bus: "virtio",
			Dev: "vda",
		},
		Driver: &libvirtxml.DomainDiskDriver{
			Name: "qemu",
			Type: "raw",
		},
		Source: &libvirtxml.DomainDiskSource{
			Index: 0,
			Volume: &libvirtxml.DomainDiskSourceVolume{
				Pool:   pool.Name,
				Volume: volume.Name,
			},
		},
	}

	disk.Driver = &libvirtxml.DomainDiskDriver{
		Name: "qemu",
		Type: "qcow2",
	}
	bootstrapDom.Devices.Disks = append(bootstrapDom.Devices.Disks, disk)

	ignitionKey := fmt.Sprintf("/var/lib/libvirt/openshift-images/%s-bootstrap/%s-bootstrap.ign", config.ClusterID, config.ClusterID)
	bootstrapDom.QEMUCommandline = &libvirtxml.DomainQEMUCommandline{
		Args: []libvirtxml.DomainQEMUCommandlineArg{
			{
				Value: "-fw_cfg",
			},
			{
				Value: fmt.Sprintf("name=%s,file=%s", "opt/com.coreos/config", ignitionKey),
			},
		},
	}

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

	logrus.Debug("  Creating base volume")
	baseVolume, err := createBaseVolume(virConn, config, pool)
	if err != nil {
		return err
	}

	logrus.Debug("  Creating main volume")
	mainVolume, err := createMainVolume(virConn, config, pool, baseVolume)
	if err != nil {
		return err
	}

	logrus.Debug("  Creating ignition")
	err = createIgnition(virConn, config, pool)
	if err != nil {
		return err
	}

	logrus.Debug("  Creating bootstrap domain")
	err = createBootstrapDomain(virConn, config, pool, mainVolume)
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

	vol, err := virConn.StorageVolLookupByName(pool, name)
	if err != nil {
		return err
	}

	logrus.Debug("  Deleting main volume")
	err = virConn.StorageVolDelete(vol, libvirt.StorageVolDeleteNormal)
	if err != nil {
		return err
	}

	vol, err = virConn.StorageVolLookupByName(pool, fmt.Sprintf("%s-bootstrap-base", config.ClusterID))
	if err != nil {
		return err
	}

	logrus.Debug("  Deleting base volume")
	err = virConn.StorageVolDelete(vol, libvirt.StorageVolDeleteNormal)
	if err != nil {
		return err
	}

	vol, err = virConn.StorageVolLookupByName(pool, fmt.Sprintf("%s-bootstrap.ign", config.ClusterID))
	if err != nil {
		return err
	}

	logrus.Debug("  Deleting ignition volume")
	err = virConn.StorageVolDelete(vol, libvirt.StorageVolDeleteNormal)
	if err != nil {
		return err
	}

	logrus.Debug("  Destroying pool")
	err = virConn.StoragePoolDestroy(pool)
	if err != nil {
		return err
	}

	logrus.Debug("  Deleting pool pool")
	err = virConn.StoragePoolDelete(pool, libvirt.StoragePoolDeleteNormal)
	if err != nil {
		return err
	}

	return nil
}
