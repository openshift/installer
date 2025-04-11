/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tele

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

// CorrIDKey is the type of the key used to store correlation
// IDs in context.Contexts.
type CorrIDKey string

// CorrIDKeyVal is the key used to store the correlation ID in
// context.Contexts, HTTP headers, and other similar locations.
const CorrIDKeyVal CorrIDKey = "x-ms-correlation-request-id"

// CorrID is a correlation ID that the cluster API provider
// sends with all API requests to Azure. Do not create one
// of these manually. Instead, use the CtxWithCorrelationID function
// to create one of these within a context.Context.
type CorrID string

// ctxWithCorrID creates a CorrID and creates a new context.Context
// with the new CorrID in it. It returns the _new_ context and the
// newly created CorrID. If there was a problem creating the correlation
// ID, the new context will not have the correlation ID in it and the
// returned CorrID will be the empty string.After you call this function, prefer to
// use the newly created context over the old one. Common usage is
// below:
//
//	ctx := context.Background()
//	ctx, newCorrID := ctxWithCorrID(ctx)
//	fmt.Println("new corr ID: ", newCorrID)
//	doSomething(ctx)
func ctxWithCorrID(ctx context.Context) (context.Context, CorrID) {
	if currentCorrIDIface := ctx.Value(CorrIDKeyVal); currentCorrIDIface != nil {
		currentCorrID, ok := currentCorrIDIface.(CorrID)
		if ok {
			return ctx, currentCorrID
		}
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, CorrID("")
	}
	newCorrID := CorrID(uid.String())
	ctx = context.WithValue(ctx, CorrIDKeyVal, newCorrID)
	return ctx, newCorrID
}

// CorrIDFromCtx attempts to fetch a correlation ID from the given
// context.Context. If none exists, returns an empty CorrID and false.
// Otherwise returns the CorrID value and true.
func CorrIDFromCtx(ctx context.Context) (CorrID, bool) {
	currentCorrIDIface := ctx.Value(CorrIDKeyVal)
	if currentCorrIDIface == nil {
		return CorrID(""), false
	}

	if corrID, ok := currentCorrIDIface.(CorrID); ok {
		return corrID, ok
	}

	return CorrID(""), false
}

// corrIDLogger attempts to fetch the correlation ID from the
// given ctx using CorrIDFromCtx. If it finds one, this function
// uses lggr.WithValues to return a new logr.Logger with the
// correlation ID in it.
func corrIDLogger(ctx context.Context, lggr logr.Logger) logr.Logger {
	corrID, ok := CorrIDFromCtx(ctx)
	if ok {
		lggr = lggr.WithValues(string(CorrIDKeyVal), string(corrID))
	}
	return lggr
}
