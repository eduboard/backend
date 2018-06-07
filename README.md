# backend [![Build Status](https://travis-ci.org/eduboard/backend.svg?branch=master)](https://travis-ci.org/eduboard/backend)
The backend for eduBoard, written in Go with [httprouter](https://github.com/julienschmidt/httprouter) and MongoDB.

## Installing
- Run `dep ensure` to install dependencies

### Traefik (optional)
The traefik reverse proxy is for https support. To use it you have to
- change email to a valid email in traefik.toml
- change domain to your domain in traefik.toml
- move the file to /opt/traefik/traefik.toml

## Endpoints

- `/` reservered for frontend static files

### Registration
- `/api/register` Register a new user.
- `/api/login` Login an existing user.
- `/api/logout` Logout current user.

### User
- `/api/v1/user/:id/courses/` GET users courses.

### Courses
- `/api/v1/courses/` GET all accessible courses
- `/api/v1/courses/:id/` GET a certain course
