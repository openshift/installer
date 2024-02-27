/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// AWS represents the values of the 'AWS' type.
//
// _Amazon Web Services_ specific settings of a cluster.
type AWS struct {
	bitmap_                                uint32
	kmsKeyArn                              string
	sts                                    *STS
	accessKeyID                            string
	accountID                              string
	additionalComputeSecurityGroupIds      []string
	additionalControlPlaneSecurityGroupIds []string
	additionalInfraSecurityGroupIds        []string
	auditLog                               *AuditLog
	billingAccountID                       string
	ec2MetadataHttpTokens                  Ec2MetadataHttpTokens
	etcdEncryption                         *AwsEtcdEncryption
	privateHostedZoneID                    string
	privateHostedZoneRoleARN               string
	privateLinkConfiguration               *PrivateLinkClusterConfiguration
	secretAccessKey                        string
	subnetIDs                              []string
	tags                                   map[string]string
	privateLink                            bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWS) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// KMSKeyArn returns the value of the 'KMS_key_arn' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Customer Managed Key to encrypt EBS Volume
func (o *AWS) KMSKeyArn() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.kmsKeyArn
	}
	return ""
}

// GetKMSKeyArn returns the value of the 'KMS_key_arn' attribute and
// a flag indicating if the attribute has a value.
//
// Customer Managed Key to encrypt EBS Volume
func (o *AWS) GetKMSKeyArn() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.kmsKeyArn
	}
	return
}

// STS returns the value of the 'STS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (o *AWS) STS() *STS {
	if o != nil && o.bitmap_&2 != 0 {
		return o.sts
	}
	return nil
}

// GetSTS returns the value of the 'STS' attribute and
// a flag indicating if the attribute has a value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (o *AWS) GetSTS() (value *STS, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.sts
	}
	return
}

// AccessKeyID returns the value of the 'access_key_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS access key identifier.
func (o *AWS) AccessKeyID() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.accessKeyID
	}
	return ""
}

// GetAccessKeyID returns the value of the 'access_key_ID' attribute and
// a flag indicating if the attribute has a value.
//
// AWS access key identifier.
func (o *AWS) GetAccessKeyID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.accessKeyID
	}
	return
}

// AccountID returns the value of the 'account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS account identifier.
func (o *AWS) AccountID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.accountID
	}
	return ""
}

// GetAccountID returns the value of the 'account_ID' attribute and
// a flag indicating if the attribute has a value.
//
// AWS account identifier.
func (o *AWS) GetAccountID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.accountID
	}
	return
}

// AdditionalComputeSecurityGroupIds returns the value of the 'additional_compute_security_group_ids' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional AWS Security Groups to be added to default worker (compute) machine pool.
func (o *AWS) AdditionalComputeSecurityGroupIds() []string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.additionalComputeSecurityGroupIds
	}
	return nil
}

// GetAdditionalComputeSecurityGroupIds returns the value of the 'additional_compute_security_group_ids' attribute and
// a flag indicating if the attribute has a value.
//
// Additional AWS Security Groups to be added to default worker (compute) machine pool.
func (o *AWS) GetAdditionalComputeSecurityGroupIds() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.additionalComputeSecurityGroupIds
	}
	return
}

// AdditionalControlPlaneSecurityGroupIds returns the value of the 'additional_control_plane_security_group_ids' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional AWS Security Groups to be added to default control plane machine pool.
func (o *AWS) AdditionalControlPlaneSecurityGroupIds() []string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.additionalControlPlaneSecurityGroupIds
	}
	return nil
}

// GetAdditionalControlPlaneSecurityGroupIds returns the value of the 'additional_control_plane_security_group_ids' attribute and
// a flag indicating if the attribute has a value.
//
// Additional AWS Security Groups to be added to default control plane machine pool.
func (o *AWS) GetAdditionalControlPlaneSecurityGroupIds() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.additionalControlPlaneSecurityGroupIds
	}
	return
}

