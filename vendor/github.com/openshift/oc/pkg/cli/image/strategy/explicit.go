package strategy

import (
	"context"
	"fmt"
	"sync"

	"k8s.io/klog/v2"

	"github.com/openshift/library-go/pkg/image/reference"
	"github.com/openshift/library-go/pkg/image/registryclient"
)

type explicitICSPStrategy struct {
	lock sync.Mutex

	alternates            map[reference.DockerImageReference][]reference.DockerImageReference
	icspFile              string
	readICSPsFromFileFunc readICSPsFromFileFunc
}

var _ registryclient.AlternateBlobSourceStrategy = &explicitICSPStrategy{}

// NewICSPExplicitStrategy returns ICSP alternate strategy which always reads
// alternate sources first rather than original requested.
func NewICSPExplicitStrategy(file string) registryclient.AlternateBlobSourceStrategy {
	return &explicitICSPStrategy{
		icspFile:              file,
		alternates:            make(map[reference.DockerImageReference][]reference.DockerImageReference),
		readICSPsFromFileFunc: readICSPsFromFile,
	}
}

func (s *explicitICSPStrategy) FirstRequest(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
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

func (s *explicitICSPStrategy) OnFailure(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.alternates) == 0 {
		return nil, fmt.Errorf("no alternative image references found for image: %s", locator.String())
	}
	return s.alternates[locator], nil
}

// resolve gathers possible image sources for a given image
// gathered from ImageContentSourcePolicy objects and user-passed image.
// Will lookup from cluster or from ImageContentSourcePolicy file passed from user.
// Image reference of user-given image may be different from original in case of mirrored images.
func (s *explicitICSPStrategy) resolve(ctx context.Context, imageRef reference.DockerImageReference) ([]reference.DockerImageReference, error) {
	if len(s.icspFile) == 0 {
		return nil, fmt.Errorf("no ImageContentSourceFile specified")
	}
	klog.V(5).Infof("Reading ICSP from file %s", s.icspFile)
	icspList, err := s.readICSPsFromFileFunc(s.icspFile)
	if err != nil {
		return nil, err
	}
	// always add the original as the last reference
	imageRefList, err := alternativeImageSourcesICSP(imageRef, icspList, true)
	if err != nil {
		return nil, err
	}
	return imageRefList, nil
}

type explicitIDMSStrategy struct {
	lock sync.Mutex

	alternates            map[reference.DockerImageReference][]reference.DockerImageReference
	idmsFile              string
	readIDMSsFromFileFunc readIDMSsFromFileFunc
}

var _ registryclient.AlternateBlobSourceStrategy = &explicitIDMSStrategy{}

// NewIDMSExplicitStrategy returns IDMS alternate strategy which always reads
// alternate sources first rather than original requested.
func NewIDMSExplicitStrategy(file string) registryclient.AlternateBlobSourceStrategy {
	return &explicitIDMSStrategy{
		idmsFile:              file,
		alternates:            make(map[reference.DockerImageReference][]reference.DockerImageReference),
		readIDMSsFromFileFunc: readIDMSsFromFile,
	}
}

func (s *explicitIDMSStrategy) FirstRequest(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
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

func (s *explicitIDMSStrategy) OnFailure(ctx context.Context, locator reference.DockerImageReference) (alternateRepositories []reference.DockerImageReference, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if len(s.alternates) == 0 {
		return nil, fmt.Errorf("no alternative image references found for image: %s", locator.String())
	}
	return s.alternates[locator], nil
}

// resolve gathers possible image sources for a given image
// gathered from ImageDigestMirrorSet objects and user-passed image.
// Will lookup from cluster or from ImageDigestMirrorSet file passed from user.
// Image reference of user-given image may be different from original in case of mirrored images.
func (s *explicitIDMSStrategy) resolve(ctx context.Context, imageRef reference.DockerImageReference) ([]reference.DockerImageReference, error) {
	if len(s.idmsFile) == 0 {
		return nil, fmt.Errorf("no ImageContentSourceFile specified")
	}
	klog.V(5).Infof("Reading IDMS from file %s", s.idmsFile)
	idmsList, err := s.readIDMSsFromFileFunc(s.idmsFile)
	if err != nil {
		return nil, err
	}
	// always add the original as the last reference
	imageRefList, err := alternativeImageSourcesIDMS(imageRef, idmsList, true)
	if err != nil {
		return nil, err
	}
	return imageRefList, nil
}
