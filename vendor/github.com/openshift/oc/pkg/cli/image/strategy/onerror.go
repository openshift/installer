package strategy

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"

	apicfgv1 "github.com/openshift/api/config/v1"
	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	apicfgv1scheme "github.com/openshift/client-go/config/clientset/versioned/scheme"
	operatorv1alpha1scheme "github.com/openshift/client-go/operator/clientset/versioned/scheme"
	"github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/library-go/pkg/image/registryclient"
)

type onErrorICSPStrategy struct {
	lock sync.Mutex

	alternates            map[reference.DockerImageReference][]reference.DockerImageReference
	icspFile              string
	readICSPsFromFileFunc readICSPsFromFileFunc
}

var _ registryclient.AlternateBlobSourceStrategy = &onErrorICSPStrategy{}

// NewICSPOnErrorStrategy returns ICSP alternate strategy which reads alternate
// sources only after getting an error from the original requested.
func NewICSPOnErrorStrategy(file string) registryclient.AlternateBlobSourceStrategy {
	return &onErrorICSPStrategy{
		icspFile:              file,
		alternates:            make(map[reference.DockerImageReference][]reference.DockerImageReference),
		readICSPsFromFileFunc: readICSPsFromFile,
	}
}

func (s *onErrorICSPStrategy) FirstRequest(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
	return nil, nil
}

