# Dev Build and Deployment Guide

## Pre Req:

1. Docker
2. Golang/migrate lib
3. Yarn package manager
4. Make
5. Migrate Cli
6. Psql Db

### Development Instance

- Use `main.sh` to run development instances.
  - ./main.sh -e
  - ./main.sh -f
  - ./main.sh -m
  - ./main.sh -g
- Run `-e` command first to get all .env variables for build.
- Run `-D` to load dev container and run migration scripts. This add extra data for developers to develop against.
- If everything is succesful, you should be able to succesfully have a development instance of psql container within docker with required migration files.
- Continue setup from MakeFile.

### Test Instance

- Use `main.sh` to run test instance.
  - This test instance will have no data. Assumption is that you view mashed app from the scratch.
  - ./main.sh -e
  - ./main.sh -T
- Run `-e` command first to get all .env variables for build.
- Run `-T` for loading test instance. This will load all migration scripts and also run unit tests.
- Run `-u` for cleanup purposes. Must be done manually.

### Production Instance

#### Using docker build targets for multiple instances

To support docker build for various deployment instances, we do not use the `docker-compose.deploy.yml` file so that we can deploy manually and as we please. To do this however, you must build the instances yourself and push to the docker registry. This step should only be done after a new rc has been cut or a new major version is deployed.

1. Login to docker with proper username and password.
2. Build the correct images. Ensure that you are in the correct release candidate and tags are valid.

```bash
# frontend
docker build -t earmuffjam/fleetwise-frontend:1.0.0 -f frontend/prod.Dockerfile .

# backend
docker build -t earmuffjam/fleetwise-backend:1.0.0 -f server/Dockerfile .

# api
docker build -t earmuffjam/fleetwise-apilayer:1.0.0 -f apilayer/Dockerfile .
```

3. Push the docker images to the registry.

```bash
# frontend
docker push earmuffjam/fleetwise-frontend:1.0.0

# backend
docker push earmuffjam/fleetwise-backend:1.0.0

# api
docker push earmuffjam/fleetwise-apilayer:1.0.0

```

4. After these changes have been pushed to the docker register, `ssh` into your production instance. If you have ssh config setup, you can use that or else you can just use the normal means of logging into the instance via `ssh`.

5. Login to docker registry with `docker login`.

6. Pull your latest images.

```bash

# frontend
docker pull earmuffjam/fleetwise-frontend:1.0.0

# backend
docker pull earmuffjam/fleetwise-backend:1.0.0

# api
docker pull earmuffjam/fleetwise-apilayer:1.0.0

```

7. After a successful pull of the images from the registry, you should be ready to run the images into the production environment. These images are already built mini containers that just need to run in your production environment. If you have existing images remove them.

```bash
docker stop fleetwise-frontend fleetwise-backend fleetwise-apilayer
docker rm fleetwise-frontend fleetwise-backend fleetwise-apilayer

```

8. Run the latest images.

```bash

# frontend
docker run -d --name fleetwise-frontend -p 8081:80 earmuffjam/fleetwise-frontend:1.0.0

# backend
docker run -d --name fleetwise-backend -p 8089:5432 --link fleetwise-frontend earmuffjam/fleetwise-backend:1.0.0

# apilayer
docker run -d --name fleetwise-apilayer -p 8087:8087 --link fleetwise-backend earmuffjam/fleetwise-apilayer:1.0.0

# minio storage
docker run -d --name minio -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=${MINIO_ROOT_USER} \
  -e MINIO_ROOT_PASSWORD=${MINIO_ROOT_PASSWORD} \
  earmuffjam/fleetwise-minio:latest

```

9. View the running logs with

```bash

docker logs fleetwise-frontend
docker logs fleetwise-backend
docker logs fleetwise-apilayer

```

#### Using docker compose for a single container instance

1. Production instance must be executed in sequence. Since there should be ability to alter data, we have to run migration scripts as well. Please be aware of this.
2. Only required env variables are copied over.
3. Since data cannot be modified here, no test data is inserted. A fresh container will have no users. Running flag `-u` will `REMOVE ALL DATA`. `NEVER RUN THIS IN PRODUCTION ENV`

- Use `mainDeploy.sh` to run production instance.
  - This is production instance. Please be careful. Data is not scrubbed here.
  - ./mainDeploy.sh -e
  - ./mainDeploy.sh -p
  - ./mainDeploy.sh -m
- Run `-m` to allow for the migration scripts to run. `must`
