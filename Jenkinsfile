/* Tips
1. Keep stages focused on producing one artifact or achieving one goal. This makes stages easier to parallelize or re-structure later.
1. Stages should simply invoke a make target or a self-contained script. Do not write testing logic in this Jenkinsfile.
3. CoreOS does not ship with `make`, so Docker builds still have to use small scripts.
*/

def creds = [
  file(credentialsId: 'tectonic-license', variable: 'TF_VAR_tectonic_license_path'),
  file(credentialsId: 'tectonic-pull', variable: 'TF_VAR_tectonic_pull_secret_path'), [
    $class: 'UsernamePasswordMultiBinding',
    credentialsId: 'jenkins-tectonic-installer',
    usernameVariable: 'AWS_ACCESS_KEY_ID',
    passwordVariable: 'AWS_SECRET_ACCESS_KEY'
  ]
]

def quay_creds = [
  usernamePassword(
    credentialsId: 'quay-robot',
    passwordVariable: 'QUAY_ROBOT_SECRET',
    usernameVariable: 'QUAY_ROBOT_USERNAME'
  )
]

def builder_image = 'quay.io/coreos/tectonic-builder:v1.22'

pipeline {
  agent none
  options {
    timeout(time:70, unit:'MINUTES')
    timestamps()
    buildDiscarder(logRotator(numToKeepStr:'100'))
  }

  stages {
    stage('Build & Test') {
      environment {
        GO_PROJECT = '/go/src/github.com/coreos/tectonic-installer'
        MAKEFLAGS = '-j4'
      }
      steps {
        node('worker && ec2') {
          withDockerContainer(builder_image) {
            checkout scm
            sh """#!/bin/bash -ex
            mkdir -p \$(dirname $GO_PROJECT) && ln -sf $WORKSPACE $GO_PROJECT

            # TODO: Remove me.
            go get github.com/segmentio/terraform-docs
            go get github.com/s-urbaniak/terraform-examples

            cd $GO_PROJECT/
            make structure-check
            make bin/smoke

            cd $GO_PROJECT/installer
            make clean
            make tools
            make build

            make dirtycheck
            make lint
            make test
            """
            stash name: 'installer', includes: 'installer/bin/linux/installer'
            stash name: 'node_modules', includes: 'installer/frontend/node_modules/**'
            stash name: 'smoke', includes: 'bin/smoke'
          }
        }
      }
    }

    stage("Tests") {
      environment {
        TECTONIC_INSTALLER_ROLE = 'tectonic-installer'
      }
      steps {
        parallel (
          "SmokeTest TerraForm: AWS": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(args: '-v /etc/passwd:/etc/passwd:ro', image: builder_image) {
                  checkout scm
                  unstash 'installer'
                  unstash 'smoke'
                  timeout(35) {
                    sh """#!/bin/bash -ex
                    . ${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE"
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-tls.tfvars
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh create vars/aws-tls.tfvars
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh test vars/aws-tls.tfvars
                    """
                  }
                  catchError {
                    timeout (5) {
                      sshagent(['aws-smoke-test-ssh-key']) {
                        sh """#!/bin/bash
                        # Running without -ex because we don't care if this fails
                        . ${WORKSPACE}/tests/smoke/aws/smoke.sh common vars/aws-tls.tfvars
                        ${WORKSPACE}/tests/smoke/aws/cluster-foreach.sh ${WORKSPACE}/tests/smoke/forensics.sh
                        """
                      }
                    }
                  }
                  retry(3) {
                    timeout(15) {
                      sh """#!/bin/bash -ex
                      . ${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE"
                      ${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws-tls.tfvars
                      """
                    }
                  }
                }
              }
            }
          },
          "SmokeTest TerraForm: AWS (non-TLS)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(builder_image) {
                  checkout scm
                  unstash 'installer'
                  timeout(5) {
                    sh """#!/bin/bash -ex
                    . ${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE"
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws.tfvars
                    """
                  }
                }
              }
            }
          },
          "SmokeTest TerraForm: AWS (experimental)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(builder_image) {
                  checkout scm
                  unstash 'installer'
                  timeout(5) {
                    sh """#!/bin/bash -ex
                    . ${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE"
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-exp.tfvars
                    """
                  }
                }
              }
            }
          },
          "SmokeTest TerraForm: AWS (custom ca)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(builder_image) {
                  checkout scm
                  unstash 'installer'
                  timeout(5) {
                    sh """#!/bin/bash -ex
                    . ${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE"
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-ca.tfvars
                    """
                  }
                }
              }
            }
          },
          "SmokeTest TerraForm: AWS (private vpc)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(image: builder_image, args: '--device=/dev/net/tun --cap-add=NET_ADMIN') {
                  checkout scm
                  unstash 'installer'
                  unstash 'smoke'
                  timeout(45) {
                    sh """#!/bin/bash -ex
                    . ${WORKSPACE}/tests/smoke/aws/smoke.sh create-vpc
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-vpc.tfvars
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh create vars/aws-vpc.tfvars
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh test vars/aws-vpc.tfvars
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws-vpc.tfvars
                    ${WORKSPACE}/tests/smoke/aws/smoke.sh destroy-vpc
                    """
                  }
                }
              }
            }
          },
          "SmokeTest Terraform: Bare Metal": {
            node('worker && bare-metal') {
              checkout scm
              unstash 'installer'
              unstash 'smoke'
              withCredentials(creds) {
                timeout(35) {
                  sh """#!/bin/bash -ex
                  ${WORKSPACE}/tests/smoke/bare-metal/smoke.sh vars/metal.tfvars
                  """
                }
              }
            }
          },
          "IntegrationTest Installer Gui": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(builder_image) {
                  checkout scm
                  unstash 'installer'
                  unstash 'node_modules'
                  sh """#!/bin/bash -ex
                  cd installer
                  make launch-installer-guitests
                  make gui-tests-cleanup
                  """
                }
              }
            }
          }
        )
      }
      post {
        failure {
          node('worker && ec2') {
            withCredentials([creds, quay_creds]) {
              withDockerContainer(builder_image) {
                checkout scm
                unstash 'installer'
                timeout(15) {
                  sh """#!/bin/bash -x
                  ${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws.tfvars
                  """
                }
                timeout(20) {
                  sh """#!/bin/bash -x
                  ${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws-vpc.tfvars
                  ${WORKSPACE}/tests/smoke/aws/smoke.sh destroy-vpc
                  """
                }
                timeout(10) {
                  sh """#!/bin/bash -x
                  docker login -u="$QUAY_ROBOT_USERNAME" -p="$QUAY_ROBOT_SECRET" quay.io
                  ${WORKSPACE}/tests/smoke/aws/smoke.sh grafiti-clean vars/aws.tfvars
                  ${WORKSPACE}/tests/smoke/aws/smoke.sh grafiti-clean vars/aws-vpc.tfvars
                  """
                }
              }
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
        node('worker && ec2') {
          withCredentials(quay_creds) {
            checkout scm
            unstash 'installer'
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
  }
}
