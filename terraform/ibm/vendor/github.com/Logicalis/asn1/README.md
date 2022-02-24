---

# WARNING

## This repo has been archived!

NO further developement will be made in the foreseen future.

---


[![Build Status](https://travis-ci.org/PromonLogicalis/asn1.svg?branch=master)](https://travis-ci.org/PromonLogicalis/asn1) [![Go Report Card](https://goreportcard.com/badge/github.com/PromonLogicalis/asn1)](https://goreportcard.com/report/github.com/PromonLogicalis/asn1) [![GoDoc](https://godoc.org/github.com/PromonLogicalis/asn1?status.svg)](https://godoc.org/github.com/PromonLogicalis/asn1)
# asn1
--
    import "github.com/PromonLogicalis/asn1"

Package asn1 implements encoding and decoding of ASN.1 data structures using
both Basic Encoding Rules (BER) or its subset, the Distinguished Encoding Rules
(BER).

This package is highly inspired by the Go standard package "encoding/asn1" while
supporting additional features such as BER encoding and decoding and ASN.1
CHOICE types.

By default and for convenience the package uses DER for encoding and BER for
decoding. However it's possible to use a Context object to set the desired
encoding and decoding rules as well other options.

Restrictions:

- BER allows STRING types, such as OCTET STRING and BIT STRING, to be encoded as
constructed types containing inner elements that should be concatenated to form
the complete string. The package does not support that, but in the future
decoding of constructed strings should be included.

## Usage

#### func  Decode

```go
func Decode(data []byte, obj interface{}) (rest []byte, err error)
```
Decode parses the given BER data into obj. The argument obj should be a
reference to the value that will hold the parsed data. Decode uses a default
Context and is equivalent to:

    rest, err := asn1.NewContext().Decode(data, &obj)

#### func  DecodeWithOptions

```go
func DecodeWithOptions(data []byte, obj interface{}, options string) (rest []byte, err error)
```
DecodeWithOptions parses the given BER data into obj using the additional
options. The argument obj should be a reference to the value that will hold the
parsed data. Decode uses a default Context and is equivalent to:

    rest, err := asn1.NewContext().DecodeWithOptions(data, &obj, options)

#### func  Encode

```go
func Encode(obj interface{}) (data []byte, err error)
```
Encode returns the DER encoding of obj. Encode uses a default Context and it's
equivalent to:

    data, err = asn1.NewContext().Encode(obj)

#### func  EncodeWithOptions

```go
func EncodeWithOptions(obj interface{}, options string) (data []byte, err error)
```
EncodeWithOptions returns the DER encoding of obj using additional options.
EncodeWithOptions uses a default Context and it's equivalent to:

    data, err = asn1.NewContext().EncodeWithOptions(obj, options)

#### type Choice

```go
type Choice struct {
	Type    reflect.Type
	Options string
}
```

Choice represents one option available for a CHOICE element.

#### type Context

```go
type Context struct {
}
```

Context keeps options that affect the ASN.1 encoding and decoding

Use the NewContext() function to create a new Context instance:

    ctx := ber.NewContext()
    // Set options, ex:
    ctx.SetDer(true, true)
    // And call decode or encode functions
    bytes, err := ctx.EncodeWithOptions(value, "explicit,application,tag:5")
    ...

#### func  NewContext

```go
func NewContext() *Context
```
NewContext creates and initializes a new context. The returned Context does not
contains any registered choice and it's set to DER encoding and BER decoding.

#### func (*Context) AddChoice

```go
func (ctx *Context) AddChoice(choice string, entries []Choice) error
```
AddChoice registers a list of types as options to a given choice.

The string choice refers to a choice name defined into an element via additional
options for DecodeWithOptions and EncodeWithOptions of via struct tags.

For example, considering that a field "Value" can be an INTEGER or an OCTET
STRING indicating two types of errors, each error with a different tag number,
the following can be used:

    // Error types
    type SimpleError string
    type ComplextError string
    // The main object
    type SomeSequence struct {
    	// ...
    	Value	interface{}	`asn1:"choice:value"`
    	// ...
    }
    // A Context with the registered choices
    ctx := asn1.NewContext()
    ctx.AddChoice("value", []asn1.Choice {
    	{
    		Type: reflect.TypeOf(int(0)),
    	},
    	{
    		Type: reflect.TypeOf(SimpleError("")),
    		Options: "tag:1",
    	},
    	{
    		Type: reflect.TypeOf(ComplextError("")),
    		Options: "tag:2",
    	},
    })

Some important notes:

1. Any choice value must be an interface. During decoding the necessary type
will be allocated to keep the parsed value.

2. The INTEGER type will be encoded using its default class and tag number and
so it's not necessary to specify any Options for it.

3. The two error types in our example are encoded as strings, in order to make
possible to differentiate both types during encoding they actually need to be
different types. This is solved by defining two alias types: SimpleError and
ComplextError.

4. Since both errors use the same encoding type, ASN.1 says they must have
distinguished tags. For that, the appropriate tag is defined for each type.

To encode a choice value, all that is necessary is to set the choice field with
the proper object. To decode a choice value, a type switch can be used to
determine which type was used.

#### func (*Context) Decode

```go
func (ctx *Context) Decode(data []byte, obj interface{}) (rest []byte, err error)
```
Decode parses the given data into obj. The argument obj should be a reference to
the value that will hold the parsed data.

See (*Context).DecodeWithOptions() for further details.

#### func (*Context) DecodeWithOptions

```go
func (ctx *Context) DecodeWithOptions(data []byte, obj interface{}, options string) (rest []byte, err error)
```
DecodeWithOptions parses the given data into obj using the additional options.
The argument obj should be a reference to the value that will hold the parsed
data.

It uses the reflect package to inspect obj and because of that only exported
struct fields (those that start with a capital letter) are considered.

The Context object defines the decoding rules (BER or DER) and the types
available for CHOICE types.

The asn1 package maps Go types to ASN.1 data structures. The package also
provides types to specific ASN.1 data structures. The default mapping is shown
in the table below:

    Go type                | ASN.1 universal tag
    -----------------------|--------------------
    bool                   | BOOLEAN
    All int and uint types | INTEGER
    *big.Int               | INTEGER
    string                 | OCTET STRING
    []byte                 | OCTET STRING
    asn1.Oid               | OBJECT INDETIFIER
    asn1.Null              | NULL
    Any array or slice     | SEQUENCE OF
    Any struct             | SEQUENCE

Arrays and slices are decoded using different rules. A slice is always appended
while an array requires an exact number of elements, otherwise a ParseError is
returned.

The default mapping can be changed using options provided in the argument
options (for the root value) or via struct tags for struct fields. Struct tags
use the namei space "asn1".

The available options for encoding and decoding are:

    tag

This option requires an numeric argument (ie: "tag:1") and indicates that a
element is encoded and decoded as a context specific element with the given tag
number. The context specific class can be overridden with the options
"application" or "universal".

    universal

Sets the tag class to universal. Requires "tag".

    application

Sets the tag class to application. Requires "tag".

    explicit

Indicates the element is encoded with an enclosing tag. It's usually used in
conjunction with "tag".

    optional

Indicates that an element can be suppressed.

A missing element that is not marked with "optional" or "default" causes a
ParseError to be returned during decoding. A missing element marked as optional
is left untouched.

During encoding, a zero value elements is suppressed from output if it's marked
as optional.

    default

This option is handled similarly to the "optional" option but requires a numeric
argument (ie: "default:1").

A missing element that is marked with "default" is set to the given default
value during decoding.

A zero value marked with "default" is suppressed from output when encoding is
set to DER or is encoded with the given default value when encoding is set to
BER.

    indefinite

This option is used only during encoding and causes a constructed element to be
encoded using the indefinite form.

    choice

Indicates that an element can be of one of several types as defined by
(*Context).AddChoice()

    set

Indicates that a struct, array or slice should be encoded and decoded as a SET
instead of a SEQUENCE.

It also affects the way that structs are encoded and decoded in DER. A struct
marked with "set" has its fields always encoded in the ascending order of its
tags, instead of following the order that the fields are defined in the struct.

Similarly, a struct marked with "set" always enforces that same order when
decoding in DER.

#### func (*Context) Encode

```go
func (ctx *Context) Encode(obj interface{}) (data []byte, err error)
```
Encode returns the ASN.1 encoding of obj.

See (*Context).EncodeWithOptions() for further details.

#### func (*Context) EncodeWithOptions

```go
func (ctx *Context) EncodeWithOptions(obj interface{}, options string) (data []byte, err error)
```
EncodeWithOptions returns the ASN.1 encoding of obj using additional options.

See (*Context).DecodeWithOptions() for further details regarding types and
options.

#### func (*Context) SetDer

```go
func (ctx *Context) SetDer(encoding bool, decoding bool)
```
SetDer sets DER mode for encofing and decoding.

#### func (*Context) SetLogger

```go
func (ctx *Context) SetLogger(logger *log.Logger)
```
SetLogger defines the logger used.

#### type Null

```go
type Null struct{}
```

Null is used to encode and decode ASN.1 NULLs.

#### type Oid

```go
type Oid []uint
```

Oid is used to encode and decode ASN.1 OBJECT IDENTIFIERs.

#### func (Oid) Cmp

```go
func (oid Oid) Cmp(other Oid) int
```
Cmp returns zero if both Oids are the same, a negative value if oid
lexicographically precedes other and a positive value otherwise.

#### func (Oid) String

```go
func (oid Oid) String() string
```
String returns the dotted representation of oid.

#### type ParseError

```go
type ParseError struct {
	Msg string
}
```

ParseError is returned by the package to indicate that the given data is
invalid.

#### func (*ParseError) Error

```go
func (e *ParseError) Error() string
```
Error returns the error message of a ParseError.

#### type SyntaxError

```go
type SyntaxError struct {
	Msg string
}
```

SyntaxError is returned by the package to indicate that the given value or
struct is invalid.

#### func (*SyntaxError) Error

```go
func (e *SyntaxError) Error() string
```
Error returns the error message of a SyntaxError.
