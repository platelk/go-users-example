# go-users-example

This project aims to provide an example of golang architecture for a simple project.

## Build

To do a local build, just do:

```
go build
```

To generate a docker image:

```
docker build -t go-users-example .
```

## Test

All the test are based on golang test package:

```
go test .
```

## Run

The service can be launch directly based on default values

```
$> ./go-users-example
```

and can be configured through env variable:

* `HTTP_ADDR`: the listen string representation like ":8080"
* `LOG_LEVEL`: define the level of log. default is `info`

## Architecture principles

This repository try to provide a possible golang service architecture which is describe [here](./doc/ddd.md)

## How to go from there ?

This repo show a simple user service which aims to be extendable in the future with business changes.
This service was written from small user base where read and writes are similar.
In the case of scale in request, we may want to do some changes:
* switch to a delayed creation and update trough event first which can be sent to a kafka and later write into the system
  to allow better handling of burst
* split read and write request and move all "search" features to a dedicated service
* Use event from creation and update to store the user to a better datastore for search capabilities

Some obvious features hasn't been implemented too because of time and complexity for my aim:
* Authentication and Authorization to create/update/delete
* Proper store
* Better data validation
* Monitoring

## Example

### Create

```
$> http POST :8080/v1/user first_name=plop last_name=test nick_name=plop email=test@test.com -v

POST /v1/user HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 90
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "email": "test@test.com",
    "first_name": "plop",
    "last_name": "test",
    "nick_name": "plop"
}

HTTP/1.1 200 OK
Content-Length: 223
Content-Type: text/plain; charset=utf-8
Date: Sun, 11 Oct 2020 21:40:17 GMT

{
    "user": {
        "country": "",
        "email": "test@test.com",
        "first_name": "plop",
        "id": "86fcf3cd-a280-4356-8fc5-abb1eef103b5",
        "last_name": "test",
        "nick_name": "plop",
        "password": "$2a$12$9ljtWwaw3TijaU8vcB0Lau/vWX8YOH3At67dr4dgKfZ9wl/jlePc2"
    }
}

```

### Update

```
$>  http PUT :8080/v1/user id=86fcf3cd-a280-4356-8fc5-abb1eef103b5 email=updated-email@test.com -v

PUT /v1/user HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 81
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "email": "updated-email@test.com",
    "id": "86fcf3cd-a280-4356-8fc5-abb1eef103b5"
}

HTTP/1.1 200 OK
Content-Length: 232
Content-Type: text/plain; charset=utf-8
Date: Sun, 11 Oct 2020 21:40:38 GMT

{
    "user": {
        "country": "",
        "email": "updated-email@test.com",
        "first_name": "plop",
        "id": "86fcf3cd-a280-4356-8fc5-abb1eef103b5",
        "last_name": "test",
        "nick_name": "plop",
        "password": "$2a$12$9ljtWwaw3TijaU8vcB0Lau/vWX8YOH3At67dr4dgKfZ9wl/jlePc2"
    }
}
```

### Search

```
$> http :8080/v1/users first_name==plop -v

GET /v1/users?first_name=plop HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8080
User-Agent: HTTPie/1.0.3



HTTP/1.1 200 OK
Content-Length: 658
Content-Type: text/plain; charset=utf-8
Date: Sun, 11 Oct 2020 21:44:12 GMT

{
    "users": [
        {
            "country": "",
            "email": "test@test.com",
            "first_name": "plop",
            "id": "2db1c029-f8d0-4cac-ae3e-b0ede6b2ea32",
            "last_name": "test",
            "nick_name": "plop",
            "password": "$2a$12$CsBJv7rwmAxtF07byuA.WeA.V/08N4iy4x72VVh./h5fGSIIttiEu"
        },
        {
            "country": "",
            "email": "test2@test.com",
            "first_name": "plop",
            "id": "0066948d-3f4a-4fdd-b3ce-b7f01841c5fb",
            "last_name": "test",
            "nick_name": "plop",
            "password": "$2a$12$V5e278peZm.fnkO94KW/UeeQUjnuSlZ7gKfFE5PlbDYGVureZ8sf."
        },
        {
            "country": "",
            "email": "test4@test.com",
            "first_name": "plop",
            "id": "c09fbd7b-1b28-48c1-9c5e-74e4bcd013be",
            "last_name": "test",
            "nick_name": "plop",
            "password": "$2a$12$OTsMAICpwkBDfP0QK.KDYOVZ.EPwZg80YSd3W7VmUAd3tRfbz3oQK"
        }
    ]
}
```

### Delete

```
$> http DELETE :8080/v1/user id=2db1c029-f8d0-4cac-ae3e-b0ede6b2ea32 -v

DELETE /v1/user HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 46
Content-Type: application/json
Host: localhost:8080
User-Agent: HTTPie/1.0.3

{
    "id": "2db1c029-f8d0-4cac-ae3e-b0ede6b2ea32"
}

HTTP/1.1 200 OK
Content-Length: 223
Content-Type: text/plain; charset=utf-8
Date: Sun, 11 Oct 2020 21:45:15 GMT

{
    "user": {
        "country": "",
        "email": "test@test.com",
        "first_name": "plop",
        "id": "2db1c029-f8d0-4cac-ae3e-b0ede6b2ea32",
        "last_name": "test",
        "nick_name": "plop",
        "password": "$2a$12$CsBJv7rwmAxtF07byuA.WeA.V/08N4iy4x72VVh./h5fGSIIttiEu"
    }
}
```
