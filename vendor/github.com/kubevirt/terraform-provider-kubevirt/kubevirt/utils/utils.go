package utils

func ConvertMap(src map[string]interface{}) map[string]string {
	result := map[string]string{}
	for k, v := range src {
		result[k] = v.(string)
	}
	return result
}
