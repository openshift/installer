#!/usr/bin/env groovy

/* Tips
1. Keep stages focused on producing one artifact or achieving one goal. This makes stages easier to parallelize or re-structure later.
1. Stages should simply invoke a make target or a self-contained script. Do not write testing logic in this Jenkinsfile.
3. CoreOS does not ship with `make`, so Docker builds still have to use small scripts.
*/

def creds = [
  file(credentialsId: 'tectonic-license', variable: 'TF_VAR_tectonic_license_path'),
  file(credentialsId: 'tectonic-pull', variable: 'TF_VAR_tectonic_pull_secret_path'),
  [
    $class: 'AmazonWebServicesCredentialsBinding',
    credentialsId: 'tectonic-jenkins-installer'
  ],
  [
    $class: 'AzureCredentialsBinding',
    credentialsId: 'azure-tectonic-test-service-principal',
    subscriptionIdVariable: 'ARM_SUBSCRIPTION_ID',
    clientIdVariable: 'ARM_CLIENT_ID',
    clientSecretVariable: 'ARM_CLIENT_SECRET',
    tenantIdVariable: 'ARM_TENANT_ID'
  ]
]

def quay_creds = [
  usernamePassword(
    credentialsId: 'quay-robot',
    passwordVariable: 'QUAY_ROBOT_SECRET',
    usernameVariable: 'QUAY_ROBOT_USERNAME'
  )
]

def default_builder_image = 'quay.io/coreos/tectonic-builder:v1.39'
def tectonic_smoke_test_env_image = 'quay.io/coreos/tectonic-smoke-test-env:v5.3'

