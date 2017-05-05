package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"github.com/dghubble/sessions"

	"path/filepath"
	"time"

	"github.com/coreos/tectonic-installer/installer/server/ctxh"
	"github.com/coreos/tectonic-installer/installer/server/defaults"
	"github.com/coreos/tectonic-installer/installer/server/terraform"
	"github.com/kardianos/osext"
)

// TerraformApplyHandlerInput describes the input expected by the
// terraformApplyHandler HTTP Handler.
type TerraformApplyHandlerInput struct {
	Platform      string                 `json:"platform"`
	Credentials   terraform.Credentials  `json:"credentials"`
	AdminPassword []byte                 `json:"adminPassword"`
	Variables     map[string]interface{} `json:"variables"`
	License       string                 `json:"license"`
	PullSecret    string                 `json:"pullSecret"`
	DryRun        bool                   `json:"dryRun"`
	Retry         bool                   `json:"retry"`
}

func terraformApplyHandler(sessionProvider sessions.Store) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		if req.Method != "POST" {
			return ctxh.NewAppError(nil, "POST method required", http.StatusMethodNotAllowed)
		}

		// Read the input from the request's body.
		var input TerraformApplyHandlerInput
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			return ctxh.NewAppError(err, "could not unmarshal request data", http.StatusBadRequest)
		}
		defer req.Body.Close()

		var ex *terraform.Executor
		var errCtx *ctxh.AppError
		if input.Retry {
			// Restore the execution environment from the session.
			_, ex, _, errCtx = restoreExecutionFromSession(req, sessionProvider, &input.Credentials)
		} else {
			// Create a new TerraForm Executor with the TF variables.
			ex, errCtx = newExecutorFromApplyHandlerInput(&input)
		}
		if errCtx != nil {
			return errCtx
		}
		tfMainDir := fmt.Sprintf("%s/platforms/%s", ex.WorkingDirectory(), input.Platform)

		// Copy the TF Templates to the Executor's working directory.
		if err := terraform.RestoreSources(ex.WorkingDirectory()); err != nil {
			return ctxh.NewAppError(err, "could not write TerraForm templates", http.StatusInternalServerError)
		}

		// Execute TerraForm get and wait for it to finish.
		_, getDone, err := ex.Execute("get", "-no-color", tfMainDir)
		if err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to run TerraForm (get): %v", err.Error()), http.StatusInternalServerError)
		}
		<-getDone

		// Store both the path to the Executor and the ID of the execution so that
		// the status can be read later on.
		session := sessionProvider.New(installerSessionName)
		session.Values["terraform_path"] = ex.WorkingDirectory()

		var id int
		var action string
		if input.DryRun {
			id, _, err = ex.Execute("show", "-no-color", tfMainDir)
			action = "show"
		} else {
			id, _, err = ex.Execute("apply", "-input=false", "-no-color", tfMainDir)
			action = "apply"
		}
		if err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to run TerraForm (%v): %v", action, err.Error()), http.StatusInternalServerError)
		}
		session.Values["terraform_id"] = id
		session.Values["action"] = action

		if err := sessionProvider.Save(w, session); err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to save session: %v", err.Error()), http.StatusInternalServerError)
		}

		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}

func terraformStatusHandler(sessionProvider sessions.Store) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		// Restore the execution environment from the session.
		session, ex, exID, errCtx := restoreExecutionFromSession(req, sessionProvider, nil)
		if errCtx != nil {
			// Error directly (rather than NewAppError, which logs) since the
			// frontend periodically calls this endpoint to advance screens
			http.Error(w, fmt.Sprintf("Could not find session data: %v", errCtx), http.StatusNotFound)
			return nil
		}

		// Retrieve the status and output.
		status, err := ex.Status(exID)
		if status == terraform.ExecutionStatusUnknown {
			return ctxh.NewAppError(err, "could not retrieve TerraForm execution's status", http.StatusInternalServerError)
		}
		output, _ := ex.Output(exID)
		outputBytes, _ := ioutil.ReadAll(output)

		// Return results.
		result := struct {
			Status          string        `json:"status"`
			Output          string        `json:"output,omitempty"`
			Error           string        `json:"error,omitempty"`
			Action          string        `json:"action"`
			TectonicConsole ServiceStatus `json:"tectonicConsole"`
		}{
			Status: string(status),
			Output: string(outputBytes),
		}
		action := session.Values["action"]
		if action != nil {
			result.Action = action.(string)
		}
		if err != nil {
			result.Error = err.Error()
		}

		client := defaultStatusClient

		type Input struct {
			TectonicDomain string `json:"tectonicDomain"`
		}
		var input Input
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			return ctxh.NewAppError(err, "could not unmarshal request data", http.StatusBadRequest)
		}

		result.TectonicConsole = TectonicConsoleHealth(client, input.TectonicDomain)

		writeJSONData(w, result)
		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}

