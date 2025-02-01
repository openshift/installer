package destroy

import (
	"context"
	"io"
	"math"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/drain"
)

// Volume represents the necessary fields to delete all the persistent volumes from
// a cluster that is about to be destroyed.
type Volume struct {
	persistentVolumeClaimList *corev1.PersistentVolumeClaimList
	persistentVolumeList      *corev1.PersistentVolumeList
	clientSet                 *kubernetes.Clientset
	logger                    logrus.FieldLogger

	mu sync.Mutex
}

func newKubeClientSet(kubeConfig string) (*kubernetes.Clientset, error) {
	auth, err := getKubeConfigAuth(kubeConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get auth kubeconfig")
	}
	config, err := clientcmd.RESTConfigFromKubeConfig(auth)
	if err != nil {
		return nil, errors.Wrapf(err, "could not get config from kubeconfig")
	}
	return kubernetes.NewForConfig(config)
}

// NewVolume sets up the Volume struct, first configuring the kubernetes client set, then
// retrieving the persistent volumes and their claims.
func NewVolume(ctx context.Context, logger logrus.FieldLogger, kubeConfig string) (*Volume, error) {
	v := &Volume{}
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	clientSet, err := newKubeClientSet(kubeConfig)
	if err != nil {
		return nil, err
	}
	v.clientSet = clientSet

	pvcList, err := clientSet.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	v.persistentVolumeClaimList = pvcList.DeepCopy()

	pvList, err := clientSet.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	v.persistentVolumeList = pvList.DeepCopy()
	v.logger = logger

	return v, nil
}

func (v *Volume) drainNodes(ctx context.Context) error {
	timeout, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	v.logger.Info("Draining and cordoning nodes, this might take a few moments.")
	drainHelper := &drain.Helper{
		Ctx:                             timeout,
		Client:                          v.clientSet,
		Force:                           true,
		GracePeriodSeconds:              0,
		IgnoreAllDaemonSets:             true,
		Timeout:                         0,
		DeleteEmptyDirData:              true,
		Selector:                        "",
		PodSelector:                     "",
		ChunkSize:                       0,
		DisableEviction:                 true,
		SkipWaitForDeleteTimeoutSeconds: 0,
		AdditionalFilters:               nil,
		// do we need any of this output?
		Out:                             io.Discard,
		ErrOut:                          io.Discard,
		DryRunStrategy:                  0,
		OnPodDeletedOrEvicted:           nil,
		OnPodDeletionOrEvictionFinished: func(pod *corev1.Pod, usingEviction bool, err error) {},
		OnPodDeletionOrEvictionStarted:  func(pod *corev1.Pod, usingEviction bool) {},
	}

	nodeList, err := v.clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Need all nodes cordoned before draining so there is no
	// available resources for the pod to restart
	for _, node := range nodeList.Items {
		if _, ok := node.Labels["node-role.kubernetes.io/worker"]; ok {
			v.logger.Debugf("cordoning node: %s", node.Name)
			if err := drain.RunCordonOrUncordon(drainHelper, &node, true); err != nil {
				return err
			}
		}
	}
	for _, node := range nodeList.Items {
		if _, ok := node.Labels["node-role.kubernetes.io/worker"]; ok {
			v.logger.Debugf("draining node: %s", node.Name)
			if err := drain.RunNodeDrain(drainHelper, node.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Volume) setupVolumeSharedInformer() (informers.SharedInformerFactory, error) {
	sharedInformer := informers.NewSharedInformerFactory(v.clientSet, time.Second*15)

	pvcInformer := sharedInformer.Core().V1().PersistentVolumeClaims()
	pvInformer := sharedInformer.Core().V1().PersistentVolumes()

	_, err := pvcInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			if pvc, ok := obj.(*corev1.PersistentVolumeClaim); ok {
				v.logger.Debugf("informer: deleting volume claim %s", pvc.Name)
				v.mu.Lock()
				defer v.mu.Unlock()
				v.persistentVolumeClaimList.Items = slices.DeleteFunc(v.persistentVolumeClaimList.Items, func(dfpvc corev1.PersistentVolumeClaim) bool {
					return dfpvc.Name == pvc.Name
				})
			}
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = pvInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			if pv, ok := obj.(*corev1.PersistentVolume); ok {
				v.logger.Debugf("informer: deleting volume %s", pv.Name)
				v.mu.Lock()
				defer v.mu.Unlock()

				v.persistentVolumeList.Items = slices.DeleteFunc(v.persistentVolumeList.Items, func(dfpv corev1.PersistentVolume) bool {
					return dfpv.Name == pv.Name
				})
			}
		},
	})
	if err != nil {
		return nil, err
	}

	return sharedInformer, nil
}

// deletePersistentVolumes sets up two shared informers, one to monitor the deletion of persistent volumes
// the other persistent volume claims. Then uses ExponentialBackoffWithContext to start the removal
// of PVCs and check when all the persistent volumes are removed.
func (v *Volume) deletePersistentVolumes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	var duration time.Duration
	duration = math.MaxInt64

	deadline, ok := ctx.Deadline()
	if ok {
		duration = time.Until(deadline)
	}

	sharedInformer, err := v.setupVolumeSharedInformer()
	if err != nil {
		return err
	}

	stopCh := make(chan struct{})
	sharedInformer.Start(stopCh)
	sharedInformer.WaitForCacheSync(stopCh)

	runOnce := 0
	backoff := wait.Backoff{
		Duration: time.Second * 20,
		Factor:   1.1,
		Jitter:   0,
		Steps:    math.MaxInt,
		Cap:      duration,
	}

	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(fctx context.Context) (bool, error) {
		v.logger.Debugf("ExponentailBackoff, Step %d", runOnce)

		v.mu.Lock()
		defer v.mu.Unlock()

		// We need to delete only once and after the shared informers are started.
		if runOnce == 0 {
			for _, pvc := range v.persistentVolumeClaimList.Items {
				v.logger.Debugf("deleting volume claim %s", pvc.Name)
				if err := v.clientSet.CoreV1().PersistentVolumeClaims(pvc.Namespace).Delete(fctx, pvc.Name, metav1.DeleteOptions{}); err != nil {
					return false, err
				}
			}
		}

		// Once the list of persistent volumes is empty we can stop the shared informers
		// and exit this loop.
		if len(v.persistentVolumeList.Items) == 0 {
			v.logger.Debugf("All volumes are deleted")
			stopCh <- struct{}{}
			return true, nil
		}

		runOnce++
		return false, nil
	})

	if err != nil {
		if wait.Interrupted(err) {
			v.logger.Error("timeout waiting for all volumes to be deleted")
			v.mu.Lock()
			defer v.mu.Unlock()

			for _, pvc := range v.persistentVolumeClaimList.Items {
				v.logger.Errorf("%s persistent volume claim was not deleted", pvc.Name)
			}

			for _, pv := range v.persistentVolumeList.Items {
				v.logger.Errorf("%s persistent volume was not deleted", pv.Name)
			}
		} else {
			return err
		}
	}

	return nil
}

func getKubeConfigAuth(kubeConfig string) ([]byte, error) {
	_, err := os.Stat(kubeConfig)
	if err != nil {
		return nil, err
	}

	auth, err := os.ReadFile(kubeConfig)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
