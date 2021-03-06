APP?=orchidgo
RELEASE?=0.0.1
GOOS?=linux

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: check
check: prepare_linter
	# golangcli-lint run -v
	# gometalinter --vendor ./...

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${GOOS} go build -o bin/${APP} \
		-ldflags "-X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \

.PHONY: clean
clean:
	@rm -f bin/${APP}

.PHONY: vendor
vendor: prepare_dep
	dep ensure

HAS_DEP := $(shell command -v dep;)
HAS_LINTER := $(shell command -v golangci-lint;)

.PHONY: prepare_dep
prepare_dep:
ifndef HAS_DEP
	go get -u -v -d github.com/golang/dep/cmd/dep && \
	go install -v github.com/golang/dep/cmd/dep
endif

.PHONY: prepare_linter
prepare_linter:
ifndef HAS_LINTER
	go get -u -v -d github.com/golangci/golangci-lint/cmd/golangci-lint && \
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint
endif
