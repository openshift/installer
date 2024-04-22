// Package serial combines several signature stores in a single store.
// Signatures are searched in each substore in turn until a match is
// found.
package serial

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift/library-go/pkg/verify/store"
)

// Store provides access to signatures stored in sub-stores.
type Store struct {
	Stores []store.Store
}

// Signatures fetches signatures for the provided digest.
func (s *Store) Signatures(ctx context.Context, name string, digest string, fn store.Callback) error {
	allDone := false

	wrapper := func(ctx context.Context, signature []byte, errIn error) (done bool, err error) {
		done, err = fn(ctx, signature, errIn)
		if done {
			allDone = true
		}
		return done, err
	}

	for _, store := range s.Stores {
		err := store.Signatures(ctx, name, digest, wrapper)
		if err != nil || allDone {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}

	_, err := fn(ctx, nil, fmt.Errorf("%s: %w", s.String(), store.ErrNotFound))
	return err
}

// String returns a description of where this store finds
// signatures.
func (s *Store) String() string {
	wrapped := "no stores"
	if len(s.Stores) > 0 {
		names := make([]string, 0, len(s.Stores))
		for _, store := range s.Stores {
			names = append(names, store.String())
		}
		wrapped = strings.Join(names, ", ")
	}
	return fmt.Sprintf("serial signature store wrapping %s", wrapped)
}
