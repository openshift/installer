/*
Copyright 2022 The Kubernetes Authors.

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

package hash

import (
	"fmt"
	"hash"
	"hash/fnv"

	"github.com/davecgh/go-spew/spew"
	"k8s.io/apimachinery/pkg/util/dump"
)

// SpewHashObject writes specified object to hash using the spew library
// which follows pointers and prints actual values of the nested objects
// ensuring the hash does not change when a pointer changes.
func SpewHashObject(hasher hash.Hash, objectToWrite interface{}) error {
	hasher.Reset()
	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}

	if _, err := printer.Fprintf(hasher, "%#v", objectToWrite); err != nil {
		return fmt.Errorf("failed to write object to hasher")
	}
	return nil
}

// ComputeSpewHash computes the hash of a InstanceSpec using the spew library.
func ComputeSpewHash(objectToWrite interface{}) (uint32, error) {
	instanceSpecHasher := fnv.New32a()
	if err := SpewHashObject(instanceSpecHasher, objectToWrite); err != nil {
		return 0, err
	}
	return instanceSpecHasher.Sum32(), nil
}

// DeepHashObject writes specified object to hash using the spew library
// which follows pointers and prints actual values of the nested objects
// ensuring the hash does not change when a pointer changes.
// This function is taken from https://github.com/kubernetes/kubernetes/blob/v1.29.2/pkg/util/hash/hash.go#L26-L32
func DeepHashObject(hasher hash.Hash, objectToWrite interface{}) {
	hasher.Reset()
	fmt.Fprintf(hasher, "%v", dump.ForHash(objectToWrite))
}
