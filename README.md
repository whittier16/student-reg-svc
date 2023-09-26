# student-reg-svc

## Objective

To provide the backend that enables teachers to efficiently manage student data

## Background
Teachers need a system where they can perform administrative functions for their students. Teachers and students are identified by their email addresses.

### Prerequisites

1. Go version 1.11 and above to support [Go Modules](https://github.com/golang/go/wiki/Modules)
2. Docker
3. Docker compose
4. MySQL

### Quick Start Guide

If you haven't used `go mod` before, make sure you set `GO111MODULE=on` in your environment variable

1. Clone the repository

```
cd ~/go/src/github.com/whittier16
git clone https://github.com/whittier16/student-reg-svc.git
cd student-reg-svc
```

2. Create the `.env` file by copying `.env-template` and replacing the values as needed

```
cp .env-template .env
```

3. Replace the values in `configs/config-*.yml` as needed

4. Create a database
```
mysql -u root
mysql>create database stdnt_reg;
```

5. Dump sql file to import into the database
```
mysql -u root
cd db/migration
mysql --protocol=tcp --host=127.0.0.1 --user=root --port=3306 --default-character-set=utf8 --comments --database=stdnt_reg --password=<password> < "db/migration/dump.sql"
```

6. Build the binary

```
make
```

7. Run the binary

```
./student-reg-svc
```

8. Run local manually [http://localhost:5005](http://localhost:5005)

Alternatively, you can run the web server via docker. Prerequesite is follow steps 1 to 3 above.
To run the api server:

1. Use docker compose to launch MySQL database

```shell
cd deployments/docker/db-only
docker-compose up
```

2. Grant permission to user

```shell
mysql -u root
mysql>ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';
```

3. On a separate terminal, build the image

```
make build-docker
```
[Dockerfile](build%2Fpackage%2FDockerfile)

4. Run the image

```
 docker run --publish 5005:5005 docker.io/library/student-reg-svc:latest
```

5. Once everything has started up, you should be able to access the webapp via http://localhost:5005/ on your host machine.

```
open http://localhost:5005/
```

### Makefile commands

- `make` build the binary
- `make linux` build the binary in linux format to be used in Docker container
- `make clean` clean the build
- `make test` run the `ginkgo` unit tests
- `make ci-test` run the test without starting the redis server for use in CI/CD integration
- `make lint` run the linter
- `make vet` run `go vet` to analyze source code
- `make startdb` start the redis server via docker-compose
- `make stopdb` stop the redis server
- `make build-docker` to build the docker image
- `make push-docker` to push the docker image to `quay.io`. WARNING: This should only be used by CI/CD
- `make swagger` view the OpenAPI Specs

### API

- Development server: `http://localhost:5005`

#### `POST /auth`

##### Sample request

```
curl --location --request POST 'localhost:5005/auth'
```

<details><summary>Success Response</summary>
<p>

```
HTTP/1.1 204 No Content
Content-Type: application/json
Token: <TOKEN>
Vary: Origin
Date: Tue, 26 Sep 2023 11:55:04 GMT
```

</p>
</details>

> **_NOTE:_** Copy the value of the `Token` header from the response headers for the rest API endpoints 
when executing API requests

#### `POST /api/commonstudents`

##### Sample request

```
curl --location 'localhost:5005/api/commonstudents?teacher=teacher1%40gmail.com&teacher=teacher2%40gmail.com' \
	--header 'Content-Type: application/json' \
	--header 'Token: {TOKEN}' 
```

<details><summary>Success Response</summary>
<p>

```
{
    "students": [
        "student2@gmail.com",
        "student3@gmail.com",
        "student4@gmail.com"
    ]
}
```

</p>
</details>
 
- You can also download the [collection file](docs%2FStudent Registration Service.postman_collection.json) from this repo, then import directly into Postman.

### TODO
- 
