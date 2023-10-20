//go:build !windows
// +build !windows

/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains the function that returns the trusted CA certificates for operating systems
// other than Windows, where Go knows how to load the system trusted CA store.

package internal

import (
	"crypto/x509"
)

// loadSystemCAs loads the trusted CA certifites from the system trusted CA store.
func loadSystemCAs() (pool *x509.CertPool, err error) {
	pool, err = x509.SystemCertPool()
	return
}
