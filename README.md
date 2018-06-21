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

- `/` reservered for frontend static files

### Registration
- `/api/register` Register a new user.

   	```json
   	{
   	    "name": "Mathias",
   	    "surname": "Hertzel (optional)",
   	    "email": "mathias.hertzel@example.com",
   	    "password": "supersecret"
   	}
   	```

    ```json
    {
        "name": "Mathias",
        "surname": "Hertzel", 
        "email": "mathias.hertzel@example.com",
        "password": "supersecret"
    }
    ```
- `/api/login` Login an existing user.

    ```json
    {
        "email": "mathias.hertzel@example.com",
        "password": "supersecret"
    }
    ```

    ```json
    {
        "name": "Mathias",
        "surname": "Hertzel",
        "email": "mathias.hertzel@example.com"
    }
    ```
- `/api/logout` Logout current user.

### User
- `/api/v1/me` GET own user (based on SessionToken).

    ```json
    {
        "id": "12345",
        "name": "Mathias",
        "surname": "Hertzel",
        "email": "mathias.hertzel@gmail.com"
    }
    ```
- `/api/v1/users/:id` GET users.

    ```json
    {
        "id": "12345",
        "name": "Mathias",
        "surname": "Hertzel",
        "email": "mathias.hertzel@gmail.com"
    }
    ```
- `/api/v1/users/:id/courses` GET all courses a member is subscribed to

    ```json
    [
        {
            "id:": "1",
            "title": "Course 1",
            "description": "a short description"
        },
        {
            "id": "2",
            "title": "Course 2",
            "description": "another short description"
        }
    ]
    ```

### Courses
- `/api/v1/courses/` GET all accessible courses

    ```json
    [
        {
            "id:": "1",
            "title": "Course 1",
            "description": "a short description"
        },
        {
            "id": "2",
            "title": "Course 2",
            "description": "another short description"
        }
    ]
    ```
- `/api/v1/courses/:id` GET a certain course

    ```json
    {
        "id": "1",
        "title": "Course 1",
        "description": "a short description",
        "members":
        [
            "12345",
            "12346"
        ],
        "labels":
        [
            "123",
            "124"
        ],
        "entries":
        [
            {
                "id": "1",
                "date": "2018-06-20 15:04:05",
                "message": "test test",
                "pictures":
                [
                    "https://example.com/picture1.png",
                    "https:/example.com/picture2.jpg"
                ]
            },
            {
                "id": "2",
                "date": "2018-06-20 15:12:48",
                "message": "hello world",
                "pictures":
                [
                    "https://example.com/picture1.png",
                    "https:/example.com/picture2.jpg"
                ]
            }
        ]
    }
    ```
- `/api/v1/courses/:id/users` GET all subscribed members from a course

    ```json
    [
        {
            "id": "12345",
            "name": "Mathias",
            "surname": "Hertzel",
            "email": "mathias.hertzel@gmail.com"
        },
        {
            "id": "12345",
            "name": "Mathias",
            "surname": "Hertzel",
            "email": "mathias.hertzel@gmail.com"
        }
    ]
    ```
- `/api/v1/courses/:id/users/subscribe` POST user ids that will be subscribed to the course
    Input
    ```json
    [
        {
            "id": "1"
        },
        {
            "id": "2"
        }
    ]
    ```

- `/api/v1/courses/:id/users/unsubscribe` POST user ids that will be unsubscribed from the course
    Input
    ```json
    [
        {
            "id": "1"
        },
        {
            "id": "2"
        }
    ]
   ```

