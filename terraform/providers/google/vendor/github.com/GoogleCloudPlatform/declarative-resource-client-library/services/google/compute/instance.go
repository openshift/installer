// Copyright 2024 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package compute

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Instance struct {
	CanIPForward           *bool                           `json:"canIPForward"`
	CpuPlatform            *string                         `json:"cpuPlatform"`
	CreationTimestamp      *string                         `json:"creationTimestamp"`
	DeletionProtection     *bool                           `json:"deletionProtection"`
	Description            *string                         `json:"description"`
	Disks                  []InstanceDisks                 `json:"disks"`
	GuestAccelerators      []InstanceGuestAccelerators     `json:"guestAccelerators"`
	Hostname               *string                         `json:"hostname"`
	Id                     *string                         `json:"id"`
	Labels                 map[string]string               `json:"labels"`
	Metadata               map[string]string               `json:"metadata"`
	MachineType            *string                         `json:"machineType"`
	MinCpuPlatform         *string                         `json:"minCpuPlatform"`
	Name                   *string                         `json:"name"`
	NetworkInterfaces      []InstanceNetworkInterfaces     `json:"networkInterfaces"`
	Scheduling             *InstanceScheduling             `json:"scheduling"`
	ServiceAccounts        []InstanceServiceAccounts       `json:"serviceAccounts"`
	ShieldedInstanceConfig *InstanceShieldedInstanceConfig `json:"shieldedInstanceConfig"`
	Status                 *InstanceStatusEnum             `json:"status"`
	StatusMessage          *string                         `json:"statusMessage"`
	Tags                   []string                        `json:"tags"`
	Zone                   *string                         `json:"zone"`
	Project                *string                         `json:"project"`
	SelfLink               *string                         `json:"selfLink"`
}

func (r *Instance) String() string {
	return dcl.SprintResource(r)
}

// The enum InstanceDisksInterfaceEnum.
type InstanceDisksInterfaceEnum string

// InstanceDisksInterfaceEnumRef returns a *InstanceDisksInterfaceEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceDisksInterfaceEnumRef(s string) *InstanceDisksInterfaceEnum {
	v := InstanceDisksInterfaceEnum(s)
	return &v
}

func (v InstanceDisksInterfaceEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCSI", "NVME"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceDisksInterfaceEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceDisksModeEnum.
type InstanceDisksModeEnum string

// InstanceDisksModeEnumRef returns a *InstanceDisksModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceDisksModeEnumRef(s string) *InstanceDisksModeEnum {
	v := InstanceDisksModeEnum(s)
	return &v
}

func (v InstanceDisksModeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"READ_WRITE", "READ_ONLY"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceDisksModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceDisksTypeEnum.
type InstanceDisksTypeEnum string

// InstanceDisksTypeEnumRef returns a *InstanceDisksTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceDisksTypeEnumRef(s string) *InstanceDisksTypeEnum {
	v := InstanceDisksTypeEnum(s)
	return &v
}

func (v InstanceDisksTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SCRATCH", "PERSISTENT"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceDisksTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceNetworkInterfacesAccessConfigsNetworkTierEnum.
type InstanceNetworkInterfacesAccessConfigsNetworkTierEnum string

// InstanceNetworkInterfacesAccessConfigsNetworkTierEnumRef returns a *InstanceNetworkInterfacesAccessConfigsNetworkTierEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceNetworkInterfacesAccessConfigsNetworkTierEnumRef(s string) *InstanceNetworkInterfacesAccessConfigsNetworkTierEnum {
	v := InstanceNetworkInterfacesAccessConfigsNetworkTierEnum(s)
	return &v
}

func (v InstanceNetworkInterfacesAccessConfigsNetworkTierEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PREMIUM", "STANDARD"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceNetworkInterfacesAccessConfigsNetworkTierEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceNetworkInterfacesAccessConfigsTypeEnum.
type InstanceNetworkInterfacesAccessConfigsTypeEnum string

// InstanceNetworkInterfacesAccessConfigsTypeEnumRef returns a *InstanceNetworkInterfacesAccessConfigsTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceNetworkInterfacesAccessConfigsTypeEnumRef(s string) *InstanceNetworkInterfacesAccessConfigsTypeEnum {
	v := InstanceNetworkInterfacesAccessConfigsTypeEnum(s)
	return &v
}

