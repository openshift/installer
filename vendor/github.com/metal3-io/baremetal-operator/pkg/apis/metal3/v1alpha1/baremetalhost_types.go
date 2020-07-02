package v1alpha1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// NOTE: json tags are required.  Any new fields you add must have
// json tags for the fields to be serialized.

// NOTE(dhellmann): Update docs/api.md when changing these data structure.

const (
	// BareMetalHostFinalizer is the name of the finalizer added to
	// hosts to block delete operations until the physical host can be
	// deprovisioned.
	BareMetalHostFinalizer string = "baremetalhost.metal3.io"

	// PausedAnnotation is the annotation that pauses the reconciliation (triggers
	// an immediate requeue)
	PausedAnnotation = "baremetalhost.metal3.io/paused"

	// StatusAnnotation is the annotation that keeps a copy of the Status of BMH
	// This is particularly useful when we pivot BMH. If the status
	// annotation is present and status is empty, BMO will reconstruct BMH Status
	// from the status annotation.
	StatusAnnotation = "baremetalhost.metal3.io/status"

	// UnhealthyAnnotation is the annotation that sets unhealthy status of BMH
	UnhealthyAnnotation = "baremetalhost.metal3.io/unhealthy"
)

// RootDeviceHints holds the hints for specifying the storage location
// for the root filesystem for the image.
type RootDeviceHints struct {
	// A Linux device name like "/dev/vda". The hint must match the
	// actual value exactly.
	DeviceName string `json:"deviceName,omitempty"`

	// A SCSI bus address like 0:0:0:0. The hint must match the actual
	// value exactly.
	HCTL string `json:"hctl,omitempty"`

	// A vendor-specific device identifier. The hint can be a
	// substring of the actual value.
	Model string `json:"model,omitempty"`

	// The name of the vendor or manufacturer of the device. The hint
	// can be a substring of the actual value.
	Vendor string `json:"vendor,omitempty"`

	// Device serial number. The hint must match the actual value
	// exactly.
	SerialNumber string `json:"serialNumber,omitempty"`

	// The minimum size of the device in Gigabytes.
	// +kubebuilder:validation:Minimum=0
	MinSizeGigabytes int `json:"minSizeGigabytes,omitempty"`

	// Unique storage identifier. The hint must match the actual value
	// exactly.
	WWN string `json:"wwn,omitempty"`

	// Unique storage identifier with the vendor extension
	// appended. The hint must match the actual value exactly.
	WWNWithExtension string `json:"wwnWithExtension,omitempty"`

	// Unique vendor storage identifier. The hint must match the
	// actual value exactly.
	WWNVendorExtension string `json:"wwnVendorExtension,omitempty"`

	// True if the device should use spinning media, false otherwise.
	Rotational *bool `json:"rotational,omitempty"`
}

// OperationalStatus represents the state of the host
type OperationalStatus string

const (
	// OperationalStatusOK is the status value for when the host is
	// configured correctly and is manageable.
	OperationalStatusOK OperationalStatus = "OK"

	// OperationalStatusDiscovered is the status value for when the
	// host is only partially configured, such as when when the BMC
	// address is known but the login credentials are not.
	OperationalStatusDiscovered OperationalStatus = "discovered"

	// OperationalStatusError is the status value for when the host
	// has any sort of error.
	OperationalStatusError OperationalStatus = "error"
)

// ErrorType indicates the class of problem that has caused the Host resource
// to enter an error state.
type ErrorType string

const (
	// RegistrationError is an error condition occurring when the
	// controller is unable to connect to the Host's baseboard management
	// controller.
	RegistrationError ErrorType = "registration error"
	// InspectionError is an error condition occurring when an attempt to
	// obtain hardware details from the Host fails.
	InspectionError ErrorType = "inspection error"
	// ProvisioningError is an error condition occuring when the controller
	// fails to provision or deprovision the Host.
	ProvisioningError ErrorType = "provisioning error"
	// PowerManagementError is an error condition occurring when the
	// controller is unable to modify the power state of the Host.
	PowerManagementError ErrorType = "power management error"
)

