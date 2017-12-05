#!/bin/env groovy

folder("triggers")

job("triggers/tectonic-installer-nightly-trigger_with_gcp") {
  logRotator(10, 1000)
  description('Tectonic Installer nightly builds against master only with gcp. Changes here will be reverted automatically.')

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
          configs {
            booleanParameters {
              configs {
                booleanParameterConfig {
                  name('RUN_CONFORMANCE_TESTS')
                  value(false)
                }
                booleanParameterConfig {
                  name('RUN_SMOKE_TESTS')
                  value(true)
                }
                booleanParameterConfig {
                  name('PLATFORM/AWS')
                  value(false)
                }
                booleanParameterConfig {
                  name('PLATFORM/AZURE')
                  value(false)
                }
                booleanParameterConfig {
                  name('PLATFORM/METAL')
                  value(false)
                }
                booleanParameterConfig {
                  name('PLATFORM/GCP')
                  value(true)
                }
              }
            }
          }
        }
      }
    }
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Installer Nightly Build with GCP - Master Branch")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      notifyRepeatedFailure(true)
      room('#team-installer')
      teamDomain('coreos')
    }
  }
}
