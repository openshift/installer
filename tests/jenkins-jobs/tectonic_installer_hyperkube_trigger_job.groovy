#!/bin/env groovy​

folder("triggers")

job("triggers/tectonic-installer-hyperkube-trigger") {
  logRotator(10, 1000)
  description('Tectonic Installer using latest the Hyperkube build against master. Changes here will be reverted automatically.')

  wrappers {
    colorizeOutput()
    timestamps()
  }

  parameters {
    stringParam('hyperkube_image', '', "Please define the param like: {hyperkube=\"HYPERKUBE_IMAGE\"")
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
      customMessage("Tectonic Installer Hyperkube Build - Master Branch - Hyperkube: \${hyperkube_image}")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      notifyRepeatedFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
