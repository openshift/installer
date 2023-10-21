package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyPlayReadyContentKeyLocation = ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier{}

type ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier struct {
	KeyId string `json:"keyId"`

	// Fields inherited from ContentKeyPolicyPlayReadyContentKeyLocation
}

var _ json.Marshaler = ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier{}

func (s ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier: %+v", err)
	}

	return encoded, nil
}
