/*
Copyright 2023 The Kubernetes Authors.

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

package conversion

import (
	"bytes"
	"encoding/json"
	"hash/crc64"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"

	hasher "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/hash"
)

// RestorerFor holds all field restorers for a given type T.
type RestorerFor[T metav1.Object] map[string]FieldRestorerFor[T]

type FieldRestorerFor[T any] interface {
	marshalState(T, T) (json.RawMessage, error)
	restore(json.RawMessage, T) error
}

// HashedFieldRestorerOpt is a modifier which adds optional behaviour to a
// HashedFieldRestorer. It can be passed as an argument to
// HashedFieldRestorer().
type HashedFieldRestorerOpt[T, F any] func(*hashedFieldRestorer[T, F])

// HashedFieldRestorer restores a field to its original pre-conversion state
// only if it was not modified while converted. It does this by doing a full
// round-trip conversion before saving its state, and storing the hash of an
// unmodified round-trip converted object.
//
// HashedFieldRestorer takes 2 type arguments:
//
//	T: The type of the object being converted. This is a pointer type.
//	F: The type of the field in T which is being restored. This is not a
//	   pointer type.
//
// The type arguments can usually be omitted, as they can be inferred from the
// arguments.
//
// HashedFieldRestorer takes 2 arguments:
//
//	getField: A function which takes the object being converted (of type T) and returns a pointer to the field being restored (of type *F).
//	restoreField: Described below.
//
// restoreField is used to restore parts of the field to their original
// pre-conversion state which were not restored by the normal conversion
// functions. This is to preserve idempotency when the normal conversion
// functions are lossy, so it is only called if the field was not modified
// during conversion.
//
// restoreField takes 2 arguments:
//
//	previous: A pointer to the saved original state of the field before conversion.
//	dst: A pointer to the field which needs to be updated.
//
// The normal conversion functions should do as much work as possible, so
// restoreField should only be used to restore state which cannot be restored
// any other way.
func HashedFieldRestorer[T, F any](getField func(T) *F, restoreField func(*F, *F), opts ...HashedFieldRestorerOpt[T, F]) FieldRestorerFor[T] {
	restorer := hashedFieldRestorer[T, F]{
		getField:     getField,
		restoreField: restoreField,
	}
	for _, opt := range opts {
		opt(&restorer)
	}

	return restorer
}

// HashedFilterField adds a field filter to a HashedFieldRestorer. The field
// filter returns a modified copy of the field to be used for comparison. This
// can be used to ignore certain changes during comparison.
func HashedFilterField[T, F any](f func(*F) *F) func(*hashedFieldRestorer[T, F]) {
	return func(r *hashedFieldRestorer[T, F]) {
		r.filterField = f
	}
}

type hashedFieldRestorer[T any, F any] struct {
	// getField returns the field to be restored
	getField func(T) *F

	// restoreField restores the field to its original pre-conversion state
	restoreField func(*F, *F)

	// filterField returns a modified copy of the field to be used for comparison. i.e. filtered fields will be ignored in comparison.
	filterField func(*F) *F
}

var _ FieldRestorerFor[any] = hashedFieldRestorer[any, any]{}

//nolint:unused
type hashedFieldRestoreState struct {
	Hash []byte          `json:"h,omitempty"`
	Data json.RawMessage `json:"d,omitempty"`
}

//nolint:unused
func (r hashedFieldRestorer[T, F]) getHash(obj T) []byte {
	f := r.getField(obj)
	if r.filterField != nil {
		f = r.filterField(f)
	}

	table := crc64.MakeTable(crc64.ECMA)
	hash := crc64.New(table)
	hasher.DeepHashObject(hash, f)
	return hash.Sum(make([]byte, 0, 8))
}

//nolint:unused
func (r hashedFieldRestorer[T, F]) marshalState(src, compare T) (json.RawMessage, error) {
	f := r.getField(src)
	if r.filterField != nil {
		f = r.filterField(f)
	}

	b, err := json.Marshal(f)
	if err != nil {
		return nil, err
	}

	srcHash := r.getHash(src)
	cmpHash := r.getHash(compare)
	if bytes.Equal(srcHash, cmpHash) {
		return nil, nil
	}

	restoreState := hashedFieldRestoreState{
		Hash: cmpHash,
		Data: b,
	}
	return json.Marshal(restoreState)
}

//nolint:unused
func (r hashedFieldRestorer[T, F]) restore(b json.RawMessage, dst T) error {
	var restoreState hashedFieldRestoreState
	if err := json.Unmarshal(b, &restoreState); err != nil {
		return err
	}

	if bytes.Equal(r.getHash(dst), restoreState.Hash) {
		var previous F
		if err := json.Unmarshal(restoreState.Data, &previous); err != nil {
			return err
		}
		r.restoreField(&previous, r.getField(dst))
	}

	return nil
}

// UnconditionalFieldRestorer restores a field to its previous value without checking to see if it was modified during conversion.
func UnconditionalFieldRestorer[T, F any](getField func(T) *F) FieldRestorerFor[T] {
	return unconditionalFieldRestorer[T, F]{
		getField: getField,
	}
}

type unconditionalFieldRestorer[T any, F any] struct {
	getField func(T) *F
}

var _ FieldRestorerFor[any] = unconditionalFieldRestorer[any, any]{}

//nolint:unused
type unconditionalFieldRestoreState json.RawMessage

//nolint:unused
func (r unconditionalFieldRestorer[T, F]) marshalState(src, _ T) (json.RawMessage, error) {
	f := r.getField(src)

	// We could do this with a comparable constraint on F, but that's too restrictive for an arbitrary struct.
	if f == nil || reflect.ValueOf(f).Elem().IsZero() {
		return nil, nil
	}

	return json.Marshal(r.getField(src))
}

//nolint:unused
func (r unconditionalFieldRestorer[T, F]) restore(b json.RawMessage, dst T) error {
	return json.Unmarshal(b, r.getField(dst))
}

// restoreData serialises state for a set of restorers.
type restoreData map[string]json.RawMessage

type pointerToObject[T any] interface {
	*T
	metav1.Object
}

// ConvertAndRestore converts src to dst. During conversion is reads and uses any restore data in the src annotation, and writes new restore data to the dst annotation.
func ConvertAndRestore[S pointerToObject[SO], SO any, D metav1.Object](src S, dst D, convert func(S, D, conversion.Scope) error, unconvert func(D, S, conversion.Scope) error, srcRestorer RestorerFor[S], dstRestorer RestorerFor[D]) error {
	var dstRestoreData restoreData
	srcAnnotations := src.GetAnnotations()
	if srcData, ok := srcAnnotations[utilconversion.DataAnnotation]; ok {
		dstRestoreData = make(restoreData)
		if err := json.Unmarshal([]byte(srcData), &dstRestoreData); err != nil {
			return err
		}
	}

	if err := convert(src, dst, nil); err != nil {
		return err
	}

	compare := new(SO)
	if err := unconvert(dst, compare, nil); err != nil {
		return err
	}

	for field, restoreData := range dstRestoreData {
		if restorer, ok := dstRestorer[field]; ok {
			if err := restorer.restore(restoreData, dst); err != nil {
				return err
			}
		}
	}

	srcRestoreData := make(restoreData)
	for field, restorer := range srcRestorer {
		restoreData, err := restorer.marshalState(src, compare)
		if err != nil {
			return err
		}
		if restoreData != nil {
			srcRestoreData[field] = restoreData
		}
	}

	if len(srcRestoreData) > 0 {
		restoreData, err := json.Marshal(srcRestoreData)
		if err != nil {
			return err
		}

		dstAnnotations := dst.GetAnnotations()
		if dstAnnotations == nil {
			dstAnnotations = map[string]string{}
		}
		dstAnnotations[utilconversion.DataAnnotation] = string(restoreData)
		dst.SetAnnotations(dstAnnotations)
	}

	return nil
}
