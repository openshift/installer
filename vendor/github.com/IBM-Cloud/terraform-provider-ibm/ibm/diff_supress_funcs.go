// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {

	if old == "" {
		return false
	}
	var oldObj, newObj []map[string]interface{}
	err := json.Unmarshal([]byte(old), &oldObj)
	if err != nil {
		log.Printf("Error unmarshalling old json :: %s", err.Error())
		return false
	}
	err = json.Unmarshal([]byte(new), &newObj)
	if err != nil {
		log.Printf("Error unmarshalling new json :: %s", err.Error())
		return false
	}

	oldm := make(map[interface{}]interface{})
	newm := make(map[interface{}]interface{})

	for _, m := range oldObj {
		oldm[m["key"]] = m["value"]
	}
	for _, m := range newObj {
		newm[m["key"]] = m["value"]
	}
	return reflect.DeepEqual(oldm, newm)
}
