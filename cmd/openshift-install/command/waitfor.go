package command

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	configv2 "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	clientwatch "k8s.io/client-go/tools/watch"
	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	configlisters "github.com/openshift/client-go/config/listers/config/v1"
	routeclient "github.com/openshift/client-go/route/clientset/versioned"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/installconfig"
	timer "github.com/openshift/installer/pkg/metrics/timer"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/baremetal"
	cov1helpers "github.com/openshift/library-go/pkg/config/clusteroperator/v1helpers"
	"github.com/openshift/library-go/pkg/route/routeapihelpers"
)

const (
	// ExitCodeInstallConfigError is used when there is a install-config error.
	ExitCodeInstallConfigError = iota + 3
	// ExitCodeInfrastructureFailed is used there is a infrastructure error.
	ExitCodeInfrastructureFailed
	// ExitCodeBootstrapFailed is used when bootstrap failed.
	ExitCodeBootstrapFailed
	// ExitCodeInstallFailed is used when cluster installation failed.
	ExitCodeInstallFailed
	// ExitCodeOperatorStabilityFailed is used when operator stability check failed.
	ExitCodeOperatorStabilityFailed
	// ExitCodeInterrupt is used when the interrupt signal was received.
	ExitCodeInterrupt

	// coStabilityThreshold is how long a cluster operator must have Progressing=False
	// in order to be considered stable. Measured in seconds.
	coStabilityThreshold float64 = 30
)

// SkipPasswordPrintFlag when true means do not print the generated user password.
var SkipPasswordPrintFlag bool

// WaitForInstallComplete waits for cluster to complete installation, checks for operator stability
// and logs cluster information when successful.
func WaitForInstallComplete(ctx context.Context, config *rest.Config, assetstore asset.Store) error {
	// FIXME: Register the worker nodes to target group of ingress LB.
	// Remove after CCM support dualstack NLB.
	if err := waitForWorkerNodesAvailability(ctx, config, assetstore); err != nil {
		return err
	}

	if err := waitForInitializedCluster(ctx, config, assetstore); err != nil {
		return err
	}

	if err := addRouterCAToClusterCA(ctx, config, RootOpts.Dir); err != nil {
		return err
	}

	if err := waitForStableOperators(ctx, config); err != nil {
		return err
	}

	consoleURL, err := getConsole(ctx, config)
	if err != nil {
		logrus.Warnf("Cluster does not have a console available: %v", err)
	}

	return logComplete(RootOpts.Dir, consoleURL)
}

