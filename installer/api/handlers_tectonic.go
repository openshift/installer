package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"

	"strings"

	"errors"

	"os"

	"github.com/coreos/tectonic-installer/installer/pkg/containerlinux"
	"github.com/coreos/tectonic-installer/installer/pkg/terraform"
	"github.com/dghubble/sessions"
)

// Status represents the individual state of a single Kubernetes service
type Status struct {
	Success   bool   `json:"success"`
	Failed    bool   `json:"failed"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Message   string `json:"message"`
}

// Services represents whether several Tectonic components have yet booted (or failed)
type Services struct {
	Kubernetes       Status `json:"kubernetes"`
	IsEtcdSelfHosted bool   `json:"isEtcdSelfHosted"`
	Etcd             Status `json:"etcd"`
	Identity         Status `json:"identity"`
	Ingress          Status `json:"ingress"`
	Console          Status `json:"console"`
	TectonicSystem   Status `json:"tectonicSystem"` // All other services in tectonic-system namespace
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
	cfgPath := filepath.Join(ex.WorkingDirectory(), kubeconfigLocation)

	if _, err := os.Stat(cfgPath); err != nil {
		return nil, err
	}

	rules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: cfgPath}
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

// setStatus is a simple helper function to convert a Kubernetes pod phase to
// a Status object as defined above
func statusFromPods(pods []v1.Pod) Status {
	status := Status{Failed: false, Success: true}

	if len(pods) == 0 {
		status.Success = false
		return status
	}

	for _, p := range pods {
		if p.Status.Phase == v1.PodFailed {
			status.Failed = true
			status.Success = false
			status.Namespace = p.ObjectMeta.Namespace
			status.Pod = p.ObjectMeta.Name
			return status
		}

		if p.Status.Phase != v1.PodRunning {
			status.Success = false
		}
	}

	return status
}

// Input represents the request body passed to the status endpoint.
type Input struct {
	TectonicDomain string `json:"tectonicDomain"`
}

// ConsoleHealth returns the Status of the Tectonic Console.
func consoleHealth(endpoint string) Status {
	// TODO: Determine Tectonic Console failure state
	status := Status{
		Success: false,
		Failed:  false,
		Message: "",
	}

	// defaultStatusClient respects proxies, sets reasonable timeouts, and allows
	// checking the status of services running with self-signed certificates.
	client := http.Client{
		Timeout: time.Duration(10 * time.Second),
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	resp, err := client.Get(fmt.Sprintf("https://%s/health", endpoint))
	if err != nil {
		status.Message = err.Error()
		return status
	}
	defer resp.Body.Close()

	type health struct {
		Status string `json:"status"`
	}
	h := new(health)
	err = json.NewDecoder(resp.Body).Decode(h)
	if err != nil {
		status.Message = err.Error()
		return status
	}

	if h.Status == "ok" {
		status.Success = true
	}
	return status
}

// tectonicStatus returns the current status (starting, running, or failed) of
// the tectonic cluster services (etcd, Tectonic Identity, Tectonic Console), etc.
func tectonicStatus(ex *terraform.Executor, input Input) (Services, error) {
	services := Services{
		Kubernetes:     Status{Success: false, Failed: false},
		Etcd:           Status{Success: false, Failed: false},
		Identity:       Status{Success: false, Failed: false},
		Ingress:        Status{Success: false, Failed: false},
		Console:        Status{Success: false, Failed: false},
		TectonicSystem: Status{Success: false, Failed: false},
	}

	client, err := newK8sClient(*ex)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			// This is because the kubeconfig hasn't been generated yet; assume services haven't started
			services.Kubernetes.Message = err.Error()
			return services, nil
		}
		services.Kubernetes.Message = err.Error()
		return services, err
	}

	pods, err := client.CoreV1().Pods("").List(v1.ListOptions{})
	if err != nil {
		// APIServer probably hasn't started yet; assume services haven't started
		return services, errors.New("API server has not started")
	}

	etcd, err := isEtcdSelfHosted(*ex)
	if err != nil {
		services.Etcd.Message = err.Error()
		return services, newInternalServerError("Failed to determine if etcd is self-hosted: %s", err)
	}

	services.IsEtcdSelfHosted = etcd

	var kubePods []v1.Pod
	var tectPods []v1.Pod
	var identiPods []v1.Pod
	var ingressPods []v1.Pod
	var etcdPods []v1.Pod

	for _, p := range pods.Items {
		name := p.ObjectMeta.Name
		namespace := p.ObjectMeta.Namespace
		switch {
		case strings.Contains(name, "tectonic-identity"):
			identiPods = append(identiPods, p)
		case strings.Contains(name, "tectonic-ingress-controller"):
			ingressPods = append(ingressPods, p)
		case etcd && strings.Contains(name, "kube-etcd"):
			etcdPods = append(etcdPods, p)
		case namespace == "kube-system":
			kubePods = append(kubePods, p)
		case namespace == "tectonic-system":
			tectPods = append(tectPods, p)
		}
	}

	services.Kubernetes = statusFromPods(kubePods)
	services.TectonicSystem = statusFromPods(tectPods)
	services.Etcd = statusFromPods(etcdPods)
	services.Identity = statusFromPods(identiPods)
	services.Ingress = statusFromPods(ingressPods)

	services.Console = consoleHealth(input.TectonicDomain)

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
	status, statusErr := ex.Status(exID)
	if status == terraform.ExecutionStatusUnknown {
		return TerraformStatus{}, newBadRequestError("Could not retrieve TerraForm execution's status: %s", statusErr)
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
	if statusErr != nil {
		response.Error = statusErr.Error()
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
			if err.Error() != "API server has not started" {
				return newInternalServerError("Could not retrieve Tectonic status: %s", err)
			}
		}
		response.Tectonic = tec
	}

	return writeJSONResponse(w, req, http.StatusOK, response)
}

// Download kubeconfig file, or return a 404 if it doesn't exist yet.
func tectonicKubeconfigHandler(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	// Restore the execution environment from the session.
	_, ex, _, err := restoreExecutionFromSession(req, ctx.Sessions, nil)
	if err != nil {
		return err
	}

	cfg, err := os.Open(filepath.Join(ex.WorkingDirectory(), kubeconfigLocation))
	if err != nil {
		http.Error(w, "Could not retrieve kubeconfig", http.StatusNotFound)
		return nil
	}
	defer cfg.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=kubeconfig")
	w.Header().Set("Content-Type", "text/plain")
	io.Copy(w, cfg)
	return nil
}

// tectonicFactsHandler gets a list of available Container Linux AMIs as well
// as the Tectonic license and pull secret if they exist.
func tectonicFactsHandler(w http.ResponseWriter, req *http.Request, ctx *Context) error {
	var amis []containerlinux.AMI
	var err error

	// We only need the AMIs list if AWS is enabled
	for _, platform := range ctx.Config.Platforms {
		if platform == "aws-tf" {
			amis, err = containerlinux.ListAMIImages(containerLinuxListTimeout)
			if err != nil {
				return newInternalServerError("Failed to query available images: %s", err)
			}
			break
		}
	}

	ex, err := os.Executable()
	if err != nil {
		return newInternalServerError("Could not retrieve Tectonic facts")
	}

	license, err := ioutil.ReadFile(filepath.Join(filepath.Dir(ex), "license.txt"))
	if err != nil {
		license = nil
	}

	pullSecret, err := ioutil.ReadFile(filepath.Join(filepath.Dir(ex), "pull_secret.json"))
	if err != nil {
		pullSecret = nil
	}

	type response struct {
		AMIs       []containerlinux.AMI `json:"amis"`
		License    string               `json:"license"`
		PullSecret string               `json:"pullSecret"`
	}
	return writeJSONResponse(w, req, http.StatusOK, response{amis, string(license), string(pullSecret)})
}
