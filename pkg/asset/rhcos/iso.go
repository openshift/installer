package rhcos

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"time"

	"github.com/coreos/stream-metadata-go/arch"
	"github.com/coreos/stream-metadata-go/stream"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/agent/workflow"
	workflowreport "github.com/openshift/installer/pkg/asset/agent/workflow/report"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/rhcos/cache"
	"github.com/openshift/installer/pkg/types"
)

// BaseIso generates the base ISO file for the image.
type BaseIso struct {
	streamGetter CoreOSBuildFetcher
	ocRelease    ReleasePayload
}

// CoreOSBuildFetcher will be to used to switch the source of the coreos metadata.
type CoreOSBuildFetcher func(ctx context.Context) (*stream.Stream, error)

// defaultCoreOSStreamGetter uses the pinned metadata.
var defaultCoreOSStreamGetter = rhcos.FetchCoreOSBuild

// NewBaseISOFetcher returns a struct that can be used to fetch a base ISO using
// the default method.
func NewBaseISOFetcher(ocRelease ReleasePayload, streamGetter CoreOSBuildFetcher) *BaseIso {
	if streamGetter == nil {
		streamGetter = defaultCoreOSStreamGetter
	}
	return &BaseIso{
		streamGetter: streamGetter,
		ocRelease:    ocRelease,
	}
}

// GetBaseISOFilename retrieves the base ISO for the given architecture
// (possibly from the cache) and returns its location on disk.
func (i *BaseIso) GetBaseISOFilename(ctx context.Context, arch string) (baseIsoFileName string, err error) {
	err = workflowreport.GetReport(ctx).Stage(workflow.StageFetchBaseISO)
	if err != nil {
		return
	}

	if urlOverride, ok := os.LookupEnv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); ok && urlOverride != "" {
		logrus.Warn("Found override for OS Image. Please be warned, this is not advised")
		baseIsoFileName, err = cache.DownloadImageFile(urlOverride, cache.AgentApplicationName)
	} else {
		baseIsoFileName, err = i.retrieveBaseIso(ctx, arch)
	}

	return
}

// GetMetalArtifact returns the CoreOS artifacts for metal for a given arch
// from a given stream.
func GetMetalArtifact(ctx context.Context, archName string, streamGetter CoreOSBuildFetcher) (stream.PlatformArtifacts, error) {
	if streamGetter == nil {
		streamGetter = defaultCoreOSStreamGetter
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get the ISO to use from rhcos.json
	st, err := streamGetter(ctx)
	if err != nil {
		return stream.PlatformArtifacts{}, err
	}

	streamArch, err := st.GetArchitecture(archName)
	if err != nil {
		return stream.PlatformArtifacts{}, err
	}

	metal, ok := streamArch.Artifacts["metal"]
	if !ok {
		return stream.PlatformArtifacts{}, fmt.Errorf("coreOs stream data not found for 'metal' artifact")
	}

	return metal, nil
}

// Download the ISO using the URL in rhcos.json.
func (i *BaseIso) downloadIso(ctx context.Context, archName string) (string, error) {
	metal, err := GetMetalArtifact(ctx, archName, i.streamGetter)
	if err != nil {
		return "", err
	}

	format, ok := metal.Formats["iso"]
	if !ok {
		return "", fmt.Errorf("no ISO found to download for %s", archName)
	}

	url := format.Disk.Location
	sha := format.Disk.Sha256
	cachedImage, err := cache.DownloadImageFileWithSha(url, cache.AgentApplicationName, sha)
	if err != nil {
		return "", fmt.Errorf("failed to download base ISO image %s: %w", url, err)
	}

	return cachedImage, nil
}

func (i *BaseIso) checkReleasePayloadBaseISOVersion(ctx context.Context, r ReleasePayload, archName string) {
	logrus.Debugf("Checking release payload base ISO version")

	// Get current release payload CoreOS version
	payloadRelease, err := r.GetBaseIsoVersion(archName)
	if err != nil {
		logrus.Warnf("unable to determine base ISO version: %s", err.Error())
		return
	}

	// Get pinned version from installer
	metal, err := GetMetalArtifact(ctx, archName, i.streamGetter)
	if err != nil {
		logrus.Warnf("unable to determine base ISO version: %s", err.Error())
		return
	}

	// Check for a mismatch
	if metal.Release != payloadRelease {
		logrus.Warnf("base ISO version mismatch in release payload. Expected version %s but found %s", metal.Release, payloadRelease)
	}
}

func (i *BaseIso) retrieveBaseIso(ctx context.Context, archName string) (string, error) {
	// Default iso archName to x86_64.
	if archName == "" {
		archName = arch.RpmArch(types.ArchitectureAMD64)
	}

	if i.ocRelease != nil {
		// If we have the image registry location and 'oc' command is available then get from release payload
		logrus.Info("Extracting base ISO from release payload")

		if err := workflowreport.GetReport(ctx).SubStage(workflow.StageFetchBaseISOExtract); err != nil {
			return "", err
		}
		baseIsoFileName, err := i.ocRelease.GetBaseIso(archName, i.streamGetter)
		if err == nil {
			if err := workflowreport.GetReport(ctx).SubStage(workflow.StageFetchBaseISOVerify); err != nil {
				return "", err
			}
			i.checkReleasePayloadBaseISOVersion(ctx, i.ocRelease, archName)

			logrus.Debugf("Extracted base ISO image %s from release payload", baseIsoFileName)
			return baseIsoFileName, nil
		}

		if errors.Is(err, fs.ErrNotExist) {
			// if image extract failed to extract the iso that architecture may be missing from release image
			return "", fmt.Errorf("base ISO for %s not found in release image, check release image architecture", archName)
		}
		if !errors.Is(err, &exec.Error{}) { // Already warned about missing oc binary
			logrus.Warning("Failed to extract base ISO from release payload - check registry configuration")
		}
	}

	logrus.Info("Downloading base ISO")
	if err := workflowreport.GetReport(ctx).SubStage(workflow.StageFetchBaseISODownload); err != nil {
		return "", err
	}
	return i.downloadIso(ctx, archName)
}
