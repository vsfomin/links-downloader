pipeline {
    environment {
        registry = "fomik2/link-downloader"
        registryCredential = credentials('docker-registry')
    }
    agent any

    stages {
        stage ('Pre-Build') {
            steps {
                cleanWs()
            }
           
        }
        stage ('Clone') {
            steps {
                git branch: 'main', credentialsId: '445e8779-a0e8-4b6c-a236-1cccd28f8ca0', url: 'git@github.com:vsfomin/links-downloader.git'
                sh 'git clone git@github.com:vsfomin/publisher.git'
                //delete all other images
                sh 'docker-compose down'
                sh "docker images | egrep \"publisher|downloader\" | awk \'{print \$3}\' | xargs -I ARG docker image rm ARG -f "   
            }
        }
        stage('Test') {
            steps {
                echo 'This is a unit tests steps'
                sh 'sonar-scanner'
                sh 'go test ./...'
            }
        }
            
            
        stage('Build image and login to dockerhub') {
            steps {
                echo 'Build link-downloader image and push it'
                sh 'docker build -t fomik2/${JOB_BASE_NAME}:version-${BUILD_NUMBER} .'
                sh 'echo $registryCredential_PSW | docker login -u $registryCredential_USR --password-stdin'
                echo 'Build publisher image'
                sh 'docker build -t fomik2/publisher:version-${BUILD_NUMBER} ./publisher'
               
            }
        }
            
        stage('Push image to dockerhub') {
            steps {
                echo 'Push image from registry'
                sh 'docker push fomik2/${JOB_BASE_NAME}:version-${BUILD_NUMBER}'
                sh 'docker push fomik2/publisher:version-${BUILD_NUMBER}'
            }
        }
        
        stage('Docker-compose down and up') {
            environment {
                DC_IMAGE_NAME = "fomik2/${JOB_BASE_NAME}"
                DC_TAG = "version-${BUILD_NUMBER}" 
            }
            steps {
                echo 'Docker-compose down and up'
                sh 'cat docker-compose.yaml'
                sh 'docker-compose down && docker-compose up -d'
            }
        }
        
    }
        post {
            always {
                sh "docker images | grep none | awk \'{print \$3}\' | xargs -I \'{}\' docker image rm \'{}\' -f "
                //sh "docker images | grep version-${currentBuild.previousBuild.number} | awk \'{print \$3}\' | xargs -I \'{}\' docker image rm \'{}\' -f "
                //sh 'docker logout'
            }
        }
}

