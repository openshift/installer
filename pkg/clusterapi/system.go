package clusterapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	azic "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	powervsic "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/clusterapi/internal/process"
	"github.com/openshift/installer/pkg/clusterapi/internal/process/addr"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	sys = &system{}
)

// SystemState is the state of the cluster-api system.
type SystemState string

const (
	// SystemStateRunning indicates the system is running.
	SystemStateRunning SystemState = "running"
	// SystemStateStopped indicates the system is stopped.
	SystemStateStopped SystemState = "stopped"

	// ArtifactsDir is the directory where output (manifests, kubeconfig, etc.)
	// related to CAPI-based installs are stored.
	ArtifactsDir = ".clusterapi_output"
)

// Interface is the interface for the cluster-api system.
type Interface interface {
	Run(ctx context.Context) error
	State() SystemState
	Client() client.Client
	Teardown()
	CleanEtcd()
}

// System returns the cluster-api system.
func System() Interface {
	return sys
}

// system creates a local capi control plane
// to use as a management cluster.
type system struct {
	sync.Mutex

	client client.Client

	componentDir string
	lcp          *localControlPlane

	wg           sync.WaitGroup
	teardownOnce sync.Once
	cancel       context.CancelFunc

	logWriter *io.PipeWriter
}

// hostHasIPv4Address verifies if the host that launches the host control plane has IPv4 address.
func hostHasIPv4Address() (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	for _, intf := range interfaces {
		if intf.Flags&net.FlagUp == 0 || intf.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := intf.Addrs()
		if err != nil {
			return false, err
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.To4() != nil {
				return true, nil
			}
		}
	}
	return false, nil
}

