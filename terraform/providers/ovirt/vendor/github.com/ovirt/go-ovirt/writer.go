//
// Copyright (c) 2017 Joey <majunjiev@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ovirtsdk

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"time"
	"unicode/utf8"
)

var (
	esc_quot = []byte("&#34;") // shorter than "&quot;"
	esc_apos = []byte("&#39;") // shorter than "&apos;"
	esc_amp  = []byte("&amp;")
	esc_lt   = []byte("&lt;")
	esc_gt   = []byte("&gt;")
	esc_tab  = []byte("&#x9;")
	esc_nl   = []byte("&#xA;")
	esc_cr   = []byte("&#xD;")
	esc_fffd = []byte("\uFFFD") // Unicode replacement character
)

// XMLWriter marshalizes the struct to XML
type XMLWriter struct {
	*bufio.Writer
}

// NewXMLWriter creates a XMLWriter instance
func NewXMLWriter(w io.Writer) *XMLWriter {
	return &XMLWriter{
		Writer: bufio.NewWriter(w),
	}
}

func (writer *XMLWriter) WriteElement(uri, name, value string, attrs map[string]string) error {
	if name == "" {
		return fmt.Errorf("xml: start tag with no name")
	}
	writer.WriteStart(uri, name, attrs)
	writer.WriteString(value)
	writer.WriteEnd(name)
	return nil
}

func (writer *XMLWriter) WriteStart(uri, name string, attrs map[string]string) error {
	if name == "" {
		return fmt.Errorf("xml: start tag with no name")
	}
	writer.WriteByte('<')
	writer.WriteString(name)
	if uri != "" {
		writer.WriteString(` xmlns="`)
		writer.EscapeString(uri)
		writer.WriteByte('"')
	}
	if attrs != nil && len(attrs) > 0 {
		for attrName, attrValue := range attrs {
			writer.WriteByte(' ')
			writer.WriteString(attrName)
			writer.WriteString(`="`)
			writer.EscapeString(attrValue)
			writer.WriteByte('"')
		}
	}
	writer.WriteByte('>')
	return nil
}

func (writer *XMLWriter) WriteEnd(name string) error {
	if name == "" {
		return fmt.Errorf("xml: end tag with no name")
	}
	writer.WriteByte('<')
	writer.WriteByte('/')
	writer.WriteString(name)
	writer.WriteByte('>')
	return nil
}

func (writer *XMLWriter) WriteCharacter(name, s string) error {
	return writer.WriteElement("", name, s, nil)
}

func (writer *XMLWriter) WriteCharacters(name string, ss []string) error {
	for _, s := range ss {
		err := writer.WriteCharacter(name, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *XMLWriter) WriteBool(name string, b bool) error {
	return writer.WriteElement("", name, writer.FormatBool(b), nil)
}

func (writer *XMLWriter) WriteBools(name string, bs []bool) error {
	for _, b := range bs {
		err := writer.WriteBool(name, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *XMLWriter) FormatBool(b bool) string {
	return strconv.FormatBool(b)
}

func (writer *XMLWriter) WriteInt64(name string, i int64) error {
	return writer.WriteElement("", name, writer.FormatInt64(i), nil)
}

func (writer *XMLWriter) WriteInt64s(name string, is []int64) error {
	for _, i := range is {
		err := writer.WriteInt64(name, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *XMLWriter) FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func (writer *XMLWriter) WriteFloat64(name string, f float64) error {
	return writer.WriteElement("", name, writer.FormatFloat64(f), nil)
}

func (writer *XMLWriter) WriteFloat64s(name string, fs []float64) error {
	for _, f := range fs {
		err := writer.WriteFloat64(name, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *XMLWriter) FormatFloat64(f float64) string {
	return strconv.FormatFloat(f, 'e', 3, 64)
}

func (writer *XMLWriter) WriteDate(name string, t time.Time) error {
	return writer.WriteElement("", name, writer.FormatDate(t), nil)
}

func (writer *XMLWriter) WriteDates(name string, ts []time.Time) error {
	for _, t := range ts {
		err := writer.WriteDate(name, t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *XMLWriter) FormatDate(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func (writer *XMLWriter) EscapeString(s string) {
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width
		switch r {
		case '"':
			esc = esc_quot
		case '\'':
			esc = esc_apos
		case '&':
			esc = esc_amp
		case '<':
			esc = esc_lt
		case '>':
			esc = esc_gt
		case '\t':
			esc = esc_tab
		case '\n':
			esc = esc_nl
		case '\r':
			esc = esc_cr
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = esc_fffd
				break
			}
			continue
		}
		writer.WriteString(s[last : i-width])
		writer.Write(esc)
		last = i
	}
	writer.WriteString(s[last:])
}

// Decide whether the given rune is in the XML Character Range, per
// the Char production of http://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xDF77 ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}
