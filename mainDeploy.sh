
#!/bin/bash

# File: _main.sh
# Description: Loads the database layer for the application

help() {
    echo "Usage: $0 [Choose any of the following options...] {loadEnv, loadDb, loadMigration, uninstall}" >&2
    echo
    echo "   -e --loadEnv              Loads the environment variables only"
    echo "   -p --loadProdEnv          Allows to load all containers and start with new fresh data. This is production env. Does not erase data"
    echo "   -m --loadMigration        Allows to load the migration in sequence. Does not erase data but requires container to be up"
    echo "   -u --uninstall            Uninstall the application erasing all data"
    echo
    exit 1
}

loadEnv() {
    echo "envOnly flag provided. only loading env variables."
    chmod +x setup/prod/_loadProdEnvVariables.sh
    ./setup/prod/_loadProdEnvVariables.sh
}

loadCronJobs() {
    echo "loadCronJobs flag provided. setting up cron jobs for existing psql functions"
    chmod +x setup/prod/_loadCron.sh
    ./setup/prod/_loadCron.sh
}

loadProdEnv() {

    echo "loadProdEnv flag provided. building all containers in production mode."

    chmod +x setup/prod/_loadProdEnvVariables.sh
    ./setup/prod/_loadProdEnvVariables.sh

    docker-compose down
    docker-compose -f docker-compose.deploy.yml up --build -d

    sleep +2
    loadMigration

    sleep +2
    loadCronJobs

    sleep +2
    loadSeedData

}

loadSeedData() {
    echo "building seed data files."
    chmod +x setup/prod/_addSeedData.sh
    ./setup/prod/_addSeedData.sh
}

loadMigration() {

    echo "migration flag provided. running migration in sequence."

    chmod +x setup/prod/_loadMigration.sh
    ./setup/prod/_loadMigration.sh
}

uninstall() {

    echo " WARNING !!!! PRODUCTION INSTANCE:: PURGING ALL DATA ... Have you backed up? Press 'Y' to confirm, or 'N' to cancel."
    read -r -p "Are you sure you want to proceed? [Y/n] " response
    response=$(echo "$response" | tr '[:upper:]' '[:lower:]') 
    if [[ "$response" =~ ^(yes|y| ) ]]; then
        echo "PRODUCTION INSTANCE:: PURGING ALL DATA ... hope you know what you are doing..."
        sleep 10
        docker-compose down --remove-orphans --volumes
    else
        echo "Uninstallation canceled."
    fi
}

if [ "$#" -eq 0 ]; then
    help
fi

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -h|--help) help; shift ;;
        -e|--loadEnv) loadEnv; shift ;;
        -p|--loadProdEnv) loadProdEnv; shift ;;
        -m|--loadMigration) loadMigration; shift ;;
        -u|--uninstall) uninstall; shift ;;
        *) help; echo "Unknown parameter passed: $1" ;;
    esac
    shift
done
