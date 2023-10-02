# Instructions to test out Auth server

## Build mysql server
```bash
cd db/auth
docker build --tag auth-db .
```

## Build the Auth server 
```bash
make auth
```

> create a folder in the project root `.env-auth`
```.env
DB_URL=auth-db:3306
DB_PASSWORD=12345
```

## Finally make the docker compose up

```bash
# for starting
docker compose up

# for stopping
docker compose down
```

> Open another terminal to test out the functionality

```bash
# check whether we can reach the backend-auth
curl -X GET 0.0.0.0:8080/account | jq -r .
```
response
```json
{
  "stdout": "",
  "error": "",
  "Account": {
    "Loc": {
      "[GET] Health status": "/account/healthz",
      "[GET] authorization bearer token": "/account/token/status",
      "[GET] token renew": "/account/renew",
      "[POST] logout": "/account/logout",
      "[POST] signin": "/account/signin",
      "[POST] signup": "/account/signup"
    }
  }
}
```

```bash
# signup
curl -i -X POST 0.0.0.0:8080/account/signup -d '{"username": "crop123", "name": "Dummy Name", "email": "123@gmail.com", "password": "1111"}'

# signin
curl -X POST 0.0.0.0:8080/account/signin -d '{"username":"crop123", "password":"1111"}'

```

> recieved a token as response
```json
{
  "stdout": "Login Successful",
  "error": "",
  "Account": {
    "token": "<JWT-Token>"
  }
}
```

how to use it for authentication

```bash
curl -i -X GET 0.0.0.0:8080/account/token -H "Authorization: Bearer <JWT-Token>"

HTTP/1.1 200 OK
Accept: application/json; charset=utf-8
Content-Type: application/json; charset=utf-8
Server: authentication-server
Date: Sat, 30 Sep 2023 15:40:05 GMT
Content-Length: 55

{"stdout":"crop123","error":"","Account":null}
```



# Instructions to test out Image server

## Build all resources
```bash
make build
```

## Build the image server 
```bash
make image
```

> create a folder in the project root `.env-image`
```.env
DB_URL=image-db:6379
DB_PASSWORD=12345
```

## Finally make the docker compose up

```bash
# for starting
make run-watch

# for stopping
docker compose down
```

> Open another terminal to test out the functionality

```bash
# check whether we can reach the backend-auth
curl -X GET 0.0.0.0:8090/image | jq -r .
```
response
```json
{
  "stdout": "",
  "error": "",
  "Account": {
    "Loc": {
      "[GET] Health status": "/account/healthz",
      "[GET] authorization bearer token": "/account/token/status",
      "[GET] token renew": "/account/renew",
      "[POST] logout": "/account/logout",
      "[POST] signin": "/account/signin",
      "[POST] signup": "/account/signup"
    }
  }
}
```

```bash
# signup
curl -i -X POST 0.0.0.0:8080/account/signup -d '{"username": "crop123", "name": "Dummy Name", "email": "123@gmail.com", "password": "1111"}'

# signin
curl -X POST 0.0.0.0:8080/account/signin -d '{"username":"crop123", "password":"1111"}'

```

> recieved a token as response
```json
{
  "stdout": "Login Successful",
  "error": "",
  "Account": {
    "token": "<JWT-Token>"
  }
}
```

how to use it for authentication

```bash
curl -i -X GET 0.0.0.0:8080/image/get -H "Authorization: Bearer <JWT-Token>"

HTTP/1.1 200 OK
Accept: application/json; charset=utf-8
Content-Type: application/json; charset=utf-8
Server: authentication-server
Date: Sat, 30 Sep 2023 15:40:05 GMT
Content-Length: 55

{"stdout":"crop123","error":"","Account":null}
```


```bash
curl -i -X POST 0.0.0.0:8090/image/upload -H "Authorization: Bearer <JWT-Token>" -d '{"raw_image_bytes": "<Image must be in []bytes format in the request payload>", "image_format": "image/jpeg | image/png"}'
```