// waitForWorkerNodesAvailability waits for worker nodes to be running and register them with the TargetGroup of the ingress NLB.
// NOTE: This should be handled by the CCM, not the installer.
func waitForWorkerNodesAvailability(ctx context.Context, config *rest.Config, assetstore asset.Store) error {
	timer.StartTimer("CCM: Worker nodes Available")
	icAsset, err := assetstore.Load(&installconfig.InstallConfig{})
	if err != nil {
		return fmt.Errorf("failed to load installconfig: %w", err)
	}
	if icAsset == nil {
		return fmt.Errorf("failed to installconfig: received nil")
	}

	ic := icAsset.(*installconfig.InstallConfig).Config

	// Nothing to do!
	if ic.Platform.AWS == nil || !ic.AWS.DualStackEnabled() {
		return nil
	}

	region := ic.Platform.AWS.Region

	// FIXME: Ignore edge compute pool
	numOfNodes := ptr.Deref(ic.Compute[0].Replicas, 0) + ptr.Deref(ic.ControlPlane.Replicas, 0)
	if numOfNodes == 0 {
		// nothing to do, but should not happen
		return nil
	}

	nodeCheckDuration := 10 * time.Minute
	nodeContext, cancel := context.WithTimeout(ctx, nodeCheckDuration)
	defer cancel()

	untilTime := time.Now().Add(nodeCheckDuration)
	timezone, _ := untilTime.Zone()
	logrus.Infof("CCM: Waiting up to %v (until %v %s) to ensure worker nodes are available and registered with ingress LB...",
		nodeCheckDuration, untilTime.Format(time.Kitchen), timezone)

	cc, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create a config client: %w", err)
	}
	configInformers := informers.NewSharedInformerFactory(cc, 0)
	nodeInformer := configInformers.Core().V1().Nodes().Informer()
	nodeLister := configInformers.Core().V1().Nodes().Lister()
	configInformers.Start(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(), nodeInformer.HasSynced) {
		return fmt.Errorf("informers never started")
	}

	// Create clients to call AWS API
	// FIXME: Let's ignore the custom endpoints for now
	cfg, err := configv2.LoadDefaultConfig(ctx, configv2.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	elbv2Client := elbv2.NewFromConfig(cfg)
	ec2Client := ec2.NewFromConfig(cfg)

	waitErr := wait.PollUntilContextCancel(nodeContext, 1*time.Second, true, func(ctx context.Context) (done bool, err error) {
		// If the expected number of nodes are running, proceed.
		// Otherwise, requeue.
		nodes, err := nodeLister.List(labels.Everything())
		if err != nil {
			return false, fmt.Errorf("failed to get nodes: %w", err)
		}
		if len(nodes) < int(numOfNodes) {
			return false, nil
		}

		// Convert nodes to EC2 instance IDs
		var instanceIDs []string
		for _, node := range nodes {
			url, err := url.Parse(node.Spec.ProviderID)
			if err != nil {
				return false, fmt.Errorf("invalid node provider ID (%s): %w", node.Spec.ProviderID, err)
			}
			if url.Scheme != "aws" {
				return false, fmt.Errorf("invalid scheme for AWS instance (%s)", node.Spec.ProviderID)
			}

			awsID := ""
			tokens := strings.Split(strings.Trim(url.Path, "/"), "/")
			// last token in the providerID is the aws resource ID for both EC2 and Fargate nodes
			if len(tokens) > 0 {
				awsID = tokens[len(tokens)-1]
			}
			instanceIDs = append(instanceIDs, awsID)
		}

		// Enable Primary IPv6 flag if dual-stack IPv6 primary
		// FIXME: This should be done by MAPI when creating the instances.
		if ic.AWS.IPFamily == aws.DualStackIPv6Primary {
			instances, err := ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
				InstanceIds: instanceIDs,
			})
			if err != nil {
				return false, fmt.Errorf("failed to describe ec2 instances for nodes: %w", err)
			}

			// FIXME: Assume all instances have at least 1 ENI with 1 as primary ENI
			for _, reservation := range instances.Reservations {
				for _, instance := range reservation.Instances {
					var primaryENI ec2types.InstanceNetworkInterface
					for _, eni := range instance.NetworkInterfaces {
						if ptr.Deref(eni.Attachment.DeviceIndex, 0) == 0 {
							primaryENI = eni
							break
						}
					}

					enabled := true
					_, err := ec2Client.ModifyNetworkInterfaceAttribute(ctx, &ec2.ModifyNetworkInterfaceAttributeInput{
						EnablePrimaryIpv6:  &enabled,
						NetworkInterfaceId: primaryENI.NetworkInterfaceId,
					})
					if err != nil {
						return false, fmt.Errorf("failed to set primary ipv6 for instance %s: %w", *instance.InstanceId, err)
					}
				}
			}
		}

		// Get the NodePort service for default ingress
		// Reference: oc -n openshift-ingress get svc router-nodeport-default -o=wide
		svc, err := cc.CoreV1().Services("openshift-ingress").Get(nodeContext, "router-nodeport-default", metav1.GetOptions{})
		if err != nil {
			// The service is not yet created by CIO, wait a bit more
			if apierrors.IsNotFound(err) {
				return false, nil
			}
			return false, fmt.Errorf("failed to get service openshift-ingress/router-nodeport-default: %w", err)
		}

		// Check if nodes are already registered. If true, nothing to do.
		targetGrp, err := elbv2Client.DescribeTargetGroups(nodeContext, &elbv2.DescribeTargetGroupsInput{
			Names: []string{fmt.Sprintf("%.20s-ingress-%d", ic.GetName(), 443), fmt.Sprintf("%.20s-ingress-%d", ic.GetName(), 80)},
		})
		if err != nil {
			return false, fmt.Errorf("failed to describe target group: %w", err)
		}

		for _, tg := range targetGrp.TargetGroups {
			// Get registered targets by querying target health API.
			targetDesc, err := elbv2Client.DescribeTargetHealth(nodeContext, &elbv2.DescribeTargetHealthInput{
				TargetGroupArn: tg.TargetGroupArn,
			})
			if err != nil {
				return false, fmt.Errorf("failed to get target health for target group %s: %w", *tg.TargetGroupArn, err)
			}

			// No targets found. We will register nodes as targets.
			if len(targetDesc.TargetHealthDescriptions) == 0 {
				var nodeport int32
				for _, port := range svc.Spec.Ports {
					if port.Port == *tg.Port {
						nodeport = port.NodePort
						break
					}
				}

				var targets []elbv2types.TargetDescription
				for _, instanceID := range instanceIDs {
					targets = append(targets, elbv2types.TargetDescription{
						Port: &nodeport,
						Id:   &instanceID,
					})
				}

				// Register worker and controlplane nodes with TargetGroup of ingress NLB.
				_, err = elbv2Client.RegisterTargets(nodeContext, &elbv2.RegisterTargetsInput{
					TargetGroupArn: tg.TargetGroupArn,
					Targets:        targets,
				})
				if err != nil {
					return false, fmt.Errorf("failed to register nodes to ingress LB: %w", err)
				}
			}
		}

		return true, nil
	})
	if waitErr != nil {
		return fmt.Errorf("failed to wait for worker node availability: %w", waitErr)
	}

	timer.StopTimer("CCM: Worker nodes Available")

	logrus.Info("CCM: Worker nodes are available and registered with ingress LB")
	return nil
}

