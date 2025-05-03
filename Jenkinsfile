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

    stage('Run Allure Tests') {
      steps {
        container('go') {
          catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
            sh '''
              echo "🧪 Running Allure-enhanced Go Tests..."
              mkdir -p ${ALLURE_RESULTS}
              chmod -R 777 ${ALLURE_RESULTS}
              API_URL=${API_URL} INPUT_DIR=${INPUT_DIR} MAP_FILE=${MAP_FILE} go test -v
            '''
          }
        }
      }
    }

    stage('Publish Allure Report') {
      steps {
        script {
          allure([
            includeProperties: false,
            jdk: '',
            results: [[path: "${ALLURE_RESULTS}"]],
            reportBuildPolicy: 'ALWAYS'
          ])
        }
      }
    }
  }

  post {
    always {
      archiveArtifacts artifacts: "${ALLURE_RESULTS}/**", fingerprint: true
    }
  }
}
