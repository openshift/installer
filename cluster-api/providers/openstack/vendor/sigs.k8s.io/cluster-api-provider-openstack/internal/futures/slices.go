package futures

import stdlib_slices "slices"

// SlicesConcat vendors the Go v1.22 slices.Concat function.
func SlicesConcat[S ~[]E, E any](slices ...S) S {
	size := 0
	for _, s := range slices {
		size += len(s)
		if size < 0 {
			panic("len out of range")
		}
	}
	newslice := stdlib_slices.Grow[S](nil, size)
	for _, s := range slices {
		newslice = append(newslice, s...)
	}
	return newslice
}
