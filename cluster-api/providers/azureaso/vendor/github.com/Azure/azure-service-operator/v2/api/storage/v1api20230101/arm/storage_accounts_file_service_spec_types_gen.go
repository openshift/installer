// Code generated by azure-service-operator-codegen. DO NOT EDIT.
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package arm

import "github.com/Azure/azure-service-operator/v2/pkg/genruntime"

type StorageAccountsFileService_Spec struct {
	Name string `json:"name,omitempty"`

	// Properties: The properties of File services in storage account.
	Properties *StorageAccounts_FileService_Properties_Spec `json:"properties,omitempty"`
}

var _ genruntime.ARMResourceSpec = &StorageAccountsFileService_Spec{}

// GetAPIVersion returns the ARM API version of the resource. This is always "2023-01-01"
func (service StorageAccountsFileService_Spec) GetAPIVersion() string {
	return "2023-01-01"
}

// GetName returns the Name of the resource
func (service *StorageAccountsFileService_Spec) GetName() string {
	return service.Name
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.Storage/storageAccounts/fileServices"
func (service *StorageAccountsFileService_Spec) GetType() string {
	return "Microsoft.Storage/storageAccounts/fileServices"
}

type StorageAccounts_FileService_Properties_Spec struct {
	// Cors: Specifies CORS rules for the File service. You can include up to five CorsRule elements in the request. If no
	// CorsRule elements are included in the request body, all CORS rules will be deleted, and CORS will be disabled for the
	// File service.
	Cors *CorsRules `json:"cors,omitempty"`

	// ProtocolSettings: Protocol settings for file service
	ProtocolSettings *ProtocolSettings `json:"protocolSettings,omitempty"`

	// ShareDeleteRetentionPolicy: The file service properties for share soft delete.
	ShareDeleteRetentionPolicy *DeleteRetentionPolicy `json:"shareDeleteRetentionPolicy,omitempty"`
}

// Protocol settings for file service
type ProtocolSettings struct {
	// Smb: Setting for SMB protocol
	Smb *SmbSetting `json:"smb,omitempty"`
}

// Setting for SMB protocol
type SmbSetting struct {
	// AuthenticationMethods: SMB authentication methods supported by server. Valid values are NTLMv2, Kerberos. Should be
	// passed as a string with delimiter ';'.
	AuthenticationMethods *string `json:"authenticationMethods,omitempty"`

	// ChannelEncryption: SMB channel encryption supported by server. Valid values are AES-128-CCM, AES-128-GCM, AES-256-GCM.
	// Should be passed as a string with delimiter ';'.
	ChannelEncryption *string `json:"channelEncryption,omitempty"`

	// KerberosTicketEncryption: Kerberos ticket encryption supported by server. Valid values are RC4-HMAC, AES-256. Should be
	// passed as a string with delimiter ';'
	KerberosTicketEncryption *string `json:"kerberosTicketEncryption,omitempty"`

	// Multichannel: Multichannel setting. Applies to Premium FileStorage only.
	Multichannel *Multichannel `json:"multichannel,omitempty"`

	// Versions: SMB protocol versions supported by server. Valid values are SMB2.1, SMB3.0, SMB3.1.1. Should be passed as a
	// string with delimiter ';'.
	Versions *string `json:"versions,omitempty"`
}

// Multichannel setting. Applies to Premium FileStorage only.
type Multichannel struct {
	// Enabled: Indicates whether multichannel is enabled
	Enabled *bool `json:"enabled,omitempty"`
}
