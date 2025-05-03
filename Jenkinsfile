pipeline {
  agent {
    kubernetes {
      yaml """
apiVersion: v1
kind: Pod
spec:
  volumes:
    - name: workspace-volume
      emptyDir: {}
  containers:
    - name: go
      image: golang:1.23.5
      command: ["cat"]
      tty: true
      volumeMounts:
        - name: workspace-volume
          mountPath: /home/jenkins/agent
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
          sh 'mkdir -p ${ALLURE_RESULTS}'
        }
      }
    }

    stage('Run Allure Tests') {
      steps {
        container('go') {
          catchError(buildResult: 'FAILURE', stageResult: 'FAILURE') {
            sh '''
              echo "Running Allure-enhanced Go Tests..."
              API_URL=${API_URL} INPUT_DIR=${INPUT_DIR} MAP_FILE=${MAP_FILE} go test -v
            '''
          }
        }
      }
    }

    stage('Archive Allure Results') {
          steps {
            archiveArtifacts artifacts: "${ALLURE_RESULTS}/**", fingerprint: true
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
}
