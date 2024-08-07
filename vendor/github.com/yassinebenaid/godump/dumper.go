package godump

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// Style defines a general interface used for styling.
type Style interface {
	Apply(string) string
}

//nolint:revive // this name is uncommon, but why not
func __(s Style, v string) string {
	if s == nil {
		return v
	}
	return s.Apply(v)
}

// RGB implements [Style] and allows you to define your style as an RGB value, it uses ANSI escape sequences under the hood.
type RGB struct {
	R, G, B int
}

func (rgb RGB) Apply(v string) string {
	return fmt.Sprintf("\033[38;2;%v;%v;%vm%s\033[0m", rgb.R, rgb.G, rgb.B, v)
}

// Theme allows you to define your preferred styling for [Dumper].
type Theme struct {
	// String defines the style used for strings
	String Style

	// Quotes defines the style used for quotes (") around strings.
	Quotes Style

	// Bool defines the style used for boolean values.
	Bool Style

	// Number defines the style used for numbers, including all types of integers, floats and complex numbers.
	Number Style

	// Types defines the style used for defined and/or structural types, eg., slices, structs, maps...
	Types Style

	// Nil defines the style used for nil.
	Nil Style

	// Func defines the style used for functions.
	Func Style

	// Chan defines the style used for channels.
	Chan Style

	// UnsafePointer defines the style used for unsafe pointers.
	UnsafePointer Style

	// Address defines the style used for address symbol '&'.
	Address Style

	// PointerTag defines the style used for pointer tags, typically the pointer id '#x' and the recursive reference '@x'.
	PointerTag Style

	// Fields defines the style used for struct fields.
	Fields Style

	// Braces defines the style used for braces '{}' in structural types.
	Braces Style
}

// DefaultTheme is the default [Theme] used by [Dump].
var DefaultTheme = Theme{
	String:        RGB{138, 201, 38},
	Quotes:        RGB{112, 214, 255},
	Bool:          RGB{249, 87, 56},
	Number:        RGB{10, 178, 242},
	Types:         RGB{0, 150, 199},
	Address:       RGB{205, 93, 0},
	PointerTag:    RGB{110, 110, 110},
	Nil:           RGB{219, 57, 26},
	Func:          RGB{160, 90, 220},
	Fields:        RGB{189, 176, 194},
	Chan:          RGB{195, 154, 76},
	UnsafePointer: RGB{89, 193, 180},
	Braces:        RGB{185, 86, 86},
}

// Dump pretty prints `v` using the default Dumper options and the default theme
func Dump(v any) error {
	return (&Dumper{
		Theme: DefaultTheme,
	}).Println(v)
}

// Dumper provides an elegant interface to pretty print any variable of any type in a colored and structured format.
//
// The zero value for Dumper is a theme-less Dumper ready to use.
type Dumper struct {
	// Indentation is an optional string used for indentation.
	// The default value is a string of three spaces.
	Indentation string

	// ShowPrimitiveNamedTypes determines whether to show primitive named types.
	ShowPrimitiveNamedTypes bool

	// HidePrivateFields allows you to optionally hide struct's unexported fields from being printed.
	HidePrivateFields bool

	// Theme allows you to define your preferred styling.
	Theme Theme

	buf    bytes.Buffer
	depth  uint
	ptrs   map[uintptr]uint
	ptrTag uint
}

// Print formats `v` and writes the result to standard output.
//
// It returns a write error if encountered while writing to standard output.
func (d *Dumper) Print(v any) error {
	return d.Fprint(os.Stdout, v)
}

// Println formats `v`, appends a new line, and writes the result to standard output.
//
// It returns a write error if encountered while writing to standard output.
func (d *Dumper) Println(v any) error {
	return d.Fprintln(os.Stdout, v)
}

// Fprint formats `v` and writes the result to `dst`.
//
// It returns a write error if encountered while writing to `dst`.
func (d *Dumper) Fprint(dst io.Writer, v any) error {
	d.init()
	d.dump(reflect.ValueOf(v))
	if _, err := d.buf.WriteTo(dst); err != nil {
		return fmt.Errorf("dumper error: encountered unexpected write error, %v", err)
	}
	return nil
}

