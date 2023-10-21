package rolesanywhere

func expandStringList(tfList []interface{}) []string {
	var result []string

	for _, rawVal := range tfList {
		if v, ok := rawVal.(string); ok && v != "" {
			result = append(result, v)
		}
	}

	return result
}
