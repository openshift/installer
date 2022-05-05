package ovirtclient

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
	ovirtsdk "github.com/ovirt/go-ovirt"
)

type vmBuilderComponent func(params OptionalVMParameters, builder *ovirtsdk.VmBuilder)

func vmBuilderComment(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if comment := params.Comment(); comment != "" {
		builder.Comment(comment)
	}
}

func vmBuilderCPU(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if cpu := params.CPU(); cpu != nil {
		cpuBuilder := ovirtsdk.NewCpuBuilder()
		if cpuTopo := cpu.Topo(); cpuTopo != nil {
			cpuBuilder.TopologyBuilder(ovirtsdk.
				NewCpuTopologyBuilder().
				Cores(int64(cpu.Topo().Cores())).
				Threads(int64(cpu.Topo().Threads())).
				Sockets(int64(cpu.Topo().Sockets())))
		}
		if mode := cpu.Mode(); mode != nil {
			cpuBuilder.Mode(ovirtsdk.CpuMode(*mode))
		}
		builder.CpuBuilder(cpuBuilder)
	}
}

func vmBuilderHugePages(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	var customProperties []*ovirtsdk.CustomProperty
	if hugePages := params.HugePages(); hugePages != nil {
		customProp, err := ovirtsdk.NewCustomPropertyBuilder().
			Name("hugepages").
			Value(strconv.FormatUint(uint64(*hugePages), 10)).
			Build()
		if err != nil {
			panic(newError(EBug, "Failed to build 'hugepages' custom property from value %d", hugePages))
		}
		customProperties = append(customProperties, customProp)
	}
	if len(customProperties) > 0 {
		builder.CustomPropertiesOfAny(customProperties...)
	}
}

func vmBuilderMemory(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if memory := params.Memory(); memory != nil {
		builder.Memory(*memory)
	}
}

func vmBuilderInitialization(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if init := params.Initialization(); init != nil {
		initBuilder := ovirtsdk.NewInitializationBuilder()

		if init.CustomScript() != "" {
			initBuilder.CustomScript(init.CustomScript())
		}
		if init.HostName() != "" {
			initBuilder.HostName(init.HostName())
		}
		builder.InitializationBuilder(initBuilder)
	}
}

func vmPlacementPolicyParameterConverter(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if pp := params.PlacementPolicy(); pp != nil {
		placementPolicyBuilder := ovirtsdk.NewVmPlacementPolicyBuilder()
		if affinity := (*pp).Affinity(); affinity != nil {
			placementPolicyBuilder.Affinity(ovirtsdk.VmAffinity(*affinity))
		}
		hosts := make([]ovirtsdk.HostBuilder, len((*pp).HostIDs()))
		for i, hostID := range (*pp).HostIDs() {
			hostBuilder := ovirtsdk.NewHostBuilder().Id(string(hostID))
			hosts[i] = *hostBuilder
		}
		placementPolicyBuilder.HostsBuilderOfAny(hosts...)
		builder.PlacementPolicyBuilder(placementPolicyBuilder)
	}
}

func (o *oVirtClient) CreateVM(clusterID ClusterID, templateID TemplateID, name string, params OptionalVMParameters, retries ...RetryStrategy) (result VM, err error) {
	retries = defaultRetries(retries, defaultLongTimeouts(o))

	if err := validateVMCreationParameters(clusterID, templateID, name, params); err != nil {
		return nil, err
	}

	if params == nil {
		params = &vmParams{}
	}

	message := fmt.Sprintf("creating VM %s", name)
	vm, err := createSDKVM(clusterID, templateID, name, params)
	if err != nil {
		return nil, err
	}

	err = retry(
		message,
		o.logger,
		retries,
		func() error {
			vmCreateRequest := o.conn.SystemService().VmsService().Add().Vm(vm)
			if clone := params.Clone(); clone != nil {
				vmCreateRequest.Clone(*clone)
			}
			response, err := vmCreateRequest.Send()
			if err != nil {
				return err
			}
			vm, ok := response.Vm()
			if !ok {
				return newError(EFieldMissing, "missing VM in VM create response")
			}
			result, err = convertSDKVM(vm, o, o.logger, "creating VM")
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert VM",
				)
			}
			return nil
		},
	)
	return result, err
}