// Run launches the cluster-api system.
func (c *system) Run(ctx context.Context) error { //nolint:gocyclo
	c.Lock()
	defer c.Unlock()

	// Setup the context with a cancel function.
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	// Create the local control plane.
	lcp := &localControlPlane{}

	ipv4, err := hostHasIPv4Address()
	if err != nil {
		return err
	}
	// If the host has no IPv4 available, the default value of service network should be modified to IPv6 CIDR.
	if !ipv4 {
		lcp.APIServerArgs = map[string]string{
			"service-cluster-ip-range": "fd02::/112",
		}
	}

	if err := lcp.Run(ctx); err != nil {
		return fmt.Errorf("failed to run local control plane: %w", err)
	}
	c.lcp = lcp
	c.client = c.lcp.Client

	// Create a temporary directory to unpack the cluster-api assets
	// and use it as the working directory for the envtest environment.
	componentDir, err := os.MkdirTemp("", "openshift-cluster-api-system-components")
	if err != nil {
		return fmt.Errorf("failed to create temporary folder for cluster api components: %w", err)
	}
	if err := data.Unpack(componentDir, "/cluster-api"); err != nil {
		return fmt.Errorf("failed to unpack cluster api components: %w", err)
	}
	c.componentDir = componentDir

	// Create the controllers, we always need to run the cluster-api core controller.
	controllers := []*controller{
		{
			Name:       "Cluster API",
			Path:       fmt.Sprintf("%s/cluster-api", c.lcp.BinDir),
			Components: []string{c.componentDir + "/core-components.yaml"},
			Args: []string{
				"-v=2",
				"--diagnostics-address=0",
				"--health-addr={{suggestHealthHostPort}}",
				"--webhook-port={{.WebhookPort}}",
				"--webhook-cert-dir={{.WebhookCertDir}}",
			},
		},
	}

	metadata, err := metadata.Load(command.RootOpts.Dir)
	if err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	platform := metadata.Platform()
	if platform == "" {
		return fmt.Errorf("no platform configured in metadata")
	}

	// Create the infrastructure controllers.
	// Only add the controllers for the platform we are deploying to.
	switch platform {
	case aws.Name:
		controller := c.getInfrastructureController(
			&AWS,
			[]string{
				"-v=4",
				"--diagnostics-address=0",
				"--health-addr={{suggestHealthHostPort}}",
				"--webhook-port={{.WebhookPort}}",
				"--webhook-cert-dir={{.WebhookCertDir}}",
				"--feature-gates=BootstrapFormatIgnition=true,ExternalResourceGC=true,TagUnmanagedNetworkResources=false,EKS=false",
			},
			map[string]string{},
		)
		if cfg := metadata.AWS; cfg != nil && len(cfg.ServiceEndpoints) > 0 {
			endpoints := make([]string, 0, len(cfg.ServiceEndpoints))
			// CAPA expects name=url pairs of service endpoints
			for _, endpoint := range cfg.ServiceEndpoints {
				endpoints = append(endpoints, fmt.Sprintf("%s=%s", endpoint.Name, endpoint.URL))
			}
			controller.Args = append(controller.Args, fmt.Sprintf("--service-endpoints=%s:%s", cfg.Region, strings.Join(endpoints, ",")))
		}
		controllers = append(controllers, controller)
	case azure.Name:
		cloudName := metadata.Azure.CloudName
		if cloudName == "" {
			cloudName = azure.PublicCloud
		}
		session, err := azic.GetSession(cloudName, metadata.Azure.ARMEndpoint)
		if err != nil {
			return fmt.Errorf("unable to retrieve azure session: %w", err)
		}
		azProvider := Azure
		var envFP string
		if cloudName == azure.StackCloud {
			// Set provider so that the Azure Stack (forked) controller and CRDs are used.
			azProvider = AzureStack

			// Lay down the environment file so that cloud-provider-azure running in
			// CAPZ & ASO controllers can load the environment.
			b, err := json.Marshal(session.Environment)
			if err != nil {
				return errors.Wrap(err, "could not serialize Azure Stack endpoints")
			}
			envFP = filepath.Join(c.componentDir, "azurestackcloud.json")
			if err = os.WriteFile(envFP, b, 0600); err != nil {
				return fmt.Errorf("failed to write Azure Stack environment file: %w", err)
			}
		}

		controllers = append(controllers,
			c.getInfrastructureController(
				&azProvider,
				[]string{
					"-v=2",
					"--health-addr={{suggestHealthHostPort}}",
					"--webhook-port={{.WebhookPort}}",
					"--webhook-cert-dir={{.WebhookCertDir}}",
					"--feature-gates=MachinePool=false",
				},
				map[string]string{
					"AZURE_ENVIRONMENT_FILEPATH": envFP,
				},
			),
			c.getInfrastructureController(
				&AzureASO,
				[]string{
					"-v=0",
					"-metrics-addr=0",
					"-health-addr={{suggestHealthHostPort}}",
					"-webhook-port={{.WebhookPort}}",
					"-webhook-cert-dir={{.WebhookCertDir}}",
					"-crd-pattern=",
					"-crd-management=none",
				}, map[string]string{
					"POD_NAMESPACE":                     "capz-system",
					"AZURE_CLIENT_ID":                   session.Credentials.ClientID,
					"AZURE_CLIENT_SECRET":               session.Credentials.ClientSecret,
					"AZURE_CLIENT_CERTIFICATE":          session.Credentials.ClientCertificatePath,
					"AZURE_CLIENT_CERTIFICATE_PASSWORD": session.Credentials.ClientCertificatePassword,
					"AZURE_TENANT_ID":                   session.Credentials.TenantID,
					"AZURE_SUBSCRIPTION_ID":             session.Credentials.SubscriptionID,
					"AZURE_RESOURCE_MANAGER_ENDPOINT":   session.Environment.ResourceManagerEndpoint,
					"AZURE_RESOURCE_MANAGER_AUDIENCE":   session.Environment.TokenAudience,
					"AZURE_ENVIRONMENT_FILEPATH":        envFP,
				},
			),
		)
	case gcp.Name:
		session, err := gcpic.GetSession(context.Background())
		if err != nil {
			return fmt.Errorf("failed to create gcp session: %w", err)
		}

		//nolint:gosec // CAPG only expects a single credentials environment variable
		gAppCredEnvVar := "GOOGLE_APPLICATION_CREDENTIALS"
		capgEnvVars := map[string]string{
			gAppCredEnvVar: session.Path,
		}

		if v, ok := capgEnvVars[gAppCredEnvVar]; ok {
			logrus.Infof("setting %q to %s for capg infrastructure controller", gAppCredEnvVar, v)
		}

		controllers = append(controllers,
			c.getInfrastructureController(
				&GCP,
				[]string{
					"-v=2",
					"--diagnostics-address=0",
					"--health-addr={{suggestHealthHostPort}}",
					"--webhook-port={{.WebhookPort}}",
					"--webhook-cert-dir={{.WebhookCertDir}}",
				},
				capgEnvVars,
			),
		)
	case ibmcloud.Name:
		ibmcloudFlags := []string{
			"--provider-id-fmt=v2",
			"-v=2",
			"--health-addr={{suggestHealthHostPort}}",
			"--leader-elect=false",
			"--webhook-port={{.WebhookPort}}",
			"--webhook-cert-dir={{.WebhookCertDir}}",
			fmt.Sprintf("--namespace=%s", capiutils.Namespace),
		}

		// Get the ServiceEndpoint overrides, along with Region, to pass on to CAPI, if any.
		if serviceEndpoints := metadata.IBMCloud.GetRegionAndEndpointsFlag(); serviceEndpoints != "" {
			ibmcloudFlags = append(ibmcloudFlags, fmt.Sprintf("--service-endpoint=%s", serviceEndpoints))
		}

		iamEndpoint := "https://iam.cloud.ibm.com"
		// Override IAM endpoint if an override was provided.
		if overrideURL := ibmcloud.CheckServiceEndpointOverride(configv1.IBMCloudServiceIAM, metadata.IBMCloud.ServiceEndpoints); overrideURL != "" {
			iamEndpoint = overrideURL
		}

		controllers = append(controllers,
			c.getInfrastructureController(
				&IBMCloud,
				ibmcloudFlags,
				map[string]string{
					"IBMCLOUD_AUTH_TYPE": "iam",
					"IBMCLOUD_APIKEY":    os.Getenv("IC_API_KEY"),
					"IBMCLOUD_AUTH_URL":  iamEndpoint,
					"LOGLEVEL":           "5",
				},
			),
		)
	case nutanix.Name:
		controllers = append(controllers,
			c.getInfrastructureController(
				&Nutanix,
				[]string{
					"--diagnostics-address=0",
					"--health-probe-bind-address={{suggestHealthHostPort}}",
					"--leader-elect=false",
				},
				map[string]string{},
			),
		)
	case openstack.Name:
		controllers = append(controllers,
			c.getInfrastructureController(
				&OpenStack,
				[]string{
					"-v=2",
					"--diagnostics-address=0",
					"--health-addr={{suggestHealthHostPort}}",
					"--webhook-port={{.WebhookPort}}",
					"--webhook-cert-dir={{.WebhookCertDir}}",
				},
				map[string]string{
					"EXP_KUBEADM_BOOTSTRAP_FORMAT_IGNITION": "true",
				},
			),
		)
	case vsphere.Name:
		controllers = append(controllers,
			c.getInfrastructureController(
				&VSphere,
				[]string{
					"-v=2",
					"--diagnostics-address=0",
					"--health-addr={{suggestHealthHostPort}}",
					"--webhook-port={{.WebhookPort}}",
					"--webhook-cert-dir={{.WebhookCertDir}}",
					"--leader-elect=false",
				},
				map[string]string{
					"EXP_KUBEADM_BOOTSTRAP_FORMAT_IGNITION": "true",
					"EXP_CLUSTER_RESOURCE_SET":              "true",
				},
			),
		)
	case powervs.Name:
		// We need to prompt for missing variables because NewPISession requires them!
		bxClient, err := powervsic.NewBxClient(true)
		if err != nil {
			return fmt.Errorf("failed to create a BxClient in Run: %w", err)
		}
		APIKey := bxClient.GetBxClientAPIKey()

		controller := c.getInfrastructureController(
			&IBMCloud,
			[]string{
				"--provider-id-fmt=v2",
				"--v=2",
				"--health-addr={{suggestHealthHostPort}}",
				"--webhook-port={{.WebhookPort}}",
				"--webhook-cert-dir={{.WebhookCertDir}}",
			},
			map[string]string{
				"IBMCLOUD_AUTH_TYPE": "iam",
				"IBMCLOUD_APIKEY":    APIKey,
				"IBMCLOUD_AUTH_URL":  "https://iam.cloud.ibm.com",
				"LOGLEVEL":           "2",
			},
		)
		if cfg := metadata.PowerVS; cfg != nil {
			overrides := bxClient.MapServiceEndpointsForCAPI(cfg)
			if len(overrides) > 0 {
				controller.Args = append(controller.Args, fmt.Sprintf("--service-endpoint=%s:%s", cfg.Region, strings.Join(overrides, ",")))
			}
		}
		controllers = append(controllers, controller)
	default:
		return fmt.Errorf("unsupported platform %q", platform)
	}

	// We only show controller logs if the log level is DEBUG or above
	c.logWriter = logrus.StandardLogger().WriterLevel(logrus.DebugLevel)

	// We create a wait group to wait for the controllers to stop,
	// this waitgroup is a global, and is used by the Teardown function
	// which is expected to be called when the program exits.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		// Stop the controllers when the context is cancelled.
		<-ctx.Done()
		logrus.Info("Shutting down local Cluster API controllers...")
		for _, ct := range controllers {
			if ct.state != nil {
				if err := ct.state.Stop(); err != nil {
					logrus.Warnf("Failed to stop controller: %s: %v", ct.Name, err)
					continue
				}
				logrus.Infof("Stopped controller: %s", ct.Name)
			}
		}
	}()

	// Run the controllers.
	for _, ct := range controllers {
		if err := c.runController(ctx, ct); err != nil {
			return fmt.Errorf("failed to run controller %q: %w", ct.Name, err)
		}
	}

	return nil
}

