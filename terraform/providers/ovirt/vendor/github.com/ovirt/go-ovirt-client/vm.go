package ovirtclient

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

// VMID is a specific type for virtual machine IDs.
type VMID string

// VMClient includes the methods required to deal with virtual machines.
type VMClient interface {
	// CreateVM creates a virtual machine.
	CreateVM(
		clusterID ClusterID,
		templateID TemplateID,
		name string,
		optional OptionalVMParameters,
		retries ...RetryStrategy,
	) (VM, error)
	// GetVM returns a single virtual machine based on an ID.
	GetVM(id VMID, retries ...RetryStrategy) (VM, error)
	// GetVMByName returns a single virtual machine based on a Name.
	GetVMByName(name string, retries ...RetryStrategy) (VM, error)
	// UpdateVM updates the virtual machine with the given parameters.
	// Use UpdateVMParams to obtain a builder for the params.
	UpdateVM(id VMID, params UpdateVMParameters, retries ...RetryStrategy) (VM, error)
	// AutoOptimizeVMCPUPinningSettings sets the CPU settings to optimized.
	AutoOptimizeVMCPUPinningSettings(id VMID, optimize bool, retries ...RetryStrategy) error
	// StartVM triggers a VM start. The actual VM startup will take time and should be waited for via the
	// WaitForVMStatus call.
	StartVM(id VMID, retries ...RetryStrategy) error
	// StopVM triggers a VM power-off. The actual VM stop will take time and should be waited for via the
	// WaitForVMStatus call. The force parameter will cause the shutdown to proceed even if a backup is currently
	// running.
	StopVM(id VMID, force bool, retries ...RetryStrategy) error
	// ShutdownVM triggers a VM shutdown. The actual VM shutdown will take time and should be waited for via the
	// WaitForVMStatus call. The force parameter will cause the shutdown to proceed even if a backup is currently
	// running.
	ShutdownVM(id VMID, force bool, retries ...RetryStrategy) error
	// WaitForVMStatus waits for the VM to reach the desired status.
	WaitForVMStatus(id VMID, status VMStatus, retries ...RetryStrategy) (VM, error)
	// ListVMs returns a list of all virtual machines.
	ListVMs(retries ...RetryStrategy) ([]VM, error)
	// SearchVMs lists all virtual machines matching a certain criteria specified in params.
	SearchVMs(params VMSearchParameters, retries ...RetryStrategy) ([]VM, error)
	// RemoveVM removes a virtual machine specified by id.
	RemoveVM(id VMID, retries ...RetryStrategy) error
	// AddTagToVM Add tag specified by id to a VM.
	AddTagToVM(id VMID, tagID TagID, retries ...RetryStrategy) error
	// AddTagToVMByName Add tag specified by Name to a VM.
	AddTagToVMByName(id VMID, tagName string, retries ...RetryStrategy) error
	// RemoveTagFromVM removes the specified tag from the specified VM.
	RemoveTagFromVM(id VMID, tagID TagID, retries ...RetryStrategy) error
	// ListVMTags lists the tags attached to a VM.
	ListVMTags(id VMID, retries ...RetryStrategy) (result []Tag, err error)
	// GetVMIPAddresses fetches the IP addresses reported by the guest agent in the VM.
	// Optional parameters can be passed to filter the result list.
	//
	// The returned result will be a map of network interface names and the list of IP addresses assigned to them,
	// excluding any IP addresses in the specified parameters.
	GetVMIPAddresses(id VMID, params VMIPSearchParams, retries ...RetryStrategy) (map[string][]net.IP, error)
	// GetVMNonLocalIPAddresses fetches the IP addresses and filters them to return only non-local IP addresses.
	//
	// The returned result will be a map of network interface names and the list of IP addresses assigned to them,
	// excluding any IP addresses in the specified parameters.
	GetVMNonLocalIPAddresses(id VMID, retries ...RetryStrategy) (map[string][]net.IP, error)
	// WaitForVMIPAddresses waits for at least one IP address to be reported that is not in specified ranges.
	//
	// The returned result will be a map of network interface names and the list of IP addresses assigned to them,
	// excluding any IP addresses in the specified parameters.
	WaitForVMIPAddresses(id VMID, params VMIPSearchParams, retries ...RetryStrategy) (map[string][]net.IP, error)
	// WaitForNonLocalVMIPAddress waits for at least one IP address to be reported that is not in the following ranges:
	//
	// - 0.0.0.0/32
	// - 127.0.0.0/8
	// - 169.254.0.0/15
	// - 224.0.0.0/4
	// - 255.255.255.255/32
	// - ::/128
	// - ::1/128
	// - fe80::/64
	// - ff00::/8
	//
	// It also excludes the following interface names:
	//
	// - lo
	// - dummy*
	//
	// The returned result will be a map of network interface names and the list of non-local IP addresses assigned to
	// them.
	WaitForNonLocalVMIPAddress(id VMID, retries ...RetryStrategy) (map[string][]net.IP, error)
}

// VMIPSearchParams contains the parameters for searching or waiting for IP addresses on a VM.
type VMIPSearchParams interface {
	// GetIncludedRanges returns a list of network ranges that the returned IP address must match.
	GetIncludedRanges() []net.IPNet
	// GetExcludedRanges returns a list of IP ranges that should not be taken into consideration when returning IP
	// addresses.
	GetExcludedRanges() []net.IPNet
	// GetIncludedInterfaces returns a list of interface names of which the interface name must match at least
	// one.
	GetIncludedInterfaces() []string
	// GetExcludedInterfaces returns a list of interface names that should be excluded from the search.
	GetExcludedInterfaces() []string
	// GetIncludedInterfacePatterns returns a list of regular expressions of which at least one must match
	// the interface name.
	GetIncludedInterfacePatterns() []*regexp.Regexp
	// GetExcludedInterfacePatterns returns a list of regular expressions that match interface names needing to be
	// excluded from the IP address search.
	GetExcludedInterfacePatterns() []*regexp.Regexp
}

// BuildableVMIPSearchParams is a buildable version of VMIPSearchParams.
type BuildableVMIPSearchParams interface {
	VMIPSearchParams

	WithIncludedRange(ipRange net.IPNet) BuildableVMIPSearchParams
	WithExcludedRange(ipRange net.IPNet) BuildableVMIPSearchParams
	WithIncludedInterface(interfaceName string) BuildableVMIPSearchParams
	WithExcludedInterface(interfaceName string) BuildableVMIPSearchParams
	WithIncludedInterfacePattern(interfaceNamePattern *regexp.Regexp) BuildableVMIPSearchParams
	WithExcludedInterfacePattern(interfaceNamePattern *regexp.Regexp) BuildableVMIPSearchParams
}

// NewVMIPSearchParams returns a buildable parameter set for VM IP searches.
func NewVMIPSearchParams() BuildableVMIPSearchParams {
	return &vmIPSearchParams{}
}

type vmIPSearchParams struct {
	excludedRanges                []net.IPNet
	includedRanges                []net.IPNet
	excludedInterfaceNames        []string
	includedInterfaceNames        []string
	excludedInterfaceNamePatterns []*regexp.Regexp
	includedInterfaceNamePatterns []*regexp.Regexp
}

func (v *vmIPSearchParams) GetIncludedRanges() []net.IPNet {
	return v.includedRanges
}

func (v *vmIPSearchParams) GetIncludedInterfaces() []string {
	return v.includedInterfaceNames
}

func (v *vmIPSearchParams) GetIncludedInterfacePatterns() []*regexp.Regexp {
	return v.includedInterfaceNamePatterns
}

func (v *vmIPSearchParams) WithIncludedRange(ipRange net.IPNet) BuildableVMIPSearchParams {
	v.includedRanges = append(v.includedRanges, ipRange)
	return v
}

func (v *vmIPSearchParams) WithIncludedInterface(interfaceName string) BuildableVMIPSearchParams {
	v.includedInterfaceNames = append(v.includedInterfaceNames, interfaceName)
	return v
}

func (v *vmIPSearchParams) WithIncludedInterfacePattern(interfaceNamePattern *regexp.Regexp) BuildableVMIPSearchParams {
	v.includedInterfaceNamePatterns = append(v.includedInterfaceNamePatterns, interfaceNamePattern)
	return v
}

func (v *vmIPSearchParams) GetExcludedRanges() []net.IPNet {
	return v.excludedRanges
}

func (v *vmIPSearchParams) GetExcludedInterfaces() []string {
	return v.excludedInterfaceNames
}

func (v *vmIPSearchParams) GetExcludedInterfacePatterns() []*regexp.Regexp {
	return v.excludedInterfaceNamePatterns
}

func (v *vmIPSearchParams) WithExcludedRange(ipRange net.IPNet) BuildableVMIPSearchParams {
	v.excludedRanges = append(v.excludedRanges, ipRange)
	return v
}

func (v *vmIPSearchParams) WithExcludedInterface(interfaceName string) BuildableVMIPSearchParams {
	v.excludedInterfaceNames = append(v.excludedInterfaceNames, interfaceName)
	return v
}

func (v *vmIPSearchParams) WithExcludedInterfacePattern(interfaceNamePattern *regexp.Regexp) BuildableVMIPSearchParams {
	v.excludedInterfaceNamePatterns = append(v.excludedInterfaceNamePatterns, interfaceNamePattern)
	return v
}

// VMData is the core of VM providing only data access functions.
type VMData interface {
	// ID returns the unique identifier (UUID) of the current virtual machine.
	ID() VMID
	// Name is the user-defined name of the virtual machine.
	Name() string
	// Comment is the comment added to the VM.
	Comment() string
	// ClusterID returns the cluster this machine belongs to.
	ClusterID() ClusterID
	// TemplateID returns the ID of the base template for this machine.
	TemplateID() TemplateID
	// Status returns the current status of the VM.
	Status() VMStatus
	// CPU returns the CPU structure of a VM.
	CPU() VMCPU
	// Memory return the Memory of a VM in Bytes.
	Memory() int64
	// MemoryPolicy returns the memory policy set on the VM.
	MemoryPolicy() MemoryPolicy
	// TagIDs returns a list of tags for this VM.
	TagIDs() []TagID
	// HugePages returns the hugepage settings for the VM, if any.
	HugePages() *VMHugePages
	// Initialization returns the virtual machine’s initialization configuration.
	Initialization() Initialization
	// HostID returns the ID of the host if available.
	HostID() *HostID
	// PlacementPolicy returns placement policy applied to this VM, if any. It may be nil if no placement policy is set.
	// The second returned value will be false if no placement policy exists.
	PlacementPolicy() (placementPolicy VMPlacementPolicy, ok bool)
	// InstanceTypeID is the source type ID for the instance parameters.
	InstanceTypeID() *InstanceTypeID
	// VMType returns the VM type for the current VM.
	VMType() VMType

	// OS returns the operating system structure.
	OS() VMOS
}

