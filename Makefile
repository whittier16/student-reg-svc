VERSION=0.1.0
LDFLAGS=-ldflags "-X main.version=`date -u +${VERSION}.build.%Y%m%d.%H%M%S`"
BUILDFLAGS=-tags=gorillamux
NAME=student-reg-svc
MAIN=cmd/${NAME}/main.go

all: app

app:
	go build -race ${BUILDFLAGS} ${LDFLAGS} -o ${NAME} ${MAIN}

dist:
	go build ${BUILDFLAGS} ${LDFLAGS} -o -v ${NAME} ${MAIN}

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o ${NAME} ${BUILDFLAGS} ${LDFLAGS} ${MAIN}

clean:
	go clean -x ${MAIN}
	rm -f ${NAME}

test: stopdb startdb
	ginkgo -v --fail-fast ${BUILDFLAGS} -r .

ci-test:
	# test without starting DB
	ginkgo --fail-fast ${BUILDFLAGS} -r --randomize-suites --trace --race --show-node-events -cover -covermode=atomic -coverprofile=coverage.out --output-dir=/tmp/

test-cover:
	go test -v ./internal/... ./cmd/... -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out

lint:
	golint internal/... cmd/... test/...

vet:
	go vet ${BUILDFLAGS} ./internal/...
	go vet ${BUILDFLAGS} ./cmd/...
	go vet ${BUILDFLAGS} ./test/...

startdb:
	docker-compose -p ${NAME}-dbonly -f test/docker/docker-compose.yml up -d

stopdb:
	docker-compose -p ${NAME}-dbonly -f test/docker/docker-compose.yml down

build-docker:
	scripts/build-docker.sh ${NAME}

push-docker:
	scripts/push-docker.sh ${NAME}

swagger:
	api/run.sh

run:
	go run ${MAIN}

.PHONY: app linux clean test lint docker publish run