func (v InstanceNetworkInterfacesAccessConfigsTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ONE_TO_ONE_NAT"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceNetworkInterfacesAccessConfigsTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum.
type InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum string

// InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnumRef returns a *InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnumRef(s string) *InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum {
	v := InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum(s)
	return &v
}

func (v InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PREMIUM", "STANDARD"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum.
type InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum string

// InstanceNetworkInterfacesIPv6AccessConfigsTypeEnumRef returns a *InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceNetworkInterfacesIPv6AccessConfigsTypeEnumRef(s string) *InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum {
	v := InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum(s)
	return &v
}

func (v InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ONE_TO_ONE_NAT"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum InstanceStatusEnum.
type InstanceStatusEnum string

// InstanceStatusEnumRef returns a *InstanceStatusEnum with the value of string s
// If the empty string is provided, nil is returned.
func InstanceStatusEnumRef(s string) *InstanceStatusEnum {
	v := InstanceStatusEnum(s)
	return &v
}

func (v InstanceStatusEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PROVISIONING", "STAGING", "RUNNING", "STOPPING", "SUSPENDING", "SUSPENDED", "TERMINATED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "InstanceStatusEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type InstanceDisks struct {
	empty             bool                            `json:"-"`
	AutoDelete        *bool                           `json:"autoDelete"`
	Boot              *bool                           `json:"boot"`
	DeviceName        *string                         `json:"deviceName"`
	DiskEncryptionKey *InstanceDisksDiskEncryptionKey `json:"diskEncryptionKey"`
	Index             *int64                          `json:"index"`
	InitializeParams  *InstanceDisksInitializeParams  `json:"initializeParams"`
	Interface         *InstanceDisksInterfaceEnum     `json:"interface"`
	Mode              *InstanceDisksModeEnum          `json:"mode"`
	Source            *string                         `json:"source"`
	Type              *InstanceDisksTypeEnum          `json:"type"`
}

type jsonInstanceDisks InstanceDisks

func (r *InstanceDisks) UnmarshalJSON(data []byte) error {
	var res jsonInstanceDisks
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceDisks
	} else {

		r.AutoDelete = res.AutoDelete

		r.Boot = res.Boot

		r.DeviceName = res.DeviceName

		r.DiskEncryptionKey = res.DiskEncryptionKey

		r.Index = res.Index

		r.InitializeParams = res.InitializeParams

		r.Interface = res.Interface

		r.Mode = res.Mode

		r.Source = res.Source

		r.Type = res.Type

	}
	return nil
}

// This object is used to assert a desired state where this InstanceDisks is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceDisks *InstanceDisks = &InstanceDisks{empty: true}

func (r *InstanceDisks) Empty() bool {
	return r.empty
}

func (r *InstanceDisks) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceDisks) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceDisksDiskEncryptionKey struct {
	empty           bool    `json:"-"`
	RawKey          *string `json:"rawKey"`
	RsaEncryptedKey *string `json:"rsaEncryptedKey"`
	Sha256          *string `json:"sha256"`
}

type jsonInstanceDisksDiskEncryptionKey InstanceDisksDiskEncryptionKey

func (r *InstanceDisksDiskEncryptionKey) UnmarshalJSON(data []byte) error {
	var res jsonInstanceDisksDiskEncryptionKey
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceDisksDiskEncryptionKey
	} else {

		r.RawKey = res.RawKey

		r.RsaEncryptedKey = res.RsaEncryptedKey

		r.Sha256 = res.Sha256

	}
	return nil
}

// This object is used to assert a desired state where this InstanceDisksDiskEncryptionKey is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceDisksDiskEncryptionKey *InstanceDisksDiskEncryptionKey = &InstanceDisksDiskEncryptionKey{empty: true}

func (r *InstanceDisksDiskEncryptionKey) Empty() bool {
	return r.empty
}

func (r *InstanceDisksDiskEncryptionKey) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceDisksDiskEncryptionKey) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceDisksInitializeParams struct {
	empty                    bool                                                   `json:"-"`
	DiskName                 *string                                                `json:"diskName"`
	DiskSizeGb               *int64                                                 `json:"diskSizeGb"`
	DiskType                 *string                                                `json:"diskType"`
	SourceImage              *string                                                `json:"sourceImage"`
	SourceImageEncryptionKey *InstanceDisksInitializeParamsSourceImageEncryptionKey `json:"sourceImageEncryptionKey"`
}

type jsonInstanceDisksInitializeParams InstanceDisksInitializeParams