// VMOS is the structure describing the virtual machine operating system, if set.
type VMOS interface {
	Type() string
}

type vmOS struct {
	t string
}

func (v vmOS) Type() string {
	return v.t
}

// VMPlacementPolicy is the structure that holds the rules for VM migration to other hosts.
type VMPlacementPolicy interface {
	Affinity() *VMAffinity
	HostIDs() []HostID
}

// VMAffinity is the affinity used in the placement policy on determining if a VM can be migrated to a different host.
type VMAffinity string

const (
	// VMAffinityMigratable allows automatic and manual VM migrations to other hosts. This is the default.
	VMAffinityMigratable VMAffinity = "migratable"
	// VMAffinityPinned disallows migrating to other hosts.
	VMAffinityPinned VMAffinity = "pinned"
	// VMAffinityUserMigratable allows only manual migrations to different hosts by a user.
	VMAffinityUserMigratable VMAffinity = "user_migratable"
)

// Validate checks the VM affinity for a valid value.
func (v VMAffinity) Validate() error {
	switch v {
	case VMAffinityMigratable:
		return nil
	case VMAffinityPinned:
		return nil
	case VMAffinityUserMigratable:
		return nil
	default:
		return newError(EBadArgument, "invalud value for VMAffinity: %s", v)
	}
}

// VMAffinityValues returns a list of all valid VMAffinity values.
func VMAffinityValues() []VMAffinity {
	return []VMAffinity{
		VMAffinityMigratable,
		VMAffinityPinned,
		VMAffinityUserMigratable,
	}
}

// VMCPU is the CPU configuration of a VM.
type VMCPU interface {
	// Topo is the desired CPU topology for this VM.
	Topo() VMCPUTopo
	// Mode returns the mode of the CPU.
	Mode() *CPUMode
}

type vmCPU struct {
	topo *vmCPUTopo
	mode *CPUMode
}

func (v vmCPU) Mode() *CPUMode {
	return v.mode
}

func (v vmCPU) Topo() VMCPUTopo {
	return v.topo
}

func (v *vmCPU) clone() *vmCPU {
	if v == nil {
		return nil
	}
	return &vmCPU{
		topo: v.topo.clone(),
	}
}

// VMHugePages is the hugepages setting of the VM in bytes.
type VMHugePages uint64

// Validate returns an error if the VM hugepages doesn't have a valid value.
func (h VMHugePages) Validate() error {
	for _, hugePages := range VMHugePagesValues() {
		if hugePages == h {
			return nil
		}
	}
	return newError(
		EBadArgument,
		"Invalid value for VM huge pages: %d must be one of: %s",
		h,
		strings.Join(VMHugePagesValues().Strings(), ", "),
	)
}

const (
	// VMHugePages2M represents the small value of supported huge pages setting which is 2048 Kib.
	VMHugePages2M VMHugePages = 2048
	// VMHugePages1G represents the large value of supported huge pages setting which is 1048576 Kib.
	VMHugePages1G VMHugePages = 1048576
)

// VMHugePagesList is a list of VMHugePages.
type VMHugePagesList []VMHugePages

// Strings creates a string list of the values.
func (l VMHugePagesList) Strings() []string {
	result := make([]string, len(l))
	for i, hugepage := range l {
		result[i] = fmt.Sprint(hugepage)
	}
	return result
}

// VMHugePagesValues returns all possible VMHugepages values.
func VMHugePagesValues() VMHugePagesList {
	return []VMHugePages{
		VMHugePages2M,
		VMHugePages1G,
	}
}

// Initialization defines to the virtual machine’s initialization configuration.
type Initialization interface {
	CustomScript() string
	HostName() string
}

// BuildableInitialization is a buildable version of Initialization.
type BuildableInitialization interface {
	Initialization
	WithCustomScript(customScript string) BuildableInitialization
	WithHostname(hostname string) BuildableInitialization
}

// initialization defines to the virtual machine’s initialization configuration.
// customScript - Cloud-init script which will be executed on Virtual Machine when deployed.
// hostname - Hostname to be set to Virtual Machine when deployed.
type initialization struct {
	customScript string
	hostname     string
}

// NewInitialization creates a new Initialization from the specified parameters.
func NewInitialization(customScript, hostname string) Initialization {
	return &initialization{
		customScript: customScript,
		hostname:     hostname,
	}
}

func (i *initialization) CustomScript() string {
	return i.customScript
}

func (i *initialization) HostName() string {
	return i.hostname
}

func (i *initialization) WithCustomScript(customScript string) BuildableInitialization {
	i.customScript = customScript
	return i
}

func (i *initialization) WithHostname(hostname string) BuildableInitialization {
	i.hostname = hostname
	return i
}

// convertSDKInitialization converts the initialization of a VM. We keep the error return in case we need it later
// as errors may happen as we extend this function and we don't want to touch other functions.
func convertSDKInitialization(sdkObject *ovirtsdk.Vm) (*initialization, error) { //nolint:unparam
	initializationSDK, ok := sdkObject.Initialization()
	if !ok {
		// This happens for some, but not all API calls if the initialization is not set.
		return &initialization{}, nil
	}

	init := initialization{}
	customScript, ok := initializationSDK.CustomScript()
	if ok {
		init.customScript = customScript
	}
	hostname, ok := initializationSDK.HostName()
	if ok {
		init.hostname = hostname
	}
	return &init, nil
}

// VM is the implementation of the virtual machine in oVirt.
type VM interface {
	VMData

	// Update updates the virtual machine with the given parameters. Use UpdateVMParams to
	// get a builder for the parameters.
	Update(params UpdateVMParameters, retries ...RetryStrategy) (VM, error)
	// Remove removes the current VM. This involves an API call and may be slow.
	Remove(retries ...RetryStrategy) error

	// Start will cause a VM to start. The actual start process takes some time and should be checked via WaitForStatus.
	Start(retries ...RetryStrategy) error
	// Stop will cause the VM to power-off. The force parameter will cause the VM to stop even if a backup is currently
	// running.
	Stop(force bool, retries ...RetryStrategy) error
	// Shutdown will cause the VM to shut down. The force parameter will cause the VM to shut down even if a backup
	// is currently running.
	Shutdown(force bool, retries ...RetryStrategy) error
	// WaitForStatus will wait until the VM reaches the desired status. If the status is not reached within the
	// specified amount of retries, an error will be returned. If the VM enters the desired state, an updated VM
	// object will be returned.
	WaitForStatus(status VMStatus, retries ...RetryStrategy) (VM, error)

	// CreateNIC creates a network interface on the current VM. This involves an API call and may be slow.
	CreateNIC(name string, vnicProfileID VNICProfileID, params OptionalNICParameters, retries ...RetryStrategy) (NIC, error)
	// GetNIC fetches a NIC with a specific ID on the current VM. This involves an API call and may be slow.
	GetNIC(id NICID, retries ...RetryStrategy) (NIC, error)
	// ListNICs fetches a list of network interfaces attached to this VM. This involves an API call and may be slow.
	ListNICs(retries ...RetryStrategy) ([]NIC, error)

	// AttachDisk attaches a disk to this VM.
	AttachDisk(
		diskID DiskID,
		diskInterface DiskInterface,
		params CreateDiskAttachmentOptionalParams,
		retries ...RetryStrategy,
	) (DiskAttachment, error)
	// GetDiskAttachment returns a specific disk attachment for the current VM by ID.
	GetDiskAttachment(diskAttachmentID DiskAttachmentID, retries ...RetryStrategy) (DiskAttachment, error)
	// ListDiskAttachments lists all disk attachments for the current VM.
	ListDiskAttachments(retries ...RetryStrategy) ([]DiskAttachment, error)
	// DetachDisk removes a specific disk attachment by the disk attachment ID.
	DetachDisk(
		diskAttachmentID DiskAttachmentID,
		retries ...RetryStrategy,
	) error
	// Tags list all tags for the current VM
	Tags(retries ...RetryStrategy) ([]Tag, error)

	// GetHost retrieves the host object for the current VM. If the VM is not running, nil will be returned.
	GetHost(retries ...RetryStrategy) (Host, error)

	// GetIPAddresses fetches the IP addresses and returns a map of the interface name and list of IP addresses.
	//
	// The optional parameters let you filter the returned interfaces and IP addresses.
	GetIPAddresses(params VMIPSearchParams, retries ...RetryStrategy) (map[string][]net.IP, error)
	// GetNonLocalIPAddresses fetches the IP addresses, filters them for non-local IP addresses, and returns a map of the
	// interface name and list of IP addresses.
	GetNonLocalIPAddresses(retries ...RetryStrategy) (map[string][]net.IP, error)
	// WaitForIPAddresses waits for at least one IP address to be reported that is not in specified ranges.
	//
	// The returned result will be a map of network interface names and the list of IP addresses assigned to them,
	// excluding any IP addresses and interfaces in the specified parameters.
	WaitForIPAddresses(params VMIPSearchParams, retries ...RetryStrategy) (map[string][]net.IP, error)
	// WaitForNonLocalIPAddress waits for at least one IP address to be reported that is not in the following ranges:
	//
	// - 0.0.0.0/32
	// - 127.0.0.0/8
	// - 169.254.0.0/15
	// - 224.0.0.0/4
	// - 255.255.255.255/32
	// - ::/128
	// - ::1/128
	// - fe80::/64
	// - ff00::/8
	//
	// It also excludes the following interface names:
	//
	// - lo
	// - dummy*
	//
	// The returned result will be a map of network interface names and the list of non-local IP addresses assigned to
	// them.
	WaitForNonLocalIPAddress(retries ...RetryStrategy) (map[string][]net.IP, error)

	// AddTag adds the specified tag to the current VM.
	AddTag(tagID TagID, retries ...RetryStrategy) (err error)
	// RemoveTag removes the tag from the current VM.
	RemoveTag(tagID TagID, retries ...RetryStrategy) (err error)
	// ListTags lists the tags attached to the current VM.
	ListTags(retries ...RetryStrategy) (result []Tag, err error)

	// ListGraphicsConsoles lists the graphics consoles on the VM.
	ListGraphicsConsoles(retries ...RetryStrategy) ([]VMGraphicsConsole, error)

	// SerialConsole returns true if the VM has a serial console.
	SerialConsole() bool
}

