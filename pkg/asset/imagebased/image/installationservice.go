package image

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/openshift/installer/pkg/asset"
)

const (
	installationService = `[Unit]
Wants=network-online.target
After=network-online.target
Description=SNO Image-based Installation
[Service]
Environment=SEED_IMAGE={{.SeedImage}}
Environment=HTTP_PROXY={{or .HTTPProxy ""}}
Environment=http_proxy={{or .HTTPProxy ""}}
Environment=HTTPS_PROXY={{or .HTTPSProxy ""}}
Environment=https_proxy={{or .HTTPSProxy ""}}
Environment=NO_PROXY={{or .NoProxy ""}}
Environment=no_proxy={{or .NoProxy ""}}
Environment=IBI_CONFIGURATION_FILE={{.IBIConfigurationPath}}
Environment=PULL_SECRET_FILE={{.PullSecretPath}}
Type=oneshot
RemainAfterExit=yes
ExecStartPre=/usr/bin/chcon -t install_exec_t {{.InstallationScriptPath}}
ExecStart={{.InstallationScriptPath}}
[Install]
WantedBy=multi-user.target
`
)

type installServiceTemplate struct {
	SeedImage              string
	IBIConfigurationPath   string
	InstallationScriptPath string
	HTTPProxy              string
	HTTPSProxy             string
	NoProxy                string
	PullSecretPath         string
}

// InstallationService is the systemd service that executes the image-based
// installation during boot.
type InstallationService struct {
	Content string
}

// Dependencies returns the assets on which the InstallationService asset depends.
func (i *InstallationService) Dependencies() []asset.Asset {
	return []asset.Asset{
		&ImageBasedInstallationConfig{},
	}
}

// Generate generates the installation systemd service for the image-based installation ISO.
func (i *InstallationService) Generate(dependencies asset.Parents) error {
	ibiConfig := &ImageBasedInstallationConfig{}
	dependencies.Get(ibiConfig)

	templateData := installServiceTemplate{
		SeedImage:              ibiConfig.Config.SeedImage,
		PullSecretPath:         pullSecretPath,
		IBIConfigurationPath:   ibiConfigurationPath,
		InstallationScriptPath: installationScriptPath,
	}

	if ibiConfig.Config.Proxy != nil {
		templateData.HTTPProxy = ibiConfig.Config.Proxy.HTTPProxy
		templateData.HTTPSProxy = ibiConfig.Config.Proxy.HTTPSProxy
		templateData.NoProxy = ibiConfig.Config.Proxy.NoProxy
	}

	installationServiceContent, err := renderTemplate(string(installationService), templateData)
	if err != nil {
		return fmt.Errorf("failed to render installation service: %w", err)
	}

	i.Content = string(installationServiceContent)

	return nil
}

// Name returns the human-friendly name of the asset.
func (i *InstallationService) Name() string {
	return "ImageBased installation systemd service"
}

func renderTemplate(templateData string, params any) ([]byte, error) {
	tmpl := template.New("template")
	tmpl = template.Must(tmpl.Parse(templateData))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}
	return buf.Bytes(), nil
}
