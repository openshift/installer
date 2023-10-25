/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

type Indexer interface {
	// Index returns the index of the Indexer. The index can be passed to a registration.Index to
	// build an index for the controller-runtime client. If Index returns nil, there is nothing to index.
	// See controller-runtime mgr.GetFieldIndexer().IndexField() for more details.
	Index() []string
}
