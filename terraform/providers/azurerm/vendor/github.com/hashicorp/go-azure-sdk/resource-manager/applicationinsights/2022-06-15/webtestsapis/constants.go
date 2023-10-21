package webtestsapis

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebTestKind string

const (
	WebTestKindMultistep WebTestKind = "multistep"
	WebTestKindPing      WebTestKind = "ping"
	WebTestKindStandard  WebTestKind = "standard"
)

func PossibleValuesForWebTestKind() []string {
	return []string{
		string(WebTestKindMultistep),
		string(WebTestKindPing),
		string(WebTestKindStandard),
	}
}

func parseWebTestKind(input string) (*WebTestKind, error) {
	vals := map[string]WebTestKind{
		"multistep": WebTestKindMultistep,
		"ping":      WebTestKindPing,
		"standard":  WebTestKindStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WebTestKind(input)
	return &out, nil
}
