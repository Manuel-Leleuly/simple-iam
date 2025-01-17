# Log In API Spec

## Log In

Endpoint: POST /api/login

Request Body:

```json
{
  "email": "johndoe@example.com",
  "password": "johndoe123"
}
```

Response Header:

- Set-Cookie: access_token, refresh_token

Response Body:

```json
{
  "data": {
    "status": "success",
    "access_token": "accessToken123",
    "refresh_token": "refreshToken123"
  }
}
```

## Refresh Token

Endpoint: POST /api/login/refresh

Request Header:

- Authorization: refresh_token

Response Header:

- Set-Cookie: access_token, refresh_token

Response Body:

```json
{
  "data": {
    "status": "success",
    "access_token": "accessToken123",
    "refresh_token": "refreshToken123"
  }
}
```