// VMSearchParameters declares the parameters that can be passed to a VM search. Each parameter
// is declared as a pointer, where a nil value will mean that parameter will not be searched for.
// All parameters are used together as an AND filter.
type VMSearchParameters interface {
	// Name will match the name of the virtual machine exactly.
	Name() *string
	// Tag will match the tag of the virtual machine.
	Tag() *string
	// Statuses will return a list of acceptable statuses for this VM search.
	Statuses() *VMStatusList
	// NotStatuses will return a list of not acceptable statuses for this VM search.
	NotStatuses() *VMStatusList
}

// BuildableVMSearchParameters is a buildable version of VMSearchParameters.
type BuildableVMSearchParameters interface {
	VMSearchParameters

	// WithName sets the name to search for.
	WithName(name string) BuildableVMSearchParameters
	// WithTag sets the tag to search for.
	WithTag(name string) BuildableVMSearchParameters
	// WithStatus adds a single status to the filter.
	WithStatus(status VMStatus) BuildableVMSearchParameters
	// WithNotStatus excludes a VM status from the search.
	WithNotStatus(status VMStatus) BuildableVMSearchParameters
	// WithStatuses will return the statuses the returned VMs should be in.
	WithStatuses(list VMStatusList) BuildableVMSearchParameters
	// WithNotStatuses will return the statuses the returned VMs should not be in.
	WithNotStatuses(list VMStatusList) BuildableVMSearchParameters
}

// VMSearchParams creates a buildable set of search parameters for easier use.
func VMSearchParams() BuildableVMSearchParameters {
	return &vmSearchParams{
		lock: &sync.Mutex{},
	}
}

type vmSearchParams struct {
	lock *sync.Mutex

	name        *string
	tag         *string
	statuses    *VMStatusList
	notStatuses *VMStatusList
}

func (v *vmSearchParams) WithStatus(status VMStatus) BuildableVMSearchParameters {
	v.lock.Lock()
	defer v.lock.Unlock()
	newStatuses := append(*v.statuses, status)
	v.statuses = &newStatuses
	return v
}

func (v *vmSearchParams) WithNotStatus(status VMStatus) BuildableVMSearchParameters {
	v.lock.Lock()
	defer v.lock.Unlock()
	newNotStatuses := append(*v.notStatuses, status)
	v.statuses = &newNotStatuses
	return v
}

func (v *vmSearchParams) Tag() *string {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.tag
}

func (v *vmSearchParams) Name() *string {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.name
}

func (v *vmSearchParams) Statuses() *VMStatusList {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.statuses
}

func (v *vmSearchParams) NotStatuses() *VMStatusList {
	v.lock.Lock()
	defer v.lock.Unlock()
	return v.notStatuses
}

func (v *vmSearchParams) WithName(name string) BuildableVMSearchParameters {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.name = &name
	return v
}

func (v *vmSearchParams) WithTag(tag string) BuildableVMSearchParameters {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.tag = &tag
	return v
}

func (v *vmSearchParams) WithStatuses(list VMStatusList) BuildableVMSearchParameters {
	v.lock.Lock()
	defer v.lock.Unlock()
	newStatuses := list.Copy()
	v.statuses = &newStatuses
	return v
}

func (v *vmSearchParams) WithNotStatuses(list VMStatusList) BuildableVMSearchParameters {
	v.lock.Lock()
	defer v.lock.Unlock()
	newNotStatuses := list.Copy()
	v.notStatuses = &newNotStatuses
	return v
}

// OptionalVMParameters are a list of parameters that can be, but must not necessarily be added on VM creation. This
// interface is expected to be extended in the future.
type OptionalVMParameters interface {
	// Comment returns the comment for the VM.
	Comment() string

	// CPU contains the CPU topology, if any.
	CPU() VMCPUParams

	// HugePages returns the optional value for the HugePages setting for VMs.
	HugePages() *VMHugePages

	// Initialization defines the virtual machine’s initialization configuration.
	Initialization() Initialization

	// Clone should return true if the VM should be cloned from the template instead of linking it. This means that the
	// template can be removed while the VM still exists.
	Clone() *bool

	// Memory returns the VM memory in Bytes.
	Memory() *int64

	// MemoryPolicy returns the memory policy configuration for this VM, if any.
	MemoryPolicy() *MemoryPolicyParameters

	// Disks returns a list of disks that are to be changed from the template.
	Disks() []OptionalVMDiskParameters

	// PlacementPolicy returns a VM placement policy to apply, if any.
	PlacementPolicy() *VMPlacementPolicyParameters

	// InstanceTypeID returns the instance type ID if set.
	InstanceTypeID() *InstanceTypeID

	// VMType is the type of the VM created.
	VMType() *VMType

	// OS returns the operating system parameters, and true if the OS parameter has been set.
	OS() (VMOSParameters, bool)

	// SerialConsole returns if a serial console should be created or not.
	SerialConsole() *bool
}

// BuildableVMParameters is a variant of OptionalVMParameters that can be changed using the supplied
// builder functions. This is placed here for future use.
type BuildableVMParameters interface {
	OptionalVMParameters

	// WithComment adds a common to the VM.
	WithComment(comment string) (BuildableVMParameters, error)
	// MustWithComment is identical to WithComment, but panics instead of returning an error.
	MustWithComment(comment string) BuildableVMParameters

	// WithCPU adds a VMCPUTopo to the VM.
	WithCPU(cpu VMCPUParams) (BuildableVMParameters, error)
	// MustWithCPU adds a VMCPUTopo and panics if an error happens.
	MustWithCPU(cpu VMCPUParams) BuildableVMParameters
	// WithCPUParameters is a simplified function that calls NewVMCPUTopo and adds the CPU topology to
	// the VM.
	// Deprecated: use WithCPU instead.
	WithCPUParameters(cores, threads, sockets uint) (BuildableVMParameters, error)
	// MustWithCPUParameters is a simplified function that calls MustNewVMCPUTopo and adds the CPU topology to
	// the VM.
	// Deprecated: use MustWithCPU instead.
	MustWithCPUParameters(cores, threads, sockets uint) BuildableVMParameters

	// WithHugePages sets the HugePages setting for the VM.
	WithHugePages(hugePages VMHugePages) (BuildableVMParameters, error)
	// MustWithHugePages is identical to WithHugePages, but panics instead of returning an error.
	MustWithHugePages(hugePages VMHugePages) BuildableVMParameters
	// WithMemory sets the Memory setting for the VM.
	WithMemory(memory int64) (BuildableVMParameters, error)
	// MustWithMemory is identical to WithMemory, but panics instead of returning an error.
	MustWithMemory(memory int64) BuildableVMParameters
	// WithMemoryPolicy sets the memory policy parameters for the VM.
	WithMemoryPolicy(memory MemoryPolicyParameters) BuildableVMParameters

	// WithInitialization sets the virtual machine’s initialization configuration.
	WithInitialization(initialization Initialization) (BuildableVMParameters, error)
	// MustWithInitialization is identical to WithInitialization, but panics instead of returning an error.
	MustWithInitialization(initialization Initialization) BuildableVMParameters
	// MustWithInitializationParameters is a simplified function that calls MustNewInitialization and adds customScript
	MustWithInitializationParameters(customScript, hostname string) BuildableVMParameters

	// WithClone sets the clone flag. If the clone flag is true the VM is cloned from the template instead of linking to
	// it. This means the template can be deleted while the VM still exists.
	WithClone(clone bool) (BuildableVMParameters, error)
	// MustWithClone is identical to WithClone, but panics instead of returning an error.
	MustWithClone(clone bool) BuildableVMParameters

	// WithDisks adds disk configurations to the VM creation to manipulate the disks inherited from templates.
	WithDisks(disks []OptionalVMDiskParameters) (BuildableVMParameters, error)
	// MustWithDisks is identical to WithDisks, but panics instead of returning an error.
	MustWithDisks(disks []OptionalVMDiskParameters) BuildableVMParameters

	// WithPlacementPolicy adds a placement policy dictating which hosts the VM can be migrated to.
	WithPlacementPolicy(placementPolicy VMPlacementPolicyParameters) BuildableVMParameters

	// WithInstanceTypeID sets the instance type ID on the VM.
	WithInstanceTypeID(instanceTypeID InstanceTypeID) (BuildableVMParameters, error)
	// MustWithInstanceTypeID is identical to WithInstanceTypeID but panics instead of returning an error.
	MustWithInstanceTypeID(instanceTypeID InstanceTypeID) BuildableVMParameters

	// WithVMType sets the virtual machine type.
	WithVMType(vmType VMType) (BuildableVMParameters, error)
	// MustWithVMType is identical to WithVMType, but panics instead of returning an error.
	MustWithVMType(vmType VMType) BuildableVMParameters

	// WithOS adds the operating system parameters to the VM creation.
	WithOS(parameters VMOSParameters) BuildableVMParameters

	// WithSerialConsole adds or removes a serial console to the VM.
	WithSerialConsole(serialConsole bool) BuildableVMParameters
}

// VMCPUParams contain the CPU parameters for a VM.
type VMCPUParams interface {
	// Mode is the mode the CPU is used in. See CPUMode for details.
	Mode() *CPUMode
	// Topo contains the topology of the CPU.
	Topo() VMCPUTopoParams
}

// BuildableVMCPUParams is a buildable version of VMCPUParams.
type BuildableVMCPUParams interface {
	VMCPUParams

	WithMode(mode CPUMode) (BuildableVMCPUParams, error)
	MustWithMode(mode CPUMode) BuildableVMCPUParams

	WithTopo(topo VMCPUTopoParams) (BuildableVMCPUParams, error)
	MustWithTopo(topo VMCPUTopoParams) BuildableVMCPUParams
}