// ProvisioningState defines the states the provisioner will report
// the host has having.
type ProvisioningState string

const (
	// StateNone means the state is unknown
	StateNone ProvisioningState = ""

	// StateRegistrationError means there was an error registering the
	// host with the backend
	StateRegistrationError ProvisioningState = "registration error"

	// StateRegistering means we are telling the backend about the host
	StateRegistering ProvisioningState = "registering"

	// StateMatchProfile means we are comparing the discovered details
	// against known hardware profiles
	StateMatchProfile ProvisioningState = "match profile"

	// StateReady means the host can be consumed
	StateReady ProvisioningState = "ready"

	// StateAvailable means the host can be consumed
	StateAvailable ProvisioningState = "available"

	// StateProvisioning means we are writing an image to the host's
	// disk(s)
	StateProvisioning ProvisioningState = "provisioning"

	// StateProvisioningError means we are writing an image to the
	// host's disk(s)
	StateProvisioningError ProvisioningState = "provisioning error"

	// StateProvisioned means we have written an image to the host's
	// disk(s)
	StateProvisioned ProvisioningState = "provisioned"

	// StateExternallyProvisioned means something else is managing the
	// image on the host
	StateExternallyProvisioned ProvisioningState = "externally provisioned"

	// StateDeprovisioning means we are removing an image from the
	// host's disk(s)
	StateDeprovisioning ProvisioningState = "deprovisioning"

	// StateInspecting means we are running the agent on the host to
	// learn about the hardware components available there
	StateInspecting ProvisioningState = "inspecting"

	// StatePowerManagementError means something went wrong trying to
	// power the server on or off.
	StatePowerManagementError ProvisioningState = "power management error"

	// StateDeleting means we are in the process of cleaning up the host
	// ready for deletion
	StateDeleting ProvisioningState = "deleting"
)

// BMCDetails contains the information necessary to communicate with
// the bare metal controller module on host.
type BMCDetails struct {

	// Address holds the URL for accessing the controller on the
	// network.
	Address string `json:"address"`

	// The name of the secret containing the BMC credentials (requires
	// keys "username" and "password").
	CredentialsName string `json:"credentialsName"`

	// DisableCertificateVerification disables verification of server
	// certificates when using HTTPS to connect to the BMC. This is
	// required when the server certificate is self-signed, but is
	// insecure because it allows a man-in-the-middle to intercept the
	// connection.
	DisableCertificateVerification bool `json:"disableCertificateVerification,omitempty"`
}

// BareMetalHostSpec defines the desired state of BareMetalHost
type BareMetalHostSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code
	// after modifying this file

	// Taints is the full, authoritative list of taints to apply to
	// the corresponding Machine. This list will overwrite any
	// modifications made to the Machine on an ongoing basis.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`

	// How do we connect to the BMC?
	BMC BMCDetails `json:"bmc,omitempty"`

	// What is the name of the hardware profile for this host? It
	// should only be necessary to set this when inspection cannot
	// automatically determine the profile.
	HardwareProfile string `json:"hardwareProfile,omitempty"`

	// Provide guidance about how to choose the device for the image
	// being provisioned.
	RootDeviceHints *RootDeviceHints `json:"rootDeviceHints,omitempty"`

	// Which MAC address will PXE boot? This is optional for some
	// types, but required for libvirt VMs driven by vbmc.
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	BootMACAddress string `json:"bootMACAddress,omitempty"`

	// Should the server be online?
	Online bool `json:"online"`

	// ConsumerRef can be used to store information about something
	// that is using a host. When it is not empty, the host is
	// considered "in use".
	ConsumerRef *corev1.ObjectReference `json:"consumerRef,omitempty"`

	// Image holds the details of the image to be provisioned.
	Image *Image `json:"image,omitempty"`

	// UserData holds the reference to the Secret containing the user
	// data to be passed to the host before it boots.
	UserData *corev1.SecretReference `json:"userData,omitempty"`

	// NetworkData holds the reference to the Secret containing network
	// configuration (e.g content of network_data.json which is passed
	// to Config Drive).
	NetworkData *corev1.SecretReference `json:"networkData,omitempty"`

	// MetaData holds the reference to the Secret containing host metadata
	// (e.g. meta_data.json which is passed to Config Drive).
	MetaData *corev1.SecretReference `json:"metaData,omitempty"`

	// Description is a human-entered text used to help identify the host
	Description string `json:"description,omitempty"`

	// ExternallyProvisioned means something else is managing the
	// image running on the host and the operator should only manage
	// the power status and hardware inventory inspection. If the
	// Image field is filled in, this field is ignored.
	ExternallyProvisioned bool `json:"externallyProvisioned,omitempty"`
}

