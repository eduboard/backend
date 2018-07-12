# Endpoints

- `/` reservered for frontend static files

## Registration
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

## User
- `/api/v1/me` GET own user (based on SessionToken).

    ```json
    {
        "id": "12345",
        "name": "Mathias",
        "surname": "Hertzel",
        "email": "mathias.hertzel@gmail.com"
    }
    ```
- `/api/v1/users` GET all users

    ```json
    [
        {
            "id": "12345",
            "name": "Mathias",
            "surname": "Herzel",
            "email": "mathias.hertzel@gmail.com",
            "profilePicture": "http://example.com"
        }
    ]

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
            "description": "another short description",
            "entries":
            [
                {
                    "id": "1",
                    "date": "2018-07-01T15:04:05-07:00",
                    "message": "test test",
                    "pictures":
                    [
                        "https://example.com/picture1.png",
                        "https:/example.com/picture2.jpg"
                    ]
                }
            ],
            "schedules": 
            [
                {
                    "day": 0,
                    "startsAt": "2018-06-24T08:52:28.764Z",
                    "room": "EN 154",
                    "title": "meme Class"
                }
            ]
        }
    ]
    ```
     _Remarks:_ Day 0 indicates Sunday, 1 is Monday and so on...
     
## Courses
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
                "date": "2018-07-01T15:04:05-07:00",
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
        ],
        "schedules": 
        [
            {
                "day": 0,
                "startsAt": "2018-06-24T08:52:28.764Z",
                "room": "EN 154",
                "title": "meme Class"
            }
        ]
    }
    ```
    _Remarks:_ Day 0 indicates Sunday, 1 is Monday and so on...
    
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
- `/api/v1/courses/:courseId/entries` POST new entry that will be added to the course

    Input
    ```json
    {
        "date": "2018-07-01T15:04:05-07:00",
        "message": "herpy derpy new messagerpy",
        "pictures":
        [
            "https://example.com/course/picture1",
            "https://example.com/course/picture2"
        ],
        "published": true
    }
    ```
    Output
    ```json
    {
        "id": "1",
        "date": "2018-07-01T15:04:05-07:00",
        "message": "herpy derpy new messagerpy",
        "pictures":
        [
            "https://example.com/course/picture1",
            "https://example.com/course/picture2"
        ],
        "published": true
    }
    ```
- `/api/v1/courses/:courseId/entries/:entryId` PUT updates in entry from a course (not yet implemented)

    Input
    ```json
    {
        "date":"2018-07-01T20:04:05-07:00",
        "message": "herpy derpy",
        "pictures":
        [
            "https://example.com/course/picture1",
            "https://example.com/course/picture2"
        ],
        "published": false
    }
    ```
    Output
    ```json
    {
        "id": "1",
        "date": "2018-07-01T20:04:05-07:00",
        "message": "herpy derpy",
        "pictures":
        [
            "https://example.com/course/picture1",
            "https://example.com/course/picture2"
        ],
        "published": false
    }
    ```
- `/api/v1/courses/:courseID/entries/:entryId` DELETE selected entry from a course
