package ovirtclient

// TestConnectionClient defines the functions related to testing the connection.
type TestConnectionClient interface {
	// Test tests if the connection is alive or not.
	Test(retries ...RetryStrategy) error
}

func (o *oVirtClient) Test(retries ...RetryStrategy) error {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	return retry(
		"testing oVirt engine connection",
		o.logger,
		retries,
		func() error {
			return o.conn.SystemService().Connection().Test()
		},
	)
}

func (m *mockClient) Test(retries ...RetryStrategy) error {
	retries = defaultRetries(retries, defaultReadTimeouts(m))
	return retry(
		"testing oVirt engine connection",
		nil,
		retries,
		func() error {
			return nil
		},
	)
}
