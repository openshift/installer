package clusterapi

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	capnv1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
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
	capgv1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	capiv1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	capov1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	capvv1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/openshift/installer/cmd/openshift-install/command"
)

var (
	// Scheme is the scheme used by the local control plane.
	Scheme = scheme.Scheme
)

func init() {
	utilruntime.Must(clusterv1.AddToScheme(Scheme))
	utilruntime.Must(capav1beta1.AddToScheme(Scheme))
	utilruntime.Must(capav1.AddToScheme(Scheme))
	utilruntime.Must(capzv1.AddToScheme(Scheme))
	utilruntime.Must(capgv1.AddToScheme(Scheme))
	utilruntime.Must(capvv1.AddToScheme(Scheme))
	utilruntime.Must(capov1.AddToScheme(Scheme))
	utilruntime.Must(capiv1.AddToScheme(Scheme))
	utilruntime.Must(capnv1.AddToScheme(Scheme))
}

// localControlPlane creates a local capi control plane
// to use as a management cluster.
type localControlPlane struct {
	Env            *envtest.Environment
	Client         client.Client
	Cfg            *rest.Config
	BinDir         string
	EtcdDataDir    string
	KubeconfigPath string
	EtcdLog        *os.File
	APIServerLog   *os.File
	APIServerArgs  map[string]string
}

// Run launches the local control plane.
func (c *localControlPlane) Run(ctx context.Context) error {
	// Create a temporary directory to unpack the cluster-api binaries.
	c.BinDir = filepath.Join(command.RootOpts.Dir, "cluster-api")
	if err := UnpackClusterAPIBinary(c.BinDir); err != nil {
		return fmt.Errorf("failed to unpack cluster-api binary: %w", err)
	}
	if err := UnpackEnvtestBinaries(c.BinDir); err != nil {
		return fmt.Errorf("failed to unpack envtest binaries: %w", err)
	}
	c.EtcdDataDir = filepath.Join(command.RootOpts.Dir, ArtifactsDir, "etcd")

	// Write etcd & kube-apiserver output to log files.
	var err error
	if err := os.MkdirAll(filepath.Join(command.RootOpts.Dir, ArtifactsDir), 0750); err != nil {
		return fmt.Errorf("error creating artifacts dir: %w", err)
	}
	if c.EtcdLog, err = os.Create(filepath.Join(command.RootOpts.Dir, ArtifactsDir, "etcd.log")); err != nil {
		return fmt.Errorf("failed to create etcd log file: %w", err)
	}
	if c.APIServerLog, err = os.Create(filepath.Join(command.RootOpts.Dir, ArtifactsDir, "kube-apiserver.log")); err != nil {
		return fmt.Errorf("failed to create kube-apiserver log file: %w", err)
	}

	log.SetLogger(klog.NewKlogr())
	logrus.Info("Started local control plane with envtest")
	c.Env = &envtest.Environment{
		Scheme:                   Scheme,
		AttachControlPlaneOutput: true,
		BinaryAssetsDirectory:    c.BinDir,
		ControlPlaneStartTimeout: 10 * time.Second,
		ControlPlaneStopTimeout:  10 * time.Second,
		ControlPlane: envtest.ControlPlane{
			Etcd: &envtest.Etcd{
				DataDir: c.EtcdDataDir,
				Out:     c.EtcdLog,
				Err:     c.EtcdLog,
			},
		},
	}
	apiServer := &envtest.APIServer{
		Out: c.APIServerLog,
		Err: c.APIServerLog,
	}
	for key, value := range c.APIServerArgs {
		apiServer.Configure().Set(key, value)
	}
	c.Env.ControlPlane.APIServer = apiServer
	c.Cfg, err = c.Env.Start()
	if err != nil {
		return err
	}

	artifactsDirPath := filepath.Join(command.RootOpts.Dir, ArtifactsDir)
	err = os.MkdirAll(artifactsDirPath, 0750)
	if err != nil {
		return fmt.Errorf("error creating cluster-api artifacts directory: %w", err)
	}

	kc := fromEnvTestConfig(c.Cfg)
	{
		kf, err := os.Create(filepath.Join(artifactsDirPath, "envtest.kubeconfig"))
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

func (c *localControlPlane) Stop() error {
	return c.Env.Stop()
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
