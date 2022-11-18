// Package parallel combines several signature stores in a single store.
// Signatures are searched in each substore simultaneously until a
// match is found.
package parallel

import (
	"context"
	"fmt"
	"strings"

	"github.com/openshift/library-go/pkg/verify/store"
)

type signatureResponse struct {
	signature []byte
	errIn     error
}

// Store provides access to signatures stored in sub-stores.
type Store struct {
	Stores []store.Store
}

// Signatures fetches signatures for the provided digest.
func (s *Store) Signatures(ctx context.Context, name string, digest string, fn store.Callback) error {
	nestedCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	responses := make(chan signatureResponse, len(s.Stores))
	errorChannelCount := 0
	errorChannel := make(chan error, 1)

	for i := range s.Stores {
		errorChannelCount++
		go func(ctx context.Context, wrappedStore store.Store, name string, digest string, responses chan signatureResponse, errorChannel chan error) {
			errorChannel <- wrappedStore.Signatures(ctx, name, digest, func(ctx context.Context, signature []byte, errIn error) (done bool, err error) {
				select {
				case <-ctx.Done():
					select {
					case responses <- signatureResponse{signature: signature, errIn: errIn}:
					default:
					}
					return false, ctx.Err()
				case responses <- signatureResponse{signature: signature, errIn: errIn}:
				}
				return false, nil
			})
		}(nestedCtx, s.Stores[i], name, digest, responses, errorChannel)
	}

	allDone := false
	var loopError error
	for errorChannelCount > 0 {
		if allDone {
			err := <-errorChannel
			errorChannelCount--
			if loopError == nil && err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				loopError = err
			}
		} else {
			select {
			case response := <-responses:
				done, err := fn(ctx, response.signature, response.errIn)
				if done || err != nil {
					allDone = true
					loopError = err
					cancel()
				}
			case err := <-errorChannel:
				errorChannelCount--
				if loopError == nil && err != nil && err != context.Canceled && err != context.DeadlineExceeded {
					loopError = err
				}
			}
		}
	}
	close(responses)
	close(errorChannel)
	if loopError != nil {
		return loopError
	}

	if err := ctx.Err(); err != nil {
		return err // because we discard context errors from the wrapped stores
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
	return fmt.Sprintf("parallel signature store wrapping %s", wrapped)
}
