package backupvaults

type BackupVaultResourceOperationPredicate struct {
	ETag     *string
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p BackupVaultResourceOperationPredicate) Matches(input BackupVaultResource) bool {

	if p.ETag != nil && (input.ETag == nil && *p.ETag != *input.ETag) {
		return false
	}

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
