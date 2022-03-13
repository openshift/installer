package asn1

import (
	"bytes"
	"io"
	"io/ioutil"
	"strconv"
)

// ASN.1 class tags.
const (
	classUniversal       = 0x00
	classApplication     = 0x01
	classContextSpecific = 0x02
	classPrivate         = 0x03
)

// ASN.1 universal tag numbers.
const (
	tagEoc             = 0x00
	tagBoolean         = 0x01
	tagInteger         = 0x02
	tagBitString       = 0x03
	tagOctetString     = 0x04
	tagNull            = 0x05
	tagOid             = 0x06
	tagSequence        = 0x10
	tagSet             = 0x11
	tagPrintableString = 0x13
	tagT61String       = 0x14
	tagIA5String       = 0x16
	tagUtcTime         = 0x17
)

// Internal consts
const (
	intBits  = strconv.IntSize
	intBytes = intBits / 8
)

type rawValue struct {
	Class       uint
	Tag         uint
	Constructed bool
	Indefinite  bool
	Content     []byte
}

func (raw *rawValue) encode() ([]byte, error) {

	if raw == nil {
		return []byte{}, nil
	}

	buf, err := encodeIdentifier(raw)
	if err != nil {
		return nil, err
	}

	// Add length information and raw data
	if !raw.Indefinite {
		length := uint(len(raw.Content))
		buf = append(buf, encodeLength(length)...)
		buf = append(buf, raw.Content...)
	} else {
		// Indefinite length uses 0x80, data..., 0x00, 0x00
		if !raw.Constructed {
			return nil, syntaxError("indefinite length is only allowed to constructed types")
		}
		buf = append(buf, 0x80)
		buf = append(buf, raw.Content...)
		buf = append(buf, 0x00, 0x00)
	}

	return buf, nil
}

func encodeIdentifier(node *rawValue) ([]byte, error) {

	if node.Class > 0x03 {
		return nil, syntaxError("invalid class value: %d", node.Class)
	}

	identifier := []byte{0x00}

	// Class (bits 7 and 6) + primitive/constructed (1 bit) + tag (5 bits)
	identifier[0] += byte((node.Class & 0x03) << 6)

	// Primitive/constructed (bit 5)
	if node.Constructed {
		identifier[0] += byte(1 << 5)
	}

	// Tag (bits 4 to 0)
	if node.Tag <= 30 {
		identifier[0] += byte(0x1f & node.Tag)
	} else {
		identifier[0] += 0x1f
		identifier = append(identifier, encodeMultiByteTag(node.Tag)...)
	}
	return identifier, nil
}

func encodeMultiByteTag(tag uint) []byte {

	// A tag is encoded in a big endian sequence of octets each one holding a 7 bit value.
	// The most significant bit of each octet must be 1, with the exception of the last octet that must be zero.
	// Example:  1xxxxxxx 1xxxxxxx ... 0xxxxxxx

	// An int32 needs 5 octets and an int64 needs 10:
	bufLen := (intBits-1)/7 + 1
	buf := make([]byte, bufLen)

	for i := range buf {
		shift := uint(7 * (len(buf) - i - 1))
		mask := uint(0x7f << shift)
		buf[i] = byte((tag & mask) >> shift)
		// Only the last byte is not marked with 0x80
		if i != len(buf)-1 {
			buf[i] |= 0x80
		}
	}
	// Discard leading zero values
	return removeLeadingBytes(buf, 0x80)
}

func encodeLength(length uint) []byte {

	// The first bit indicates if length is encoded in a single byte
	if length < 0x80 {
		return []byte{byte(length)}
	}

	// Multi byte length follow the rules:
	// - First byte: 0x80 + N (number of following bytes where length is encoded).
	// - N bytes

	// A byte slice length is an int. So we just need at most 4 bytes
	buf := make([]byte, intBytes)
	for i := range buf {
		shift := uint((intBytes - i - 1) * 8)
		mask := uint(0xff << shift)
		buf[i] = byte((mask & length) >> shift)
	}

	// Ignore leading zeros
	buf = removeLeadingBytes(buf, 0x00)

	// Add leading byte with the number of following bytes
	buf = append([]byte{0x80 + byte(len(buf))}, buf...)
	return buf
}

