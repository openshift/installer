package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"k8s.io/apimachinery/pkg/util/wait"
)

const jobTypeName = "job"

// listJobs lists jobs in the vpc.
func (o *ClusterUninstaller) listJobs() (cloudResources, error) {
	var jobs *models.Jobs
	var job *models.Job
	var err error

	o.Logger.Debugf("Listing jobs")

	if o.jobClient == nil {
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listJobs: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	jobs, err = o.jobClient.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}

	result := []cloudResource{}
	for _, job = range jobs.Jobs {
		// https://github.com/IBM-Cloud/power-go-client/blob/master/power/models/job.go
		if strings.Contains(*job.Operation.ID, o.InfraID) {
			if *job.Status.State == "completed" {
				continue
			}
			o.Logger.Debugf("listJobs: FOUND: %s (%s) (%s)", *job.Operation.ID, *job.ID, *job.Status.State)
			result = append(result, cloudResource{
				key:      *job.Operation.ID,
				name:     *job.Operation.ID,
				status:   *job.Status.State,
				typeName: jobTypeName,
				id:       *job.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

// DeleteJobResult The different states deleting a job can take.
type DeleteJobResult int

const (
	// DeleteJobSuccess A job has finished successfully.
	DeleteJobSuccess DeleteJobResult = iota

	// DeleteJobRunning A job is currently running.
	DeleteJobRunning

	// DeleteJobError A job has resulted in an error.
	DeleteJobError
)

func (o *ClusterUninstaller) deleteJob(item cloudResource) (DeleteJobResult, error) {
	var job *models.Job
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteJob: case <-ctx.Done()")
		return DeleteJobError, ctx.Err() // we're cancelled, abort
	default:
	}

	job, err = o.jobClient.Get(item.id)
	if err != nil {
		o.Logger.Debugf("listJobs: deleteJob: job %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Job %q", item.name)
		return DeleteJobSuccess, nil
	}

	switch *job.Status.State {
	case "completed":
		//		err = o.jobClient.Delete(item.id)
		//		if err != nil {
		//			return DeleteJobError, fmt.Errorf("failed to delete job %s: %w", item.name, err)
		//		}

		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Debugf("Deleting job %q", item.name)

		return DeleteJobSuccess, nil

	case "active":
		o.Logger.Debugf("Waiting for job %q to delete (status is %q)", item.name, *job.Status.State)
		return DeleteJobRunning, nil

	case "failed":
		err = fmt.Errorf("@TODO we cannot query error message inside the job")
		return DeleteJobError, fmt.Errorf("job %v has failed: %w", item.id, err)

	default:
		o.Logger.Debugf("Default waiting for job %q to delete (status is %q)", item.name, *job.Status.State)
		return DeleteJobRunning, nil
	}
}

// destroyJobs removes all job resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyJobs() error {
	firstPassList, err := o.listJobs()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(jobTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyJobs: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			result, err2 := o.deleteJob(item)
			switch result {
			case DeleteJobSuccess:
				o.Logger.Debugf("destroyJobs: deleteJob returns DeleteJobSuccess")
				return true, nil
			case DeleteJobRunning:
				o.Logger.Debugf("destroyJobs: deleteJob returns DeleteJobRunning")
				return false, nil
			case DeleteJobError:
				o.Logger.Debugf("destroyJobs: deleteJob returns DeleteJobError: %v", err2)
				return false, err2
			default:
				return false, fmt.Errorf("destroyJobs: deleteJob unknown result enum %v", result)
			}
		})
		if err != nil {
			o.Logger.Fatal("destroyJobs: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(jobTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyJobs: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyJobs: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listJobs()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyJobs: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyJobs: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