func (r *InstanceDisksInitializeParams) UnmarshalJSON(data []byte) error {
	var res jsonInstanceDisksInitializeParams
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceDisksInitializeParams
	} else {

		r.DiskName = res.DiskName

		r.DiskSizeGb = res.DiskSizeGb

		r.DiskType = res.DiskType

		r.SourceImage = res.SourceImage

		r.SourceImageEncryptionKey = res.SourceImageEncryptionKey

	}
	return nil
}

// This object is used to assert a desired state where this InstanceDisksInitializeParams is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceDisksInitializeParams *InstanceDisksInitializeParams = &InstanceDisksInitializeParams{empty: true}

func (r *InstanceDisksInitializeParams) Empty() bool {
	return r.empty
}

func (r *InstanceDisksInitializeParams) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceDisksInitializeParams) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceDisksInitializeParamsSourceImageEncryptionKey struct {
	empty  bool    `json:"-"`
	RawKey *string `json:"rawKey"`
	Sha256 *string `json:"sha256"`
}

type jsonInstanceDisksInitializeParamsSourceImageEncryptionKey InstanceDisksInitializeParamsSourceImageEncryptionKey

func (r *InstanceDisksInitializeParamsSourceImageEncryptionKey) UnmarshalJSON(data []byte) error {
	var res jsonInstanceDisksInitializeParamsSourceImageEncryptionKey
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceDisksInitializeParamsSourceImageEncryptionKey
	} else {

		r.RawKey = res.RawKey

		r.Sha256 = res.Sha256

	}
	return nil
}

// This object is used to assert a desired state where this InstanceDisksInitializeParamsSourceImageEncryptionKey is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceDisksInitializeParamsSourceImageEncryptionKey *InstanceDisksInitializeParamsSourceImageEncryptionKey = &InstanceDisksInitializeParamsSourceImageEncryptionKey{empty: true}

func (r *InstanceDisksInitializeParamsSourceImageEncryptionKey) Empty() bool {
	return r.empty
}

func (r *InstanceDisksInitializeParamsSourceImageEncryptionKey) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceDisksInitializeParamsSourceImageEncryptionKey) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceGuestAccelerators struct {
	empty            bool    `json:"-"`
	AcceleratorCount *int64  `json:"acceleratorCount"`
	AcceleratorType  *string `json:"acceleratorType"`
}

type jsonInstanceGuestAccelerators InstanceGuestAccelerators

func (r *InstanceGuestAccelerators) UnmarshalJSON(data []byte) error {
	var res jsonInstanceGuestAccelerators
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceGuestAccelerators
	} else {

		r.AcceleratorCount = res.AcceleratorCount

		r.AcceleratorType = res.AcceleratorType

	}
	return nil
}

// This object is used to assert a desired state where this InstanceGuestAccelerators is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceGuestAccelerators *InstanceGuestAccelerators = &InstanceGuestAccelerators{empty: true}

func (r *InstanceGuestAccelerators) Empty() bool {
	return r.empty
}

func (r *InstanceGuestAccelerators) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceGuestAccelerators) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceNetworkInterfaces struct {
	empty             bool                                         `json:"-"`
	AccessConfigs     []InstanceNetworkInterfacesAccessConfigs     `json:"accessConfigs"`
	IPv6AccessConfigs []InstanceNetworkInterfacesIPv6AccessConfigs `json:"ipv6AccessConfigs"`
	AliasIPRanges     []InstanceNetworkInterfacesAliasIPRanges     `json:"aliasIPRanges"`
	Name              *string                                      `json:"name"`
	Network           *string                                      `json:"network"`
	NetworkIP         *string                                      `json:"networkIP"`
	Subnetwork        *string                                      `json:"subnetwork"`
}

type jsonInstanceNetworkInterfaces InstanceNetworkInterfaces

func (r *InstanceNetworkInterfaces) UnmarshalJSON(data []byte) error {
	var res jsonInstanceNetworkInterfaces
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceNetworkInterfaces
	} else {

		r.AccessConfigs = res.AccessConfigs

		r.IPv6AccessConfigs = res.IPv6AccessConfigs

		r.AliasIPRanges = res.AliasIPRanges

		r.Name = res.Name

		r.Network = res.Network

		r.NetworkIP = res.NetworkIP

		r.Subnetwork = res.Subnetwork

	}
	return nil
}

// This object is used to assert a desired state where this InstanceNetworkInterfaces is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceNetworkInterfaces *InstanceNetworkInterfaces = &InstanceNetworkInterfaces{empty: true}

