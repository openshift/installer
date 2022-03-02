package openstack

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/yaml.v2"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/certificates"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
)

const (
	rsaPrivateKeyBlockType      = "RSA PRIVATE KEY"
	certificateRequestBlockType = "CERTIFICATE REQUEST"
)

func expandContainerInfraV1LabelsMap(v map[string]interface{}) (map[string]string, error) {
	m := make(map[string]string)
	for key, val := range v {
		labelValue, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf("label %s value should be string", key)
		}
		m[key] = labelValue
	}
	return m, nil
}

func expandContainerInfraV1LabelsString(v map[string]interface{}) (string, error) {
	var formattedLabels string
	for key, val := range v {
		labelValue, ok := val.(string)
		if !ok {
			return "", fmt.Errorf("label %s value should be string", key)
		}
		formattedLabels = strings.Join([]string{
			formattedLabels,
			fmt.Sprintf("%s=%s", key, labelValue),
		}, ",")
	}
	formattedLabels = strings.Trim(formattedLabels, ",")

	return formattedLabels, nil
}

func containerInfraClusterTemplateV1AppendUpdateOpts(updateOpts []clustertemplates.UpdateOptsBuilder, attribute, value string) []clustertemplates.UpdateOptsBuilder {
	if value == "" {
		updateOpts = append(updateOpts, clustertemplates.UpdateOpts{
			Op:   clustertemplates.RemoveOp,
			Path: strings.Join([]string{"/", attribute}, ""),
		})
	} else {
		updateOpts = append(updateOpts, clustertemplates.UpdateOpts{
			Op:    clustertemplates.ReplaceOp,
			Path:  strings.Join([]string{"/", attribute}, ""),
			Value: value,
		})
	}
	return updateOpts
}

// ContainerInfraClusterV1StateRefreshFunc returns a resource.StateRefreshFunc
// that is used to watch a container infra Cluster.
func containerInfraClusterV1StateRefreshFunc(client *gophercloud.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		c, err := clusters.Get(client, clusterID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return c, "DELETE_COMPLETE", nil
			}
			return nil, "", err
		}

		errorStatuses := []string{
			"CREATE_FAILED",
			"UPDATE_FAILED",
			"DELETE_FAILED",
			"RESUME_FAILED",
			"ROLLBACK_FAILED",
		}
		for _, errorStatus := range errorStatuses {
			if c.Status == errorStatus {
				err = fmt.Errorf("openstack_containerinfra_cluster_v1 is in an error state: %s", c.StatusReason)
				return c, c.Status, err
			}
		}

		return c, c.Status, nil
	}
}

// containerInfraClusterV1Flavor will determine the flavor for a container infra
// cluster based on either what was set in the configuration or environment
// variable.
func containerInfraClusterV1Flavor(d *schema.ResourceData) (string, error) {
	if flavor := d.Get("flavor").(string); flavor != "" {
		return flavor, nil
	}
	// Try the OS_MAGNUM_FLAVOR environment variable
	if v := os.Getenv("OS_MAGNUM_FLAVOR"); v != "" {
		return v, nil
	}

	return "", nil
}

// containerInfraClusterV1Flavor will determine the master flavor for a
// container infra cluster based on either what was set in the configuration
// or environment variable.
func containerInfraClusterV1MasterFlavor(d *schema.ResourceData) (string, error) {
	if flavor := d.Get("master_flavor").(string); flavor != "" {
		return flavor, nil
	}

	// Try the OS_MAGNUM_MASTER_FLAVOR environment variable
	if v := os.Getenv("OS_MAGNUM_MASTER_FLAVOR"); v != "" {
		return v, nil
	}

	return "", nil
}

type kubernetesConfig struct {
	APIVersion     string                    `yaml:"apiVersion"`
	Kind           string                    `yaml:"kind"`
	Clusters       []kubernetesConfigCluster `yaml:"clusters"`
	Contexts       []kubernetesConfigContext `yaml:"contexts"`
	CurrentContext string                    `yaml:"current-context"`
	Users          []kubernetesConfigUser    `yaml:"users"`
}

