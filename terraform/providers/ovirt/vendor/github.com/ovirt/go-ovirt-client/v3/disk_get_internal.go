package ovirtclient

// getDisk is the internal copy of GetDisk which returns a diskWithData. When code generation becomes better,
// this should be unified with GetDisk.
func (m *mockClient) getDisk(id DiskID, _ ...RetryStrategy) (*diskWithData, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.disks[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "disk with ID %s not found", id)
}