func removeLeadingBytes(buf []byte, target byte) []byte {
	start := 0
	for start < len(buf)-1 && buf[start] == target {
		start++
	}
	return buf[start:]
}

func decodeRawValue(reader io.Reader) (*rawValue, error) {

	class, tag, constructed, err := decodeIdentifier(reader)
	if err != nil {
		return nil, err
	}

	length, indefinite, err := decodeLength(reader)
	if err != nil {
		return nil, err
	}
	if indefinite && !constructed {
		return nil, parseError("primitive node with indefinite length")
	}

	// Indefinite form
	var content []byte
	if !indefinite {
		content = make([]byte, length)
		_, err = io.ReadFull(reader, content)
		if err != nil {
			return nil, err
		}
	} else {
		buffer := bytes.NewBuffer([]byte{})
		childrenReader := io.TeeReader(reader, buffer)
		err := readEoc(childrenReader)
		if err != nil {
			return nil, err
		}
		// At this point, buffer also contains the EoC bytes
		content = buffer.Bytes()
		content = content[:len(content)-2]
	}

	raw := rawValue{class, tag, constructed, indefinite, content}
	return &raw, nil
}

func readEoc(reader io.Reader) error {

	for {
		class, tag, constructed, err := decodeIdentifier(reader)
		if err != nil {
			return err
		}

		length, indefinite, err := decodeLength(reader)
		if err != nil {
			return err
		}

		if indefinite && !constructed {
			return parseError("primitive node with indefinite length")
		}

		if class == 0 && tag == 0 && indefinite == false && length == 0 {
			break
		}

		if indefinite {
			err = readEoc(reader)
		} else {
			err = skipBytes(reader, int64(length))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeMultiByteTag(reader io.Reader) (uint, error) {
	// Tag is encoded in one or more following bytes
	tag := uint(0)
	for {
		// Read a byte
		b, err := readByte(reader)
		if err != nil {
			return 0, err
		}
		// if we need to shift out non zeros bits, so the tag is too big for an uint
		msb := uint64(0xfe) << (intBits - 8) // 7 most significant bits
		if uint64(tag)&msb != 0 {
			return 0, parseError("multi byte tag too big")
		}
		// Shift the previous value and add the new 7 bits
		tag = (tag << 7) | uint(b&0x7f)
		// The last byte is indicated by the most significant bit equals to zero
		if b&0x80 == 0 {
			break
		}
	}
	return tag, nil
}

func decodeIdentifier(reader io.Reader) (class uint, tag uint, constructed bool, err error) {

	b, err := readByte(reader)
	if err != nil {
		return
	}

	// Read class and constructed flag
	class = uint((b & 0xc0) >> 6)
	if b&0x20 != 0 {
		constructed = true
	}

	// Read the tag number
	tag = uint(b & 0x1f)
	if tag == 0x1f {
		// Tag is encoded in one or more following bytes
		tag, err = decodeMultiByteTag(reader)
		if err != nil {
			return
		}
	}
	return
}

func decodeLength(reader io.Reader) (length uint, indefinite bool, err error) {

	// Read length
	b, err := readByte(reader)
	if err != nil {
		return
	}

	// Short form
	if b&0x80 == 0 {
		length = uint(b)
		return
	}

	// Indefinite form
	if b == 0x80 {
		indefinite = true
		return
	}

	// Long form
	if b == 0xff {
		err = parseError("invalid number of length octets: %x", b)
		return
	}
	octets := make([]byte, int(b&0x7f))
	_, err = io.ReadFull(reader, octets)
	if err != nil {
		return
	}
	for _, b = range octets {
		msb := uint64(0xff) << (intBits - 8)
		if uint64(length)&msb != 0 {
			err = parseError("multi byte length too big")
			return
		}
		length = (length << 8) | uint(b)
	}

	return
}

func readByte(reader io.Reader) (byte, error) {
	buf := []byte{0x00}
	_, err := io.ReadFull(reader, buf)
	return buf[0], err
}

func skipBytes(reader io.Reader, count int64) error {
	_, err := io.CopyN(ioutil.Discard, reader, count)
	return err
}
