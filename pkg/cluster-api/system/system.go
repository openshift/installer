package system

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset/installconfig"
	providers "github.com/openshift/installer/pkg/cluster-api"
	"github.com/openshift/installer/pkg/cluster-api/system/process"
	"github.com/openshift/installer/pkg/cluster-api/system/process/addr"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func init() {
	// Set verbose logging for controller-runtime.
	flag.Set("v", "2")
	flag.Parse()
}

var (
	wg          = &sync.WaitGroup{}
	ctx, cancel = context.WithCancel(signals.SetupSignalHandler())
	once        = &sync.Once{}
)

// Teardown shuts down the local capi control plane and all its controllers.
func Teardown() {
	once.Do(func() {
		cancel()
		logrus.Info("Shutting down local Cluster API control plane...")
		ch := make(chan struct{})
		go func() {
			wg.Wait()
			close(ch)
		}()
		select {
		case <-ch:
			logrus.Info("Local control plane has completed operations")
		case <-time.After(30 * time.Second):
			logrus.Warn("Timed out waiting for local control plane to shut down")
		}
	})
}

// System creates a local capi control plane
// to use as a management cluster.
// TODO: Add support for existing management cluster.
type System struct {
	manifestDir string
	lcp         *localControlPlane
	Client      client.Client
}

// Name returns the human-friendly name of the asset.
func (c *System) Name() string {
	return "Cluster API System"
}

// Run launches the cluster-api system.
func (c *System) Run(clusterID *installconfig.ClusterID, installConfig *installconfig.InstallConfig) (err error) {
	c.lcp = &localControlPlane{}
	if err := c.lcp.Run(clusterID, installConfig); err != nil {
		return fmt.Errorf("failed to run local control plane: %w", err)
	}
	c.Client = c.lcp.Client

	// Create a temporary directory to unpack the cluster-api assets
	// and use it as the working directory for the envtest environment.
	manifestDir, err := os.MkdirTemp("", "openshift-cluster-api-manifests")
	if err != nil {
		return err
	}
	if err := data.Unpack(manifestDir, "/cluster-api"); err != nil {
		return err
	}
	c.manifestDir = manifestDir

	controllers := []*controller{
		{
			Name:      "Cluster API",
			Path:      fmt.Sprintf("%s/cluster-api", c.lcp.BinDir),
			Manifests: []string{c.manifestDir + "/core-components.yaml"},
			Args: []string{
				"-v=2",
				"--metrics-bind-addr=0",
				"--health-addr={{suggestHealthHostPort}}",
				"--webhook-port={{.WebhookPort}}",
				"--webhook-cert-dir={{.WebhookCertDir}}",
			},
		},
	}

	platform := installConfig.Config.Platform.Name()

	switch platform {
	case aws.Name:
		controllers = append(controllers,
			c.getInfrastructureController(
				&providers.AWS,
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
				&providers.Azure,
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
				&providers.AzureASO,
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
	case ibmcloud.Name:
	case nutanix.Name:
	case vsphere.Name:
	default:
		return fmt.Errorf("unsupported platform %q", platform)
	}

	// Run the controllers.
	for _, ct := range controllers {
		if err := c.runController(ctx, ct); err != nil {
			return fmt.Errorf("failed to run controller %q: %w", ct.Name, err)
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Stop the controllers when the context is cancelled.
		<-ctx.Done()
		for _, ct := range controllers {
			if ct.state != nil {
				if err := ct.state.Stop(); err != nil {
					logrus.Warnf("Failed to stop local manager: %s: %v", ct.Name, err)
					continue
				}
				logrus.Infof("Stopped local manager: %s", ct.Name)
			}
		}
	}()

	return nil
}

func (c *System) getInfrastructureController(provider *providers.Provider, args []string, env map[string]string) *controller {
	manifests := []string{}
	defaultManifestPath := filepath.Join(c.manifestDir, fmt.Sprintf("/%s-infrastructure-components.yaml", provider.Name))
	if _, err := os.Stat(defaultManifestPath); err == nil {
		manifests = append(manifests, defaultManifestPath)
	}
	return &controller{
		Provider:  provider,
		Name:      fmt.Sprintf("%s infrastructure provider", provider.Name),
		Path:      fmt.Sprintf("%s/cluster-api-provider-%s_%s_%s", filepath.Join(c.lcp.BinDir, provider.Source), provider.Name, runtime.GOOS, runtime.GOARCH),
		Manifests: manifests,
		Args:      args,
		Env:       env,
	}
}

type controller struct {
	Provider *providers.Provider
	state    *process.State

	Name      string
	Dir       string
	Path      string
	Manifests []string
	Args      []string
	Env       map[string]string
}

func (c *System) runController(ctx context.Context, ct *controller) error {
	if ct.Provider != nil {
		if err := ct.Provider.Extract(c.lcp.BinDir); err != nil {
			logrus.Fatal(err)
		}
	}

	wh := envtest.WebhookInstallOptions{
		Paths:                   ct.Manifests,
		IgnoreSchemeConvertible: true,
	}
	if err := wh.Install(c.lcp.Cfg); err != nil {
		return fmt.Errorf("failed to prepare controller %q webhook options: %w", ct.Name, err)
	}

	var healthHost string
	var healthPort int
	funcs := template.FuncMap{
		"suggestHealthHostPort": func() (string, error) {
			var err error
			healthPort, healthHost, err = addr.Suggest("")
			if err != nil {
				return "", fmt.Errorf("unable to grab random port: %w", err)
			}
			return fmt.Sprintf("%s:%d", healthHost, healthPort), nil
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

	// Override KUBECONFIG to point to the local control plane.
	env := []string{}
	if ct.Env == nil {
		ct.Env = map[string]string{}
	}
	ct.Env["KUBECONFIG"] = c.lcp.KubeconfigPath
	for key, value := range ct.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Install the manifests for the controller, if any.
	if len(ct.Manifests) > 0 {
		opts := envtest.CRDInstallOptions{
			Scheme:         c.lcp.Env.Scheme,
			Paths:          ct.Manifests,
			WebhookOptions: wh,
		}
		if _, err := envtest.InstallCRDs(c.lcp.Cfg, opts); err != nil {
			return fmt.Errorf("failed to install controller %q manifests in local control plane: %w", ct.Name, err)
		}
	}

	logrus.Infof("Running process: %s with args %v and env %v", ct.Name, ct.Args, env)
	pr := &process.State{
		Path:         ct.Path,
		Args:         ct.Args,
		Dir:          ct.Dir,
		Env:          env,
		StartTimeout: 60 * time.Second,
		StopTimeout:  10 * time.Second,
		HealthCheck: process.HealthCheck{
			URL: url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", healthHost, healthPort),
				Path:   "/healthz",
			},
		},
	}
	if err := pr.Init(ct.Name); err != nil {
		return fmt.Errorf("failed to initialize process state for controller %q: %w", ct.Name, err)
	}
	if err := pr.Start(ctx, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("failed to start controller %q: %w", ct.Name, err)
	}
	ct.state = pr
	return nil
}
