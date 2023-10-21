package user

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InvitationState string

const (
	InvitationStateFailed  InvitationState = "Failed"
	InvitationStateNotSent InvitationState = "NotSent"
	InvitationStateSending InvitationState = "Sending"
	InvitationStateSent    InvitationState = "Sent"
)

func PossibleValuesForInvitationState() []string {
	return []string{
		string(InvitationStateFailed),
		string(InvitationStateNotSent),
		string(InvitationStateSending),
		string(InvitationStateSent),
	}
}

func parseInvitationState(input string) (*InvitationState, error) {
	vals := map[string]InvitationState{
		"failed":  InvitationStateFailed,
		"notsent": InvitationStateNotSent,
		"sending": InvitationStateSending,
		"sent":    InvitationStateSent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InvitationState(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateLocked    ProvisioningState = "Locked"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateLocked),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"locked":    ProvisioningStateLocked,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type RegistrationState string

const (
	RegistrationStateNotRegistered RegistrationState = "NotRegistered"
	RegistrationStateRegistered    RegistrationState = "Registered"
)

func PossibleValuesForRegistrationState() []string {
	return []string{
		string(RegistrationStateNotRegistered),
		string(RegistrationStateRegistered),
	}
}

func parseRegistrationState(input string) (*RegistrationState, error) {
	vals := map[string]RegistrationState{
		"notregistered": RegistrationStateNotRegistered,
		"registered":    RegistrationStateRegistered,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegistrationState(input)
	return &out, nil
}
