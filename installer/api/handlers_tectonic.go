package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"

	"strings"

	"errors"

	"os"

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
	"github.com/coreos/tectonic-installer/installer/pkg/terraform"
	"github.com/dghubble/sessions"
	"github.com/hashicorp/terraform/helper/resource"
)

// Status represents the individual state of a single Kubernetes service
type Status struct {
	Success bool `json:"success"`
	Failed  bool `json:"failed"`
}

// Services represents whether several Tectonic components have yet booted (or failed)
type Services struct {
	Kubernetes       Status `json:"kubernetes"`
	IsEtcdSelfHosted bool   `json:"isEtcdSelfHosted"`
	Etcd             Status `json:"etcd"`
	Identity         Status `json:"identity"`
	Ingress          Status `json:"ingress"`
	Console          Status `json:"console"`
	Tectonic         Status `json:"tectonic"` // All other Tectonic services
	DnsName          string `json:"dnsName"`  // Stores the instance of the Tectonic Console
}

// isEtcdSelfHosted determines if the etcd service will be hosted on the cluster
// being terraformed; if it is, its status will be tracked like the other services
func isEtcdSelfHosted(ex terraform.Executor) (bool, error) {
	data, err := ex.LoadVars()
	if err != nil {
		return false, err
	}
	val, ok := data["tectonic_experimental"]
	if !ok {
		return false, nil
	}
	etcd, ok := val.(bool)
	if !ok {
		return false, errors.New("Could not parse experimental flag as bool")
	}
	return etcd, nil
}

// kubeconfigLocation is the location of the kubeconfig file, relative to the
// current working directory
const kubeconfigLocation = "generated/auth/kubeconfig"

// newK8sClient uses the generated kubeconfig file, if it exists, to create an
// API client that can be used for checking pod statuses
func newK8sClient(ex terraform.Executor) (*kubernetes.Clientset, error) {
	kCfgPath := filepath.Join(ex.WorkingDirectory(), kubeconfigLocation)

	if _, err := os.Stat(kCfgPath); err != nil {
		return nil, err
	}

	rules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kCfgPath}
	cfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})
	restConfig, err := cfg.ClientConfig()
	if err != nil {
		return nil, err
	}

	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// getStatus is a simple helper function to convert a Kubernetes pod phase to
// a Status object as defined above
func getStatus(p v1.PodPhase) Status {
	if p == v1.PodRunning {
		return Status{Failed: false, Success: true}
	} else if p == v1.PodFailed {
		return Status{Failed: true, Success: false}
	}
	return Status{Failed: false, Success: false}
}

// Input represents the request body passed to the status endpoint.
type Input struct {
	TectonicDomain string `json:"tectonicDomain"`
}