func createSDKVM(
	clusterID ClusterID,
	templateID TemplateID,
	name string,
	params OptionalVMParameters,
) (*ovirtsdk.Vm, error) {
	builder := ovirtsdk.NewVmBuilder()
	builder.Cluster(ovirtsdk.NewClusterBuilder().Id(string(clusterID)).MustBuild())
	builder.Template(ovirtsdk.NewTemplateBuilder().Id(string(templateID)).MustBuild())
	builder.Name(name)
	parts := []vmBuilderComponent{
		vmBuilderComment,
		vmBuilderCPU,
		vmBuilderHugePages,
		vmBuilderInitialization,
		vmBuilderMemory,
		vmPlacementPolicyParameterConverter,
		vmBuilderMemoryPolicy,
		vmInstanceTypeID,
		vmTypeCreator,
		vmOSCreator,
		vmSerialConsoleCreator,
	}

	for _, part := range parts {
		part(params, builder)
	}

	if params != nil {
		var diskAttachments []*ovirtsdk.DiskAttachment
		for i, d := range params.Disks() {
			diskAttachment := ovirtsdk.NewDiskAttachmentBuilder()
			diskBuilder := ovirtsdk.NewDiskBuilder()
			diskBuilder.Id(string(d.DiskID()))
			if sparse := d.Sparse(); sparse != nil {
				diskBuilder.Sparse(*sparse)
			}
			if format := d.Format(); format != nil {
				diskBuilder.Format(ovirtsdk.DiskFormat(*format))
			}
			if storageDomainID := d.StorageDomainID(); storageDomainID != nil {
				diskBuilder.StorageDomainsBuilderOfAny(*ovirtsdk.NewStorageDomainBuilder().Id(string(*storageDomainID)))
			}
			diskAttachment.DiskBuilder(diskBuilder)
			sdkDisk, err := diskAttachment.Build()
			if err != nil {
				return nil, wrap(err, EBadArgument, "Failed to convert disk %d.", i)
			}
			diskAttachments = append(diskAttachments, sdkDisk)
		}
		builder.DiskAttachmentsOfAny(diskAttachments...)
	}

	vm, err := builder.Build()
	if err != nil {
		return nil, wrap(err, EBug, "failed to build VM")
	}
	return vm, nil
}

func vmSerialConsoleCreator(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	serial := params.SerialConsole()
	if serial == nil {
		return
	}
	builder.ConsoleBuilder(ovirtsdk.NewConsoleBuilder().Enabled(*serial))
}

func vmOSCreator(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if os, ok := params.OS(); ok {
		osBuilder := ovirtsdk.NewOperatingSystemBuilder()
		if t := os.Type(); t != nil {
			osBuilder.Type(*t)
		}
		builder.OsBuilder(osBuilder)
	}
}

func vmTypeCreator(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if vmType := params.VMType(); vmType != nil {
		builder.Type(ovirtsdk.VmType(*vmType))
	}
}

func vmInstanceTypeID(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if instanceTypeID := params.InstanceTypeID(); instanceTypeID != nil {
		builder.InstanceTypeBuilder(ovirtsdk.NewInstanceTypeBuilder().Id(string(*instanceTypeID)))
	}
}

func vmBuilderMemoryPolicy(params OptionalVMParameters, builder *ovirtsdk.VmBuilder) {
	if memPolicyParams := params.MemoryPolicy(); memPolicyParams != nil {
		memoryPolicyBuilder := ovirtsdk.NewMemoryPolicyBuilder()
		if guaranteed := (*memPolicyParams).Guaranteed(); guaranteed != nil {
			memoryPolicyBuilder.Guaranteed(*guaranteed)
		}
		if max := (*memPolicyParams).Max(); max != nil {
			memoryPolicyBuilder.Max(*max)
		}
		if ballooning := (*memPolicyParams).Ballooning(); ballooning != nil {
			memoryPolicyBuilder.Ballooning(*ballooning)
		}
		builder.MemoryPolicyBuilder(memoryPolicyBuilder)
	}
}

