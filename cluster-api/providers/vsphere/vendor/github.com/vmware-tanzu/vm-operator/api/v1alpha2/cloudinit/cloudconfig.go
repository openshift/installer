// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

// +kubebuilder:object:generate=true

package cloudinit

import (
	"encoding/json"

	vmopv1a2common "github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
)

// CloudConfig is the VM Operator API subset of a Cloud-Init CloudConfig and
// contains several of the CloudConfig's frequently used modules.
type CloudConfig struct {
	// Timezone describes the timezone represented in /usr/share/zoneinfo.
	//
	// +optional
	Timezone string `json:"timezone,omitempty"`

	// DefaultUserEnabled may be set to true to ensure even if the Users field
	// is not empty, the default user is still created on systems that have one
	// defined. By default, Cloud-Init ignores the default user if the
	// CloudConfig provides one or more non-default users via the Users field.
	//
	// +optional
	DefaultUserEnabled bool `json:"defaultUserEnabled,omitempty"`

	// Users allows adding/configuring one or more users on the guest.
	//
	// +optional
	// +listType=map
	// +listMapKey=name
	Users []User `json:"users,omitempty"`

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
	//
	// +optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	RunCmd json.RawMessage `json:"runcmd,omitempty"`

	// WriteFiles
	//
	// +optional
	// +listType=map
	// +listMapKey=path
	WriteFiles []WriteFile `json:"write_files,omitempty"`

	// SSHPwdAuth sets whether or not to accept password authentication. ``true`` will enable password
	// auth. ``false`` will disable. Default: leave the value unchanged. In order for this
	// config to be applied, SSH may need to be restarted. On systemd systems, this restart will
	// only happen if the SSH service has already been started. On non-systemd systems, a
	// restart will be attempted regardless of the service state.
	//
	// +optional
	SSHPwdAuth *bool `json:"ssh_pwauth,omitempty"`
}

