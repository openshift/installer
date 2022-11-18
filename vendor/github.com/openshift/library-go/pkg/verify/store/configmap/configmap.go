package configmap

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"

	"github.com/openshift/library-go/pkg/verify/store"
	"github.com/openshift/library-go/pkg/verify/util"
)

const (
	// NamespaceLabelConfigMap is the Namespace label applied to a configmap
	// containing signatures.
	NamespaceLabelConfigMap = "openshift-config-managed"

	// ReleaseLabelConfigMap is a label applied to a configmap inside the
	// openshift-config-managed namespace that indicates it contains signatures
	// for release image digests. Any binaryData key that starts with the digest
	// is added to the list of signatures checked.
	ReleaseLabelConfigMap = "release.openshift.io/verification-signatures"
)

// Store abstracts retrieving signatures from config maps on a cluster.
type Store struct {
	client corev1client.ConfigMapsGetter
	ns     string

	limiter *rate.Limiter
	lock    sync.Mutex
	last    []corev1.ConfigMap
}

// NewStore returns a store that can retrieve or persist signatures on a
// cluster. If limiter is not specified it defaults to one call every 30 seconds.
func NewStore(client corev1client.ConfigMapsGetter, limiter *rate.Limiter) *Store {
	if limiter == nil {
		limiter = rate.NewLimiter(rate.Every(30*time.Second), 1)
	}
	return &Store{
		client:  client,
		ns:      NamespaceLabelConfigMap,
		limiter: limiter,
	}
}

// String displays information about this source for human review.
func (s *Store) String() string {
	return fmt.Sprintf("config maps in %s with label %q", s.ns, ReleaseLabelConfigMap)
}

// rememberMostRecentConfigMaps stores a set of config maps containing
// signatures.
func (s *Store) rememberMostRecentConfigMaps(last []corev1.ConfigMap) {
	names := make([]string, 0, len(last))
	for _, cm := range last {
		names = append(names, cm.ObjectMeta.Name)
	}
	sort.Strings(names)
	s.lock.Lock()
	defer s.lock.Unlock()
	klog.V(4).Infof("remember most recent signature config maps: %s", strings.Join(names, " "))
	s.last = last
}

// mostRecentConfigMaps returns the last cached version of config maps
// containing signatures.
func (s *Store) mostRecentConfigMaps() []corev1.ConfigMap {
	s.lock.Lock()
	defer s.lock.Unlock()
	klog.V(4).Info("use cached most recent signature config maps")
	return s.last
}

// Signatures fetches signatures for the provided digest
// out of config maps labelled with ReleaseLabelConfigMap in the
// NamespaceLabelConfigMap namespace.
func (s *Store) Signatures(ctx context.Context, name string, digest string, fn store.Callback) error {
	// avoid repeatedly reloading config maps
	items := s.mostRecentConfigMaps()
	r := s.limiter.Reserve()
	if items == nil || r.OK() {
		configMaps, err := s.client.ConfigMaps(s.ns).List(ctx, metav1.ListOptions{
			LabelSelector: ReleaseLabelConfigMap,
		})
		if err != nil {
			s.rememberMostRecentConfigMaps([]corev1.ConfigMap{})
			return err
		}
		items = configMaps.Items
		s.rememberMostRecentConfigMaps(configMaps.Items)
	}

	prefix, err := util.DigestToKeyPrefix(digest, "-")
	if err != nil {
		return err
	}

	for _, cm := range items {
		klog.V(4).Infof("searching for %s in signature config map %s", prefix, cm.ObjectMeta.Name)
		for k, v := range cm.BinaryData {
			if strings.HasPrefix(k, prefix) {
				klog.V(4).Infof("key %s from signature config map %s matches %s", k, cm.ObjectMeta.Name, digest)
				done, err := fn(ctx, v, nil)
				if err != nil || done {
					return err
				}
				if err := ctx.Err(); err != nil {
					return err
				}
			}
		}
		if done, err := fn(ctx, nil, fmt.Errorf("prefix %s in config map %s: %w", prefix, cm.ObjectMeta.Name, store.ErrNotFound)); err != nil || done {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	return nil
}

// Store attempts to persist the provided signatures into a form Signatures will
// retrieve.
func (s *Store) Store(ctx context.Context, signaturesByDigest map[string][][]byte) error {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: s.ns,
			Name:      "signatures-managed",
			Labels: map[string]string{
				ReleaseLabelConfigMap: "",
			},
		},
		BinaryData: make(map[string][]byte),
	}
	count := 0
	for digest, signatures := range signaturesByDigest {
		prefix, err := util.DigestToKeyPrefix(digest, "-")
		if err != nil {
			return err
		}
		for i := 0; i < len(signatures); i++ {
			cm.BinaryData[fmt.Sprintf("%s-%d", prefix, i)] = signatures[i]
			count += 1
		}
	}
	return retry.OnError(
		retry.DefaultRetry,
		func(err error) bool { return errors.IsConflict(err) || errors.IsAlreadyExists(err) },
		func() error {
			existing, err := s.client.ConfigMaps(s.ns).Get(ctx, cm.Name, metav1.GetOptions{})
			if errors.IsNotFound(err) {
				_, err := s.client.ConfigMaps(s.ns).Create(ctx, cm, metav1.CreateOptions{})
				if err != nil {
					klog.V(4).Infof("create signature cache config map %s in namespace %s with %d signatures", cm.ObjectMeta.Name, s.ns, count)
				}
				return err
			}
			if err != nil {
				return err
			}
			existing.Labels = cm.Labels
			existing.BinaryData = cm.BinaryData
			existing.Data = cm.Data
			_, err = s.client.ConfigMaps(s.ns).Update(ctx, existing, metav1.UpdateOptions{})
			if err != nil {
				klog.V(4).Infof("update signature cache config map %s in namespace %s with %d signatures", cm.ObjectMeta.Name, s.ns, count)
			}
			return err
		},
	)
}
