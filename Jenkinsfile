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
        stage ('Cleanup') {
            steps {
                dir('directoryToDelete') {
                    deleteDir()
                }
            }
        }

        stage ('Pull App Config') {
            steps {
                withCredentials([usernamePassword(credentialsId: '6a0ef684-a441-4262-937f-d9a7a0602b56', passwordVariable: 'gitlabPassword', usernameVariable: 'gitlabUsername')]) {
                    sh '''
                    echo "Pull App Config" 
                    git clone -b master https://${gitlabUsername}:${gitlabPassword}@gitlab.com/velo-labs/velo-app-configs.git
                '''
                }
            }
        }

        stage ('Pull App Deployment') {
            steps {
                withCredentials([usernamePassword(credentialsId: '6a0ef684-a441-4262-937f-d9a7a0602b56', passwordVariable: 'gitlabPassword', usernameVariable: 'gitlabUsername')]) {
                    sh '''
                    echo "Pull App Deployment"
                    git clone -b master https://${gitlabUsername}:${gitlabPassword}@gitlab.com/velo-labs/velo-app-deployment.git
                '''
                }
            }
        }

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
                    docker run --rm -v $(pwd)/reports:/reports ${CONTAINER_IMAGE}:${dockerTag}-test sh -c "go test ./node/app/... -v -coverprofile .coverage.txt | go-junit-report > /reports/coverage-tasks.xml"
                    docker run --rm -v $(pwd)/reports:/reports ${CONTAINER_IMAGE}:${dockerTag}-test sh -c "go test ./node/app/... -v -coverprofile .coverage.txt; gocov convert .coverage.txt | gocov-xml > /reports/coverage.xml && cp .coverage.txt /reports/"
                    sudo chown -R jenkins:jenkins $(pwd)/reports
                '''
            }
        }

        stage('SonarQube Code Analysis') {
            steps {
                script {
                    echo "SonarQube Code Analysis"
                    withSonarQubeEnv('sonarqube') {
                        sh  "sonar-scanner " +
                                "-Dsonar.host.url=${sonarqubeUrl} " +
                                "-Dsonar.projectKey=${appName} " +
                                "-Dsonar.projectName=${appName} " +
                                "-Dsonar.projectVersion=${dockerTag} " +
                                "-Dsonar.sonar.tests=node/** " +
                                "-Dsonar.exclusions=**/*_test.go,**/vendor/** " +
                                "-Dsonar.go.coverage.reportPaths=reports/.coverage.txt "
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

        stage('Deployment') {
            parallel {
                stage ('Deploy to Develop Environment') {
                    when {
                        branch 'develop'
                    }
                    steps {
                        withCredentials([sshUserPrivateKey(credentialsId: '7a468a86-b55c-4bc2-a0de-e3f9e77568be', usernameVariable: 'ec2user', keyFileVariable: 'ec2keyfile')]) {
                        sh '''
                            echo "Deployment"
                            ec2_json=$(aws ec2 describe-instances --filters "Name=tag-key,Values=Service" "Name=tag-value,Values=velo-cen-node" "Name=instance-state-name,Values=running" | jq -r '.Reservations[]')
                            ec2_tag_output(){
                              ec2_output=$(aws ec2 describe-instances --instance-id $1 --query 'Reservations[*].Instances[*].[PrivateDnsName,Tags[?Key==`ConfigMap`]|[0].Value]' --output text )
                            }
                            instances=$(echo $ec2_json | jq -r '.Instances[].InstanceId');
                            for row in ${instances}; do
                              ec2_tag_output $row
                              addr=$(echo $ec2_output | awk '{print $1}')
                              configmap=$(echo $ec2_output | awk '{print $2}')
                              echo "address:"$addr", configmap:"$configmap
                              anslog=$(ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i ${addr}, -u ${ec2user} --private-key ${ec2keyfile} -e configdir=${WORKSPACE}/velo-app-configs/develop -e host=${addr} -e configmap=${configmap} -e docker_image=${dockerImage} ${WORKSPACE}/velo-app-deployment/deploy-app.yaml)
                              echo $anslog
                            done;
                        '''
                        }
                    }
                }
                stage ('Deploy to Test Environment') {
                    when {
                        branch 'release/*'
                    }
                    steps {
                        withCredentials([sshUserPrivateKey(credentialsId: '7a468a86-b55c-4bc2-a0de-e3f9e77568be', usernameVariable: 'ec2user', keyFileVariable: 'ec2keyfile')]) {
                            sh '''
                            echo "Deployment"
                            ec2_json=$(aws ec2 describe-instances --filters "Name=tag-key,Values=Service" "Name=tag-value,Values=velo-cen-node" "Name=instance-state-name,Values=running" | jq -r '.Reservations[]')
                            ec2_tag_output(){
                              ec2_output=$(aws ec2 describe-instances --instance-id $1 --query 'Reservations[*].Instances[*].[PrivateDnsName,Tags[?Key==`ConfigMap`]|[0].Value]' --output text )
                            }
                            instances=$(echo $ec2_json | jq -r '.Instances[].InstanceId');
                            for row in ${instances}; do
                              ec2_tag_output $row
                              addr=$(echo $ec2_output | awk '{print $1}')
                              configmap=$(echo $ec2_output | awk '{print $2}')
                              echo "address:"$addr", configmap:"$configmap
                              anslog=$(ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i ${addr}, -u ${ec2user} --private-key ${ec2keyfile} -e configdir=${WORKSPACE}/velo-app-configs/test -e host=${addr} -e configmap=${configmap} -e docker_image=${dockerImage} ${WORKSPACE}/velo-app-deployment/deploy-app.yaml)
                              echo $anslog
                            done;
                        '''
                        }
                    }
                }
                stage ('Deploy to Staging Environment') {
                    when {
                        branch 'master'
                    }
                    steps {
                        withCredentials([sshUserPrivateKey(credentialsId: '7a468a86-b55c-4bc2-a0de-e3f9e77568be', usernameVariable: 'ec2user', keyFileVariable: 'ec2keyfile')]) {
                            sh '''
                            echo "Deployment"
                            ec2_json=$(aws ec2 describe-instances --filters "Name=tag-key,Values=Service" "Name=tag-value,Values=velo-cen-node" "Name=instance-state-name,Values=running" | jq -r '.Reservations[]')
                            ec2_tag_output(){
                              ec2_output=$(aws ec2 describe-instances --instance-id $1 --query 'Reservations[*].Instances[*].[PrivateDnsName,Tags[?Key==`ConfigMap`]|[0].Value]' --output text )
                            }
                            instances=$(echo $ec2_json | jq -r '.Instances[].InstanceId');
                            for row in ${instances}; do
                              ec2_tag_output $row
                              addr=$(echo $ec2_output | awk '{print $1}')
                              configmap=$(echo $ec2_output | awk '{print $2}')
                              echo "address:"$addr", configmap:"$configmap
                              anslog=$(ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i ${addr}, -u ${ec2user} --private-key ${ec2keyfile} -e configdir=${WORKSPACE}/velo-app-configs/staging -e host=${addr} -e configmap=${configmap} -e docker_image=${dockerImage} ${WORKSPACE}/velo-app-deployment/deploy-app.yaml)
                              echo $anslog
                            done;
                        '''
                        }
                    }
                }
            }
        }

        stage('Check Deployment Status') {
            steps {
                sh '''
                    echo "Check Deployment Status"
                    ec2_json=$(aws ec2 describe-instances --filters "Name=tag-key,Values=Service" "Name=tag-value,Values=velo-cen-node" "Name=instance-state-name,Values=running" | jq -r '.Reservations[]')
                    ec2_tag_output(){
                      ec2_output=$(aws ec2 describe-instances --instance-id $1 --query 'Reservations[*].Instances[*].[PrivateDnsName,Tags[?Key==`ConfigMap`]|[0].Value]' --output text )
                    }
                    instances=$(echo $ec2_json | jq -r '.Instances[].InstanceId');
                    for row in ${instances}; do
                      ec2_tag_output $row
                      addr=$(echo $ec2_output | awk '{print $1}')
                      grpc-health-probe -addr=${addr}:6666
                    done;
                '''
            }
        }
    }
    post {
            always {
                junit "reports/coverage-tasks.xml"
                cobertura coberturaReportFile: "reports/coverage.xml"
            }
            cleanup {
                deleteDir()
            }
    }
}