func (s *onErrorICSPStrategy) OnFailure(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
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
func (s *onErrorICSPStrategy) resolve(ctx context.Context, imageRef reference.DockerImageReference) ([]reference.DockerImageReference, error) {
	if len(s.icspFile) == 0 {
		return nil, fmt.Errorf("no ImageContentSourceFile specified")
	}
	klog.V(5).Infof("Reading ICSP from file %s", s.icspFile)
	icspList, err := s.readICSPsFromFileFunc(s.icspFile)
	if err != nil {
		return nil, err
	}
	// always add the original as the first reference
	imageRefList, err := alternativeImageSourcesICSP(imageRef, icspList, false)
	if err != nil {
		return nil, err
	}
	return imageRefList, nil
}

// alternativeImageSourcesICSP returns unique list of DockerImageReference objects from list of ImageContentSourcePolicy objects
// addSourceAsLastAlternate decides whether the original imageRef is first or the last element in the result
func alternativeImageSourcesICSP(imageRef reference.DockerImageReference, icspList []operatorv1alpha1.ImageContentSourcePolicy, addSourceAsLastAlternate bool) ([]reference.DockerImageReference, error) {
	var imageSources []reference.DockerImageReference
	klog.V(5).Infof("%v ImageReference added to potential ImageSourcePrefixes from ImageContentSourcePolicy", imageRef.AsRepository().AsV2())
	if !addSourceAsLastAlternate {
		imageSources = append(imageSources, imageRef.AsRepository().AsV2())
	}
	repo := imageRef.AsRepository().Exact()
	for _, icsp := range icspList {
		repoDigestMirrors := icsp.Spec.RepositoryDigestMirrors
		for _, rdm := range repoDigestMirrors {
			var suffix string
			var err error
			rdmSourceRef, err := reference.Parse(rdm.Source)
			if err != nil {
				return nil, fmt.Errorf("invalid source %q: %w", rdm.Source, err)
			}
			// AsV2 in the right call is required to ensure we transform docker registry
			// from docker.io to registry-1.docker.io
			if imageRef.AsRepository().AsV2() != rdmSourceRef.AsRepository().AsV2() {
				if !isSubrepo(repo, rdm.Source) {
					continue
				}
				suffix = repo[len(rdm.Source):]
			}
			klog.V(5).Infof("%v RepositoryDigestMirrors source matches given image", imageRef.AsRepository().AsV2())
			for _, m := range rdm.Mirrors {
				mRef, err := reference.Parse(m + suffix)
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
	icspData, err := os.ReadFile(icspFile)
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

type onErrorIDMSStrategy struct {
	lock sync.Mutex

	alternates            map[reference.DockerImageReference][]reference.DockerImageReference
	idmsFile              string
	readIDMSsFromFileFunc readIDMSsFromFileFunc
}

var _ registryclient.AlternateBlobSourceStrategy = &onErrorIDMSStrategy{}

// NewIDMSOnErrorStrategy returns IDMS alternate strategy which reads alternate
// sources only after getting an error from the original requested.
func NewIDMSOnErrorStrategy(file string) registryclient.AlternateBlobSourceStrategy {
	return &onErrorIDMSStrategy{
		idmsFile:              file,
		alternates:            make(map[reference.DockerImageReference][]reference.DockerImageReference),
		readIDMSsFromFileFunc: readIDMSsFromFile,
	}
}

func (s *onErrorIDMSStrategy) FirstRequest(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
	return nil, nil
}

func (s *onErrorIDMSStrategy) OnFailure(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
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
// gathered from ImageDigestMirrorSet file.
// Image reference of user-given image may be different from original in case of mirrored images.
func (s *onErrorIDMSStrategy) resolve(ctx context.Context, imageRef reference.DockerImageReference) ([]reference.DockerImageReference, error) {
	if len(s.idmsFile) == 0 {
		return nil, fmt.Errorf("no ImageDigestMirrorSet specified")
	}
	klog.V(5).Infof("Reading IDMS from file %s", s.idmsFile)
	idmsList, err := s.readIDMSsFromFileFunc(s.idmsFile)
	if err != nil {
		return nil, err
	}
	// always add the original as the first reference
	imageRefList, err := alternativeImageSourcesIDMS(imageRef, idmsList, false)
	if err != nil {
		return nil, err
	}
	return imageRefList, nil
}

// alternativeImageSourcesIDMS returns unique list of DockerImageReference objects from list of ImageDigestMirrorSet objects
// addSourceAsLastAlternate decides whether the original imageRef is first or the last element in the result
func alternativeImageSourcesIDMS(imageRef reference.DockerImageReference, idmsList []apicfgv1.ImageDigestMirrorSet, addSourceAsLastAlternate bool) ([]reference.DockerImageReference, error) {
	var imageSources []reference.DockerImageReference
	var addSource bool
	repo := imageRef.AsRepository().Exact()
	for _, idms := range idmsList {
		repoDigestMirrors := idms.Spec.ImageDigestMirrors
		for _, rdm := range repoDigestMirrors {
			var suffix string
			var err error
			rdmSourceRef, err := reference.Parse(rdm.Source)
			if err != nil {
				return nil, fmt.Errorf("invalid source %q: %w", rdm.Source, err)
			}
			// AsV2 in the right call is required to ensure we transform docker registry
			// from docker.io to registry-1.docker.io
			if imageRef.AsRepository().AsV2() != rdmSourceRef.AsRepository().AsV2() {
				if !isSubrepo(repo, rdm.Source) {
					continue
				}
				suffix = repo[len(rdm.Source):]
			}
			klog.V(5).Infof("%v ImageDigestMirrors source matches given image", imageRef.AsRepository().AsV2())
			// check valid mirrorSourcePolicy
			addSource, err = isAddSource(idmsList, rdm.Source)
			if err != nil {
				return nil, err
			}
			for _, m := range rdm.Mirrors {
				mRef, err := reference.Parse(string(m) + suffix)
				if err != nil {
					return nil, fmt.Errorf("invalid mirror %q: %w", m, err)
				}
				imageSources = append(imageSources, mRef)
				klog.V(5).Infof("%v RepositoryDigestMirrors mirror added to potential ImageSourcePrefixes from ImageDigestMirrorSet", m)
			}
		}
	}

	if addSource || len(imageSources) == 0 {
		klog.V(5).Infof("%v ImageReference added to potential ImageSourcePrefixes from ImageDigestMirrorSet", imageRef.AsRepository().AsV2())
		if addSourceAsLastAlternate {
			imageSources = append(imageSources, imageRef.AsRepository().AsV2())
		} else {
			imageSources = append([]reference.DockerImageReference{imageRef.AsRepository().AsV2()}, imageSources...)
		}
	} else {
		klog.V(5).Infof("%v ImageReference not added to potential ImageSourcePrefixes from ImageDigestMirrorSet", imageRef.AsRepository().AsV2())
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

func isAddSource(idmsList []apicfgv1.ImageDigestMirrorSet, source string) (bool, error) {
	var found bool
	var mirrorSourcePolicy apicfgv1.MirrorSourcePolicy
	for _, idms := range idmsList {
		for _, mirror := range idms.Spec.ImageDigestMirrors {
			if mirror.Source != source {
				continue
			}
			if !found {
				if mirror.MirrorSourcePolicy != "" {
					mirrorSourcePolicy = mirror.MirrorSourcePolicy
				} else {
					mirrorSourcePolicy = apicfgv1.AllowContactingSource
				}
				found = true
				continue
			}
			if mirrorSourcePolicy == apicfgv1.AllowContactingSource && mirror.MirrorSourcePolicy == "" {
				continue
			}
			if mirrorSourcePolicy != mirror.MirrorSourcePolicy {
				return found, fmt.Errorf("ImageDigestMirrorSet can only contain one MirrorSourcePolicy for source %s", source)
			}
		}
	}
	if found {
		return mirrorSourcePolicy == apicfgv1.AllowContactingSource, nil
	}

	return true, nil
}

// readIDMSsFromFileFunc is used for testing to be able to inject IDMS data
type readIDMSsFromFileFunc func(string) ([]apicfgv1.ImageDigestMirrorSet, error)

// readIDMSsFromFile appends to list of alternative image sources from IDMS file
// returns error if no idms object decoded from file data
func readIDMSsFromFile(idmsFile string) ([]apicfgv1.ImageDigestMirrorSet, error) {
	idmsData, err := os.ReadFile(idmsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read ImageDigestMirrorSet %s: %v", idmsFile, err)
	}
	if len(idmsData) == 0 {
		return nil, fmt.Errorf("no data found in ImageDigestMirrorSet %s", idmsFile)
	}
	idmsObj, err := runtime.Decode(apicfgv1scheme.Codecs.UniversalDeserializer(), idmsData)
	if err != nil {
		return nil, fmt.Errorf("error decoding ImageDigestMirrorSet from %s: %v", idmsFile, err)
	}
	idms, ok := idmsObj.(*apicfgv1.ImageDigestMirrorSet)
	if !ok {
		return nil, fmt.Errorf("could not decode ImageDigestMirrorSet from %s", idmsFile)
	}
	return []apicfgv1.ImageDigestMirrorSet{*idms}, nil
}

func isSubrepo(repo, ancestor string) bool {
	if repo == ancestor {
		return true
	}
	if len(repo) > len(ancestor) {
		return strings.HasPrefix(repo, ancestor) && repo[len(ancestor)] == '/'
	}
	return false
}
