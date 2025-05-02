pipeline {
  agent any

  environment {
    SHARD_COUNT = 4
    OUTPUT_DIR = "junit-results"
    API_URL = "http://localhost:8080/transit"
  }

  stages {
    stage('Checkout Code') {
      steps {
        git url: 'https://github.com/babaYaga451/go-tnt-automation.git', branch: 'main'
      }
    }

    stage('Split Input Files') {
      steps {
        script {
          sh "rm -rf ${env.OUTPUT_DIR} && mkdir -p ${env.OUTPUT_DIR}"

          def files = sh(
            script: "find ./data -name '*.txt' | sort -R",
            returnStdout: true
          ).trim().split('\n').toList()

          for (int i = 1; i <= SHARD_COUNT.toInteger(); i++) {
            def shardId = i
            def shardFiles = files.findAll { idx -> (files.indexOf(idx) % SHARD_COUNT.toInteger() + 1) == shardId }
            echo "Shard-${shardId} has ${shardFiles.size()} files"
            writeFile file: "shard-${shardId}.list", text: shardFiles.join('\n')
          }
        }
      }
    }

    stage('Run Tests in Parallel') {
      steps {
        script {
          def branches = [:]

          for (int i = 1; i <= SHARD_COUNT.toInteger(); i++) {
            def shardId = i
            def label = "go-shard-${shardId}"

            branches["shard${shardId}"] = {
              podTemplate(inheritFrom: 'default', label: label, yaml: '''
apiVersion: v1
kind: Pod
spec:
  containers:
    - name: go
      image: golang:1.21
      command:
        - bash
        - -c
        - cat
      tty: true
''') {
                node(label) {
                  container('go') {
                  def cmd = """
                    echo "🔹 Executing shard: shard${shardId}"
                    echo "🔹 Go version:"
                    go version
                    echo "🔹 Go location:"
                    which go
                    echo "🔹 Running test runner"
                    go run cmd/test-transit/main.go \\
                      -inputFiles=\$(cat shard-${shardId}.list | tr '\\n' ',') \\
                      -mapFile=dest.csv \\
                      -apiURL=${API_URL} \\
                      -k=10 \\
                      -workers=4 \\
                      -outputFile=${OUTPUT_DIR}/shard${shardId}.xml
                  """

                  sh cmd
                  }
                }
              }
            }
          }

          parallel branches
        }
      }
    }

    stage('Merge JUnit Reports') {
      steps {
        sh '''
          echo '<?xml version="1.0"?><testsuites>' > ${OUTPUT_DIR}/results.xml
          grep -h '<testsuite' ${OUTPUT_DIR}/shard*.xml >> ${OUTPUT_DIR}/results.xml
          echo '</testsuites>' >> ${OUTPUT_DIR}/results.xml
        '''
      }
    }

    stage('Publish Single Report') {
      steps {
        junit "${OUTPUT_DIR}/results.xml"
      }
    }
  }

  post {
    always {
      archiveArtifacts artifacts: "${OUTPUT_DIR}/*.xml", fingerprint: true
    }
  }
}
