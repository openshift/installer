package asn1

import (
	"strconv"
	"strings"
)

type fieldOptions struct {
	universal    bool
	application  bool
	explicit     bool
	indefinite   bool
	optional     bool
	set          bool
	tag          *int
	defaultValue *int
	choice       *string
}

// validate returns an error if any option is invalid.
func (opts *fieldOptions) validate() error {
	tagError := func(class string) error {
		return syntaxError(
			"'tag' must be specified when '%s' is used", class)
	}
	if opts.universal && opts.tag == nil {
		return tagError("universal")
	}
	if opts.application && opts.tag == nil {
		return tagError("application")
	}
	if opts.tag != nil && *opts.tag < 0 {
		return syntaxError("'tag' cannot be negative: %d", *opts.tag)
	}
	if opts.choice != nil && *opts.choice == "" {
		return syntaxError("'choice' cannot be empty")
	}
	return nil
}

// parseOption returns a parsed fieldOptions or an error.
func parseOptions(s string) (*fieldOptions, error) {
	var opts fieldOptions
	for _, token := range strings.Split(s, ",") {
		args := strings.Split(strings.TrimSpace(token), ":")
		err := parseOption(&opts, args)
		if err != nil {
			return nil, err
		}
	}
	if err := opts.validate(); err != nil {
		return nil, err
	}
	return &opts, nil
}

// parseOption parse a single option.
func parseOption(opts *fieldOptions, args []string) error {
	var err error
	switch args[0] {
	case "":
		// ignore

	case "universal":
		opts.universal, err = parseBoolOption(args)

	case "application":
		opts.application, err = parseBoolOption(args)

	case "explicit":
		opts.explicit, err = parseBoolOption(args)

	case "indefinite":
		opts.indefinite, err = parseBoolOption(args)

	case "optional":
		opts.optional, err = parseBoolOption(args)

	case "set":
		opts.set, err = parseBoolOption(args)

	case "tag":
		opts.tag, err = parseIntOption(args)

	case "default":
		opts.defaultValue, err = parseIntOption(args)

	case "choice":
		opts.choice, err = parseStringOption(args)

	default:
		err = syntaxError("Invalid option: %s", args[0])
	}
	return err
}

// parseBoolOption just checks if no arguments were given.
func parseBoolOption(args []string) (bool, error) {
	if len(args) > 1 {
		return false, syntaxError("option '%s' does not have arguments.",
			args[0])
	}
	return true, nil
}

// parseIntOption parses an integer argument.
func parseIntOption(args []string) (*int, error) {
	if len(args) != 2 {
		return nil, syntaxError("option '%s' does not have arguments.")
	}
	num, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, syntaxError("invalid value '%s' for option '%s'.",
			args[1], args[0])
	}
	return &num, nil
}

// parseStringOption parses a string argument.
func parseStringOption(args []string) (*string, error) {
	if len(args) != 2 {
		return nil, syntaxError("option '%s' does not have arguments.")
	}
	return &args[1], nil
}
