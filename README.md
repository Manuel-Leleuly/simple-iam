# Simple IAM

## Overview

This is just a simple IAM application made with Go. I'm not even sure if I can call this IAM since there also endpoints to get the user list and detail. This is my first time learning Go. Therefore, I am open to any suggestions and feedbacks. Have fun.

## How to use

### Local

If you want to run this project locally, you can do so by first installing all the dependencies required:

```
go mod download
```

Then, start the development environment by running the following command:

```
go run main.go

// or

go run .
```

### Local (with docker)

If you want to run this project using docker, you can do so by running the following command:

```
make run
```

## Development Requirements

- Go (this project was created using Go 1.23.3)
- MySQL
- IDE
- optional
  - Docker
  - Docker Compose
  - GNU Make

## Environment variables

| Name          | Optional | Description                          |
| ------------- | -------- | ------------------------------------ |
| DB_USER       | no       | username used in MySQL DB            |
| DB_PASSWORD   | no       | password used in MySQL DB            |
| DB_HOST       | no       | host used in MySQL                   |
| DB_PORT       | no       | port used in MySQL                   |
| DB_NAME       | no       | the name of the database             |
| DB_TEST_NAME  | yes      | the name of the database for testing |
| CLIENT_SECRET | no       | the secret code for client           |
