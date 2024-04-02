# event planner API

This is a simple event planner API written in Go. It currently only works with
http/1.1.

## CGo

This project uses a CGo SQLite library for its DB

## API Endpoints

This api currently provides 9 endpoints. The root address and path of for the API
is `http://<host_address>/api`. The endpoints are:
- `GET <api_root>/events/`                  Gets all the stored events in the system
- `GET <api_root>/events/<id>`              Gets information pertaining to a specific event
- `POST <api_root>/signup/`                 Creates a new user.
- `POST <api_root>/login`                   Logs in using email+password. The response will contain a JWT token for
  authentication.

The routes below require authentication acquired from the `/api/login` endpoint.
- `POST <api_root>/events/`                 Creates a new event. Requires authentication.
- `PUT <api_root>/events/<id>`              Updates the information for an event. Only permitted for the event creator.
- `DELETE <api_root>/events/<id>`           Deletes the event from the system. Only permitted for the event creator.
- `POST <api_root>/events/<id>/register`    Registers the logged-in user to attend the event.
- `DELETE <api_root>/events/<id>/register`  Deletes registration for the logged-in user for the event

For all the endpoints, the `<id>` segment signifies the ID for an event
stored in the system database. These IDs can be viewed in the payload
by calling the `GET <api_root>/events/` endpoint.

## Authentication

All non-get endpoints (other than /signup and /login) require authentication. This API uses
JWT tokens to authenticate users. The token has a lifespan of 2 hours. The token is aquired
by making a successful request to the `POST /login` endpoint. A request to `/login` should
look as follows:
```http request
GET /api/login HTTP/1.1
Content-Type: application/json; charset=utf-8
Accept: application/json

{"email": "example@test.com", "password": "test123"}
```
The response will look as follows if the credentials matched a user:
```http request
200 OK

{"message": "login successful!", "token": "<TOKEN>"}
```

The token returned in the response is then used for all the endpoints that require
authentication. For example with the endpoint to create events:

```http request
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

## Mandatory Request Payloads

Some endpoints have mandatory payloads. All endpoints that require authentication
expect a JWT token in the `Authorization` header.

The following endpoints require some payload in the request body:

### `POST <api_root>/events`

The `/api/events` endpoint expects the body to contain a JSON payload. The JSON
should conform to the following:
```
{
  "name": String,
  "description": String,
  "location": String,
  "date_time": RFC3339 compliant date-time String
}
```

If the request is successful the server gives a 200 OK response with a JSON body
payload. The payload will have two fields, a `message` and a `event` field. The
message field is just a simple confirmation that the request worked, and the event
field will contain the created event.

### `PUT <api_root>/events/<id>`

This endpoint updates the event with the id number of the last segment of the
path. The valid JSON payload should contain any subset of the fields for the
`POST /api/events` endpoint. So you can update any fields, from just a single
one to all of them.

The response on a successful request will contain the same payload as the
`POST /api/events/` endpoint.

### `POST <api_root>/api/register`

This endpoint is used to register as a user with the service. The endpoint
expects the body to have a JSON payload with the following fields:
```json5
{
  "email": "email address string",
  "password": "password with maximum length 30"
}
```

If successfull the response will contain two fields, a simple `message` field
and a `token` field containing an JWT token to use for authentication for other
endpoints.

### `POST <api_root>/api/login`

This endpoint expects the same payload as `/api/register` but requires
the user to have already been registered with the service. The response is also
identical.
