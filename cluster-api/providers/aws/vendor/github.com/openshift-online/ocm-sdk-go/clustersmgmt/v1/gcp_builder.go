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

// GCPBuilder contains the data and logic needed to build 'GCP' objects.
//
// Google cloud platform settings of a cluster.
type GCPBuilder struct {
	bitmap_                 uint32
	authURI                 string
	authProviderX509CertURL string
	authentication          *GcpAuthenticationBuilder
	clientID                string
	clientX509CertURL       string
	clientEmail             string
	privateKey              string
	privateKeyID            string
	privateServiceConnect   *GcpPrivateServiceConnectBuilder
	projectID               string
	security                *GcpSecurityBuilder
	tokenURI                string
	type_                   string
}

// NewGCP creates a new builder of 'GCP' objects.
func NewGCP() *GCPBuilder {
	return &GCPBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AuthURI sets the value of the 'auth_URI' attribute to the given value.
func (b *GCPBuilder) AuthURI(value string) *GCPBuilder {
	b.authURI = value
	b.bitmap_ |= 1
	return b
}

// AuthProviderX509CertURL sets the value of the 'auth_provider_X509_cert_URL' attribute to the given value.
func (b *GCPBuilder) AuthProviderX509CertURL(value string) *GCPBuilder {
	b.authProviderX509CertURL = value
	b.bitmap_ |= 2
	return b
}

// Authentication sets the value of the 'authentication' attribute to the given value.
//
// Google cloud platform authentication method of a cluster.
func (b *GCPBuilder) Authentication(value *GcpAuthenticationBuilder) *GCPBuilder {
	b.authentication = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// ClientID sets the value of the 'client_ID' attribute to the given value.
func (b *GCPBuilder) ClientID(value string) *GCPBuilder {
	b.clientID = value
	b.bitmap_ |= 8
	return b
}

// ClientX509CertURL sets the value of the 'client_X509_cert_URL' attribute to the given value.
func (b *GCPBuilder) ClientX509CertURL(value string) *GCPBuilder {
	b.clientX509CertURL = value
	b.bitmap_ |= 16
	return b
}

// ClientEmail sets the value of the 'client_email' attribute to the given value.
func (b *GCPBuilder) ClientEmail(value string) *GCPBuilder {
	b.clientEmail = value
	b.bitmap_ |= 32
	return b
}

// PrivateKey sets the value of the 'private_key' attribute to the given value.
func (b *GCPBuilder) PrivateKey(value string) *GCPBuilder {
	b.privateKey = value
	b.bitmap_ |= 64
	return b
}

// PrivateKeyID sets the value of the 'private_key_ID' attribute to the given value.
func (b *GCPBuilder) PrivateKeyID(value string) *GCPBuilder {
	b.privateKeyID = value
	b.bitmap_ |= 128
	return b
}

// PrivateServiceConnect sets the value of the 'private_service_connect' attribute to the given value.
//
// Google cloud platform private service connect configuration of a cluster.
func (b *GCPBuilder) PrivateServiceConnect(value *GcpPrivateServiceConnectBuilder) *GCPBuilder {
	b.privateServiceConnect = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// ProjectID sets the value of the 'project_ID' attribute to the given value.
func (b *GCPBuilder) ProjectID(value string) *GCPBuilder {
	b.projectID = value
	b.bitmap_ |= 512
	return b
}

// Security sets the value of the 'security' attribute to the given value.
//
// Google cloud platform security settings of a cluster.
func (b *GCPBuilder) Security(value *GcpSecurityBuilder) *GCPBuilder {
	b.security = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// TokenURI sets the value of the 'token_URI' attribute to the given value.
func (b *GCPBuilder) TokenURI(value string) *GCPBuilder {
	b.tokenURI = value
	b.bitmap_ |= 2048
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *GCPBuilder) Type(value string) *GCPBuilder {
	b.type_ = value
	b.bitmap_ |= 4096
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPBuilder) Copy(object *GCP) *GCPBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.authURI = object.authURI
	b.authProviderX509CertURL = object.authProviderX509CertURL
	if object.authentication != nil {
		b.authentication = NewGcpAuthentication().Copy(object.authentication)
	} else {
		b.authentication = nil
	}
	b.clientID = object.clientID
	b.clientX509CertURL = object.clientX509CertURL
	b.clientEmail = object.clientEmail
	b.privateKey = object.privateKey
	b.privateKeyID = object.privateKeyID
	if object.privateServiceConnect != nil {
		b.privateServiceConnect = NewGcpPrivateServiceConnect().Copy(object.privateServiceConnect)
	} else {
		b.privateServiceConnect = nil
	}
	b.projectID = object.projectID
	if object.security != nil {
		b.security = NewGcpSecurity().Copy(object.security)
	} else {
		b.security = nil
	}
	b.tokenURI = object.tokenURI
	b.type_ = object.type_
	return b
}

// Build creates a 'GCP' object using the configuration stored in the builder.
func (b *GCPBuilder) Build() (object *GCP, err error) {
	object = new(GCP)
	object.bitmap_ = b.bitmap_
	object.authURI = b.authURI
	object.authProviderX509CertURL = b.authProviderX509CertURL
	if b.authentication != nil {
		object.authentication, err = b.authentication.Build()
		if err != nil {
			return
		}
	}
	object.clientID = b.clientID
	object.clientX509CertURL = b.clientX509CertURL
	object.clientEmail = b.clientEmail
	object.privateKey = b.privateKey
	object.privateKeyID = b.privateKeyID
	if b.privateServiceConnect != nil {
		object.privateServiceConnect, err = b.privateServiceConnect.Build()
		if err != nil {
			return
		}
	}
	object.projectID = b.projectID
	if b.security != nil {
		object.security, err = b.security.Build()
		if err != nil {
			return
		}
	}
	object.tokenURI = b.tokenURI
	object.type_ = b.type_
	return
}
