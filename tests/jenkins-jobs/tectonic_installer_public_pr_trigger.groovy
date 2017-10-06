#!/bin/env groovy​

folder("triggers")

job("triggers/tectonic-installer-pr-trigger") {
  description('Tectonic Installer PR Trigger. Changes here will be reverted automatically.')

  concurrentBuild()

  logRotator(30, 100)
  label("worker && ec2")

  properties {
    githubProjectUrl('https://github.com/coreos/tectonic-installer')
  }

  wrappers {
    colorizeOutput()
    timestamps()
    buildInDocker {
      image('quay.io/coreos/tectonic-smoke-test-env:v5.6')
    }
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
      blackListLabels("")
      whiteListLabels("")
      includedRegions("")
      excludedRegions("")
    }
  }

  steps {
    shell """#!/bin/bash -ex
      curl "https://api.github.com/repos/coreos/tectonic-installer/labels" > repoLabels
      repoLabels=\$(jq ".[] | .name" repoLabels)
      repoLabels=\$(echo \$repoLabels | tr -d "\\"" | tr [a-z] [A-Z] | tr - _)
      for label in \$repoLabels
      do
        echo \$label=false >> env_vars
      done


      curl "https://api.github.com/repos/coreos/tectonic-installer/issues/\${ghprbPullId}" > pr
      labels=\$(jq ".labels | .[] | .name" pr)
      labels=\$(echo \$labels | tr -d "\\"" | tr [a-z] [A-Z] | tr - _)
      for label in \$labels
      do
        echo \$label=true >> env_vars
      done
    """

    downstreamParameterized {
      trigger('tectonic-installer/PR-\${ghprbPullId}') {
        parameters {
          propertiesFile("env_vars", true)
        }
      }
    }
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Installer PR Trigger")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      notifyRepeatedFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
