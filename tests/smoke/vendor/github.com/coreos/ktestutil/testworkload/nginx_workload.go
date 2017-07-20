package testworkload

import (
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	utilrand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	batchv1 "k8s.io/client-go/pkg/apis/batch/v1"
	extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var (
	// PollTimeoutForNginx is the max duration for polling when using Nginx testworkload.
	PollTimeoutForNginx = 2 * time.Minute
	// PollIntervalForNginx is the interval between each condition check of poolling when using Nginx testworkload.
	PollIntervalForNginx = 5 * time.Second
)

// Nginx creates a temp nginx deployment/service pair
// that can be used as a test workload
type Nginx struct {
	Namespace string
	Name      string
	// List of pods that belong to the deployment
	Pods []*v1.Pod

	client          kubernetes.Interface
	podSelector     *metav1.LabelSelector
	nodeSelector    *metav1.LabelSelector
	pingPodSelector *metav1.LabelSelector
}

// NginxOpts defines func that applies custom options for Nginx
type NginxOpts func(*Nginx) error

// NewNginx create this nginx deployment/service pair.
// It waits until all the pods in the deployment are running.
func NewNginx(kc kubernetes.Interface, namespace string, options ...NginxOpts) (*Nginx, error) {
	//create random suffix
	name := fmt.Sprintf("nginx-%s", utilrand.String(5))

	n := &Nginx{
		Namespace: namespace,
		Name:      name,
		podSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"app": name}},
		nodeSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{},
		},
		pingPodSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{},
		},
		client: kc,
	}

	//apply options
	for _, option := range options {
		if err := option(n); err != nil {
			return nil, fmt.Errorf("error invalid options: %v", err)
		}
	}

	//create nginx deployment
	if err := n.newNginxDeployment(); err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, fmt.Errorf("error creating deployment %s: %v", n.Name, err)
	}
	if err := wait.PollImmediate(PollIntervalForNginx, PollTimeoutForNginx, func() (bool, error) {
		d, err := kc.ExtensionsV1beta1().Deployments(n.Namespace).Get(n.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if d.Status.UpdatedReplicas != d.Status.AvailableReplicas && d.Status.UnavailableReplicas != 0 {
			return false, nil
		}

		return true, nil
	}); err != nil {
		return nil, fmt.Errorf("deployment %s is not ready: %v", n.Name, err)
	}

	//wait for all pods to enter running phase
	if err := wait.PollImmediate(PollIntervalForNginx, PollTimeoutForNginx, func() (bool, error) {
		pl, err := kc.CoreV1().Pods(n.Namespace).List(metav1.ListOptions{
			LabelSelector: metav1.FormatLabelSelector(n.podSelector),
		})
		if err != nil {
			return false, err
		}

		if len(pl.Items) == 0 {
			return false, nil
		}

		var pods []*v1.Pod
		for i := range pl.Items {
			p := &pl.Items[i]
			if p.Status.Phase != v1.PodRunning {
				return false, nil
			}

			pods = append(pods, p)
		}

		n.Pods = pods
		return true, nil
	}); err != nil {
		return nil, fmt.Errorf("pods in deployment %s not ready: %v", n.Name, err)
	}

	//create nginx service
	if err := n.newNginxService(); err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, fmt.Errorf("error creating service %s: %v", n.Name, err)
	}

	return n, nil
}

// WithNginxSelector adds custom labels for Deployment's Selector field.
// Affects only Deployment pods.
func WithNginxSelector(labels map[string]string) NginxOpts {
	return func(n *Nginx) error {
		for k, v := range labels {
			n.podSelector.MatchLabels[k] = v
		}
		return nil
	}
}

// WithNginxNodeSelector adds custom labels for Pod's NodeSelector field
// Affects only Deployment's pods.
func WithNginxNodeSelector(labels map[string]string) NginxOpts {
	return func(n *Nginx) error {
		for k, v := range labels {
			n.nodeSelector.MatchLabels[k] = v
		}
		return nil
	}
}

