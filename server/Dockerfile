FROM postgres:15.4
RUN apt-get update && apt-get install -y curl


RUN apt-get -y install postgresql-15-cron
RUN echo "shared_preload_libraries='pg_cron'" >> /usr/share/postgresql/postgresql.conf.sample
RUN echo "cron.database_name='community'" >> /usr/share/postgresql/postgresql.conf.sample

# Install extension in postgres database
RUN echo "CREATE EXTENSION IF NOT EXISTS pg_cron;" >> /docker-entrypoint-initdb.d/init.sql

# Expose port 5432
EXPOSE 5432