type kubernetesConfigCluster struct {
	Cluster kubernetesConfigClusterData `yaml:"cluster"`
	Name    string                      `yaml:"name"`
}
type kubernetesConfigClusterData struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type kubernetesConfigContext struct {
	Context kubernetesConfigContextData `yaml:"context"`
	Name    string                      `yaml:"name"`
}
type kubernetesConfigContextData struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type kubernetesConfigUser struct {
	Name string                   `yaml:"name"`
	User kubernetesConfigUserData `yaml:"user"`
}

type kubernetesConfigUserData struct {
	ClientKeyData         string `yaml:"client-key-data"`
	ClientCertificateData string `yaml:"client-certificate-data"`
}

func flattenContainerInfraV1Kubeconfig(d *schema.ResourceData, containerInfraClient *gophercloud.ServiceClient) (map[string]string, error) {
	clientSert, ok := d.Get("kubeconfig.client_certificate").(string)
	if ok && clientSert != "" {
		return d.Get("kubeconfig").(map[string]string), nil
	}

	certificateAuthority, err := certificates.Get(containerInfraClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error getting certificate authority: %s", err)
	}

	clientKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("Error generating client key: %s", err)
	}

	csrTemplate := x509.CertificateRequest{
		PublicKey:          clientKey.Public,
		SignatureAlgorithm: x509.SHA512WithRSA,
		Subject: pkix.Name{
			CommonName:         "admin",
			Organization:       []string{"system:masters"},
			OrganizationalUnit: []string{"terraform"},
		},
	}

	clientCsr, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, clientKey)
	if err != nil {
		return nil, fmt.Errorf("Error generating client CSR: %s", err)
	}

	pemClientKey := pem.EncodeToMemory(
		&pem.Block{
			Type:  rsaPrivateKeyBlockType,
			Bytes: x509.MarshalPKCS1PrivateKey(clientKey),
		},
	)

	pemClientCsr := pem.EncodeToMemory(
		&pem.Block{
			Type:  certificateRequestBlockType,
			Bytes: clientCsr,
		},
	)

	certificateCreateOpts := certificates.CreateOpts{
		ClusterUUID: d.Id(),
		CSR:         string(pemClientCsr),
	}

	clientCertificate, err := certificates.Create(containerInfraClient, certificateCreateOpts).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error requesting client certificate: %s", err)
	}

	name := d.Get("name").(string)
	host := d.Get("api_address").(string)
	rawKubeconfig, err := renderKubeconfig(name, host, []byte(certificateAuthority.PEM), []byte(clientCertificate.PEM), pemClientKey)
	if err != nil {
		return nil, fmt.Errorf("Error rendering kubeconfig: %s", err)
	}

	return map[string]string{
		"raw_config":             string(rawKubeconfig),
		"host":                   host,
		"cluster_ca_certificate": certificateAuthority.PEM,
		"client_certificate":     clientCertificate.PEM,
		"client_key":             string(pemClientKey),
	}, nil
}

func renderKubeconfig(name string, host string, clusterCaCertificate []byte, clientCertificate []byte, clientKey []byte) ([]byte, error) {
	userName := fmt.Sprintf("%s-admin", name)

	config := kubernetesConfig{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []kubernetesConfigCluster{
			{
				Name: name,
				Cluster: kubernetesConfigClusterData{
					CertificateAuthorityData: base64.StdEncoding.EncodeToString(clusterCaCertificate),
					Server:                   host,
				},
			},
		},
		Contexts: []kubernetesConfigContext{
			{
				Context: kubernetesConfigContextData{
					Cluster: name,
					User:    userName,
				},
				Name: name,
			},
		},
		CurrentContext: name,
		Users: []kubernetesConfigUser{
			{
				Name: userName,
				User: kubernetesConfigUserData{
					ClientCertificateData: base64.StdEncoding.EncodeToString(clientCertificate),
					ClientKeyData:         base64.StdEncoding.EncodeToString(clientKey),
				},
			},
		},
	}

	return yaml.Marshal(config)
}
