package machine

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/drain"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	machinev1 "github.com/openshift/api/machine/v1beta1"

	"github.com/openshift/machine-api-operator/pkg/util/conditions"
)

const (
	nodeControlPlaneLabel = "node-role.kubernetes.io/control-plane"
	nodeMasterLabel       = "node-role.kubernetes.io/master"
)

// DrainController performs pods eviction for deleting node
type machineDrainController struct {
	client.Client
	config *rest.Config
	scheme *runtime.Scheme

	eventRecorder record.EventRecorder
}

// newDrainController returns a new reconcile.Reconciler for machine-drain-controller
func newDrainController(mgr manager.Manager) reconcile.Reconciler {
	d := &machineDrainController{
		Client:        mgr.GetClient(),
		eventRecorder: mgr.GetEventRecorderFor("machine-drain-controller"),
		config:        mgr.GetConfig(),
		scheme:        mgr.GetScheme(),
	}
	return d
}

// newDrainRateLimiter is based on the workqueue.DefaultControllerRateLimiter.
// The default rate limiter used by controller-runtime has a base delay of 5 milliseconds.
// As we know drains are a slower operation then traditional reconciles, we start with a
// larger base delay to allow the pods time for graceful shutdown.
// We cap out at 1000 seconds as with the default queue.
func newDrainRateLimiter() workqueue.RateLimiter {
	return workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Second, 1000*time.Second),
		// 10 qps, 100 bucket size.  This is only for retry speed and its only the overall factor (not per item)
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
	)
}

