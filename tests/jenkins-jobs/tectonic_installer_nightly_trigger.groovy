#!/bin/env groovy​

job("tectonic-installer-nightly-trigger") {
  logRotator(10, 10)
  description('Tectonic Installer nightly builds against master. Changes here will be reverted automatically.')

  wrappers {
    colorizeOutput()
    timestamps()
  }

  triggers {
    cron('H 3 * * *')
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
        }
      }
    }
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Installer Nightly Build")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
