package powervs

import (
	"strings"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/pkg/errors"
)

const jobTypeName = "job"

// listJobs lists jobs in the vpc.
func (o *ClusterUninstaller) listJobs() (cloudResources, error) {
	var jobs *models.Jobs
	var job *models.Job
	var err error

	o.Logger.Debugf("Listing jobs")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listJobs: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	jobs, err = o.jobClient.GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list jobs")
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

func (o *ClusterUninstaller) deleteJob(item cloudResource) error {
	var job *models.Job
	var err error

	job, err = o.jobClient.Get(item.id)
	if err != nil {
		o.Logger.Debugf("listJobs: deleteJob: job %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted job %q", item.name)
		return nil
	}

	if !strings.EqualFold(*job.Status.State, "active") {
		o.Logger.Debugf("Waiting for job %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting job %q", item.name)

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteJob: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	err = o.jobClient.Delete(item.id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete job %s", item.name)
	}

	return nil
}

// destroyJobs removes all job resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyJobs() error {
	found, err := o.listJobs()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(jobTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyJobs: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted job %q", item.name)
				continue
			}
			err := o.deleteJob(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(jobTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(jobTypeName); len(items) > 0 {
		return errors.Errorf("destroyJobs: %d undeleted items pending", len(items))
	}
	return nil
}