func (d *machineDrainController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	m := &machinev1.Machine{}
	if err := d.Client.Get(ctx, request.NamespacedName, m); err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return.
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	existingDrainedCondition := conditions.Get(m, machinev1.MachineDrained)
	alreadyDrained := existingDrainedCondition != nil && existingDrainedCondition.Status == corev1.ConditionTrue

	if !m.ObjectMeta.DeletionTimestamp.IsZero() && ptr.Deref(m.Status.Phase, "") == machinev1.PhaseDeleting && !alreadyDrained {
		drainFinishedCondition := conditions.TrueCondition(machinev1.MachineDrained)

		if _, exists := m.ObjectMeta.Annotations[ExcludeNodeDrainingAnnotation]; !exists && m.Status.NodeRef != nil {
			// pre-drain.delete lifecycle hook
			// Return early without error, will requeue if/when the hook owner removes the annotation.
			if len(m.Spec.LifecycleHooks.PreDrain) > 0 {
				klog.Infof("%v: not draining machine: lifecycle blocked by pre-drain hook", m.Name)
				d.eventRecorder.Eventf(m, corev1.EventTypeNormal, "DrainBlocked", "Drain blocked by pre-drain hook")
				return reconcile.Result{}, nil
			}
			d.eventRecorder.Eventf(m, corev1.EventTypeNormal, "DrainProceeds", "Node drain proceeds")
			if err := d.drainNode(ctx, m); err != nil {
				klog.Errorf("%v: failed to drain node for machine: %v", m.Name, err)
				conditions.Set(m, conditions.FalseCondition(
					machinev1.MachineDrained,
					machinev1.MachineDrainError,
					machinev1.ConditionSeverityWarning,
					"could not drain machine: %v", err,
				))
				d.eventRecorder.Eventf(m, corev1.EventTypeNormal, "DrainRequeued", "Node drain requeued: %v", err.Error())
				return delayIfRequeueAfterError(err)
			}
			d.eventRecorder.Eventf(m, corev1.EventTypeNormal, "DrainSucceeded", "Node drain succeeded")
			drainFinishedCondition.Message = "Drain finished successfully"
		} else {
			d.eventRecorder.Eventf(m, corev1.EventTypeNormal, "DrainSkipped", "Node drain skipped")
			drainFinishedCondition.Message = "Node drain skipped"
		}

		conditions.Set(m, drainFinishedCondition)
		// requeue request in case of failed update
		if err := d.Client.Status().Update(ctx, m); err != nil {
			return reconcile.Result{}, fmt.Errorf("could not update machine status: %w", err)
		}
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (d *machineDrainController) drainNode(ctx context.Context, machine *machinev1.Machine) error {
	kubeClient, err := kubernetes.NewForConfig(d.config)
	if err != nil {
		return fmt.Errorf("unable to build kube client: %v", err)
	}
	node, err := kubeClient.CoreV1().Nodes().Get(ctx, machine.Status.NodeRef.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If an admin deletes the node directly, we'll end up here.
			klog.Infof("Could not find node from noderef, it may have already been deleted: %v", machine.Status.NodeRef.Name)
			return nil
		}
		return fmt.Errorf("unable to get node %q: %v", machine.Status.NodeRef.Name, err)
	}

	if err := d.isDrainAllowed(ctx, node); err != nil {
		return fmt.Errorf("drain not permitted: %w", err)
	}

	drainer := &drain.Helper{
		Ctx:                 ctx,
		Client:              kubeClient,
		Force:               true,
		IgnoreAllDaemonSets: true,
		DeleteEmptyDirData:  true,
		GracePeriodSeconds:  -1,
		// If a pod is not evicted in 20 seconds, retry the eviction next time the
		// machine gets reconciled again (to allow other machines to be reconciled).
		Timeout: 20 * time.Second,
		OnPodDeletedOrEvicted: func(pod *corev1.Pod, usingEviction bool) {
			verbStr := "Deleted"
			if usingEviction {
				verbStr = "Evicted"
			}
			klog.Info(fmt.Sprintf("%s pod from Node", verbStr),
				"pod", fmt.Sprintf("%s/%s", pod.Name, pod.Namespace))
		},
		Out:    writer{klog.Info},
		ErrOut: writer{klog.Error},
	}

	if nodeIsUnreachable(node) {
		klog.Infof("%q: Node %q is unreachable, draining will ignore gracePeriod. PDBs are still honored.",
			machine.Name, node.Name)
		// Since kubelet is unreachable, pods will never disappear and we still
		// need SkipWaitForDeleteTimeoutSeconds so we don't wait for them.
		drainer.SkipWaitForDeleteTimeoutSeconds = skipWaitForDeleteTimeoutSeconds
		drainer.GracePeriodSeconds = 1
	}

	if err := drain.RunCordonOrUncordon(drainer, node, true); err != nil {
		// Can't cordon a node
		klog.Warningf("cordon failed for node %q: %v", node.Name, err)
		return &RequeueAfterError{RequeueAfter: 20 * time.Second}
	}

	if err := drain.RunNodeDrain(drainer, node.Name); err != nil {
		klog.Warningf("drain failed for machine %q: %v", machine.Name, err)

		// Make sure we return a regular error to take advantage of exponential backoff.
		// This will allow certain pods that need to finish work (eg static
		// installer pods) to complete even when being drained.
		// If we never allow the pods to complete, this can cause a deadlock between the
		// drain controller and installer pods.
		return err
	}

	klog.Infof("drain successful for machine %q", machine.Name)
	d.eventRecorder.Eventf(machine, corev1.EventTypeNormal, "Deleted", "Node %q drained", node.Name)

	return nil
}

// isDrainAllowed checks whether the drain is permitted at this time.
// It checks the following:
// - Is the node cordoned, if so allow draining to complete any previous attempt to drain.
// - Is the node a control plane node, if so, only allow draining if no other control plane node is already being drained.
func (d *machineDrainController) isDrainAllowed(ctx context.Context, node *corev1.Node) error {
	if node.Spec.Unschedulable {
		// If the node has already been cordoned, continue to drain.
		return nil
	}

	if !isControlPlaneNode(*node) {
		// We always allow draining of worker nodes.
		return nil
	}

	nodes := &corev1.NodeList{}
	if err := d.Client.List(ctx, nodes); err != nil {
		return fmt.Errorf("could not list control plane nodes: %v", err)
	}

	for _, otherNode := range nodes.Items {
		if isControlPlaneNode(otherNode) && otherNode.Spec.Unschedulable {
			klog.Warningf("Drain not permitted for node %q: found other control plane node (%s) already cordoned: other node may be being drained, do not continue until the other node is removed", node.Name, otherNode.Name)
			return &RequeueAfterError{RequeueAfter: 20 * time.Second}
		}
	}

	return nil
}

// isControlPlaneNode checks if the Node is labelled as a control plane node.
func isControlPlaneNode(node corev1.Node) bool {
	_, controlPlane := node.Labels[nodeControlPlaneLabel]
	_, master := node.Labels[nodeMasterLabel]

	return controlPlane || master
}
