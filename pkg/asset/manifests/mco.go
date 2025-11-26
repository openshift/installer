package manifests

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/gcp"
)

func generateMCOManifest(installConfig *types.InstallConfig, template []*asset.File) []*asset.File {
	_, customWImg := customBootImages(installConfig)

	// If there are no custom images, skip creating the manifest
	// to defer to the MCO's default behavior.
	if !customWImg {
		return nil
	}

	tmplData := mcoTemplateData{DisableMachinesetBootMgmt: customWImg}

	mcoCfg := applyTemplateData(template[0].Data, tmplData)
	return []*asset.File{
		{
			Filename: path.Join(manifestDir, strings.TrimSuffix(filepath.Base(template[0].Filename), ".template")),
			Data:     mcoCfg,
		},
	}
}

func customBootImages(ic *types.InstallConfig) (customCPImg, customWImg bool) {
	switch ic.Platform.Name() {
	case aws.Name:
		customCPImg, customWImg = awsBootImages(ic)
	case gcp.Name:
		customCPImg, customWImg = gcpBootImages(ic)
	default:
		// We do not need to consider other platforms, because default boot image management has not been enabled yet.
		return
	}
	return
}

func awsBootImages(ic *types.InstallConfig) (cpImg bool, wImg bool) {
	if dmp := ic.AWS.DefaultMachinePlatform; dmp != nil && dmp.AMIID != "" {
		return true, true
	}

	if cp := ic.ControlPlane; cp != nil && cp.Platform.AWS != nil && cp.Platform.AWS.AMIID != "" {
		cpImg = true
	}

	// On AWS, we need to check both compute and edge compute machine pool.
	for _, computeMP := range ic.Compute {
		if awsPlatform := computeMP.Platform.AWS; awsPlatform != nil && awsPlatform.AMIID != "" {
			wImg = true
		}
	}
	return
}

func gcpBootImages(ic *types.InstallConfig) (cpImg bool, wImg bool) {
	if dmp := ic.GCP.DefaultMachinePlatform; dmp != nil && dmp.OSImage != nil {
		return true, true
	}

	if cp := ic.ControlPlane; cp != nil && cp.Platform.GCP != nil && cp.Platform.GCP.OSImage != nil {
		cpImg = true
	}

	if w := ic.Compute; len(w) > 0 && w[0].Platform.GCP != nil && w[0].Platform.GCP.OSImage != nil {
		wImg = true
	}
	return
}
