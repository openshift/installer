package systemd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Entry holds a systemd log entry.
// https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html
type Entry struct {
	// Hostname is the name of the originating host.
	Hostname string

	// InvocationID is a randomized, unique 128-bit ID identifying each
	// runtime cycle of the unit.
	InvocationID string

	// MessageID is a 128bit message ID for recognizing certain message
	// types.
	MessageID string

	// Message it the human-readable message string for this entry. This
	// is supposed to be the primary text shown to the user. It is
	// usually not translated (but might be in some cases), and is
	// not supposed to be parsed for metadata.
	Message string

	// PID is the process ID of the process the journal entry originates
	// from.
	PID int

	// RealtimeTimestamp is the wallclock time (CLOCK_REALTIME) at the
	// point in time the entry was received by the journal.  This has
	// different properties from SourceRealtimeTimestamp, as it is usually
	// a bit later but more likely to be monotonic.
	RealtimeTimestamp time.Time

	// Restarts is a counter of unit restarts, as long as the unit isn't
	// fully stopped, i.e. as long as it remains up or remains in
	// auto-start states.  This is only set on entries that schedule
	// restarts.
	Restarts int

	// SourceRealtimeTimestamp is the earliest trusted timestamp of the
	// message, if any is known that is different from the reception time
	// of the journal.
	SourceRealtimeTimestamp time.Time

	// SyslogIdentifier is the identifier string (i.e.  "tag") as
	// specified in the original datagram.
	SyslogIdentifier string

	// SystemdInvocationID is the invocation ID for the runtime cycle of
	// the unit the message was generated in, as available to processes of
	// the unit in InvocationID.
	SystemdInvocationID string

	// SystemdUnit is the systemd unit name of the process the journal
	// entry originates from.
	SystemdUnit string

	// Unit is the systemd unit name of the unit being operated on.
	// This is only set on entries where systemd logging its management
	// of other units.
	Unit string
}

// UnmarshalJSON unmarshals entry JSON.
func (e *Entry) UnmarshalJSON(b []byte) error {
	var data map[string]string
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	e.Hostname = data["_HOSTNAME"]
	e.InvocationID = data["INVOCATION_ID"]
	e.MessageID = data["MESSAGE_ID"]
	e.Message = data["MESSAGE"]
	e.SyslogIdentifier = data["SYSLOG_IDENTIFIER"]
	e.SystemdInvocationID = data["_SYSTEMD_INVOCATION_ID"]
	e.SystemdUnit = data["_SYSTEMD_UNIT"]
	e.Unit = data["UNIT"]

	var err error
	e.PID, err = strconv.Atoi(data["_PID"])
	if err != nil {
		return fmt.Errorf("parse _PID: %w", err)
	}

	if data["N_RESTARTS"] != "" {
		e.Restarts, err = strconv.Atoi(data["N_RESTARTS"])
		if err != nil {
			return fmt.Errorf("parse N_RESTARTS: %w", err)
		}
	}

	e.RealtimeTimestamp, err = parseTime(data["__REALTIME_TIMESTAMP"])
	if err != nil {
		return fmt.Errorf("parse __REALTIME_TIMESTAMP: %w", err)
	}

	if data["_SOURCE_REALTIME_TIMESTAMP"] != "" {
		e.SourceRealtimeTimestamp, err = parseTime(data["_SOURCE_REALTIME_TIMESTAMP"])
		if err != nil {
			return fmt.Errorf("parse _SOURCE_REALTIME_TIMESTAMP: %w", err)
		}
	}

	return nil
}

// parseTime parses a time in microseconds since the epoch UTC,
// formatted as a decimal string.
func parseTime(data string) (time.Time, error) {
	milliseconds, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(milliseconds/1000000, (milliseconds%1000000)*1000).UTC(), nil
}

// String formats the entry similarly to a syslog log line.
func (e *Entry) String() string {
	timestamp := e.SourceRealtimeTimestamp
	if timestamp.IsZero() {
		timestamp = e.RealtimeTimestamp
	}
	return fmt.Sprintf("%s %s %s[%d]: %s", timestamp.Format(time.RFC3339), e.Hostname, e.SyslogIdentifier, e.PID, e.Message)
}
