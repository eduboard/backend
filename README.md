# backend [![Build Status](https://travis-ci.org/eduboard/backend.svg?branch=master)](https://travis-ci.org/eduboard/backend) [![codecov](https://codecov.io/gh/eduboard/backend/branch/master/graph/badge.svg)](https://codecov.io/gh/eduboard/backend)

The backend for eduBoard, written in Go with [httprouter](https://github.com/julienschmidt/httprouter) and MongoDB.

## Installing
- Install `dep`
- Run `dep ensure` to install dependencies

## Running
The backend needs MongoDB to be connected. Connection parameters can be changed using ENV.

### Docker
The easiest way to run the backend is using Docker. Just run `docker-compose up` and you are done.

### Traefik (optional)
The traefik reverse proxy is used for https support. Some changes need to be made to traefik.toml 
- change email to a valid email
- change domain to your domain
- move the file to /opt/traefik/traefik.toml

### Watchtower (optional)
The watchtower makes sure that you that you always have the newest eduboard version running.
You can start it over the command line (your container name might be different from `eduboard_server_1`)

- `docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  v2tec/watchtower eduboard_server_1`

## Endpoints
See the [Spec](./ENDPOINTS.md).