// NewVMCPUParams creates a new VMCPUParams object.
func NewVMCPUParams() BuildableVMCPUParams {
	return &vmCPUParams{}
}

type vmCPUParams struct {
	mode *CPUMode
	topo VMCPUTopoParams
}

func (v *vmCPUParams) WithMode(mode CPUMode) (BuildableVMCPUParams, error) {
	if err := mode.Validate(); err != nil {
		return nil, err
	}
	v.mode = &mode
	return v, nil
}

func (v *vmCPUParams) MustWithMode(mode CPUMode) BuildableVMCPUParams {
	builder, err := v.WithMode(mode)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmCPUParams) WithTopo(topo VMCPUTopoParams) (BuildableVMCPUParams, error) {
	v.topo = topo
	return v, nil
}

func (v *vmCPUParams) MustWithTopo(topo VMCPUTopoParams) BuildableVMCPUParams {
	b, err := v.WithTopo(topo)
	if err != nil {
		panic(err)
	}
	return b
}

func (v vmCPUParams) Mode() *CPUMode {
	return v.mode
}

func (v vmCPUParams) Topo() VMCPUTopoParams {
	return v.topo
}

// VMCPUTopoParams contain the CPU topology parameters for a VM.
type VMCPUTopoParams interface {
	// Sockets returns the number of sockets to be added to the VM. Must be at least 1.
	Sockets() uint
	// Threads returns the number of CPU threads to be added to the VM. Must be at least 1.
	Threads() uint
	// Cores returns the number of CPU cores to be added to the VM. Must be at least 1.
	Cores() uint
}

// BuildableVMCPUTopoParams is a buildable version of VMCPUTopoParams.
type BuildableVMCPUTopoParams interface {
	VMCPUTopoParams

	WithSockets(sockets uint) (BuildableVMCPUTopoParams, error)
	MustWithSockets(sockets uint) BuildableVMCPUTopoParams
	WithCores(cores uint) (BuildableVMCPUTopoParams, error)
	MustWithCores(cores uint) BuildableVMCPUTopoParams
	WithThreads(threads uint) (BuildableVMCPUTopoParams, error)
	MustWithThreads(threads uint) BuildableVMCPUTopoParams
}

// NewVMCPUTopoParams creates a new BuildableVMCPUTopoParams.
func NewVMCPUTopoParams() BuildableVMCPUTopoParams {
	return &vmCPUTopoParams{
		sockets: 1,
		cores:   1,
		threads: 1,
	}
}

type vmCPUTopoParams struct {
	sockets uint
	cores   uint
	threads uint
}

func (v vmCPUTopoParams) Sockets() uint {
	return v.sockets
}

func (v vmCPUTopoParams) Threads() uint {
	return v.threads
}

func (v vmCPUTopoParams) Cores() uint {
	return v.cores
}

func (v *vmCPUTopoParams) WithSockets(sockets uint) (BuildableVMCPUTopoParams, error) {
	if sockets == 0 {
		return nil, newError(EBadArgument, "sockets must be at least 1")
	}
	v.sockets = sockets
	return v, nil
}

func (v *vmCPUTopoParams) MustWithSockets(sockets uint) BuildableVMCPUTopoParams {
	builder, err := v.WithSockets(sockets)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmCPUTopoParams) WithCores(cores uint) (BuildableVMCPUTopoParams, error) {
	if cores == 0 {
		return nil, newError(EBadArgument, "cores must be at least 1")
	}
	v.cores = cores
	return v, nil
}

func (v *vmCPUTopoParams) MustWithCores(cores uint) BuildableVMCPUTopoParams {
	builder, err := v.WithCores(cores)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmCPUTopoParams) WithThreads(threads uint) (BuildableVMCPUTopoParams, error) {
	if threads == 0 {
		return nil, newError(EBadArgument, "threads must be at least 1")
	}
	v.threads = threads
	return v, nil
}

func (v *vmCPUTopoParams) MustWithThreads(threads uint) BuildableVMCPUTopoParams {
	builder, err := v.WithThreads(threads)
	if err != nil {
		panic(err)
	}
	return builder
}

// VMOSParameters contains the VM parameters pertaining to the operating system.
type VMOSParameters interface {
	// Type returns the type-string for the operating system.
	Type() *string
}

// BuildableVMOSParameters is a buildable version of VMOSParameters.
type BuildableVMOSParameters interface {
	VMOSParameters

	WithType(t string) (BuildableVMOSParameters, error)
	MustWithType(t string) BuildableVMOSParameters
}

// NewVMOSParameters creates a new VMOSParameters structure.
func NewVMOSParameters() BuildableVMOSParameters {
	return &vmOSParameters{}
}

type vmOSParameters struct {
	t *string
}

func (v *vmOSParameters) Type() *string {
	return v.t
}

func (v *vmOSParameters) WithType(t string) (BuildableVMOSParameters, error) {
	v.t = &t
	return v, nil
}

func (v *vmOSParameters) MustWithType(t string) BuildableVMOSParameters {
	builder, err := v.WithType(t)
	if err != nil {
		panic(err)
	}
	return builder
}

// CPUMode is the mode of the CPU on a VM.
type CPUMode string

const (
	// CPUModeCustom contains custom settings for the CPU.
	CPUModeCustom CPUMode = "custom"
	// CPUModeHostModel copies the host CPU make and model.
	CPUModeHostModel CPUMode = "host_model"
	// CPUModeHostPassthrough passes through the host CPU for nested virtualization.
	CPUModeHostPassthrough CPUMode = "host_passthrough"
)

// Validate checks if the CPU mode is valid.
func (c CPUMode) Validate() error {
	switch c {
	case CPUModeCustom:
		return nil
	case CPUModeHostModel:
		return nil
	case CPUModeHostPassthrough:
		return nil
	default:
		return newError(EBadArgument, "invalid CPU mode: %s", c)
	}
}

// CPUModeValues lists all valid CPU modes.
func CPUModeValues() []CPUMode {
	return []CPUMode{
		CPUModeCustom,
		CPUModeHostModel,
		CPUModeHostPassthrough,
	}
}

// VMType contains some preconfigured settings, such as the availability of remote desktop, for the VM.
type VMType string

const (
	// VMTypeDesktop indicates that the virtual machine is intended to be used as a desktop and enables the SPICE
	// virtual console, and a sound device will be added to the VM.
	VMTypeDesktop VMType = "desktop"
	// VMTypeServer indicates that the VM will be used as a server. In this case a sound device will not be added to the
	// VM.
	VMTypeServer VMType = "server"
	// VMTypeHighPerformance indicates that the VM should be configured for high performance. This entails the following
	// options:
	//
	// - Enable headless mode.
	// - Enable pass-through host CPU.
	// - Enable I/O threads.
	// - Enable I/O threads pinning and set the pinning topology.
	// - Enable the paravirtualized random number generator PCI (virtio-rng) device.
	// - Disable all USB devices.
	// - Disable the soundcard device.
	// - Disable the smartcard device.
	// - Disable the memory balloon device.
	// - Disable the watchdog device.
	// - Disable migration.
	// - Disable high availability.
	VMTypeHighPerformance VMType = "high_performance"
)

// Validate checks if the VMType value is valid.
func (v VMType) Validate() error {
	switch v {
	case VMTypeDesktop:
		return nil
	case VMTypeServer:
		return nil
	case VMTypeHighPerformance:
		return nil
	default:
		return newError(EBadArgument, "Invalid VM type: %s", v)
	}
}

// VMTypeValues returns all possible values for VM types.
func VMTypeValues() []VMType {
	return []VMType{
		VMTypeDesktop,
		VMTypeServer,
		VMTypeHighPerformance,
	}
}

// VMPlacementPolicyParameters contains the optional parameters on VM placement.
type VMPlacementPolicyParameters interface {
	// Affinity dictates how a VM can be migrated to a different host. This can be nil, in which case the engine
	// default is to set the policy to migratable.
	Affinity() *VMAffinity
	// HostIDs returns a list of host IDs to apply as possible migration targets. The default is an empty list,
	// which means the VM can be migrated to any host.
	HostIDs() []HostID
}

// BuildableVMPlacementPolicyParameters is a buildable version of the VMPlacementPolicyParameters.
type BuildableVMPlacementPolicyParameters interface {
	VMPlacementPolicyParameters

	// WithAffinity sets the way VMs can be migrated to other hosts.
	WithAffinity(affinity VMAffinity) (BuildableVMPlacementPolicyParameters, error)
	// MustWithAffinity is identical to WithAffinity, but panics instead of returning an error.
	MustWithAffinity(affinity VMAffinity) BuildableVMPlacementPolicyParameters

	// WithHostIDs sets the list of hosts this VM can be migrated to.
	WithHostIDs(hostIDs []HostID) (BuildableVMPlacementPolicyParameters, error)
	// MustWithHostIDs is identical to WithHostIDs, but panics instead of returning an error.
	MustWithHostIDs(hostIDs []HostID) BuildableVMPlacementPolicyParameters
}

// NewVMPlacementPolicyParameters creates a new BuildableVMPlacementPolicyParameters for use on VM creation.
func NewVMPlacementPolicyParameters() BuildableVMPlacementPolicyParameters {
	return &vmPlacementPolicyParameters{}
}

type vmPlacementPolicyParameters struct {
	affinity *VMAffinity
	hostIDs  []HostID
}

func (v vmPlacementPolicyParameters) Affinity() *VMAffinity {
	return v.affinity
}

func (v vmPlacementPolicyParameters) HostIDs() []HostID {
	return v.hostIDs
}

func (v vmPlacementPolicyParameters) WithAffinity(affinity VMAffinity) (BuildableVMPlacementPolicyParameters, error) {
	if err := affinity.Validate(); err != nil {
		return nil, err
	}
	v.affinity = &affinity
	return v, nil
}

func (v vmPlacementPolicyParameters) MustWithAffinity(affinity VMAffinity) BuildableVMPlacementPolicyParameters {
	builder, err := v.WithAffinity(affinity)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v vmPlacementPolicyParameters) WithHostIDs(hostIDs []HostID) (BuildableVMPlacementPolicyParameters, error) {
	v.hostIDs = hostIDs
	return v, nil
}

