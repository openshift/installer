package strategy

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	operatorv1alpha1scheme "github.com/openshift/client-go/operator/clientset/versioned/scheme"
	"github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/library-go/pkg/image/registryclient"
)

type onErrorStrategy struct {
	lock sync.Mutex

	alternates            map[reference.DockerImageReference][]reference.DockerImageReference
	icspFile              string
	readICSPsFromFileFunc readICSPsFromFileFunc
}

var _ registryclient.AlternateBlobSourceStrategy = &onErrorStrategy{}

// NewICSPOnErrorStrategy returns ICSP alternate strategy which reads alternate
// sources only after getting an error from the original requested.
func NewICSPOnErrorStrategy(file string) registryclient.AlternateBlobSourceStrategy {
	return &onErrorStrategy{
		icspFile:              file,
		alternates:            make(map[reference.DockerImageReference][]reference.DockerImageReference),
		readICSPsFromFileFunc: readICSPsFromFile,
	}
}

func (s *onErrorStrategy) FirstRequest(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
	return nil, nil
}

func (s *onErrorStrategy) OnFailure(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if alternates, ok := s.alternates[locator]; ok {
		return alternates, nil
	}
	alternates, err := s.resolve(ctx, locator)
	if err != nil {
		return nil, err
	}
	if len(alternates) == 0 {
		return nil, fmt.Errorf("no alternative image references found for image: %s", locator.String())
	}
	s.alternates[locator] = alternates
	return s.alternates[locator], nil
}

// resolve gathers possible image sources for a given image
// gathered from ImageContentSourcePolicy file.
// Image reference of user-given image may be different from original in case of mirrored images.
func (s *onErrorStrategy) resolve(ctx context.Context, imageRef reference.DockerImageReference) ([]reference.DockerImageReference, error) {
	if len(s.icspFile) == 0 {
		return nil, fmt.Errorf("no ImageContentSourceFile specified")
	}
	klog.V(5).Infof("Reading ICSP from file %s", s.icspFile)
	icspList, err := s.readICSPsFromFileFunc(s.icspFile)
	if err != nil {
		return nil, err
	}
	// always add the original as the first reference
	imageRefList, err := alternativeImageSources(imageRef, icspList, false)
	if err != nil {
		return nil, err
	}
	return imageRefList, nil
}

// alternativeImageSources returns unique list of DockerImageReference objects from list of ImageContentSourcePolicy objects
// addSourceAsLastAlternate decides whether the original imageRef is first or the last element in the result
func alternativeImageSources(imageRef reference.DockerImageReference, icspList []operatorv1alpha1.ImageContentSourcePolicy, addSourceAsLastAlternate bool) ([]reference.DockerImageReference, error) {
	var imageSources []reference.DockerImageReference
	klog.V(5).Infof("%v ImageReference added to potential ImageSourcePrefixes from ImageContentSourcePolicy", imageRef.AsRepository().AsV2())
	if !addSourceAsLastAlternate {
		imageSources = append(imageSources, imageRef.AsRepository().AsV2())
	}
	for _, icsp := range icspList {
		repoDigestMirrors := icsp.Spec.RepositoryDigestMirrors
		for _, rdm := range repoDigestMirrors {
			var err error
			rdmSourceRef, err := reference.Parse(rdm.Source)
			if err != nil {
				return nil, fmt.Errorf("invalid source %q: %w", rdm.Source, err)
			}
			// AsV2 in the right call is required to ensure we transform docker registry
			// from docker.io to registry-1.docker.io
			if imageRef.AsRepository().AsV2() != rdmSourceRef.AsRepository().AsV2() {
				continue
			}
			klog.V(5).Infof("%v RepositoryDigestMirrors source matches given image", imageRef.AsRepository().AsV2())
			for _, m := range rdm.Mirrors {
				mRef, err := reference.Parse(m)
				if err != nil {
					return nil, fmt.Errorf("invalid mirror %q: %w", m, err)
				}
				imageSources = append(imageSources, mRef)
				klog.V(5).Infof("%v RepositoryDigestMirrors mirror added to potential ImageSourcePrefixes from ImageContentSourcePolicy", m)
			}
		}
	}
	if addSourceAsLastAlternate {
		imageSources = append(imageSources, imageRef.AsRepository().AsV2())
	}
	uniqueMirrors := make([]reference.DockerImageReference, 0, len(imageSources))
	uniqueMap := make(map[reference.DockerImageReference]bool)
	for _, imageSourceMirror := range imageSources {
		if _, ok := uniqueMap[imageSourceMirror]; !ok {
			uniqueMap[imageSourceMirror] = true
			uniqueMirrors = append(uniqueMirrors, imageSourceMirror)
		}
	}
	klog.V(2).Infof("Found sources: %v for image: %v", uniqueMirrors, imageRef)
	return uniqueMirrors, nil
}

// readICSPsFromFileFunc is used for testing to be able to inject ICSP data
type readICSPsFromFileFunc func(string) ([]operatorv1alpha1.ImageContentSourcePolicy, error)

// readICSPsFromFile appends to list of alternative image sources from ICSP file
// returns error if no icsp object decoded from file data
func readICSPsFromFile(icspFile string) ([]operatorv1alpha1.ImageContentSourcePolicy, error) {
	icspData, err := ioutil.ReadFile(icspFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read ImageContentSourceFile %s: %v", icspFile, err)
	}
	if len(icspData) == 0 {
		return nil, fmt.Errorf("no data found in ImageContentSourceFile %s", icspFile)
	}
	icspObj, err := runtime.Decode(operatorv1alpha1scheme.Codecs.UniversalDeserializer(), icspData)
	if err != nil {
		return nil, fmt.Errorf("error decoding ImageContentSourcePolicy from %s: %v", icspFile, err)
	}
	icsp, ok := icspObj.(*operatorv1alpha1.ImageContentSourcePolicy)
	if !ok {
		return nil, fmt.Errorf("could not decode ImageContentSourcePolicy from %s", icspFile)
	}
	return []operatorv1alpha1.ImageContentSourcePolicy{*icsp}, nil
}