// User is a CloudConfig user data structure.
type User struct {
	// CreateGroups is a flag that may be set to false to disable creation of
	// specified user groups.
	//
	// Defaults to true when Name is not "default".
	//
	// +optional
	CreateGroups *bool `json:"create_groups,omitempty"`

	// ExpireData is the date on which the user's account will be disabled.
	//
	// +optional
	ExpireDate *string `json:"expiredate,omitempty"`

	// Gecos is an optional comment about the user, usually a comma-separated
	// string of the user's real name and contact information.
	//
	// +optional
	Gecos *string `json:"gecos,omitempty"`

	// Groups is an optional list of groups to add to the user.
	//
	// +optional
	Groups []string `json:"groups,omitempty"`

	// HashedPasswd is a hash of the user's password that will be applied even
	// if the specified user already exists.
	//
	// +optional
	HashedPasswd *vmopv1a2common.SecretKeySelector `json:"hashed_passwd,omitempty"`

	// Homedir is the optional home directory for the user.
	//
	// Defaults to "/home/<username>" when Name is not "default".
	//
	// +optional
	Homedir *string `json:"homedir,omitempty"`

	// Inactive optionally represents the number of days until the user is
	// disabled.
	//
	// +optional
	Inactive *int32 `json:"inactive,omitempty"`

	// LockPasswd disables password login.
	//
	// Defaults to true when Name is not "default".
	//
	// +optional
	LockPasswd *bool `json:"lock_passwd,omitempty"`

	// Name is the user's login name.
	//
	// Please note this field may be set to the special value of "default" when
	// this User is the first element in the Users list from the CloudConfig.
	// When set to "default", all other fields from this User must be nil.
	Name string `json:"name"`

	// NoCreateHome prevents the creation of the home directory.
	//
	// Defaults to false when Name is not "default".
	//
	// +optional
	NoCreateHome *bool `json:"no_create_home,omitempty"`

	// NoLogInit prevents the initialization of lastlog and faillog for the
	// user.
	//
	// Defaults to false when Name is not "default".
	//
	// +optional
	NoLogInit *bool `json:"no_log_init,omitempty"`

	// NoUserGroup prevents the creation of the group named after the user.
	//
	// Defaults to false when Name is not "default".
	//
	// +optional
	NoUserGroup *bool `json:"no_user_group,omitempty"`

	// Passwd is a hash of the user's password that will be applied only to
	// a newly created user. To apply a new, hashed password to an existing user
	// please use HashedPasswd instead.
	//
	// +optional
	Passwd *vmopv1a2common.SecretKeySelector `json:"passwd,omitempty"`

	// PrimaryGroup is the primary group for the user.
	//
	// Defaults to the value of the Name field when it is not "default".
	//
	// +optional
	PrimaryGroup *string `json:"primary_group,omitempty"`

	// SELinuxUser is the SELinux user for the user's login.
	//
	// +optional
	SELinuxUser *string `json:"selinux_user,omitempty"`

	// Shell is the path to the user's login shell.
	//
	// Please note the default is to set no shell, which results in a
	// system-specific default being used.
	//
	// +optional
	Shell *string `json:"shell,omitempty"`

	// SnapUser specifies an e-mail address to create the user as a Snappy user
	// through "snap create-user".
	//
	// If an Ubuntu SSO account is associated with the address, the username and
	// SSH keys will be requested from there.
	//
	// +optional
	SnapUser *string `json:"snapuser,omitempty"`

	// SSHAuthorizedKeys is a list of SSH keys to add to the user's authorized
	// keys file.
	//
	// Please note this field may not be combined with SSHRedirectUser.
	//
	// +optional
	SSHAuthorizedKeys []string `json:"ssh_authorized_keys,omitempty"`

	// SSHImportID is a list of SSH IDs to import for the user.
	//
	// Please note this field may not be combined with SSHRedirectUser.
	//
	// +optional
	SSHImportID []string `json:"ssh_import_id,omitempty"`

	// SSHRedirectUser may be set to true to disable SSH logins for this user.
	//
	// Please note that when specified, all SSH keys from cloud meta-data will
	// be configured in a disabled state for this user. Any SSH login as this
	// user will timeout with a message to login instead as the default user.
	//
	// This field may not be combined with SSHAuthorizedKeys or SSHImportID.
	//
	// Defaults to false when Name is not "default".
	//
	// +optional
	SSHRedirectUser *bool `json:"ssh_redirect_user,omitempty"`

	// Sudo is a sudo rule to apply to the user.
	//
	// When omitted, no sudo rules will be applied to the user.
	//
	// +optional
	Sudo *string `json:"sudo,omitempty"`

	// System is an optional flag that indicates the user should be created as
	// a system user with no home directory.
	//
	// Defaults to false when Name is not "default".
	//
	// +optional
	System *bool `json:"system,omitempty"`

	// UID is the user's ID.
	//
	// When omitted the guest will default to the next available number.
	//
	// +optional
	UID *int64 `json:"uid,omitempty"`
}

// WriteFileEncoding specifies the encoding type of a file's content.
//
// +kubebuilder:validation:Enum=b64;base64;gz;gzip;"gz+b64";"gz+base64";"gzip+b64";"gzip+base64";"text/plain"
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

// WriteFile is a CloudConfig
// write_file data structure.
type WriteFile struct {
	// Append specifies whether or not to append the content to an existing file
	// if the file specified by Path already exists.
	//
	// +optional
	Append bool `json:"append,omitempty"`

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
	//
	// +optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	Content json.RawMessage `json:"content,omitempty"`

	// Defer indicates to defer writing the file until Cloud-Init's "final"
	// stage, after users are created and packages are installed.
	//
	// +optional
	Defer bool `json:"defer,omitempty"`

	// Encoding is an optional encoding type of the content.
	//
	// +optional
	// +kubebuilder:default="text/plain"
	Encoding WriteFileEncoding `json:"encoding,omitempty"`

	// Owner is an optional "owner:group" to chown the file.
	//
	// +optional
	// +kubebuilder:default="root:root"
	Owner string `json:"owner,omitempty"`

	// Path is the path of the file to which the content is decoded and written.
	Path string `json:"path"`

	// Permissions an optional set of file permissions to set.
	//
	// Please note the permissions should be specified as an octal string, ex.
	// "0###".
	//
	// When omitted the guest will default this value to "0644".
	//
	// +optional
	// +kubebuilder:default="0644"
	Permissions string `json:"permissions,omitempty"`
}
