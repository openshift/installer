package system

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/klog/v2"
	capav1beta1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1"
	capav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	capzv1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	clusterv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3" //nolint:staticcheck
	clusterv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4" //nolint:staticcheck
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/asset/installconfig"
	providers "github.com/openshift/installer/pkg/cluster-api"
)

func init() {
	utilruntime.Must(clusterv1alpha3.AddToScheme(scheme.Scheme))
	utilruntime.Must(clusterv1alpha4.AddToScheme(scheme.Scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme.Scheme))
	utilruntime.Must(capav1beta1.AddToScheme(scheme.Scheme))
	utilruntime.Must(capav1.AddToScheme(scheme.Scheme))
	utilruntime.Must(capzv1.AddToScheme(scheme.Scheme))
}

// localControlPlane creates a local capi control plane
// to use as a management cluster.
// TODO: Add support for existing management cluster.
type localControlPlane struct {
	Env            *envtest.Environment
	Client         client.Client
	Cfg            *rest.Config
	BinDir         string
	KubeconfigPath string
}

// Name returns the human-friendly name of the asset.
func (c *localControlPlane) Name() string {
	return "Local Control Plane"
}

// Run launches the local control plane.
func (c *localControlPlane) Run(clusterID *installconfig.ClusterID, installConfig *installconfig.InstallConfig) error {
	// Create a temporary directory to unpack the cluster-api binaries.
	binDir, err := os.MkdirTemp("", "openshift-cluster-api-bins")
	if err != nil {
		return err
	}
	c.BinDir = binDir
	if err := providers.UnpackClusterAPIBinary(binDir); err != nil {
		return err
	}
	if err := providers.UnpackEnvtestBinaries(binDir); err != nil {
		return err
	}

	log.SetLogger(klog.NewKlogr())
	logrus.Info("Started local control plane with envtest")
	c.Env = &envtest.Environment{
		Scheme:                   scheme.Scheme,
		AttachControlPlaneOutput: false,
		BinaryAssetsDirectory:    c.BinDir,
		ControlPlaneStartTimeout: 10 * time.Second,
		ControlPlaneStopTimeout:  10 * time.Second,
	}
	c.Cfg, err = c.Env.Start()
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Stop the controllers when the context is cancelled.
		<-ctx.Done()
		if err := c.Env.Stop(); err != nil {
			logrus.Warnf("failed to stop envtest: %v", err)
			return
		}
		logrus.Info("Stopped local control plane with envtest")
	}()

	kc := fromEnvTestConfig(c.Cfg)
	{
		dir := filepath.Join(command.RootOpts.Dir, "auth")
		kf, err := os.Create(filepath.Join(dir, "envtest.kubeconfig"))
		if err != nil {
			return err
		}
		if _, err := kf.Write(kc); err != nil {
			return err
		}
		if err := kf.Close(); err != nil {
			return err
		}
		c.KubeconfigPath, err = filepath.Abs(kf.Name())
		if err != nil {
			return err
		}
	}

	// Create a new client to interact with the cluster.
	cl, err := client.New(c.Cfg, client.Options{
		Scheme: c.Env.Scheme,
	})
	if err != nil {
		return err
	}
	c.Client = cl

	logrus.Infof("Stored kubeconfig for envtest in: %v", c.KubeconfigPath)
	return nil
}

// fromEnvTestConfig returns a new Kubeconfig in byte form when running in envtest.
func fromEnvTestConfig(cfg *rest.Config) []byte {
	clusterName := "envtest"
	contextName := fmt.Sprintf("%s@%s", cfg.Username, clusterName)
	c := api.Config{
		Clusters: map[string]*api.Cluster{
			clusterName: {
				Server:                   cfg.Host,
				CertificateAuthorityData: cfg.CAData,
			},
		},
		Contexts: map[string]*api.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: cfg.Username,
			},
		},
		AuthInfos: map[string]*api.AuthInfo{
			cfg.Username: {
				ClientKeyData:         cfg.KeyData,
				ClientCertificateData: cfg.CertData,
			},
		},
		CurrentContext: contextName,
	}
	data, err := clientcmd.Write(c)
	if err != nil {
		logrus.Fatalf("failed to write kubeconfig: %v", err)
	}
	return data
}
