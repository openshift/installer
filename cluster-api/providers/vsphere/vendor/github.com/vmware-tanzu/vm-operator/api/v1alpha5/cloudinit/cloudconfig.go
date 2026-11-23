// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

// +kubebuilder:object:generate=true

package cloudinit

import (
	"encoding/json"

	vmopv1a5common "github.com/vmware-tanzu/vm-operator/api/v1alpha5/common"
)

// CloudConfig is the VM Operator API subset of a Cloud-Init CloudConfig and
// contains several of the CloudConfig's frequently used modules.
type CloudConfig struct {
	// +optional

	// Timezone describes the timezone represented in /usr/share/zoneinfo.
	Timezone string `json:"timezone,omitempty"`

	// +optional

	// DefaultUserEnabled may be set to true to ensure even if the Users field
	// is not empty, the default user is still created on systems that have one
	// defined. By default, Cloud-Init ignores the default user if the
	// CloudConfig provides one or more non-default users via the Users field.
	DefaultUserEnabled bool `json:"defaultUserEnabled,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=name

	// Users allows adding/configuring one or more users on the guest.
	Users []User `json:"users,omitempty"`

	// +optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields

	// RunCmd allows running one or more commands on the guest.
	// The entries in this list can adhere to two, different formats:
	//
	// Format 1 -- a string that contains the command and its arguments, ex.
	//
	//     runcmd:
	//     - "ls -al"
	//
	// Format 2 -- a list of the command and its arguments, ex.
	//
	//     runcmd:
	//     - - echo
	//       - "Hello, world."
	RunCmd json.RawMessage `json:"runcmd,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=path

	// WriteFiles allows adding files to the guest file system.
	WriteFiles []WriteFile `json:"write_files,omitempty"`

	// +optional

	// SSHPwdAuth sets whether or not to accept password authentication.
	// In order for this config to be applied, SSH may need to be restarted.
	// On systemd systems, this restart will only happen if the SSH service has
	// already been started. On non-systemd systems, a restart will be attempted
	// regardless of the service state.
	SSHPwdAuth *bool `json:"ssh_pwauth,omitempty"`
}

// User is a CloudConfig user data structure.
type User struct {
	// +optional

	// CreateGroups is a flag that may be set to false to disable creation of
	// specified user groups.
	//
	// Defaults to true when Name is not "default".
	CreateGroups *bool `json:"create_groups,omitempty"`

	// +optional

	// ExpireData is the date on which the user's account will be disabled.
	ExpireDate *string `json:"expiredate,omitempty"`

	// +optional

	// Gecos is an optional comment about the user, usually a comma-separated
	// string of the user's real name and contact information.
	Gecos *string `json:"gecos,omitempty"`

	// +optional

	// Groups is an optional list of groups to add to the user.
	Groups []string `json:"groups,omitempty"`

	// +optional

	// HashedPasswd is a hash of the user's password that will be applied even
	// if the specified user already exists.
	HashedPasswd *vmopv1a5common.SecretKeySelector `json:"hashed_passwd,omitempty"`

	// +optional

	// Homedir is the optional home directory for the user.
	//
	// Defaults to "/home/<username>" when Name is not "default".
	Homedir *string `json:"homedir,omitempty"`

	// +optional

	// Inactive optionally represents the number of days until the user is
	// disabled.
	Inactive *int32 `json:"inactive,omitempty"`

	// +optional

	// LockPasswd disables password login.
	//
	// Defaults to true when Name is not "default".
	LockPasswd *bool `json:"lock_passwd,omitempty"`

	// Name is the user's login name.
	//
	// Please note this field may be set to the special value of "default" when
	// this User is the first element in the Users list from the CloudConfig.
	// When set to "default", all other fields from this User must be nil.
	Name string `json:"name"`

	// +optional

	// NoCreateHome prevents the creation of the home directory.
	//
	// Defaults to false when Name is not "default".
	NoCreateHome *bool `json:"no_create_home,omitempty"`

	// +optional

	// NoLogInit prevents the initialization of lastlog and faillog for the
	// user.
	//
	// Defaults to false when Name is not "default".
	NoLogInit *bool `json:"no_log_init,omitempty"`

	// +optional

	// NoUserGroup prevents the creation of the group named after the user.
	//
	// Defaults to false when Name is not "default".
	NoUserGroup *bool `json:"no_user_group,omitempty"`

	// +optional

	// Passwd is a hash of the user's password that will be applied only to
	// a newly created user. To apply a new, hashed password to an existing user
	// please use HashedPasswd instead.
	Passwd *vmopv1a5common.SecretKeySelector `json:"passwd,omitempty"`

	// +optional

	// PrimaryGroup is the primary group for the user.
	//
	// Defaults to the value of the Name field when it is not "default".
	PrimaryGroup *string `json:"primary_group,omitempty"`

	// +optional

	// SELinuxUser is the SELinux user for the user's login.
	SELinuxUser *string `json:"selinux_user,omitempty"`

	// +optional

	// Shell is the path to the user's login shell.
	//
	// Please note the default is to set no shell, which results in a
	// system-specific default being used.
	Shell *string `json:"shell,omitempty"`

	// +optional

	// SnapUser specifies an e-mail address to create the user as a Snappy user
	// through "snap create-user".
	//
	// If an Ubuntu SSO account is associated with the address, the username and
	// SSH keys will be requested from there.
	SnapUser *string `json:"snapuser,omitempty"`

	// +optional

	// SSHAuthorizedKeys is a list of SSH keys to add to the user's authorized
	// keys file.
	//
	// Please note this field may not be combined with SSHRedirectUser.
	SSHAuthorizedKeys []string `json:"ssh_authorized_keys,omitempty"`

	// +optional

	// SSHImportID is a list of SSH IDs to import for the user.
	//
	// Please note this field may not be combined with SSHRedirectUser.
	SSHImportID []string `json:"ssh_import_id,omitempty"`

	// +optional

	// SSHRedirectUser may be set to true to disable SSH logins for this user.
	//
	// Please note that when specified, all SSH keys from cloud meta-data will
	// be configured in a disabled state for this user. Any SSH login as this
	// user will timeout with a message to login instead as the default user.
	//
	// This field may not be combined with SSHAuthorizedKeys or SSHImportID.
	//
	// Defaults to false when Name is not "default".
	SSHRedirectUser *bool `json:"ssh_redirect_user,omitempty"`

	// +optional

	// Sudo is a sudo rule to apply to the user.
	//
	// When omitted, no sudo rules will be applied to the user.
	Sudo *string `json:"sudo,omitempty"`

	// +optional

	// System is an optional flag that indicates the user should be created as
	// a system user with no home directory.
	//
	// Defaults to false when Name is not "default".
	System *bool `json:"system,omitempty"`

	// +optional

	// UID is the user's ID.
	//
	// When omitted the guest will default to the next available number.
	UID *int64 `json:"uid,omitempty"`
}

