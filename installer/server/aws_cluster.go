package server

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/session"
	bootkube "github.com/kubernetes-incubator/bootkube/pkg/asset"
	bootkubeTLS "github.com/kubernetes-incubator/bootkube/pkg/tlsutil"
	"golang.org/x/net/context"

	"github.com/coreos/tectonic-installer/installer/binassets"
	"github.com/coreos/tectonic-installer/installer/server/asset"
	"github.com/coreos/tectonic-installer/installer/server/aws/cloudforms"
)

const (
	awsCloudFormation = "cloud-formation.json"
)

var (
	errMissingAccessKeyID     = errors.New("installer: Missing AWS AccessKeyID")
	errMissingSecretAccessKey = errors.New("installer: Missing AWS SecretAccessKey")
	errNoClusterCreated       = errors.New("installer: No cluster created")

	controllerTmpl    = template.Must(template.New("").Parse(string(binassets.MustAsset("provisioning/tectonic-aws-controller.yaml.tmpl"))))
	workerTmpl        = template.Must(template.New("").Parse(string(binassets.MustAsset("provisioning/tectonic-aws-worker.yaml.tmpl"))))
	etcdTmpl          = template.Must(template.New("").Parse(string(binassets.MustAsset("provisioning/tectonic-aws-etcd.yaml.tmpl"))))
	stackTemplateTmpl = template.Must(template.New("").Parse(string(binassets.MustAsset("provisioning/stack-template.json.tmpl"))))
)

// TectonicAWSCluster provisions a Tectonic self-hosted Kubernetes cluster on
// AWS EC2.
type TectonicAWSCluster struct {
	CloudForm *cloudforms.Config `json:"cloudForm"`

	// AWS api credential
	AccessKeyID     string `json:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey"`
	SessionToken    string `json:"sessionToken"`

	// Custom Certificate Authority (optional)
	CACertificate string `json:"caCertificate"`
	CAPrivateKey  string `json:"caPrivateKey"`

	// Tectonic
	Tectonic *TectonicConfig `json:"tectonic"`

	// AWS client session
	awsSession *session.Session

	// Generated assets
	cluster *cloudforms.Cluster
}

// Initialize validates cluster data and sets defaults.
func (c *TectonicAWSCluster) Initialize() error {
	if c.AccessKeyID == "" {
		return errMissingAccessKeyID
	}
	if c.SecretAccessKey == "" {
		return errMissingSecretAccessKey
	}

	var err error
	c.awsSession, err = getAWSSession(c.AccessKeyID, c.SecretAccessKey, c.SessionToken, c.CloudForm.Region)
	if err != nil {
		return err
	}

	// set defaults and populate computed values
	c.CloudForm.SetDefaults()
	err = c.CloudForm.SetComputed(c.awsSession)
	if err != nil {
		return err
	}

	// Tectonic asset pipeline requires both domains
	c.Tectonic.ControllerDomain = c.CloudForm.ControllerDomain
	c.Tectonic.TectonicDomain = c.CloudForm.TectonicDomain
	if err := c.CloudForm.Valid(); err != nil {
		return err
	}
	return c.Tectonic.AssertValid()
}

// GenerateAssets generates cluster provisioning assets.
func (c *TectonicAWSCluster) GenerateAssets() ([]asset.Asset, error) {
	config, err := c.getBootkubeConfig(c.CloudForm)
	if err != nil {
		return nil, err
	}

	// Self-hosted (kube-system) manifests
	assets, err := NewBootkubeAssets(config)
	if err != nil {
		return nil, err
	}

	// Tectonic (tectonic-system) manifests and add-ons
	m := metrics{
		certificatesStrategy:   getCertificatesStrategy(c.CACertificate),
		installerPlatform:      c.Kind(),
		tectonicUpdaterEnabled: c.Tectonic.Updater.Enabled,
	}
	assets, err = NewTectonicAssets(assets, c.Tectonic, m)
	if err != nil {
		return nil, err
	}

	c.CloudForm.ControllerTemplate = controllerTmpl
	c.CloudForm.WorkerTemplate = workerTmpl
	c.CloudForm.EtcdTemplate = etcdTmpl
	c.CloudForm.StackTemplate = stackTemplateTmpl

	// decompose kubeconfig credentials, AWS KMS can't handle entire kubeconfig
	secretAssets, err := kubeconfigToSecretAssets(assets)
	if err != nil {
		return nil, fmt.Errorf("could not create AWS session: %v", err)
	}

	// Create location for assets to be stored in S3
	s3bucket := cloudforms.NewAwsBucket(c.awsSession, c.CloudForm.HostedZoneName)
	assetsFileName := c.CloudForm.ClusterName + "/assets.zip"

	// Save location of assets before rendering cloudformation ignition templates
	c.CloudForm.AssetsS3Bucket = s3bucket.Bucket()
	c.CloudForm.AssetsS3File = assetsFileName

	// Render the cloud-formation 'manifest'
	c.cluster, err = cloudforms.NewCloudFormation(c.CloudForm, c.awsSession, secretAssets)
	if err != nil {
		return nil, err
	}
	assets = append(assets, asset.New(awsCloudFormation, []byte(c.cluster.StackBody)))

	// Render Terraform's tfvars.
	terraformAssets, err := newAWSTerraformVars(c)
	if err != nil {
		return nil, err
	}
	assets = append(assets, terraformAssets...)

	// Zip assets to place in s3
	buf, err := asset.ZipAssets(assets)
	if err != nil {
		return nil, fmt.Errorf("Failed to zip assets: %v", err)
	}

	// Upload assets to S3
	err = s3bucket.Upload(assetsFileName, buf)
	if err != nil {
		return nil, fmt.Errorf("Failed to upload assets: %v", err)
	}

	return assets, nil
}

