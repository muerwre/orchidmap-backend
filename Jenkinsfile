pipeline {
    agent any
    
    environment {        
        ENV = "${env.BRANCH_NAME == "master" ? env.ORCHID_STABLE_BACKEND_ENV : env.ORCHID_STAGING_BACKEND_ENV}"
        TZ = "\$(readlink /etc/localtime | sed 's#/usr/share/zoneinfo/##')"
    }

    stages {
        stage('check') {
            steps {
                echo "ENV: ${ENV}"
                echo "WORKSPACE: ${WORKSPACE}"

                script {
                    if("${ENV}" == "" || ("${env.BRANCH_NAME}" != "master" && "${env.BRANCH_NAME}" != "develop")) {
                        error "INCORRECT VARIABLES"
                        return
                    }
                }
            }
        }    

        stage('copy env') {
            steps {
                sh "cp -a ${ENV}/. ${WORKSPACE}"
                sh "echo -en \"\nTZ=${TZ}\" >> ${WORKSPACE}/.env"
                sh "cat ${WORKSPACE}/.env"
            }
        }

        stage('Build (docker)') {
            steps {
                sh "docker-compose build"
            }
        }

        stage('deploy') {
            steps {
                sh "docker-compose up -d"
            }
        }
    }
}
