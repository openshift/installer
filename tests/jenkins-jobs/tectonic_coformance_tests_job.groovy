#!/bin/env groovy​

folder("tectonic-tests")

job("tectonic-tests/tectonic-conformance-test") {
  description("This job runs the Conformance tests for the tectonic-installer.\nThis job is manage by tectonic-installer.\nChanges here will be reverted automatically")
  logRotator(30, 100)
  parameters {
      stringParam('BRANCH', 'master', 'Which branch to run the conformance tests.')
  }

  label('worker&&ec2')

  scm {
    git {
      remote {
        url('https://github.com/coreos/tectonic-installer')
      }
      branch("origin/\${BRANCH}")
    }
  }

  wrappers {
    colorizeOutput()
    timestamps()
    credentialsBinding {
      usernamePassword("QUAY_ROBOT_USERNAME", "QUAY_ROBOT_SECRET", "quay-robot")
      usernamePassword("TF_VAR_tectonic_admin_email", "TF_VAR_tectonic_admin_password_hash", "tectonic-console-login")
      amazonWebServicesCredentialsBinding {
        accessKeyVariable("AWS_ACCESS_KEY_ID")
        secretKeyVariable("AWS_SECRET_ACCESS_KEY")
        credentialsId("tectonic-jenkins-installer")
      }
      file("TF_VAR_tectonic_pull_secret_path", "tectonic-pull")
      file("TF_VAR_tectonic_license_path", "tectonic-license")
    }
    sshAgent('aws-smoke-test-ssh-key')
  }

  triggers {
    cron('H 0 * * *')
  }

  steps {
    def cmd = """#!/bin/bash -ex
./tests/conformance/conformance.sh
    """.stripIndent()
    shell(cmd)
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Tectonic Conformance Tests")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      notifySuccess(true)
      notifyRepeatedFailure(true)
      room('#forum-installer')
      teamDomain('coreos')
    }
  }
}
