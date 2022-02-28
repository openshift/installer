package openstack

import (
	"fmt"
	"strconv"
)

// blockStorageQuotasetVolTypeQuotaToInt converts block storage vol type quota from map of strings to map of integers.
func blockStorageQuotasetVolTypeQuotaToInt(raw map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{}, len(raw))

	for k, v := range raw {
		strVal, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%v is not a string", v)
		}

		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			return nil, fmt.Errorf("%s can't be converted to int", strVal)
		}

		res[k] = intVal
	}

	return res, nil
}

// blockStorageQuotasetVolTypeQuotaToStr converts block storage vol type quota from map of interfaces to map of strings.
func blockStorageQuotasetVolTypeQuotaToStr(raw map[string]interface{}) (map[string]string, error) {
	res := make(map[string]string, len(raw))

	for k, v := range raw {
		switch value := v.(type) {
		case int:
			res[k] = strconv.Itoa(value)
		case float32, float64:
			res[k] = fmt.Sprintf("%.0f", value)
		case string:
			res[k] = value
		default:
			return nil, fmt.Errorf("got unknown type for quota volume type %s value: %+v", k, v)
		}
	}

	return res, nil
}
