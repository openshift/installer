package local

import (
	"fmt"
	"strconv"
)

func validateMode(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)

	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) > 4 || len(v) < 3 {
		es = append(es, fmt.Errorf("bad mode for file - string length should be 3 or 4 digits: %s", v))
	}

	fileMode, err := strconv.ParseInt(v, 8, 64)

	if err != nil || fileMode > 0777 || fileMode < 0 {
		es = append(es, fmt.Errorf("bad mode for file - must be three octal digits: %s", v))
	}

	return
}
