# Jenkins Jobs

This folder contains Jenkins jobs which will be used to build, test and support the Tectonic Installer.


## tectonic-builder Jenkins Job

This file (`tectonic_builder_docker_image.groovy`) creates a Jenkins job called `tectonic-builder-docker-image` to build the `tectonic-builder` docker image which is used to execute our tests.

This job has 3 input parameters:

* `TERRAFORM_UPSTREAM_URL`: If you need to build the image using the upstream `Terraform`
* `TECTONIC_BUILDER_VERSION`: The version of the `tectonic-builder` you are building.
* `DRY_RUN`: Just to build the image.

If you don't set the `TERRAFORM_UPSTREAM_URL` it will build the image using the tectonic custom terraform version.


## Tectonic Installer Nightly Trigger

This file (`tectonic_installer_nightly_trigger.groovy`) creates a Jenkins job called `tectonic-installer-nightly-trigger` under `triggers` folder to run the tests against the `Tectonic Installer` in the `master` branch.
This job will run everyday around 3AM UTC time.

Parameters:

* No input parameters are required.

## Tectonic Installer Public PR Trigger

This file (`tectonic_installer_public_pr_trigger.groovy`) creates a Jenkins job called `tectonic-installer-pr-trigger` under `triggers` folder to run the tests against the `Tectonic Installer` using the PR branch.

Parameters:

* No input parameters are required.

## Tectonic Installer Upstream Terraform Trigger

This file creates a Jenkins job called `upstream-terraform-trigger` under `triggers` folder to run the tests against the `Tectonic Installer` in the `master` branch using the `upstream Terraform`
This job will run everyday.

To change the default `builder_image` please update the code in the `tectonic_installer_upstream_terraform_trigger.groovy` file and submit a PR.

Parameters:

* `builder_image`: Tectonic-builder docker image with the upstream Terraform

## Tag/Clean AWS

This file (`maintenance/grafiti_tag_clean_aws.groovy`) creates a Jenkins job called `tag_clean_aws_resources_grafiti` under `maintenance` folder to tag and clean AWS resources across all AWS regions every 6 hours.

Parameters:

* `TAG_CLEAN`: Select if want to Tag or Clean the AWS resources.
* `START_HOUR`: Used only when tagging. Number of hours prior to now to start parsing logs from CloudTrail. Default 7.
* `END_HOUR`: Used only when tagging. Number of hours prior to now to stop parsing logs from CloudTrail. Default 0.
* `AWS_REGION`: Optional. Specific AWS region in which to tag resources.
* `DATE_VALUE_OVERRIDE`: Optional. YYYY-MM-DD formatted tag value of resources to delete.
* `GRAFITI_IMAGE`: Grafiti docker image to use in the job.
* `SCRIPT_DIR`: Folder which contains the grafiti scripts in the tectonic-installer repository.
IPT_DIR`: Folder which contains the grafiti scripts in the tectonic-installer repository.