func (r *InstanceNetworkInterfaces) Empty() bool {
	return r.empty
}

func (r *InstanceNetworkInterfaces) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceNetworkInterfaces) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceNetworkInterfacesAccessConfigs struct {
	empty                    bool                                                   `json:"-"`
	Name                     *string                                                `json:"name"`
	NatIP                    *string                                                `json:"natIP"`
	ExternalIPv6             *string                                                `json:"externalIPv6"`
	ExternalIPv6PrefixLength *string                                                `json:"externalIPv6PrefixLength"`
	SetPublicPtr             *bool                                                  `json:"setPublicPtr"`
	PublicPtrDomainName      *string                                                `json:"publicPtrDomainName"`
	NetworkTier              *InstanceNetworkInterfacesAccessConfigsNetworkTierEnum `json:"networkTier"`
	Type                     *InstanceNetworkInterfacesAccessConfigsTypeEnum        `json:"type"`
}

type jsonInstanceNetworkInterfacesAccessConfigs InstanceNetworkInterfacesAccessConfigs

func (r *InstanceNetworkInterfacesAccessConfigs) UnmarshalJSON(data []byte) error {
	var res jsonInstanceNetworkInterfacesAccessConfigs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceNetworkInterfacesAccessConfigs
	} else {

		r.Name = res.Name

		r.NatIP = res.NatIP

		r.ExternalIPv6 = res.ExternalIPv6

		r.ExternalIPv6PrefixLength = res.ExternalIPv6PrefixLength

		r.SetPublicPtr = res.SetPublicPtr

		r.PublicPtrDomainName = res.PublicPtrDomainName

		r.NetworkTier = res.NetworkTier

		r.Type = res.Type

	}
	return nil
}

// This object is used to assert a desired state where this InstanceNetworkInterfacesAccessConfigs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceNetworkInterfacesAccessConfigs *InstanceNetworkInterfacesAccessConfigs = &InstanceNetworkInterfacesAccessConfigs{empty: true}

func (r *InstanceNetworkInterfacesAccessConfigs) Empty() bool {
	return r.empty
}

func (r *InstanceNetworkInterfacesAccessConfigs) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceNetworkInterfacesAccessConfigs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceNetworkInterfacesIPv6AccessConfigs struct {
	empty                    bool                                                       `json:"-"`
	Name                     *string                                                    `json:"name"`
	NatIP                    *string                                                    `json:"natIP"`
	ExternalIPv6             *string                                                    `json:"externalIPv6"`
	ExternalIPv6PrefixLength *string                                                    `json:"externalIPv6PrefixLength"`
	SetPublicPtr             *bool                                                      `json:"setPublicPtr"`
	PublicPtrDomainName      *string                                                    `json:"publicPtrDomainName"`
	NetworkTier              *InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum `json:"networkTier"`
	Type                     *InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum        `json:"type"`
}

type jsonInstanceNetworkInterfacesIPv6AccessConfigs InstanceNetworkInterfacesIPv6AccessConfigs

func (r *InstanceNetworkInterfacesIPv6AccessConfigs) UnmarshalJSON(data []byte) error {
	var res jsonInstanceNetworkInterfacesIPv6AccessConfigs
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceNetworkInterfacesIPv6AccessConfigs
	} else {

		r.Name = res.Name

		r.NatIP = res.NatIP

		r.ExternalIPv6 = res.ExternalIPv6

		r.ExternalIPv6PrefixLength = res.ExternalIPv6PrefixLength

		r.SetPublicPtr = res.SetPublicPtr

		r.PublicPtrDomainName = res.PublicPtrDomainName

		r.NetworkTier = res.NetworkTier

		r.Type = res.Type

	}
	return nil
}

// This object is used to assert a desired state where this InstanceNetworkInterfacesIPv6AccessConfigs is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceNetworkInterfacesIPv6AccessConfigs *InstanceNetworkInterfacesIPv6AccessConfigs = &InstanceNetworkInterfacesIPv6AccessConfigs{empty: true}

func (r *InstanceNetworkInterfacesIPv6AccessConfigs) Empty() bool {
	return r.empty
}

func (r *InstanceNetworkInterfacesIPv6AccessConfigs) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceNetworkInterfacesIPv6AccessConfigs) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceNetworkInterfacesAliasIPRanges struct {
	empty               bool    `json:"-"`
	IPCidrRange         *string `json:"ipCidrRange"`
	SubnetworkRangeName *string `json:"subnetworkRangeName"`
}

