/* Tips
1. Keep stages focused on producing one artifact or achieving one goal. This makes stages easier to parallelize or re-structure later.
1. Stages should simply invoke a make target or a self-contained script. Do not write testing logic in this Jenkinsfile.
3. CoreOS does not ship with `make`, so Docker builds still have to use small scripts.
*/

def creds = [
  file(credentialsId: 'tectonic-license', variable: 'TF_VAR_tectonic_license_path'),
  file(credentialsId: 'tectonic-pull', variable: 'TF_VAR_tectonic_pull_secret_path'), [
    $class: 'UsernamePasswordMultiBinding',
    credentialsId: 'tectonic-aws',
    usernameVariable: 'AWS_ACCESS_KEY_ID',
    passwordVariable: 'AWS_SECRET_ACCESS_KEY'
  ]
]

pipeline {
  agent {
    label 'worker'
  }

  options {
    timeout(time:60, unit:'MINUTES')
    timestamps()
    buildDiscarder(logRotator(numToKeepStr:'20'))
  }

  environment {
    TECTONIC_INSTALLER_ROLE= 'tectonic-installer'
    GO_PROJECT = '/go/src/github.com/coreos/tectonic-installer'
    MAKEFLAGS = '-j4'
  }

  stages {
    stage('Installer: Build & Test') {
      agent {
        docker {
          image 'quay.io/coreos/tectonic-builder:v1.12'
        }
      }
      steps {
        checkout scm
        sh "mkdir -p \$(dirname $GO_PROJECT) && ln -sf $WORKSPACE $GO_PROJECT"
        sh "go get github.com/golang/lint/golint"
        sh """#!/bin/bash -ex
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
        """
        stash name: 'installer', includes: 'installer/bin/linux/installer'
        stash name: 'sanity', includes: 'installer/bin/sanity'
      }
    }

    stage("Smoke Tests") {
      agent {
        docker {
          image 'quay.io/coreos/tectonic-builder:v1.12'
        }
      }
      steps {
        parallel (
          "TerraForm: AWS": {
            unstash 'installer'
            unstash 'sanity'
            withCredentials(creds) {
              timeout(30) {
                sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                sh '${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws.tfvars'
                sh '${WORKSPACE}/tests/smoke/aws/smoke.sh create vars/aws.tfvars'
                sh '${WORKSPACE}/tests/smoke/aws/smoke.sh test vars/aws.tfvars'
              }
              timeout(10) {
                sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                sh '${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws.tfvars'
              }
            }
          },
          "TerraForm: AWS-experimental": {
            unstash 'installer'
            unstash 'sanity'
            withCredentials(creds) {
              timeout(5) {
                sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                sh '${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-exp.tfvars'
              }
            }
          },
          "TerraForm: AWS-custom-ca": {
            unstash 'installer'
            unstash 'sanity'
            withCredentials(creds) {
              timeout(5) {
                sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                sh '${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-ca.tfvars'
              }
            }
          }
        )
      }
      post {
        failure {
          unstash 'installer'
          withCredentials(creds) {
            timeout(10) {
              sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
              sh '${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws.tfvars'
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
        unstash 'installer'
        unstash 'sanity'
        withCredentials([
            usernamePassword(
              credentialsId: 'quay-robot',
              passwordVariable: 'QUAY_ROBOT_SECRET',
              usernameVariable: 'QUAY_ROBOT_USERNAME'
            )
          ]) {
          sh """
            docker build -t quay.io/coreos/tectonic-installer:master -f images/tectonic-installer/Dockerfile .
            docker login -u="$QUAY_ROBOT_USERNAME" -p="$QUAY_ROBOT_SECRET" quay.io
            docker push quay.io/coreos/tectonic-installer:master
            docker logout quay.io
          """
        }
      }
    }
  }

  post {
    always {
      // Cleanup workspace
      deleteDir()
    }
  }
}
