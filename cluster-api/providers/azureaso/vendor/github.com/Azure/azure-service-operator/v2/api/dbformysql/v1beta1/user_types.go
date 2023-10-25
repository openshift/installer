// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package v1beta1

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	v1 "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

// +kubebuilder:rbac:groups=dbformysql.azure.com,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbformysql.azure.com,resources={users/status,users/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// User is a MySQL user
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

// +kubebuilder:webhook:path=/mutate-dbformysql-azure-com-v1beta1-user,mutating=true,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=dbformysql.azure.com,resources=users,verbs=create;update,versions=v1beta1,name=default.v1beta1.users.dbformysql.azure.com,admissionReviewVersions=v1

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

// +kubebuilder:webhook:path=/validate-dbformysql-azure-com-v1beta1-user,mutating=false,sideEffects=None,matchPolicy=Exact,failurePolicy=fail,groups=dbformysql.azure.com,resources=users,verbs=create;update,versions=v1beta1,name=validate.v1beta1.users.dbformysql.azure.com,admissionReviewVersions=v1beta1

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

var _ conversion.Convertible = &User{}

// ConvertFrom populates our ConfigurationStore from the provided hub ConfigurationStore
func (user *User) ConvertFrom(hub conversion.Hub) error {
	source, ok := hub.(*v1.User)
	if !ok {
		return fmt.Errorf("expected dbformysql/v1/User but received %T instead", hub)
	}

	return user.AssignProperties_From_User(source)
}

// ConvertTo populates the provided hub ConfigurationStore from our ConfigurationStore
func (user *User) ConvertTo(hub conversion.Hub) error {
	destination, ok := hub.(*v1.User)
	if !ok {
		return fmt.Errorf("expected dbformysql/v1/User but received %T instead", hub)
	}

	return user.AssignProperties_To_User(destination)
}

// AssignProperties_To_User populates the provided destination User from our User
func (user *User) AssignProperties_To_User(destination *v1.User) error {
	// ObjectMeta
	destination.ObjectMeta = *user.ObjectMeta.DeepCopy()

	// Spec
	var spec v1.UserSpec
	err := user.Spec.AssignProperties_To_UserSpec(&spec)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_To_User_Spec() to populate field Spec")
	}
	destination.Spec = spec

	// Status
	var status v1.UserStatus
	err = user.Status.AssignProperties_To_UserStatus(&status)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_To_User_Status() to populate field Status")
	}
	destination.Status = status

	// No error
	return nil
}

func (user *User) AssignProperties_From_User(source *v1.User) error {
	// ObjectMeta
	user.ObjectMeta = *source.ObjectMeta.DeepCopy()

	// Spec
	var spec UserSpec
	err := spec.AssignProperties_From_UserSpec(&source.Spec)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_From_UserSpec() to populate field Spec")
	}
	user.Spec = spec

	// Status
	var status UserStatus
	err = status.AssignProperties_From_UserStatus(&source.Status)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_From_UserStatus() to populate field Status")
	}
	user.Status = status

	// No error
	return nil
}

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
	//reference to a dbformysql.azure.com/FlexibleServer resource
	Owner *genruntime.KubernetesOwnerReference `group:"dbformysql.azure.com" json:"owner,omitempty" kind:"FlexibleServer"`

	// Hostname is the host the user will connect from. If omitted, the default is to allow connection from any hostname.
	Hostname string `json:"hostname,omitempty"`

	// The server-level roles assigned to the user.
	// Privileges include the following: RELOAD, PROCESS, SHOW
	// DATABASES, REPLICATION SLAVE, REPLICATION CLIENT, CREATE USER
	Privileges []string `json:"privileges,omitempty"`

	// The database-level roles assigned to the user (keyed by
	// database name). Privileges include the following: SELECT,
	// INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, INDEX,
	// ALTER, CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, CREATE
	// VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, EVENT, TRIGGER
	DatabasePrivileges map[string][]string `json:"databasePrivileges,omitempty"`

	// TODO: Note this is required right now but will move to be optional in the future when we have AAD support
	// +kubebuilder:validation:Required
	// LocalUser contains details for creating a standard (non-aad) MySQL User
	LocalUser *LocalUserSpec `json:"localUser,omitempty"`
}

func (userSpec *UserSpec) AssignProperties_To_UserSpec(destination *v1.UserSpec) error {
	// AzureName
	destination.AzureName = userSpec.AzureName

	// Owner
	if userSpec.Owner != nil {
		owner := userSpec.Owner.Copy()
		destination.Owner = &owner
	} else {
		destination.Owner = nil
	}

	// Hostname
	destination.Hostname = userSpec.Hostname

	// Privileges
	destination.Privileges = genruntime.CloneSliceOfString(userSpec.Privileges)

	// DatabasePrivileges
	if userSpec.DatabasePrivileges != nil {
		result := make(map[string][]string, len(userSpec.DatabasePrivileges))
		for k, v := range userSpec.DatabasePrivileges {
			result[k] = genruntime.CloneSliceOfString(v)
		}

		destination.DatabasePrivileges = result
	}

	// LocalUser
	if userSpec.LocalUser != nil {
		destination.LocalUser = &v1.LocalUserSpec{
			ServerAdminUsername: userSpec.LocalUser.ServerAdminUsername,
			Password:            userSpec.LocalUser.Password.DeepCopy(),
			ServerAdminPassword: userSpec.LocalUser.ServerAdminPassword.DeepCopy(),
		}
	}

	// No error
	return nil
}

func (userSpec *UserSpec) AssignProperties_From_UserSpec(source *v1.UserSpec) error {
	// AzureName
	userSpec.AzureName = source.AzureName

	// Owner
	if source.Owner != nil {
		owner := source.Owner.Copy()
		userSpec.Owner = &owner
	} else {
		userSpec.Owner = nil
	}

	// Hostname
	userSpec.Hostname = source.Hostname

	// Privileges
	userSpec.Privileges = genruntime.CloneSliceOfString(source.Privileges)

	// DatabasePrivileges
	if source.DatabasePrivileges != nil {
		result := make(map[string][]string, len(source.DatabasePrivileges))
		for k, v := range source.DatabasePrivileges {
			result[k] = genruntime.CloneSliceOfString(v)
		}

		userSpec.DatabasePrivileges = result
	}

	// LocalUser
	if source.LocalUser != nil {
		userSpec.LocalUser = &LocalUserSpec{
			ServerAdminUsername: source.LocalUser.ServerAdminUsername,
			Password:            source.LocalUser.Password.DeepCopy(),
			ServerAdminPassword: source.LocalUser.ServerAdminPassword.DeepCopy(),
		}
	}

	// No error
	return nil
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

type UserStatus struct {
	//Conditions: The observed state of the resource
	Conditions []conditions.Condition `json:"conditions,omitempty"`
}

func (userSpec *UserStatus) AssignProperties_To_UserStatus(destination *v1.UserStatus) error {
	destination.Conditions = genruntime.CloneSliceOfCondition(userSpec.Conditions)

	// No error
	return nil
}

func (userSpec *UserStatus) AssignProperties_From_UserStatus(source *v1.UserStatus) error {
	userSpec.Conditions = genruntime.CloneSliceOfCondition(source.Conditions)

	// No error
	return nil
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}
