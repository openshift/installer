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
  }
}
