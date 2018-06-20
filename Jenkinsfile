#!/usr/bin/env groovy

commonCreds = [
  file(credentialsId: 'tectonic-license', variable: 'LICENSE_PATH'),
  file(credentialsId: 'tectonic-pull', variable: 'PULL_SECRET_PATH'),
  [
    $class: 'StringBinding',
    credentialsId: 'github-coreosbot',
    variable: 'GITHUB_CREDENTIALS'
  ]
]

credsLog = commonCreds.collect()
credsLog.push(
  [
    $class: 'AmazonWebServicesCredentialsBinding',
    credentialsId: 'TF-TECTONIC-JENKINS'
  ]
)

creds = commonCreds.collect()
creds.push(
  [
    $class: 'AmazonWebServicesCredentialsBinding',
    credentialsId: 'TF-TECTONIC-JENKINS-NO-SESSION'
  ]
)

quayCreds = [
  usernamePassword(
    credentialsId: 'quay-robot',
    passwordVariable: 'QUAY_ROBOT_SECRET',
    usernameVariable: 'QUAY_ROBOT_USERNAME'
  )
]

tectonicSmokeTestEnvImage = 'quay.io/coreos/tectonic-smoke-test-env:v6.0'
tectonicBazelImage = 'quay.io/coreos/tectonic-builder:bazel-v0.3'
originalCommitId = 'UNKNOWN'

pipeline {
  agent { label 'worker && ec2' }
  options {
    // Individual steps have stricter timeouts. 360 minutes should be never reached.
    timeout(time:6, unit:'HOURS')
    timestamps()
    buildDiscarder(logRotator(numToKeepStr:'20', artifactNumToKeepStr: '20'))
  }
  parameters {
    booleanParam(
      name: 'RUN_SMOKE_TESTS',
      defaultValue: false,
      description: ''
    )
    booleanParam(
      name: 'PLATFORM/AWS',
      defaultValue: false,
      description: ''
    )
    booleanParam(
      name: 'NOTIFY_SLACK',
      defaultValue: false,
      description: 'Notify slack channel on failure.'
    )
    string(
      name: 'SLACK_CHANNEL',
      defaultValue: '#team-installer',
      description: 'Slack channel to notify on failure.'
    )
  }

  stages {
    stage("Smoke Tests") {
      when {
        anyOf {
          branch "master"
          expression { return params.RUN_SMOKE_TESTS }
        }
      }
      options {
        timeout(time: 70, unit: 'MINUTES')
      }
      steps {
        withDockerContainer(tectonicSmokeTestEnvImage) {
          withCredentials(creds) {
            ansiColor('xterm') {
              sh """#!/bin/bash -e
                   export HOME=/home/jenkins
                   ./tests/run.sh
                   cp bazel-bin/tectonic-dev.tar.gz .
                 """
              // Produce an artifact which can be downloaded via web UI
              stash name: 'tectonic-tarball', includes: 'tectonic-dev.tar.gz'
            }
          }
        }
      }
    }

    stage('Build docker image')  {
      when {
        branch 'master'
      }
      steps {
        withCredentials(quayCreds) {
          ansiColor('xterm') {
            sh """
              docker build -t quay.io/coreos/tectonic-installer:master -f images/tectonic-installer/Dockerfile .
              docker login -u="$QUAY_ROBOT_USERNAME" -p="$QUAY_ROBOT_SECRET" quay.io
              docker push quay.io/coreos/tectonic-installer:master
              docker logout quay.io
            """
            cleanWs notFailBuild: true
          }
        }
      }
    }

  }
  post {
    always {
      forcefullyCleanWorkspace()
      cleanWs notFailBuild: true
    }
    failure {
      script {
        if (params.NOTIFY_SLACK) {
          slackSend color: 'danger', channel: params.SLACK_CHANNEL, message: "Job ${env.JOB_NAME}, build no. #${BUILD_NUMBER} failed! (<${env.BUILD_URL}|Open>)"
        }
      }
    }
  }
}

def forcefullyCleanWorkspace() {
  return withDockerContainer(
    image: tectonicBazelImage,
    args: '-u root'
  ) {
    ansiColor('xterm') {
      sh """#!/bin/bash -e
        if [ -d "\$WORKSPACE" ]
        then
          rm -rf \$WORKSPACE/*
        fi
      """
    }
  }
}
