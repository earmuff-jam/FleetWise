#!/bin/bash

# File: _main.sh
# Description: Loads the application and runs associated scripts
#
# Most common usage would be to run dev instance - run script with -D flag and run script with -T flag for test instance.

help() {
    echo "Usage: $0 [Choose any of the following options...] {loadEnv, loadDb, loadDevEnv, loadUnitTest, loadMigration, loadData, uninstall}" >&2
    echo
    echo "   -e --loadEnv              Loads the environment variables only"
    echo "   -d --loadDb               Loads the psql container with fresh installation"
    echo "   -D --loadDevEnv           Loads the psql container with existing data. Use for development instance"
    echo "   -T --loadUnitTest         Loads the psql container with test data for api tests"
    echo "   -m --loadMigration        Loads the migration in sequence. Does not erase data but requires container to be up"
    echo "   -t --loadData             Loads the test data for application"
    echo "   -u --uninstall            Uninstall the application erasing all data"
    echo
    exit 1
}

loadEnv() {
    echo "envOnly flag provided. refreshing env variables."
    chmod +x setup/dev/_loadEnvVariables.sh
    ./setup/dev/_loadEnvVariables.sh
}

loadDb() {
    echo "loadDb flag provided. building psql docker image and running migrations."
    docker-compose down 
    docker-compose -f docker-compose-db.yml up --build -d
}

loadDevEnv() {
    echo "load dev env flag provided. starting psql container. please wait..."
    loadEnv

    sleep +2
    loadDb

    sleep +2
    loadMigration

    sleep +2
    loadSeedData

    sleep +2
    loadData
    echo "adding generated test data for development purpose only..."
    cd apilayer/dataLake && go run main.go

}

loadSeedData() {
    echo "building seed data files."
    chmod +x setup/dev/_addSeedData.sh
    ./setup/dev/_addSeedData.sh
}

loadUnitTest() {
    echo "load test environment flag provided. using test env variables and building all containers. please wait ..."
    
    chmod +x setup/dev/_loadTestEnvVariables.sh
    ./setup/dev/_loadTestEnvVariables.sh
    
    docker-compose down --remove-orphans --volumes
    docker-compose -f docker-compose-db.yml up --build -d

    sleep +2
    loadMigration

    sleep +2
    loadSeedData

    sleep +2
    loadData

    sleep +2
    echo "running unit test on test db instance $POSTGRES_DB"
    cd apilayer && go test -coverprofile=logs/coverage.out ./... && go tool cover -html=logs/coverage.out

    echo "activity complete. close running containers manually."

}

loadMigration() {
    echo "migration flag provided. running migration in sequence."
    chmod +x setup/dev/_loadMigration.sh
    ./setup/dev/_loadMigration.sh
}

loadData() {
    echo "loadData flag provided. loading all test data."
    chmod +x setup/dev/_addTestData.sh
	./setup/dev/_addTestData.sh
}

uninstall() {
    echo "uninstalling software. removing all data..."
    docker-compose down --remove-orphans --volumes
}

if [ "$#" -eq 0 ]; then
    help
fi

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) help; shift ;;
        -e|--loadEnv) loadEnv; shift ;;
        -d|--loadDb) loadDb; shift ;;
        -D|--loadDevEnv) loadDevEnv; shift ;;
        -T|--loadUnitTest) loadUnitTest; shift ;;
        -m|--loadMigration) loadMigration; shift ;;
        -t|--loadData) loadData; shift ;;
        -u|--uninstall) uninstall; shift ;;
        *) help; echo "Unknown parameter passed: $1" ;;
    esac
    shift
done