pipeline {
    agent any

    environment {
        DOCKER_NS     = "twbc"
        EXTRA_VERSION = "build-${BUILD_NUMBER}"
    }

    stages {
        stage('Build Image') {
            steps {
                sh '''
                make docker
                '''
            }
        }
        stage('Upload Image') {
            steps {
                sh 'aws ecr get-login-password | docker login --username AWS --password-stdin ${DOCKER_REGISTRY}'
                sh '''
                make docker-list 2>/dev/null | grep ^twbc | while read line
                do
                   docker tag $line $DOCKER_REGISTRY/$line
                   docker tag $line $DOCKER_REGISTRY/${line/:*/:latest}
                   docker push $DOCKER_REGISTRY/$line
                   docker push $DOCKER_REGISTRY/${line/:*/:latest}
                   docker rmi $line $DOCKER_REGISTRY/$line
                done
                '''
            }
        }
        stage('Test Fabcar') {
            steps {
                script {
                    try {
                        build(
                            job: 'fabric-sample-gm',
                            parameters: [
                                [$class: 'StringParameterValue', name: 'IMAGE_PEER', value: sh(script: 'make peer-docker-list 2>/dev/null ', returnStdout: true).trim()],
                                [$class: 'StringParameterValue', name: 'IMAGE_ORDERER', value: sh(script: 'make orderer-docker-list 2>/dev/null ', returnStdout: true).trim()],
                                [$class: 'StringParameterValue', name: 'IMAGE_TOOLS', value: sh(script: 'make tools-docker-list 2>/dev/null ', returnStdout: true).trim()],
                            ]
                        )
                    } catch (Exception e) {
                        currentBuild.result = 'UNSTABLE'
                    }
                }
            }
        }
    }
}