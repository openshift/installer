# Unified logging interface for Go client libraries

This repository contains a simple unified logging interface for all oVirt Go client libraries. This is *not* a logger itself, just an interface definition coupled with default loggers.

This library is used as a dependency, so you will most likely not need to rely on it. However, if you need to fetch it you can do so using `go get`:

```bash
go get github.com/ovirt/go-ovirt-client-log/v2
```

You can then reference this library using the `ovirtclientlog` package name.

## Default loggers

This library providers 3 default loggers:

- Go logging
- Go test logging
- "NOOP" logging

### Go logging

A Go logger can be created using the `NewGoLogger()` function.

```go
logger := ovirtclientlog.NewGoLogger()
```

Optionally, a [`*log.Logger`](https://pkg.go.dev/log#Logger) instance can be passed. If it is not passed, the log is written to the globally configured log destination.

```go
buf := &bytes.Buffer{}
backingLogger := log.New(buf, "", 0)
logger := ovirtclientlog.NewGoLogger(backingLogger)
```

### Test logging

This library also contains the ability to log via [`Logf` in `*testing.T`](https://pkg.go.dev/testing#T.Logf). You can create a logger like this:

```go
func TestYourFeature(t *testing.T) {
	logger := ovirtclientlog.NewTestLogger(t)
}
```

Using the test logger will have the benefit that the log messages will be properly attributed to the test that wrote them even if multiple tests are executed in parallel.

### NOOP logging

If you need a logger that doesn't do anything simply use `ovirtclientlog.NewNOOPLogger()`.

### klog logging

We provide klog-based logging [in a separate library](https://github.com/oVirt/go-ovirt-client-log-klog). You can use it as follows:

```
go get github.com/oVirt/go-ovirt-client-log-klog
```

```go
package main

import (
	kloglogger "github.com/ovirt/go-ovirt-client-log-klog"
)

func main() {
	logger := kloglogger.New()
	
	// Pass logger to other library as needed.
}
```

You can also specify separate verbosity levels:

```go
package main

import (
	kloglogger "github.com/ovirt/go-ovirt-client-log-klog"
	"k8s.io/klog/v2"
)

func main() {
	logger := kloglogger.NewVerbose(
		klog.V(4),
		klog.V(3),
		klog.V(2),
		klog.V(1),
    )

	// Pass logger to other library as needed.
}
```

## Adding your own logger

You can easily integrate your own logger too. Loggers must satisfy the following interface:

```go
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})

    WithContext(ctx context.Context) Logger
}
```

For example, you can adapt logging to [klog](https://github.com/kubernetes/klog) like this:

```go
type klogLogger struct {
}

func (k klogLogger) Debugf(format string, args ...interface{}) {
	// klog doesn't have a debug level
	klog.Infof(format, args...)
}

func (k klogLogger) Infof(format string, args ...interface{}) {
	klog.Infof(format, args...)
}

func (k klogLogger) Warningf(format string, args ...interface{}) {
	klog.Warningf(format, args...)
}

func (k klogLogger) Errorf(format string, args ...interface{}) {
	klog.Errorf(format, args...)
}

func (k klogLogger) WithContext(_ context.Context) Logger {
    return k
}
```

You can then create a new logger copy like this:

```go
logger := &klogLogger{}
```
