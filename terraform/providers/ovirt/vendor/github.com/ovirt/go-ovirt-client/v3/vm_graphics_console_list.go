package ovirtclient

import "fmt"

func (o *oVirtClient) ListVMGraphicsConsoles(vmID VMID, retries ...RetryStrategy) (result []VMGraphicsConsole, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("listing graphics consoles for VM %s", vmID),
		o.logger,
		retries,
		func() error {
			resp, err := o.conn.SystemService().VmsService().VmService(string(vmID)).GraphicsConsolesService().List().Send()
			if err != nil {
				return err
			}
			consolesList, ok := resp.Consoles()
			if !ok {
				return newFieldNotFound("graphics consoles list response", "consoles")
			}
			result = make([]VMGraphicsConsole, len(consolesList.Slice()))
			for i, c := range consolesList.Slice() {
				result[i], err = convertSDKGraphicsConsole(c, o)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	return result, err
}

func (m *mockClient) ListVMGraphicsConsoles(vmID VMID, retries ...RetryStrategy) ([]VMGraphicsConsole, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	graphicsConsoles, ok := m.graphicsConsolesByVM[vmID]
	if !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", vmID)
	}
	result := make([]VMGraphicsConsole, len(graphicsConsoles))
	for i, graphicsConsole := range graphicsConsoles {
		result[i] = graphicsConsole
	}
	return result, nil
}
