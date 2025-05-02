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
          ).trim().split('\n')

          for (int i = 1; i <= SHARD_COUNT.toInteger(); i++) {
            def shardId = i
            def shardFiles = files.findAll { idx -> (files.indexOf(idx) % SHARD_COUNT.toInteger() + 1) == shardId }
            writeFile file: "shard-${shardId}.list", text: shardFiles.join('\n')
          }
        }
      }
    }

    stage('Run Tests in Parallel') {
      parallel {
        script {
          def branches = [:]
          for (int i = 1; i <= SHARD_COUNT.toInteger(); i++) {
            def shardId = i
            branches["shard${shardId}"] = {
              stage("Run shard${shardId}") {
                podTemplate(yaml: """
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: go
    image: golang:1.21
    command:
    - cat
    tty: true
""") {
                  node(POD_LABEL) {
                    container('go') {
                      sh """
                        go run cmd/test-transit/main.go \\
                          -inputFiles=$(cat shard-${shardId}.list | tr '\\n' ',') \\
                          -mapFile=dest.csv \\
                          -apiURL=${API_URL} \\
                          -k=10 \\
                          -workers=4 \\
                          -outputFile=${OUTPUT_DIR}/shard${shardId}.xml
                      """
                    }
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
      agent any
      steps {
        sh """
          echo '<?xml version="1.0"?><testsuites>' > ${OUTPUT_DIR}/results.xml
          grep -h '<testsuite' ${OUTPUT_DIR}/shard*.xml >> ${OUTPUT_DIR}/results.xml
          echo '</testsuites>' >> ${OUTPUT_DIR}/results.xml
        """
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
