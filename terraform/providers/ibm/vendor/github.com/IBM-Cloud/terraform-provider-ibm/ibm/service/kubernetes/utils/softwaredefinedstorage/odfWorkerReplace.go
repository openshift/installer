package softwaredefinedstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	v2 "github.com/IBM-Cloud/bluemix-go/api/container/containerv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	templatev1 "github.com/openshift/api/template/v1"
	templatev1client "github.com/openshift/client-go/template/clientset/versioned/typed/template/v1"
	cephv1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	corev1 "k8s.io/api/core/v1"
	policy "k8s.io/api/policy/v1beta1"
	apierror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var deploymentList []string

// ODF Struct Defined
type odf struct{}

func NewSdsOdf() Sds {
	return &odf{}
}

// Steps before Worker Replace for ODF
func (o odf) PreWorkerReplace(worker v2.Worker) error {
	log.Println("Inside preWorkerReplace for ODF")
	workerName := worker.NetworkInterfaces[0].IpAddress
	log.Println("This is the Worker to be replaced", workerName)
	// Scale Deployments to 0
	for label, id := range AppLabelId {
		err := scalePods(workerName, app, label, id, 0)
		if err != nil {
			return err
		}
	}
	// Check if Replica Set is 0
	for _, v := range deploymentList {
		_, err := waitForOdfDeploymentStatus(0, v)
		if err != nil {
			return err
		}
	}

	log.Println("Deployments have been successfully scaled!")
	node := workerName

	// Cordon the node
	_, err := waitForNodeCordonStatus(node)
	if err != nil {
		return err
	}

	//Drain the node
	pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{FieldSelector: "spec.nodeName=" + workerName})
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting Pods from worker node %s - %s", workerName, err)
	}
	// Evict the pods from the given node
	for _, pod := range pods.Items {
		evictPod(pod.Name, pod.Namespace)
		log.Printf("Pod %s has been evicted\n", pod.Name)
	}
	log.Printf("Node %s has been drained\n", node)
	return nil

}

func evictPod(name, namespace string) error {

	return clientSet.PolicyV1beta1().Evictions(namespace).Evict(context.TODO(), &policy.Eviction{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace},
	})
}

// Function to Scale the deployments based on the pods present in the current worker node
func scalePods(workerName string, appLabel string, appLabelValue string, idLabel string, replicas int32) error {
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{appLabel: appLabelValue}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
		FieldSelector: "spec.nodeName=" + workerName}
	pods, err := clientSet.CoreV1().Pods(odfNamespace).List(context.TODO(), listOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting Pods in the openshift-storage namespace: %s", err)
	}
	if len(pods.Items) == 0 {
		return nil
	}
	for _, pod := range pods.Items {
		log.Printf("Pod name: %v\n", pod.Name)
		deploymentName := pod.Labels[appLabel] + "-" + pod.Labels[idLabel]
		deploymentsClient := clientSet.AppsV1().Deployments(odfNamespace)
		deploymentList = append(deploymentList, deploymentName)
		result, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("[ERROR] Failed to get latest version of deployment: %v - %v", deploymentName, err)
		}

		result.Spec.Replicas = int32Ptr(replicas) // reduce replica count
		scaleDeployment, err := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating scale object in deployment: %v - %v", deploymentName, err)
		}

		log.Printf("Successfully scaled deployment %s to %d replicas", deploymentName, *scaleDeployment.Spec.Replicas)
	}
	return nil
}

