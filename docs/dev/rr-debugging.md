# RR (Record and Replay) Debugging

This document describes the implementation of RR (Record and Replay) debugging capabilities in the OpenShift installer's Cluster API system. RR is a powerful debugging tool that allows you to record program execution and then replay it deterministically for debugging purposes.


## Overview

The implementation includes modifications to enable RR debugging for Cluster API controllers. This allows developers to:

1. Record the execution of Cluster API controllers with deterministic replay
2. Debug complex timing-dependent issues that are difficult to reproduce
3. Step through execution multiple times with identical behavior
4. Analyze race conditions and concurrency issues

## Implementation Details

### Key Changes

The implementation modifies two main components:

1. **Process Management** (`pkg/clusterapi/internal/process/process.go`)
2. **System Controller** (`pkg/clusterapi/system.go`)

### Process Management Modifications

The process management system has been modified to:

- Use process group signaling (`syscall.Kill(-ps.Cmd.Process.Pid, syscall.SIGTERM)`) instead of direct process signaling
- Add enhanced logging for process exit states
- Temporarily disable timeout-based process termination for debugging purposes

### Controller Execution Modifications

The controller execution system has been modified to:

- Replace capi provider with `rr` for CPU binding and execution recording
- Add RR-specific flags for optimal debugging:
  - `--wait`: Ensures RR waits for the recorded process
  - `--disable-avx-512`: Disables AVX-512 instructions for compatibility
  - `--bind-to-cpu=0`: Binds execution to CPU 0 for deterministic behavior (or CPUs with P and E cores)

## Setup

### Installing RR

```bash
# On Fedora/RHEL/CentOS
sudo dnf install rr

# On Ubuntu/Debian
sudo apt install rr
```

### Installing Delve

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```
### Applying the Patch

Apply the RR debugging patch to enable RR recording:

```bash
# Apply the patch (patch file location: docs/dev/rr-debugging.patch)
git apply docs/dev/rr-debugging.patch
```

### Building the Installer

Build the installer with RR debugging enabled:

```bash
MODE=dev TAGS="release" ./hack/build.sh
```

## Usage

### Recording Execution

`rr` requires these kernel parameters to trace:

```bash
sudo sysctl kernel.perf_event_paranoid=-1;
sudo sysctl kernel.kptr_restrict=0
```

When running the installer with RR debugging enabled, Cluster API controllers will automatically be recorded using RR. The recording process:

1. Captures all system calls, memory accesses, and timing information
2. Stores the trace in `~/.local/share/rr/latest-trace`
3. Maintains deterministic replay capability

### Replaying with Delve

To replay a recorded trace using Delve (dlv), use the following command:

```bash
dlv replay --listen=:2345 --headless=true --api-version=2 --accept-multiclient ~/.local/share/rr/latest-trace
```

**Command Options:**
- `--listen=:2345`: Listens on port 2345 for debugger connections
- `--headless=true`: Runs in headless mode without requiring a terminal
- `--api-version=2`: Uses Delve API version 2
- `--accept-multiclient`: Allows multiple debugger clients to connect
- `~/.local/share/rr/latest-trace`: Path to the recorded RR trace

### Connecting a Debugger

After starting the replay session, you can connect your preferred debugger:

- **VS Code**: Configure launch.json to connect to `localhost:2345`
- **GoLand/IntelliJ**: [Use the remote debugging configuration](https://www.jetbrains.com/help/go/attach-to-running-go-processes-with-debugger.html#attach-to-a-process-on-a-remote-machine) 
- **Command line**: Use `dlv connect localhost:2345`

## References

- [RR Documentation](https://rr-project.org/)
- [Delve Documentation](https://github.com/go-delve/delve)
- [Cluster API Documentation](https://cluster-api.sigs.k8s.io/)
