#### Bootstrap Service Records ####

For the purposes of diagnosing installation failures that occur during bootstrapping, the
progresses of services running on the bootstrap machines are tracked in json files in the
/var/log/openshift directory. The progress for each service is tracked in its own file. For
example, the bootkube service progress is tracked in the /var/log/openshift/bootkube.json
file. The following progress events are tracked.
* A service adds an entry when the service starts.
* A service adds an entry when the service ends. The entry includes the result of the service
invocation, either success or failure. If the invocation failed, then the entry includes
the line number of the error and the last three lines from the service's journal log.
* A service adds an entry when a stage of the service starts. An example of a service stage
is the cvo-bootstrap stage of the bootkube service. During that stage, the
cluster-version-operator renders its manifests.
* A service  adds an entry when a stage of the service ends. This is similar to the entry
added when a service ends, including having a result and error information if applicable.

##### Managing Service Records in a Service #####

To track its progress, a service should source the /usr/local/bin/bootstrap-service-record.sh
script. When a service sources the script, the script will add an entry to the json file for
the service indicating that the service started. When the service ends, either successfully
or due to an error, the script will add an entry to the json file for the service indicating that
the service ended. The script will consider whether the last command executed was successful or
not in order to determine whether the service was successful.

For tracking stages, the service should call functions from the sourced script.
* The service should call the `record_service_stage_start` function when a stage of the service
starts. The function takes the name of the stage as its single argument.
* The service should call the `record_service_stage_success` function when a stage of the service
ends successfully.
* The service should call the `record_service_stage_failure` function when a stage of the service
ends due to a failure. The script will automatically record an entry for a stage failure if the
service ends during the execution of a stage.

###### Pre- and Post-Commands ######

If a service has pre- or post-commands that could either run for significant periods or could
potentially fail, then those commands should add to the json file as well. Such a command should
source the same /usr/local/bin/bootstrap-service-record.sh script. It should also set either the
`PRE_COMMAND` or `POST_COMMAND` environment variable with a value that identifies the command.
For example, kubelet.service has a pre-command of /usr/local/bin/kubelet-pause-image.sh. The
kubelet-pause-image.sh script sets the `PRE_COMMAND` environment variable to "kubelet-pause-image"
before sourcing the bootstrap-service-record.sh script. All of the entries for the pre-command
will contain a `preCommand` field with the "kubelet-pause-image" value.

###### Sample Script #######

```shell script
#!/usr/bin/env bash
set -euoE pipefail

# Source the script to record service entries.
# This will create en entry for the start of the service.
. /usr/local/bin/bootstrap-service-record.sh

# Record the start of the "first" stage.
record_service_stage_start "first"

# Record the successful end of the "first" stage.
record_service_stage_success

while true
do
    # Record the start of the "second" stage.
    record_service_stage_start "second-stage"
    
    if [ some_check ]
    then
        # Record the successful end of the "second" stage.
        record_service_stage_success
        break
    else
        # Record the failing end of the "second" stage.
        record_service_stage_failure  
    fi
done

# Record the start of the third stage.
record_service_stage_start "third"

# If the command fails, then an entry will be recorded for the failing end
# of the third stage and an entry will be recorded for the failing end of
# the service.
some_command_that_may_fail

# Record the end of the third stage.
record_service_stage_success

# Since this is the end of the script, an entry will be recorded for the
# successful end of the service.
```