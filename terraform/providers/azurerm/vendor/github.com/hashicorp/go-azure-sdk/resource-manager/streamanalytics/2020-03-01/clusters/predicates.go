package clusters

type ClusterOperationPredicate struct {
	Etag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ClusterOperationPredicate) Matches(input Cluster) bool {

	if p.Etag != nil && (input.Etag == nil && *p.Etag != *input.Etag) {
		return false
	}

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil && *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}

type ClusterJobOperationPredicate struct {
	Id             *string
	StreamingUnits *int64
}

func (p ClusterJobOperationPredicate) Matches(input ClusterJob) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.StreamingUnits != nil && (input.StreamingUnits == nil && *p.StreamingUnits != *input.StreamingUnits) {
		return false
	}

	return true
}
