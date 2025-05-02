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

          for (int i = 0; i < SHARD_COUNT.toInteger(); i++) {
            sh "mkdir -p ${SHARD_DIR}/shard${i+1}"
          }

          for (int i = 0; i < files.size(); i++) {
            int shardNum = (i % SHARD_COUNT) + 1
            sh "cp ${files[i]} ${SHARD_DIR}/shard${shardNum}/"
          }
        }
      }
    }

    stage('Run Shards in Parallel') {
      parallel {
        stage('Shard 1') {
          steps {
            sh 'echo "🔹 Files in shard1:" && ls -l shards/shard1'
          }
        }
        stage('Shard 2') {
          steps {
            sh 'echo "🔹 Files in shard2:" && ls -l shards/shard2'
          }
        }
        stage('Shard 3') {
          steps {
            sh 'echo "🔹 Files in shard3:" && ls -l shards/shard3'
          }
        }
        stage('Shard 4') {
          steps {
            sh 'echo "🔹 Files in shard4:" && ls -l shards/shard4'
          }
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
