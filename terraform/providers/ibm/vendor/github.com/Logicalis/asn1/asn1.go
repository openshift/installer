// Package asn1 implements encoding and decoding of ASN.1 data structures using
// both Basic Encoding Rules (BER) or its subset, the Distinguished Encoding
// Rules (BER).
//
// This package is highly inspired by the Go standard package "encoding/asn1"
// while supporting additional features such as BER encoding and decoding and
// ASN.1 CHOICE types.
//
// By default and for convenience the package uses DER for encoding and BER for
// decoding. However it's possible to use a Context object to set the desired
// encoding and decoding rules as well other options.
//
// Restrictions:
//
// - BER allows STRING types, such as OCTET STRING and BIT STRING, to be
// encoded as constructed types containing inner elements that should be
// concatenated to form the complete string. The package does not support that,
// but in the future decoding of constructed strings should be included.
package asn1

// TODO add a mechanism for extendability
// TODO proper checking of the constructed flag
// TODO support for constructed encoding and decoding of string types in BER

import (
	"fmt"
	"reflect"
)

// Internal constants.
const (
	tagKey = "asn1"
)

// Encode returns the DER encoding of obj. Encode uses a default Context and
// it's equivalent to:
//
//	data, err = asn1.NewContext().Encode(obj)
//
func Encode(obj interface{}) (data []byte, err error) {
	ctx := NewContext()
	return ctx.EncodeWithOptions(obj, "")
}

// EncodeWithOptions returns the DER encoding of obj using additional options.
// EncodeWithOptions uses a default Context and it's equivalent to:
//
//	data, err = asn1.NewContext().EncodeWithOptions(obj, options)
//
func EncodeWithOptions(obj interface{}, options string) (data []byte, err error) {
	ctx := NewContext()
	return ctx.EncodeWithOptions(obj, options)
}

// Decode parses the given BER data into obj. The argument obj should be a
// reference to the value that will hold the parsed data. Decode uses a
// default Context and is equivalent to:
//
//	rest, err := asn1.NewContext().Decode(data, &obj)
//
func Decode(data []byte, obj interface{}) (rest []byte, err error) {
	ctx := NewContext()
	return ctx.DecodeWithOptions(data, obj, "")
}

// DecodeWithOptions parses the given BER data into obj using the additional
// options. The argument obj should be a reference to the value that will hold
// the parsed data. Decode uses a default Context and is equivalent to:
//
//	rest, err := asn1.NewContext().DecodeWithOptions(data, &obj, options)
//
func DecodeWithOptions(data []byte, obj interface{}, options string) (rest []byte, err error) {
	ctx := NewContext()
	return ctx.DecodeWithOptions(data, obj, options)
}

// ParseError is returned by the package to indicate that the given data is
// invalid.
type ParseError struct {
	Msg string
}

// Error returns the error message of a ParseError.
func (e *ParseError) Error() string {
	return e.Msg
}

// parseError allocates a new ParseError.
func parseError(msg string, args ...interface{}) *ParseError {
	return &ParseError{fmt.Sprintf(msg, args...)}
}

// SyntaxError is returned by the package to indicate that the given value or
// struct is invalid.
type SyntaxError struct {
	Msg string
}

// Error returns the error message of a SyntaxError.
func (e *SyntaxError) Error() string {
	return e.Msg
}

// syntaxError allocates a new ParseError,
func syntaxError(msg string, args ...interface{}) *SyntaxError {
	return &SyntaxError{fmt.Sprintf(msg, args...)}
}

// setDefaultValue sets a reflected value to its default value based on the
// field options.
func (ctx *Context) setDefaultValue(value reflect.Value, opts *fieldOptions) error {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value.SetInt(int64(*opts.defaultValue))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value.SetUint(uint64(*opts.defaultValue))

	default:
		return syntaxError("default value is only allowed to integers")
	}
	return nil
}

// newDefaultValue creates a new reflected value and sets it to its default value.
func (ctx *Context) newDefaultValue(objType reflect.Type, opts *fieldOptions) (reflect.Value, error) {
	value := reflect.New(objType).Elem()
	if opts.defaultValue == nil {
		return value, nil
	}
	return value, ctx.setDefaultValue(value, opts)
}
