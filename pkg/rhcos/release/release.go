// Package release extracts the RHCOS build version from an
// OpenShift release image.
package release

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"

	"github.com/containers/image/docker"
	"github.com/containers/image/image"
	"github.com/containers/image/pkg/blobinfocache"
	"github.com/containers/image/signature"
	"github.com/containers/image/transports"
	"github.com/containers/image/types"
	"github.com/docker/distribution/reference"
	imagev1 "github.com/openshift/api/image/v1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// RHCOSBuild extracts the RHCOS build version from an OpenShift
// release image.
func RHCOSBuild(ctx context.Context, releaseImagePullSpec string, pullSecret []byte) (string, error) {
	imageReferencesBytes, err := imageReferences(ctx, releaseImagePullSpec, pullSecret)
	if err != nil {
		return "", err
	}

	pullSpec, err := osContentPullSpec(imageReferencesBytes)
	if err != nil {
		return "", err
	}

	imageSource, image, err := pull(ctx, pullSpec, pullSecret)
	if err != nil {
		return "", err
	}
	defer imageSource.Close()

	inspect, err := image.Inspect(ctx)
	if err != nil {
		return "", err
	}

	version, ok := inspect.Labels["version"]
	if !ok {
		return "", errors.Errorf("%s has no version annotation", pullSpec)
	}

	return version, nil
}

// pull pulls an image.  If the pull secret includes credentials for
// the pull-spec authority, those are used to authenticate requests.
// The signature on the returned image, if any, is validated agaist
// the signature policy from /etc/containers/policy.json.
func pull(ctx context.Context, pullSpec string, pullSecret []byte) (types.ImageSource, types.Image, error) {
	named, err := reference.ParseNamed(pullSpec)
	if err != nil {
		return nil, nil, err
	}

	reference, err := docker.NewReference(named)
	if err != nil {
		return nil, nil, err
	}

	sys := &types.SystemContext{}
	err = addPullSecret(sys, pullSecret, named)
	if err != nil {
		return nil, nil, err
	}

	imageSource, err := reference.NewImageSource(ctx, sys)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			imageSource.Close()
		}
	}()

	unparsed := image.UnparsedInstance(imageSource, nil)

	policy, err := signature.DefaultPolicy(sys)
	if err != nil {
		return nil, nil, err
	}

	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return nil, nil, err
	}

	if allowed, err := policyContext.IsRunningImageAllowed(ctx, unparsed); !allowed || err != nil {
		return nil, nil, errors.Wrapf(err, "%s rejected", pullSpec)
	}

	img, err := image.FromUnparsedImage(ctx, sys, unparsed)
	return imageSource, img, err
}

func imageReferences(ctx context.Context, releaseImagePullSpec string, pullSecret []byte) ([]byte, error) {
	imageSource, image, err := pull(ctx, releaseImagePullSpec, pullSecret)
	if err != nil {
		return nil, err
	}
	defer imageSource.Close()

	layers := image.LayerInfos()
	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]
		blob, _, err := imageSource.GetBlob(ctx, layer, blobinfocache.NoCache)
		if err != nil {
			return nil, err
		}
		defer blob.Close()

		// fmt.Println(layer.MediaType) // TODO: blank from quay?  Just assume this is a gzipped tar?
		gzipReader, err := gzip.NewReader(blob)
		if err != nil {
			return nil, err
		}

		tarReader := tar.NewReader(gzipReader)
		for {
			hdr, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			if hdr.Name == "release-manifests/image-references" {
				return ioutil.ReadAll(tarReader)
			}
		}
	}

	return nil, errors.Errorf("no image reference found in %s", transports.ImageName(image.Reference()))
}

func osContentPullSpec(imageReferences []byte) (string, error) {
	scheme := runtime.NewScheme()
	imagev1.Install(scheme)
	decoder := serializer.NewCodecFactory(scheme).UniversalDecoder(
		imagev1.GroupVersion,
	)
	obj, _, err := decoder.Decode(imageReferences, nil, nil)
	if err != nil {
		return "", err
	}

	imageStream, ok := obj.(*imagev1.ImageStream)
	if !ok {
		return "", errors.Errorf("image references is a %q, not an image stream", obj.GetObjectKind())
	}

	for _, tag := range imageStream.Spec.Tags {
		if tag.Name == "machine-os-content" {
			if tag.From == nil {
				return "", errors.New("machine-os-content tag has an empty 'from' property")
			}

			if tag.From.Kind != "DockerImage" {
				return "", errors.Errorf("unrecognized machine-os-content tag from kind %q", tag.From.Kind)
			}

			return tag.From.Name, nil
		}
	}

	return "", errors.New("no machine-os-content tag found")
}
