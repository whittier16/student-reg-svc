FROM golang:1.19 AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -a -o student-reg-svc ./cmd/student-reg-svc/main.go

EXPOSE 5005
CMD [ "/app/student-reg-svc" ]
