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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// AWSBuilder contains the data and logic needed to build 'AWS' objects.
//
// _Amazon Web Services_ specific settings of a cluster.
type AWSBuilder struct {
	bitmap_                                uint32
	kmsKeyArn                              string
	sts                                    *STSBuilder
	accessKeyID                            string
	accountID                              string
	additionalAllowedPrincipals            []string
	additionalComputeSecurityGroupIds      []string
	additionalControlPlaneSecurityGroupIds []string
	additionalInfraSecurityGroupIds        []string
	auditLog                               *AuditLogBuilder
	billingAccountID                       string
	ec2MetadataHttpTokens                  Ec2MetadataHttpTokens
	etcdEncryption                         *AwsEtcdEncryptionBuilder
	hcpInternalCommunicationHostedZoneId   string
	privateHostedZoneID                    string
	privateHostedZoneRoleARN               string
	privateLinkConfiguration               *PrivateLinkClusterConfigurationBuilder
	secretAccessKey                        string
	subnetIDs                              []string
	tags                                   map[string]string
	vpcEndpointRoleArn                     string
	privateLink                            bool
}

// NewAWS creates a new builder of 'AWS' objects.
func NewAWS() *AWSBuilder {
	return &AWSBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// KMSKeyArn sets the value of the 'KMS_key_arn' attribute to the given value.
func (b *AWSBuilder) KMSKeyArn(value string) *AWSBuilder {
	b.kmsKeyArn = value
	b.bitmap_ |= 1
	return b
}

// STS sets the value of the 'STS' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *AWSBuilder) STS(value *STSBuilder) *AWSBuilder {
	b.sts = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// AccessKeyID sets the value of the 'access_key_ID' attribute to the given value.
func (b *AWSBuilder) AccessKeyID(value string) *AWSBuilder {
	b.accessKeyID = value
	b.bitmap_ |= 4
	return b
}

// AccountID sets the value of the 'account_ID' attribute to the given value.
func (b *AWSBuilder) AccountID(value string) *AWSBuilder {
	b.accountID = value
	b.bitmap_ |= 8
	return b
}

// AdditionalAllowedPrincipals sets the value of the 'additional_allowed_principals' attribute to the given values.
func (b *AWSBuilder) AdditionalAllowedPrincipals(values ...string) *AWSBuilder {
	b.additionalAllowedPrincipals = make([]string, len(values))
	copy(b.additionalAllowedPrincipals, values)
	b.bitmap_ |= 16
	return b
}

// AdditionalComputeSecurityGroupIds sets the value of the 'additional_compute_security_group_ids' attribute to the given values.
func (b *AWSBuilder) AdditionalComputeSecurityGroupIds(values ...string) *AWSBuilder {
	b.additionalComputeSecurityGroupIds = make([]string, len(values))
	copy(b.additionalComputeSecurityGroupIds, values)
	b.bitmap_ |= 32
	return b
}

// AdditionalControlPlaneSecurityGroupIds sets the value of the 'additional_control_plane_security_group_ids' attribute to the given values.
func (b *AWSBuilder) AdditionalControlPlaneSecurityGroupIds(values ...string) *AWSBuilder {
	b.additionalControlPlaneSecurityGroupIds = make([]string, len(values))
	copy(b.additionalControlPlaneSecurityGroupIds, values)
	b.bitmap_ |= 64
	return b
}

// AdditionalInfraSecurityGroupIds sets the value of the 'additional_infra_security_group_ids' attribute to the given values.
func (b *AWSBuilder) AdditionalInfraSecurityGroupIds(values ...string) *AWSBuilder {
	b.additionalInfraSecurityGroupIds = make([]string, len(values))
	copy(b.additionalInfraSecurityGroupIds, values)
	b.bitmap_ |= 128
	return b
}

// AuditLog sets the value of the 'audit_log' attribute to the given value.
//
// Contains the necessary attributes to support audit log forwarding
func (b *AWSBuilder) AuditLog(value *AuditLogBuilder) *AWSBuilder {
	b.auditLog = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// BillingAccountID sets the value of the 'billing_account_ID' attribute to the given value.
func (b *AWSBuilder) BillingAccountID(value string) *AWSBuilder {
	b.billingAccountID = value
	b.bitmap_ |= 512
	return b
}

// Ec2MetadataHttpTokens sets the value of the 'ec_2_metadata_http_tokens' attribute to the given value.
//
// Which Ec2MetadataHttpTokens to use for metadata service interaction options for EC2 instances
func (b *AWSBuilder) Ec2MetadataHttpTokens(value Ec2MetadataHttpTokens) *AWSBuilder {
	b.ec2MetadataHttpTokens = value
	b.bitmap_ |= 1024
	return b
}

// EtcdEncryption sets the value of the 'etcd_encryption' attribute to the given value.
//
// Contains the necessary attributes to support etcd encryption for AWS based clusters.
func (b *AWSBuilder) EtcdEncryption(value *AwsEtcdEncryptionBuilder) *AWSBuilder {
	b.etcdEncryption = value
	if value != nil {
		b.bitmap_ |= 2048
	} else {
		b.bitmap_ &^= 2048
	}
	return b
}

// HcpInternalCommunicationHostedZoneId sets the value of the 'hcp_internal_communication_hosted_zone_id' attribute to the given value.
func (b *AWSBuilder) HcpInternalCommunicationHostedZoneId(value string) *AWSBuilder {
	b.hcpInternalCommunicationHostedZoneId = value
	b.bitmap_ |= 4096
	return b
}

// PrivateHostedZoneID sets the value of the 'private_hosted_zone_ID' attribute to the given value.
func (b *AWSBuilder) PrivateHostedZoneID(value string) *AWSBuilder {
	b.privateHostedZoneID = value
	b.bitmap_ |= 8192
	return b
}

// PrivateHostedZoneRoleARN sets the value of the 'private_hosted_zone_role_ARN' attribute to the given value.
func (b *AWSBuilder) PrivateHostedZoneRoleARN(value string) *AWSBuilder {
	b.privateHostedZoneRoleARN = value
	b.bitmap_ |= 16384
	return b
}

// PrivateLink sets the value of the 'private_link' attribute to the given value.
func (b *AWSBuilder) PrivateLink(value bool) *AWSBuilder {
	b.privateLink = value
	b.bitmap_ |= 32768
	return b
}

// PrivateLinkConfiguration sets the value of the 'private_link_configuration' attribute to the given value.
//
// Manages the configuration for the Private Links.
func (b *AWSBuilder) PrivateLinkConfiguration(value *PrivateLinkClusterConfigurationBuilder) *AWSBuilder {
	b.privateLinkConfiguration = value
	if value != nil {
		b.bitmap_ |= 65536
	} else {
		b.bitmap_ &^= 65536
	}
	return b
}

// SecretAccessKey sets the value of the 'secret_access_key' attribute to the given value.
func (b *AWSBuilder) SecretAccessKey(value string) *AWSBuilder {
	b.secretAccessKey = value
	b.bitmap_ |= 131072
	return b
}

// SubnetIDs sets the value of the 'subnet_IDs' attribute to the given values.
func (b *AWSBuilder) SubnetIDs(values ...string) *AWSBuilder {
	b.subnetIDs = make([]string, len(values))
	copy(b.subnetIDs, values)
	b.bitmap_ |= 262144
	return b
}

// Tags sets the value of the 'tags' attribute to the given value.
func (b *AWSBuilder) Tags(value map[string]string) *AWSBuilder {
	b.tags = value
	if value != nil {
		b.bitmap_ |= 524288
	} else {
		b.bitmap_ &^= 524288
	}
	return b
}

// VpcEndpointRoleArn sets the value of the 'vpc_endpoint_role_arn' attribute to the given value.
func (b *AWSBuilder) VpcEndpointRoleArn(value string) *AWSBuilder {
	b.vpcEndpointRoleArn = value
	b.bitmap_ |= 1048576
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSBuilder) Copy(object *AWS) *AWSBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.kmsKeyArn = object.kmsKeyArn
	if object.sts != nil {
		b.sts = NewSTS().Copy(object.sts)
	} else {
		b.sts = nil
	}
	b.accessKeyID = object.accessKeyID
	b.accountID = object.accountID
	if object.additionalAllowedPrincipals != nil {
		b.additionalAllowedPrincipals = make([]string, len(object.additionalAllowedPrincipals))
		copy(b.additionalAllowedPrincipals, object.additionalAllowedPrincipals)
	} else {
		b.additionalAllowedPrincipals = nil
	}
	if object.additionalComputeSecurityGroupIds != nil {
		b.additionalComputeSecurityGroupIds = make([]string, len(object.additionalComputeSecurityGroupIds))
		copy(b.additionalComputeSecurityGroupIds, object.additionalComputeSecurityGroupIds)
	} else {
		b.additionalComputeSecurityGroupIds = nil
	}
	if object.additionalControlPlaneSecurityGroupIds != nil {
		b.additionalControlPlaneSecurityGroupIds = make([]string, len(object.additionalControlPlaneSecurityGroupIds))
		copy(b.additionalControlPlaneSecurityGroupIds, object.additionalControlPlaneSecurityGroupIds)
	} else {
		b.additionalControlPlaneSecurityGroupIds = nil
	}
	if object.additionalInfraSecurityGroupIds != nil {
		b.additionalInfraSecurityGroupIds = make([]string, len(object.additionalInfraSecurityGroupIds))
		copy(b.additionalInfraSecurityGroupIds, object.additionalInfraSecurityGroupIds)
	} else {
		b.additionalInfraSecurityGroupIds = nil
	}
	if object.auditLog != nil {
		b.auditLog = NewAuditLog().Copy(object.auditLog)
	} else {
		b.auditLog = nil
	}
	b.billingAccountID = object.billingAccountID
	b.ec2MetadataHttpTokens = object.ec2MetadataHttpTokens
	if object.etcdEncryption != nil {
		b.etcdEncryption = NewAwsEtcdEncryption().Copy(object.etcdEncryption)
	} else {
		b.etcdEncryption = nil
	}
	b.hcpInternalCommunicationHostedZoneId = object.hcpInternalCommunicationHostedZoneId
	b.privateHostedZoneID = object.privateHostedZoneID
	b.privateHostedZoneRoleARN = object.privateHostedZoneRoleARN
	b.privateLink = object.privateLink
	if object.privateLinkConfiguration != nil {
		b.privateLinkConfiguration = NewPrivateLinkClusterConfiguration().Copy(object.privateLinkConfiguration)
	} else {
		b.privateLinkConfiguration = nil
	}
	b.secretAccessKey = object.secretAccessKey
	if object.subnetIDs != nil {
		b.subnetIDs = make([]string, len(object.subnetIDs))
		copy(b.subnetIDs, object.subnetIDs)
	} else {
		b.subnetIDs = nil
	}
	if len(object.tags) > 0 {
		b.tags = map[string]string{}
		for k, v := range object.tags {
			b.tags[k] = v
		}
	} else {
		b.tags = nil
	}
	b.vpcEndpointRoleArn = object.vpcEndpointRoleArn
	return b
}

// Build creates a 'AWS' object using the configuration stored in the builder.
func (b *AWSBuilder) Build() (object *AWS, err error) {
	object = new(AWS)
	object.bitmap_ = b.bitmap_
	object.kmsKeyArn = b.kmsKeyArn
	if b.sts != nil {
		object.sts, err = b.sts.Build()
		if err != nil {
			return
		}
	}
	object.accessKeyID = b.accessKeyID
	object.accountID = b.accountID
	if b.additionalAllowedPrincipals != nil {
		object.additionalAllowedPrincipals = make([]string, len(b.additionalAllowedPrincipals))
		copy(object.additionalAllowedPrincipals, b.additionalAllowedPrincipals)
	}
	if b.additionalComputeSecurityGroupIds != nil {
		object.additionalComputeSecurityGroupIds = make([]string, len(b.additionalComputeSecurityGroupIds))
		copy(object.additionalComputeSecurityGroupIds, b.additionalComputeSecurityGroupIds)
	}
	if b.additionalControlPlaneSecurityGroupIds != nil {
		object.additionalControlPlaneSecurityGroupIds = make([]string, len(b.additionalControlPlaneSecurityGroupIds))
		copy(object.additionalControlPlaneSecurityGroupIds, b.additionalControlPlaneSecurityGroupIds)
	}
	if b.additionalInfraSecurityGroupIds != nil {
		object.additionalInfraSecurityGroupIds = make([]string, len(b.additionalInfraSecurityGroupIds))
		copy(object.additionalInfraSecurityGroupIds, b.additionalInfraSecurityGroupIds)
	}
	if b.auditLog != nil {
		object.auditLog, err = b.auditLog.Build()
		if err != nil {
			return
		}
	}
	object.billingAccountID = b.billingAccountID
	object.ec2MetadataHttpTokens = b.ec2MetadataHttpTokens
	if b.etcdEncryption != nil {
		object.etcdEncryption, err = b.etcdEncryption.Build()
		if err != nil {
			return
		}
	}
	object.hcpInternalCommunicationHostedZoneId = b.hcpInternalCommunicationHostedZoneId
	object.privateHostedZoneID = b.privateHostedZoneID
	object.privateHostedZoneRoleARN = b.privateHostedZoneRoleARN
	object.privateLink = b.privateLink
	if b.privateLinkConfiguration != nil {
		object.privateLinkConfiguration, err = b.privateLinkConfiguration.Build()
		if err != nil {
			return
		}
	}
	object.secretAccessKey = b.secretAccessKey
	if b.subnetIDs != nil {
		object.subnetIDs = make([]string, len(b.subnetIDs))
		copy(object.subnetIDs, b.subnetIDs)
	}
	if b.tags != nil {
		object.tags = make(map[string]string)
		for k, v := range b.tags {
			object.tags[k] = v
		}
	}
	object.vpcEndpointRoleArn = b.vpcEndpointRoleArn
	return
}
