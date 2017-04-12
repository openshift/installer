pipeline {
  agent {
    docker {
      image 'quay.io/coreos/tectonic-terraform:v0.0.2'
      label 'worker'
    }
  }

  options {
    timeout(time:35, unit:'MINUTES')
    buildDiscarder(logRotator(numToKeepStr:'20'))
  }

  stages {
    stage('Syntax Check') {
      steps {
        sh 'make structure-check'
      }
    }

    stage('Smoke Tests') {
      steps {
        parallel (
          "AWS": {
              withCredentials([file(credentialsId: 'tectonic-license', variable: 'TF_VAR_tectonic_pull_secret_path'),
                               file(credentialsId: 'tectonic-pull', variable: 'TF_VAR_tectonic_license_path'),
                               [
                                 $class: 'UsernamePasswordMultiBinding',
                                 credentialsId: 'tectonic-aws-creds',
                                 usernameVariable: 'AWS_ACCESS_KEY_ID',
                                 passwordVariable: 'AWS_SECRET_ACCESS_KEY'
                               ]
                               ]) {
              sh '''
              # Set required configuration
              export PLATFORM=aws
              export CLUSTER="tf-${PLATFORM}-${BRANCH_NAME}-${BUILD_ID}"

              # s3 buckets require lowercase names
              export TF_VAR_tectonic_cluster_name=$(echo ${CLUSTER} | awk '{print tolower($0)}')

              # AWS specific configuration
              export AWS_REGION="us-west-2"

              # make core utils accessible to make
              export PATH=/bin:${PATH}

              # Create local config
              make localconfig

              # Use smoke test configuration for deployment
              ln -sf ${WORKSPACE}/test/aws.tfvars ${WORKSPACE}/build/${CLUSTER}/terraform.tfvars

              make plan

              # always cleanup cluster
              function shutdown() {
                make destroy
              }
              trap shutdown EXIT

              make apply
              '''
            }
          }
        )
      }
    }
  }
}
