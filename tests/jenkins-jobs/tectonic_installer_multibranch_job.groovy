#!/bin/env groovy​

multibranchPipelineJob("tectonic-installer") {
  description("Install a Kubernetes cluster the CoreOS Tectonic Way: HA, self-hosted, RBAC, etcd Operator, and more\nThis job is manage by tectonic-installer.\nChanges here will be reverted automatically.")
  branchSources {
    branchSource {
      source {
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
      strategy {
        defaultBranchPropertyStrategy {
          props {
            noTriggerBranchProperty()
          }
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
