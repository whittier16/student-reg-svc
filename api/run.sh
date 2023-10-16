#!/bin/sh

echo "listening on http://localhost:8060"

open_browser() {
  sleep 1
  open "http://localhost:8060"
}

open_browser &

docker run --name swagger-ui --rm -p 8060:8080 -e "SWAGGER_JSON=/api/openapi.yaml" -v "${HOME}/Github/whittier16/student-reg-svc/api:/api" swaggerapi/swagger-ui

