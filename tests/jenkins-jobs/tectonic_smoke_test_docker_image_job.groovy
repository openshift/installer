#!/bin/env groovy​

folder("builders")

job("builders/tectonic-smoke-env-docker-image") {
  logRotator(-1, 10)
  description('Build quay.io/coreos/tectonic-smoke-test-env Docker image. Changes here will be reverted automatically.')

  label 'worker&&ec2'

  parameters {
    stringParam('TECTONIC_SMOKE_TAG', '', 'Tectonic Smoke Docker tag')
    booleanParam('DRY_RUN', true, 'Just build the docker image')
  }

  wrappers {
    colorizeOutput()
    timestamps()
    credentialsBinding {
      usernamePassword("QUAY_USERNAME", "QUAY_PASSWD", "quay-robot")
    }
  }

  scm {
    git {
      remote {
        url('https://github.com/coreos/tectonic-installer')
      }
      branch('origin/master')
    }
  }


  steps {
    def cmd = """#!/bin/bash -ex
export TECTONIC_SMOKE_IMAGE=quay.io/coreos/tectonic-smoke-test-env:\${TECTONIC_SMOKE_TAG}
docker build -t \${TECTONIC_SMOKE_IMAGE} -f images/tectonic-smoke-test-env/Dockerfile .

if \${DRY_RUN};
then
  echo "Just build the image"
else
  echo "Pushing the Image to quay"
  docker login quay.io -u \${QUAY_USERNAME} -p \${QUAY_PASSWD}
  docker push \${TECTONIC_SMOKE_IMAGE}
fi
  """.stripIndent()
    shell(cmd)
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Jenkins Builder: tectonic-smoke-test-env - tag: \${TECTONIC_SMOKE_TAG}")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      notifyRepeatedFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
