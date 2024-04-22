package version

import (
	"errors"
	"regexp"
	"strconv"

	apimachineryversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	restclient "k8s.io/client-go/rest"
)

type ServerVersion struct {
	apimachineryversion.Info
	MajorNumber int
	MinorNumber int
}

var (
	firstNumberRegex = regexp.MustCompile("[0-9]+")
)

func NewServerVersion(version apimachineryversion.Info) (*ServerVersion, error) {
	major, err := strconv.Atoi(firstNumberRegex.FindString(version.Major))
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(firstNumberRegex.FindString(version.Minor))
	if err != nil {
		return nil, err
	}
	return &ServerVersion{version, major, minor}, nil
}

func RetrieveServerVersion(clientConfig *restclient.Config) (*ServerVersion, error) {
	if clientConfig == nil {
		return nil, errors.New("clientConfig is nil")
	}
	client, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	serverVersion, err := client.ServerVersion()
	if err != nil {
		return nil, err
	}
	return NewServerVersion(*serverVersion)
}

type ServerVersionRetriever interface {
	RetrieveServerVersion() (*ServerVersion, error)
}

func NewServerVersionRetriever(clientConfig *restclient.Config) ServerVersionRetriever {
	return &serverVersionDetector{clientConfig}
}

type serverVersionDetector struct {
	clientConfig *restclient.Config
}

func (s *serverVersionDetector) RetrieveServerVersion() (*ServerVersion, error) {
	return RetrieveServerVersion(s.clientConfig)
}
