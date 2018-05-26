# backend
The backend for eduBoard, written in Go with [httprouter](https://github.com/julienschmidt/httprouter) and PostgreSQL.

## Installing
- Run `dep eunsure` to install dependencies

## Endpoints

- `/` reservered for frontend static files

### User
- `/api/v1/register` Register a new user.
- `/api/v1/login` Login an existing user.
- `/api/v1/user/:id/courses/` GET users courses.

### Courses
- `/api/v1/courses/` GET all accessible courses
- `/api/v1/courses/:id/` GET a certain course
