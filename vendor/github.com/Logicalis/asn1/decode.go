package asn1

import (
	"bytes"
	"io"
	"reflect"
	"sort"
)

// Expected values
type expectedElement struct {
	class   uint
	tag     uint
	decoder decoderFunction
}

// Expected values for fields
type expectedFieldElement struct {
	expectedElement
	value reflect.Value
	opts  *fieldOptions
}

// Decode parses the given data into obj. The argument obj should be a reference
// to the value that will hold the parsed data.
//
// See (*Context).DecodeWithOptions() for further details.
//
func (ctx *Context) Decode(data []byte, obj interface{}) (rest []byte, err error) {
	return ctx.DecodeWithOptions(data, obj, "")
}

// DecodeWithOptions parses the given data into obj using the additional
// options. The argument obj should be a reference to the value that will hold
// the parsed data.
//
// It uses the reflect package to inspect obj and because of that only exported
// struct fields (those that start with a capital letter) are considered.
//
// The Context object defines the decoding rules (BER or DER) and the types
// available for CHOICE types.
//
// The asn1 package maps Go types to ASN.1 data structures. The package also
// provides types to specific ASN.1 data structures. The default mapping is
// shown in the table below:
//
//	Go type                | ASN.1 universal tag
//	-----------------------|--------------------
//	bool                   | BOOLEAN
//	All int and uint types | INTEGER
//	*big.Int               | INTEGER
//	string                 | OCTET STRING
//	[]byte                 | OCTET STRING
//	asn1.Oid               | OBJECT INDETIFIER
//	asn1.Null              | NULL
//	Any array or slice     | SEQUENCE OF
//	Any struct             | SEQUENCE
//
// Arrays and slices are decoded using different rules. A slice is always
// appended while an array requires an exact number of elements, otherwise a
// ParseError is returned.
//
// The default mapping can be changed using options provided in the argument
// options (for the root value) or via struct tags for struct fields. Struct
// tags use the namei space "asn1".
//
// The available options for encoding and decoding are:
//
//	tag
//
// This option requires an numeric argument (ie: "tag:1") and indicates that a
// element is encoded and decoded as a context specific element with the given
// tag number. The context specific class can be overridden with the options
// "application" or "universal".
//
//	universal
//
// Sets the tag class to universal. Requires "tag".
//
//	application
//
// Sets the tag class to application. Requires "tag".
//
//	explicit
//
// Indicates the element is encoded with an enclosing tag. It's usually
// used in conjunction with "tag".
//
//	optional
//
// Indicates that an element can be suppressed.
//
// A missing element that is not marked with "optional" or "default" causes a
// ParseError to be returned during decoding. A missing element marked as
// optional is left untouched.
//
// During encoding, a zero value elements is suppressed from output if it's
// marked as optional.
//
//	default
//
// This option is handled similarly to the "optional" option but requires a
// numeric argument (ie: "default:1").
//
// A missing element that is marked with "default" is set to the given default
// value during decoding.
//
// A zero value marked with "default" is suppressed from output  when encoding
// is set to DER or is encoded with the given default value when encoding is set
// to BER.
//
//	indefinite
//
// This option is used only during encoding and causes a constructed element to
// be encoded using the indefinite form.
//
//	choice
//
// Indicates that an element can be of one of several types as defined by
// (*Context).AddChoice()
//
//	set
//
// Indicates that a struct, array or slice should be encoded and decoded as a
// SET instead of a SEQUENCE.
//
// It also affects the way that structs are encoded and decoded in DER. A struct
// marked with "set" has its fields always encoded in the ascending order of its
// tags, instead of following the order that the fields are defined in the
// struct.
//
// Similarly, a struct marked with "set" always enforces that same order when
// decoding in DER.
//
func (ctx *Context) DecodeWithOptions(data []byte, obj interface{}, options string) (rest []byte, err error) {

	opts, err := parseOptions(options)
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr, reflect.Interface:
		value = value.Elem()
	}

	if !value.CanSet() {
		return nil, syntaxError("go type '%s' is read-only", value.Type())
	}

	reader := bytes.NewBuffer(data)
	err = ctx.decode(reader, value, opts)
	if err != nil {
		return nil, err
	}

	return reader.Bytes(), nil
}

// Main decode function
func (ctx *Context) decode(reader io.Reader, value reflect.Value, opts *fieldOptions) error {

	// Parse an Asn.1 element
	raw, err := decodeRawValue(reader)
	if err != nil {
		return err
	}
	if ctx.der.decoding && raw.Indefinite {
		return parseError("indefinite length form is not supported by DER mode")
	}

	elem, err := ctx.getExpectedElement(raw, value.Type(), opts)
	if err != nil {
		return err
	}

	// And tag must match
	if raw.Class != elem.class || raw.Tag != elem.tag {
		ctx.log.Printf("%#v\n", opts)
		return parseError("expected tag (%d,%d) but found (%d,%d)",
			elem.class, elem.tag, raw.Class, raw.Tag)
	}

	return elem.decoder(raw.Content, value)
}