func terraformAssetsHandler(sessionProvider sessions.Store) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		// Restore the execution environment from the session.
		_, ex, _, err := restoreExecutionFromSession(req, sessionProvider, nil)
		if err != nil {
			return err
		}

		// Stream the assets as a ZIP.
		w.Header().Set("Content-Type", "application/zip")
		if err := ex.Zip(w, true); err != nil {
			return ctxh.NewAppError(err, "could not archive assets", http.StatusInternalServerError)
		}
		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}

// TerraformDestroyHandlerInput describes the input expected by the
// terraformDestroyHandler HTTP Handler.
type TerraformDestroyHandlerInput struct {
	Platform    string                `json:"platform"`
	Credentials terraform.Credentials `json:"credentials"`
}

func terraformDestroyHandler(sessionProvider sessions.Store) ctxh.ContextHandler {
	fn := func(ctx context.Context, w http.ResponseWriter, req *http.Request) *ctxh.AppError {
		// Read the input from the request's body.
		var input TerraformDestroyHandlerInput
		if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
			return ctxh.NewAppError(err, "could not unmarshal request data", http.StatusBadRequest)
		}
		defer req.Body.Close()

		// Restore the execution environment from the session.
		_, ex, _, errCtx := restoreExecutionFromSession(req, sessionProvider, &input.Credentials)
		if errCtx != nil {
			return errCtx
		}
		tfMainDir := fmt.Sprintf("%s/platforms/%s", ex.WorkingDirectory(), input.Platform)

		// Execute TerraForm apply in the background.
		id, _, err := ex.Execute("destroy", "-force", "-no-color", tfMainDir)
		if err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to run TerraForm (apply): %v", err.Error()), http.StatusInternalServerError)
		}

		// Store both the path to the Executor and the ID of the execution so that
		// the status can be read later on.
		session := sessionProvider.New(installerSessionName)
		session.Values["action"] = "destroy"
		session.Values["terraform_path"] = ex.WorkingDirectory()
		session.Values["terraform_id"] = id
		if err := sessionProvider.Save(w, session); err != nil {
			return ctxh.NewAppError(err, fmt.Sprintf("Failed to save session: %v", err.Error()), http.StatusInternalServerError)
		}

		return nil
	}
	return ctxh.ContextHandlerFuncWithError(fn)
}