// tectonicStatus returns the current status (starting, running, or failed) of
// the tectonic cluster services (etcd, Tectonic Identity, Tectonic Console), etc.
func tectonicStatus(ex *terraform.Executor, input Input) (Services, error) {
	services := Services{
		Kubernetes: Status{Success: false, Failed: false},
		Etcd:       Status{Success: false, Failed: false},
		Identity:   Status{Success: false, Failed: false},
		Ingress:    Status{Success: false, Failed: false},
		Console:    Status{Success: false, Failed: false},
		Tectonic:   Status{Success: false, Failed: false},
	}

	client, err := newK8sClient(*ex)
	if err != nil {
		if _, ok := err.(os.PathError); ok {
			// This is because the kubeconfig hasn't been generated yet; assume services haven't started
			return services, nil
		} else {
			return nil, err
		}
	}

	pods, err := client.CoreV1().Pods("").List(v1.ListOptions{})
	if err != nil {
		// APIServer probably hasn't started yet; assume services haven't started
		return services, newNotFoundError("API server has not started")
	}

	etcd, err := isEtcdSelfHosted(*ex)
	if err != nil {
		return services, newInternalServerError("Failed to determine if etcd is self-hosted: %s", err)
	}

	console := tectonic.ConsoleHealth(nil, input.TectonicDomain)

	services.DnsName = console.Instance

	services.Console.Success = console.Ready
	// TODO: Determine Tectonic Console failure state

	services.IsEtcdSelfHosted = etcd

	services.Kubernetes.Success = true
	services.Tectonic.Success = true

	var kubePods []v1.Pod
	var tectPods []v1.Pod

	for _, p := range pods.Items {
		name := p.ObjectMeta.Name
		switch {
		case strings.Contains(name, "tectonic-identity"):
			services.Identity = getStatus(p.Status.Phase)
		case strings.Contains(name, "tectonic-ingress-controller"):
			services.Ingress = getStatus(p.Status.Phase)
		case etcd && strings.Contains(name, "kube-etcd"):
			services.Etcd = getStatus(p.Status.Phase)
		default:
			if p.ObjectMeta.Namespace == "kube-system" {
				kubePods = append(kubePods, p)
			} else if p.ObjectMeta.Namespace == "tectonic-system" {
				tectPods = append(tectPods, p)
			}
		}
	}

	if len(kubePods) == 0 {
		services.Kubernetes = Status{Success: false, Failed: false}
	} else {
		for _, p := range kubePods {
			if p.Status.Phase == v1.PodFailed || services.Kubernetes.Failed {
				services.Kubernetes = Status{Failed: true, Success: false}
			} else if p.Status.Phase == v1.PodRunning && services.Kubernetes.Success {
				services.Kubernetes = Status{Failed: false, Success: true}
			} else {
				services.Kubernetes = Status{Failed: false, Success: false}
			}
		}
	}

	if len(tectPods) == 0 {
		services.Tectonic = Status{Success: false, Failed: false}
	} else {
		for _, p := range tectPods {
			if p.Status.Phase == v1.PodFailed || services.Tectonic.Failed {
				services.Tectonic = Status{Failed: true, Success: false}
			} else if p.Status.Phase == v1.PodRunning && services.Tectonic.Success {
				services.Tectonic = Status{Failed: false, Success: true}
			} else {
				services.Tectonic = Status{Failed: false, Success: false}
			}
		}
	}

	return services, nil
}

// TerraformStatus represents the current status of the Terraform run.
type TerraformStatus struct {
	Status string `json:"status"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
	Action string `json:"action"`
}

// terraformStatus returns the current status of the terraform run, if it has begun
func terraformStatus(session *sessions.Session, ex *terraform.Executor, exID int) (TerraformStatus, error) {
	// Retrieve Terraform's status and output.
	status, err := ex.Status(exID)
	if status == terraform.ExecutionStatusUnknown {
		return TerraformStatus{}, newBadRequestError("Could not retrieve TerraForm execution's status: %s", err)
	}
	output, err := ex.Output(exID)
	if err != nil {
		return TerraformStatus{}, newInternalServerError("Could not retrieve Terraform output: %s", err)
	}
	outputBytes, err := ioutil.ReadAll(output)
	if err != nil {
		return TerraformStatus{}, newInternalServerError("Could not read Terraform output: %s", err)
	}

	// Return results.
	response := TerraformStatus{
		Status: string(status),
		Output: string(outputBytes),
	}
	action := session.Values["action"]
	if action != nil {
		response.Action = action.(string)
	}
	if err != nil {
		response.Error = err.Error()
	}

	return response, nil
}

// Response combines the terraform and tectonic statuses into one object
type Response struct {
	Terraform TerraformStatus `json:"terraform"`
	Tectonic  Services        `json:"tectonic"`
}

// tectonicStatusHandler retrieves the terraform and tectonic statuses, then
// combines them into one response and returns it. If the terraform apply has
// not yet begun, it will return a 404.
func tectonicStatusHandler(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	input := Input{}
	// Read the input from the request's body.
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return newBadRequestError("Could not unmarshal input: %s", err)
	}

	// Restore the execution environment from the session.
	session, ex, exID, errCtx := restoreExecutionFromSession(req, ctx.Sessions, nil)
	if errCtx != nil {
		// Error directly (rather than NewAppError, which logs) since the
		// frontend periodically calls this endpoint to advance screens
		http.Error(w, fmt.Sprintf("Could not find session data: %v", errCtx), http.StatusNotFound)
		return nil
	}

	response := Response{}

	tf, err := terraformStatus(session, ex, exID)
	if err != nil {
		return newInternalServerError("Could not retrieve Terraform status: %s", err)
	}
	response.Terraform = tf

	if tf.Action != "destroy" {
		tec, err := tectonicStatus(ex, input)
		if err != nil {
			if _, ok := err.(resource.NotFoundError); !ok {
				return newInternalServerError("Could not retrieve Tectonic status: %s", err)
			}
		}
		response.Tectonic = tec
	}

	return writeJSONResponse(w, req, http.StatusOK, response)
}
