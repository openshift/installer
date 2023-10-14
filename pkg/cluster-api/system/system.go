package system

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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
)

var (
	wg          = &sync.WaitGroup{}
	ctx, cancel = context.WithCancel(signals.SetupSignalHandler())
)

// Teardown shuts down the local capi control plane and all its controllers.
func Teardown() {
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
}

// System creates a local capi control plane
// to use as a management cluster.
// TODO: Add support for existing management cluster.
type System struct {
	lcp    *localControlPlane
	Client client.Client
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

	controllers := []*controller{
		{
			Name:      "Cluster API",
			Path:      fmt.Sprintf("%s/cluster-api", c.lcp.BinDir),
			Manifests: []string{manifestDir + "/core-components.yaml"},
		},
		{
			Name:      "AWS Infrastructure Provider",
			Path:      fmt.Sprintf("%s/cluster-api-provider-%s_%s_%s", filepath.Join(c.lcp.BinDir, providers.AWS.Source), providers.AWS.Name, runtime.GOOS, runtime.GOARCH),
			Manifests: []string{manifestDir + "/aws-infrastructure-components.yaml"},
			Args:      []string{"--feature-gates=BootstrapFormatIgnition=true,ExternalResourceGC=true"},
		},
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

type controller struct {
	state *process.State

	Name      string
	Path      string
	Manifests []string
	Args      []string
}

func (c *System) runController(ctx context.Context, ct *controller) error {
	wh := envtest.WebhookInstallOptions{
		Paths: ct.Manifests,
	}
	if err := wh.Install(c.lcp.Cfg); err != nil {
		return fmt.Errorf("failed to prepare controller %q webhook options: %w", ct.Name, err)
	}
	port, host, err := addr.Suggest("")
	if err != nil {
		return fmt.Errorf("unable to grab random port for serving health checks on: %w", err)
	}

	// TODO(vincepri): Check if these args have already been set, and overwrite.
	ct.Args = append(ct.Args,
		"-v=2",
		"--metrics-bind-addr=:0",
		fmt.Sprintf("--health-addr=%s:%d", host, port),
		fmt.Sprintf("--kubeconfig=%s", c.lcp.KubeconfigPath),
		fmt.Sprintf("--webhook-port=%d", wh.LocalServingPort),
		fmt.Sprintf("--webhook-cert-dir=%s", wh.LocalServingCertDir),
	)
	opts := envtest.CRDInstallOptions{
		Scheme:         c.lcp.Env.Scheme,
		Paths:          ct.Manifests,
		WebhookOptions: wh,
	}
	if _, err := envtest.InstallCRDs(c.lcp.Cfg, opts); err != nil {
		return fmt.Errorf("failed to install controller %q manifests in local control plane: %w", ct.Name, err)
	}
	pr := &process.State{
		Path:         ct.Path,
		Args:         ct.Args,
		StartTimeout: 10 * time.Second,
		StopTimeout:  10 * time.Second,
		HealthCheck: process.HealthCheck{
			URL: url.URL{
				Scheme: "http",
				Host:   fmt.Sprintf("%s:%d", host, port),
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
