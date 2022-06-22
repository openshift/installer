package ovirtclient

import "fmt"

func (o *oVirtClient) RemoveVMGraphicsConsole(
	vmID VMID,
	graphicsConsoleID VMGraphicsConsoleID,
	retries ...RetryStrategy,
) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	return retry(
		fmt.Sprintf("removing graphics consoles %s from VM %s", graphicsConsoleID, vmID),
		o.logger,
		retries,
		func() error {
			_, err = o.conn.
				SystemService().
				VmsService().
				VmService(string(vmID)).
				GraphicsConsolesService().
				ConsoleService(string(graphicsConsoleID)).
				Remove().
				Send()
			return err
		},
	)
}

func (m *mockClient) RemoveVMGraphicsConsole(
	vmID VMID,
	graphicsConsoleID VMGraphicsConsoleID,
	_ ...RetryStrategy,
) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	graphicsConsoles, ok := m.graphicsConsolesByVM[vmID]
	if !ok {
		return newError(ENotFound, "VM with ID %s not found", vmID)
	}

	foundIndex := -1
	for i, graphicsConsole := range graphicsConsoles {
		if graphicsConsole.ID() == graphicsConsoleID {
			foundIndex = i
			break
		}
	}
	if foundIndex == -1 {
		return newError(ENotFound, "Graphics console with ID %s not found on VM %s", graphicsConsoleID, vmID)
	}

	m.graphicsConsolesByVM[vmID] = append(
		m.graphicsConsolesByVM[vmID][:foundIndex],
		m.graphicsConsolesByVM[vmID][foundIndex+1:]...,
	)

	return nil
}