// Steps after worker replace has been complete
func (o odf) PostWorkerReplace(worker v2.Worker) error {
	log.Println("In Post Worker Replace")
	deploymentsClient := clientSet.AppsV1().Deployments(odfNamespace)
	log.Println("Scaling Up Deployments")
	// Scale Up Deployments that have been Brought down in pre worker replace
	for _, v := range deploymentList {
		if !strings.Contains(v, crashcollectorLabel) {
			result, err := deploymentsClient.Get(context.TODO(), v, metav1.GetOptions{})
			if err != nil {
				return fmt.Errorf("[ERROR] Failed to get latest version of deployment: %v - %v", v, err)
			}
			result.Spec.Replicas = int32Ptr(1) // increase replica count
			scaleDeployment, err := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("[Error] Error updating scale object in deployment: %v - %v", v, err)
			}
			log.Printf("Successfully scaled deployment %s to %d replicas", v, *scaleDeployment.Spec.Replicas)
		}
	}

	log.Println("Checking if Number of Replicas is 1")
	// Check if the deployment replicas are 1
	for _, v := range deploymentList {
		if !strings.Contains(v, crashcollectorLabel) {
			_, err := waitForOdfDeploymentStatus(1, v)
			if err != nil {
				return err
			}
		}
	}

	log.Println("Deployments have been successfully scaled")
	deploymentList = nil
	workerName := worker.NetworkInterfaces[0].IpAddress
	log.Println("This is the New Worker Name", workerName)

	// Check the status of the OSD pods, if failed then run job
	labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{app: osdLabel}}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
		FieldSelector: "spec.nodeName=" + workerName}

	pods, err := clientSet.CoreV1().Pods(odfNamespace).List(context.TODO(), listOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting Pods in the openshift-storage namespace %v", err)
	}

	osdID := ""
	for _, pod := range pods.Items {
		log.Printf("Pod Name %s is %s/n", pod.Name, pod.Status.Phase)
		if !(pod.Status.Phase == "Running" || pod.Status.Phase == "Pending") {
			osdID = osdID + pod.Labels["ceph-osd-id"] + ","
		}
	}
	osdID = strings.TrimRight(osdID, ",")
	if len(osdID) > 0 {
		err := executeTemplate(osdID)
		if err != nil {
			return err
		}
	}

	// Check Ceph Status if HEALTH OK then return success!
	log.Println("Fetching Ceph Cluster Status")
	_, err = waitForCephClusterStatus()
	if err != nil {
		return err
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }

// To execute a given template using the openshift go client
func executeTemplate(osdID string) error {
	// Create an OpenShift template/v1 client.
	templateclient, err := templatev1client.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating templateClient %s", err)
	}

	// Get the "ocs-osd-removal" Template from the "openshift-storage" Namespace.
	template, err := templateclient.Templates(odfNamespace).Get(context.Background(),
		"ocs-osd-removal", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("[ERROR] Error getting ocs-osd-removal template %s", err)
	}

	// INSTANTIATE THE TEMPLATE.
	// To set Template parameters, create a Secret holding overridden parameters
	// and their values.
	secret, err := clientSet.CoreV1().Secrets(odfNamespace).Create(context.Background(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "parameters",
		},
		StringData: map[string]string{
			"FAILED_OSD_IDS": osdID,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating secret for ocs-osd-removal template %s", err)
	}

	// Create a TemplateInstance object, linking the Template and a reference to
	// the Secret object created above.
	templateInstance, err := templateclient.TemplateInstances(odfNamespace).Create(context.Background(),
		&templatev1.TemplateInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "templateinstance",
			},
			Spec: templatev1.TemplateInstanceSpec{
				Template: *template,
				Secret: &corev1.LocalObjectReference{
					Name: secret.Name,
				},
			},
		}, metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("[ERROR] Error deploying ocs-osd-removal template %s", err)
	}

	_, err = waitForTemplateInstanceStatus(templateInstance)
	if err != nil {
		return err
	}

	// DELETE THE INSTANTIATED TEMPLATE.
	// We use the foreground propagation policy to ensure that the garbage
	// collector removes all instantiated objects before the TemplateInstance
	// itself disappears.
	foreground := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{PropagationPolicy: &foreground}
	err = templateclient.TemplateInstances(odfNamespace).Delete(context.TODO(), templateInstance.Name,
		deleteOptions)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting TemplateInstance resource %s", err)
	}

	// Finally delete the "parameters" Secret.
	err = clientSet.CoreV1().Secrets(odfNamespace).Delete(context.TODO(), secret.Name,
		metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting parameter Secret resource %s", err)
	}

	return nil
}

func waitForOdfDeploymentStatus(replicas int32, deploymentName string) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        odfDeploymentRefreshFunc(replicas, deploymentName),
		Timeout:        time.Duration(sdsTimeout),
		Delay:          5 * time.Second,
		MinTimeout:     5 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func odfDeploymentRefreshFunc(replicas int32, deploymentName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Println("Checking Deployment Status....")
		deploymentsClient := clientSet.AppsV1().Deployments(odfNamespace)
		result, err := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return nil, "NotReady", fmt.Errorf("[ERROR] Failed to get latest version of deployment: %v - %v", deploymentName, err)
		}
		log.Printf("Deployment Name:%s Replica Number:%d\n", deploymentName, result.Status.ReadyReplicas)
		if result.Status.ReadyReplicas == replicas {
			log.Println("Deployment has been successfully scaled")
			return true, "Ready", nil
		}
		return nil, "NotReady", nil
	}
}