// ChecksumType holds the algorithm name for the checksum
// +kubebuilder:validation:Enum=md5;sha256;sha512
type ChecksumType string

const (
	// MD5 checksum type
	MD5 ChecksumType = "md5"

	// SHA256 checksum type
	SHA256 ChecksumType = "sha256"

	// SHA512 checksum type
	SHA512 ChecksumType = "sha512"
)

// Image holds the details of an image either to provisioned or that
// has been provisioned.
type Image struct {
	// URL is a location of an image to deploy.
	URL string `json:"url"`

	// Checksum is the checksum for the image.
	Checksum string `json:"checksum"`

	// ChecksumType is the checksum algorithm for the image.
	// e.g md5, sha256, sha512
	ChecksumType ChecksumType `json:"checksumType,omitempty"`

	// DiskFormat contains the format of the image (raw, qcow2, ...)
	// Needs to be set to raw for raw images streaming
	// +kubebuilder:validation:Enum=raw;qcow2;vdi;vmdk
	DiskFormat *string `json:"format,omitempty"`
}

// FIXME(dhellmann): We probably want some other module to own these
// data structures.

// ClockSpeed is a clock speed in MHz
type ClockSpeed float64

// ClockSpeed multipliers
const (
	MegaHertz ClockSpeed = 1.0
	GigaHertz            = 1000 * MegaHertz
)

// Capacity is a disk size in Bytes
type Capacity int64

// Capacity multipliers
const (
	Byte     Capacity = 1
	KibiByte          = Byte * 1024
	KiloByte          = Byte * 1000
	MebiByte          = KibiByte * 1024
	MegaByte          = KiloByte * 1000
	GibiByte          = MebiByte * 1024
	GigaByte          = MegaByte * 1000
	TebiByte          = GibiByte * 1024
	TeraByte          = GigaByte * 1000
)

// CPU describes one processor on the host.
type CPU struct {
	Arch           string     `json:"arch"`
	Model          string     `json:"model"`
	ClockMegahertz ClockSpeed `json:"clockMegahertz"`
	Flags          []string   `json:"flags"`
	Count          int        `json:"count"`
}

// Storage describes one storage device (disk, SSD, etc.) on the host.
type Storage struct {
	// A name for the disk, e.g. "disk 1 (boot)"
	Name string `json:"name"`

	// Whether this disk represents rotational storage
	Rotational bool `json:"rotational"`

	// The size of the disk in Bytes
	SizeBytes Capacity `json:"sizeBytes"`

	// The name of the vendor of the device
	Vendor string `json:"vendor,omitempty"`

	// Hardware model
	Model string `json:"model,omitempty"`

	// The serial number of the device
	SerialNumber string `json:"serialNumber"`

	// The WWN of the device
	WWN string `json:"wwn,omitempty"`

	// The WWN Vendor extension of the device
	WWNVendorExtension string `json:"wwnVendorExtension,omitempty"`

	// The WWN with the extension
	WWNWithExtension string `json:"wwnWithExtension,omitempty"`

	// The SCSI location of the device
	HCTL string `json:"hctl,omitempty"`
}

