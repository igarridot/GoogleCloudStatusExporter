tag_name = norbega/
golang_version := 1.17.1
VERSION ?= v2.0.0-rc2
PROJECT ?= gcp-status-exporter

build-osx:
	cd src && env GOOS=darwin GOARCH=amd64 go build -o ../bin/$(PROJECT)-osx

build-linux:
	cd src && env GOOS=linux GOARCH=amd64 go build -o ../bin/$(PROJECT)-linux

build:
	@docker build -t $(tag_name)$(PROJECT):$(VERSION) -f docker/Dockerfile \
		--build-arg GOLANG_VERSION=$(golang_version) \
		--build-arg PROJECT_NAME=$(PROJECT) \
		.
push:
	docker push $(tag_name)$(PROJECT):$(VERSION)

install-requirements:
	cd src && go get -u
	cd src && go mod init || True
	cd src && go mod tidy

tests:
	cd src && go test

run-local:
	docker run -d --name gcp-exporter -p '9118:9118' $(tag_name)$(PROJECT):$(VERSION) '--exporter.collect-resolved-incidents' '--exporter.save-last-update'

stop-local:
	docker rm -f  gcp-exporter

create-tag:
	git tag -a $(VERSION)
	git push origin --tags

.PHONY: build push run-local stop-local create-tag install-requirements tests
