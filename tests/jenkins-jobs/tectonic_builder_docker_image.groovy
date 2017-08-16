#!/bin/env groovy​

folder("builders")

job("builders/tectonic-builder-docker-image") {
  logRotator(-1, 10)
  description('Build quay.io/coreos/tectonic-builder Docker image. Changes here will be reverted automaticaly.')

  label 'worker&&ec2'

  parameters {
    stringParam('TERRAFORM_UPSTREAM_URL', '', 'upstream Terraform download url, defaults to CoreOS custom Terraform release (github.com/coreos/terraform)')
    stringParam('TECTONIC_BUILDER_TAG', '', 'Tectonic Builder Docker tag')
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
    def cmd = """
    #!/bin/bash -e

    if [ -z "\${TERRAFORM_UPSTREAM_URL}" ]
    then
      export TECTONIC_BUILDER_IMAGE=quay.io/coreos/tectonic-builder:\${TECTONIC_BUILDER_TAG}
      docker build -t \${TECTONIC_BUILDER_IMAGE} -f images/tectonic-builder/Dockerfile .
    else
      export TECTONIC_BUILDER_IMAGE=quay.io/coreos/tectonic-builder:\${TECTONIC_BUILDER_TAG}-upstream-terraform
      docker build -t \${TECTONIC_BUILDER_IMAGE} --build-arg TERRAFORM_URL=\${TERRAFORM_UPSTREAM_URL} -f images/tectonic-builder/Dockerfile .
    fi

    if [ !\${DRY_RUN} = true  ] ; then
      docker login quay.io -u \${QUAY_USERNAME} -p \${QUAY_PASSWD}
      docker push \${TECTONIC_BUILDER_IMAGE}
    fi
  """.stripIndent()
    shell(cmd)
  }

  publishers {
    wsCleanup()
    slackNotifier {
      authTokenCredentialId('tectonic-slack-token')
      customMessage("Jenkins Builder: tectonic-builder")
      includeCustomMessage(true)
      notifyBackToNormal(true)
      notifyFailure(true)
      room('#tectonic-installer-ci')
      teamDomain('coreos')
    }
  }
}