// waitForInitializedCluster watches the ClusterVersion waiting for confirmation
// that the cluster has been initialized.
func waitForInitializedCluster(ctx context.Context, config *rest.Config, assetstore asset.Store) error {
	// TODO revert this value back to 30 minutes.  It's currently at the end of 4.6 and we're trying to see if the
	timeout := 40 * time.Minute

	// Wait longer for baremetal, due to length of time it takes to boot
	if installConfig, err := assetstore.Load(&installconfig.InstallConfig{}); err == nil && installConfig != nil {
		if installConfig.(*installconfig.InstallConfig).Config.Platform.Name() == baremetal.Name {
			timeout = 60 * time.Minute
		}
	}

	untilTime := time.Now().Add(timeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for the cluster at %s to initialize...",
		timeout, untilTime.Format(time.Kitchen), timezone, config.Host)
	cc, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create a config client")
	}
	clusterVersionContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	failing := configv1.ClusterStatusConditionType("Failing")
	timer.StartTimer("Cluster Operators Available")
	var lastError string
	_, err = clientwatch.UntilWithSync(
		clusterVersionContext,
		cache.NewListWatchFromClient(cc.ConfigV1().RESTClient(), "clusterversions", "", fields.OneTermEqualSelector("metadata.name", "version")),
		&configv1.ClusterVersion{},
		nil,
		func(event watch.Event) (bool, error) {
			switch event.Type {
			case watch.Added, watch.Modified:
				cv, ok := event.Object.(*configv1.ClusterVersion)
				if !ok {
					logrus.Warnf("Expected a ClusterVersion object but got a %q object instead", event.Object.GetObjectKind().GroupVersionKind())
					return false, nil
				}
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorAvailable) &&
					cov1helpers.IsStatusConditionFalse(cv.Status.Conditions, failing) &&
					cov1helpers.IsStatusConditionFalse(cv.Status.Conditions, configv1.OperatorProgressing) {
					timer.StopTimer("Cluster Operators Available")
					return true, nil
				}
				if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, failing) {
					lastError = cov1helpers.FindStatusCondition(cv.Status.Conditions, failing).Message
				} else if cov1helpers.IsStatusConditionTrue(cv.Status.Conditions, configv1.OperatorProgressing) {
					lastError = cov1helpers.FindStatusCondition(cv.Status.Conditions, configv1.OperatorProgressing).Message
				}
				logrus.Debugf("Still waiting for the cluster to initialize: %s", lastError)
				return false, nil
			}
			logrus.Debug("Still waiting for the cluster to initialize...")
			return false, nil
		},
	)

	if err == nil {
		logrus.Debug("Cluster is initialized")
		return nil
	}

	if lastError != "" {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return errors.Errorf("failed to initialize the cluster: %s", lastError)
		}

		return errors.Wrapf(err, "failed to initialize the cluster: %s", lastError)
	}

	return errors.Wrap(err, "failed to initialize the cluster")
}