// VLANID is a 12-bit 802.1Q VLAN identifier
// +kubebuilder:validation:Type=integer
// +kubebuilder:validation:Minimum=0
// +kubebuilder:validation:Maximum=4094
type VLANID int32

// VLAN represents the name and ID of a VLAN
type VLAN struct {
	ID VLANID `json:"id"`

	Name string `json:"name,omitempty"`
}

// NIC describes one network interface on the host.
type NIC struct {
	// The name of the NIC, e.g. "nic-1"
	Name string `json:"name"`

	// The name of the model, e.g. "virt-io"
	Model string `json:"model"`

	// The device MAC addr
	// +kubebuilder:validation:Pattern=`[0-9a-fA-F]{2}(:[0-9a-fA-F]{2}){5}`
	MAC string `json:"mac"`

	// The IP address of the device
	IP string `json:"ip"`

	// The speed of the device
	SpeedGbps int `json:"speedGbps"`

	// The VLANs available
	VLANs []VLAN `json:"vlans,omitempty"`

	// The untagged VLAN ID
	VLANID VLANID `json:"vlanId"`

	// Whether the NIC is PXE Bootable
	PXE bool `json:"pxe"`
}

// Firmware describes the firmware on the host.
type Firmware struct {
	// The BIOS for this firmware
	BIOS BIOS `json:"bios"`
}

// BIOS describes the BIOS version on the host.
type BIOS struct {
	// The release/build date for this BIOS
	Date string `json:"date"`

	// The vendor name for this BIOS
	Vendor string `json:"vendor"`

	// The version of the BIOS
	Version string `json:"version"`
}

// HardwareDetails collects all of the information about hardware
// discovered on the host.
type HardwareDetails struct {
	SystemVendor HardwareSystemVendor `json:"systemVendor"`
	Firmware     Firmware             `json:"firmware"`
	RAMMebibytes int                  `json:"ramMebibytes"`
	NIC          []NIC                `json:"nics"`
	Storage      []Storage            `json:"storage"`
	CPU          CPU                  `json:"cpu"`
	Hostname     string               `json:"hostname"`
}

// HardwareSystemVendor stores details about the whole hardware system.
type HardwareSystemVendor struct {
	Manufacturer string `json:"manufacturer"`
	ProductName  string `json:"productName"`
	SerialNumber string `json:"serialNumber"`
}

// CredentialsStatus contains the reference and version of the last
// set of BMC credentials the controller was able to validate.
type CredentialsStatus struct {
	Reference *corev1.SecretReference `json:"credentials,omitempty"`
	Version   string                  `json:"credentialsVersion,omitempty"`
}

// Match compares the saved status information with the name and
// content of a secret object.
func (cs CredentialsStatus) Match(secret corev1.Secret) bool {
	switch {
	case cs.Reference == nil:
		return false
	case cs.Reference.Name != secret.ObjectMeta.Name:
		return false
	case cs.Reference.Namespace != secret.ObjectMeta.Namespace:
		return false
	case cs.Version != secret.ObjectMeta.ResourceVersion:
		return false
	}
	return true
}

// OperationMetric contains metadata about an operation (inspection,
// provisioning, etc.) used for tracking metrics.
type OperationMetric struct {
	// +nullable
	Start metav1.Time `json:"start,omitempty"`
	// +nullable
	End metav1.Time `json:"end,omitempty"`
}

// Duration returns the length of time that was spent on the
// operation. If the operation is not finished, it returns 0.
func (om OperationMetric) Duration() time.Duration {
	if om.Start.IsZero() {
		return 0
	}
	return om.End.Time.Sub(om.Start.Time)
}

// OperationHistory holds information about operations performed on a
// host.
type OperationHistory struct {
	Register    OperationMetric `json:"register,omitempty"`
	Inspect     OperationMetric `json:"inspect,omitempty"`
	Provision   OperationMetric `json:"provision,omitempty"`
	Deprovision OperationMetric `json:"deprovision,omitempty"`
}

