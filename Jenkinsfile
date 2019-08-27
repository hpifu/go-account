pipeline {
    agent any
    stages {
        stage('behave') {
            steps {
                sh 'pwd'
                sh 'pwd'
            }
        }
        stage('image') {
            steps {
                sh 'make image'
            }
        }
        stage('deploy') {
            steps {
                sh 'make deploy'
            }
        }
    }
}

