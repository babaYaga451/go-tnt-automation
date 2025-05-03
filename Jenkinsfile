pipeline {
  agent {
    kubernetes {
      yaml """
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: go
      image: golang:1.23.5
      command:
        - cat
      tty: true
"""
    }
  }

  environment {
    API_URL = "http://transit-api.jenkins-cluster.svc.cluster.local/transit"
    INPUT_DIR = "./data"
    MAP_FILE = "dest.csv"
    ALLURE_RESULTS = "allure-results"
  }

  stages {
    stage('Checkout Code') {
      steps {
        git url: 'https://github.com/babaYaga451/go-tnt-automation.git', branch: 'main'
      }
    }

    stage('Prepare Directories') {
      steps {
        container('go') {
          sh '''
            mkdir -p ${ALLURE_RESULTS}
          '''
        }
      }
    }

    stage('Run Allure Tests') {
      steps {
        container('go') {
          sh '''
            echo "Running Allure-enhanced Go Tests..."
            go test -v
          '''
        }
      }
    }

    stage('Publish Allure Report') {
      steps {
        allure([
          includeProperties: false,
          jdk: '',
          results: [[path: "${ALLURE_RESULTS}"]],
          reportBuildPolicy: 'ALWAYS'
        ])
      }
    }
  }
  
  post {
    always {
      archiveArtifacts artifacts: "${ALLURE_RESULTS}/**", fingerprint: true
    }
  }
}
