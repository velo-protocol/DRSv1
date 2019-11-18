pipeline {
    agent { label "slave" }
    environment{
        sonarqubeURL="http://sonarqube.lightnet-nonprod.com:9000"
        branchName = sh(
                script: "printf \$(git rev-parse --abbrev-ref HEAD | sed -e \"s|/|-|g\")",
                returnStdout: true
        )
        GIT_COMMIT_SHORT = sh(
                script: "printf \$(git rev-parse --short=8 ${GIT_COMMIT})",
                returnStdout: true
        )
        newVersion="${env.BUILD_NUMBER}-${env.GIT_COMMIT_SHORT}"
        dockerTag="${env.branchName}-${env.newVersion}"
        dockerImage="${env.CONTAINER_IMAGE}:${env.dockerTag}"
        appName="cen"
        CONTAINER_IMAGE="registry.gitlab.com/velo-labs/${appName}"
    }
    stages {
        stage('Build Image Test') {
            steps {
                withCredentials([usernamePassword(credentialsId: '6a0ef684-a441-4262-937f-d9a7a0602b56', passwordVariable: 'gitlabPassword', usernameVariable: 'gitlabUsername')]) {
                    sh '''
                        echo "Build Image"
                        docker login -u ${gitlabUsername} -p ${gitlabPassword} registry.gitlab.com
                        docker build --pull --target development -t ${CONTAINER_IMAGE}:${dockerTag}-test -f node/Dockerfile .
                    '''
                }
            }
        }

        stage('Unit Test') {
            steps {
                sh '''
                    echo "Run unit test -> ${CONTAINER_IMAGE}:${dockerTag}-test"
                    mkdir -p $(pwd)/reports
                    docker run --rm -v $(pwd)/reports:/reports ${CONTAINER_IMAGE}:${dockerTag}-test sh -c "make ci_test | go-junit-report > /reports/coverage-tasks.xml"
                    docker run --rm -v $(pwd)/reports:/reports ${CONTAINER_IMAGE}:${dockerTag}-test sh -c "make ci_test; gocov convert .coverage.txt | gocov-xml > /reports/coverage.xml && cp .coverage.txt /reports/"
                    sudo chown -R jenkins:jenkins $(pwd)/reports
                '''
            }
        }

        stage('SonarQube Code Analysis') {
            steps {
                script {
                    echo "SonarQube Code Analysis"
                    withSonarQubeEnv('sonarqube') {
                        sh '''
                            sed -i s~#SONARQUBE_URL#~${sonarqubeURL}~g Makefile
                            sed -i s~#APP_VERSION#~${dockerTag}~g Makefile
                            make ci_sonarqube
                        '''
                    }
                }
            }
        }

        stage('SonarQube Quality Gate') {
            steps {
                sleep(10)
                script {
                    echo "SonarQube Quality Gate"
                    timeout(time: 1, unit: 'MINUTES') {
                        def qg = waitForQualityGate()
                        if (qg.status != 'OK') {
                            error "Pipeline aborted due to quality gate failure: ${qg.status}"
                        }
                    }
                }
            }
        }

        stage('Build and Push to Registry') {
            when {
                anyOf {
                    branch 'develop';
                    branch 'release/*';
                    branch 'master'
                }
            }
            steps {
                withCredentials([usernamePassword(credentialsId: '6a0ef684-a441-4262-937f-d9a7a0602b56', passwordVariable: 'gitlabPassword', usernameVariable: 'gitlabUsername')]) {
                    sh '''
                        echo "Build and Push to Registry"
                        docker login -u ${gitlabUsername} -p ${gitlabPassword} registry.gitlab.com
                        docker build --pull -t ${CONTAINER_IMAGE}:${dockerTag} -f node/Dockerfile .
                        docker push ${CONTAINER_IMAGE}:${dockerTag}
                        docker tag ${CONTAINER_IMAGE}:${dockerTag} ${CONTAINER_IMAGE}:${branchName}
                        docker push ${CONTAINER_IMAGE}:${branchName}
                    '''
                }
            }
        }

        stage('Trigger to Deployment job') {
            parallel {
                stage ('Deploy to Develop Environment') {
                    when {
                        branch 'develop'
                    }
                    steps {
                        build job: 'DRSv1-deploy', parameters: [string(name: 'dockerVersion', value: env.dockerTag),string(name: 'environment', value: 'develop')]
                    }
                }
                stage ('Deploy to Staging Environment') {
                    when {
                        branch 'master'
                    }
                    steps {
                        build job: 'DRSv1-deploy', parameters: [string(name: 'dockerVersion', value: env.dockerTag),string(name: 'environment', value: 'staging')]
                    }
                }
            }
        }
    }

    post {
        always {
            junit "reports/coverage-tasks.xml"
            cobertura coberturaReportFile: "reports/coverage.xml"
        }
    }
}