// Client returns the client for the local control plane.
func (c *system) Client() client.Client {
	c.Lock()
	defer c.Unlock()

	return c.client
}

// Teardown shuts down the local capi control plane and all its controllers.
func (c *system) Teardown() {
	c.Lock()
	defer c.Unlock()

	if c.lcp == nil {
		return
	}

	// Clean up the binary directory.
	defer os.RemoveAll(c.lcp.BinDir)

	// Clean up log file handles.
	defer c.lcp.EtcdLog.Close()
	defer c.lcp.APIServerLog.Close()

	// Proceed to shutdown.
	c.teardownOnce.Do(func() {
		c.cancel()
		ch := make(chan struct{})
		go func() {
			c.wg.Wait()
			logrus.Info("Shutting down local Cluster API control plane...")
			if err := c.lcp.Stop(); err != nil {
				logrus.Warnf("Failed to stop local Cluster API control plane: %v", err)
			}
			close(ch)
		}()
		select {
		case <-ch:
			logrus.Info("Local Cluster API system has completed operations")
		case <-time.After(60 * time.Second):
			logrus.Warn("Timed out waiting for local Cluster API system to shut down")
		}

		c.logWriter.Close()
	})
}

// CleanEtcd removes the etcd database from the host.
func (c *system) CleanEtcd() {
	c.Lock()
	defer c.Unlock()

	if c.lcp == nil {
		return
	}

	// Clean up the etcd directory.
	if err := os.RemoveAll(c.lcp.EtcdDataDir); err != nil {
		logrus.Warnf("Unable to delete local etcd data directory %s. It is safe to remove the directory manually", c.lcp.EtcdDataDir)
	}
}

