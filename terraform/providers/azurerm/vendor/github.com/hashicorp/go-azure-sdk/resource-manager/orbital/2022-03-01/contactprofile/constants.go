package contactprofile

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoTrackingConfiguration string

const (
	AutoTrackingConfigurationDisabled AutoTrackingConfiguration = "disabled"
	AutoTrackingConfigurationSBand    AutoTrackingConfiguration = "sBand"
	AutoTrackingConfigurationXBand    AutoTrackingConfiguration = "xBand"
)

func PossibleValuesForAutoTrackingConfiguration() []string {
	return []string{
		string(AutoTrackingConfigurationDisabled),
		string(AutoTrackingConfigurationSBand),
		string(AutoTrackingConfigurationXBand),
	}
}

func parseAutoTrackingConfiguration(input string) (*AutoTrackingConfiguration, error) {
	vals := map[string]AutoTrackingConfiguration{
		"disabled": AutoTrackingConfigurationDisabled,
		"sband":    AutoTrackingConfigurationSBand,
		"xband":    AutoTrackingConfigurationXBand,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoTrackingConfiguration(input)
	return &out, nil
}

type Direction string

const (
	DirectionDownlink Direction = "downlink"
	DirectionUplink   Direction = "uplink"
)

func PossibleValuesForDirection() []string {
	return []string{
		string(DirectionDownlink),
		string(DirectionUplink),
	}
}

func parseDirection(input string) (*Direction, error) {
	vals := map[string]Direction{
		"downlink": DirectionDownlink,
		"uplink":   DirectionUplink,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Direction(input)
	return &out, nil
}

type Polarization string

const (
	PolarizationLHCP             Polarization = "LHCP"
	PolarizationLinearHorizontal Polarization = "linearHorizontal"
	PolarizationLinearVertical   Polarization = "linearVertical"
	PolarizationRHCP             Polarization = "RHCP"
)

func PossibleValuesForPolarization() []string {
	return []string{
		string(PolarizationLHCP),
		string(PolarizationLinearHorizontal),
		string(PolarizationLinearVertical),
		string(PolarizationRHCP),
	}
}

func parsePolarization(input string) (*Polarization, error) {
	vals := map[string]Polarization{
		"lhcp":             PolarizationLHCP,
		"linearhorizontal": PolarizationLinearHorizontal,
		"linearvertical":   PolarizationLinearVertical,
		"rhcp":             PolarizationRHCP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Polarization(input)
	return &out, nil
}

type Protocol string

const (
	ProtocolTCP Protocol = "TCP"
	ProtocolUDP Protocol = "UDP"
)

func PossibleValuesForProtocol() []string {
	return []string{
		string(ProtocolTCP),
		string(ProtocolUDP),
	}
}

func parseProtocol(input string) (*Protocol, error) {
	vals := map[string]Protocol{
		"tcp": ProtocolTCP,
		"udp": ProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Protocol(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
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
