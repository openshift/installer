package webpubsub

type PrivateEndpointConnectionOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p PrivateEndpointConnectionOperationPredicate) Matches(input PrivateEndpointConnection) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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

type PrivateLinkResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p PrivateLinkResourceOperationPredicate) Matches(input PrivateLinkResource) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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

type SharedPrivateLinkResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p SharedPrivateLinkResourceOperationPredicate) Matches(input SharedPrivateLinkResource) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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

type SignalRServiceUsageOperationPredicate struct {
	CurrentValue *int64
	Id           *string
	Limit        *int64
	Unit         *string
}

func (p SignalRServiceUsageOperationPredicate) Matches(input SignalRServiceUsage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil && *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil && *p.Limit != *input.Limit) {
		return false
	}

	if p.Unit != nil && (input.Unit == nil && *p.Unit != *input.Unit) {
		return false
	}

	return true
}

type WebPubSubHubOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p WebPubSubHubOperationPredicate) Matches(input WebPubSubHub) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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

type WebPubSubResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p WebPubSubResourceOperationPredicate) Matches(input WebPubSubResource) bool {

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
