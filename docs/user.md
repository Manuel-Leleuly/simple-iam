# User API Spec

## Create User

Endpoint: POST /api/users

Request Body:

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe",
  "email": "johndoe@example.com",
  "password": "johndoe123"
}
```

Response Body:

```json
{
  "data": {
    "id": "c285969d442347b999b591879ba062cd",
    "first_name": "John",
    "last_name": "Doe",
    "username": "johndoe",
    "email": "johndoe@example.com"
  }
}
```

## Get User List

Endpoint: GET /api/users

Request Header:

- Authorization: access_token

Query Parameter:

- id: string
- first_name: string
- last_name: string
- email: string
- offset: number, default to 0
- limit: number, default to 10

Response Body:

```json
{
  "data": [
    {
      "id": "c285969d442347b999b591879ba062cd",
      "first_name": "John",
      "last_name": "Doe",
      "username": "johndoe",
      "email": "johndoe@example.com"
    },
    {
      "id": "20dca272cb564d2fa851ba7df6586b65",
      "first_name": "Jane",
      "last_name": "Doe",
      "username": "janedoe",
      "email": "janedoe@example.com"
    }
  ]
}
```

## Get User Detail

Endpoint: GET /api/users/{userId}

Request Header:

- Authorization: access_token

Response Body:

```json
{
  "data": {
    "id": "c285969d442347b999b591879ba062cd",
    "first_name": "John",
    "last_name": "Doe",
    "username": "johndoe",
    "email": "johndoe@example.com"
  }
}
```

## Update User

Endpoint: PATCH /api/users/{userId}

Request Header:

- Authorization: access_token

Request Body:

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe",
  "email": "johndoe@example.com"
}
```

Response Body:

```json
{
  "data": {
    "id": "c285969d442347b999b591879ba062cd",
    "first_name": "John",
    "last_name": "Doe",
    "username": "johndoe",
    "email": "johndoe@example.com"
  }
}
```

## Delete User

Endpoint: DELETE /api/users/{userId}

Request Header:

- Authorization: access_token

Response Body:

```json
{
  "data": {
    "success": true
  }
}
```
