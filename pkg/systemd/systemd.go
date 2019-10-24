// Package systemd contains handlers for systemd logs.
package systemd

import (
	"encoding/json"
	"io"
)

// Log holds systemd unit logs.
type Log struct {
	// Entries holds entries from the log.
	Entries []Entry
}

// NewLog parsers systemd's JSON log format.
// https://www.freedesktop.org/wiki/Software/systemd/json/
func NewLog(r io.Reader) (*Log, error) {
	decoder := json.NewDecoder(r)
	log := &Log{
		Entries: []Entry{},
	}
	for {
		var entry Entry
		if err := decoder.Decode(&entry); err == io.EOF {
			return log, nil
		} else if err != nil {
			return log, err
		}
		log.Entries = append(log.Entries, entry)
	}
}

// Format formats unit logs similarly to syslog log lines for the
// given number of runtime cycle invocations.  Zero will return all
// invocations, positive numbers will return the initial count
// invocations, and negative numbers will return the final count
// invocations.  For example, log.Render("my.unit", -2) will format
// the final two invocations of my.unit.
func (log *Log) Format(unit string, invocations int) []string {
	allInvocations := []string{}
	invocation := ""
	for _, entry := range log.Entries {
		if entry.Unit == unit || entry.SystemdUnit == unit {
			if entry.InvocationID != "" && entry.InvocationID != invocation {
				allInvocations = append(allInvocations, entry.InvocationID)
			} else if entry.SystemdInvocationID != "" && entry.SystemdInvocationID != invocation {
				allInvocations = append(allInvocations, entry.SystemdInvocationID)
			}
		}
	}

	var selectedInvocations []string
	switch {
	case invocations == 0:
		selectedInvocations = allInvocations
	case invocations > 0:
		selectedInvocations = allInvocations[:invocations]
	default: // invocations < 0
		selectedInvocations = allInvocations[len(allInvocations)+invocations:]
	}

	selected := make(map[string]struct{}, len(selectedInvocations))
	exists := struct{}{}
	for _, invocation := range selectedInvocations {
		selected[invocation] = exists
	}

	lines := []string{}
	for _, entry := range log.Entries {
		if entry.Unit == unit || entry.SystemdUnit == unit {
			match := false
			if _, ok := selected[entry.InvocationID]; ok {
				match = true
			} else if _, ok := selected[entry.SystemdInvocationID]; ok {
				match = true
			}
			if match {
				lines = append(lines, entry.String())
			}
		}
	}
	return lines
}

// Restarts returns the number of unit restarts, as long as the unit
// isn't fully stopped, i.e. as long as it remains up or remains in
// auto-start states.
func (log *Log) Restarts(unit string) int {
	for i := len(log.Entries) - 1; i >= 0; i-- {
		if log.Entries[i].Unit == unit && log.Entries[i].Restarts > 0 {
			return log.Entries[i].Restarts
		}
	}
	return 0
}
