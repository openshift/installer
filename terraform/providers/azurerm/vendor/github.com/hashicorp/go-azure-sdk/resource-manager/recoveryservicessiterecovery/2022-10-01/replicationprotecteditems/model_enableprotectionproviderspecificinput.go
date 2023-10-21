package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnableProtectionProviderSpecificInput interface {
}

func unmarshalEnableProtectionProviderSpecificInputImplementation(input []byte) (EnableProtectionProviderSpecificInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling EnableProtectionProviderSpecificInput into map[string]interface: %+v", err)
	}

	value, ok := temp["instanceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "A2ACrossClusterMigration") {
		var out A2ACrossClusterMigrationEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2ACrossClusterMigrationEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "A2A") {
		var out A2AEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2AEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzureEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzureEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2EnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2EnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMageEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmEnableProtectionInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmEnableProtectionInput: %+v", err)
		}
		return out, nil
	}

	type RawEnableProtectionProviderSpecificInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawEnableProtectionProviderSpecificInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
