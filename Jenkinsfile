#!/usr/bin/env groovy

/* Tips
1. Keep stages focused on producing one artifact or achieving one goal. This makes stages easier to parallelize or re-structure later.
1. Stages should simply invoke a make target or a self-contained script. Do not write testing logic in this Jenkinsfile.
3. CoreOS does not ship with `make`, so Docker builds still have to use small scripts.
*/

commonCreds = [
  file(credentialsId: 'tectonic-license', variable: 'TF_VAR_tectonic_license_path'),
  file(credentialsId: 'tectonic-pull', variable: 'TF_VAR_tectonic_pull_secret_path'),
  usernamePassword(
    credentialsId: 'jenkins-log-analyzer-user',
    passwordVariable: 'LOG_ANALYZER_PASSWORD',
    usernameVariable: 'LOG_ANALYZER_USER'
  ),
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
creds.push(
  string(credentialsId: 'AWS-TECTONIC-TRACK-2-ROLE-NAME', variable: 'TF_VAR_tectonic_aws_installer_role')
)

quayCreds = [
  usernamePassword(
    credentialsId: 'quay-robot',
    passwordVariable: 'QUAY_ROBOT_SECRET',
    usernameVariable: 'QUAY_ROBOT_USERNAME'
  )
]

defaultBuilderImage = 'quay.io/coreos/tectonic-builder:v1.45'
tectonicSmokeTestEnvImage = 'quay.io/coreos/tectonic-smoke-test-env:v5.16'
tectonicBazelImage = 'quay.io/coreos/tectonic-builder:bazel-v0.3'
originalCommitId = 'UNKNOWN'

pipeline {
  agent none
  environment {
    KUBE_CONFORMANCE_IMAGE = 'quay.io/coreos/kube-conformance:v1.9.1_coreos.0'
    LOGSTASH_BUCKET= "log-analyzer-tectonic-installer"
  }
  options {
    // Individual steps have stricter timeouts. 360 minutes should be never reached.
    timeout(time:6, unit:'HOURS')
    timestamps()
    buildDiscarder(logRotator(numToKeepStr:'20', artifactNumToKeepStr: '20'))
  }
  parameters {
    string(
      name: 'builder_image',
      defaultValue: defaultBuilderImage,
      description: 'tectonic-builder docker image to use for builds'
    )
    string(
      name: 'hyperkube_image',
      defaultValue: '',
      description: 'Hyperkube image. Please define the param like: {hyperkube="<HYPERKUBE_IMAGE>"}'
    )
    booleanParam(
      name: 'RUN_CONFORMANCE_TESTS',
      defaultValue: false,
      description: ''
    )
    booleanParam(
      name: 'RUN_SMOKE_TESTS',
      defaultValue: true,
      description: ''
    )
    booleanParam(
      name: 'RUN_GUI_TESTS',
      defaultValue: true,
      description: ''
    )
    string(
      name: 'COMPONENT_TEST_IMAGES',
      defaultValue: '',
      description: 'List of container images for component tests to run (comma-separated)'
    )
    booleanParam(
      name: 'PLATFORM/AWS',
      defaultValue: true,
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
    string(
      name: 'GITHUB_REPO',
      defaultValue: 'coreos/tectonic-installer',
      description: 'Github repository'
    )
    string(
      name: 'SPECIFIC_GIT_COMMIT',
      description: 'Checkout a specific git ref (e.g. sha or tag). If not set, Jenkins uses the most recent commit of the triggered branch.',
      defaultValue: '',
    )
  }

  stages {
    stage('Build & Test') {
      environment {
        GO_PROJECT = "/go/src/github.com/${params.GITHUB_REPO}"
        MAKEFLAGS = '-j4'
      }
      steps {
        node('worker && ec2') {
          ansiColor('xterm') {
            script {
              def err = null
              try {
                timeout(time: 20, unit: 'MINUTES') {
                  forcefullyCleanWorkspace()

                  /*
                    This supports users who require builds at a specific git ref
                    instead of the branch tip.
                  */
                  if (params.SPECIFIC_GIT_COMMIT == '') {
                    checkout scm
                    originalCommitId = sh(returnStdout: true, script: 'git rev-parse "origin/${BRANCH_NAME}"').trim()
                  } else {
                    checkout([
                      $class: 'GitSCM',
                      branches: [[name: params.SPECIFIC_GIT_COMMIT]],
                      userRemoteConfigs: [[url: "https://github.com/${params.GITHUB_REPO}.git"]]
                    ])
                    // In case params.SPECIFIC_GIT_COMMIT is a mutable tag instead
                    // of a sha
                    originalCommitId = sh(returnStdout: true, script: 'git rev-parse "${SPECIFIC_GIT_COMMIT}"').trim()
                  }

                  echo "originalCommitId: ${originalCommitId}"
                  stash name: 'clean-repo', excludes: 'installer/vendor/**,tests/smoke/vendor/**'

                  withDockerContainer(tectonicBazelImage) {
                    sh "bazel test terraform_fmt --test_output=all"
                    sh "bazel test installer:cli_units --test_output=all"
                    sh"""#!/bin/bash -ex
                      bazel build tarball tests/smoke

                      # Jenkins `stash` does not follow symlinks - thereby temporarily copy the files to the root dir
                      cp bazel-bin/tectonic.tar.gz .
                      cp bazel-bin/tests/smoke/linux_amd64_stripped/smoke .
                    """
                    stash name: 'tectonic.tar.gz', includes: 'tectonic.tar.gz'
                    stash name: 'smoke-tests', includes: 'smoke'
                    archiveArtifacts allowEmptyArchive: true, artifacts: 'tectonic.tar.gz'
                  }

                  withDockerContainer(tectonicSmokeTestEnvImage) {
                    sh"""#!/bin/bash -ex
                      cd tests/rspec
                      rubocop --cache false spec lib
                    """
                  }
                }
              } catch (error) {
                err = error
                throw error
              } finally {
                reportStatusToGithub((err == null) ? 'success' : 'failure', 'basic-tests', originalCommitId)
              }
            }
          }
        }
      }
    }

    stage("Smoke Tests") {
      when {
        expression {
          return params.RUN_SMOKE_TESTS || params.RUN_CONFORMANCE_TESTS || params.COMPONENT_TEST_IMAGES != ''
        }
      }
      environment {
        TECTONIC_INSTALLER_ROLE = 'tf-tectonic-installer-track-2'
        GRAFITI_DELETER_ROLE = 'tf-grafiti'
        TF_VAR_tectonic_container_images = "${params.hyperkube_image}"
        TF_VAR_tectonic_kubelet_debug_config = "--minimum-container-ttl-duration=8h --maximum-dead-containers-per-container=9999 --maximum-dead-containers=9999"
      }
      steps {
        script {
          def builds = [:]
          def aws = [
            [file: 'basic_spec.rb', args: ''],
            // [file: 'vpc_internal_spec.rb', args: '--device=/dev/net/tun --cap-add=NET_ADMIN -u root'],
            // [file: 'network_flannel_spec.rb', args: ''],
            // [file: 'exp_spec.rb', args: ''],
            // [file: 'ca_spec.rb', args: ''],
            // [file: 'custom_tls_spec.rb', args: '']
          ]

          if (params."PLATFORM/AWS") {
            aws.each { build ->
              filepath = 'spec/aws/' + build.file
              builds['aws/' + build.file] = runRSpecTest(filepath, build.args, creds)
            }
          }

          parallel builds
        }
      }
    }

    stage('Build docker image')  {
      when {
        branch 'master'
      }
      steps {
        node('worker && ec2') {
          forcefullyCleanWorkspace()
          withCredentials(quayCreds) {
            ansiColor('xterm') {
              unstash 'clean-repo'
              unstash 'tectonic.tar.gz'
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
  post {
    always {
      node('worker && ec2') {
        forcefullyCleanWorkspace()
        echo "Starting with streaming the logfile to the S3 bucket"
        withDockerContainer(params.builder_image) {
          withCredentials(credsLog) {
            unstash 'clean-repo'
            script {
              try {
                sh """#!/bin/bash -xe
                export BUILD_RESULT=${currentBuild.currentResult}
                ./tests/jenkins-jobs/scripts/log-analyzer-copy.sh jenkins-logs
                """
              } catch (Exception e) {
                notifyBuildSlack()
              } finally {
                cleanWs notFailBuild: true
              }
            }
          }
        }
      }
    }

    failure {
      script {
        notifyBuildSlack()
      }
    }
  }
}

def forcefullyCleanWorkspace() {
  return withDockerContainer(
    image: tectonicSmokeTestEnvImage,
    args: '-u root'
  ) {
    ansiColor('xterm') {
      sh """#!/bin/bash -e
        if [ -d "\$WORKSPACE" ]
        then
          rm -rfv \$WORKSPACE/*
        fi
      """
    }
  }
}

def unstashCleanRepoTectonicTarGZSmokeTests() {
  unstash 'clean-repo'
  unstash 'tectonic.tar.gz'
  unstash 'smoke-tests'
  sh """#!/bin/bash -ex
    # Jenkins `stash` does not follow symlinks - thereby temporarily copy the files to the root dir
    mkdir -p bazel-bin/tests/smoke/linux_amd64_stripped/
    cp tectonic.tar.gz bazel-bin/.
    cp smoke bazel-bin/tests/smoke/linux_amd64_stripped/.
  """
}

def runRSpecTest(testFilePath, dockerArgs, credentials) {
  return {
    node('worker && ec2') {
      def err = null
      try {
        timeout(time: 5, unit: 'HOURS') {
          forcefullyCleanWorkspace()
          ansiColor('xterm') {
            withCredentials(credentials + quayCreds) {
              withDockerContainer(
                image: tectonicSmokeTestEnvImage,
                args: '-u root -v /var/run/docker.sock:/var/run/docker.sock ' + dockerArgs
              ) {
                unstashCleanRepoTectonicTarGZSmokeTests()
                sh """#!/bin/bash -ex
                  mkdir -p templogfiles && chmod 777 templogfiles
                  cd tests/rspec

                  # Directing test output both to stdout as well as a log file
                  rspec ${testFilePath} --format RspecTap::Formatter --format RspecTap::Formatter --out ../../templogfiles/tap.log
                """
              }
            }
          }
        }
      } catch (error) {
        err = error
        throw error
      } finally {
        reportStatusToGithub((err == null) ? 'success' : 'failure', testFilePath, originalCommitId)
        step([$class: "TapPublisher", testResults: "templogfiles/*", outputTapToConsole: true, planRequired: false])
        archiveArtifacts allowEmptyArchive: true, artifacts: 'bazel-bin/tectonic/**/logs/**'
        withDockerContainer(params.builder_image) {
         withCredentials(credsLog) {
          script {
            try {
              sh """#!/bin/bash -xe
              ./tests/jenkins-jobs/scripts/log-analyzer-copy.sh smoke-test-logs ${testFilePath}
              """
            } catch (Exception e) {
              notifyBuildSlack()
            } finally {
              cleanWs notFailBuild: true
            }
          }
         }
        }
        cleanWs notFailBuild: true
      }

    }
  }
}

def reportStatusToGithub(status, context, commitId) {
  withCredentials(creds) {
    sh """#!/bin/bash -ex
      ./tests/jenkins-jobs/scripts/report-status-to-github.sh ${status} ${context} ${commitId} ${params.GITHUB_REPO}
    """
  }
}

def notifyBuildSlack() {
  if (!params.NOTIFY_SLACK) {
    return
  }
  def colorCode = '#FF0000'
  def subject = "Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'"
  def summary = "${subject} (${env.BUILD_URL})"

  // Send notifications
  echo 'Sending notification to slack...'
  slackSend(color: colorCode, message: summary, channel: params.SLACK_CHANNEL)
  echo 'Slack notification sent.'
}
