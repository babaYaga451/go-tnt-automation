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
    OUTPUT_FILE = "shard1/results.xml"
  }

  stages {
    stage('Checkout Code') {
      steps {
        git url: 'https://github.com/babaYaga451/go-tnt-automation.git', branch: 'main'
      }
    }

    stage('Prepare Output Dir') {
      steps {
        sh 'mkdir -p shard1'
      }
    }

    stage('Run Transit Test') {
      steps {
        container('go') {
          sh '''
            go version
            echo "🔹 Running Go Transit Test..."
            go run cmd/test-transit/main.go \
              -inputDir=./data \
              -mapFile=./dest.csv \
              -apiURL=$API_URL \
              -k=10 \
              -workers=4 \
              -outputFile=$OUTPUT_FILE
          '''
        }
      }
    }

    stage('Publish Report') {
      steps {
        sh 'cat shard1/results.xml'
        junit "${OUTPUT_FILE}"
      }
    }
  }

  post {
    always {
      archiveArtifacts artifacts: 'shard1/results.xml', fingerprint: true
    }
  }
}
