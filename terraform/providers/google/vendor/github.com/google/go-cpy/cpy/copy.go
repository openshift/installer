// Copyright 2020, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cpy copies Go objects.
//
// This package provides a generic way of deep copying Go objects.
// It is designed with performance in mind and is suitable for production use.
// By design, it does not handle unexported fields. If such fields need copying,
// it is the responsibility of the user to provide a custom copy function
// to specify how a specific type should be copied.
//
// WARNING: This package's API is currently unstable and may change without
// warning. If this matters to you, you should wait until version
// 1.0 is released before using it.
package cpy

import (
	"fmt"
	"reflect"
	"sync"
)

// A Copier copies Go objects.
type Copier struct {
	// concFuncs is a list of copy functions that operate on concrete types.
	concFuncs []reflect.Value // []func(T) T

	// ifaceFuncs is a list of copy functions that operate on interface types.
	ifaceFuncs []reflect.Value // []func(I) I

	// lookupFuncCache is a mapping from reflect.Type
	// to a reflect.Value representing a function operating on that type.
	lookupFuncCache sync.Map // map[reflect.Type]reflect.Value

	// exportedFieldsCache is a mapping from reflect.Type
	// to a list of exported struct field indexes.
	exportedFieldsCache sync.Map // map[reflect.Type][]int

	// ignoreAllUnexported specifies whether to ignore unxported fields
	// as opposed to panicking when encountering them.
	ignoreAllUnexported bool
}

// New initializes a new Copier according to the provided options.
//
// Example usage:
//
//	// As a global variable.
//	var copier = cpy.New(
//		cpy.Func(proto.Clone),
//		cpy.Shallow(time.Time{}),
//	)
//
//	// Elsewhere in application code.
//	dst := copier.Copy(src)
//
// It is recommended that the Copier returned by New
// be stored in a global variable so that it can be reused.
func New(opts ...Option) *Copier {
	// Process options in reverse order since latter arguments take precedence.
	// Separate out functions that operate on concrete and interface types.
	var c Copier
	for i := len(opts) - 1; i >= 0; i-- {
		opt := opts[i]
		for _, fnc := range opt.copyFuncs {
			if fnc.Type().In(0).Kind() != reflect.Interface {
				c.concFuncs = append(c.concFuncs, fnc)
			} else {
				c.ifaceFuncs = append(c.ifaceFuncs, fnc)
			}
		}
		if opt.ignoreAllUnexported {
			c.ignoreAllUnexported = true
		}
	}

	// TODO: There is no obviously right behavior to take with regard to
	// unexported fields in a struct. Possible approaches:
	//
	//
	// 1. Panic if the type contains unexported fields or interfaces anywhere
	// in the type tree (that are not handled by a Func option).
	// Furthermore, the same concrete type must always be passed to Copy.
	// This ensures that a unit test exercising Copy will fail loudly if
	// an unexported field is ever introduced.
	//
	// It notifies the owner of the Copy call that their copy semantics may
	// be broken and needs adjustment to figure out what the right behavior
	// should be (whether it is actually safe to ignore the unexported fields)
	// or whether a custom Func option should be provided to properly copy
	// the type with unexported fields.
	//
	// The detriment of this approach is a higher probability that adding an
	// unexported field to a type causes some remote target to suddenly fail.
	// Also, there is a higher chance of false-positives where an unexported
	// field occurs in the type-tree, but is functionally never copied at
	// runtime since it (or some higher value in the type tree) is always zero.
	//
	//
	// 2. Panic when Copy tries to copy a struct value with unexported fields.
	// This approach avoids unnecessarily panicking just because an unexported
	// field is part of the type tree, but only does so at the very moment
	// an unexported field is being copied.
	//
	// The advantage of this approach relative to previous approach is the
	// reduction in false-positives and avoiding more cases where adding
	// unexported fields causing an unrelated target to fail. The disadvantage
	// of this approach is that uses of Copy with insufficient test coverage
	// may fail to detect that an unexported field does get copied,
	// resulting in a spurious panic at runtime in production code.
	//
	//
	// 3. Always ignore unexported fields.
	// This approach avoids spurious panics that may occur at runtime,
	// but may present silent data corruption. Fundamentally, it means that the
	// output value may not be identical to the input value.
	//
	//
	// 4. Shallow copy unexported fields, but deep copy exported fields.
	// It's possible to shallow copy unexported fields by shallow copying the
	// entire struct. While this approach avoids any panics, it does causes
	// a inconsistency where exported fields are deep copied, but unexported
	// fields are not. Furthermore, it may not be safe to shallow copy the
	// unexported fields since they may contain mutexes or other data structures
	// that shouldn't be shallow copied.
	//
	//
	// 5. Always copy unexported fields (with the use of unsafe).
	// Not only does this require importing unsafe, it is semantically unsafe
	// since we have no guarantees whether the unexported fields of some remote
	// type an even safe to copy. This is the least attractive option.
	// An AllowUnexported option is palatable as it explicitly declares that
	// a specific type's unexported fields are safe to copy.
	//
	//
	// Since it is unclear what the default behavior should be with regard to
	// unexported fields, require that users specify IgnoreAllUnexported so
	// that it is obvious up front what the behavior is. This makes cpy mostly
	// backwards compatible with deepcopy, which it seeks to replace.
	//
	// See the discussion on cl/333563483 for more details.
	if !c.ignoreAllUnexported {
		panic("cpy.IgnoreAllUnexported must be specified; this requirement may change in the future")
	}

	return &c
}

