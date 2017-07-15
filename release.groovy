// Associated job: https://jenkins-tectonic-installer.prod.coreos.systems/job/tectonic-release/

/* Tips
1. Keep stages focused on producing one artifact or achieving one goal. This makes stages easier to parallelize or re-structure later.
1. Stages should simply invoke a make target or a self-contained script. Do not write testing logic in this Jenkinsfile.
3. CoreOS does not ship with `make`, so Docker builds still have to use small scripts.
*/

def creds = [
  [
    $class: 'UsernamePasswordMultiBinding',
    credentialsId: 'tectonic-release-s3-upload',
    usernameVariable: 'AWS_ACCESS_KEY_ID',
    passwordVariable: 'AWS_SECRET_ACCESS_KEY'
  ],
  [
    $class: 'StringBinding',
    credentialsId: 'github-coreosbot',
    variable: 'GITHUB_CREDENTIALS'
  ]
]

def builder_image = 'quay.io/coreos/tectonic-builder:v1.33'

pipeline {
  agent none

  options {
    timeout(time:25, unit:'MINUTES')
  }

  parameters {
    string(name: 'releaseTag')
    string(name: 'preRelease')
  }

  stages {
    stage('Release') {
      agent none
      environment {
        GO_PROJECT = '/go/src/github.com/coreos/tectonic-installer'
        MAKEFLAGS = '-j4'
      }
      steps {
        node('worker && ec2') {
          withCredentials(creds) {
            withDockerContainer(builder_image) {
              checkout([$class: 'GitSCM', branches: [[name: "refs/tags/${params.releaseTag}"]], userRemoteConfigs: [[credentialsId: 'github-coreosbot', url: 'https://github.com/coreos/tectonic-installer.git']]])
              sh """#!/bin/bash -ex
                mkdir -p \$(dirname $GO_PROJECT) && ln -sf $WORKSPACE $GO_PROJECT

                export VERSION=${params.releaseTag}
                export PRE_RELEASE=${params.preRelease}
                go version

                # TODO: Remove me.
                go get github.com/segmentio/terraform-docs
                go get github.com/s-urbaniak/terraform-examples

                cd $GO_PROJECT/
                make structure-check

                cd $GO_PROJECT/installer
                make clean
                make tools
                make build

                make dirtycheck
                make lint
                make test

                make release
              """
            }
          }
        }
      }
    }
  }
}
