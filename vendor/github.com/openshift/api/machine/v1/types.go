package v1

// ManagementState denotes whether the resource is expected to be managed by the controller or by the user.
type ManagementState string

const (
	// ManagedManagementState denotes that the resource is expected to be managed by the controller.
	ManagedManagementState ManagementState = "Managed"

	// UnmanagedManagementState denotes that the resource is expected to be managed by the user.
	UnmanagedManagementState ManagementState = "Unmanaged"
)

// LocalSecretReference contains enough information to let you locate the
// referenced Secret inside the same namespace.
// +structType=atomic
type LocalSecretReference struct {
	// Name of the Secret.
	// +kubebuilder:validation:=Required
	Name string `json:"name"`
}
