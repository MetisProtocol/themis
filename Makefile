# Fetch git latest tag
LATEST_GIT_TAG:=$(shell git describe --tags $(git rev-list --tags --max-count=1))
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/metis-seq/themis/version.Name=themis \
		  -X github.com/metis-seq/themis/version.ServerName=themisd \
		  -X github.com/metis-seq/themis/version.ClientName=themiscli \
		  -X github.com/metis-seq/themis/version.Version=$(VERSION) \
		  -X github.com/metis-seq/themis/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/version.Name=themis \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=themisd \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=themiscli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

clean:
	rm -rf build

tests:
	# go test  -v ./...

	go test -v ./app/ ./auth/ ./sidechannel/ ./bank/ ./chainmanager/  ./staking/ -cover -coverprofile=cover.out -parallel 1

# make build
build: clean
	mkdir -p build
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o build/themisd ./cmd/themisd
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o build/themiscli ./cmd/themiscli
	CGO_ENABLED=1 go build $(BUILD_FLAGS) -o build/themis-bridge ./bridge
	@echo "====================================================\n==================Build Successful==================\n===================================================="

# make install
install:
	go install $(BUILD_FLAGS) ./cmd/themisd
	go install $(BUILD_FLAGS) ./cmd/themiscli
	go install $(BUILD_FLAGS) ./bridge

contracts:
	abigen --abi=contracts/stakemanager/stakemanager.abi --pkg=stakemanager --out=contracts/stakemanager/stakemanager.go
	abigen --abi=contracts/stakinginfo/stakinginfo.abi --pkg=stakinginfo --out=contracts/stakinginfo/stakinginfo.go
	abigen --abi=contracts/sequencerset/sequencerset.abi --pkg=sequencerset --out=contracts/sequencerset/sequencerset.go
	abigen --abi=contracts/erc20/erc20.abi --pkg=erc20 --out=contracts/erc20/erc20.go

build-arm: clean
	mkdir -p build
	env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ go build $(BUILD_FLAGS) -o build/themisd ./cmd/themisd
	env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ go build $(BUILD_FLAGS) -o build/themiscli ./cmd/themiscli
	env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ go build $(BUILD_FLAGS) -o build/themis-bridge ./bridge
	@echo "====================================================\n==================Build Successful==================\n===================================================="

#
# Code quality
#

LINT_COMMAND := $(shell command -v golangci-lint 2> /dev/null)
lint:
ifndef LINT_COMMAND
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2
endif
	golangci-lint run --config ./.golangci.yml

#
# docker commands
#

build-docker:
	@echo Fetching latest tag: $(LATEST_GIT_TAG)
	git checkout $(LATEST_GIT_TAG)
	docker build -t "metis-seq/themis:$(LATEST_GIT_TAG)" -f docker/Dockerfile .

push-docker:
	@echo Pushing docker tag image: $(LATEST_GIT_TAG)
	docker push "metis-seq/themis:$(LATEST_GIT_TAG)"

build-docker-develop:
	docker build -t "metis-seq/themis:develop" .

.PHONY: contracts build

PACKAGE_NAME          := github.com/metis-seq/themis
GOLANG_CROSS_VERSION  ?= v1.19.1

.PHONY: release-dry-run
release-dry-run:
	@docker run \
		--platform linux/amd64 \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e CGO_CFLAGS=-Wno-unused-function \
		-e GITHUB_TOKEN \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

.PHONY: release
release:
	@docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e GITHUB_TOKEN \
		-e DOCKER_USERNAME \
		-e DOCKER_PASSWORD \
		-e SLACK_WEBHOOK \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v $(HOME)/.docker/config.json:/root/.docker/config.json \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate
