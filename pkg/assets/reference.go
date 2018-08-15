package assets

import (
	"fmt"
)

// Reference holds an Asset reference.
type Reference struct {
	// Name holds the the name of the referenced asset.  This value can
	// be used to retrieve the latest value of that asset's state.
	Name string `json:"name,omitempty"`

	// Hash holds the hash of the referenced asset.  This value pins a
	// specific value of that asset's state.
	Hash []byte `json:"hash,omitempty"`
}

type referenceSlice []Reference

func (ref *Reference) String() string {
	return fmt.Sprintf("%q (%x)", ref.Name, ref.Hash)
}

func (refs referenceSlice) Len() int {
	return len(refs)
}

func (refs referenceSlice) Less(i, j int) bool {
	return refs[i].Name < refs[i].Name
}

func (refs referenceSlice) Swap(i, j int) {
	refs[i], refs[j] = refs[j], refs[i]
}
