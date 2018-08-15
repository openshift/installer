package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

func cloudConfigRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	ssn := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	creds, err := ssn.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	return installerassets.TemplateRebuilder(
		"files/opt/tectonic/tectonic/aws/99_cloud-creds-secret.yaml",
		nil,
		map[string]interface{}{
			"AccessKeyID":     creds.AccessKeyID,
			"SecretAccessKey": creds.SecretAccessKey,
		},
	)(ctx, getByName)
}

func init() {
	installerassets.Rebuilders["files/opt/tectonic/tectonic/aws/99_cloud-creds-secret.yaml"] = cloudConfigRebuilder
}
