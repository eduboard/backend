# backend [![Build Status](https://travis-ci.org/eduboard/backend.svg?branch=master)](https://travis-ci.org/eduboard/backend)
The backend for eduBoard, written in Go with [httprouter](https://github.com/julienschmidt/httprouter) and MongoDB.

## Installing
- Install `dep`
- Run `dep ensure` to install dependencies

## Running
The backend needs MongoDB to be connected. Connection parameters can be changed using ENV.

### Docker
The easiest way to run the backend is using Docker. Just run `docker-compose up` and you are done.

## Endpoints

- `/` reservered for frontend static files

### Registration
- `/api/register` Register a new user.
- `/api/login` Login an existing user.
- `/api/logout` Logout current user.

### User
- `/api/v1/user/:id/courses` GET users courses.

### Courses
- `/api/v1/courses/` GET all accessible courses
- `/api/v1/courses/:id` GET a certain course