// getExpectedElement returns the expected element for a given type. raw is only
// used as hint when decoding choices.
// TODO: consider replacing raw for class and tag number.
func (ctx *Context) getExpectedElement(raw *rawValue, elemType reflect.Type, opts *fieldOptions) (elem expectedElement, err error) {

	// Get the expected universal tag and its decoder for the given Go type
	elem, err = ctx.getUniversalTag(elemType, opts)
	if err != nil {
		return
	}

	// Modify the expected tag and decoder function based on the given options
	if opts.tag != nil {
		elem.class = classContextSpecific
		elem.tag = uint(*opts.tag)
	}
	if opts.universal {
		elem.class = classUniversal
	}
	if opts.application {
		elem.class = classApplication
	}

	if opts.explicit {
		elem.decoder = func(data []byte, value reflect.Value) error {
			// Unset previous flags
			opts.explicit = false
			opts.tag = nil
			opts.application = false
			// Parse child
			reader := bytes.NewBuffer(data)
			return ctx.decode(reader, value, opts)
		}
		return
	}

	if opts.choice != nil {
		// Get the registered choices
		var entry choiceEntry
		entry, err = ctx.getChoiceByTag(*opts.choice, raw.Class, raw.Tag)
		if err != nil {
			return
		}

		// Get the decoder for the new value
		elem.class, elem.tag = raw.Class, raw.Tag
		elem.decoder = func(data []byte, value reflect.Value) error {
			// Allocate a new value and set to the current one
			nestedValue := reflect.New(entry.typ).Elem()
			err = entry.decoder(data, nestedValue)
			if err != nil {
				return err
			}
			value.Set(nestedValue)
			return nil
		}
	}

	// At this point a decoder function already be found
	if elem.decoder == nil {
		err = parseError("go type not supported '%s'", elemType)
	}
	return
}

// getUniversalTag maps an type to a Asn.1 universal type.
func (ctx *Context) getUniversalTag(objType reflect.Type, opts *fieldOptions) (elem expectedElement, err error) {

	elem.class = classUniversal

	// Special types:
	switch objType {
	case bigIntType:
		elem.tag = tagInteger
		elem.decoder = ctx.decodeBigInt
	case oidType:
		elem.tag = tagOid
		elem.decoder = ctx.decodeOid
	case nullType:
		elem.tag = tagNull
		elem.decoder = ctx.decodeNull
	default:
		// Generic types:
		elem = ctx.getUniversalTagByKind(objType, opts)
	}

	// Check options for universal types
	if opts.set {
		if elem.tag != tagSequence {
			err = syntaxError(
				"'set' cannot be used with Go type '%s'", objType)
		}
		elem.tag = tagSet
	}
	return
}

// getUniversalTagByKind uses type kind to defined the decoder.
func (ctx *Context) getUniversalTagByKind(objType reflect.Type, opts *fieldOptions) (elem expectedElement) {

	switch objType.Kind() {
	case reflect.Bool:
		elem.tag = tagBoolean
		elem.decoder = ctx.decodeBool

	case reflect.String:
		elem.tag = tagOctetString
		elem.decoder = ctx.decodeString

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		elem.tag = tagInteger
		elem.decoder = ctx.decodeInt

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		elem.tag = tagInteger
		elem.decoder = ctx.decodeUint

	case reflect.Struct:
		elem.tag = tagSequence
		elem.decoder = ctx.decodeStruct
		if opts.set {
			elem.decoder = ctx.decodeStructAsSet
		}

	case reflect.Array:
		if objType.Elem().Kind() == reflect.Uint8 {
			elem.tag = tagOctetString
			elem.decoder = ctx.decodeOctetString
		} else {
			elem.tag = tagSequence
			elem.decoder = ctx.decodeArray
		}

	case reflect.Slice:
		if objType.Elem().Kind() == reflect.Uint8 {
			elem.tag = tagOctetString
			elem.decoder = ctx.decodeOctetString
		} else {
			elem.tag = tagSequence
			elem.decoder = ctx.decodeSlice
		}
	}
	return
}

//
func (ctx *Context) getExpectedFieldElements(value reflect.Value) ([]expectedFieldElement, error) {
	expectedValues := []expectedFieldElement{}
	for i := 0; i < value.NumField(); i++ {
		if value.CanSet() {
			// Get field and options
			field := value.Field(i)
			opts, err := parseOptions(value.Type().Field(i).Tag.Get(tagKey))
			if err != nil {
				return nil, err
			}
			// Expand choices
			raw := &rawValue{}
			if opts.choice == nil {
				elem, err := ctx.getExpectedElement(raw, field.Type(), opts)
				if err != nil {
					return nil, err
				}
				expectedValues = append(expectedValues,
					expectedFieldElement{elem, field, opts})
			} else {
				entries, err := ctx.getChoices(*opts.choice)
				if err != nil {
					return nil, err
				}
				for _, entry := range entries {
					raw.Class = entry.class
					raw.Tag = entry.tag
					elem, err := ctx.getExpectedElement(raw, field.Type(), opts)
					if err != nil {
						return nil, err
					}
					expectedValues = append(expectedValues,
						expectedFieldElement{elem, field, opts})
				}
			}
		}
	}
	return expectedValues, nil
}

