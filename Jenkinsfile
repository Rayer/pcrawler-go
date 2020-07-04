def PROJECT_NAME = "Rayer/pcrawler-go"
pipeline {
   agent any

   stages {
    stage('Fetch from github') {
        steps {
            slackSend message: "Project ${PROJECT_NAME} start to build."
            git credentialsId: '26c5c0a0-d02d-4d77-af28-761ffb97c5cc', url: 'https://github.com/Rayer/pcrawler-go.git'
        }
    }
    stage('Unit test') {
        steps {
            sh label: 'go version', script: 'go version'
            sh label: 'install gocover-cobertura', script: 'go get github.com/t-yuki/gocover-cobertura'
            sh label: 'go unit test', script: 'go test --coverprofile=cover.out'
            sh label: 'convert coverage xml', script: '~/go/bin/gocover-cobertura < cover.out > coverage.xml'
        }
    }
    stage ("Extract test results") {
        steps {
            cobertura coberturaReportFile: 'coverage.xml'
        }
    }

    stage('Benchmark') {
        steps {
            sh label: 'go benchmark', script: 'go test -bench=.'
        }
    }

    stage('build executable') {
        steps {
            sh label: 'show version', script: 'go version'
            sh label: 'build library', script: 'go build'
        }
    }
   }

   post {
        aborted {
            slackSend message: "Project ${PROJECT_NAME} aborted."
        }
        success {
            slackSend message: "Project ${PROJECT_NAME} is built successfully."
        }
        failure {
            slackSend message: "Project ${PROJECT_NAME} is failed to be built."
        }
    }
}