type jsonInstanceNetworkInterfacesAliasIPRanges InstanceNetworkInterfacesAliasIPRanges

func (r *InstanceNetworkInterfacesAliasIPRanges) UnmarshalJSON(data []byte) error {
	var res jsonInstanceNetworkInterfacesAliasIPRanges
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceNetworkInterfacesAliasIPRanges
	} else {

		r.IPCidrRange = res.IPCidrRange

		r.SubnetworkRangeName = res.SubnetworkRangeName

	}
	return nil
}

// This object is used to assert a desired state where this InstanceNetworkInterfacesAliasIPRanges is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceNetworkInterfacesAliasIPRanges *InstanceNetworkInterfacesAliasIPRanges = &InstanceNetworkInterfacesAliasIPRanges{empty: true}

func (r *InstanceNetworkInterfacesAliasIPRanges) Empty() bool {
	return r.empty
}

func (r *InstanceNetworkInterfacesAliasIPRanges) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceNetworkInterfacesAliasIPRanges) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceScheduling struct {
	empty             bool    `json:"-"`
	AutomaticRestart  *bool   `json:"automaticRestart"`
	OnHostMaintenance *string `json:"onHostMaintenance"`
	Preemptible       *bool   `json:"preemptible"`
}

type jsonInstanceScheduling InstanceScheduling

func (r *InstanceScheduling) UnmarshalJSON(data []byte) error {
	var res jsonInstanceScheduling
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceScheduling
	} else {

		r.AutomaticRestart = res.AutomaticRestart

		r.OnHostMaintenance = res.OnHostMaintenance

		r.Preemptible = res.Preemptible

	}
	return nil
}

// This object is used to assert a desired state where this InstanceScheduling is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceScheduling *InstanceScheduling = &InstanceScheduling{empty: true}

func (r *InstanceScheduling) Empty() bool {
	return r.empty
}

func (r *InstanceScheduling) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceScheduling) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceServiceAccounts struct {
	empty  bool     `json:"-"`
	Email  *string  `json:"email"`
	Scopes []string `json:"scopes"`
}

type jsonInstanceServiceAccounts InstanceServiceAccounts

func (r *InstanceServiceAccounts) UnmarshalJSON(data []byte) error {
	var res jsonInstanceServiceAccounts
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceServiceAccounts
	} else {

		r.Email = res.Email

		r.Scopes = res.Scopes

	}
	return nil
}

// This object is used to assert a desired state where this InstanceServiceAccounts is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceServiceAccounts *InstanceServiceAccounts = &InstanceServiceAccounts{empty: true}

func (r *InstanceServiceAccounts) Empty() bool {
	return r.empty
}

func (r *InstanceServiceAccounts) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceServiceAccounts) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type InstanceShieldedInstanceConfig struct {
	empty                     bool  `json:"-"`
	EnableSecureBoot          *bool `json:"enableSecureBoot"`
	EnableVtpm                *bool `json:"enableVtpm"`
	EnableIntegrityMonitoring *bool `json:"enableIntegrityMonitoring"`
}

type jsonInstanceShieldedInstanceConfig InstanceShieldedInstanceConfig

func (r *InstanceShieldedInstanceConfig) UnmarshalJSON(data []byte) error {
	var res jsonInstanceShieldedInstanceConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyInstanceShieldedInstanceConfig
	} else {

		r.EnableSecureBoot = res.EnableSecureBoot

		r.EnableVtpm = res.EnableVtpm

		r.EnableIntegrityMonitoring = res.EnableIntegrityMonitoring

	}
	return nil
}

// This object is used to assert a desired state where this InstanceShieldedInstanceConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyInstanceShieldedInstanceConfig *InstanceShieldedInstanceConfig = &InstanceShieldedInstanceConfig{empty: true}

func (r *InstanceShieldedInstanceConfig) Empty() bool {
	return r.empty
}

func (r *InstanceShieldedInstanceConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *InstanceShieldedInstanceConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Instance) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "compute",
		Type:    "Instance",
		Version: "compute",
	}
}

