#!/bin/env groovy​

folder("triggers")

job("triggers/upstream-terraform-trigger") {
  description('Tectonic Installer using Terraform Upstream. Changes here will be reverted automatically.')

  logRotator(10, 1000)
  wrappers {
    colorizeOutput()
    timestamps()
  }

  triggers {
    cron('@daily')
  }

  parameters {
    stringParam('builder_image', 'quay.io/coreos/tectonic-builder:v1.36-upstream-terraform', 'tectonic-builder docker image with upstream Terraform')
  }

  steps {
    triggerBuilder {
      configs {
        blockableBuildTriggerConfig {
          projects("tectonic-installer/master")
          block {
            buildStepFailureThreshold("FAILURE")
            unstableThreshold("UNSTABLE")
            failureThreshold("FAILURE")
          }
          configs {
            currentBuildParameters()
          }
        }
      }
    }
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Installer Upstream Terraform Build")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      notifyRepeatedFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
