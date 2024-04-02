# event planner API

This is a simple event planner API written in Go. It currently only works with
http/1.1.

## CGo

This project uses a CGo SQLite library for its DB

## API Endpoints

This api currently provides 9 endpoints. The root address and path of for the API
is `http://<host_address>/api`. The endpoints are:
- `GET <api_root>/events/`                  Gets all the stored events in the system
- `POST <api_root>/events/`                 Creates a new event. Requires authentication.
- `POST <api_root>/signup/`                 Creates a new user.
- `POST <api_root>/login`                   Logs in using email+password. The response will contain a JWT token for
                                            authentication.

The routes below require authentication acquired from the `/api/login` endpoint.
- `GET <api_root>/events/<id>`              Gets information pertaining to a specific event
- `PUT <api_root>/events/<id>`              Updates the information for an event. Only permitted for the event creator.
- `DELETE <api_root>/events/<id>`           Deletes the event from the system. Only permitted for the event creator.
- `POST <api_root>/events/<id>/register`    Registers the logged in user to attend the event.
- `DELETE <api_root>/events/<id>/register`  Deletes registration for the logged in user for the event

For all the endpoints, the `<id>` segment signifies the ID for an event
stored in the system database. These IDS can be viewed in the payload
by calling the `GET <api_root>/events/` endpoint.

## Authentication

All non-get endpoints (other than /signup and /login) require authentication. This API uses
JWT tokens to authenticate users. The token has a lifespan of 2 hours. The token is aquired
by making a successful request to the `POST /login` endpoint. A request to `/login` should
look as follows:
```http
GET /api/login HTTP/1.1
Content-Type: application/json; charset=utf-8
Accept: application/json

{"email": "example@test.com", "password": "test123"}
```
The response will look as follows if the credentials matched a user:
```http
200 OK

{"message": "login successful!", "token": "<TOKEN>"}
```

The token returned in the response is then used for all the endpoints that require
authentication. For example with the endpoint to create events:

```http
POST /api/events HTTP/1.1
Content-Type: application/json; charset=utf-8
Accept: application/json
Authorization: <TOKEN>

{
    "name":"Test event",
    "description":"This event is for show",
    "location":"Somewhere",
    "time":"2025-01-25T14:00:23.000Z"
}
```