// BareMetalHostStatus defines the observed state of BareMetalHost
type BareMetalHostStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code
	// after modifying this file

	// OperationalStatus holds the status of the host
	// +kubebuilder:validation:Enum="";OK;discovered;error
	OperationalStatus OperationalStatus `json:"operationalStatus"`

	// ErrorType indicates the type of failure encountered when the
	// OperationalStatus is OperationalStatusError
	// +kubebuilder:validation:Enum=registration error;inspection error;provisioning error;power management error
	ErrorType ErrorType `json:"errorType,omitempty"`

	// LastUpdated identifies when this status was last observed.
	// +optional
	LastUpdated *metav1.Time `json:"lastUpdated,omitempty"`

	// The name of the profile matching the hardware details.
	HardwareProfile string `json:"hardwareProfile"`

	// The hardware discovered to exist on the host.
	HardwareDetails *HardwareDetails `json:"hardware,omitempty"`

	// Information tracked by the provisioner.
	Provisioning ProvisionStatus `json:"provisioning"`

	// the last credentials we were able to validate as working
	GoodCredentials CredentialsStatus `json:"goodCredentials,omitempty"`

	// the last credentials we sent to the provisioning backend
	TriedCredentials CredentialsStatus `json:"triedCredentials,omitempty"`

	// the last error message reported by the provisioning subsystem
	ErrorMessage string `json:"errorMessage"`

	// indicator for whether or not the host is powered on
	PoweredOn bool `json:"poweredOn"`

	// OperationHistory holds information about operations performed
	// on this host.
	OperationHistory OperationHistory `json:"operationHistory"`
}

