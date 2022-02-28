# Contributing to IBM Cloud Terraform Provider

**Table of contents**
- [Contributing to IBM Cloud Terraform Provider](#contributing-to-ibm-cloud-terraform-provider)
  - [Issues](#issues)
    - [Issue reporting checklists](#issue-reporting-checklists)
      - [Bug Reports](#bug-reports)
      - [Feature requests](#feature-requests)
      - [Questions](#questions)
  - [Pull requests](#pull-requests)
    - [Checklists for contribution](#checklists-for-contribution)
      - [Enhancement or Bugfix to a resource](#enhancement-or-bugfix-to-a-resource)
      - [New resource](#new-resource)
    - [Writing acceptance tests](#writing-acceptance-tests)
      - [Acceptance tests often cost money to run](#acceptance-tests-often-cost-money-to-run)
      - [Running an acceptance test](#running-an-acceptance-test)
      - [Writing an acceptance test](#writing-an-acceptance-test)
  - [Release management](#release-management)
    - [Production release](#production-release)
    - [Pre-production release](#pre-production-release)
    - [Dev release](#dev-release)

---
**Change history**
| Change date | Change description | 
|-------------|--------------------|
| 03-Jan-2022 | Updated the release management process to include pre-production & production releases of the IBM Cloud Provider |

---
**First:** if you are unsure  _anything_, just ask or submit the issue or create pull request anyways. We appreciate any sort of contributions.

However, for those individuals who want a bit more guidance on the best way to contribute to the project. This document will cover what we are looking for. By addressing all the points we are looking for,  it raises the chances to quickly merge or address your contributions.

Specifically, we have provided checklists below for each type of an issue and pull request that can happen on the project. These checklists represent everything we need to review and respond quickly.

## Issues

### Issue reporting checklists

We welcome issues of all kinds including feature requests, bug reports, and general questions. Following checklists provides the guidelines for the well-formed issue of specific type.

#### Bug Reports

 - [ ] __Test against latest release__: Make sure you test against the latest released version. It is possible we already fixed the bug that you are experiencing.

 - [ ] __Search for possible duplicate reports__: It is helpful to keep bug reports consolidated to one thread, so do a quick search on existing bug reports to check if anybody else has reported the same thing. You can scope searches by the label `bug` to quickly narrow down the search.

 - [ ] __Include steps to reproduce__: Provide steps to reproduce the issue, along with your `.tf` files, with secrets removed, so we can try to reproduce it. Without this, it makes much harder to fix an issue.

 - [ ] __For panics, include `crash.log`__: If you experience a panic, please create a [gist](https://gist.github.com) of the **entire** generated crash log to us to look at. Double check no sensitive items are present in the log.

#### Feature requests

 - [ ] __Search for possible duplicate requests__: It is helpful to keep requests consolidated to one thread, so do a quick search on an existing requests to check if anybody else has reported the same requests. You can scope searches by the label `enhancement` to quickly narrow down the search.

 - [ ] __Include a use case description__: In addition to describing the behavior of the feature you would like to see added, it is helpful you to lay out the reason why the feature would be important and how it would benefit Terraform users?

#### Questions

 - [ ] __Search for answers in Terraform documentation__: We are happy to answer questions in GitHub issues, but it helps reduce issue churn and maintainer workload if you work to find answers to common questions in the documentation. Usually question issues results in documentation updates to help future users, so if you do not find an answer, you can give us pointers for where you would expect to see it in the documentation.

## Pull requests

Thank you for contributing! Here you will find the information on what to include in your pull request (PRs) to ensure it is accepted quickly.

 * For PRs that follow the guidelines, we expect to review and merge very quickly.
 * PRs that do not follow the guidelines is annotated with what they are missing. A community or core team member may be able to swing around and help finish up the work, but these PRs generally hang out much longer until completed and merged.

### Checklists for contribution

There are different kinds of contribution, each of which has its own standards for a speedy review. The following sections describe the guidelines for each type of the contribution.

#### Enhancement or Bugfix to a resource

Working on an existing resources is a great way to start as a Terraform contributor because you can work within an existing code and tests to get a feel for what to do.

 - [ ] __Acceptance test coverage of new behavior__: Existing resources each have a set of [acceptance tests][acctests] covering their functionality. These tests  exercises all the behavior of the resource. Whether you are adding something or fixing a bug, the idea is to have an acceptance test that fails if your code are removed. Sometimes it is sufficient to **enhance** an existing test by adding an assertion or tweaking the configuration that are used, but often a new test is better to add. You can copy or paste an existing test and follow the conventions you see there, modifying the test to exercise the behavior of your code.

 - [ ] __Documentation updates__: If your code makes any changes that need to be documented, you should include those documentation updates in the same PR. 
   
 - [ ] __Well-formed Code__: Do your best to follow an existing conventions you see in the codebase, and ensure your code is formatted with **go fmt**. (The Travis CI build fails if **go fmt** has not been run on incoming code.) The PR reviewers can help out on this front, and may provide comments with suggestions on how to improve the code.

#### New resource

Implementing a new resource is a good way to learn more about how Terraform interacts with upstream APIs. There are plenty of examples to draw from in the existing resources, but you still get to implement something completely new.

 - [ ] __Minimal LOC__: It can be inefficient for both the reviewer and author to go through long feedback cycles on a big PR with many resources. We therefore encourage you to only submit **one resource at a time**.
 - [ ] __Acceptance tests__: New resources should include acceptance tests covering their behavior. See [Writing Acceptance Tests](#writing-acceptance-tests) below for a detailed guide on how to approach these.
 - [ ] __Documentation__: Each resource gets a page in the Terraform documentation. The [Terraform website](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs) source is in this repository and includes instructions for getting a local copy of the site up and running if you would like to preview your changes. For a resource, you will want to add a new file in the appropriate place and add a link to the sidebar for that page.
 - [ ] __Well-formed Code__: Do your best to follow an existing conventions you see in the codebase, and ensure your code is formatted with **go fmt**. (The Travis CI build fail if **go fmt** has not been run on incoming code.) The PR reviewers help out on this front, and may provide comments with suggestions on how to improve the code.

### Writing acceptance tests

Terraform includes an acceptance test harness that does most of the repetitive work involved in testing a resource.

#### Acceptance tests often cost money to run

Acceptance tests create real resources, they often cost money to run. Because the resources only exist for a short period of time, the total amount of money required is usually a relatively small. Nevertheless, we do not want financial limitations to be a barrier for contribution, so if you are unable to pay to run acceptance tests for your contribution, simply mention this in your PR. We happily accept **best effort** implementations of acceptance tests and run them for you on our side. This might mean that your PR takes a bit longer to merge, but it definitely is not a blocker for contributions.

#### Running an acceptance test

Acceptance tests can be run by using the **testacc** target in the **Makefile**. The individual tests to run is controlled by using a regular expression. Prior to running the tests provider configuration details such as access keys must be made available as an environment variables.

For example, to run an acceptance test, the following environment variables must be set:

```sh
export IC_API_KEY=...
export IAAS_CLASSIC_API_KEY=...
export IAAS_CLASSIC_USERNAME=...
```

For certain tests, the following values may also needs to be set:

```sh
export IBM_ORG=...
export IBM_SPACE=...
export IBM_ID1=...
export IBM_ID2=...
export IBM_IAMUSER=...
```

You can enable the Terraform logs by setting the following environment variable:

```sh
export TF_LOG=DEBUG
```

Tests can then be run by specifying the target provider and a regular expression defining the tests to run:

```sh
$ make testacc TEST=./ibm TESTARGS='-run=TestAccIBMComputeVmInstance_basic'
==> Checking that code complies with gofmt requirements...
go generate ./...
TF_ACC=1 go test ./ibm -v -run=TestAccIBMComputeVmInstance_basic -timeout 700m
=== RUN   TestAccIBMComputeVmInstance_basic
--- PASS: TestAccIBMComputeVmInstance_basic (177.48s)
PASS
ok      github.com/terraform-providers/terraform-provider-ibm/ibm   177.504s
```

Entire resource test suites can be targeted by using the naming convention to write the regular expression. For example, to run all tests of the `ibm_compute_vm_instance` resource rather than just the update test, you can start testing like this:

```sh
$ make testacc TEST=./ibm TESTARGS='-run=TestAccIBMComputeVmInstance'
==> Checking that code complies with gofmt requirements...
go generate ./...
TF_ACC=1 go test ./builtin/providers/azurerm -v -run=TestAccIBMComputeVmInstance -timeout 700m
=== RUN   TestAccIBMComputeVmInstance_basic
--- PASS: TestAccIBMComputeVmInstance_basic (137.74s)
=== RUN   TestAccIBMComputeVmInstance_basic_import
--- PASS: TestAccIBMComputeVmInstance_basic_import (180.63s)
PASS
ok      github.com/terraform-providers/terraform-provider-ibm/ibm   318.392s
```

#### Writing an acceptance test

Terraform has a framework for writing acceptance tests which minimises the amount of boilerplate code necessary to use common testing patterns. The entry point to the framework is the `resource.Test()` function.

Tests are divided into `TestSteps`. Each `TestStep` proceeds by applying some
Terraform configuration by using the provider under test, and then verifying that results are as expected by making assertions by using the provider API. It is common for a single test function to exercise both the creation of and updates to a single resource. Most tests follow a similar structure.

1. Pre-flight checks are made to ensure that sufficient provider configuration is available to proceed. For example, in an acceptance test `IAAS_CLASSIC_API_KEY` , `IAAS_CLASSIC_USERNAME` and `IC_API_KEY` must be set prior to running acceptance tests. This is common to all tests exercising a single provider.

Each `TestStep` is defined in the call to `resource.Test()`. Most assertion functions are defined out of band with the tests. This keeps the tests readable, and allows reuse of assertion functions across different tests of the same type of resource. The definition of a complete test looks as follows:

```go
func TestAccIBMComputeVmInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIBMComputeVmInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccIBMComputeVmInstanceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccIBMComputeVmInstanceExists("ibm_compute_vm_instance.terraform-acceptance-test-1"),
				),
			},
        },
    })
}
```

When executing the test, the following steps are taken for each `TestStep`:

1. The Terraform configuration required for the test is applied. This is responsible for configuring the resource under test, and any dependencies it may have. For example, to test the `ibm_compute_vm_instance` resource. This results in configuration which looks like this:

```terraform
resource "ibm_compute_vm_instance" "terraform-acceptance-test-1" {
   hostname = "terraform-sample-blockDeviceTemplateGroup"
   domain = "bar.example.com"
   datacenter = "ams01"
   public_network_speed = 10
   hourly_billing = false
   cores = 1
   memory = 1024
   local_disk = false
   image_id = 12345
   tags = [
     "collectd",
     "mesos-master"
   ]
   public_subnet = "50.97.46.160/28"
   private_subnet = "10.56.109.128/26"
} 
```

2. Assertions are run by using the provider API. These use the provider API directly rather than asserting against the resource state. For example, to verify that the `ibm_compute_vm_instance` described above was created successfully, a test function like this is used:

```go
    func resourceIBMComputeVmInstanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	service := services.GetVirtualGuestService(meta.(ClientSession).SoftLayerSession())
	guestID, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Not a valid ID, must be an integer: %s", err)
	}

	result, err := service.Id(guestID).GetObject()
	if err != nil {
		if apiErr, ok := err.(sl.Error); ok {
			if apiErr.StatusCode == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return result.Id != nil && *result.Id == guestID, nil
	}
```

Notice that the only information used from the Terraform state is the ID of the resource - though in this case it is necessary to split the ID into constituent parts in order to use the provider API. For computed properties, we instead assert that the value saved in the Terraform state was the expected value if possible. The testing framework provides helper functions for several common types of check. For example:

```go
    resource.TestCheckResourceAttr("ibm_compute_vm_instance.terraform-test-1", "hourly_billing", "true"),
```

1. The resources created by the test are destroyed. This step happens automatically, and is the equivalent of calling `terraform destroy`.

2. Assertions are made against the provider API to verify that the resources have indeed been removed. If these checks fail, the test fails and reports **dangling resources**. The code to ensure that the `ibm_compute_vm_instance` shown above looks as follow:

```go
    go func testAccIBMComputeVmInstanceDestroy(s *terraform.State) error {
	service := services.GetVirtualGuestService(testAccProvider.Meta().(ClientSession).SoftLayerSession())

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_compute_vm_instance" {
			continue
		}

		guestID, _ := strconv.Atoi(rs.Primary.ID)

		// Try to find the guest
		_, err := service.Id(guestID).GetObject()

		// Wait

		if err == nil {
			return fmt.Errorf("Virtual guest still exists: %s", rs.Primary.ID)
		else if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for virtual guest (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}
	return nil
	}
```

These functions usually test only for the resource directly under test.

## Release management

The `IBM Cloud Provider for Terraform` release can be mainly classified in to three types:
- Production release
- Pre-production release
- Dev release

### Production release
Typically, the production release of the `IBM Cloud Provider for Terraform` will be made, once in a month. The release can be major or minor based on the PR's commited. The production release is targetted from branch **release**. Once the release is published, users can download the binary from [Terraform registry](https://registry.terraform.io/providers/IBM-Cloud/ibm/latest).

#### How to use this provider

To install Terraform 0.13 or higher version provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```terraform
  terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = "release version number"
    }
  }
}

provider "ibm" {
  # Configuration options
}
```

### Pre-production release

Typically, the pre-production releases of the `IBM Cloud Provider for Terraform` will be made every two weeks or on-demand.  This release will be tagged with **pre** tag.  For example, **1.38.0-pre0**. The pre-production release is targeted from branch **master** and marked as pre-release (release is identified as non-production ready).

~> **Note** A pre-release version is a version number that contains a suffix introduced by a dash, such as **1.38.0-pre**. A pre-release version can be selected only by an exact version constraint (the = operator or no operator). Pre-release versions do not match inexact operators such as `>=`, `~>`, etc.

#### How to use this provider
To install Terraform 0.13 or higher version provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```terraform
  terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = "1.38.0-pre0"
    }
  }
}

provider "ibm" {
  # Configuration options
}
```

### Dev release

The individual developers or the IBM Cloud Service team can make their own `dev` releases, from their respective Git repository (forked from https://github.com/IBM-Cloud/terraform-provider-ibm).  

Note: You can use the existing GitHub actions to run the release workflows, in your forked repository. You can prepare a `dev` release by adding a new version tag in your repository.

#### How to use this provider (dev release)

- Download the respective binary from `dev` release repository
- Unzip the folder
- Depending on particular target platform create a directory structure as stated
   ```
   mkdir -p $TERRAFORMHOME/registry.terraform.io/ibm-cloud/ibm/#version#/#target# 
   ```
    where, TERRAFORMHOME 
    - Windows: %APPDATA%/terraform.d/plugins
    - Mac OS X: $HOME/.terraform.d/plugins
    - Linux: $HOME/.terraform.d/plugins
    - #version# : The release version of the binary
    - #target# : specifies a particular target platform using a format like darwin_amd64, linux_amd64, windows_amd64.
- Create a CLI configuration file 
  - On Windows, the file must be named **terraform.rc** and placed in the relevant user's %APPDATA% directory. 
  - On all other systems, the file must be named **.terraformrc** (note the leading period) and placed directly in the home directory of the relevant user.
- Add below content to CLI configuration file
  ```terraform
  provider_installation {
    filesystem_mirror {
        path = "<TERRAFORMHOME>"
        include = ["ibm-cloud/ibm"]
    }
    direct {
        include = ["*/*"]
    }
  } 
  ```
