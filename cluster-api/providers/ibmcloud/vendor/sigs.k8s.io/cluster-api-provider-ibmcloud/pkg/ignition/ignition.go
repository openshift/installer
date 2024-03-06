/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ignition

// CaReference holds the CaReference specific information.
type CaReference struct {
	HTTPHeaders  HTTPHeaders  `json:"httpHeaders,omitempty"`
	Source       string       `json:"source"`
	Verification Verification `json:"verification,omitempty"`
}

// Config holds the Config specific information.
type Config struct {
	Ignition Ignition `json:"ignition"`
	Networkd Networkd `json:"networkd,omitempty"`
	Passwd   Passwd   `json:"passwd,omitempty"`
	Storage  Storage  `json:"storage,omitempty"`
	Systemd  Systemd  `json:"systemd,omitempty"`
}

// ConfigReference holds the ConfigReference specific information.
type ConfigReference struct {
	HTTPHeaders  HTTPHeaders  `json:"httpHeaders,omitempty"`
	Source       string       `json:"source"`
	Verification Verification `json:"verification,omitempty"`
}

// Create holds the Create specific information.
type Create struct {
	Force   bool           `json:"force,omitempty"`
	Options []CreateOption `json:"options,omitempty"`
}

// CreateOption holds the CreateOption specific information.
type CreateOption string

// Device holds the Device specific information.
type Device string

// Directory holds the Directory specific information.
type Directory struct {
	Node
	DirectoryEmbedded1
}

// DirectoryEmbedded1 holds the DirectoryEmbedded1 specific information.
type DirectoryEmbedded1 struct {
	Mode *int `json:"mode,omitempty"`
}

// Disk holds the Disk specific information.
type Disk struct {
	Device     string      `json:"device"`
	Partitions []Partition `json:"partitions,omitempty"`
	WipeTable  bool        `json:"wipeTable,omitempty"`
}

// File holds the File specific information.
type File struct {
	Node
	FileEmbedded1
}

// FileContents holds the FileContents specific information.
type FileContents struct {
	Compression  string       `json:"compression,omitempty"`
	HTTPHeaders  HTTPHeaders  `json:"httpHeaders,omitempty"`
	Source       string       `json:"source,omitempty"`
	Verification Verification `json:"verification,omitempty"`
}

// FileEmbedded1 holds the FileEmbedded1 specific information.
type FileEmbedded1 struct {
	Append   bool         `json:"append,omitempty"`
	Contents FileContents `json:"contents,omitempty"`
	Mode     *int         `json:"mode,omitempty"`
}

// Filesystem holds the Filesystem specific information.
type Filesystem struct {
	Mount *Mount  `json:"mount,omitempty"`
	Name  string  `json:"name,omitempty"`
	Path  *string `json:"path,omitempty"`
}

// Group holds the Group specific information.
type Group string

// HTTPHeader holds the HTTPHeader specific information.
type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// HTTPHeaders holds the HTTPHeaders specific information.
type HTTPHeaders []HTTPHeader

// Ignition holds the Ignition specific information.
type Ignition struct {
	Config   IgnitionConfig `json:"config,omitempty"`
	Proxy    Proxy          `json:"proxy,omitempty"`
	Security Security       `json:"security,omitempty"`
	Timeouts Timeouts       `json:"timeouts,omitempty"`
	Version  string         `json:"version,omitempty"`
}

// IgnitionConfig holds the IgnitionConfig specific information.
type IgnitionConfig struct { //nolint:revive
	Append  []ConfigReference `json:"append,omitempty"`
	Replace *ConfigReference  `json:"replace,omitempty"`
}

// Link holds the Link specific information.
type Link struct {
	Node
	LinkEmbedded1
}

// LinkEmbedded1 holds the LinkEmbedded1 specific information.
type LinkEmbedded1 struct {
	Hard   bool   `json:"hard,omitempty"`
	Target string `json:"target"`
}

// Mount holds the Mount specific information.
type Mount struct {
	Create         *Create       `json:"create,omitempty"`
	Device         string        `json:"device"`
	Format         string        `json:"format"`
	Label          *string       `json:"label,omitempty"`
	Options        []MountOption `json:"options,omitempty"`
	UUID           *string       `json:"uuid,omitempty"`
	WipeFilesystem bool          `json:"wipeFilesystem,omitempty"`
}

// MountOption holds the MountOption specific information.
type MountOption string

// Networkd holds the Networkd specific information.
type Networkd struct {
	Units []Networkdunit `json:"units,omitempty"`
}

// NetworkdDropin holds the NetworkdDropin specific information.
type NetworkdDropin struct {
	Contents string `json:"contents,omitempty"`
	Name     string `json:"name"`
}

// Networkdunit holds the Networkdunit specific information.
type Networkdunit struct {
	Contents string           `json:"contents,omitempty"`
	Dropins  []NetworkdDropin `json:"dropins,omitempty"`
	Name     string           `json:"name"`
}

// NoProxyItem holds the NoProxyItem specific information.
type NoProxyItem string

// Node holds the Node specific information.
type Node struct {
	Filesystem string     `json:"filesystem"`
	Group      *NodeGroup `json:"group,omitempty"`
	Overwrite  *bool      `json:"overwrite,omitempty"`
	Path       string     `json:"path"`
	User       *NodeUser  `json:"user,omitempty"`
}

