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

def quay_creds = [
  usernamePassword(
    credentialsId: 'quay-robot',
    passwordVariable: 'QUAY_ROBOT_SECRET',
    usernameVariable: 'QUAY_ROBOT_USERNAME'
  )
]

def builder_image = 'quay.io/coreos/tectonic-builder:v1.13'

pipeline {
  agent none
  options {
    timeout(time:60, unit:'MINUTES')
    timestamps()
    buildDiscarder(logRotator(numToKeepStr:'20'))
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
            stash name: 'sanity', includes: 'installer/bin/sanity'
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
                withDockerContainer(builder_image) {
                  checkout scm
                  unstash 'installer'
                  unstash 'sanity'
                  timeout(30) {
                    sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                    sh '${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws.tfvars'
                    sh '${WORKSPACE}/tests/smoke/aws/smoke.sh create vars/aws.tfvars'
                    sh '${WORKSPACE}/tests/smoke/aws/smoke.sh test vars/aws.tfvars'
                    sh '${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws.tfvars'
                  }
                }
              }
            }
          },
          "SmokeTest TerraForm: AWS (Experimental)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(builder_image) {
                  checkout scm
                  unstash 'installer'
                  timeout(5) {
                    sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                    sh '${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-exp.tfvars'
                  }
                }
              }
            }
          },
          "SmokeTest TerraForm: AWS (custom-ca)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(builder_image) {
                  checkout scm
                  unstash 'installer'
                  timeout(5) {
                    sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                    sh '${WORKSPACE}/tests/smoke/aws/smoke.sh plan vars/aws-ca.tfvars'
                  }
                }
              }
            }
          },
          "SmokeTest Terraform: Bare Metal": {
            node('worker && bare-metal') {
              checkout scm
              unstash 'installer'
              unstash 'sanity'
              withCredentials(creds) {
                timeout(30) {
                  sh '${WORKSPACE}/tests/smoke/bare-metal/smoke.sh vars/metal.tfvars'
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
            withCredentials(creds) {
              withDockerContainer(builder_image) {
                checkout scm
                unstash 'installer'
                sh 'set +x -e && eval "$(${WORKSPACE}/tests/smoke/aws/smoke.sh assume-role "$TECTONIC_INSTALLER_ROLE")"'
                retry(3) {
                  timeout(15) {
                    sh '${WORKSPACE}/tests/smoke/aws/smoke.sh destroy vars/aws.tfvars'
                  }
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