// AdditionalInfraSecurityGroupIds returns the value of the 'additional_infra_security_group_ids' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional AWS Security Groups to be added to default infra machine pool.
func (o *AWS) AdditionalInfraSecurityGroupIds() []string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.additionalInfraSecurityGroupIds
	}
	return nil
}

// GetAdditionalInfraSecurityGroupIds returns the value of the 'additional_infra_security_group_ids' attribute and
// a flag indicating if the attribute has a value.
//
// Additional AWS Security Groups to be added to default infra machine pool.
func (o *AWS) GetAdditionalInfraSecurityGroupIds() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.additionalInfraSecurityGroupIds
	}
	return
}

// AuditLog returns the value of the 'audit_log' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Audit log forwarding configuration
func (o *AWS) AuditLog() *AuditLog {
	if o != nil && o.bitmap_&128 != 0 {
		return o.auditLog
	}
	return nil
}

// GetAuditLog returns the value of the 'audit_log' attribute and
// a flag indicating if the attribute has a value.
//
// Audit log forwarding configuration
func (o *AWS) GetAuditLog() (value *AuditLog, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.auditLog
	}
	return
}

// BillingAccountID returns the value of the 'billing_account_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// BillingAccountID is the account used for billing subscriptions purchased via the marketplace
func (o *AWS) BillingAccountID() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.billingAccountID
	}
	return ""
}

// GetBillingAccountID returns the value of the 'billing_account_ID' attribute and
// a flag indicating if the attribute has a value.
//
// BillingAccountID is the account used for billing subscriptions purchased via the marketplace
func (o *AWS) GetBillingAccountID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.billingAccountID
	}
	return
}

// Ec2MetadataHttpTokens returns the value of the 'ec_2_metadata_http_tokens' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Which Ec2MetadataHttpTokens to use for metadata service interaction options for EC2 instances
func (o *AWS) Ec2MetadataHttpTokens() Ec2MetadataHttpTokens {
	if o != nil && o.bitmap_&512 != 0 {
		return o.ec2MetadataHttpTokens
	}
	return Ec2MetadataHttpTokens("")
}

// GetEc2MetadataHttpTokens returns the value of the 'ec_2_metadata_http_tokens' attribute and
// a flag indicating if the attribute has a value.
//
// Which Ec2MetadataHttpTokens to use for metadata service interaction options for EC2 instances
func (o *AWS) GetEc2MetadataHttpTokens() (value Ec2MetadataHttpTokens, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.ec2MetadataHttpTokens
	}
	return
}

// EtcdEncryption returns the value of the 'etcd_encryption' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Related etcd encryption configuration
func (o *AWS) EtcdEncryption() *AwsEtcdEncryption {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.etcdEncryption
	}
	return nil
}

// GetEtcdEncryption returns the value of the 'etcd_encryption' attribute and
// a flag indicating if the attribute has a value.
//
// Related etcd encryption configuration
func (o *AWS) GetEtcdEncryption() (value *AwsEtcdEncryption, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.etcdEncryption
	}
	return
}

// PrivateHostedZoneID returns the value of the 'private_hosted_zone_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of private hosted zone.
func (o *AWS) PrivateHostedZoneID() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.privateHostedZoneID
	}
	return ""
}

// GetPrivateHostedZoneID returns the value of the 'private_hosted_zone_ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID of private hosted zone.
func (o *AWS) GetPrivateHostedZoneID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.privateHostedZoneID
	}
	return
}

// PrivateHostedZoneRoleARN returns the value of the 'private_hosted_zone_role_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Role ARN for private hosted zone.
func (o *AWS) PrivateHostedZoneRoleARN() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.privateHostedZoneRoleARN
	}
	return ""
}