// WithNginxPingJobLabels adds custom labels for PinJob's pods.
// Affects only PingJob's pods.
func WithNginxPingJobLabels(labels map[string]string) NginxOpts {
	return func(n *Nginx) error {
		for k, v := range labels {
			n.pingPodSelector.MatchLabels[k] = v
		}
		return nil
	}
}

// IsReachable pings the nginx service.
// Expects the nginx service to be reachable.
func (n *Nginx) IsReachable() error {
	if err := n.newPingPod(true); err != nil {
		return fmt.Errorf("error svc wasn't reachable: %v", err)
	}

	return nil
}

// IsUnReachable pings the nginx service.
// Expects the nginx service to be unreachable.
func (n *Nginx) IsUnReachable() error {
	if err := n.newPingPod(false); err != nil {
		return fmt.Errorf("error svc was reachable: %v", err)
	}

	return nil
}

// Delete deletes the deployment and service
func (n *Nginx) Delete() error {
	delPropPolicy := metav1.DeletePropagationForeground
	if err := wait.PollImmediate(PollIntervalForNginx, PollTimeoutForNginx, func() (bool, error) {
		if err := n.client.ExtensionsV1beta1().Deployments(n.Namespace).Delete(n.Name, &metav1.DeleteOptions{
			PropagationPolicy: &delPropPolicy,
		}); err != nil && !apierrors.IsNotFound(err) {
			return false, nil
		}

		if err := n.client.CoreV1().Services(n.Namespace).Delete(n.Name, &metav1.DeleteOptions{
			PropagationPolicy: &delPropPolicy,
		}); err != nil && !apierrors.IsNotFound(err) {
			return false, nil
		}
		return true, nil
	}); err != nil {
		return fmt.Errorf("error deleting %s deployment and serivce: %v", n.Name, err)
	}

	return nil
}

func (n *Nginx) newNginxDeployment() error {
	var (
		repl  int32 = 2
		cPort int32 = 80
	)
	d := &extensionsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.Name,
			Namespace: n.Namespace,
		},
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: &repl,
			Selector: n.podSelector,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: n.podSelector.MatchLabels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.12-alpine",
							Ports: []v1.ContainerPort{
								{
									ContainerPort: cPort,
								},
							},
						},
					},
					NodeSelector: n.nodeSelector.MatchLabels,
				},
			},
		},
	}
	if _, err := n.client.ExtensionsV1beta1().Deployments(n.Namespace).Create(d); err != nil {
		return err
	}

	return nil
}

func (n *Nginx) newNginxService() error {
	var (
		cPort int32 = 80
		tPort int32 = 80
	)
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      n.Name,
			Namespace: n.Namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: n.podSelector.MatchLabels,
			Ports: []v1.ServicePort{
				{
					Protocol:   v1.ProtocolTCP,
					Port:       cPort,
					TargetPort: intstr.FromInt(int(tPort)),
				},
			},
		},
	}
	if _, err := n.client.CoreV1().Services(n.Namespace).Create(svc); err != nil {
		return err
	}

	return nil
}

func (n *Nginx) newPingPod(reachable bool) error {
	name := fmt.Sprintf("%s-ping-job-%s", n.Name, utilrand.String(5))
	deadline := int64(PollTimeoutForNginx.Seconds())

	cmd := fmt.Sprintf("wget --timeout 5 %s", n.Name)
	if !reachable {
		cmd = fmt.Sprintf("! %s", cmd)
	}
	runcmd := []string{"/bin/sh", "-c", cmd}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: n.Namespace,
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds: &deadline,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: n.pingPodSelector.MatchLabels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    "ping-container",
							Image:   "alpine:3.6",
							Command: runcmd,
						},
					},
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
			},
		},
	}

	if _, err := n.client.BatchV1().Jobs(n.Namespace).Create(job); err != nil {
		return err
	}

	// wait for pod state
	if err := wait.PollImmediate(PollIntervalForNginx, PollTimeoutForNginx, func() (bool, error) {
		j, err := n.client.BatchV1().Jobs(n.Namespace).Get(job.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if j.Status.Succeeded < 1 {
			return false, nil
		}

		return true, nil
	}); err != nil {
		return fmt.Errorf("ping job didn't succeed: %v", err)
	}

	return nil
}