// StatusChecker returns a StatusChecker for Tectonic AWS clusters.
func (c *TectonicAWSCluster) StatusChecker() (StatusChecker, error) {
	if c.cluster == nil {
		return nil, errNoClusterCreated
	}
	return TectonicAWSChecker{
		AccessKeyID:      c.AccessKeyID,
		SecretAccessKey:  c.SecretAccessKey,
		SessionToken:     c.SessionToken,
		Region:           c.CloudForm.Region,
		ControllerDomain: c.CloudForm.ControllerDomain,
		TectonicDomain:   c.CloudForm.TectonicDomain,
		// TODO(dghubbe): cluster is big, only need cluster name, formalize this
		Cluster: &cloudforms.Cluster{
			ClusterName: c.cluster.ClusterName,
		},
	}, nil
}

// Kind returns the kind name of a cluster.
func (c *TectonicAWSCluster) Kind() string {
	return "tectonic-aws"
}

// Publish pushes the Cloud Formation template to AWS.
func (c *TectonicAWSCluster) Publish(ctx context.Context) error {
	if c.cluster == nil {
		return fmt.Errorf("error finding cluster to publish")
	}

	if _, err := c.cluster.Deploy(c.awsSession, c.CloudForm.Tags); err != nil {
		return fmt.Errorf("Error deploying cluster: %v", err)
	}
	return nil
}

// getBootkubeConfig returns a bootkube Config for asset generation.
func (c *TectonicAWSCluster) getBootkubeConfig(cluster *cloudforms.Config) (BootkubeConfig, error) {
	// determine the etcd endpoint
	var e string
	if cluster.ExternalETCDClient == "" {
		if len(cluster.ETCDInstances) == 0 {
			return BootkubeConfig{}, errors.New("no etcd instance available for bootkube")
		}
		e = fmt.Sprintf("http://%s:2379", cluster.ETCDInstances[0].DomainName)
	} else {
		e = cluster.ExternalETCDClient
	}

	// verify that the etcd endpoint is a valid URL
	etcd, err := url.Parse(e)
	if err != nil {
		return BootkubeConfig{}, err
	}

	// kube-apiserver net.URLs (e.g. https://external.domain.com:443)
	controllerURL, err := url.Parse(fmt.Sprintf("https://%s:443", cluster.ControllerDomain))
	if err != nil {
		return BootkubeConfig{}, err
	}

	_, podNet, err := net.ParseCIDR(cluster.PodCIDR)
	if err != nil {
		return BootkubeConfig{}, err
	}

	_, serviceNet, err := net.ParseCIDR(cluster.ServiceCIDR)
	if err != nil {
		return BootkubeConfig{}, err
	}

	// (optional) user-provided certificate authority
	caCert, err := getCACertificate(c.CACertificate)
	if err != nil {
		return BootkubeConfig{}, err
	}
	caPrivateKey, err := getCAPrivateKey(c.CAPrivateKey)
	if err != nil {
		return BootkubeConfig{}, err
	}

	// APIServers - bootkube uses the first endpoint to template --api-servers
	// flags in generated manifests
	// AltNames - bootkube generates TLS server certs for the temporary and
	// self-hosted apiserver. For both, we use external DNS.
	return BootkubeConfig{
		Config: bootkube.Config{
			EtcdServers: []*url.URL{etcd},
			CACert:      caCert,
			CAPrivKey:   caPrivateKey,
			APIServers:  []*url.URL{controllerURL},
			AltNames: &bootkubeTLS.AltNames{
				DNSNames: []string{cluster.ControllerDomain},
			},
			PodCIDR:       podNet,
			ServiceCIDR:   serviceNet,
			APIServiceIP:  cluster.APIServiceIP,
			DNSServiceIP:  cluster.DNSServiceIP,
			CloudProvider: "aws",
		},
		OIDCIssuer: &OIDCIssuer{
			IssuerURL:     fmt.Sprintf("https://%s/identity", cluster.TectonicDomain),
			ClientID:      oidcKubectlClientID,
			UsernameClaim: oidcUsernameClaim,
			// CACert is written as a secret and mounted by the kube-apiserver
			CAPath: "/etc/kubernetes/secrets/ca.crt",
		},
	}, nil
}
