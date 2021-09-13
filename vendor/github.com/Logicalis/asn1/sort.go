package asn1

import "sort"

// isTagLessThan compares two tags (class + tag number)
// TODO: maybe a common Tag type can simplify that.
func isTagLessThan(c1, t1, c2, t2 uint) bool {
	if c1 == c2 {
		return t1 < t2
	}
	return c1 < c2
}

// rawValueSlice is a helper type to sort an slice of RawValues
type rawValueSlice []*rawValue

var _ sort.Interface = rawValueSlice{}

func (s rawValueSlice) Len() int      { return len(s) }
func (s rawValueSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s rawValueSlice) Less(i, j int) bool {
	return isTagLessThan(s[i].Class, s[i].Tag, s[j].Class, s[j].Tag)
}

// expectedFieldElementSlice is a helper type to sort an slice of expectedFieldElement
type expectedFieldElementSlice []expectedFieldElement

var _ sort.Interface = expectedFieldElementSlice{}

func (s expectedFieldElementSlice) Len() int      { return len(s) }
func (s expectedFieldElementSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s expectedFieldElementSlice) Less(i, j int) bool {
	return isTagLessThan(s[i].class, s[i].tag, s[j].class, s[j].tag)
}
