tag_name = norbega/
golang_version := 1.25
VERSION ?= v2.0.0-rc4
PROJECT ?= gcp-status-exporter
GOARCH ?= amd64

build-osx:
	cd src && env GOOS=darwin GOARCH=amd64 go build -o ../bin/$(PROJECT)-osx

build-osx-arm64:
	cd src && env GOOS=darwin GOARCH=arm64 go build -o ../bin/$(PROJECT)-osx-arm64

build-linux:
	cd src && env GOOS=linux GOARCH=amd64 go build -o ../bin/$(PROJECT)-linux

build:
	@docker build -t $(tag_name)$(PROJECT):$(VERSION) -f docker/Dockerfile \
		--build-arg GOLANG_VERSION=$(golang_version) \
		--build-arg PROJECT_NAME=$(PROJECT) \
		--build-arg GOARCH=$(GOARCH) \
		.

build-arm64:
	@docker build --platform linux/arm64 -t $(tag_name)$(PROJECT):$(VERSION) -f docker/Dockerfile \
		--build-arg GOLANG_VERSION=$(golang_version) \
		--build-arg PROJECT_NAME=$(PROJECT) \
		--build-arg GOARCH=arm64 \
		.

push:
	docker push $(tag_name)$(PROJECT):$(VERSION)

push-arm64:
	docker push $(tag_name)$(PROJECT):$(VERSION)

install-requirements:
	cd src && go get -u
	cd src && go mod init || True
	cd src && go mod tidy

tests:
	cd src && go test ./...

run-local:
	docker run -d --name gcp-exporter -p '9118:9118' $(tag_name)$(PROJECT):$(VERSION) '--exporter.collect-resolved-incidents' '-exporter.filtered-products Apigee' '--exporter.save-last-update'

stop-local:
	docker rm -f  gcp-exporter

create-tag:
	git tag -a $(VERSION)
	git push origin --tags

.PHONY: build build-arm64 push push-arm64 run-local run-local-arm64 stop-local stop-local-arm64 create-tag install-requirements tests build-osx-arm64