func (v vmPlacementPolicyParameters) MustWithHostIDs(hostIDs []HostID) BuildableVMPlacementPolicyParameters {
	builder, err := v.WithHostIDs(hostIDs)
	if err != nil {
		panic(err)
	}
	return builder
}

// MemoryPolicyParameters contain the parameters for the memory policy setting on the VM.
type MemoryPolicyParameters interface {
	Guaranteed() *int64
	Max() *int64
	Ballooning() *bool
}

// BuildableMemoryPolicyParameters is a buildable version of MemoryPolicyParameters.
type BuildableMemoryPolicyParameters interface {
	MemoryPolicyParameters

	WithGuaranteed(guaranteed int64) (BuildableMemoryPolicyParameters, error)
	MustWithGuaranteed(guaranteed int64) BuildableMemoryPolicyParameters

	WithMax(max int64) (BuildableMemoryPolicyParameters, error)
	MustWithMax(max int64) BuildableMemoryPolicyParameters

	WithBallooning(ballooning bool) (BuildableMemoryPolicyParameters, error)
	MustWithBallooning(ballooning bool) BuildableMemoryPolicyParameters
}

// NewMemoryPolicyParameters creates a new instance of BuildableMemoryPolicyParameters.
func NewMemoryPolicyParameters() BuildableMemoryPolicyParameters {
	return &memoryPolicyParameters{}
}

type memoryPolicyParameters struct {
	guaranteed *int64
	max        *int64
	ballooning *bool
}

func (m *memoryPolicyParameters) Ballooning() *bool {
	return m.ballooning
}

func (m *memoryPolicyParameters) WithBallooning(ballooning bool) (BuildableMemoryPolicyParameters, error) {
	m.ballooning = &ballooning
	return m, nil
}

func (m *memoryPolicyParameters) MustWithBallooning(ballooning bool) BuildableMemoryPolicyParameters {
	builder, err := m.WithBallooning(ballooning)
	if err != nil {
		panic(err)
	}
	return builder
}

func (m *memoryPolicyParameters) MustWithGuaranteed(guaranteed int64) BuildableMemoryPolicyParameters {
	builder, err := m.WithGuaranteed(guaranteed)
	if err != nil {
		panic(err)
	}
	return builder
}

func (m *memoryPolicyParameters) Guaranteed() *int64 {
	return m.guaranteed
}

func (m *memoryPolicyParameters) WithGuaranteed(guaranteed int64) (BuildableMemoryPolicyParameters, error) {
	m.guaranteed = &guaranteed
	return m, nil
}

func (m *memoryPolicyParameters) MustWithMax(max int64) BuildableMemoryPolicyParameters {
	builder, err := m.WithMax(max)
	if err != nil {
		panic(err)
	}
	return builder
}

func (m *memoryPolicyParameters) Max() *int64 {
	return m.max
}

func (m *memoryPolicyParameters) WithMax(max int64) (BuildableMemoryPolicyParameters, error) {
	m.max = &max
	return m, nil
}

// MemoryPolicy is the memory policy set on the VM.
type MemoryPolicy interface {
	// Guaranteed returns the number of guaranteed bytes to the VM.
	Guaranteed() *int64
	// Max returns the maximum amount of memory given to the VM.
	Max() *int64
	// Ballooning returns true if the VM can give back the memory it is not using to the host OS.
	Ballooning() bool
}

type memoryPolicy struct {
	guaranteed *int64
	max        *int64
	ballooning bool
}

func (m memoryPolicy) Ballooning() bool {
	return m.ballooning
}

func (m memoryPolicy) Max() *int64 {
	return m.max
}

func (m memoryPolicy) Guaranteed() *int64 {
	return m.guaranteed
}

// OptionalVMDiskParameters describes the disk parameters that can be given to VM creation. These manipulate the
// disks inherited from the template.
type OptionalVMDiskParameters interface {
	// DiskID returns the identifier of the disk that is being changed.
	DiskID() DiskID
	// Sparse sets the sparse parameter if set. Note, that Sparse is only supported in oVirt on block devices with QCOW2
	// images. On NFS you MUST use raw disks to use sparse.
	Sparse() *bool
	// Format returns the image format to be used for the specified disk.
	Format() *ImageFormat
	// StorageDomainID returns the optional storage domain ID to use for this disk.
	StorageDomainID() *StorageDomainID
}

// BuildableVMDiskParameters is a buildable version of OptionalVMDiskParameters.
type BuildableVMDiskParameters interface {
	OptionalVMDiskParameters

	// WithSparse enables or disables sparse disk provisioning. Note, that Sparse is only supported in oVirt on block
	// devices with QCOW2 images. On NFS you MUST use raw images to use sparse. See WithFormat.
	WithSparse(sparse bool) (BuildableVMDiskParameters, error)
	// MustWithSparse is identical to WithSparse, but panics instead of returning an error.
	MustWithSparse(sparse bool) BuildableVMDiskParameters

	// WithFormat adds a disk format to the VM on creation. Note, that QCOW2 is only supported in conjunction with
	// Sparse on block devices. On NFS you MUST use raw images to use sparse. See WithSparse.
	WithFormat(format ImageFormat) (BuildableVMDiskParameters, error)
	// MustWithFormat is identical to WithFormat, but panics instead of returning an error.
	MustWithFormat(format ImageFormat) BuildableVMDiskParameters

	// WithStorageDomainID adds a storage domain to use for the disk.
	WithStorageDomainID(storageDomainID StorageDomainID) (BuildableVMDiskParameters, error)
	// MustWithStorageDomainID is identical to WithStorageDomainID but panics instead of returning an error.
	MustWithStorageDomainID(storageDomainID StorageDomainID) BuildableVMDiskParameters
}

// NewBuildableVMDiskParameters creates a new buildable OptionalVMDiskParameters.
func NewBuildableVMDiskParameters(diskID DiskID) (BuildableVMDiskParameters, error) {
	return &vmDiskParameters{
		diskID,
		nil,
		nil,
		nil,
	}, nil
}

// MustNewBuildableVMDiskParameters is identical to NewBuildableVMDiskParameters but panics instead of returning an
// error.
func MustNewBuildableVMDiskParameters(diskID DiskID) BuildableVMDiskParameters {
	builder, err := NewBuildableVMDiskParameters(diskID)
	if err != nil {
		panic(err)
	}
	return builder
}

type vmDiskParameters struct {
	diskID          DiskID
	sparse          *bool
	format          *ImageFormat
	storageDomainID *StorageDomainID
}

func (v *vmDiskParameters) StorageDomainID() *StorageDomainID {
	return v.storageDomainID
}

func (v *vmDiskParameters) WithStorageDomainID(storageDomainID StorageDomainID) (BuildableVMDiskParameters, error) {
	v.storageDomainID = &storageDomainID
	return v, nil
}

func (v *vmDiskParameters) MustWithStorageDomainID(storageDomainID StorageDomainID) BuildableVMDiskParameters {
	b, err := v.WithStorageDomainID(storageDomainID)
	if err != nil {
		panic(err)
	}
	return b
}

func (v *vmDiskParameters) Format() *ImageFormat {
	return v.format
}

func (v *vmDiskParameters) WithFormat(format ImageFormat) (BuildableVMDiskParameters, error) {
	if err := format.Validate(); err != nil {
		return nil, err
	}
	v.format = &format
	return v, nil
}

func (v *vmDiskParameters) MustWithFormat(format ImageFormat) BuildableVMDiskParameters {
	builder, err := v.WithFormat(format)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmDiskParameters) DiskID() DiskID {
	return v.diskID
}

func (v *vmDiskParameters) Sparse() *bool {
	return v.sparse
}

func (v *vmDiskParameters) WithSparse(sparse bool) (BuildableVMDiskParameters, error) {
	v.sparse = &sparse
	return v, nil
}

func (v *vmDiskParameters) MustWithSparse(sparse bool) BuildableVMDiskParameters {
	builder, err := v.WithSparse(sparse)
	if err != nil {
		panic(err)
	}
	return builder
}

// UpdateVMParameters returns a set of parameters to change on a VM.
type UpdateVMParameters interface {
	// Name returns the name for the VM. Return nil if the name should not be changed.
	Name() *string
	// Comment returns the comment for the VM. Return nil if the name should not be changed.
	Comment() *string
}

// VMCPUTopo contains the CPU topology information about a VM.
type VMCPUTopo interface {
	// Cores is the number of CPU cores.
	Cores() uint
	// Threads is the number of CPU threads in a core.
	Threads() uint
	// Sockets is the number of sockets.
	Sockets() uint
}

// NewVMCPUTopo creates a new VMCPUTopo from the specified parameters.
func NewVMCPUTopo(cores uint, threads uint, sockets uint) (VMCPUTopo, error) {
	if cores == 0 {
		return nil, newError(EBadArgument, "number of cores must be positive")
	}
	if threads == 0 {
		return nil, newError(EBadArgument, "number of threads must be positive")
	}
	if sockets == 0 {
		return nil, newError(EBadArgument, "number of sockets must be positive")
	}
	return &vmCPUTopo{
		cores:   cores,
		threads: threads,
		sockets: sockets,
	}, nil
}

// MustNewVMCPUTopo is equivalent to NewVMCPUTopo, but panics instead of returning an error.
func MustNewVMCPUTopo(cores uint, threads uint, sockets uint) VMCPUTopo {
	topo, err := NewVMCPUTopo(cores, threads, sockets)
	if err != nil {
		panic(err)
	}
	return topo
}

type vmCPUTopo struct {
	cores   uint
	threads uint
	sockets uint
}

func (v *vmCPUTopo) Cores() uint {
	return v.cores
}

func (v *vmCPUTopo) Threads() uint {
	return v.threads
}

func (v *vmCPUTopo) Sockets() uint {
	return v.sockets
}

func (v *vmCPUTopo) clone() *vmCPUTopo {
	if v == nil {
		return nil
	}
	return &vmCPUTopo{
		cores:   v.cores,
		threads: v.threads,
		sockets: v.sockets,
	}
}