// waitForStableOperators ensures that each cluster operator is "stable", i.e. the
// operator has not been in a progressing state for at least a certain duration,
// 30 seconds by default. Returns an error if any operator does meet this threshold
// after a deadline, 30 minutes by default.
func waitForStableOperators(ctx context.Context, config *rest.Config) error {
	timer.StartTimer("Cluster Operators Stable")

	stabilityCheckDuration := 30 * time.Minute
	stabilityContext, cancel := context.WithTimeout(ctx, stabilityCheckDuration)
	defer cancel()

	untilTime := time.Now().Add(stabilityCheckDuration)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) to ensure each cluster operator has finished progressing...",
		stabilityCheckDuration, untilTime.Format(time.Kitchen), timezone)

	cc, err := configclient.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create a config client")
	}
	configInformers := configinformers.NewSharedInformerFactory(cc, 0)
	clusterOperatorInformer := configInformers.Config().V1().ClusterOperators().Informer()
	clusterOperatorLister := configInformers.Config().V1().ClusterOperators().Lister()
	configInformers.Start(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(), clusterOperatorInformer.HasSynced) {
		return fmt.Errorf("informers never started")
	}

	waitErr := wait.PollUntilContextCancel(stabilityContext, 1*time.Second, true, waitForAllClusterOperators(clusterOperatorLister))
	if waitErr != nil {
		logrus.Errorf("Error checking cluster operator Progressing status: %q", waitErr)
		stableOperators, unstableOperators, err := currentOperatorStability(clusterOperatorLister)
		if err != nil {
			logrus.Errorf("Error checking final cluster operator Progressing status: %q", err)
		}
		logrus.Debugf("These cluster operators were stable: [%s]", strings.Join(sets.List(stableOperators), ", "))
		logrus.Errorf("These cluster operators were not stable: [%s]", strings.Join(sets.List(unstableOperators), ", "))

		logrus.Exit(ExitCodeOperatorStabilityFailed)
	}

	timer.StopTimer("Cluster Operators Stable")

	logrus.Info("All cluster operators have completed progressing")

	return nil
}

// getConsole returns the console URL from the route 'console' in namespace openshift-console.
func getConsole(ctx context.Context, config *rest.Config) (string, error) {
	url := ""
	// Need to keep these updated if they change
	consoleNamespace := "openshift-console"
	consoleRouteName := "console"
	rc, err := routeclient.NewForConfig(config)
	if err != nil {
		return "", errors.Wrap(err, "creating a route client")
	}

	consoleRouteTimeout := 2 * time.Minute
	logrus.Infof("Checking to see if there is a route at %s/%s...", consoleNamespace, consoleRouteName)
	consoleRouteContext, cancel := context.WithTimeout(ctx, consoleRouteTimeout)
	defer cancel()
	// Poll quickly but only log when the response
	// when we've seen 15 of the same errors or output of
	// no route in a row (to show we're still alive).
	logDownsample := 15
	silenceRemaining := logDownsample
	timer.StartTimer("Console")
	wait.Until(func() {
		route, err := rc.RouteV1().Routes(consoleNamespace).Get(ctx, consoleRouteName, metav1.GetOptions{})
		if err == nil {
			logrus.Debugf("Route found in openshift-console namespace: %s", consoleRouteName)
			if uri, _, err2 := routeapihelpers.IngressURI(route, ""); err2 == nil {
				url = uri.String()
				logrus.Debug("OpenShift console route is admitted")
				cancel()
			} else {
				err = err2
			}
		} else if apierrors.IsNotFound(err) {
			logrus.Debug("OpenShift console route does not exist")
			cancel()
		}

		if err != nil {
			silenceRemaining--
			if silenceRemaining == 0 {
				logrus.Debugf("Still waiting for the console route: %v", err)
				silenceRemaining = logDownsample
			}
		}
	}, 2*time.Second, consoleRouteContext.Done())
	err = consoleRouteContext.Err()
	if err != nil && !errors.Is(err, context.Canceled) {
		return url, errors.Wrap(err, "waiting for openshift-console URL")
	}
	if url == "" {
		return url, errors.New("could not get openshift-console URL")
	}
	timer.StopTimer("Console")
	return url, nil
}

// logComplete prints info upon completion.
func logComplete(directory, consoleURL string) error {
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	kubeconfig := filepath.Join(absDir, "auth", "kubeconfig")
	pwFile := filepath.Join(absDir, "auth", "kubeadmin-password")
	pw, err := os.ReadFile(pwFile)
	if err != nil {
		return err
	}
	logrus.Info("Install complete!")
	logrus.Infof("To access the cluster as the system:admin user when using 'oc', run\n    export KUBECONFIG=%s", kubeconfig)
	if consoleURL != "" {
		logrus.Infof("Access the OpenShift web-console here: %s", consoleURL)
		if SkipPasswordPrintFlag {
			logrus.Infof("Credentials omitted, if necessary verify the %s file", pwFile)
		} else {
			logrus.Infof("Login to the console with user: %q, and password: %q", "kubeadmin", pw)
		}
	}
	return nil
}