pipeline {
  agent none
  options {
    timeout(time:120, unit:'MINUTES')
    timestamps()
    buildDiscarder(logRotator(numToKeepStr:'100'))
  }
  parameters {
    string(
      name: 'builder_image',
      defaultValue: default_builder_image,
      description: 'tectonic-builder docker image to use for builds'
    )
    string(
      name: 'hyperkube_image',
      defaultValue: '',
      description: 'Hyperkube image. Please define the param like: {hyperkube="<HYPERKUBE_IMAGE>"}'
    )
    booleanParam(
      name: 'run_smoke_tests',
      defaultValue: true,
      description: ''
    )
  }

  stages {
    stage('Build & Test') {
      environment {
        GO_PROJECT = '/go/src/github.com/coreos/tectonic-installer'
        MAKEFLAGS = '-j4'
      }
      steps {
        node('worker && ec2') {
          withDockerContainer(params.builder_image) {
            ansiColor('xterm') {
              checkout scm
              sh """#!/bin/bash -ex
              mkdir -p \$(dirname $GO_PROJECT) && ln -sf $WORKSPACE $GO_PROJECT

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
              rm -fr frontend/tests_output
              """
              stash name: 'repository'
              cleanWs notFailBuild: true
            }
          }
          withDockerContainer(tectonic_smoke_test_env_image) {
            unstash 'repository'
            sh"""#!/bin/bash -ex
              cd tests/rspec
              bundler exec rubocop --cache false spec lib
            """
            cleanWs notFailBuild: true
          }
        }
      }
    }

    stage('GUI Tests') {
      environment {
        TECTONIC_INSTALLER_ROLE = 'tectonic-installer'
        GRAFITI_DELETER_ROLE = 'grafiti-deleter'
        TF_VAR_tectonic_container_images = "${params.hyperkube_image}"
      }
      steps {
        parallel (
          "IntegrationTest AWS Installer Gui": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(params.builder_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                    cd installer
                    make launch-aws-installer-guitests
                    make gui-aws-tests-cleanup
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "IntegrationTest Baremetal Installer Gui": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(image: params.builder_image, args: '-u root') {
                  ansiColor('xterm') {
                    unstash 'repository'
                    script {
                      try {
                        sh """#!/bin/bash -ex
                        cd installer
                        make launch-baremetal-installer-guitests
                        """
                      }
                      catch (error) {
                        throw error
                      }
                      finally {
                        sh """#!/bin/bash -x
                        cd installer
                        make gui-baremetal-tests-cleanup
                        make clean
                        """
                        cleanWs notFailBuild: true
                      }
                    }
                  }
                }
              }
            }
          }
        )
      }
    }

    stage("Smoke Tests") {
      when {
        expression {
          return params.run_smoke_tests
        }
      }
      environment {
        TECTONIC_INSTALLER_ROLE = 'tectonic-installer'
        GRAFITI_DELETER_ROLE = 'grafiti-deleter'
        TF_VAR_tectonic_container_images = "${params.hyperkube_image}"
      }
      steps {
        parallel (
          "SmokeTest AWS RSpec": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundler exec rspec spec/aws_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest AWS VPC RSpec": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(
                    image: tectonic_smoke_test_env_image,
                    args: '--device=/dev/net/tun --cap-add=NET_ADMIN -u root'
                ) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundler exec rspec spec/aws_vpc_internal_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest AWS Network Policy RSpec": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundler exec rspec spec/aws_network_policy_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest AWS Exp RSpec": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundler exec rspec spec/aws_exp_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest AWS custom ca RSpec": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundler exec rspec spec/aws_ca_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest Azure Basic": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundle exec rspec spec/azure_basic_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest Azure (Experimental)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundle exec rspec spec/azure_experimental_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest Azure Private Cluster": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(
                    image: tectonic_smoke_test_env_image,
                    args: '--device=/dev/net/tun --cap-add=NET_ADMIN -u root'
                ) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundle exec rspec spec/azure_private_external_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
/*
 * Test temporarily disabled
 *
          "SmokeTest: Azure (existing DNS)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                    ansiColor('xterm') {
                      unstash 'repository'
                      sh """#!/bin/bash -ex
                        cd tests/rspec
                        bundle exec rspec spec/azure_dns_spec.rb
                      """
                      cleanWs notFailBuild: true
                    }
                }
              }
            }
          },
*/
          "SmokeTest Azure External Network": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundle exec rspec spec/azure_external_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest: Azure External Network, Experimental)": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundle exec rspec spec/azure_external_experimental_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          "SmokeTest: Azure Example File": {
            node('worker && ec2') {
              withCredentials(creds) {
                withDockerContainer(tectonic_smoke_test_env_image) {
                  ansiColor('xterm') {
                    unstash 'repository'
                    sh """#!/bin/bash -ex
                      cd tests/rspec
                      bundle exec rspec spec/azure_example_spec.rb
                    """
                    cleanWs notFailBuild: true
                  }
                }
              }
            }
          },
          /* Temporarily disabled for consolidation
           * Fails very often due to Packet flakiness
           *
          "SmokeTest: Bare Metal": {
            node('worker && bare-metal') {
              ansiColor('xterm') {
                unstash 'repository'
                withCredentials(creds) {
                  timeout(35) {
                    sh """#!/bin/bash -ex
                    ${WORKSPACE}/tests/smoke/bare-metal/smoke.sh vars/metal.tfvars
                    """
                  }
                  cleanWs notFailBuild: true
                }
              }
            }
          }, */
        )
      }
    }

    stage('Build docker image')  {
      when {
        branch 'master'
      }
      steps {
        node('worker && ec2') {
          withCredentials(quay_creds) {
            ansiColor('xterm') {
              unstash 'repository'
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
  }
}

def notifySlack() {
    def link = "<${env.BUILD_URL}>"
    def msg = "Tectonic Installer failed to destroy azure resources (${link})"
    slackSend(
        channel: '#tectonic-installer-ci',
        color: "warning",
        message: msg,
        teamDomain: 'coreos',
        tokenCredentialId: 'tectonic-slack-token',
        failOnError: true,
    )
}
