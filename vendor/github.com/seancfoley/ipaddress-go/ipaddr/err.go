//
// Copyright 2020-2022 Sean C Foley
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
//

package ipaddr

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
)

type addressError struct {
	// key to look up the error message
	key string

	// an optional string with the address
	str string
}

func (a *addressError) Error() string {
	return getStr(a.str) + lookupStr("ipaddress.address.error") + " " + lookupStr(a.key)
}

func getStr(str string) (res string) {
	if len(str) > 0 {
		res = str + " "
	}
	return
}

// GetKey can be used to internationalize the error strings in the IPAddress library.
// The list of keys and their English translations are listed in IPAddressResources.properties.
// Use your own preferred method to map the key to your translations.
// One such option is golang.org/x/text which provides language tags (https://pkg.go.dev/golang.org/x/text/language?utm_source=godoc#Tag),
// which can then be mapped to catalogs, each catalog a list of translations for the set of keys provided here.
// In the code you supply a language key to use the right catalog.
// You can use the gotext tool to integrate those translations with your application.
func (a *addressError) GetKey() string {
	return a.key
}

type mergedError struct {
	addrerr.AddressError
	merged []addrerr.AddressError
}

func (a *mergedError) GetMerged() []addrerr.AddressError {
	return a.merged
}

type addressStringError struct {
	addressError
}

type addressStringNestedError struct {
	addressStringError
	nested addrerr.AddressStringError
}

func (a *addressStringNestedError) Error() string {
	return a.addressError.Error() + ": " + a.nested.Error()
}

type addressStringIndexError struct {
	addressStringError

	// byte index location in string of the error
	index int
}

func (a *addressStringIndexError) Error() string {
	return lookupStr("ipaddress.address.error") + " " + lookupStr(a.key) + " " + strconv.Itoa(a.index)
}

type hostNameError struct {
	addressError
}

// GetAddrError returns the nested address error which is nil for a host name error
func (a *hostNameError) GetAddrError() addrerr.AddressError {
	return nil
}

func (a *hostNameError) Error() string {
	return getStr(a.str) + lookupStr("ipaddress.host.error") + " " + lookupStr(a.key)
}

type hostNameNestedError struct {
	hostNameError
	nested error
}

type hostAddressNestedError struct {
	hostNameIndexError
	nested addrerr.AddressError
}

// GetAddrError returns the nested address error
func (a *hostAddressNestedError) GetAddrError() addrerr.AddressError {
	return a.nested
}

func (a *hostAddressNestedError) Error() string {
	if a.hostNameIndexError.key != "" {
		return getStr(a.str) + lookupStr("ipaddress.host.error") + " " + a.hostNameIndexError.Error() + " " + a.nested.Error()
	}
	return getStr(a.str) + lookupStr("ipaddress.host.error") + " " + a.nested.Error()
}

type hostNameIndexError struct {
	hostNameError

	// byte index location in string of the error
	index int
}

func (a *hostNameIndexError) Error() string {
	return getStr(a.str) + lookupStr("ipaddress.host.error") + " " + lookupStr(a.key) + " " + strconv.Itoa(a.index)
}

type incompatibleAddressError struct {
	addressError
}

type sizeMismatchError struct {
	incompatibleAddressError
}

type addressValueError struct {
	addressError

	// the value
	val int
}

///////////////////////////////////////////////

type wrappedErr struct {
	// root cause
	cause error

	// wrapper
	err error

	str string
}

func (wrappedErr *wrappedErr) Error() string {
	str := wrappedErr.str
	if len(str) > 0 {
		return str
	}
	str = wrappedErr.err.Error() + ": " + wrappedErr.cause.Error()
	wrappedErr.str = str
	return str
}

func newError(str string) error {
	return errors.New(str)
}

// errorF returns a formatted error
func errorF(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

// wrapErrf wraps the given error, but only if it is not nil.
func wrapErrf(err error, format string, a ...interface{}) error {
	return wrapper(true, err, format, a...)
}

// wrapToErrf is like wrapErrf but always returns an error
func wrapToErrf(err error, format string, a ...interface{}) error {
	return wrapper(false, err, format, a...)
}

func wrapper(nilIfFirstNil bool, err error, format string, a ...interface{}) error {
	if err == nil {
		if nilIfFirstNil {
			return nil
		}
		return errorF(format, a...)
	}
	return &wrappedErr{
		cause: err,
		err:   errorF(format, a...),
	}
}

type mergedErr struct {
	mergedErrs []error
	str        string
}

func (merged *mergedErr) Error() (str string) {
	str = merged.str
	if len(str) > 0 {
		return
	}
	mergedErrs := merged.mergedErrs
	errLen := len(mergedErrs)
	strs := make([]string, errLen)
	totalLen := 0
	for i, err := range mergedErrs {
		str := err.Error()
		strs[i] = str
		totalLen += len(str)
	}
	format := strings.Builder{}
	format.Grow(totalLen + errLen*2)
	format.WriteString(strs[0])
	for _, str := range strs[1:] {
		format.WriteString(", ")
		format.WriteString(str)
	}
	str = format.String()
	merged.str = str
	return
}

// mergeErrs merges an existing error with a new one
func mergeErrs(err error, format string, a ...interface{}) error {
	newErr := errorF(format, a...)
	if err == nil {
		return newErr
	}
	var merged []error
	if merge, isMergedErr := err.(*mergedErr); isMergedErr {
		merged = append(append([]error(nil), merge.mergedErrs...), newErr)
	} else {
		merged = []error{err, newErr}
	}
	return &mergedErr{mergedErrs: merged}
}

// mergeErrors merges multiple errors
func mergeAllErrs(errs ...error) error {
	var all []error
	allLen := len(errs)
	if allLen <= 1 {
		if allLen == 0 {
			return nil
		}
		return errs[0]
	}
	for _, err := range errs {
		if err != nil {
			if merge, isMergedErr := err.(*mergedErr); isMergedErr {
				all = append(all, merge.mergedErrs...)
			} else {
				all = append(all, err)
			}
		}
	}
	allLen = len(all)
	if allLen <= 1 {
		if allLen == 0 {
			return nil
		}
		return all[0]
	}
	return &mergedErr{mergedErrs: all}
}
