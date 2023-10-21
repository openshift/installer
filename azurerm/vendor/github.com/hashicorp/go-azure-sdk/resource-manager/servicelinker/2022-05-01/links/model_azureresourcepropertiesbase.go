package links

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourcePropertiesBase interface {
}

func unmarshalAzureResourcePropertiesBaseImplementation(input []byte) (AzureResourcePropertiesBase, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureResourcePropertiesBase into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "KeyVault") {
		var out AzureKeyVaultProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureKeyVaultProperties: %+v", err)
		}
		return out, nil
	}

	type RawAzureResourcePropertiesBaseImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawAzureResourcePropertiesBaseImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
