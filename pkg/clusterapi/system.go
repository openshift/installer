package clusterapi

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/clusterapi/internal/process"
	"github.com/openshift/installer/pkg/clusterapi/internal/process/addr"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var (
	sys = &system{}
)

// Interface is the interface for the cluster-api system.
type Interface interface {
	Run(ctx context.Context, installConfig *installconfig.InstallConfig) error
	Client() client.Client
	Teardown()
}

// System returns the cluster-api system.
func System() Interface {
	return sys
}

// system creates a local capi control plane
// to use as a management cluster.
type system struct {
	client client.Client

	componentDir string
	lcp          *localControlPlane

	wg           sync.WaitGroup
	teardownOnce sync.Once
	cancel       context.CancelFunc
}

// Run launches the cluster-api system.
func (c *system) Run(ctx context.Context, installConfig *installconfig.InstallConfig) error {
	// Setup the context with a cancel function.
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	// Create the local control plane.
	c.lcp = &localControlPlane{}
	if err := c.lcp.Run(ctx); err != nil {
		return fmt.Errorf("failed to run local control plane: %w", err)
	}
	c.client = c.lcp.Client

	// Create a temporary directory to unpack the cluster-api assets
	// and use it as the working directory for the envtest environment.
	componentDir, err := os.MkdirTemp("", "openshift-cluster-api-system-components")
	if err != nil {
		return err
	}
	if err := data.Unpack(componentDir, "/cluster-api"); err != nil {
		return err
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
				"--metrics-bind-addr=0",
				"--health-addr={{suggestHealthHostPort}}",
				"--webhook-port={{.WebhookPort}}",
				"--webhook-cert-dir={{.WebhookCertDir}}",
			},
		},
	}

	// Create the infrastructure controllers.
	// Only add the controllers for the platform we are deploying to.
	switch platform := installConfig.Config.Platform.Name(); platform {
	case aws.Name:
		controllers = append(controllers,
			c.getInfrastructureController(
				&AWS,
				[]string{
					"-v=2",
					"--metrics-bind-addr=0",
					"--health-addr={{suggestHealthHostPort}}",
					"--webhook-port={{.WebhookPort}}",
					"--webhook-cert-dir={{.WebhookCertDir}}",
					"--feature-gates=BootstrapFormatIgnition=true,ExternalResourceGC=true",
				},
				map[string]string{},
			),
		)
	case azure.Name:
		session, err := installConfig.Azure.Session()
		if err != nil {
			return fmt.Errorf("failed to create azure session: %w", err)
		}

		controllers = append(controllers,
			c.getInfrastructureController(
				&Azure,
				[]string{
					"-v=2",
					"--metrics-bind-addr=0",
					"--health-addr={{suggestHealthHostPort}}",
					"--webhook-port={{.WebhookPort}}",
					"--webhook-cert-dir={{.WebhookCertDir}}",
				},
				map[string]string{},
			),
			c.getInfrastructureController(
				&AzureASO,
				[]string{
					"--v=0",
					"--metrics-addr=0",
					"--health-addr={{suggestHealthHostPort}}",
					"--webhook-port={{.WebhookPort}}",
					"--webhook-cert-dir={{.WebhookCertDir}}",
					"--crd-pattern=",
					"--enable-crd-management=false",
				}, map[string]string{
					"POD_NAMESPACE":                     "capz-system",
					"AZURE_CLIENT_ID":                   session.Credentials.ClientID,
					"AZURE_CLIENT_SECRET":               session.Credentials.ClientSecret,
					"AZURE_CLIENT_CERTIFICATE":          session.Credentials.ClientCertificatePath,
					"AZURE_CLIENT_CERTIFICATE_PASSWORD": session.Credentials.ClientCertificatePassword,
					"AZURE_TENANT_ID":                   session.Credentials.TenantID,
					"AZURE_SUBSCRIPTION_ID":             session.Credentials.SubscriptionID,
				},
			),
		)
	case gcp.Name:
		// TODO
	case ibmcloud.Name:
		// TODO
	case nutanix.Name:
		// TODO
	case vsphere.Name:
		// TODO
	default:
		return fmt.Errorf("unsupported platform %q", platform)
	}

	// Run the controllers.
	for _, ct := range controllers {
		if err := c.runController(ctx, ct); err != nil {
			return fmt.Errorf("failed to run controller %q: %w", ct.Name, err)
		}
	}

	// We create a wait group to wait for the controllers to stop,
	// this waitgroup is a global, and is used by the Teardown function
	// which is expected to be called when the program exits.
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		// Stop the controllers when the context is cancelled.
		<-ctx.Done()
		for _, ct := range controllers {
			if ct.state != nil {
				if err := ct.state.Stop(); err != nil {
					logrus.Warnf("Failed to stop controller: %s: %v", ct.Name, err)
					continue
				}
				logrus.Infof("Stopped controller: %s", ct.Name)
			}
		}

		// Stop the local control plane.
		if err := c.lcp.Stop(); err != nil {
			logrus.Warnf("Failed to stop local Cluster API control plane: %v", err)
		}
	}()

	return nil
}

// Client returns the client for the local control plane.
func (c *system) Client() client.Client {
	return c.client
}

// Teardown shuts down the local capi control plane and all its controllers.
func (c *system) Teardown() {
	if c.lcp == nil {
		return
	}

	// Clean up the binary directory.
	defer os.RemoveAll(c.lcp.BinDir)

	// Proceed to shutdown.
	c.teardownOnce.Do(func() {
		c.cancel()
		logrus.Info("Shutting down local Cluster API control plane...")
		ch := make(chan struct{})
		go func() {
			c.wg.Wait()
			close(ch)
		}()
		select {
		case <-ch:
			logrus.Info("Local Cluster API system has completed operations")
		case <-time.After(60 * time.Second):
			logrus.Warn("Timed out waiting for local Cluster API system to shut down")
		}
	})
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
			logrus.Fatal(err)
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
		ct.Env["KUBECONFIG"] = c.lcp.KubeconfigPath
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
	logrus.Infof("Running process: %s with args %v and env %v", ct.Name, ct.Args, env)
	if err := pr.Start(ctx, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("failed to start controller %q: %w", ct.Name, err)
	}
	ct.state = pr
	return nil
}
