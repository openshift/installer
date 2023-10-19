// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

// +kubebuilder:rbac:groups=dbforpostgresql.azure.com,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbforpostgresql.azure.com,resources={users/status,users/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// +kubebuilder:storageversion
// User is a postgresql user
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              UserSpec   `json:"spec,omitempty"`
	Status            UserStatus `json:"status,omitempty"`
}

var _ conditions.Conditioner = &User{}

// GetConditions returns the conditions of the resource
func (user *User) GetConditions() conditions.Conditions {
	return user.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (user *User) SetConditions(conditions conditions.Conditions) {
	user.Status.Conditions = conditions
}

// +kubebuilder:webhook:path=/mutate-dbforpostgresql-azure-com-v1-user,mutating=true,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=dbforpostgresql.azure.com,resources=users,verbs=create;update,versions=v1,name=default.v1.users.dbforpostgresql.azure.com,admissionReviewVersions=v1

var _ admission.Defaulter = &User{}

// Default applies defaults to the FlexibleServer resource
func (user *User) Default() {
	user.defaultImpl()
	var temp interface{} = user
	if runtimeDefaulter, ok := temp.(genruntime.Defaulter); ok {
		runtimeDefaulter.CustomDefault()
	}
}

// defaultAzureName defaults the Azure name of the resource to the Kubernetes name
func (user *User) defaultAzureName() {
	if user.Spec.AzureName == "" {
		user.Spec.AzureName = user.Name
	}
}

// defaultImpl applies the code generated defaults to the FlexibleServer resource
func (user *User) defaultImpl() { user.defaultAzureName() }

var _ genruntime.ARMOwned = &User{}

// AzureName returns the Azure name of the resource
func (user *User) AzureName() string {
	return user.Spec.AzureName
}

// Owner returns the ResourceReference of the owner, or nil if there is no owner
func (user *User) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(user.Spec)
	return user.Spec.Owner.AsResourceReference(group, kind)
}

// +kubebuilder:webhook:path=/validate-dbforpostgresql-azure-com-v1-user,mutating=false,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=dbforpostgresql.azure.com,resources=users,verbs=create;update,versions=v1,name=validate.v1.users.dbforpostgresql.azure.com,admissionReviewVersions=v1

var _ admission.Validator = &User{}

// ValidateCreate validates the creation of the resource
func (user *User) ValidateCreate() (admission.Warnings, error) {
	validations := user.createValidations()
	var temp interface{} = user
	if runtimeValidator, ok := temp.(genruntime.Validator); ok {
		validations = append(validations, runtimeValidator.CreateValidations()...)
	}
	return genruntime.ValidateCreate(validations)
}

// ValidateDelete validates the deletion of the resource
func (user *User) ValidateDelete() (admission.Warnings, error) {
	validations := user.deleteValidations()
	var temp interface{} = user
	if runtimeValidator, ok := temp.(genruntime.Validator); ok {
		validations = append(validations, runtimeValidator.DeleteValidations()...)
	}
	return genruntime.ValidateDelete(validations)
}

// ValidateUpdate validates an update of the resource
func (user *User) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	validations := user.updateValidations()
	var temp interface{} = user
	if runtimeValidator, ok := temp.(genruntime.Validator); ok {
		validations = append(validations, runtimeValidator.UpdateValidations()...)
	}
	return genruntime.ValidateUpdate(old, validations)
}

// createValidations validates the creation of the resource
func (user *User) createValidations() []func() (admission.Warnings, error) {
	return nil
}

// deleteValidations validates the deletion of the resource
func (user *User) deleteValidations() []func() (admission.Warnings, error) {
	return nil
}

// updateValidations validates the update of the resource
func (user *User) updateValidations() []func(old runtime.Object) (admission.Warnings, error) {
	return nil
}

var _ conversion.Hub = &User{}

// Hub marks that this userSpec is the hub type for conversion
func (user *User) Hub() {}

// +kubebuilder:object:root=true
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

type UserSpec struct {
	//AzureName: The name of the resource in Azure. This is often the same as the name of the resource in Kubernetes but it
	//doesn't have to be.
	AzureName string `json:"azureName,omitempty"`

	// +kubebuilder:validation:Required
	//Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	//controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	//reference to a dbforpostgresql.azure.com/FlexibleServer resource
	Owner *genruntime.KubernetesOwnerReference `group:"dbforpostgresql.azure.com" json:"owner,omitempty" kind:"FlexibleServer"`

	// The Azure Database for PostgreSQL server is created with the 3 default roles defined.
	// azure_pg_admin
	// azure_superuser
	// your server admin user
	Roles []string `json:"roles,omitempty"`

	// +kubebuilder:default={login: true}
	// The with options of the user role.
	RoleOptions *RoleOptionsSpec `json:"roleOptions,omitempty"`

	// +kubebuilder:validation:Required
	// LocalUser contains details for creating a standard (non-aad) postgresql User
	LocalUser *LocalUserSpec `json:"localUser,omitempty"`
}

// OriginalVersion returns the original API version used to create the resource.
func (userSpec *UserSpec) OriginalVersion() string {
	return GroupVersion.Version
}

// SetAzureName sets the Azure name of the resource
func (userSpec *UserSpec) SetAzureName(azureName string) { userSpec.AzureName = azureName }

type LocalUserSpec struct {
	// +kubebuilder:validation:Required
	// ServerAdminUsername is the user name of the Server administrator
	ServerAdminUsername string `json:"serverAdminUsername,omitempty"`

	// +kubebuilder:validation:Required
	// ServerAdminPassword is a reference to a secret containing the servers administrator password
	ServerAdminPassword *genruntime.SecretReference `json:"serverAdminPassword,omitempty"`

	// +kubebuilder:validation:Required
	// Password is the password to use for the user
	Password *genruntime.SecretReference `json:"password,omitempty"`
}

type RoleOptionsSpec struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=true
	// WITH LOGIN or NOLOGIN
	Login bool `json:"login,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// WITH CREATEROLE or NOCREATEROLE
	CreateRole bool `json:"createRole,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// WITH CREATEDB or NOCREATEDB
	CreateDb bool `json:"createDb,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// WITH REPLICATION or NOREPLICATION
	Replication bool `json:"replication,omitempty"`
}

type UserStatus struct {
	//Conditions: The observed state of the resource
	Conditions []conditions.Condition `json:"conditions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}
