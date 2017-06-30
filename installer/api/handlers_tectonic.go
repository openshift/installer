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

	"github.com/coreos/tectonic-installer/installer/pkg/tectonic"
	"github.com/coreos/tectonic-installer/installer/pkg/terraform"
	"github.com/dghubble/sessions"
)

// Status represents the individual state of a single service
type Status struct {
	Success bool `json:"success"`
	Failed  bool `json:"failed"`
}

// Services represents whether several Tectonic components have yet booted (or failed)
type Services struct {
	Kubernetes Status                 `json:"kubernetes"`
	HasEtcd    bool                   `json:"hasetcd"`
	Etcd       Status                 `json:"etcd"`
	Identity   Status                 `json:"identity"`
	Ingress    Status                 `json:"ingress"`
	Console    tectonic.ServiceStatus `json:"console"`
	Tectonic   Status                 `json:"tectonic"` // All other Tectonic services
}

func isEtcdHosted(ex terraform.Executor) (bool, error) {
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

func newK8sClient(ex terraform.Executor) (*kubernetes.Clientset, error) {
	kCfgPath := filepath.Join(ex.WorkingDirectory(), "generated/auth/kubeconfig")

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

func getStatus(p v1.PodPhase) Status {
	if p == v1.PodRunning {
		return Status{Failed: false, Success: true}
	} else if p == v1.PodFailed {
		return Status{Failed: true, Success: false}
	}
	return Status{Failed: false, Success: false}
}

type Input struct {
	TectonicDomain string `json:"tectonicDomain"`
}

func tectonicStatus(ex *terraform.Executor, input Input) (Services, error) {
	services := Services{
		Kubernetes: Status{Success: false, Failed: false},
		Etcd:       Status{Success: false, Failed: false},
		Identity:   Status{Success: false, Failed: false},
		Ingress:    Status{Success: false, Failed: false},
		Console:    tectonic.ConsoleHealth(nil, input.TectonicDomain),
		Tectonic:   Status{Success: false, Failed: false},
	}

	client, err := newK8sClient(*ex)
	if err != nil {
		// This is probably because the kubeconfig hasn't been generated yet; assume services haven't started
		// TODO: Better error handling to make sure it's not a different issue
		return services, nil
	}

	pods, err := client.CoreV1().Pods("").List(v1.ListOptions{})
	if err != nil {
		// APIServer probably hasn't started yet; assume services haven't started
		// TODO: Better error handling to make sure it's not a different issue
		return services, nil
	}

	etcd, err := isEtcdHosted(*ex)
	if err != nil {
		return services, newInternalServerError("Error reading TF Vars: %s", err)
	}

	services.HasEtcd = etcd

	services.Kubernetes.Success = true
	services.Tectonic.Success = true

	var kubeServs []v1.Pod
	var tectServs []v1.Pod

	for _, p := range pods.Items {
		name := p.ObjectMeta.Name
		switch {
		case strings.Contains(name, "tectonic-identity"):
			services.Identity = getStatus(p.Status.Phase)
		case strings.Contains(name, "tectonic-ingress-controller"):
			services.Ingress = getStatus(p.Status.Phase)
		case etcd && strings.Contains(name, "kube-etcd"):
			services.Etcd = getStatus(p.Status.Phase)
		case strings.Contains(name, "kube-"):
			kubeServs = append(kubeServs, p)
		default:
			tectServs = append(tectServs, p)
		}
	}

	if len(kubeServs) == 0 {
		services.Kubernetes = Status{Success: false, Failed: false}
	} else {
		for _, p := range kubeServs {
			if p.Status.Phase == v1.PodFailed || services.Kubernetes.Failed {
				services.Kubernetes = Status{Failed: true, Success: false}
			} else if p.Status.Phase == v1.PodRunning && services.Kubernetes.Success {
				services.Kubernetes = Status{Failed: false, Success: true}
			} else {
				services.Kubernetes = Status{Failed: false, Success: false}
			}
		}
	}

	if len(tectServs) == 0 {
		services.Tectonic = Status{Success: false, Failed: false}
	} else {
		for _, p := range tectServs {
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

type TerraformStatus struct {
	Status string `json:"status"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
	Action string `json:"action"`
}

func terraformStatus(session *sessions.Session, ex *terraform.Executor, exID int) (TerraformStatus, error) {
	// Retrieve Terraform's status and output.
	status, err := ex.Status(exID)
	if status == terraform.ExecutionStatusUnknown {
		return TerraformStatus{}, newBadRequestError("Could not retrieve TerraForm execution's status: %s", err)
	}
	output, _ := ex.Output(exID)
	outputBytes, _ := ioutil.ReadAll(output)

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

type Response struct {
	Terraform TerraformStatus `json:"terraform"`
	Tectonic  Services        `json:"tectonic"`
}

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
			return newInternalServerError("Could not retrieve Tectonic status: %s", err)
		}
		response.Tectonic = tec
	}

	return writeJSONResponse(w, req, http.StatusOK, response)
}
