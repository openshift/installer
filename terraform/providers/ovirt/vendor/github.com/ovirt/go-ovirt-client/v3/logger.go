package ovirtclient

import (
	ovirtclientlog "github.com/ovirt/go-ovirt-client-log/v3"
)

// Logger is a thin wrapper around ovirtclientlog.Logger for convenience.
type Logger interface {
	ovirtclientlog.Logger
}
