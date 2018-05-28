# backend [![Build Status](https://travis-ci.org/eduboard/backend.svg?branch=master)](https://travis-ci.org/eduboard/backend)
The backend for eduBoard, written in Go with [httprouter](https://github.com/julienschmidt/httprouter) and MongoDB.

## Installing
- Run `dep eunsure` to install dependencies

## Endpoints

- `/` reservered for frontend static files

### Registration
- `/api/v1/register` Register a new user.
- `/api/v1/login` Login an existing user.
- `/api/v1/logout` Logout current user.

### User
- `/api/v1/user/:id/courses/` GET users courses.

### Courses
- `/api/v1/courses/` GET all accessible courses
- `/api/v1/courses/:id/` GET a certain course
