package asn1

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
)

// Pre-calculated types for convenience
var (
	bigIntType = reflect.TypeOf((*big.Int)(nil))
	oidType    = reflect.TypeOf(Oid{})
	nullType   = reflect.TypeOf(Null{})
)

/*
 * Basic encoders and decoders
 */

// A function that encodes data.
type encoderFunction func(reflect.Value) ([]byte, error)

// A function that decodes data.
type decoderFunction func([]byte, reflect.Value) error

func (ctx *Context) encodeBool(value reflect.Value) ([]byte, error) {
	if value.Kind() != reflect.Bool {
		return nil, wrongType(reflect.Bool.String(), value)
	}
	if value.Bool() {
		return []byte{0xff}, nil
	}
	return []byte{0x00}, nil
}

func (ctx *Context) decodeBool(data []byte, value reflect.Value) error {
	// TODO check value type
	if !ctx.der.decoding {
		boolValue := parseBigInt(data).Cmp(big.NewInt(0)) != 0
		value.SetBool(boolValue)
		return nil
	}

	// DER is more restrict regarding valid booleans
	if len(data) == 1 {
		switch data[0] {
		case 0x00:
			value.SetBool(false)
			return nil
		case 0xff:
			value.SetBool(true)
			return nil
		}
	}
	return parseError("invalid bool value")
}

func (ctx *Context) encodeBigInt(value reflect.Value) ([]byte, error) {
	num, ok := value.Interface().(*big.Int)
	if !ok {
		return nil, wrongType(bigIntType.String(), value)
	}
	if num == nil {
		return []byte{0x00}, nil
	}
	var buf []byte
	if num.Sign() >= 0 {
		buf = append([]byte{0x00}, num.Bytes()...)
	} else {
		// Set absolute value
		convertedNum := big.NewInt(0)
		convertedNum.Abs(num)
		// Subtract One
		convertedNum.Sub(convertedNum, big.NewInt(1))
		// Invert bytes
		buf = convertedNum.Bytes()
		for i, b := range buf {
			buf[i] = ^b
		}
		buf = append([]byte{0xff}, buf...)
	}
	return removeIntLeadingBytes(buf), nil
}

func (ctx *Context) decodeBigInt(data []byte, value reflect.Value) error {
	// TODO check value type
	err := checkInt(ctx, data)
	if err != nil {
		return err
	}
	i := parseBigInt(data)
	value.Set(reflect.ValueOf(i))
	return nil
}

func (ctx *Context) encodeInt(value reflect.Value) ([]byte, error) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return nil, wrongType("signed integer", value)
	}
	n := value.Int()
	buf := make([]byte, 8)
	for i := range buf {
		shift := 8 * uint(len(buf)-i-1)
		mask := int64(0xff) << shift
		buf[i] = byte((n & mask) >> shift)
	}
	return removeIntLeadingBytes(buf), nil
}

func (ctx *Context) decodeInt(data []byte, value reflect.Value) error {
	// TODO check value type
	err := checkInt(ctx, data)
	if err != nil {
		return err
	}
	if len(data) > 8 {
		return parseError("integer too large for Go type '%s'", value.Type())
	}
	// Sign extend the value
	extensionByte := byte(0x00)
	if len(data) > 0 && data[0]&0x80 != 0 {
		extensionByte = byte(0xff)
	}
	extension := make([]byte, 8-len(data))
	for i := range extension {
		extension[i] = extensionByte
	}
	data = append(extension, data...)
	// Decode binary
	num := int64(0)
	for i := 0; i < len(data); i++ {
		num <<= 8
		num |= int64(data[i])
	}
	value.SetInt(num)
	return nil
}

func (ctx *Context) encodeUint(value reflect.Value) ([]byte, error) {
	switch value.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	default:
		return nil, wrongType("unsigned integer", value)
	}
	n := value.Uint()
	buf := make([]byte, 9)
	for i := range buf {
		shift := 8 * uint(len(buf)-i-1)
		mask := uint64(0xff) << shift
		buf[i] = byte((n & mask) >> shift)
	}
	return removeIntLeadingBytes(buf), nil
}

func (ctx *Context) decodeUint(data []byte, value reflect.Value) error {
	// TODO check value type
	err := checkInt(ctx, data)
	if err != nil {
		return err
	}
	if len(data) > 8 {
		return parseError("integer too large for Go type '%s'", value.Type())
	}
	if len(data) > 0 && data[0]&0x80 != 0 {
		return parseError("negative integer can't be assigned to Go type '%s'", value.Type())
	}
	num := uint64(0)
	for i := 0; i < len(data); i++ {
		num <<= 8
		num |= uint64(data[i])
	}
	value.SetUint(num)
	return nil
}

func (ctx *Context) encodeOctetString(value reflect.Value) ([]byte, error) {
	// Check type
	kind := value.Kind()
	if !(kind == reflect.Array || kind == reflect.Slice) &&
		value.Type().Elem().Kind() == reflect.Uint8 {
		// Invalid type or element type
		return nil, wrongType("array or slice of bytes", value)
	}
	if kind == reflect.Slice {
		return value.Interface().([]byte), nil
	}
	data := make([]byte, value.Len())
	for i := 0; i < value.Len(); i++ {
		data[i] = value.Index(i).Interface().(byte)
	}
	return data, nil
}

