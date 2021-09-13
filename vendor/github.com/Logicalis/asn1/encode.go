package asn1

import (
	"reflect"
	"sort"
	"unicode"
)

// Encode returns the ASN.1 encoding of obj.
//
// See (*Context).EncodeWithOptions() for further details.
func (ctx *Context) Encode(obj interface{}) (data []byte, err error) {
	return ctx.EncodeWithOptions(obj, "")
}

// EncodeWithOptions returns the ASN.1 encoding of obj using additional options.
//
// See (*Context).DecodeWithOptions() for further details regarding types and
// options.
func (ctx *Context) EncodeWithOptions(obj interface{}, options string) (data []byte, err error) {

	opts, err := parseOptions(options)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(obj)
	raw, err := ctx.encode(value, opts)
	if err != nil {
		return
	}
	data, err = raw.encode()
	return
}

// Main encode function
func (ctx *Context) encode(value reflect.Value, opts *fieldOptions) (*rawValue, error) {

	// Skip the interface type
	switch value.Kind() {
	case reflect.Interface:
		value = value.Elem()
	}

	// If a value is missing the default value is used
	empty := isEmpty(value)
	if opts.defaultValue != nil {
		if empty && !ctx.der.encoding {
			defaultValue, err := ctx.newDefaultValue(value.Type(), opts)
			if err != nil {
				return nil, err
			}
			value = defaultValue
			empty = false
		}
	}

	// Encode data
	raw, err := ctx.encodeValue(value, opts)
	if err != nil {
		return nil, err
	}

	// Since the empty flag is already calculated, check if it's optional
	if (opts.optional || opts.defaultValue != nil) && empty {
		return nil, nil
	}

	// Modify the data generated based on the given tags
	raw, err = ctx.applyOptions(value, raw, opts)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (ctx *Context) encodeValue(value reflect.Value, opts *fieldOptions) (raw *rawValue, err error) {

	raw = &rawValue{}
	encoder := encoderFunction(nil)

	// Special types:
	objType := value.Type()
	switch objType {
	case bigIntType:
		raw.Tag = tagInteger
		encoder = ctx.encodeBigInt
	case oidType:
		raw.Tag = tagOid
		encoder = ctx.encodeOid
	case nullType:
		raw.Tag = tagNull
		encoder = ctx.encodeNull
	}

	if encoder == nil {
		// Generic types:
		switch value.Kind() {
		case reflect.Bool:
			raw.Tag = tagBoolean
			encoder = ctx.encodeBool

		case reflect.String:
			raw.Tag = tagOctetString
			encoder = ctx.encodeString

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			raw.Tag = tagInteger
			encoder = ctx.encodeInt

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			raw.Tag = tagInteger
			encoder = ctx.encodeUint

		case reflect.Struct:
			raw.Tag = tagSequence
			raw.Constructed = true
			encoder = ctx.encodeStruct
			if opts.set {
				encoder = ctx.encodeStructAsSet
			}

		case reflect.Array, reflect.Slice:
			if objType.Elem().Kind() == reflect.Uint8 {
				raw.Tag = tagOctetString
				encoder = ctx.encodeOctetString
			} else {
				raw.Tag = tagSequence
				raw.Constructed = true
				encoder = ctx.encodeSlice
			}
		}
	}

	if encoder == nil {
		return nil, syntaxError("invalid Go type: %s", value.Type())
	}
	raw.Content, err = encoder(value)
	return
}

// applyOptions modifies a raw value based on the given options.
func (ctx *Context) applyOptions(value reflect.Value, raw *rawValue, opts *fieldOptions) (*rawValue, error) {

	// Change sequence to set
	if opts.set {
		if raw.Class != classUniversal || raw.Tag != tagSequence {
			return nil, syntaxError("Go type '%s' does not accept the flag 'set'", value.Type())
		}
		raw.Tag = tagSet
	}

	// Check if this type is an Asn.1 choice
	if opts.choice != nil {
		entry, err := ctx.getChoiceByType(*opts.choice, value.Type())
		if err != nil {
			return nil, err
		}
		raw, err = ctx.applyOptions(value, raw, entry.opts)
		raw.Class = entry.class
		raw.Tag = entry.tag
	}

	// Add an enclosing tag
	if opts.explicit {
		if opts.tag == nil {
			return nil, syntaxError(
				"invalid flag 'explicit' without tag on Go type '%s'",
				value.Type())
		}
		content, err := raw.encode()
		if err != nil {
			return nil, err
		}
		raw = &rawValue{}
		raw.Constructed = true
		raw.Content = content
	}

	// Change tag
	if opts.tag != nil {
		raw.Class = classContextSpecific
		raw.Tag = uint(*opts.tag)
	}
	// Change class
	if opts.universal {
		raw.Class = classUniversal
	}
	if opts.application {
		raw.Class = classApplication
	}

	// Use the indefinite length encoding
	if opts.indefinite {
		if !raw.Constructed {
			return nil, syntaxError(
				"invalid flag 'indefinite' on Go type: %s",
				value.Type())
		}
		raw.Indefinite = true
	}

	return raw, nil
}

// isEmpty checks is a value is empty.
func isEmpty(value reflect.Value) bool {
	defaultValue := reflect.Zero(value.Type())
	return reflect.DeepEqual(value.Interface(), defaultValue.Interface())
}

// isFieldExported checks is the field name starts with a capital letter.
func isFieldExported(field reflect.StructField) bool {
	return unicode.IsUpper([]rune(field.Name)[0])
}

// getRawValuesFromFields encodes each valid field ofa struct value and returns
// a slice of raw values.
func (ctx *Context) getRawValuesFromFields(value reflect.Value) ([]*rawValue, error) {
	// Encode each child to a raw value
	children := []*rawValue{}
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldStruct := value.Type().Field(i)
		// Ignore field that are not exported (that starts with lowercase)
		if isFieldExported(fieldStruct) {
			tag := fieldStruct.Tag.Get(tagKey)
			opts, err := parseOptions(tag)
			if err != nil {
				return nil, err
			}
			raw, err := ctx.encode(fieldValue, opts)
			if err != nil {
				return nil, err
			}
			children = append(children, raw)
		}
	}
	return children, nil
}

