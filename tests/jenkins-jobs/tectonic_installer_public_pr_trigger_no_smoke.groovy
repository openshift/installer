#!/bin/env groovy​

folder("triggers")

job("triggers/tectonic-installer-pr-trigger-no-smoke") {
  description('Tectonic Installer PR Trigger. Changes here will be reverted automatically.')

  concurrentBuild()

  logRotator(30, 100)
  label("master")

  properties {
    githubProjectUrl('https://github.com/coreos/tectonic-installer')
  }

  wrappers {
    colorizeOutput()
    timestamps()
  }

  triggers {
    ghprbTrigger {
      gitHubAuthId("")
      adminlist("")
      orgslist("coreos\ncoreos-inc")
      whitelist("")
      cron("H/5 * * * *")
      triggerPhrase("ok to test")
      onlyTriggerPhrase(false)
      useGitHubHooks(true)
      permitAll(false)
      autoCloseFailedPullRequests(false)
      displayBuildErrorsOnDownstreamBuilds(false)
      commentFilePath("")
      skipBuildPhrase(".*\\[skip\\W+ci\\].*")
      blackListCommitAuthor("")
      allowMembersOfWhitelistedOrgsAsAdmin(true)
      msgSuccess("")
      msgFailure("")
      commitStatusContext("Jenkins-Tectonic-Installer")
      buildDescTemplate("#\$pullId: \$abbrTitle")
      blackListLabels("run-smoke-tests")
      whiteListLabels("")
      includedRegions("")
      excludedRegions("")
    }
  }

  steps {
    downstreamParameterized {
      trigger('tectonic-installer/PR-\${ghprbPullId}') {
        parameters {
          booleanParam('run_smoke_tests', false)
        }
      }
    }
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Installer PR Trigger No Smoke Tests")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
