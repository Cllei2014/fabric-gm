def getRepoURL() {
  sh "git config --get remote.origin.url > .git/remote-url"
  return readFile(".git/remote-url").trim()
}

void setBuildStatus(String message, String state) {
  repoUrl = getRepoURL()

  step([
      $class: "GitHubCommitStatusSetter",
      reposSource: [$class: "ManuallyEnteredRepositorySource", url: repoUrl],
      contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
      errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
      statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
  ]);
}

pipeline {
    agent any

    environment {
        DOCKER_NS     = "${DOCKER_REGISTRY}/twbc"
        EXTRA_VERSION = "build-${BUILD_NUMBER}"
        GOPATH        = "${WORKSPACE}/gopath"
    }

    stages {
        stage('Build Image') {
            steps {
                setBuildStatus("Build Started", "PENDING");

                dir('gopath/src/github.com/hyperledger/fabric') {
                    checkout scm

                    sh '''
                    go mod vendor
                    make docker
                    '''
                }
            }
        }

        stage('Upload Image') {
            steps {
                dir('gopath/src/github.com/hyperledger/fabric') {
                    sh 'aws ecr get-login-password | docker login --username AWS --password-stdin ${DOCKER_REGISTRY}'
                    sh '''
                    make docker-list 2>/dev/null | grep "$DOCKER_NS" | while read line
                    do
                       docker tag $line ${line/:*/:latest}
                       docker push $line
                       docker push ${line/:*/:latest}
                       docker rmi $line
                    done
                    '''
                }
            }
        }
        stage('Test Fabcar') {
            steps {
                dir('gopath/src/github.com/hyperledger/fabric') {
                    catchError(buildResult: 'UNSTABLE', stageResult: 'FAILURE') {
                        script {
                            def result = build(
                                job: 'fabric-sample-gm',
                                propagate: false,
                                parameters: [
                                    [$class: 'StringParameterValue', name: 'IMAGE_PEER', value: sh(script: 'make peer-docker-list 2>/dev/null ', returnStdout: true).trim()],
                                    [$class: 'StringParameterValue', name: 'IMAGE_ORDERER', value: sh(script: 'make orderer-docker-list 2>/dev/null ', returnStdout: true).trim()],
                                    [$class: 'StringParameterValue', name: 'IMAGE_TOOLS', value: sh(script: 'make tools-docker-list 2>/dev/null ', returnStdout: true).trim()],
                                ]
                            )
                            if (result.result.equals("SUCCESS")) {
                                echo "Passed Test Fabcar"
                            } else {
                                error "Failed Test Fabcar"
                            }
                        }
                    }
                }
            }
        }
    }

    post {
        success {
            setBuildStatus("Build succeeded", "SUCCESS");
        }
        unsuccessful {
            setBuildStatus("Build failed", "FAILURE");
        }
    }
}

