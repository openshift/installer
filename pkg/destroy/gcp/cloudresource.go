package gcp

import "github.com/openshift/installer/pkg/types/gcp"

// cloudResource hold various fields for any given cloud resource
type cloudResource struct {
	key      string
	name     string
	project  string
	status   string
	typeName string
	url      string
	zone     string
	quota    []gcp.QuotaUsage
	scope    resourceScope
}

type cloudResources map[string]cloudResource

// resourceScope captures whether a resource is global or regional.
// Currently this is only used for resources handled by the compute service
// in order to perform wait operations to ensure they are successfully deleted.
type resourceScope string

const (
	global   resourceScope = "global"
	regional resourceScope = "regional"
	zonal    resourceScope = "zonal"
)

func (r cloudResources) insert(resources ...cloudResource) cloudResources {
	for _, resource := range resources {
		r[resource.key] = resource
	}
	return r
}

func (r cloudResources) delete(resources ...cloudResource) cloudResources {
	for _, resource := range resources {
		delete(r, resource.key)
	}
	return r
}

func (r cloudResources) list() []cloudResource {
	values := []cloudResource{}
	for _, value := range r {
		values = append(values, value)
	}
	return values
}
