package verify

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
)

// SignatureSource provides a set of signatures by digest to save.
type SignatureSource interface {
	// Signatures returns a list of valid signatures for release digests.
	Signatures() map[string][][]byte
}

// PersistentSignatureStore is a store that can save signatures for
// later recovery.
type PersistentSignatureStore interface {
	// Store saves the provided signatures or return an error. If context
	// reaches its deadline the store should be cancelled.
	Store(ctx context.Context, signatures map[string][][]byte) error
}

// StorePersister saves signatures into store periodically.
type StorePersister struct {
	store      PersistentSignatureStore
	signatures SignatureSource
}

// NewSignatureStorePersister creates an instance that can save signatures into the destination
// store.
func NewSignatureStorePersister(dst PersistentSignatureStore, src SignatureSource) *StorePersister {
	return &StorePersister{
		store:      dst,
		signatures: src,
	}
}

// Run flushes signatures to the provided store every interval or until the context is finished.
// After context is done, it runs one more time to attempt to flush the current state. It does not
// return until that last store completes.
func (p *StorePersister) Run(ctx context.Context, interval time.Duration) {
	wait.Until(func() {
		if err := p.store.Store(ctx, p.signatures.Signatures()); err != nil {
			klog.Warningf("Unable to save signatures: %v", err)
		}
	}, interval, ctx.Done())

	if err := p.store.Store(context.Background(), p.signatures.Signatures()); err != nil {
		klog.Warningf("Unable to save signatures during final flush: %v", err)
	}
}