// encodeRawValues is a helper function to encode raw value in sequence.
func (ctx *Context) encodeRawValues(values ...*rawValue) ([]byte, error) {
	content := []byte{}
	for _, raw := range values {
		buf, err := raw.encode()
		if err != nil {
			return nil, err
		}
		content = append(content, buf...)
	}
	return content, nil
}

// encodeStruct encodes structs fields in order.
func (ctx *Context) encodeStruct(value reflect.Value) ([]byte, error) {
	// Encode each child to a raw value
	children, err := ctx.getRawValuesFromFields(value)
	if err != nil {
		return nil, err
	}
	return ctx.encodeRawValues(children...)
}

// encodeStructAsSet works similarly to encodeStruct, but in Der mode the
// fields are encoded in ascending order of their tags.
func (ctx *Context) encodeStructAsSet(value reflect.Value) ([]byte, error) {
	// Encode each child to a raw value
	children, err := ctx.getRawValuesFromFields(value)
	if err != nil {
		return nil, err
	}
	// Sort if necessary
	if ctx.der.encoding {
		sort.Sort(rawValueSlice(children))
	}
	return ctx.encodeRawValues(children...)
}

// encodeSlice encodes a slice or array as a sequence of values.
func (ctx *Context) encodeSlice(value reflect.Value) ([]byte, error) {
	content := []byte{}
	for i := 0; i < value.Len(); i++ {
		itemValue := value.Index(i)
		childBytes, err := ctx.EncodeWithOptions(itemValue.Interface(), "")
		if err != nil {
			return nil, err
		}
		content = append(content, childBytes...)
	}
	return content, nil
}
