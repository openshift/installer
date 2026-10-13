/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package kubeclient

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type andIndexer struct {
	indexers []client.FieldIndexer
}

func NewAndIndexer(indexers ...client.FieldIndexer) client.FieldIndexer {
	return &andIndexer{
		indexers: indexers,
	}
}

var _ client.FieldIndexer = &andIndexer{}

func (a *andIndexer) IndexField(ctx context.Context, obj client.Object, field string, extractValue client.IndexerFunc) error {
	for _, indexer := range a.indexers {
		err := indexer.IndexField(ctx, obj, field, extractValue)
		if err != nil {
			return err
		}
	}

	return nil
}
