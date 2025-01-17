# Simple IAM

## Overview

This is just a simple IAM application made with Go. I'm not even sure if I can call this IAM since there also endpoints to get the user list and detail. This is my first time learning Go. Therefore, I am open to any suggestions and feedbacks. Have fun.

## How to use

### Local

If you want to run this project locally, you can do so by first installing all the dependencies required:

```
go get .
```

Then, start the development environment by running the following command:

```
go run main.go

// or

go run .
```

### Production

If you want to run the production version of this project locally, you can do so by first building it:

```
make build
```

And then starting the environment by running the following command:

```
make run
```

## Development Requirements

- MySQL
- Docker
- Docker Compose
- GNU Make
- IDE

## Environment variables

| Name          | Optional | Description                                         |
| ------------- | -------- | --------------------------------------------------- |
| DB_USER       | no       | username used in MySQL DB                           |
| DB_PASSWORD   | no       | password used in MySQL DB                           |
| DB_URL        | no       | url for the MySQL DB used (example: localhost:3306) |
| DB_NAME       | no       | the name of the database                            |
| CLIENT_SECRET | no       | the secret code for client                          |
