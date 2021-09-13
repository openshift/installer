package asn1

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
)

// Context keeps options that affect the ASN.1 encoding and decoding
//
// Use the NewContext() function to create a new Context instance:
//
//	ctx := ber.NewContext()
//	// Set options, ex:
//	ctx.SetDer(true, true)
//	// And call decode or encode functions
//	bytes, err := ctx.EncodeWithOptions(value, "explicit,application,tag:5")
//	...
//
type Context struct {
	log     *log.Logger
	choices map[string][]choiceEntry
	der     struct {
		encoding bool
		decoding bool
	}
}

// Choice represents one option available for a CHOICE element.
type Choice struct {
	Type    reflect.Type
	Options string
}

// Internal register with information about the each CHOICE.
type choiceEntry struct {
	expectedElement
	typ  reflect.Type
	opts *fieldOptions
}

// NewContext creates and initializes a new context. The returned Context does
// not contains any registered choice and it's set to DER encoding and BER
// decoding.
func NewContext() *Context {
	ctx := &Context{}
	ctx.log = defaultLogger()
	ctx.choices = make(map[string][]choiceEntry)
	ctx.SetDer(true, false)
	return ctx
}

// getChoices returns a list of choices for a given name.
func (ctx *Context) getChoices(choice string) ([]choiceEntry, error) {
	entries := ctx.choices[choice]
	if entries == nil {
		return nil, syntaxError("invalid choice '%s'", choice)
	}
	return entries, nil
}

// getChoiceByType returns the choice associated to a given name and type.
func (ctx *Context) getChoiceByType(choice string, t reflect.Type) (entry choiceEntry, err error) {
	entries, err := ctx.getChoices(choice)
	if err != nil {
		return
	}

	for _, current := range entries {
		if current.typ == t {
			entry = current
			return
		}
	}
	err = syntaxError("invalid Go type '%s' for choice '%s'", t, choice)
	return
}

// getChoiceByTag returns the choice associated to a given tag.
func (ctx *Context) getChoiceByTag(choice string, class, tag uint) (entry choiceEntry, err error) {
	entries, err := ctx.getChoices(choice)
	if err != nil {
		return
	}

	for _, current := range entries {
		if current.class == class && current.tag == tag {
			entry = current
			return
		}
	}
	// TODO convert tag to text
	err = syntaxError("invalid tag [%d,%d] for choice '%s'", class, tag, choice)
	return
}

// addChoiceEntry adds a single choice to the list associated to a given name.
func (ctx *Context) addChoiceEntry(choice string, entry choiceEntry) error {
	for _, current := range ctx.choices[choice] {
		if current.class == entry.class && current.tag == entry.tag {
			return fmt.Errorf(
				"choice already registered: %s{%d, %d}",
				choice, entry.class, entry.tag)
		}
	}
	ctx.choices[choice] = append(ctx.choices[choice], entry)
	return nil
}

// AddChoice registers a list of types as options to a given choice.
//
// The string choice refers to a choice name defined into an element via
// additional options for DecodeWithOptions and EncodeWithOptions of via
// struct tags.
//
// For example, considering that a field "Value" can be an INTEGER or an OCTET
// STRING indicating two types of errors, each error with a different tag
// number, the following can be used:
//
//	// Error types
//	type SimpleError string
//	type ComplextError string
//	// The main object
//	type SomeSequence struct {
//		// ...
//		Value	interface{}	`asn1:"choice:value"`
//		// ...
//	}
//	// A Context with the registered choices
//	ctx := asn1.NewContext()
//	ctx.AddChoice("value", []asn1.Choice {
//		{
//			Type: reflect.TypeOf(int(0)),
//		},
//		{
//			Type: reflect.TypeOf(SimpleError("")),
//			Options: "tag:1",
//		},
//		{
//			Type: reflect.TypeOf(ComplextError("")),
//			Options: "tag:2",
//		},
//	})
//
// Some important notes:
//
// 1. Any choice value must be an interface. During decoding the necessary type
// will be allocated to keep the parsed value.
//
// 2. The INTEGER type will be encoded using its default class and tag number
// and so it's not necessary to specify any Options for it.
//
// 3. The two error types in our example are encoded as strings, in order to
// make possible to differentiate both types during encoding they actually need
// to be different types. This is solved by defining two alias types:
// SimpleError and ComplextError.
//
// 4. Since both errors use the same encoding type, ASN.1 says they must have
// distinguished tags. For that, the appropriate tag is defined for each type.
//
// To encode a choice value, all that is necessary is to set the choice field
// with the proper object. To decode a choice value, a type switch can be used
// to determine which type was used.
//
func (ctx *Context) AddChoice(choice string, entries []Choice) error {
	for _, e := range entries {
		opts, err := parseOptions(e.Options)
		if err != nil {
			return err
		}
		if opts.choice != nil {
			// TODO Add support for nested choices.
			return syntaxError(
				"nested choices are not allowed: '%s' inside '%s'",
				*opts.choice, choice)
		}
		raw := rawValue{}
		elem, err := ctx.getExpectedElement(&raw, e.Type, opts)
		if err != nil {
			return err
		}
		err = ctx.addChoiceEntry(choice, choiceEntry{
			expectedElement: elem,
			typ:             e.Type,
			opts:            opts,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// defaultLogger returns the default Logger. It's used to initialize a new context
// or when the logger is set to nil.
func defaultLogger() *log.Logger {
	return log.New(ioutil.Discard, "", 0)
}

// SetLogger defines the logger used.
func (ctx *Context) SetLogger(logger *log.Logger) {
	if logger == nil {
		logger = defaultLogger()
	}
	ctx.log = logger
}

// SetDer sets DER mode for encofing and decoding.
func (ctx *Context) SetDer(encoding bool, decoding bool) {
	ctx.der.encoding = encoding
	ctx.der.decoding = decoding
}
