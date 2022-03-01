package cfschema

import (
	"encoding/json"
	"fmt"
)

type Type string

// String is a string representation of Type.
func (t *Type) String() string {
	if t == nil {
		return ""
	}

	return string(*t)
}

// UnmarshalJSON is a custom JSON handler for Type.
func (t *Type) UnmarshalJSON(b []byte) error {
	var tmp string

	err := json.Unmarshal(b, &tmp)

	if err != nil {
		var tmpTypes []string

		err := json.Unmarshal(b, &tmpTypes)

		if err != nil {
			return err
		}

		if len(tmpTypes) != 2 {
			return fmt.Errorf("type arrays with less or more than 2 elements are not supported")
		}

		for _, tmpType := range tmpTypes {
			if tmpType == PropertyTypeObject {
				continue
			}

			*t = Type(tmpType)

			return nil
		}
	}

	*t = Type(tmp)

	return nil
}