// BuildableUpdateVMParameters is a buildable version of UpdateVMParameters.
type BuildableUpdateVMParameters interface {
	UpdateVMParameters

	// WithName adds an updated name to the request.
	WithName(name string) (BuildableUpdateVMParameters, error)

	// MustWithName is identical to WithName, but panics instead of returning an error
	MustWithName(name string) BuildableUpdateVMParameters

	// WithComment adds a comment to the request
	WithComment(comment string) (BuildableUpdateVMParameters, error)

	// MustWithComment is identical to WithComment, but panics instead of returning an error.
	MustWithComment(comment string) BuildableUpdateVMParameters
}

// UpdateVMParams returns a buildable set of update parameters.
func UpdateVMParams() BuildableUpdateVMParameters {
	return &updateVMParams{}
}

type updateVMParams struct {
	name    *string
	comment *string
}

func (u *updateVMParams) MustWithName(name string) BuildableUpdateVMParameters {
	builder, err := u.WithName(name)
	if err != nil {
		panic(err)
	}
	return builder
}

func (u *updateVMParams) MustWithComment(comment string) BuildableUpdateVMParameters {
	builder, err := u.WithComment(comment)
	if err != nil {
		panic(err)
	}
	return builder
}

func (u *updateVMParams) Name() *string {
	return u.name
}

func (u *updateVMParams) Comment() *string {
	return u.comment
}

func (u *updateVMParams) WithName(name string) (BuildableUpdateVMParameters, error) {
	if err := validateVMName(name); err != nil {
		return nil, err
	}
	u.name = &name
	return u, nil
}

func (u *updateVMParams) WithComment(comment string) (BuildableUpdateVMParameters, error) {
	u.comment = &comment
	return u, nil
}

// NewCreateVMParams creates a set of BuildableVMParameters that can be used to construct the optional VM parameters.
func NewCreateVMParams() BuildableVMParameters {
	return &vmParams{
		lock: &sync.Mutex{},
	}
}

// CreateVMParams creates a set of BuildableVMParameters that can be used to construct the optional VM parameters.
// Deprecated: use NewCreateVMParams instead.
func CreateVMParams() BuildableVMParameters {
	return NewCreateVMParams()
}

type vmParams struct {
	lock *sync.Mutex

	name    string
	comment string
	cpu     VMCPUParams

	hugePages *VMHugePages

	initialization Initialization
	memory         *int64
	memoryPolicy   *MemoryPolicyParameters

	clone *bool

	disks []OptionalVMDiskParameters

	placementPolicy *VMPlacementPolicyParameters

	instanceTypeID *InstanceTypeID

	vmType *VMType

	os    VMOSParameters
	osSet bool

	serialConsole *bool
}

func (v *vmParams) SerialConsole() *bool {
	return v.serialConsole
}

func (v *vmParams) WithSerialConsole(serialConsole bool) BuildableVMParameters {
	v.serialConsole = &serialConsole
	return v
}

func (v *vmParams) OS() (VMOSParameters, bool) {
	return v.os, v.osSet
}

func (v *vmParams) WithOS(os VMOSParameters) BuildableVMParameters {
	v.os = os
	v.osSet = true
	return v
}

func (v *vmParams) VMType() *VMType {
	return v.vmType
}

func (v *vmParams) WithVMType(vmType VMType) (BuildableVMParameters, error) {
	if err := vmType.Validate(); err != nil {
		return nil, err
	}
	v.vmType = &vmType
	return v, nil
}

func (v *vmParams) MustWithVMType(vmType VMType) BuildableVMParameters {
	builder, err := v.WithVMType(vmType)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) InstanceTypeID() *InstanceTypeID {
	return v.instanceTypeID
}

func (v *vmParams) WithInstanceTypeID(instanceTypeID InstanceTypeID) (BuildableVMParameters, error) {
	v.instanceTypeID = &instanceTypeID
	return v, nil
}

func (v *vmParams) MustWithInstanceTypeID(instanceTypeID InstanceTypeID) BuildableVMParameters {
	builder, err := v.WithInstanceTypeID(instanceTypeID)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) WithPlacementPolicy(placementPolicy VMPlacementPolicyParameters) BuildableVMParameters {
	v.placementPolicy = &placementPolicy
	return v
}

func (v *vmParams) PlacementPolicy() *VMPlacementPolicyParameters {
	return v.placementPolicy
}

func (v *vmParams) MemoryPolicy() *MemoryPolicyParameters {
	return v.memoryPolicy
}

func (v *vmParams) WithMemoryPolicy(memory MemoryPolicyParameters) BuildableVMParameters {
	v.memoryPolicy = &memory
	return v
}

func (v *vmParams) Clone() *bool {
	return v.clone
}

func (v *vmParams) WithClone(clone bool) (BuildableVMParameters, error) {
	v.clone = &clone
	return v, nil
}

func (v *vmParams) MustWithClone(clone bool) BuildableVMParameters {
	builder, err := v.WithClone(clone)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) Disks() []OptionalVMDiskParameters {
	return v.disks
}

func (v *vmParams) WithDisks(disks []OptionalVMDiskParameters) (BuildableVMParameters, error) {
	diskIDs := map[DiskID]int{}
	for i, d := range disks {
		if previousID, ok := diskIDs[d.DiskID()]; ok {
			return nil, newError(
				EBadArgument,
				"Disk %s appears twice, in position %d and %d.",
				d.DiskID(),
				previousID,
				i,
			)
		}
	}
	v.disks = disks
	return v, nil
}

func (v *vmParams) MustWithDisks(disks []OptionalVMDiskParameters) BuildableVMParameters {
	builder, err := v.WithDisks(disks)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) HugePages() *VMHugePages {
	return v.hugePages
}

func (v *vmParams) WithHugePages(hugePages VMHugePages) (BuildableVMParameters, error) {
	if err := hugePages.Validate(); err != nil {
		return v, err
	}
	v.hugePages = &hugePages
	return v, nil
}

func (v *vmParams) MustWithHugePages(hugePages VMHugePages) BuildableVMParameters {
	builder, err := v.WithHugePages(hugePages)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) Memory() *int64 {
	return v.memory
}

func (v *vmParams) WithMemory(memory int64) (BuildableVMParameters, error) {
	v.memory = &memory
	return v, nil
}

func (v *vmParams) MustWithMemory(memory int64) BuildableVMParameters {
	builder, err := v.WithMemory(memory)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) Initialization() Initialization {
	return v.initialization
}

func (v *vmParams) WithInitialization(initialization Initialization) (BuildableVMParameters, error) {
	v.initialization = initialization
	return v, nil
}

func (v *vmParams) MustWithInitialization(initialization Initialization) BuildableVMParameters {
	builder, err := v.WithInitialization(initialization)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) MustWithInitializationParameters(customScript, hostname string) BuildableVMParameters {
	init := NewInitialization(customScript, hostname)
	return v.MustWithInitialization(init)
}

func (v *vmParams) CPU() VMCPUParams {
	return v.cpu
}

func (v *vmParams) WithCPU(cpu VMCPUParams) (BuildableVMParameters, error) {
	v.cpu = cpu
	return v, nil
}

func (v *vmParams) MustWithCPU(cpu VMCPUParams) BuildableVMParameters {
	builder, err := v.WithCPU(cpu)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) WithCPUParameters(cores, threads, sockets uint) (BuildableVMParameters, error) {
	params := NewVMCPUTopoParams()
	params, err := params.WithCores(cores)
	if err != nil {
		return nil, err
	}
	params, err = params.WithThreads(threads)
	if err != nil {
		return nil, err
	}
	params, err = params.WithSockets(sockets)
	if err != nil {
		return nil, err
	}

	topo, err := NewVMCPUParams().WithTopo(params)
	if err != nil {
		return nil, err
	}

	return v.WithCPU(topo)
}

func (v *vmParams) MustWithCPUParameters(cores, threads, sockets uint) BuildableVMParameters {
	b, err := v.WithCPUParameters(cores, threads, sockets)
	if err != nil {
		panic(err)
	}
	return b
}

func (v *vmParams) MustWithName(name string) BuildableVMParameters {
	builder, err := v.WithName(name)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) MustWithComment(comment string) BuildableVMParameters {
	builder, err := v.WithComment(comment)
	if err != nil {
		panic(err)
	}
	return builder
}

func (v *vmParams) WithName(name string) (BuildableVMParameters, error) {
	if err := validateVMName(name); err != nil {
		return nil, err
	}
	v.name = name
	return v, nil
}

func (v *vmParams) WithComment(comment string) (BuildableVMParameters, error) {
	v.comment = comment
	return v, nil
}

func (v vmParams) Name() string {
	return v.name
}

func (v vmParams) Comment() string {
	return v.comment
}

type vm struct {
	client Client

	id              VMID
	name            string
	comment         string
	clusterID       ClusterID
	templateID      TemplateID
	status          VMStatus
	cpu             *vmCPU
	memory          int64
	tagIDs          []TagID
	hugePages       *VMHugePages
	initialization  Initialization
	hostID          *HostID
	placementPolicy *vmPlacementPolicy
	memoryPolicy    *memoryPolicy
	instanceTypeID  *InstanceTypeID
	vmType          VMType
	os              *vmOS
	serialConsole   bool
}

func (v *vm) SerialConsole() bool {
	return v.serialConsole
}

func (v *vm) ListGraphicsConsoles(retries ...RetryStrategy) ([]VMGraphicsConsole, error) {
	return v.client.ListVMGraphicsConsoles(v.id, retries...)
}

func (v *vm) OS() VMOS {
	return v.os
}

func (v *vm) VMType() VMType {
	return v.vmType
}

func (v *vm) InstanceTypeID() *InstanceTypeID {
	return v.instanceTypeID
}

func (v *vm) AddTag(tagID TagID, retries ...RetryStrategy) (err error) {
	return v.client.AddTagToVM(v.id, tagID, retries...)
}

func (v *vm) RemoveTag(tagID TagID, retries ...RetryStrategy) (err error) {
	return v.client.RemoveTagFromVM(v.id, tagID, retries...)
}

func (v *vm) ListTags(retries ...RetryStrategy) (result []Tag, err error) {
	return v.client.ListVMTags(v.id, retries...)
}

func (v *vm) PlacementPolicy() (VMPlacementPolicy, bool) {
	return v.placementPolicy, v.placementPolicy != nil
}