// Copy copies v according to the Copier presets.
//
// Values are copied according to the following rules:
//
// • Zero values (according to reflect.Value.IsZero) are returned as is.
//
// • If the current type matches a Func provided to New,
// then the specialized copy function is used to copy the value.
// Precedence is given to Funcs that operate on concrete types, then
// Funcs that operate on interface types, then the behavior listed below.
// For Funcs operating on the same type, those passed later to New
// take precedence over any preceding Func arguments.
//
// • Pointers are copied by allocating a new value of the same type and
// recursively calling Copy on the pointed-at value.
//
// • Interfaces are copied by recursively calling Copy on the
// underlying source value and casting the resulting concrete value
// to the same interface type.
//
// • Arrays and slices are copied by making a new array or slice
// of the same type (with the same length and capacity for slices).
// Every element in the source is recursively copied by calling Copy and
// storing the result into the destination array or slice.
//
// • Maps are copied by making a new map of the same type and
// recursively calling Copy on each map key and map value in the source and
// storing the result into the newly made destination map.
//
// • Structs are copied by creating a new struct of the same type and
// recursively calling Copy for each field in the source and
// storing the result into the destination struct. It panics when trying
// to copy a struct type with unexported fields unless an IgnoreAllUnexported
// option was passed to New, in which case unexported fields are ignored.
// Alternatively, a custom Func may be specified to provide a specialized
// implemention of deep-copying for the type with unexported fields based
// on the exported API for that type.
//
// • Lastly, all other types (e.g., int, string, etc.) are shallow copied.
// Note that unsafe.Pointer and channels are shallow copied since there is
// no obvious behavior to use to deep copy such types.
//
// The output type is guaranteed to be the same as the input type.
// Copy will panic if that invariant is violated by a provided Func.
// Copy presently does not handle cycles in the value and will overflow.
func (c *Copier) Copy(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	return c.copy(reflect.ValueOf(v)).Interface()
}
func (c *Copier) copy(src reflect.Value) (dst reflect.Value) {
	t := src.Type()

	// Return zero values as is.
	if src.IsZero() {
		return src
	}

	// Check if there is a specialized copy function for this type.
	if fnc := c.lookupFunc(t); fnc.IsValid() {
		ft := fnc.Type().In(0)
		if ft.Kind() != reflect.Interface {
			if t == ft {
				dst = fnc.Call([]reflect.Value{src})[0]
			} else {
				dst = fnc.Call([]reflect.Value{makeAddr(src)})[0].Elem()
			}
		} else {
			if t.Implements(ft) {
				dst = fnc.Call([]reflect.Value{src.Convert(ft)})[0].Elem().Convert(t)
			} else {
				dst = fnc.Call([]reflect.Value{makeAddr(src).Convert(ft)})[0].Elem().Elem().Convert(t)
			}
		}
		return dst
	}

	// Deep copy pointers, interfaces, arrays, slices, maps, and structs.
	dst = src // shallow copy the value by default
	switch t.Kind() {
	case reflect.Ptr:
		dst = reflect.New(src.Elem().Type())
		dst.Elem().Set(c.copy(src.Elem()))
	case reflect.Interface:
		dst = c.copy(src.Elem()).Convert(t)
	case reflect.Array:
		dst = reflect.New(t).Elem()
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(c.copy(src.Index(i)))
		}
	case reflect.Slice:
		dst = reflect.MakeSlice(t, src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(c.copy(src.Index(i)))
		}
	case reflect.Map:
		dst = reflect.MakeMap(t)
		for iter := src.MapRange(); iter.Next(); {
			dst.SetMapIndex(c.copy(iter.Key()), c.copy(iter.Value()))
		}
	case reflect.Struct:
		dst = reflect.New(t).Elem()
		for _, i := range c.exportedFields(src.Type()) {
			dst.Field(i).Set(c.copy(src.Field(i)))
		}
	}
	return dst
}

