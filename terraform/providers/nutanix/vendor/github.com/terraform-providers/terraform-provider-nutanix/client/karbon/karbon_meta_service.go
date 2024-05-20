package karbon

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

// MetaOperations ...
type MetaOperations struct {
	client *client.Client
}

// Service ...
type MetaService interface {
	// karbon v2.1
	GetVersion() (*MetaVersionResponse, error)
	GetSemanticVersion() (*MetaSemanticVersionResponse, error)
}

func (op MetaOperations) GetVersion() (*MetaVersionResponse, error) {
	ctx := context.TODO()
	path := "/v1-alpha.1/version"
	req, err := op.client.NewRequest(ctx, http.MethodGet, path, nil)
	karbonMetaVersionResponse := new(MetaVersionResponse)
	if err != nil {
		return nil, err
	}
	return karbonMetaVersionResponse, op.client.Do(ctx, req, karbonMetaVersionResponse)
}

func (op MetaOperations) GetSemanticVersion() (*MetaSemanticVersionResponse, error) {
	const expectedVersionLength int = 3
	metaSemanticVersionResponse := new(MetaSemanticVersionResponse)
	rawVersion, err := op.GetVersion()
	if err != nil {
		return nil, err
	}
	splitted := strings.Split(*rawVersion.Version, ".")
	if len(splitted) != expectedVersionLength {
		return nil, fmt.Errorf("expected karbon version to be consisting out of 3 elements but was %v", len(splitted))
	}

	major, err := strconv.Atoi(splitted[0])
	if err != nil {
		return nil, fmt.Errorf("could not convert major version to int64")
	}
	minor, err := strconv.Atoi(splitted[1])
	if err != nil {
		return nil, fmt.Errorf("could not convert minor version to int64")
	}
	rev, err := strconv.Atoi(splitted[2])
	if err != nil {
		return nil, fmt.Errorf("could not convert rev version to int64")
	}
	metaSemanticVersionResponse.MajorVersion = int64(major)
	metaSemanticVersionResponse.MinorVersion = int64(minor)
	metaSemanticVersionResponse.RevisionVersion = int64(rev)
	return metaSemanticVersionResponse, nil
}
