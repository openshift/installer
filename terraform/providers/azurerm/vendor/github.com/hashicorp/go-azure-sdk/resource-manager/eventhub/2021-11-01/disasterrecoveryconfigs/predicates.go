package disasterrecoveryconfigs

type ArmDisasterRecoveryOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ArmDisasterRecoveryOperationPredicate) Matches(input ArmDisasterRecovery) bool {

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
