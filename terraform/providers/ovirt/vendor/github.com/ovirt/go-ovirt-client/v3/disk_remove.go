package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveDisk(diskID DiskID, retries ...RetryStrategy) error {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	return retry(
		fmt.Sprintf("removing disk %s", diskID),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().DisksService().DiskService(string(diskID)).Remove().Send()
			return err
		},
	)
}
