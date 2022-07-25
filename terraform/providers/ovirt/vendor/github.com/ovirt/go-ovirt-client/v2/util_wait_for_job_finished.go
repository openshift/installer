package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

// waitForJobFinished waits for a job to truly finish. This is especially important when disks are involved as their
// status changes to OK prematurely.
//
// correlationID is a query parameter assigned to a job before it is sent to the ovirt engine, it must be unique and
// under 30 chars. To set a correlationID add `Query("correlation_id", correlationID)` to the engine API call, for
// example:
//
//     correlationID := fmt.Sprintf("image_transfer_%s", utilrand.String(5))
//     conn.
//         SystemService().
//         DisksService().
//         DiskService(diskId).
//         Update().
//         Query("correlation_id", correlationID).
//         Send()
func (o *oVirtClient) waitForJobFinished(correlationID string, retries []RetryStrategy) error {
	return retry(
		fmt.Sprintf("waiting for job with correlation ID %s to finish", correlationID),
		o.logger,
		retries,
		func() error {
			jobResp, err := o.conn.SystemService().JobsService().List().Search(fmt.Sprintf("correlation_id=%s", correlationID)).Send()
			if err != nil {
				return err
			}
			if jobSlice, ok := jobResp.Jobs(); ok {
				allJobsFinished := true
				for _, job := range jobSlice.Slice() {
					if status, _ := job.Status(); status == ovirtsdk.JOBSTATUS_STARTED {
						allJobsFinished = false
					}
				}
				if allJobsFinished {
					return nil
				}
			}
			return newError(EPending, "job for correlation ID %s still pending", correlationID)
		},
	)
}
