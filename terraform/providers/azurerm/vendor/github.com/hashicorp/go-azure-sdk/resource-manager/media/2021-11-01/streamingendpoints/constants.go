package streamingendpoints

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingEndpointResourceState string

const (
	StreamingEndpointResourceStateDeleting StreamingEndpointResourceState = "Deleting"
	StreamingEndpointResourceStateRunning  StreamingEndpointResourceState = "Running"
	StreamingEndpointResourceStateScaling  StreamingEndpointResourceState = "Scaling"
	StreamingEndpointResourceStateStarting StreamingEndpointResourceState = "Starting"
	StreamingEndpointResourceStateStopped  StreamingEndpointResourceState = "Stopped"
	StreamingEndpointResourceStateStopping StreamingEndpointResourceState = "Stopping"
)

func PossibleValuesForStreamingEndpointResourceState() []string {
	return []string{
		string(StreamingEndpointResourceStateDeleting),
		string(StreamingEndpointResourceStateRunning),
		string(StreamingEndpointResourceStateScaling),
		string(StreamingEndpointResourceStateStarting),
		string(StreamingEndpointResourceStateStopped),
		string(StreamingEndpointResourceStateStopping),
	}
}

func parseStreamingEndpointResourceState(input string) (*StreamingEndpointResourceState, error) {
	vals := map[string]StreamingEndpointResourceState{
		"deleting": StreamingEndpointResourceStateDeleting,
		"running":  StreamingEndpointResourceStateRunning,
		"scaling":  StreamingEndpointResourceStateScaling,
		"starting": StreamingEndpointResourceStateStarting,
		"stopped":  StreamingEndpointResourceStateStopped,
		"stopping": StreamingEndpointResourceStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StreamingEndpointResourceState(input)
	return &out, nil
}