// Fprintln formats `v`, appends a new line, and writes the result to `dst`.
//
// It returns a write error if encountered while writing to `dst`.
func (d *Dumper) Fprintln(dst io.Writer, v any) error {
	d.init()
	d.dump(reflect.ValueOf(v))
	d.buf.WriteString("\n")
	if _, err := d.buf.WriteTo(dst); err != nil {
		return fmt.Errorf("dumper error: encountered unexpected write error, %v", err)
	}
	return nil
}

// Sprint formats `v` and returns the resulting string.
func (d *Dumper) Sprint(v any) string {
	d.init()
	d.dump(reflect.ValueOf(v))
	return d.buf.String()
}

// Sprintln formats `v`, appends a new line, and returns the resulting string.
func (d *Dumper) Sprintln(v any) string {
	d.init()
	d.dump(reflect.ValueOf(v))
	d.buf.WriteString("\n")
	return d.buf.String()
}

func (d *Dumper) init() {
	d.buf.Reset()
	d.ptrs = make(map[uintptr]uint)
	if d.Indentation == "" {
		d.Indentation = "   "
	}
}

func (d *Dumper) dump(val reflect.Value, ignoreDepth ...bool) {
	if len(ignoreDepth) <= 0 || !ignoreDepth[0] {
		d.indent()
	}

	switch val.Kind() {
	case reflect.String:
		d.wrapType(val, __(d.Theme.Quotes, `"`)+__(d.Theme.String, val.String())+__(d.Theme.Quotes, `"`))
	case reflect.Bool:
		d.wrapType(val, __(d.Theme.Bool, fmt.Sprintf("%t", val.Bool())))
	case reflect.Slice, reflect.Array:
		d.dumpSlice(val)
	case reflect.Map:
		d.dumpMap(val)
	case reflect.Func:
		d.buf.WriteString(__(d.Theme.Func, val.Type().String()))
	case reflect.Chan:
		d.buf.WriteString(__(d.Theme.Chan, val.Type().String()))
		if vCap := val.Cap(); vCap > 0 {
			d.buf.WriteString(__(d.Theme.Chan, fmt.Sprintf("<%d>", vCap)))
		}
	case reflect.Struct:
		d.dumpStruct(val)
	case reflect.Pointer:
		d.dumpPointer(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.wrapType(val, __(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.wrapType(val, __(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Float32, reflect.Float64:
		d.wrapType(val, __(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Complex64, reflect.Complex128:
		d.wrapType(val, __(d.Theme.Number, fmt.Sprint(val)))
	case reflect.Uintptr:
		d.wrapType(val, __(d.Theme.Number, fmt.Sprintf("0x%x", val.Uint())))
	case reflect.Invalid:
		d.buf.WriteString(__(d.Theme.Nil, "nil"))
	case reflect.Interface:
		d.dump(val.Elem(), true)
	case reflect.UnsafePointer:
		d.buf.WriteString(
			__(d.Theme.Types, val.Type().String()) +
				__(d.Theme.Braces, "(") +
				__(d.Theme.UnsafePointer, fmt.Sprintf("0x%x", uintptr(val.UnsafePointer()))) +
				__(d.Theme.Braces, ")"),
		)
	}
}

func (d *Dumper) dumpSlice(v reflect.Value) {
	length := v.Len()

	var tag string
	if d.ptrTag != 0 {
		tag = __(d.Theme.PointerTag, fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.buf.WriteString(__(d.Theme.Types, fmt.Sprintf("%s:%d:%d", v.Type(), length, v.Cap())))
	d.buf.WriteString(__(d.Theme.Braces, fmt.Sprintf(" {%s", tag)))

	d.depth++
	for i := 0; i < length; i++ {
		d.buf.WriteString("\n")
		d.dump(v.Index(i))
		d.buf.WriteString(",")
	}
	d.depth--

	if length > 0 {
		d.buf.WriteString("\n")
		d.indent()
	}

	d.buf.WriteString(__(d.Theme.Braces, "}"))
}

func (d *Dumper) dumpMap(v reflect.Value) {
	keys := v.MapKeys()

	var tag string
	if d.ptrTag != 0 {
		tag = __(d.Theme.PointerTag, fmt.Sprintf("#%d", d.ptrTag))
		d.ptrTag = 0
	}

	d.buf.WriteString(__(d.Theme.Types, fmt.Sprintf("%s:%d", v.Type(), len(keys))))
	d.buf.WriteString(__(d.Theme.Braces, fmt.Sprintf(" {%s", tag)))

	d.depth++
	for _, key := range keys {
		d.buf.WriteString("\n")
		d.dump(key)
		d.buf.WriteString((": "))
		d.dump(v.MapIndex(key), true)
		d.buf.WriteString((","))
	}
	d.depth--

	if len(keys) > 0 {
		d.buf.WriteString("\n")
		d.indent()
	}

	d.buf.WriteString(__(d.Theme.Braces, "}"))
}

func (d *Dumper) dumpPointer(v reflect.Value) {
	elem := v.Elem()

	if isPrimitive(elem) {
		if elem.IsValid() {
			d.buf.WriteString(__(d.Theme.Address, "&"))
		}
		d.dump(elem, true)
		return
	}

	d.buf.WriteString(__(d.Theme.Address, "&"))
	addr := uintptr(v.UnsafePointer())

	if id, ok := d.ptrs[addr]; ok {
		d.buf.WriteString(__(d.Theme.PointerTag, fmt.Sprintf("@%d", id)))
		return
	}

	d.ptrs[addr] = uint(len(d.ptrs) + 1)
	d.ptrTag = uint(len(d.ptrs))
	d.dump(elem, true)
	d.ptrTag = 0
}

func (d *Dumper) dumpStruct(v reflect.Value) {
	vtype := v.Type()

	var tag string
	if d.ptrTag != 0 {
		tag = fmt.Sprintf("#%d", d.ptrTag)
		d.ptrTag = 0
	}

	if t := vtype.String(); strings.HasPrefix(t, "struct") {
		d.buf.WriteString(__(d.Theme.Types, "struct"))
	} else {
		d.buf.WriteString(__(d.Theme.Types, t))
	}
	d.buf.WriteString(__(d.Theme.Braces, " {"))
	d.buf.WriteString(__(d.Theme.PointerTag, tag))

	var hasFields bool

	d.depth++
	for i := 0; i < v.NumField(); i++ {
		key := vtype.Field(i)
		if !key.IsExported() && d.HidePrivateFields {
			continue
		}

		hasFields = true

		d.buf.WriteString("\n")
		d.indent()

		d.buf.WriteString(__(d.Theme.Fields, key.Name))
		d.buf.WriteString((": "))
		d.dump(v.Field(i), true)
		d.buf.WriteString((","))
	}
	d.depth--

	if hasFields {
		d.buf.WriteString("\n")
		d.indent()
	}

	d.buf.WriteString(__(d.Theme.Braces, "}"))
}

func (d *Dumper) indent() {
	d.buf.WriteString(strings.Repeat(d.Indentation, int(d.depth)))
}

func isPrimitive(val reflect.Value) bool {
	v := val
	for {
		switch v.Kind() {
		case reflect.String, reflect.Bool, reflect.Func, reflect.Chan,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Invalid, reflect.UnsafePointer:
			return true
		case reflect.Pointer:
			v = v.Elem()
		default:
			return false
		}
	}
}

func (d *Dumper) wrapType(v reflect.Value, str string) {
	if d.ShowPrimitiveNamedTypes {
		if t := v.Type(); t.PkgPath() != "" {
			str = __(d.Theme.Types, t.String()) + __(d.Theme.Braces, "(") + str + __(d.Theme.Braces, ")")
		}
	}

	d.buf.WriteString(str)
}
