pipeline {
  agent any

  environment {
    SHARD_COUNT = 4
    SHARD_DIR = "shards"
  }

  stages {
    stage('Checkout SCM') {
      steps {
        git url: 'https://github.com/babaYaga451/go-tnt-automation.git', branch: 'main'
      }
    }

    stage('Split Files Into Shards') {
      steps {
        script {
          sh "rm -rf ${SHARD_DIR} && mkdir -p ${SHARD_DIR}"
          def files = sh(script: "find ./data -name '*.txt' | sort", returnStdout: true).trim().split('\n')

          int shardCount = SHARD_COUNT.toInteger()
          int filesPerShard = Math.ceil(files.size() / (double)shardCount) as int

          for (int i = 0; i < shardCount; i++) {
            sh "mkdir -p ${SHARD_DIR}/shard${i+1}"
            def start = i * filesPerShard
            def end = Math.min(start + filesPerShard, files.size())
            def shardFiles = files.subList(start, end)

            for (f in shardFiles) {
              sh "cp '${f}' ${SHARD_DIR}/shard${i+1}/"
            }
          }
        }
      }
    }

    stage('Run Shards in Parallel') {
          steps {
            script {
              def branches = [:]
              for (int i = 1; i <= SHARD_COUNT.toInteger(); i++) {
                def shardNum = i // must capture `i` in a final variable for Groovy closure
                branches["Shard ${shardNum}"] = {
                  sh "echo '🔹 Files in shard${shardNum}:' && ls -l shards/shard${shardNum}"
                }
              }
              parallel branches
            }
          }
        }
      }

  post {
    always {
      archiveArtifacts artifacts: 'shards/**/*.txt', fingerprint: true
    }
  }
}