// NodeGroup holds the NodeGroup specific information.
type NodeGroup struct {
	ID   *int   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NodeUser holds the NodeUser specific information.
type NodeUser struct {
	ID   *int   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Partition holds the Partition specific information.
type Partition struct {
	GUID               string  `json:"guid,omitempty"`
	Label              *string `json:"label,omitempty"`
	Number             int     `json:"number,omitempty"`
	ShouldExist        *bool   `json:"shouldExist,omitempty"`
	Size               *int    `json:"size,omitempty"`
	SizeMiB            *int    `json:"sizeMiB,omitempty"`
	Start              *int    `json:"start,omitempty"`
	StartMiB           *int    `json:"startMiB,omitempty"`
	TypeGUID           string  `json:"typeGuid,omitempty"`
	WipePartitionEntry bool    `json:"wipePartitionEntry,omitempty"`
}

// Passwd holds the Passwd specific information.
type Passwd struct {
	Groups []PasswdGroup `json:"groups,omitempty"`
	Users  []PasswdUser  `json:"users,omitempty"`
}

// PasswdGroup holds the PasswdGroup specific information.
type PasswdGroup struct {
	Gid          *int   `json:"gid,omitempty"` //nolint:stylecheck
	Name         string `json:"name"`
	PasswordHash string `json:"passwordHash,omitempty"`
	System       bool   `json:"system,omitempty"`
}

// PasswdUser holds the PasswdUser specific information.
type PasswdUser struct {
	Create            *Usercreate        `json:"create,omitempty"`
	Gecos             string             `json:"gecos,omitempty"`
	Groups            []Group            `json:"groups,omitempty"`
	HomeDir           string             `json:"homeDir,omitempty"`
	Name              string             `json:"name"`
	NoCreateHome      bool               `json:"noCreateHome,omitempty"`
	NoLogInit         bool               `json:"noLogInit,omitempty"`
	NoUserGroup       bool               `json:"noUserGroup,omitempty"`
	PasswordHash      *string            `json:"passwordHash,omitempty"`
	PrimaryGroup      string             `json:"primaryGroup,omitempty"`
	SSHAuthorizedKeys []SSHAuthorizedKey `json:"sshAuthorizedKeys,omitempty"`
	Shell             string             `json:"shell,omitempty"`
	System            bool               `json:"system,omitempty"`
	UID               *int               `json:"uid,omitempty"`
}

// Proxy holds the Proxy specific information.
type Proxy struct {
	HTTPProxy  string        `json:"httpProxy,omitempty"`
	HTTPSProxy string        `json:"httpsProxy,omitempty"`
	NoProxy    []NoProxyItem `json:"noProxy,omitempty"`
}

// Raid holds the Raid specific information.
type Raid struct {
	Devices []Device     `json:"devices"`
	Level   string       `json:"level"`
	Name    string       `json:"name"`
	Options []RaidOption `json:"options,omitempty"`
	Spares  int          `json:"spares,omitempty"`
}

// RaidOption holds the RaidOption specific information.
type RaidOption string

// SSHAuthorizedKey holds the SSHAuthorizedKey specific information.
type SSHAuthorizedKey string

// Security holds the Security specific information.
type Security struct {
	TLS TLS `json:"tls,omitempty"`
}

// Storage holds the Storage specific information.
type Storage struct {
	Directories []Directory  `json:"directories,omitempty"`
	Disks       []Disk       `json:"disks,omitempty"`
	Files       []File       `json:"files,omitempty"`
	Filesystems []Filesystem `json:"filesystems,omitempty"`
	Links       []Link       `json:"links,omitempty"`
	Raid        []Raid       `json:"raid,omitempty"`
}

// Systemd holds the Systemd specific information.
type Systemd struct {
	Units []Unit `json:"units,omitempty"`
}

// SystemdDropin holds the SystemdDropin specific information.
type SystemdDropin struct {
	Contents string `json:"contents,omitempty"`
	Name     string `json:"name"`
}

// TLS holds the TLS specific information.
type TLS struct {
	CertificateAuthorities []CaReference `json:"certificateAuthorities,omitempty"`
}

// Timeouts holds the Timeouts specific information.
type Timeouts struct {
	HTTPResponseHeaders *int `json:"httpResponseHeaders,omitempty"`
	HTTPTotal           *int `json:"httpTotal,omitempty"`
}

// Unit holds the Unit specific information.
type Unit struct {
	Contents string          `json:"contents,omitempty"`
	Dropins  []SystemdDropin `json:"dropins,omitempty"`
	Enable   bool            `json:"enable,omitempty"`
	Enabled  *bool           `json:"enabled,omitempty"`
	Mask     bool            `json:"mask,omitempty"`
	Name     string          `json:"name"`
}

// Usercreate holds the Usercreate specific information.
type Usercreate struct {
	Gecos        string            `json:"gecos,omitempty"`
	Groups       []UsercreateGroup `json:"groups,omitempty"`
	HomeDir      string            `json:"homeDir,omitempty"`
	NoCreateHome bool              `json:"noCreateHome,omitempty"`
	NoLogInit    bool              `json:"noLogInit,omitempty"`
	NoUserGroup  bool              `json:"noUserGroup,omitempty"`
	PrimaryGroup string            `json:"primaryGroup,omitempty"`
	Shell        string            `json:"shell,omitempty"`
	System       bool              `json:"system,omitempty"`
	UID          *int              `json:"uid,omitempty"`
}

// UsercreateGroup holds the UsercreateGroup specific information.
type UsercreateGroup string

// Verification holds the Verification specific information.
type Verification struct {
	Hash *string `json:"hash,omitempty"`
}
