package namespaces

type AuthorizationRuleOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p AuthorizationRuleOperationPredicate) Matches(input AuthorizationRule) bool {

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

type RelayNamespaceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p RelayNamespaceOperationPredicate) Matches(input RelayNamespace) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
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
