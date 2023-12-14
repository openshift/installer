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
	hashutils "k8s.io/kubernetes/pkg/util/hash"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

// RestorerFor holds all field restorers for a given type T.
type RestorerFor[T metav1.Object] map[string]fieldRestorerFor[T]

type fieldRestorerFor[T any] interface {
	marshalState(T, T) (json.RawMessage, error)
	restore(json.RawMessage, T) error
}

// HashedFieldRestorer restores a field to its original pre-conversion state
// only if it was not modified while converted. It does this by doing a full
// round-trip conversion before saving its state, and storing the hash of an
// unmodified round-trip converted object.
type HashedFieldRestorer[T any, F any] struct {
	// GetField returns the field to be restored
	GetField func(T) *F

	// FilterField returns a modified copy of the field to be used for comparison. i.e. filtered fields will be ignored in comparison.
	FilterField func(*F) *F

	// RestoreField restores the field to its original pre-conversion state
	RestoreField func(*F, *F)
}

var _ fieldRestorerFor[any] = HashedFieldRestorer[any, any]{}

//nolint:unused
type hashedFieldRestoreState struct {
	Hash []byte          `json:"h,omitempty"`
	Data json.RawMessage `json:"d,omitempty"`
}

//nolint:unused
func (r HashedFieldRestorer[T, F]) getHash(obj T) []byte {
	f := r.GetField(obj)
	if r.FilterField != nil {
		f = r.FilterField(f)
	}

	table := crc64.MakeTable(crc64.ECMA)
	hash := crc64.New(table)
	hashutils.DeepHashObject(hash, f)
	return hash.Sum(make([]byte, 0, 8))
}

//nolint:unused
func (r HashedFieldRestorer[T, F]) marshalState(src, compare T) (json.RawMessage, error) {
	b, err := json.Marshal(r.GetField(src))
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
func (r HashedFieldRestorer[T, F]) restore(b json.RawMessage, dst T) error {
	var restoreState hashedFieldRestoreState
	if err := json.Unmarshal(b, &restoreState); err != nil {
		return err
	}

	if bytes.Equal(r.getHash(dst), restoreState.Hash) {
		var previous F
		if err := json.Unmarshal(restoreState.Data, &previous); err != nil {
			return err
		}
		r.RestoreField(&previous, r.GetField(dst))
	}

	return nil
}

// UnconditionalFieldRestorer restores a field to its previous value without checking to see if it was modified during conversion.
//

type UnconditionalFieldRestorer[T any, F any] struct {
	GetField func(T) *F
}

var _ fieldRestorerFor[any] = UnconditionalFieldRestorer[any, any]{}

//nolint:unused
type unconditionalFieldRestoreState json.RawMessage

//nolint:unused
func (r UnconditionalFieldRestorer[T, F]) marshalState(src, _ T) (json.RawMessage, error) {
	f := r.GetField(src)

	// We could do this with a comparable constraint on F, but that's too restrictive for an arbitrary struct.
	if reflect.ValueOf(f).Elem().IsZero() {
		return nil, nil
	}

	return json.Marshal(r.GetField(src))
}

//nolint:unused
func (r UnconditionalFieldRestorer[T, F]) restore(b json.RawMessage, dst T) error {
	return json.Unmarshal(b, r.GetField(dst))
}

// restoreData serialises state for a set of restorers.
type restoreData map[string]json.RawMessage

// ConvertAndRestore converts src to dst. During conversion is reads and uses any restore data in the src annotation, and writes new restore data to the dst annotation.
func ConvertAndRestore[S, D metav1.Object](src S, dst D, compare S, convert func(S, D, conversion.Scope) error, unconvert func(D, S, conversion.Scope) error, srcRestorer RestorerFor[S], dstRestorer RestorerFor[D]) error {
	// NOTE(mdbooth): passing compare in here is a wart. If there's any way
	// to redefine the signature of this function such that we can create a
	// compare object ourselves we should do it.

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