// State returns the state of the cluster-api system.
func (c *system) State() SystemState {
	c.Lock()
	defer c.Unlock()

	if c.lcp == nil {
		return SystemStateStopped
	}
	return SystemStateRunning
}

// getInfrastructureController returns a controller for the given provider,
// most of the configuration is by convention.
//
// The provider is expected to be compiled as part of the release process, and packaged in the binaries directory
// and have the name `cluster-api-provider-<name>`.
//
// While the manifests can be optional, we expect them to be in the manifests directory and named `<name>-infrastructure-components.yaml`.
func (c *system) getInfrastructureController(provider *Provider, args []string, env map[string]string) *controller {
	manifests := []string{}
	defaultManifestPath := filepath.Join(c.componentDir, fmt.Sprintf("/%s-infrastructure-components.yaml", provider.Name))
	if _, err := os.Stat(defaultManifestPath); err == nil {
		manifests = append(manifests, defaultManifestPath)
	} else {
		logrus.Infof("Failed to find manifests for provider %s at %s", provider.Name, defaultManifestPath)
	}
	return &controller{
		Provider:   provider,
		Name:       fmt.Sprintf("%s infrastructure provider", provider.Name),
		Path:       fmt.Sprintf("%s/cluster-api-provider-%s", c.lcp.BinDir, provider.Name),
		Components: manifests,
		Args:       args,
		Env:        env,
	}
}

// controller encapsulates the state of a controller, its process state, and its configuration.
type controller struct {
	Provider *Provider
	state    *process.State

	Name       string
	Dir        string
	Path       string
	Components []string
	Args       []string
	Env        map[string]string
}

