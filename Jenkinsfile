pipeline {
    agent any
    stages {
        stage('behave') {
            steps {
                sh 'pwd'
                sh 'make behave'
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

