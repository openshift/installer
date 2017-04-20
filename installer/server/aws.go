package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	bootkube "github.com/kubernetes-incubator/bootkube/pkg/asset"

	"github.com/coreos/tectonic-installer/installer/server/asset"
	"github.com/coreos/tectonic-installer/installer/server/aws/cloudforms"
)

// getAWSSession returns an AWS client session which should be re-used across
// service clients and cached. It is safe for concurrent use, but not
// modification.
func getAWSSession(accessKeyID, secretAccessKey, sessionToken, region string) (*session.Session, error) {
	// create an AWS client Config
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken)
	awsConfig := aws.NewConfig().
		WithCredentials(creds).
		WithRegion(region).
		WithCredentialsChainVerboseErrors(true)
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// kubeconfigToSecretAssets reads bootkube kubeconfig credential assets and
// returns SecretAssets which can be encrypted by Amazon KMS.
func kubeconfigToSecretAssets(assets []asset.Asset) (*cloudforms.SecretAssets, error) {
	caCert, err := asset.Find(assets, bootkube.AssetPathCACert)
	if err != nil {
		return nil, err
	}

	clientCert, err := asset.Find(assets, bootkube.AssetPathKubeletCert)
	if err != nil {
		return nil, err
	}

	clientKey, err := asset.Find(assets, bootkube.AssetPathKubeletKey)
	if err != nil {
		return nil, err
	}

	return &cloudforms.SecretAssets{
		CACert:     caCert.Data(),
		ClientCert: clientCert.Data(),
		ClientKey:  clientKey.Data(),
	}, nil
}
