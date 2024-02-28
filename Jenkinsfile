pipeline {
    // Specifies that this pipeline can run on any available agent
    agent any

    // Defines the sequence of stages that will be executed
    stages {
        // This stage checks out the source code from the SCM (Source Code Management) system
        stage('Checkout') {
            steps {
                // This command checks out the source code from the SCM into the Jenkins workspace
                checkout scm
            }
        }
}