// addRouterCAToClusterCA adds router CA to cluster CA in kubeconfig.
func addRouterCAToClusterCA(ctx context.Context, config *rest.Config, directory string) (err error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "creating a Kubernetes client")
	}

	// Configmap may not exist. log and accept not-found errors with configmap.
	caConfigMap, err := client.CoreV1().ConfigMaps("openshift-config-managed").Get(ctx, "default-ingress-cert", metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "fetching default-ingress-cert configmap from openshift-config-managed namespace")
	}

	routerCrtBytes := []byte(caConfigMap.Data["ca-bundle.crt"])
	kubeconfig := filepath.Join(directory, "auth", "kubeconfig")
	kconfig, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return errors.Wrap(err, "loading kubeconfig")
	}

	if kconfig == nil || len(kconfig.Clusters) == 0 {
		return errors.New("kubeconfig is missing expected data")
	}

	for _, c := range kconfig.Clusters {
		clusterCABytes := c.CertificateAuthorityData
		if len(clusterCABytes) == 0 {
			return errors.New("kubeconfig CertificateAuthorityData not found")
		}
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(clusterCABytes) {
			return errors.New("cluster CA found in kubeconfig not valid PEM format")
		}
		if !certPool.AppendCertsFromPEM(routerCrtBytes) {
			return errors.New("ca-bundle.crt from default-ingress-cert configmap not valid PEM format")
		}

		routerCrtBytes := append(routerCrtBytes, clusterCABytes...)
		c.CertificateAuthorityData = routerCrtBytes
	}
	if err := clientcmd.WriteToFile(*kconfig, kubeconfig); err != nil {
		return errors.Wrap(err, "writing kubeconfig")
	}
	return nil
}

// CheckIfAgentCommand logs a warning if an agent configuration was detected
// and the current command is not an agent wait-for command.
func CheckIfAgentCommand(assetStore asset.Store) {
	if agentConfig, err := assetStore.Load(&agentconfig.AgentConfig{}); err == nil && agentConfig != nil {
		logrus.Warning("An agent configuration was detected but this command is not the agent wait-for command")
	}
}

func waitForAllClusterOperators(clusterOperatorLister configlisters.ClusterOperatorLister) func(ctx context.Context) (bool, error) {
	previouslyStableOperators := sets.Set[string]{}

	return func(ctx context.Context) (bool, error) {
		stableOperators, unstableOperators, err := currentOperatorStability(clusterOperatorLister)
		if err != nil {
			return false, err
		}
		if newlyStableOperators := stableOperators.Difference(previouslyStableOperators); len(newlyStableOperators) > 0 {
			for _, name := range sets.List(newlyStableOperators) {
				logrus.Debugf("Cluster Operator %s is stable", name)
			}
		}
		if newlyUnstableOperators := previouslyStableOperators.Difference(stableOperators); len(newlyUnstableOperators) > 0 {
			for _, name := range sets.List(newlyUnstableOperators) {
				logrus.Debugf("Cluster Operator %s became unstable", name)
			}
		}
		previouslyStableOperators = stableOperators

		if len(unstableOperators) == 0 {
			return true, nil
		}

		return false, nil
	}
}

func currentOperatorStability(clusterOperatorLister configlisters.ClusterOperatorLister) (sets.Set[string], sets.Set[string], error) {
	clusterOperators, err := clusterOperatorLister.List(labels.Everything())
	if err != nil {
		return nil, nil, err // lister should never fail
	}

	stableOperators := sets.Set[string]{}
	unstableOperators := sets.Set[string]{}
	for _, clusterOperator := range clusterOperators {
		name := clusterOperator.Name
		progressing := cov1helpers.FindStatusCondition(clusterOperator.Status.Conditions, configv1.OperatorProgressing)
		if progressing == nil {
			logrus.Debugf("Cluster Operator %s progressing == nil", name)
			unstableOperators.Insert(name)
			continue
		}
		if meetsStabilityThreshold(progressing) {
			stableOperators.Insert(name)
		} else {
			logrus.Debugf("Cluster Operator %s is Progressing=%s LastTransitionTime=%v DurationSinceTransition=%.fs Reason=%s Message=%s", name, progressing.Status, progressing.LastTransitionTime.Time, time.Since(progressing.LastTransitionTime.Time).Seconds(), progressing.Reason, progressing.Message)
			unstableOperators.Insert(name)
		}
	}

	return stableOperators, unstableOperators, nil
}

func meetsStabilityThreshold(progressing *configv1.ClusterOperatorStatusCondition) bool {
	return progressing.Status == configv1.ConditionFalse && time.Since(progressing.LastTransitionTime.Time).Seconds() > coStabilityThreshold
}

// LogTroubleshootingLink displays a link for additional troubleshooting help when installation is not successful.
func LogTroubleshootingLink() {
	logrus.Error(`Cluster initialization failed because one or more operators are not functioning properly.
The cluster should be accessible for troubleshooting as detailed in the documentation linked below,
https://docs.openshift.com/container-platform/latest/support/troubleshooting/troubleshooting-installations.html
The 'wait-for install-complete' subcommand can then be used to continue the installation`)
}