// newExecutorFromApplyHandlerInput creates a new Executor based on the given
// TerraformApplyHandlerInput.
func newExecutorFromApplyHandlerInput(input *TerraformApplyHandlerInput) (*terraform.Executor, *ctxh.AppError) {
	// Construct the path where the Executor should run based on the the cluster
	// name and current's binary path.
	binaryPath, err := osext.ExecutableFolder()
	if err != nil {
		return nil, ctxh.NewAppError(err, fmt.Sprintf("Could not determine executable's folder: %v", err.Error()), http.StatusInternalServerError)
	}
	clusterName := input.Variables["tectonic_cluster_name"].(string)
	if len(clusterName) == 0 {
		return nil, ctxh.NewAppError(err, "Tectonic cluster name not provided", http.StatusBadRequest)
	}
	exPath := filepath.Join(binaryPath, "clusters", clusterName+time.Now().Format("_2006-01-02_15-04-05"))

	// Create a new Executor.
	ex, err := terraform.NewExecutor(exPath)
	if err != nil {
		return nil, ctxh.NewAppError(err, fmt.Sprintf("Could not create TerraForm executor: %v", err.Error()), http.StatusInternalServerError)
	}

	// Write the License and Pull Secret to disk, and wire these files in the
	// variables.
	if input.License == "" {
		return nil, ctxh.NewAppError(err, "Tectonic license not provided", http.StatusBadRequest)
	}
	ex.AddFile("license.txt", []byte(input.License))
	if input.PullSecret == "" {
		return nil, ctxh.NewAppError(err, "Tectonic pull secret not provided", http.StatusBadRequest)
	}
	ex.AddFile("pull_secret.json", []byte(input.PullSecret))
	input.Variables["tectonic_license_path"] = "./license.txt"
	input.Variables["tectonic_pull_secret_path"] = "./pull_secret.json"
	serviceCidr := input.Variables["tectonic_service_cidr"].(string)

	ip, ok := input.Variables["tectonic_kube_apiserver_service_ip"].(string)
	if !ok || len(ip) == 0 {
		input.Variables["tectonic_kube_apiserver_service_ip"], err = defaults.APIServiceIP(serviceCidr)
		if err != nil {
			return nil, ctxh.NewAppError(err, fmt.Sprintf("Error calculating service IP: %v", err.Error()), http.StatusInternalServerError)
		}
	}

	ip, ok = input.Variables["tectonic_kube_dns_service_ip"].(string)
	if !ok || len(ip) == 0 {
		input.Variables["tectonic_kube_dns_service_ip"], err = defaults.DNSServiceIP(serviceCidr)
		if err != nil {
			return nil, ctxh.NewAppError(err, fmt.Sprintf("Error calculating DNS IP: %v", err.Error()), http.StatusInternalServerError)
		}
	}

	ip, ok = input.Variables["tectonic_kube_etcd_service_ip"].(string)
	if !ok || len(ip) == 0 {
		input.Variables["tectonic_kube_etcd_service_ip"], err = defaults.EtcdServiceIP(serviceCidr)
		if err != nil {
			return nil, ctxh.NewAppError(err, fmt.Sprintf("Error calculating etcd service IP: %v", err.Error()), http.StatusInternalServerError)
		}
	}

	if len(input.AdminPassword) > 0 {
		passwordHash, err := bcrypt.GenerateFromPassword(input.AdminPassword, bcryptCost)
		if err != nil {
			return nil, ctxh.NewAppError(err, fmt.Sprintf("Error bcrypt()ing admin password: %v", err.Error()), http.StatusBadRequest)
		}
		input.Variables["tectonic_admin_password_hash"] = passwordHash
	}

	// Add variables and the required environment variables.
	if variables, err := json.Marshal(input.Variables); err == nil {
		ex.AddVariables(variables)
	} else {
		return nil, ctxh.NewAppError(err, fmt.Sprintf("Could not marshal TerraForm variables: %v", err.Error()), http.StatusBadRequest)
	}
	if err := ex.AddCredentials(&input.Credentials); err != nil {
		return nil, ctxh.NewAppError(err, fmt.Sprintf("Could not validate TerraForm credentials: %v", err.Error()), http.StatusBadRequest)
	}

	return ex, nil
}

// restoreExecutionFromSession tries to re-create an existing Executor based on
// the data available in session and the provided credentials.
func restoreExecutionFromSession(req *http.Request, sessionProvider sessions.Store, credentials *terraform.Credentials) (*sessions.Session, *terraform.Executor, int, *ctxh.AppError) {
	session, err := sessionProvider.Get(req, installerSessionName)
	if err != nil {
		return nil, nil, -1, ctxh.NewAppError(err, "Could not find session data. Run terraform apply first.", http.StatusNotFound)
	}
	executionPath, ok := session.Values["terraform_path"]
	if !ok {
		return nil, nil, -1, ctxh.NewAppError(err, "Could not find terraform_path in session. Run terraform apply first.", http.StatusNotFound)
	}
	executionID, ok := session.Values["terraform_id"]
	if !ok {
		return nil, nil, -1, ctxh.NewAppError(err, "Could not find terraform_id in session. Run terraform apply first.", http.StatusNotFound)
	}
	ex, err := terraform.NewExecutor(executionPath.(string))
	if err != nil {
		return nil, nil, -1, ctxh.NewAppError(err, "could not create TerraForm executor", http.StatusInternalServerError)
	}
	if err := ex.AddCredentials(credentials); err != nil {
		return nil, nil, -1, ctxh.NewAppError(err, "could not validate TerraForm credentials", http.StatusBadRequest)
	}
	return session, ex, executionID.(int), nil
}