func validateVMCreationParameters(clusterID ClusterID, templateID TemplateID, name string, params OptionalVMParameters) error {
	if name == "" {
		return newError(EBadArgument, "name cannot be empty for VM creation")
	}
	if clusterID == "" {
		return newError(EBadArgument, "cluster ID cannot be empty for VM creation")
	}
	if templateID == "" {
		return newError(EBadArgument, "template ID cannot be empty for VM creation")
	}
	if params == nil {
		return nil
	}

	memory := params.Memory()
	if memory == nil {
		mem := int64(1024 * 1024 * 1024)
		memory = &mem
	}
	guaranteedMemory := memory
	if memPolicy := params.MemoryPolicy(); memPolicy != nil {
		guaranteed := (*memPolicy).Guaranteed()
		if guaranteed != nil {
			guaranteedMemory = guaranteed
		}
	}
	if *guaranteedMemory > *memory {
		return newError(
			EBadArgument,
			"guaranteed memory is larger than the VM memory (%d > %d)",
			*guaranteedMemory,
			*memory,
		)
	}

	disks := params.Disks()
	diskIDs := map[DiskID]int{}
	for i, d := range disks {
		if previousID, ok := diskIDs[d.DiskID()]; ok {
			return newError(
				EBadArgument,
				"Disk %s appears twice, in position %d and %d.",
				d.DiskID(),
				previousID,
				i,
			)
		}
	}

	if vmType := params.VMType(); vmType != nil {
		if err := vmType.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (m *mockClient) CreateVM(
	clusterID ClusterID,
	templateID TemplateID,
	name string,
	params OptionalVMParameters,
	retries ...RetryStrategy,
) (result VM, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(m))

	if err := validateVMCreationParameters(clusterID, templateID, name, params); err != nil {
		return nil, err
	}
	if params == nil {
		params = &vmParams{}
	}
	if name == "" {
		return nil, newError(EBadArgument, "The name parameter is required for VM creation.")
	}
	err = retry(
		fmt.Sprintf("creating VM %s", name),
		m.logger,
		retries,
		func() error {
			m.lock.Lock()
			defer m.lock.Unlock()
			if _, ok := m.clusters[clusterID]; !ok {
				return newError(ENotFound, "cluster with ID %s not found", clusterID)
			}
			tpl, ok := m.templates[templateID]
			if !ok {
				return newError(ENotFound, "template with ID %s not found", templateID)
			}
			if tpl.status != TemplateStatusOK {
				return newError(EConflict, "template in status \"%s\"", tpl.status)
			}

			for _, vm := range m.vms {
				if vm.name == name {
					return newError(EConflict, "A VM with the name \"%s\" already exists.", name)
				}
			}

			cpu := m.createVMCPU(params, tpl)

			vm := m.createVM(name, params, clusterID, templateID, cpu)

			m.attachVMDisksFromTemplate(tpl, vm, params)

			if clone := params.Clone(); clone != nil && *clone {
				vm.templateID = DefaultBlankTemplateID
			}

			m.vmIPs[vm.id] = map[string][]net.IP{}
			m.addGraphicsConsoles(vm)

			result = vm
			return nil
		},
	)

	return result, err
}

func (m *mockClient) addGraphicsConsoles(vm *vm) {
	m.graphicsConsolesByVM[vm.id] = []*vmGraphicsConsole{
		{
			client: m,
			id:     VMGraphicsConsoleID(m.GenerateUUID()),
			vmID:   vm.id,
		},
		{
			client: m,
			id:     VMGraphicsConsoleID(m.GenerateUUID()),
			vmID:   vm.id,
		},
	}
}

func (m *mockClient) createVM(
	name string,
	params OptionalVMParameters,
	clusterID ClusterID,
	templateID TemplateID,
	cpu *vmCPU,
) *vm {
	id := uuid.Must(uuid.NewUUID()).String()
	init := params.Initialization()
	if init == nil {
		init = &initialization{}
	}

	vmType := m.createVMType(params)
	console := false
	if serialConsole := params.SerialConsole(); serialConsole != nil {
		console = *serialConsole
	}

	vm := &vm{
		m,
		VMID(id),
		name,
		params.Comment(),
		clusterID,
		templateID,
		VMStatusDown,
		cpu,
		m.createVMMemory(params),
		nil,
		params.HugePages(),
		init,
		nil,
		m.createPlacementPolicy(params),
		m.createVMMemoryPolicy(params),
		params.InstanceTypeID(),
		vmType,
		m.createVMOS(params),
		console,
	}
	m.vms[VMID(id)] = vm
	return vm
}

func (m *mockClient) createVMMemory(params OptionalVMParameters) int64 {
	memory := int64(1073741824)
	if params.Memory() != nil {
		memory = *params.Memory()
	}
	return memory
}

func (m *mockClient) createVMMemoryPolicy(params OptionalVMParameters) *memoryPolicy {
	memPolicy := &memoryPolicy{
		ballooning: true,
	}
	if memoryPolicyParams := params.MemoryPolicy(); memoryPolicyParams != nil {
		if guaranteedMemory := (*memoryPolicyParams).Guaranteed(); guaranteedMemory != nil {
			memPolicy.guaranteed = guaranteedMemory
		}

		if maxMemory := (*memoryPolicyParams).Max(); maxMemory != nil {
			memPolicy.max = maxMemory
		}

		if memBallooning := (*memoryPolicyParams).Ballooning(); memBallooning != nil {
			memPolicy.ballooning = *memBallooning
		}
	}
	return memPolicy
}

func (m *mockClient) createVMOS(params OptionalVMParameters) *vmOS {
	os := &vmOS{
		t: "other",
	}
	if osParams, ok := params.OS(); ok {
		if osType := osParams.Type(); osType != nil {
			os.t = *osType
		}
	}
	return os
}

func (m *mockClient) createVMType(params OptionalVMParameters) VMType {
	vmType := VMTypeServer
	if paramVMType := params.VMType(); paramVMType != nil {
		vmType = *paramVMType
	}
	return vmType
}

func (m *mockClient) createPlacementPolicy(params OptionalVMParameters) *vmPlacementPolicy {
	var pp *vmPlacementPolicy
	if params.PlacementPolicy() != nil {
		placementPolicyParams := *params.PlacementPolicy()
		pp = &vmPlacementPolicy{
			placementPolicyParams.Affinity(),
			placementPolicyParams.HostIDs(),
		}
	}
	return pp
}

func (m *mockClient) attachVMDisksFromTemplate(tpl *template, vm *vm, params OptionalVMParameters) {
	m.vmDiskAttachmentsByVM[vm.id] = make(
		map[DiskAttachmentID]*diskAttachment,
		len(m.templateDiskAttachmentsByTemplate[tpl.id]),
	)
	for _, attachment := range m.templateDiskAttachmentsByTemplate[tpl.id] {
		disk := m.disks[attachment.diskID]
		var sparse *bool
		for _, diskParam := range params.Disks() {
			if diskParam.DiskID() == disk.ID() {
				sparse = m.updateDiskParams(diskParam, disk, params)
				break
			}
		}
		newDisk := disk.clone(sparse)
		_ = newDisk.Lock()
		newDisk.alias = fmt.Sprintf("disk-%s", generateRandomID(5, m.nonSecureRandom))
		m.disks[newDisk.ID()] = newDisk

		go func() {
			time.Sleep(time.Second)
			newDisk.Unlock()
		}()

		diskAttachment := &diskAttachment{
			client:        m,
			id:            DiskAttachmentID(m.GenerateUUID()),
			vmid:          vm.id,
			diskID:        newDisk.ID(),
			diskInterface: attachment.diskInterface,
			bootable:      attachment.bootable,
			active:        attachment.active,
		}
		m.vmDiskAttachmentsByVM[vm.id][diskAttachment.id] = diskAttachment
		m.vmDiskAttachmentsByDisk[disk.id] = diskAttachment
	}
}

func (m *mockClient) updateDiskParams(
	diskParam OptionalVMDiskParameters,
	disk *diskWithData,
	params OptionalVMParameters,
) *bool {
	sparse := diskParam.Sparse()
	if format := diskParam.Format(); format != nil {
		if *format != disk.Format() {
			m.logger.Warningf(
				"the VM creation client requested a conversion from from %s to %s; the mock library does not support this and the source image data will be used unmodified which may lead to errors",
				disk.format,
				*format,
			)
			disk.format = *format
		}
	}
	if sd := diskParam.StorageDomainID(); sd != nil {
		if params.Clone() != nil && *params.Clone() {
			disk.storageDomainIDs = []StorageDomainID{*sd}
		} else {
			for _, diskSD := range disk.storageDomainIDs {
				if diskSD == *sd {
					disk.storageDomainIDs = []StorageDomainID{*sd}
					break
				}
			}
			// If the SD is not found then we leave the SD unchanged, just as the engine does.
		}
	}
	return sparse
}

func (m *mockClient) createVMCPU(params OptionalVMParameters, tpl *template) *vmCPU {
	var cpu *vmCPU
	cpuParams := params.CPU()
	switch {
	case cpuParams != nil:
		cpu = &vmCPU{}
		if topo := cpuParams.Topo(); topo != nil {
			cpu.topo = &vmCPUTopo{
				cores:   topo.Cores(),
				sockets: topo.Sockets(),
				threads: topo.Threads(),
			}
		}
		if mode := cpuParams.Mode(); mode != nil {
			cpu.mode = mode
		}
	case tpl.cpu != nil:
		cpu = tpl.cpu.clone()
	default:
		cpu = &vmCPU{
			topo: &vmCPUTopo{
				cores:   1,
				sockets: 1,
				threads: 1,
			},
		}
	}
	return cpu
}