// +kubebuilder:validation:Enum=b64;base64;gz;gzip;"gz+b64";"gz+base64";"gzip+b64";"gzip+base64";"text/plain"

// WriteFileEncoding specifies the encoding type of a file's content.
type WriteFileEncoding string

const (
	WriteFileEncodingFluffyB64    WriteFileEncoding = "b64"
	WriteFileEncodingFluffyBase64 WriteFileEncoding = "base64"
	WriteFileEncodingFluffyGz     WriteFileEncoding = "gz"
	WriteFileEncodingFluffyGzip   WriteFileEncoding = "gzip"
	WriteFileEncodingGzB64        WriteFileEncoding = "gz+b64"
	WriteFileEncodingGzBase64     WriteFileEncoding = "gz+base64"
	WriteFileEncodingGzipB64      WriteFileEncoding = "gzip+b64"
	WriteFileEncodingGzipBase64   WriteFileEncoding = "gzip+base64"
	WriteFileEncodingTextPlain    WriteFileEncoding = "text/plain"
)

// WriteFile is a CloudConfig write_file data structure.
type WriteFile struct {
	// +optional

	// Append specifies whether or not to append the content to an existing file
	// if the file specified by Path already exists.
	Append bool `json:"append,omitempty"`

	// +optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields

	// Content is the optional content to write to the provided Path.
	//
	// When omitted an empty file will be created or existing file will be
	// modified.
	//
	// The value for this field can adhere to two, different formats:
	//
	// Format 1 -- a string that contains the command and its arguments, ex.
	//
	//     content: Hello, world.
	//
	// Please note that format 1 supports all of the manners of specifying a
	// YAML string.
	//
	// Format 2 -- a secret reference with the name of the key that contains
	//             the content for the file, ex.
	//
	//     content:
	//       name: my-bootstrap-secret
	//       key: my-file-content
	Content json.RawMessage `json:"content,omitempty"`

	// +optional

	// Defer indicates to defer writing the file until Cloud-Init's "final"
	// stage, after users are created and packages are installed.
	Defer bool `json:"defer,omitempty"`

	// +optional
	// +kubebuilder:default="text/plain"

	// Encoding is an optional encoding type of the content.
	Encoding WriteFileEncoding `json:"encoding,omitempty"`

	// +optional
	// +kubebuilder:default="root:root"

	// Owner is an optional "owner:group" to chown the file.
	Owner string `json:"owner,omitempty"`

	// Path is the path of the file to which the content is decoded and written.
	Path string `json:"path"`

	// +optional
	// +kubebuilder:default="0644"

	// Permissions an optional set of file permissions to set.
	//
	// Please note the permissions should be specified as an octal string, ex.
	// "0###".
	//
	// When omitted the guest will default this value to "0644".
	Permissions string `json:"permissions,omitempty"`
}