func waitForCephClusterStatus() (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        cephClusterRefreshFunc(),
		Timeout:        time.Duration(sdsTimeout),
		Delay:          5 * time.Second,
		MinTimeout:     5 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func cephClusterRefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		scheme := runtime.NewScheme()
		utilruntime.Must(cephv1.AddToScheme(scheme))
		cntrlClient, err := client.New(restConfig, client.Options{Scheme: scheme})
		if err != nil {
			return nil, "NotReady", fmt.Errorf("[ERROR] Error getting config")
		}
		cephcluster := &cephv1.CephCluster{}
		var CephclusterGet []string
		CephclusterGet = append(CephclusterGet, "\n Cephcluster Details \n \n", "NAME  DATADIRHOSTPATH   MONCOUNT   AGE   PHASE   MESSAGE  HEALTH\n")
		err = cntrlClient.Get(context.Background(), types.NamespacedName{Name: "ocs-storagecluster-cephcluster", Namespace: odfNamespace}, cephcluster)
		if err != nil && !apierror.IsNotFound(err) {
			return nil, "NotReady", err
		} else {
			if apierror.IsNotFound(err) {
				return nil, "NotReady", fmt.Errorf("[ERROR] Cephcluster not found")
			} else {
				_, err := json.MarshalIndent(cephcluster, "", "  ")
				if err != nil {
					return nil, "NotReady", err
				}

			}
			CephclusterGet = append(CephclusterGet, "\n", cephcluster.Name, cephcluster.Spec.DataDirHostPath, strconv.Itoa(cephcluster.Spec.Mon.Count), time.Now().Sub(cephcluster.CreationTimestamp.Time).String(), string(cephcluster.Status.Phase), cephcluster.Status.Message, cephcluster.Status.CephStatus.Health, "\n")
			if err != nil {
				return nil, "NotReady", err
			}
		}
		log.Println("Ceph Cluster Status", CephclusterGet[9])
		if CephclusterGet[9] == "HEALTH_OK" {
			log.Println("Ceph Cluster is OK, Worker Replace Done!")
			return true, "Ready", nil
		}

		return nil, "NotReady", nil
	}
}

func waitForNodeCordonStatus(node string) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        nodeCordonRefreshFunc(node),
		Timeout:        time.Duration(sdsTimeout),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func nodeCordonRefreshFunc(node string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := clientSet.CoreV1().Nodes().Get(context.Background(), node, metav1.GetOptions{})
		if err != nil {
			return nil, "NotReady", fmt.Errorf("[ERROR] Node %s not found\n", node)
		}
		n.Spec.Unschedulable = true
		_, err = clientSet.CoreV1().Nodes().Update(context.Background(), n, metav1.UpdateOptions{})
		if err != nil {
			return nil, "NotReady", fmt.Errorf("[ERROR] Unable to update the node %s\n", node)
		}
		n, err = clientSet.CoreV1().Nodes().Get(context.Background(), node, metav1.GetOptions{})
		if err != nil {
			return nil, "NotReady", fmt.Errorf("[ERROR] Unable to cordon node %s\n", node)
		}
		if n.Spec.Unschedulable {
			log.Println("Node has been successfully Cordoned")
			return true, "Ready", nil
		}
		return nil, "NotReady", nil
	}
}

func waitForTemplateInstanceStatus(templateInstance *templatev1.TemplateInstance) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:        []string{"NotReady"},
		Target:         []string{"Ready"},
		Refresh:        templateInstanceRefreshFunc(templateInstance),
		Timeout:        time.Duration(sdsTimeout),
		Delay:          10 * time.Second,
		MinTimeout:     10 * time.Second,
		NotFoundChecks: 100,
	}
	return stateConf.WaitForState()
}

func templateInstanceRefreshFunc(templateInstance *templatev1.TemplateInstance) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		for _, cond := range templateInstance.Status.Conditions {
			// If the TemplateInstance contains a status condition
			// Ready == True, stop watching.
			if cond.Type == templatev1.TemplateInstanceReady &&
				cond.Status == corev1.ConditionTrue {
				return true, "Ready", nil
			}

			// If the TemplateInstance contains a status condition
			// InstantiateFailure == True, indicate failure.
			if cond.Type ==
				templatev1.TemplateInstanceInstantiateFailure &&
				cond.Status == corev1.ConditionTrue {
				return nil, "NotReady", fmt.Errorf("[ERROR] Error TemplateInstance instantiation failed")
			}
		}
		return nil, "NotReady", nil
	}
}