func (v *vm) MemoryPolicy() MemoryPolicy {
	return v.memoryPolicy
}

func (v *vm) WaitForIPAddresses(params VMIPSearchParams, retries ...RetryStrategy) (map[string][]net.IP, error) {
	return v.client.WaitForVMIPAddresses(v.id, params, retries...)
}

func (v *vm) WaitForNonLocalIPAddress(retries ...RetryStrategy) (map[string][]net.IP, error) {
	return v.client.WaitForNonLocalVMIPAddress(v.id, retries...)
}

func (v *vm) GetIPAddresses(params VMIPSearchParams, retries ...RetryStrategy) (map[string][]net.IP, error) {
	return v.client.GetVMIPAddresses(v.id, params, retries...)
}

func (v *vm) GetNonLocalIPAddresses(retries ...RetryStrategy) (map[string][]net.IP, error) {
	return v.client.GetVMNonLocalIPAddresses(v.id, retries...)
}

func (v *vm) HostID() *HostID {
	return v.hostID
}

func (v *vm) GetHost(retries ...RetryStrategy) (Host, error) {
	hostID := v.hostID
	if hostID == nil {
		return nil, nil
	}
	return v.client.GetHost(*hostID, retries...)
}

func (v *vm) HugePages() *VMHugePages {
	return v.hugePages
}

func (v *vm) Start(retries ...RetryStrategy) error {
	return v.client.StartVM(v.id, retries...)
}

func (v *vm) Stop(force bool, retries ...RetryStrategy) error {
	return v.client.StopVM(v.id, force, retries...)
}

func (v *vm) Shutdown(force bool, retries ...RetryStrategy) error {
	return v.client.ShutdownVM(v.id, force, retries...)
}

func (v *vm) WaitForStatus(status VMStatus, retries ...RetryStrategy) (VM, error) {
	return v.client.WaitForVMStatus(v.id, status, retries...)
}

func (v *vm) CPU() VMCPU {
	return v.cpu
}

func (v *vm) Memory() int64 {
	return v.memory
}

func (v *vm) Initialization() Initialization {
	return v.initialization
}

// withName returns a copy of the VM with the new name. It does not change the original copy to avoid
// shared state issues.
func (v *vm) withName(name string) *vm {
	return &vm{
		v.client,
		v.id,
		name,
		v.comment,
		v.clusterID,
		v.templateID,
		v.status,
		v.cpu,
		v.memory,
		v.tagIDs,
		v.hugePages,
		v.initialization,
		v.hostID,
		v.placementPolicy,
		v.memoryPolicy,
		v.instanceTypeID,
		v.vmType,
		v.os,
		v.serialConsole,
	}
}

// withComment returns a copy of the VM with the new comment. It does not change the original copy to avoid
// shared state issues.
func (v *vm) withComment(comment string) *vm {
	return &vm{
		v.client,
		v.id,
		v.name,
		comment,
		v.clusterID,
		v.templateID,
		v.status,
		v.cpu,
		v.memory,
		v.tagIDs,
		v.hugePages,
		v.initialization,
		v.hostID,
		v.placementPolicy,
		v.memoryPolicy,
		v.instanceTypeID,
		v.vmType,
		v.os,
		v.serialConsole,
	}
}

func (v *vm) Update(params UpdateVMParameters, retries ...RetryStrategy) (VM, error) {
	return v.client.UpdateVM(v.id, params, retries...)
}

func (v *vm) Status() VMStatus {
	return v.status
}

func (v *vm) AttachDisk(
	diskID DiskID,
	diskInterface DiskInterface,
	params CreateDiskAttachmentOptionalParams,
	retries ...RetryStrategy,
) (DiskAttachment, error) {
	return v.client.CreateDiskAttachment(v.id, diskID, diskInterface, params, retries...)
}

func (v *vm) GetDiskAttachment(diskAttachmentID DiskAttachmentID, retries ...RetryStrategy) (DiskAttachment, error) {
	return v.client.GetDiskAttachment(v.id, diskAttachmentID, retries...)
}

func (v *vm) ListDiskAttachments(retries ...RetryStrategy) ([]DiskAttachment, error) {
	return v.client.ListDiskAttachments(v.id, retries...)
}

func (v *vm) DetachDisk(diskAttachmentID DiskAttachmentID, retries ...RetryStrategy) error {
	return v.client.RemoveDiskAttachment(v.id, diskAttachmentID, retries...)
}

func (v *vm) Remove(retries ...RetryStrategy) error {
	return v.client.RemoveVM(v.id, retries...)
}

func (v *vm) CreateNIC(name string, vnicProfileID VNICProfileID, params OptionalNICParameters, retries ...RetryStrategy) (
	NIC,
	error,
) {
	return v.client.CreateNIC(v.id, vnicProfileID, name, params, retries...)
}

func (v *vm) GetNIC(id NICID, retries ...RetryStrategy) (NIC, error) {
	return v.client.GetNIC(v.id, id, retries...)
}

func (v *vm) ListNICs(retries ...RetryStrategy) ([]NIC, error) {
	return v.client.ListNICs(v.id, retries...)
}

func (v *vm) Comment() string {
	return v.comment
}

func (v *vm) ClusterID() ClusterID {
	return v.clusterID
}

func (v *vm) TemplateID() TemplateID {
	return v.templateID
}

func (v *vm) ID() VMID {
	return v.id
}

func (v *vm) Name() string {
	return v.name
}

func (v *vm) TagIDs() []TagID {
	return v.tagIDs
}

func (v *vm) Tags(retries ...RetryStrategy) ([]Tag, error) {
	tags := make([]Tag, len(v.tagIDs))
	for i, id := range v.tagIDs {
		tag, err := v.client.GetTag(id, retries...)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}
	return tags, nil
}

func (v *vm) AddTagToVM(tagID TagID, retries ...RetryStrategy) error {
	return v.client.AddTagToVM(v.id, tagID, retries...)
}
func (v *vm) AddTagToVMByName(tagName string, retries ...RetryStrategy) error {
	return v.client.AddTagToVMByName(v.id, tagName, retries...)
}

var vmNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_\-.]*$`)

func validateVMName(name string) error {
	if !vmNameRegexp.MatchString(name) {
		return newError(EBadArgument, "invalid VM name: %s", name)
	}
	return nil
}

func convertSDKVM(sdkObject *ovirtsdk.Vm, client Client, logger Logger, action string) (VM, error) {
	vmObject := &vm{
		client: client,
	}
	vmConverters := []func(sdkObject *ovirtsdk.Vm, vm *vm) error{
		vmIDConverter,
		vmNameConverter,
		vmCommentConverter,
		vmClusterConverter,
		vmStatusConverter,
		vmTemplateConverter,
		vmCPUConverter,
		vmHugePagesConverter,
		vmTagsConverter,
		vmInitializationConverter,
		vmPlacementPolicyConverter,
		vmHostConverter,
		vmMemoryConverter,
		vmMemoryPolicyConverter,
		vmInstanceTypeIDConverter,
		vmTypeConverter,
		vmOSConverter,
	}
	for _, converter := range vmConverters {
		if err := converter(sdkObject, vmObject); err != nil {
			return nil, err
		}
	}

	if err := vmSerialConsoleConverter(sdkObject, vmObject, logger, action); err != nil {
		return nil, err
	}

	return vmObject, nil
}

func vmSerialConsoleConverter(object *ovirtsdk.Vm, v *vm, logger Logger, action string) error {
	console, ok := object.Console()
	if !ok {
		// This sometimes (?) happens, and we don't know why. We will assume the serial console does not exist,
		// otherwise we create a flaky behavior.
		//
		// See these issues:
		// - https://github.com/oVirt/terraform-provider-ovirt/issues/411
		// - https://github.com/oVirt/go-ovirt-client/issues/211
		logger.Warningf(
			"If you see this message, please open a Bugzilla entry with your oVirt version! The virtual machine object was returned from the oVirt Engine while %s without a console sub-object. VM: %s status: %s",
			action,
			object.MustId(),
			object.MustStatus(),
		)
		v.serialConsole = false
		return nil
	}
	enabled, ok := console.Enabled()
	if !ok {
		return newFieldNotFound("serial console", "enabled")
	}
	v.serialConsole = enabled
	return nil
}

func vmOSConverter(object *ovirtsdk.Vm, v *vm) error {
	sdkOS, ok := object.Os()
	if !ok {
		return newFieldNotFound("vm", "os")
	}
	v.os = &vmOS{}
	osType, ok := sdkOS.Type()
	if !ok {
		return newFieldNotFound("os on vm", "type")
	}
	v.os.t = osType
	return nil
}

func vmTypeConverter(object *ovirtsdk.Vm, v *vm) error {
	vmType, ok := object.Type()
	if !ok {
		return newFieldNotFound("vm", "vm type")
	}
	v.vmType = VMType(vmType)
	return nil
}

func vmInstanceTypeIDConverter(object *ovirtsdk.Vm, v *vm) error {
	if instanceType, ok := object.InstanceType(); ok {
		instanceTypeID := InstanceTypeID(instanceType.MustId())
		v.instanceTypeID = &instanceTypeID
	}
	return nil
}

func vmMemoryPolicyConverter(object *ovirtsdk.Vm, v *vm) error {
	memPolicy, ok := object.MemoryPolicy()
	if !ok {
		return newFieldNotFound("vm", "memory policy")
	}
	resultMemPolicy := &memoryPolicy{}
	if guaranteed, ok := memPolicy.Guaranteed(); ok {
		if guaranteed < -1 {
			return newError(
				EBug,
				"the engine returned a negative guaranteed memory value for VM %s (%d)",
				object.MustId(),
				guaranteed,
			)
		}
		resultMemPolicy.guaranteed = &guaranteed
	}
	if max, ok := memPolicy.Max(); ok {
		resultMemPolicy.max = &max
	}
	if ballooning, ok := memPolicy.Ballooning(); ok {
		resultMemPolicy.ballooning = ballooning
	}
	v.memoryPolicy = resultMemPolicy
	return nil
}

func vmHostConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	if host, ok := sdkObject.Host(); ok {
		if hostID, ok := host.Id(); ok && hostID != "" {
			v.hostID = (*HostID)(&hostID)
		}
	}
	return nil
}

func vmPlacementPolicyConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	if pp, ok := sdkObject.PlacementPolicy(); ok {
		placementPolicy := &vmPlacementPolicy{}
		affinity, ok := pp.Affinity()
		if ok {
			a := VMAffinity(affinity)
			placementPolicy.affinity = &a
		}
		hosts, ok := pp.Hosts()
		if ok {
			hostIDs := make([]HostID, len(hosts.Slice()))
			for i, host := range hosts.Slice() {
				hostIDs[i] = HostID(host.MustId())
			}
			placementPolicy.hostIDs = hostIDs
		}
		v.placementPolicy = placementPolicy
	}
	return nil
}

func vmIDConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	id, ok := sdkObject.Id()
	if !ok {
		return newError(EFieldMissing, "id field missing from VM object")
	}
	v.id = VMID(id)
	return nil
}

func vmNameConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	name, ok := sdkObject.Name()
	if !ok {
		return newError(EFieldMissing, "name field missing from VM object")
	}
	v.name = name
	return nil
}

func vmCommentConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	comment, ok := sdkObject.Comment()
	if !ok {
		return newError(EFieldMissing, "comment field missing from VM object")
	}
	v.comment = comment
	return nil
}

func vmClusterConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	cluster, ok := sdkObject.Cluster()
	if !ok {
		return newError(EFieldMissing, "cluster field missing from VM object")
	}
	clusterID, ok := cluster.Id()
	if !ok {
		return newError(EFieldMissing, "ID field missing from cluster in VM object")
	}
	v.clusterID = ClusterID(clusterID)
	return nil
}

func vmStatusConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	status, ok := sdkObject.Status()
	if !ok {
		return newFieldNotFound("vm", "status")
	}
	v.status = VMStatus(status)
	return nil
}

func vmTemplateConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	template, ok := sdkObject.Template()
	if !ok {
		return newFieldNotFound("VM", "template")
	}
	templateID, ok := template.Id()
	if !ok {
		return newFieldNotFound("template in VM", "template ID")
	}
	v.templateID = TemplateID(templateID)
	return nil
}

func vmCPUConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	cpu, err := convertSDKVMCPU(sdkObject)
	if err != nil {
		return err
	}
	v.cpu = cpu
	return nil
}

func vmHugePagesConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	hugePages, err := hugePagesFromSDKVM(sdkObject)
	if err != nil {
		return err
	}
	v.hugePages = hugePages
	return nil
}

func vmMemoryConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	memory, ok := sdkObject.Memory()
	if !ok {
		return newFieldNotFound("vm", "memory")
	}
	v.memory = memory
	return nil
}

func vmInitializationConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	var vmInitialization *initialization
	vmInitialization, err := convertSDKInitialization(sdkObject)
	if err != nil {
		return err
	}
	v.initialization = vmInitialization
	return nil
}

func vmTagsConverter(sdkObject *ovirtsdk.Vm, v *vm) error {
	var tagIDs []TagID
	if sdkTags, ok := sdkObject.Tags(); ok {
		for _, tag := range sdkTags.Slice() {
			tagID := TagID(tag.MustId())
			tagIDs = append(tagIDs, tagID)
		}
	}
	v.tagIDs = tagIDs
	return nil
}

func convertSDKVMCPU(sdkObject *ovirtsdk.Vm) (*vmCPU, error) {
	sdkCPU, ok := sdkObject.Cpu()
	if !ok {
		return nil, newFieldNotFound("VM", "CPU")
	}
	cpuTopo, ok := sdkCPU.Topology()
	if !ok {
		return nil, newFieldNotFound("CPU in VM", "CPU topo")
	}
	cores, ok := cpuTopo.Cores()
	if !ok {
		return nil, newFieldNotFound("CPU topo in CPU in VM", "cores")
	}
	threads, ok := cpuTopo.Threads()
	if !ok {
		return nil, newFieldNotFound("CPU topo in CPU in VM", "threads")
	}
	sockets, ok := cpuTopo.Sockets()
	if !ok {
		return nil, newFieldNotFound("CPU topo in CPU in VM", "sockets")
	}
	sdkCPUMode, ok := sdkCPU.Mode()
	var cpuMode *CPUMode
	if ok {
		cpuMode = (*CPUMode)(&sdkCPUMode)
	}
	cpu := &vmCPU{
		topo: &vmCPUTopo{
			uint(cores),
			uint(threads),
			uint(sockets),
		},
		mode: cpuMode,
	}
	return cpu, nil
}

// VMStatus represents the status of a VM.
type VMStatus string

const (
	// VMStatusDown indicates that the VM is not running.
	VMStatusDown VMStatus = "down"
	// VMStatusImageLocked indicates that the virtual machine process is not running and there is some operation on the
	// disks of the virtual machine that prevents it from being started.
	VMStatusImageLocked VMStatus = "image_locked"
	// VMStatusMigrating indicates that the virtual machine process is running and the virtual machine is being migrated
	// from one host to another.
	VMStatusMigrating VMStatus = "migrating"
	// VMStatusNotResponding indicates that the hypervisor detected that the virtual machine is not responding.
	VMStatusNotResponding VMStatus = "not_responding"
	// VMStatusPaused indicates that the virtual machine process is running and the virtual machine is paused.
	// This may happen in two cases: when running a virtual machine is paused mode and when the virtual machine is being
	// automatically paused due to an error.
	VMStatusPaused VMStatus = "paused"
	// VMStatusPoweringDown indicates that the virtual machine process is running and it is about to stop running.
	VMStatusPoweringDown VMStatus = "powering_down"
	// VMStatusPoweringUp  indicates that the virtual machine process is running and the guest operating system is being
	// loaded. Note that if no guest-agent is installed, this status is set for a predefined period of time, that is by
	// default 60 seconds, when running a virtual machine.
	VMStatusPoweringUp VMStatus = "powering_up"
	// VMStatusRebooting indicates that the virtual machine process is running and the guest operating system is being
	// rebooted.
	VMStatusRebooting VMStatus = "reboot_in_progress"
	// VMStatusRestoringState indicates that the virtual machine process is about to run and the virtual machine is
	// going to awake from hibernation. In this status, the running state of the virtual machine is being restored.
	VMStatusRestoringState VMStatus = "restoring_state"
	// VMStatusSavingState indicates that the virtual machine process is running and the virtual machine is being
	// hibernated. In this status, the running state of the virtual machine is being saved. Note that this status does
	// not mean that the guest operating system is being hibernated.
	VMStatusSavingState VMStatus = "saving_state"
	// VMStatusSuspended indicates that the virtual machine process is not running and a running state of the virtual
	// machine was saved. This status is similar to Down, but when the VM is started in this status its saved running
	// state is restored instead of being booted using the normal procedure.
	VMStatusSuspended VMStatus = "suspended"
	// VMStatusUnassigned means an invalid status was received.
	VMStatusUnassigned VMStatus = "unassigned"
	// VMStatusUnknown indicates that the system failed to determine the status of the virtual machine.
	// The virtual machine process may be running or not running in this status.
	// For instance, when host becomes non-responsive the virtual machines that ran on it are set with this status.
	VMStatusUnknown VMStatus = "unknown"
	// VMStatusUp indicates that the virtual machine process is running and the guest operating system is loaded.
	// Note that if no guest-agent is installed, this status is set after a predefined period of time, that is by
	// default 60 seconds, when running a virtual machine.
	VMStatusUp VMStatus = "up"
	// VMStatusWaitForLaunch indicates that the virtual machine process is about to run.
	// This status is set when a request to run a virtual machine arrives to the host.
	// It is possible that the virtual machine process will fail to run.
	VMStatusWaitForLaunch VMStatus = "wait_for_launch"
)

type vmPlacementPolicy struct {
	affinity *VMAffinity
	hostIDs  []HostID
}

func (v vmPlacementPolicy) Affinity() *VMAffinity {
	return v.affinity
}

func (v vmPlacementPolicy) HostIDs() []HostID {
	return v.hostIDs
}

// Validate validates if a VMStatus has a valid value.
func (s VMStatus) Validate() error {
	for _, v := range VMStatusValues() {
		if v == s {
			return nil
		}
	}
	return newError(EBadArgument, "invalid value for VM status: %s", s)
}

// VMStatusList is a list of VMStatus.
type VMStatusList []VMStatus

// Copy creates a separate copy of the current status list.
func (l VMStatusList) Copy() VMStatusList {
	result := make([]VMStatus, len(l))
	// nolint:gosimple
	for i, s := range l {
		result[i] = s
	}
	return result
}

// Validate validates the list of statuses.
func (l VMStatusList) Validate() error {
	for _, s := range l {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// VMStatusValues returns all possible VMStatus values.
func VMStatusValues() VMStatusList {
	return []VMStatus{
		VMStatusDown,
		VMStatusImageLocked,
		VMStatusMigrating,
		VMStatusNotResponding,
		VMStatusPaused,
		VMStatusPoweringDown,
		VMStatusPoweringUp,
		VMStatusRebooting,
		VMStatusRestoringState,
		VMStatusSavingState,
		VMStatusSuspended,
		VMStatusUnassigned,
		VMStatusUnknown,
		VMStatusUp,
		VMStatusWaitForLaunch,
	}
}

// Strings creates a string list of the values.
func (l VMStatusList) Strings() []string {
	result := make([]string, len(l))
	for i, status := range l {
		result[i] = string(status)
	}
	return result
}

func hugePagesFromSDKVM(vm *ovirtsdk.Vm) (*VMHugePages, error) {
	var hugePagesText string
	customProperties, ok := vm.CustomProperties()
	if !ok {
		return nil, nil
	}
	for _, c := range customProperties.Slice() {
		customPropertyName, ok := c.Name()
		if !ok {
			return nil, nil
		}
		if customPropertyName == "hugepages" {
			hugePagesText, ok = c.Value()
			if !ok {
				return nil, nil
			}
			break
		}
	}
	hugepagesUint, err := strconv.ParseUint(hugePagesText, 10, 64)
	if err != nil {
		return nil, wrap(err, EBug, "Failed to parse 'hugepages' custom property into a number: %s", hugePagesText)
	}
	hugepages := VMHugePages(hugepagesUint)
	return &hugepages, nil
}