// GetPrivateHostedZoneRoleARN returns the value of the 'private_hosted_zone_role_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// Role ARN for private hosted zone.
func (o *AWS) GetPrivateHostedZoneRoleARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.privateHostedZoneRoleARN
	}
	return
}

// PrivateLink returns the value of the 'private_link' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Sets cluster to be inaccessible externally.
func (o *AWS) PrivateLink() bool {
	if o != nil && o.bitmap_&8192 != 0 {
		return o.privateLink
	}
	return false
}

// GetPrivateLink returns the value of the 'private_link' attribute and
// a flag indicating if the attribute has a value.
//
// Sets cluster to be inaccessible externally.
func (o *AWS) GetPrivateLink() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8192 != 0
	if ok {
		value = o.privateLink
	}
	return
}

// PrivateLinkConfiguration returns the value of the 'private_link_configuration' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Manages additional configuration for Private Links.
func (o *AWS) PrivateLinkConfiguration() *PrivateLinkClusterConfiguration {
	if o != nil && o.bitmap_&16384 != 0 {
		return o.privateLinkConfiguration
	}
	return nil
}

// GetPrivateLinkConfiguration returns the value of the 'private_link_configuration' attribute and
// a flag indicating if the attribute has a value.
//
// Manages additional configuration for Private Links.
func (o *AWS) GetPrivateLinkConfiguration() (value *PrivateLinkClusterConfiguration, ok bool) {
	ok = o != nil && o.bitmap_&16384 != 0
	if ok {
		value = o.privateLinkConfiguration
	}
	return
}

// SecretAccessKey returns the value of the 'secret_access_key' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AWS secret access key.
func (o *AWS) SecretAccessKey() string {
	if o != nil && o.bitmap_&32768 != 0 {
		return o.secretAccessKey
	}
	return ""
}

// GetSecretAccessKey returns the value of the 'secret_access_key' attribute and
// a flag indicating if the attribute has a value.
//
// AWS secret access key.
func (o *AWS) GetSecretAccessKey() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32768 != 0
	if ok {
		value = o.secretAccessKey
	}
	return
}

// SubnetIDs returns the value of the 'subnet_IDs' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The subnet ids to be used when installing the cluster.
func (o *AWS) SubnetIDs() []string {
	if o != nil && o.bitmap_&65536 != 0 {
		return o.subnetIDs
	}
	return nil
}

// GetSubnetIDs returns the value of the 'subnet_IDs' attribute and
// a flag indicating if the attribute has a value.
//
// The subnet ids to be used when installing the cluster.
func (o *AWS) GetSubnetIDs() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&65536 != 0
	if ok {
		value = o.subnetIDs
	}
	return
}

// Tags returns the value of the 'tags' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates
func (o *AWS) Tags() map[string]string {
	if o != nil && o.bitmap_&131072 != 0 {
		return o.tags
	}
	return nil
}

// GetTags returns the value of the 'tags' attribute and
// a flag indicating if the attribute has a value.
//
// Optional keys and values that the installer will add as tags to all AWS resources it creates
func (o *AWS) GetTags() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&131072 != 0
	if ok {
		value = o.tags
	}
	return
}

// AWSListKind is the name of the type used to represent list of objects of
// type 'AWS'.
const AWSListKind = "AWSList"

// AWSListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS'.
const AWSListLinkKind = "AWSListLink"

// AWSNilKind is the name of the type used to nil lists of objects of
// type 'AWS'.
const AWSListNilKind = "AWSListNil"

// AWSList is a list of values of the 'AWS' type.
type AWSList struct {
	href  string
	link  bool
	items []*AWS
}

// Len returns the length of the list.
func (l *AWSList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AWSList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSList) Get(i int) *AWS {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *AWSList) Slice() []*AWS {
	var slice []*AWS
	if l == nil {
		slice = make([]*AWS, 0)
	} else {
		slice = make([]*AWS, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSList) Each(f func(item *AWS) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *AWSList) Range(f func(index int, item *AWS) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
