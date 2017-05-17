// Update https://jenkins-tectonic.prod.coreos.systems/job/tectonic-release
pipeline {
  agent none
  
  options {
    timeout(time:25, unit:'MINUTES')
    buildDiscarder(logRotator(numToKeepStr:'20'))
  }

  parameters {
    string(name: 'releaseTag')
    boolean(name: 'pre-release')
  }
  
  stages { 
    stage('Release') {
      agent none
      environment {
        GO_PROJECT = '/go/src/github.com/coreos/tectonic-installer'
      }
      steps {
        script {
          podTemplate(
            label: 'webapp-pod',
            containers: [
              containerTemplate(
                name: 'webapp-agent',
                image: 'quay.io/coreos/tectonic-builder:v1.7',
                ttyEnabled: true,
                command: 'cat',
              )
            ],
            volumes: []
          ) {
            node('webapp-pod') {
              withCredentials([[
                $class: 'UsernamePasswordMultiBinding',
                credentialsId: 'tectonic-release-s3-upload',
                usernameVariable: 'AWS_ACCESS_KEY_ID',
                passwordVariable: 'AWS_SECRET_ACCESS_KEY'
              ]]) {
                withCredentials([[
                  $class: 'UsernamePasswordBinding',
                  credentialsId: 'github-coreosbot',
                  variable: 'GITHUB_CREDENTIALS'
                ]]) {
                  container('webapp-agent') {
                    checkout([$class: 'GitSCM', branches: [[name: "refs/tags/${params.releaseTag}"]], userRemoteConfigs: [[credentialsId: 'github-coreosbot', url: 'https://github.com/coreos/tectonic-installer.git']]])
                    sh "mkdir -p \$(dirname $GO_PROJECT) && ln -sf $WORKSPACE $GO_PROJECT"
                    sh "go get github.com/golang/lint/golint"
                    sh """#!/bin/bash -ex
                    export AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID
                    export AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY
                    export GITHUB_CREDENTIALS=$GITHUB_CREDENTIALS
                    export PRE_RELEASE=${params.pre-release}
                    go version
                    cd $GO_PROJECT/installer
                    make build
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
    }
  }
}