func (ctx *Context) decodeOctetString(data []byte, value reflect.Value) error {
	// Check type
	kind := value.Kind()
	if !(kind == reflect.Array || kind == reflect.Slice) &&
		value.Type().Elem().Kind() == reflect.Uint8 {
		// Invalid type or element type
		return wrongType("array or slice of bytes", value)
	}
	if value.Kind() == reflect.Array {
		// Check array length
		if len(data) != value.Len() {
			t := fmt.Sprintf("[%d]uint8", value.Len())
			return wrongType(t, value)
		}
		// Get reference to the array as a slice
		dest := value.Slice(0, value.Len()).Interface().([]byte)
		// Copy data
		copy(dest, data)
	} else {
		// Set value with a copy of the array data
		value.Set(reflect.ValueOf(append([]byte{}, data...)))
	}
	return nil
}

func (ctx *Context) encodeString(value reflect.Value) ([]byte, error) {
	if value.Kind() != reflect.String {
		return nil, wrongType(reflect.String.String(), value)
	}
	return []byte(value.String()), nil
}

func (ctx *Context) decodeString(data []byte, value reflect.Value) error {
	// TODO check value type
	s := string(data)
	value.SetString(s)
	return nil
}

/*
 * Custom types
 */

// Oid is used to encode and decode ASN.1 OBJECT IDENTIFIERs.
type Oid []uint

// Cmp returns zero if both Oids are the same, a negative value if oid
// lexicographically precedes other and a positive value otherwise.
func (oid Oid) Cmp(other Oid) int {
	for i, n := range oid {
		if i >= len(other) {
			return 1
		}
		if n != other[i] {
			return int(n) - int(other[i])
		}
	}
	return len(oid) - len(other)
}

// String returns the dotted representation of oid.
func (oid Oid) String() string {
	if len(oid) == 0 {
		return ""
	}
	s := fmt.Sprintf(".%d", oid[0])
	for i := 1; i < len(oid); i++ {
		s += fmt.Sprintf(".%d", oid[i])
	}
	return s
}

func (ctx *Context) encodeOid(value reflect.Value) ([]byte, error) {
	// Check values
	oid, ok := value.Interface().(Oid)
	if !ok {
		return nil, wrongType(oidType.String(), value)
	}

	value1 := uint(0)
	if len(oid) >= 1 {
		value1 = oid[0]
		if value1 > 2 {
			return nil, parseError("invalid value for first element of OID: %d", value1)
		}
	}

	value2 := uint(0)
	if len(oid) >= 2 {
		value2 = oid[1]
		if value2 > 39 {
			return nil, parseError("invalid value for first element of OID: %d", value2)
		}
	}

	bytes := []byte{byte(40*value1 + value2)}
	for i := 2; i < len(oid); i++ {
		bytes = append(bytes, encodeMultiByteTag(oid[i])...)
	}
	return bytes, nil
}

func (ctx *Context) decodeOid(data []byte, value reflect.Value) error {
	// TODO check value type
	if len(data) == 0 {
		value.Set(reflect.ValueOf(Oid{}))
		return nil
	}

	value1 := uint(data[0] / 40)
	value2 := uint(data[0]) - 40*value1
	oid := Oid{value1, value2}

	reader := bytes.NewBuffer(data[1:])
	for reader.Len() > 0 {
		valueN, err := decodeMultiByteTag(reader)
		if err != nil {
			return parseError("invalid value element in Object Identifier")
		}
		oid = append(oid, valueN)
	}

	value.Set(reflect.ValueOf(oid))
	return nil
}

// Null is used to encode and decode ASN.1 NULLs.
type Null struct{}

func (ctx *Context) encodeNull(value reflect.Value) ([]byte, error) {
	_, ok := value.Interface().(Null)
	if !ok {
		return nil, wrongType(nullType.String(), value)
	}
	return []byte{}, nil
}

func (ctx *Context) decodeNull(data []byte, value reflect.Value) error {
	// TODO check value type
	_, ok := value.Interface().(Null)
	if !ok {
		return syntaxError("invalid type: %s", value.Type())
	}
	if len(data) != 0 {
		return parseError("invalid data for Null type")
	}
	return nil
}

/*
 * Helper functions
 */

func checkInt(ctx *Context, data []byte) error {
	if ctx.der.decoding {
		if len(data) >= 2 {
			if data[0] == 0xff || data[0] == 0x00 {
				if data[0]&0x80 == data[1]&0x80 {
					return parseError("integer not encoded in the short form")
				}
			}
		}
	}
	return nil
}

func parseBigInt(data []byte) *big.Int {
	data = append([]byte{}, data...)
	neg := false
	if data[0]&0x80 != 0 {
		neg = true
		for i, b := range data {
			data[i] = ^b
		}
	}
	i := new(big.Int).SetBytes(data)
	if neg {
		i = i.Add(i, big.NewInt(1)).Neg(i)
	}
	return i
}

func removeIntLeadingBytes(buf []byte) []byte {
	start := 0
	for start < len(buf)-1 {
		// Removes a leading byte when the first NINE bits are the same

		// If the first 8 bits are not the same, skip
		if buf[start] != 0x00 && buf[start] != 0xff {
			break
		}

		// And if the 9th bit (1st bit of the second byte) is not the same as the previous 8 bits, also skip
		if (buf[start] & 0x80) != (buf[start+1] & 0x80) {
			break
		}

		// Remove a leading byte
		start++
	}
	return buf[start:]
}

func wrongType(typeName string, value reflect.Value) error {
	return syntaxError(
		"invalid Go type '%s' found when expecting '%s'",
		value.Type(), typeName)
}