// lookupFunc returns a custom copy function for the provided type
// if there is one. Otherwise, it returns an invalid value.
func (c *Copier) lookupFunc(t reflect.Type) reflect.Value {
	v, ok := c.lookupFuncCache.Load(t)
	if !ok {
		fnc := c.lookupFuncSlow(t)
		v, _ = c.lookupFuncCache.LoadOrStore(t, fnc)
	}
	return v.(reflect.Value)
}
func (c *Copier) lookupFuncSlow(t reflect.Type) reflect.Value {
	// Check for exact match with functions operating on concrete types.
	for _, t := range []reflect.Type{t, reflect.PtrTo(t)} {
		for _, fnc := range c.concFuncs {
			if t == fnc.Type().In(0) {
				return fnc
			}
		}
	}
	// Check for assignability to functions operating on interface types.
	for _, t := range []reflect.Type{t, reflect.PtrTo(t)} {
		for _, fnc := range c.ifaceFuncs {
			if strictImplements(t, fnc.Type().In(0)) {
				return fnc
			}
		}
	}
	return reflect.Value{}
}

// strictImplements is identical to reflect.Type.Implements,
// but reports false if the non-pointer version of t also implements ti.
//
// The purpose of this check is to ensure that we consistently call "func(I) I"
// with the value form of the reciever if it implements I,
// rather than an inconsistent mixture of both value and pointer receivers.
func strictImplements(t, ti reflect.Type) bool {
	if !t.Implements(ti) {
		return false
	}
	if t.Kind() == reflect.Ptr && t.Elem().Implements(ti) {
		return false
	}
	return true
}

// makeAddr returns the address of v if possible,
// otherwise it shallow copies v into a new instance and returns that.
func makeAddr(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v.Addr()
	}
	p := reflect.New(v.Type())
	p.Elem().Set(v)
	return p
}

// exportedFields returns a list of exported field indexes in struct t.
// This method caches the result since reflect.Type.Field is slow
// since every call always allocates reflect.Type.StructField.Index.
func (c *Copier) exportedFields(t reflect.Type) []int {
	v, ok := c.exportedFieldsCache.Load(t)
	if !ok {
		index := c.exportedFieldsSlow(t)
		v, _ = c.exportedFieldsCache.LoadOrStore(t, index)
	}
	return v.([]int)
}
func (c *Copier) exportedFieldsSlow(t reflect.Type) []int {
	index := make([]int, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath == "" {
			index = append(index, i) // record index of exported field
		} else if !c.ignoreAllUnexported {
			var name string
			if t.Name() != "" {
				// Named type with unexported fields.
				name = fmt.Sprintf("%q.%v", t.PkgPath(), t.Name()) // e.g., "path/to/package".MyType
			} else {
				// Unnamed type with unexported fields.
				name = fmt.Sprintf("%q.(%v)", f.PkgPath, t.String()) // e.g., "path/to/package".(struct { a int })
			}
			panic(fmt.Sprintf("unable to copy unexported field: %v.%v", name, f.Name))
		}
	}
	return index
}

