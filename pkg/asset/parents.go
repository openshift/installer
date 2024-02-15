package asset

import (
	"fmt"
	"reflect"
)

// Parents is the collection of assets upon which another asset is directly
// dependent.
type Parents map[reflect.Type]Asset

// Add adds the specified assets to the parents collection.
func (p Parents) Add(assets ...Asset) {
	for _, a := range assets {
		p[reflect.TypeOf(a)] = a
	}
}

// Get populates the state of the specified assets with the state stored in the
// parents collection.
func (p Parents) Get(assets ...Asset) {
	for _, a := range assets {
		fmt.Printf("finding asset for type: %s\n", reflect.TypeOf(a))
		if val, exists := p[reflect.TypeOf(a)]; exists {
			reflect.ValueOf(a).Elem().Set(reflect.ValueOf(val).Elem())
		} else {
			panic(fmt.Sprintf("unable to find asset of type: %v", reflect.TypeOf(a)))
		}
	}
}
