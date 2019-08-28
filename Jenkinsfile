pipeline {
    agent any
    stages {
        stage('unit') {
            steps {
                sh 'pwd'
                sh 'make dockertest'
            }
        }
        stage('image') {
            steps {
                sh 'make image'
            }
        }
        stage('behave') {
            steps {
                sh 'make dockerbehave'
            }
        }
        stage('deploy') {
            steps {
                sh 'make deploy'
            }
        }
    }
}