// runController configures the controller, and waits for it to be ready.
func (c *system) runController(ctx context.Context, ct *controller) error {
	// If the provider is not empty, we extract it to the binaries directory.
	if ct.Provider != nil {
		if err := ct.Provider.Extract(c.lcp.BinDir); err != nil {
			return fmt.Errorf("failed to extract provider %q: %w", ct.Name, err)
		}
	}

	// Create the WebhookInstallOptions from envtest, and pass the manifests we've been given as input.
	// Once built, we install them in the local control plane using the rest.Config available.
	// Envtest takes care of a few things needed to run webhooks locally:
	// - Creates a self-signed certificate for the webhook server.
	// - Tries to allocate a host:port for the webhook server to listen on.
	// - Modifies the webhook manifests to point to the local webhook server through a URL and a CABundle.
	wh := envtest.WebhookInstallOptions{
		Paths:                   ct.Components,
		IgnoreSchemeConvertible: true,
	}
	if err := wh.Install(c.lcp.Cfg); err != nil {
		return fmt.Errorf("failed to prepare controller %q webhook options: %w", ct.Name, err)
	}

	// Most providers allocate a host:port configuration for the health check,
	// which responds to a simple http request on /healthz and /readyz.
	// When an argument is configured to use the suggestHealthHostPort function,
	// we record the value, so we can pass it to
	var healthCheckHostPort string

	// Build the arguments, using go templating to render the values.
	{
		funcs := template.FuncMap{
			"suggestHealthHostPort": func() (string, error) {
				healthPort, healthHost, err := addr.Suggest("")
				if err != nil {
					return "", fmt.Errorf("unable to grab random port: %w", err)
				}
				healthCheckHostPort = fmt.Sprintf("%s:%d", healthHost, healthPort)
				return healthCheckHostPort, nil
			},
		}

		templateData := map[string]string{
			"WebhookPort":    fmt.Sprintf("%d", wh.LocalServingPort),
			"WebhookCertDir": wh.LocalServingCertDir,
			"KubeconfigPath": c.lcp.KubeconfigPath,
		}

		// We cannot override KUBECONFIG, e.g., in case the user supplies a callback that needs to access the cluster,
		// such as via credential_process in the AWS config file. The kubeconfig path is set in the controller instead.
		if ct.Provider == nil || ct.Provider.Name != "azureaso" {
			ct.Args = append(ct.Args, "--kubeconfig={{.KubeconfigPath}}")
		}

		args := make([]string, 0, len(ct.Args))
		for _, arg := range ct.Args {
			final := new(bytes.Buffer)
			tmpl := template.Must(template.New("arg").Funcs(funcs).Parse(arg))
			if err := tmpl.Execute(final, templateData); err != nil {
				return fmt.Errorf("failed to render controller %q arg %q: %w", ct.Name, arg, err)
			}
			args = append(args, strings.TrimSpace(final.String()))
		}
		ct.Args = args
	}

	// Build the environment variables.
	env := []string{}
	{
		if ct.Env == nil {
			ct.Env = map[string]string{}
		}
		// Override KUBECONFIG to point to the local control plane.
		// azureaso doesn't support the --kubeconfig parameter.
		if ct.Provider != nil && ct.Provider.Name == "azureaso" {
			ct.Env["KUBECONFIG"] = c.lcp.KubeconfigPath
		}
		for key, value := range ct.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Install the manifests for the controller, if any.
	if len(ct.Components) > 0 {
		opts := envtest.CRDInstallOptions{
			Scheme:         c.lcp.Env.Scheme,
			Paths:          ct.Components,
			WebhookOptions: wh,
		}
		if _, err := envtest.InstallCRDs(c.lcp.Cfg, opts); err != nil {
			return fmt.Errorf("failed to install controller %q manifests in local control plane: %w", ct.Name, err)
		}
	}

	// Create the process state.
	pr := &process.State{
		Path:         ct.Path,
		Args:         ct.Args,
		Dir:          ct.Dir,
		Env:          env,
		StartTimeout: 60 * time.Second,
		StopTimeout:  10 * time.Second,
	}

	// If the controller has a health check, we configure it, and wait for it to be ready.
	if healthCheckHostPort != "" {
		pr.HealthCheck = &process.HealthCheck{
			URL: url.URL{
				Scheme: "http",
				Host:   healthCheckHostPort,
				Path:   "/healthz",
			},
		}
	}

	// Initialize the process state.
	if err := pr.Init(ct.Name); err != nil {
		return fmt.Errorf("failed to initialize process state for controller %q: %w", ct.Name, err)
	}

	// Run the controller and store its state.
	logrus.Infof("Running process: %s with args %v", ct.Name, ct.Args)
	if err := pr.Start(ctx, c.logWriter, c.logWriter); err != nil {
		return fmt.Errorf("failed to start controller %q: %w", ct.Name, err)
	}
	ct.state = pr
	return nil
}
