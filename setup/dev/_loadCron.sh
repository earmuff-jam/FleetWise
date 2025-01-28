#!/bin/bash

echo "Setting up cron jobs for psql"
source .env

PGPASSWORD=home psql \
  -h "$DATABASE_DOCKER_CONTAINER_IP_ADDRESS" \
  -p "$DATABASE_DOCKER_CONTAINER_PORT" \
  -U "$POSTGRES_USER" \
  -d $POSTGRES_DB \
  -c "SELECT cron.schedule(
        'refresh-alerts'::text, 
        '0 16 * * *'::text, 
        'SELECT community.populate_maintenance_alerts();'::text
      );"