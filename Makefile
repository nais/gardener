SHELL   := bash
NAME    := navikt/gardener
LATEST  := ${NAME}:latest
DEP_IMG := navikt/dep:1.0.0
DEP     := docker run --rm -v ${PWD}:/go/src/github.com/nais/naisd -w /go/src/github.com/nais/naisd ${DEP_IMG} dep
GO_IMG  := golang:1.8
GO      := docker run --rm -v ${PWD}:/go/src/github.com/nais/gardener -w /go/src/github.com/nais/gardener ${GO_IMG} go
LDFLAGS := -X github.com/nais/gardener/version.Revision=$(shell git rev-parse --short HEAD) -X github.com/nais/gardener/version.Version=$(shell /bin/cat ./vrs)

dockerhub-release: install test linux bump tag docker-build push-dockerhub

tag:
	git tag -a $(shell /bin/cat ./vrs) -m "auto-tag from Makefile [skip ci]" && git push --tags

install:
	${DEP} ensure

test:
	${GO} test .

build:
	${GO} build -o gardener

linux:
	docker run --rm \
		-e GOOS=linux \
		-e CGO_ENABLED=0 \
		-v ${PWD}:/go/src/github.com/nais/gardener \
		-w /go/src/github.com/nais/gardener ${GO_IMG} \
		go build -a -installsuffix cgo -ldflags '-s $(LDFLAGS)' -o gardener

docker-build:
	docker image build -t ${NAME}:$(shell /bin/cat ./vrs) -t gardener -t ${NAME} -t ${LATEST} -f Dockerfile .

push-dockerhub:
	docker image push ${NAME}:$(shell /bin/cat ./vrs)