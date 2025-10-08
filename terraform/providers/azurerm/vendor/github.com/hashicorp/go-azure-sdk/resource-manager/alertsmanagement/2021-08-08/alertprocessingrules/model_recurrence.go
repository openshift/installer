package alertprocessingrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Recurrence interface {
}

func unmarshalRecurrenceImplementation(input []byte) (Recurrence, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Recurrence into map[string]interface: %+v", err)
	}

	value, ok := temp["recurrenceType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Daily") {
		var out DailyRecurrence
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DailyRecurrence: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Monthly") {
		var out MonthlyRecurrence
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MonthlyRecurrence: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Weekly") {
		var out WeeklyRecurrence
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WeeklyRecurrence: %+v", err)
		}
		return out, nil
	}

	type RawRecurrenceImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawRecurrenceImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
