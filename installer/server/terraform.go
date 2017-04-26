package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"

	"github.com/dghubble/sessions"

	"github.com/coreos/tectonic-installer/installer/server/asset"
	"github.com/coreos/tectonic-installer/installer/server/ctxh"
	"github.com/coreos/tectonic-installer/installer/server/defaults"
	"github.com/coreos/tectonic-installer/installer/server/terraform"
)

func newAWSTerraformVars(c *TectonicAWSCluster) ([]asset.Asset, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(c.Tectonic.IdentityAdminPassword, bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt failed: %v", err)
	}

	etcdServers := []string{}
	etcdCount := c.CloudForm.ETCDCount
	if c.CloudForm.ExternalETCDClient != "" {
		etcdServers = append(etcdServers, strings.Split(c.CloudForm.ExternalETCDClient, ":")[0])
		etcdCount = 0
	}

	availabilityZones := make(map[string]struct{})
	for _, az := range c.CloudForm.ControllerSubnets {
		availabilityZones[az.AvailabilityZone] = struct{}{}
	}
	for _, az := range c.CloudForm.WorkerSubnets {
		availabilityZones[az.AvailabilityZone] = struct{}{}
	}

	controllerSubnetIDs := []string{}
	workerSubnetIDs := []string{}
	if c.CloudForm.VPCID != "" {
		for _, subnet := range c.CloudForm.ControllerSubnets {
			controllerSubnetIDs = append(controllerSubnetIDs, subnet.ID)
		}
		for _, subnet := range c.CloudForm.WorkerSubnets {
			workerSubnetIDs = append(workerSubnetIDs, subnet.ID)
		}
	}

	variables := map[string]interface{}{
		"tectonic_cluster_name": c.CloudForm.ClusterName,
		"tectonic_base_domain":  c.CloudForm.HostedZoneName,

		"tectonic_license_path":     "./license",
		"tectonic_pull_secret_path": "./pull_secret",

		"tectonic_admin_email":         c.Tectonic.IdentityAdminUser,
		"tectonic_admin_password_hash": string(passwordHash),

		"tectonic_ca_cert":    c.CACertificate,
		"tectonic_ca_key":     c.CAPrivateKey,
		"tectonic_ca_key_alg": "RSA",

		"tectonic_cl_channel":   c.CloudForm.Channel,
		"tectonic_cluster_cidr": c.CloudForm.PodCIDR,
		"tectonic_service_cidr": c.CloudForm.ServiceCIDR,

		"tectonic_etcd_count":   etcdCount,
		"tectonic_etcd_servers": etcdServers,

		"tectonic_update_app_id":  c.Tectonic.Updater.AppID,
		"tectonic_update_channel": c.Tectonic.Updater.Channel,
		"tectonic_update_server":  c.Tectonic.Updater.Server,

		"tectonic_master_count": c.CloudForm.ControllerCount,
		"tectonic_worker_count": c.CloudForm.WorkerCount,

		"tectonic_kube_apiserver_service_ip": c.CloudForm.APIServiceIP.String(),
		"tectonic_kube_dns_service_ip":       c.CloudForm.DNSServiceIP.String(),

		"tectonic_aws_external_vpc_id":            c.CloudForm.VPCID,
		"tectonic_aws_external_master_subnet_ids": controllerSubnetIDs,
		"tectonic_aws_external_worker_subnet_ids": workerSubnetIDs,
		"tectonic_aws_vpc_cidr_block":             c.CloudForm.VPCCIDR,
		"tectonic_aws_az_count":                   len(availabilityZones),
		"tectonic_aws_master_ec2_type":            c.CloudForm.ControllerInstanceType,
		"tectonic_aws_worker_ec2_type":            c.CloudForm.WorkerInstanceType,
		"tectonic_aws_etcd_ec2_type":              c.CloudForm.ETCDInstanceType,
		"tectonic_aws_ssh_key":                    c.CloudForm.KeyName,
	}

	tfVars, err := mapVarsToTFVars(variables)
	if err != nil {
		return []asset.Asset{}, err
	}

	return []asset.Asset{
		asset.New("terraform/terraform.tfvars", []byte(tfVars)),
		asset.New("terraform/license", []byte(c.Tectonic.License)),
		asset.New("terraform/pull_secret", []byte(c.Tectonic.Dockercfg)),
	}, nil
}

func mapVarsToTFVars(variables map[string]interface{}) (string, error) {
	tfVars := ""

	for key, value := range variables {
		var stringValue string

		switch value := value.(type) {
		case string:
			trimmedValue := strings.Trim(value, "\n")
			if !strings.Contains(trimmedValue, "\n") {
				stringValue = fmt.Sprintf("\"%s\"", trimmedValue)
			} else {
				stringValue = fmt.Sprintf("<<EOD\n%s\nEOD", trimmedValue)
			}
		case []string:
			qValue := make([]string, len(value))
			for i := 0; i < len(value); i++ {
				qValue[i] = fmt.Sprintf("\"%s\"", strings.Trim(value[i], "\n"))
			}
			stringValue = fmt.Sprintf("[%s]", strings.Join(qValue, ", "))
		case int:
			stringValue = strconv.Itoa(value)
		default:
			return "", fmt.Errorf("unsupported type %T (%s) for TFVars\n", value, key)
		}

		tfVars = fmt.Sprintf("%s%s = %s\n", tfVars, key, stringValue)
	}

	return tfVars, nil
}

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
			return errCtx
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
	// Create a new Executor.
	ex, err := terraform.NewExecutor()
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
		return nil, nil, -1, ctxh.NewAppError(err, "could not find execution data, apply first", http.StatusNotFound)
	}
	executionPath, ok := session.Values["terraform_path"]
	if !ok {
		return nil, nil, -1, ctxh.NewAppError(err, "could not find execution data, apply first", http.StatusNotFound)
	}
	executionID, ok := session.Values["terraform_id"]
	if !ok {
		return nil, nil, -1, ctxh.NewAppError(err, "could not find execution data, apply first", http.StatusNotFound)
	}
	ex, err := terraform.NewExecutorFromPath(executionPath.(string))
	if err != nil {
		return nil, nil, -1, ctxh.NewAppError(err, "could not create TerraForm executor", http.StatusInternalServerError)
	}
	if err := ex.AddCredentials(credentials); err != nil {
		return nil, nil, -1, ctxh.NewAppError(err, "could not validate TerraForm credentials", http.StatusBadRequest)
	}
	return session, ex, executionID.(int), nil
}
