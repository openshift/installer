package image

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/rhcos"
)

// BaseIso generates the base ISO file for the image
type BaseIso struct {
	asset.DefaultFileWriter
}

var (
	baseIsoFilename = ""
)

var _ asset.WritableAsset = (*BaseIso)(nil)

// Name returns the human-friendly name of the asset.
func (i *BaseIso) Name() string {
	return "BaseIso Image"
}

// getIsoFile is a pluggable function that gets the base ISO file
type getIsoFile func() (string, error)

type getIso struct {
	getter getIsoFile
}

func newGetIso(getter getIsoFile) *getIso {
	return &getIso{getter: getter}
}

// GetIsoPluggable defines the method to use get the baseIso file
var GetIsoPluggable = downloadIso

// Download the ISO using the URL in rhcos.json
func downloadIso() (string, error) {

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	// Get the ISO to use from rhcos.json
	st, err := rhcos.FetchCoreOSBuild(ctx)
	if err != nil {
		return "", err
	}

	// Defaults to using the x86_64 baremetal ISO for all platforms
	// archName := arch.RpmArch(string(config.ControlPlane.Architecture))
	archName := "x86_64"
	streamArch, err := st.GetArchitecture(archName)
	if err != nil {
		return "", err
	}
	if artifacts, ok := streamArch.Artifacts["metal"]; ok {
		if format, ok := artifacts.Formats["iso"]; ok {
			url := format.Disk.Location

			cachedImage, err := DownloadImageFile(url)
			if err != nil {
				return "", errors.Wrapf(err, "failed to download base iso image %s", url)
			}
			return cachedImage, nil
		}
	} else {
		return "", errors.Wrap(err, "invalid artifact")
	}

	return "", fmt.Errorf("no iso found to download for %s", archName)
}

func getIsoFromReleasePayload() (string, error) {

	// TODO
	return "", nil
}

// Dependencies returns dependencies used by the asset.
func (i *BaseIso) Dependencies() []asset.Asset {
	return []asset.Asset{
		// TODO - will need to depend on installConfig for disconnected image registry
		// &installconfig.InstallConfig{},
	}
}

// Generate the baseIso
func (i *BaseIso) Generate(p asset.Parents) error {

	// TODO - if image registry location is defined in InstallConfig,
	// ic := &installconfig.InstallConfig{}
	// p.Get(ic)
	// use the GetIso function to get the BaseIso from the release payload
	isoGetter := newGetIso(GetIsoPluggable)
	baseIsoFileName, err := isoGetter.getter()
	if err != nil {
		return errors.Wrap(err, "failed to get base iso image")
	}

	i.File = &asset.File{Filename: baseIsoFileName}

	return nil
}

// Load returns the cached baseIso
func (i *BaseIso) Load(f asset.FileFetcher) (bool, error) {

	if baseIsoFilename == "" {
		return false, nil
	}

	baseIso, err := f.FetchByName(baseIsoFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", baseIsoFilename))
	}

	i.File = baseIso
	return true, nil
}
