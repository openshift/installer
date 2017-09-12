#!/bin/env groovy​

multibranchPipelineJob("tectonic-installer") {
  description("Install a Kubernetes cluster the CoreOS Tectonic Way: HA, self-hosted, RBAC, etcd Operator, and more\nThis job is manage by tectonic-installer.\nChanges here will be reverted automatically.")
  branchSources {
    github {
      scanCredentialsId("37477e0c-2ab6-46fe-a83b-64b1add4777d")
      checkoutCredentialsId("37477e0c-2ab6-46fe-a83b-64b1add4777d")
      apiUri("")
      repoOwner("coreos")
      repository("tectonic-installer")
      buildForkPRHead(false)
      buildForkPRMerge(true)
      buildOriginBranch(true)
      buildOriginBranchWithPR(false)
      buildOriginPRHead(false)
      buildOriginPRMerge(true)
    }
  }

// Post a message in this group https://groups.google.com/forum/#!topic/job-dsl-plugin/amHzGEanrco to see if there is another way to get the strategy
// back again
  configure { project ->
    project / 'sources' / 'data' / 'jenkins.branch.BranchSource'/ strategy(class: 'jenkins.branch.DefaultBranchPropertyStrategy') {
      properties(class: 'java.util.Arrays$ArrayList') {
        a(class: 'jenkins.branch.BranchProperty-array'){
          'jenkins.branch.NoTriggerBranchProperty'()
        }
      }
    }
  }

  orphanedItemStrategy {
    discardOldItems {
      daysToKeep(10)
      numToKeep(1000)
    }
  }
  triggers {
    periodic(1440)
  }
}
