package mccpv2

import (
	"errors"
	"fmt"
	"strings"
)

var (
	//ErrFilterNameMissing ...
	ErrFilterNameMissing = errors.New("Filter must have a name")

	//ErrFilterMissingOp ..
	ErrFilterMissingOp = errors.New("Filter must have an operator")
)

//Filter ...
type Filter struct {
	name  string
	op    string
	value string
}

//Name ...
func (f Filter) Name(name string) Filter {
	f.name = name
	return f
}

//Eq ...
func (f Filter) Eq(target string) Filter {
	f.op = ":"
	f.value = target
	return f
}

//In ...
func (f Filter) In(targets ...string) Filter {
	f.op = " IN "
	f.value = strings.Join(targets, ",")
	return f
}

//Ge ...
func (f Filter) Ge(target string) Filter {
	f.op = ":"
	f.value = target
	return f
}

//Le ...
func (f Filter) Le(target string) Filter {
	f.op = "<="
	f.value = target
	return f
}

//Gt ...
func (f Filter) Gt(target string) Filter {
	f.op = ">"
	f.value = target
	return f
}

//Lt ...
func (f Filter) Lt(target string) Filter {
	f.op = "<"
	f.value = target
	return f
}

func (f Filter) validate() error {
	if f.name == "" {
		return ErrFilterNameMissing
	}
	if f.op == "" {
		return ErrFilterMissingOp
	}
	return nil
}

//Build ...
func (f Filter) Build() (string, error) {
	err := f.validate()
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%s%s%s;", f.name, f.op, f.value), nil
}