// ProvisionStatus holds the state information for a single target.
type ProvisionStatus struct {
	// An indiciator for what the provisioner is doing with the host.
	State ProvisioningState `json:"state"`

	// The machine's UUID from the underlying provisioning tool
	ID string `json:"ID"`

	// Image holds the details of the last image successfully
	// provisioned to the host.
	Image Image `json:"image,omitempty"`

	// The RootDevicehints set by the user
	RootDeviceHints *RootDeviceHints `json:"rootDeviceHints,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BareMetalHost is the Schema for the baremetalhosts API
// +k8s:openapi-gen=true
// +kubebuilder:resource:shortName=bmh;bmhost
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.operationalStatus",description="Operational status"
// +kubebuilder:printcolumn:name="Provisioning Status",type="string",JSONPath=".status.provisioning.state",description="Provisioning status"
// +kubebuilder:printcolumn:name="Consumer",type="string",JSONPath=".spec.consumerRef.name",description="Consumer using this host"
// +kubebuilder:printcolumn:name="BMC",type="string",JSONPath=".spec.bmc.address",description="Address of management controller"
// +kubebuilder:printcolumn:name="Hardware Profile",type="string",JSONPath=".status.hardwareProfile",description="The type of hardware detected"
// +kubebuilder:printcolumn:name="Online",type="string",JSONPath=".spec.online",description="Whether the host is online or not"
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.errorMessage",description="Most recent error"
type BareMetalHost struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BareMetalHostSpec   `json:"spec,omitempty"`
	Status BareMetalHostStatus `json:"status,omitempty"`
}

// Available returns true if the host is available to be provisioned.
func (host *BareMetalHost) Available() bool {
	if host.Spec.ConsumerRef != nil {
		return false
	}
	if host.GetDeletionTimestamp() != nil {
		return false
	}
	if host.HasError() {
		return false
	}
	return true
}

// SetErrorMessage updates the ErrorMessage in the host Status struct
// when necessary and returns true when a change is made or false when
// no change is made.
func (host *BareMetalHost) SetErrorMessage(errType ErrorType, message string) (dirty bool) {
	if host.Status.OperationalStatus != OperationalStatusError {
		host.Status.OperationalStatus = OperationalStatusError
		dirty = true
	}
	if host.Status.ErrorType != errType {
		host.Status.ErrorType = errType
		dirty = true
	}
	if host.Status.ErrorMessage != message {
		host.Status.ErrorMessage = message
		dirty = true
	}
	return dirty
}

// ClearError removes any existing error message.
func (host *BareMetalHost) ClearError() (dirty bool) {
	dirty = host.SetOperationalStatus(OperationalStatusOK)
	var emptyErrType ErrorType = ""
	if host.Status.ErrorType != emptyErrType {
		host.Status.ErrorType = emptyErrType
		dirty = true
	}
	if host.Status.ErrorMessage != "" {
		host.Status.ErrorMessage = ""
		dirty = true
	}
	return dirty
}

// setLabel updates the given label when necessary and returns true
// when a change is made or false when no change is made.
func (host *BareMetalHost) setLabel(name, value string) bool {
	if host.Labels == nil {
		host.Labels = make(map[string]string)
	}
	if host.Labels[name] != value {
		host.Labels[name] = value
		return true
	}
	return false
}

// getLabel returns the value associated with the given label. If
// there is no value, an empty string is returned.
func (host *BareMetalHost) getLabel(name string) string {
	if host.Labels == nil {
		return ""
	}
	return host.Labels[name]
}

// NeedsHardwareProfile returns true if the profile is not set
func (host *BareMetalHost) NeedsHardwareProfile() bool {
	return host.Status.HardwareProfile == ""
}

// HardwareProfile returns the hardware profile name for the host.
func (host *BareMetalHost) HardwareProfile() string {
	return host.Status.HardwareProfile
}

// SetHardwareProfile updates the hardware profile name and returns
// true when a change is made or false when no change is made.
func (host *BareMetalHost) SetHardwareProfile(name string) (dirty bool) {
	if host.Status.HardwareProfile != name {
		host.Status.HardwareProfile = name
		dirty = true
	}
	return dirty
}

// SetOperationalStatus updates the OperationalStatus field and returns
// true when a change is made or false when no change is made.
func (host *BareMetalHost) SetOperationalStatus(status OperationalStatus) bool {
	if host.Status.OperationalStatus != status {
		host.Status.OperationalStatus = status
		return true
	}
	return false
}

// OperationalStatus returns the contents of the OperationalStatus
// field.
func (host *BareMetalHost) OperationalStatus() OperationalStatus {
	return host.Status.OperationalStatus
}

// HasError returns a boolean indicating whether there is an error
// set for the host.
func (host *BareMetalHost) HasError() bool {
	return host.Status.ErrorMessage != ""
}

// CredentialsKey returns a NamespacedName suitable for loading the
// Secret containing the credentials associated with the host.
func (host *BareMetalHost) CredentialsKey() types.NamespacedName {
	return types.NamespacedName{
		Name:      host.Spec.BMC.CredentialsName,
		Namespace: host.ObjectMeta.Namespace,
	}
}

// NeedsHardwareInspection looks at the state of the host to determine
// if hardware inspection should be run.
func (host *BareMetalHost) NeedsHardwareInspection() bool {
	if host.Spec.ExternallyProvisioned {
		// Never perform inspection if we already know something is
		// using the host and we didn't provision it.
		return false
	}
	if host.WasProvisioned() {
		// Never perform inspection if we have already provisioned
		// this host, because we don't want to reboot it.
		return false
	}
	return host.Status.HardwareDetails == nil
}

// NeedsProvisioning compares the settings with the provisioning
// status and returns true when more work is needed or false
// otherwise.
func (host *BareMetalHost) NeedsProvisioning() bool {
	if !host.Spec.Online {
		// The host is not supposed to be powered on.
		return false
	}
	if host.Spec.Image == nil {
		// Without an image, there is nothing to provision.
		return false
	}
	if host.Spec.Image.URL == "" {
		// We have an Image struct but it is empty
		return false
	}
	if host.Status.Provisioning.Image.URL == "" {
		// We have an image set, but not provisioned.
		return true
	}
	return false
}

// WasProvisioned returns true when we think we have placed an image
// on the host.
func (host *BareMetalHost) WasProvisioned() bool {
	if host.Spec.ExternallyProvisioned {
		return false
	}
	if host.Status.Provisioning.Image.URL != "" {
		// We have an image provisioned.
		return true
	}
	return false
}

// NeedsDeprovisioning compares the settings with the provisioning
// status and returns true when the host should be deprovisioned.
func (host *BareMetalHost) NeedsDeprovisioning() bool {
	if host.Spec.Image == nil {
		return true
	}
	if host.Spec.Image.URL == "" {
		return true
	}
	if host.Status.Provisioning.Image.URL == "" {
		return false
	}
	if host.Spec.Image.URL != host.Status.Provisioning.Image.URL {
		return true
	}
	return false
}

// UpdateGoodCredentials modifies the GoodCredentials portion of the
// Status struct to record the details of the secret containing
// credentials known to work.
func (host *BareMetalHost) UpdateGoodCredentials(currentSecret corev1.Secret) {
	host.Status.GoodCredentials.Version = currentSecret.ObjectMeta.ResourceVersion
	host.Status.GoodCredentials.Reference = &corev1.SecretReference{
		Name:      currentSecret.ObjectMeta.Name,
		Namespace: currentSecret.ObjectMeta.Namespace,
	}
}

// UpdateTriedCredentials modifies the TriedCredentials portion of the
// Status struct to record the details of the secret containing
// credentials known to work.
func (host *BareMetalHost) UpdateTriedCredentials(currentSecret corev1.Secret) {
	host.Status.TriedCredentials.Version = currentSecret.ObjectMeta.ResourceVersion
	host.Status.TriedCredentials.Reference = &corev1.SecretReference{
		Name:      currentSecret.ObjectMeta.Name,
		Namespace: currentSecret.ObjectMeta.Namespace,
	}
}

// NewEvent creates a new event associated with the object and ready
// to be published to the kubernetes API.
func (host *BareMetalHost) NewEvent(reason, message string) corev1.Event {
	t := metav1.Now()
	return corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: reason + "-",
			Namespace:    host.ObjectMeta.Namespace,
		},
		InvolvedObject: corev1.ObjectReference{
			Kind:       "BareMetalHost",
			Namespace:  host.Namespace,
			Name:       host.Name,
			UID:        host.UID,
			APIVersion: SchemeGroupVersion.String(),
		},
		Reason:  reason,
		Message: message,
		Source: corev1.EventSource{
			Component: "metal3-baremetal-controller",
		},
		FirstTimestamp:      t,
		LastTimestamp:       t,
		Count:               1,
		Type:                corev1.EventTypeNormal,
		ReportingController: "metal3.io/baremetal-controller",
		Related:             host.Spec.ConsumerRef,
	}
}

// OperationMetricForState returns a pointer to the metric for the given
// provisioning state.
func (host *BareMetalHost) OperationMetricForState(operation ProvisioningState) (metric *OperationMetric) {
	history := &host.Status.OperationHistory
	switch operation {
	case StateRegistering:
		metric = &history.Register
	case StateInspecting:
		metric = &history.Inspect
	case StateProvisioning:
		metric = &history.Provision
	case StateDeprovisioning:
		metric = &history.Deprovision
	}
	return
}

// GetImageChecksum returns the hash value and its algo.
func (host *BareMetalHost) GetImageChecksum() (string, string, bool) {
	if host.Spec.Image == nil {
		return "", "", false
	}

	checksum := host.Spec.Image.Checksum
	checksumType := host.Spec.Image.ChecksumType

	if checksum == "" {
		// Return empty if checksum is not provided
		return "", "", false
	}
	if checksumType == "" {
		// If only checksum is specified. Assume type is md5
		return checksum, string(MD5), true
	}
	switch checksumType {
	case MD5, SHA256, SHA512:
		return checksum, string(checksumType), true
	default:
		return "", "", false
	}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BareMetalHostList contains a list of BareMetalHost
type BareMetalHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BareMetalHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BareMetalHost{}, &BareMetalHostList{})
}
