app_tag_name = norbega/gcp-status-exporter
version = 0.3.1

build:
	docker build \
		-t $(app_tag_name):$(version) -f docker/Dockerfile .

push:
	docker push $(app_tag_name):$(version)

install-test-requirements:
	pip install -r src/requirements.txt

tests:
	python -m unittest test.test_main

run-local:
	docker run -d --name gcp-exporter -p '9118:9118' $(app_tag_name):$(version)

bash-local:
	docker exec -ti gcp-exporter sh

stop-local:
	docker rm -f  gcp-exporter

create-tag:
	git tag -a v$(version)
	git push origin --tags

.PHONY: build push run-local stop-local create-tag install-test-requirements tests