// getRawValuesFromBytes reads up to max values from the byte sequence.
func (ctx *Context) getRawValuesFromBytes(data []byte, max int) ([]*rawValue, error) {
	// Raw values
	rawValues := []*rawValue{}
	reader := bytes.NewBuffer(data)
	for i := 0; i < max; i++ {
		// Parse an Asn.1 element
		raw, err := decodeRawValue(reader)
		if err != nil {
			return nil, err
		}
		rawValues = append(rawValues, raw)
		if reader.Len() == 0 {
			return rawValues, nil
		}
	}
	return nil, parseError("too many items for Sequence")
}

// matchExpectedValues tries to decode a sequence of raw values based on the
// expected elements.
func (ctx *Context) matchExpectedValues(eValues []expectedFieldElement, rValues []*rawValue) error {
	// Try to match expected and raw values
	rIndex := 0
	for eIndex := 0; eIndex < len(eValues); eIndex++ {
		e := eValues[eIndex]
		// Using nil decoder to skip matched choices
		if e.decoder == nil {
			continue
		}

		missing := true
		if rIndex < len(rValues) {
			raw := rValues[rIndex]
			if e.class == raw.Class && e.tag == raw.Tag {
				err := e.decoder(raw.Content, e.value)
				if err != nil {
					return err
				}
				// Mark as found and advance raw values index
				missing = false
				rIndex++
				// Remove other options for the matched choice
				if e.opts.choice != nil {
					for i := eIndex + 1; i < len(eValues); i++ {
						c := eValues[i].opts.choice
						if c != nil && *c == *e.opts.choice {
							eValues[i].decoder = nil
						}
					}
				}
			}
		}

		if missing {
			if err := ctx.setMissingFieldValue(e); err != nil {
				return err
			}
		}
	}
	return nil
}

// setMissingFieldValue uses opts values to set the default value.
func (ctx *Context) setMissingFieldValue(e expectedFieldElement) error {
	if e.opts.optional || e.opts.choice != nil {
		return nil
	}
	if e.opts.defaultValue != nil {
		err := ctx.setDefaultValue(e.value, e.opts)
		if err != nil {
			return err
		}
		return nil
	}
	return parseError("missing value for [%d %d]", e.class, e.tag)
}

// decodeStruct decodes struct fields in order
func (ctx *Context) decodeStruct(data []byte, value reflect.Value) error {

	expectedValues, err := ctx.getExpectedFieldElements(value)
	if err != nil {
		return err
	}

	rawValues, err := ctx.getRawValuesFromBytes(data, len(expectedValues))
	if err != nil {
		return err
	}

	return ctx.matchExpectedValues(expectedValues, rawValues)
}

// Decode a struct as an Asn.1 Set.
//
// The order doesn't matter for set. However DER dictates that a Set should be
// encoded in the ascending order of the tags. So when decoding with DER, we
// simply do not sort the raw values and use them in their natural order.
func (ctx *Context) decodeStructAsSet(data []byte, value reflect.Value) error {

	// Get the expected values
	expectedElements, err := ctx.getExpectedFieldElements(value)
	if err != nil {
		return err
	}
	sort.Sort(expectedFieldElementSlice(expectedElements))

	// Check duplicated tags
	for i := 1; i < len(expectedElements); i++ {
		curr := expectedElements[i]
		prev := expectedElements[i-1]
		if curr.class == prev.class &&
			curr.tag == prev.tag {
			return syntaxError("duplicated tag (%d,%d)", curr.class, curr.tag)
		}
	}

	// Get the raw values
	rawValues, err := ctx.getRawValuesFromBytes(data, len(expectedElements))
	if err != nil {
		return err
	}
	if !ctx.der.decoding {
		sort.Sort(rawValueSlice(rawValues))
	}

	return ctx.matchExpectedValues(expectedElements, rawValues)
}

// decodeSlice decodes a SET(OF) as a slice
func (ctx *Context) decodeSlice(data []byte, value reflect.Value) error {
	slice := reflect.New(value.Type()).Elem()
	var err error
	for len(data) > 0 {
		elem := reflect.New(value.Type().Elem()).Elem()
		data, err = ctx.DecodeWithOptions(data, elem.Addr().Interface(), "")
		if err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, elem))
	}
	value.Set(slice)
	return nil
}

// decodeArray decodes a SET(OF) as an array
func (ctx *Context) decodeArray(data []byte, value reflect.Value) error {
	var err error
	for i := 0; i < value.Len(); i++ {
		if len(data) == 0 {
			return parseError("missing elements")
		}
		elem := reflect.New(value.Type().Elem()).Elem()
		data, err = ctx.DecodeWithOptions(data, elem.Addr().Interface(), "")
		if err != nil {
			return err
		}
		value.Index(i).Set(elem)
	}
	if len(data) > 0 {
		return parseError("too many elements")
	}
	return nil
}
