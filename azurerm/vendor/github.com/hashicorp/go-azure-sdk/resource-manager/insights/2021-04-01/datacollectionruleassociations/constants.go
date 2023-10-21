package datacollectionruleassociations

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KnownDataCollectionRuleAssociationProvisioningState string

const (
	KnownDataCollectionRuleAssociationProvisioningStateCreating  KnownDataCollectionRuleAssociationProvisioningState = "Creating"
	KnownDataCollectionRuleAssociationProvisioningStateDeleting  KnownDataCollectionRuleAssociationProvisioningState = "Deleting"
	KnownDataCollectionRuleAssociationProvisioningStateFailed    KnownDataCollectionRuleAssociationProvisioningState = "Failed"
	KnownDataCollectionRuleAssociationProvisioningStateSucceeded KnownDataCollectionRuleAssociationProvisioningState = "Succeeded"
	KnownDataCollectionRuleAssociationProvisioningStateUpdating  KnownDataCollectionRuleAssociationProvisioningState = "Updating"
)

func PossibleValuesForKnownDataCollectionRuleAssociationProvisioningState() []string {
	return []string{
		string(KnownDataCollectionRuleAssociationProvisioningStateCreating),
		string(KnownDataCollectionRuleAssociationProvisioningStateDeleting),
		string(KnownDataCollectionRuleAssociationProvisioningStateFailed),
		string(KnownDataCollectionRuleAssociationProvisioningStateSucceeded),
		string(KnownDataCollectionRuleAssociationProvisioningStateUpdating),
	}
}

func parseKnownDataCollectionRuleAssociationProvisioningState(input string) (*KnownDataCollectionRuleAssociationProvisioningState, error) {
	vals := map[string]KnownDataCollectionRuleAssociationProvisioningState{
		"creating":  KnownDataCollectionRuleAssociationProvisioningStateCreating,
		"deleting":  KnownDataCollectionRuleAssociationProvisioningStateDeleting,
		"failed":    KnownDataCollectionRuleAssociationProvisioningStateFailed,
		"succeeded": KnownDataCollectionRuleAssociationProvisioningStateSucceeded,
		"updating":  KnownDataCollectionRuleAssociationProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownDataCollectionRuleAssociationProvisioningState(input)
	return &out, nil
}
