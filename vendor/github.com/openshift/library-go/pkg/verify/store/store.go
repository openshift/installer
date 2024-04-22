// Package store defines generic interfaces for signature stores.
package store

import (
	"context"
	"errors"
)

// Callback returns true if an acceptable signature has been found, or
// an error if the loop should be aborted.  If there was a problem
// retrieving the signature, the incoming error will describe the
// problem and the function can decide how to handle that error.
type Callback func(ctx context.Context, signature []byte, errIn error) (done bool, err error)

// ErrNotFound is a base error for Callback, to be used when the store
// decides one signature-retrieval avenue is exhausted.
var ErrNotFound = errors.New("no more signatures to check")

// Store provides access to signatures by digest.
type Store interface {

	// Signatures fetches signatures for the provided digest, feeding
	// them into the provided callback until an acceptable signature is
	// found or an error occurs.
	//
	// Not finding additional signatures should result in a callback
	// call with an error wrapping ErrNotFound, to allow the caller to
	// figure out when and why the store was unable to find a signature.
	// When a store has several lookup mechanisms, this may result in
	// several callback calls with different ErrNotFound.  Signatures
	// itself should return nil in this case, because eventually running
	// out of signatures is an expected part of any invocation where the
	// callback calls never return done=true.
	Signatures(ctx context.Context, name string, digest string, fn Callback) error

	// String returns a description of where this store finds
	// signatures.  The string is a short clause intended for display in
	// a description of the verifier.
	String() string
}
