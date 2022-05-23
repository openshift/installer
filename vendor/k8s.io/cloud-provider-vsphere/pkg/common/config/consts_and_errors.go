/*
Copyright 2019 The Kubernetes Authors.

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

package config

import (
	"errors"
)

const (
	// DefaultRoundTripperCount is the number of allowed round trips
	// before an error is returned.
	DefaultRoundTripperCount uint = 3

	// DefaultAPIBinding is the default ADDRESS:PORT binding used for
	// exposing the API service.
	DefaultAPIBinding string = ":43001"

	// DefaultVCenterPortStr is the default port used to access vCenter in string form
	DefaultVCenterPortStr string = "443"
	// DefaultVCenterPort is the default port used to access vCenter in uint form
	DefaultVCenterPort uint = 443

	// DefaultSecretDirectory is the default path to the secrets directory.
	DefaultSecretDirectory string = "/etc/cloud/secrets"

	// IPv6Family string representation for IPv6
	IPv6Family = "ipv6"
	// IPv4Family string representation for IPv4
	IPv4Family = "ipv4"

	// DefaultIPFamily is the default IP addressing to use for networking
	DefaultIPFamily = IPv4Family

	// DefaultCredentialManager used for the Global CredMgr/Lister
	DefaultCredentialManager string = "Global"
)

var (
	// ErrUsernameMissing is returned when the provided username is empty.
	ErrUsernameMissing = errors.New("Username is missing")

	// ErrPasswordMissing is returned when the provided password is empty.
	ErrPasswordMissing = errors.New("Password is missing")

	// ErrInvalidVCenterIP is returned when the provided vCenter IP address is
	// missing from the provided configuration.
	ErrInvalidVCenterIP = errors.New("vsphere.conf does not have the VirtualCenter IP address specified")

	// ErrMissingVCenter is returned when the provided configuration does not
	// define any vCenters.
	ErrMissingVCenter = errors.New("No Virtual Center hosts defined")

	// ErrInvalidIPFamilyType is returned when an invalid IPFamily type is encountered
	ErrInvalidIPFamilyType = errors.New("Invalid IP Family type")
)