// Option is an option that configures a Copier.
// An option must be obtained using a constructor (e.g., Func or Shallow).
type Option option

// Keep the exact representation of Option opaque.
// We may change it to be an interface in the future.
type option struct {
	copyFuncs           []reflect.Value
	ignoreAllUnexported bool
}

// Func provides specialized copy behavior for specific types.
//
// The copy function f must be a function "func(T) T",
// where T must be a pointer, interface, array, slice, map, or struct;
// otherwise it will panic. If T is an interface type,
// it must have at least one method, otherwise Func will panic.
// Futhermore, the function must return a concrete type
// identical to the input type, otherwise Copier.Copy will panic.
//
// The Func will be used if the current type of the value being copied
// exactly matches the input type (for concrete types) or if it
// implements the input type (for interface types).
// Both the type itself (e.g., T) and a pointer to the type (e.g., *T)
// are checked when evaluating whether a given Func can be used.
//
// Example usage:
//
//	cpy.Func(proto.Clone)
//
// This option specifies that proto.Clone is used to copy all types that
// are assignable to the proto.Message interface.
func Func(fn interface{}) Option {
	v := reflect.ValueOf(fn)
	if !v.IsValid() || v.Kind() != reflect.Func ||
		v.Type().NumIn() != 1 || v.Type().NumOut() != 1 || v.Type().In(0) != v.Type().Out(0) || v.Type().IsVariadic() {
		panic(fmt.Sprintf("cpy.Func: input function %T must be a func(T) T", fn))
	}
	if t := v.Type().In(0); !validKind(t.Kind()) {
		panic(fmt.Sprintf("cpy.Func: input type %v must be a pointer, interface, array, slice, map, or struct", t))
	}
	if t := v.Type().In(0); t.Kind() == reflect.Interface && t.NumMethod() == 0 {
		panic(fmt.Sprintf("cpy.Func: interface type %v must have methods", t))
	}
	return Option{copyFuncs: []reflect.Value{v}}
}

// Shallow specifies that the provided type should be shallow copied.
// The provided type must be a pointer, interface, array, slice, map, or struct;
// otherwise it will panic.
//
// Shallow is equivalent to:
//
//	cpy.Func(func(v V) V { return v })
//
// Example usage:
//
//	cpy.Shallow(time.Time{})
//
// Since the fields of time.Time are unexported, the default behavior
// of Copier.Copy will functionally avoid copying the time value.
// This option specifies that time.Time is a value that is safe to shallow copy.
func Shallow(typs ...interface{}) Option {
	var opt Option
	for _, typ := range typs {
		t := reflect.TypeOf(typ)
		if t == nil || !validKind(t.Kind()) {
			panic(fmt.Sprintf("cpy.Shallow: input type %v must be a pointer, interface, array, slice, map, or struct", t))
		}
		v := reflect.MakeFunc(
			reflect.FuncOf([]reflect.Type{t}, []reflect.Type{t}, false), // func(T) T
			func(in []reflect.Value) []reflect.Value { return in },      // shallow copy
		)
		opt.copyFuncs = append(opt.copyFuncs, v)
	}
	return opt
}

// TODO: Add AllowUnexported(typs ...interface{}) option.
// TODO: Add IgnoreUnexported(typs ...interface{}) option.

// IgnoreAllUnexported specifies that Copy should ignore all unexported fields
// as opposed to panicking when encountering an unexported field.
func IgnoreAllUnexported() Option {
	return Option{ignoreAllUnexported: true}
}

func validKind(k reflect.Kind) bool {
	switch k {
	case reflect.Ptr, reflect.Interface, reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return true
	default:
		return false
	}
}