func (r *Instance) ID() (string, error) {
	if err := extractInstanceFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"can_ip_forward":           dcl.ValueOrEmptyString(nr.CanIPForward),
		"cpu_platform":             dcl.ValueOrEmptyString(nr.CpuPlatform),
		"creation_timestamp":       dcl.ValueOrEmptyString(nr.CreationTimestamp),
		"deletion_protection":      dcl.ValueOrEmptyString(nr.DeletionProtection),
		"description":              dcl.ValueOrEmptyString(nr.Description),
		"disks":                    dcl.ValueOrEmptyString(nr.Disks),
		"guest_accelerators":       dcl.ValueOrEmptyString(nr.GuestAccelerators),
		"hostname":                 dcl.ValueOrEmptyString(nr.Hostname),
		"id":                       dcl.ValueOrEmptyString(nr.Id),
		"labels":                   dcl.ValueOrEmptyString(nr.Labels),
		"metadata":                 dcl.ValueOrEmptyString(nr.Metadata),
		"machine_type":             dcl.ValueOrEmptyString(nr.MachineType),
		"min_cpu_platform":         dcl.ValueOrEmptyString(nr.MinCpuPlatform),
		"name":                     dcl.ValueOrEmptyString(nr.Name),
		"network_interfaces":       dcl.ValueOrEmptyString(nr.NetworkInterfaces),
		"scheduling":               dcl.ValueOrEmptyString(nr.Scheduling),
		"service_accounts":         dcl.ValueOrEmptyString(nr.ServiceAccounts),
		"shielded_instance_config": dcl.ValueOrEmptyString(nr.ShieldedInstanceConfig),
		"status":                   dcl.ValueOrEmptyString(nr.Status),
		"status_message":           dcl.ValueOrEmptyString(nr.StatusMessage),
		"tags":                     dcl.ValueOrEmptyString(nr.Tags),
		"zone":                     dcl.ValueOrEmptyString(nr.Zone),
		"project":                  dcl.ValueOrEmptyString(nr.Project),
		"self_link":                dcl.ValueOrEmptyString(nr.SelfLink),
	}
	return dcl.Nprintf("projects/{{project}}/zones/{{zone}}/instances/{{name}}", params), nil
}

const InstanceMaxPage = -1

type InstanceList struct {
	Items []*Instance

	nextToken string

	pageSize int32

	resource *Instance
}

func (l *InstanceList) HasNext() bool {
	return l.nextToken != ""
}

func (l *InstanceList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listInstance(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListInstance(ctx context.Context, project, zone string) (*InstanceList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListInstanceWithMaxResults(ctx, project, zone, InstanceMaxPage)

}

func (c *Client) ListInstanceWithMaxResults(ctx context.Context, project, zone string, pageSize int32) (*InstanceList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Instance{
		Project: &project,
		Zone:    &zone,
	}
	items, token, err := c.listInstance(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &InstanceList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetInstance(ctx context.Context, r *Instance) (*Instance, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractInstanceFields(r)

	b, err := c.getInstanceRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalInstance(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Zone = r.Zone
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeInstanceNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractInstanceFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteInstance(ctx context.Context, r *Instance) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Instance resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Instance...")
	deleteOp := deleteInstanceOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllInstance deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllInstance(ctx context.Context, project, zone string, filter func(*Instance) bool) error {
	listObj, err := c.ListInstance(ctx, project, zone)
	if err != nil {
		return err
	}

	err = c.deleteAllInstance(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllInstance(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyInstance(ctx context.Context, rawDesired *Instance, opts ...dcl.ApplyOption) (*Instance, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Instance
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyInstanceHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyInstanceHelper(c *Client, ctx context.Context, rawDesired *Instance, opts ...dcl.ApplyOption) (*Instance, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyInstance...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractInstanceFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.instanceDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToInstanceDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []instanceApiOperation
	if create {
		ops = append(ops, &createInstanceOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyInstanceDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyInstanceDiff(c *Client, ctx context.Context, desired *Instance, rawDesired *Instance, ops []instanceApiOperation, opts ...dcl.ApplyOption) (*Instance, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetInstance(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createInstanceOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapInstance(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeInstanceNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeInstanceNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeInstanceDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractInstanceFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractInstanceFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffInstance(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}

func (r *Instance) GetPolicy(basePath string) (string, string, *bytes.Buffer, error) {
	u := r.getPolicyURL(basePath)
	body := &bytes.Buffer{}
	u, err := dcl.AddQueryParams(u, map[string]string{"optionsRequestedPolicyVersion": fmt.Sprintf("%d", r.IAMPolicyVersion())})
	if err != nil {
		return "", "", nil, err
	}
	return u, "GET", body, nil
}
