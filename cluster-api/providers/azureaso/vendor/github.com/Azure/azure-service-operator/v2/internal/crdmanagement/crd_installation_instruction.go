/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package crdmanagement

import (
	"fmt"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type DiffResult string

const (
	NoDifference     = DiffResult("NoDifference")
	SpecDifferent    = DiffResult("SpecDifferent")
	VersionDifferent = DiffResult("VersionDifferent")
)

func (i DiffResult) DiffReason(crd apiextensions.CustomResourceDefinition) string {
	switch i {
	case NoDifference:
		return fmt.Sprintf("No difference between existing and goal CRD %q", makeMatchString(crd))
	case SpecDifferent:
		return fmt.Sprintf("The spec was different between existing and goal CRD %q", makeMatchString(crd))
	case VersionDifferent:
		return fmt.Sprintf("The version was different between existing and goal CRD %q", makeMatchString(crd))
	default:
		return fmt.Sprintf("Unknown DiffResult %q", i)
	}
}

type FilterResult string

const (
	MatchedExistingCRD = FilterResult("MatchedExistingCRD")
	MatchedPattern     = FilterResult("MatchedPattern")
	Excluded           = FilterResult("Excluded")
)

type CRDInstallationInstruction struct {
	CRD apiextensions.CustomResourceDefinition

	// FilterResult contains the result of if the CRD was considered for installation or not
	FilterResult FilterResult
	// FilterReason contains a user-facing reason for why the CRD was or was not considered for installation
	FilterReason string

	// DiffResult contains the result of the diff between the existing CRD (if any) and the goal CRD. This may
	// be NoDifference if the CRD was filtered from consideration before the diff phase.
	DiffResult DiffResult
}

func (i *CRDInstallationInstruction) ShouldApply() (bool, string) {
	excluded := i.FilterResult == Excluded
	if excluded {
		return false, i.FilterReason
	}

	isSame := i.DiffResult == NoDifference
	if isSame {
		return false, i.DiffResult.DiffReason(i.CRD)
	}

	return true, i.DiffResult.DiffReason(i.CRD)
}

func IncludedCRDs(instructions []*CRDInstallationInstruction) []*CRDInstallationInstruction {
	// Prealloc false positive: https://github.com/alexkohler/prealloc/issues/16
	//nolint:prealloc
	var result []*CRDInstallationInstruction
	for _, instruction := range instructions {
		if instruction.FilterResult == Excluded {
			continue
		}
		result = append(result, instruction)
	}

	return result
}
