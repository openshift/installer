// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// OperatorMode determines whether we'll run watchers and/or webhooks.
type OperatorMode int

const (
	OperatorModeWatchers = OperatorMode(1 << iota)
	OperatorModeWebhooks

	OperatorModeBoth = OperatorModeWatchers | OperatorModeWebhooks
)

// IncludesWebhooks returns whether an operator running in this mode
// should register webhooks.
func (m OperatorMode) IncludesWebhooks() bool {
	return m&OperatorModeWebhooks > 0
}

// IncludesWatchers returns whether an operator running in this mode
// should register reconcilers.
func (m OperatorMode) IncludesWatchers() bool {
	return m&OperatorModeWatchers > 0
}

// String converts the mode into a readable value.
func (m OperatorMode) String() string {
	switch m {
	case OperatorModeWatchers:
		return "watchers"
	case OperatorModeWebhooks:
		return "webhooks"
	case OperatorModeBoth:
		return "watchers-and-webhooks"
	default:
		panic(fmt.Sprintf("invalid operator mode value %d", m))
	}
}

// ParseOperatorMode converts a string value into the corresponding
// operator mode.
func ParseOperatorMode(value string) (OperatorMode, error) {
	switch strings.ToLower(value) {
	case "watchers":
		return OperatorModeWatchers, nil
	case "webhooks":
		return OperatorModeWebhooks, nil
	case "watchers-and-webhooks":
		return OperatorModeBoth, nil
	default:
		return OperatorMode(0), errors.Errorf(`operator mode value must be one of "watchers-and-webhooks", "webhooks" or "watchers" but was %q`, value)
	}
}
