package nutanix

import (
	"log"
	"sort"
)

func resourceNutanixCategoriesMigrateState(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if len(rawState) == 0 || rawState == nil {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return rawState, nil
	}

	keys := make([]string, 0, len(rawState))
	for k := range rawState {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	log.Printf("[DEBUG] meta: %#v", meta)
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)

	if l, ok := rawState["categories"]; ok {
		if assertedL, ok := l.(map[string]interface{}); ok {
			c := make([]interface{}, 0)
			keys := make([]string, 0, len(assertedL))
			for k := range assertedL {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, name := range keys {
				value := assertedL[name]
				c = append(c, map[string]interface{}{
					"name":  name,
					"value": value.(string),
				})
			}
			rawState["categories"] = c
		}
	}
	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}
