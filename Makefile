.PHONY: build pull push

IMAGE_NAME = gtirloni/prometheus-json-exporter

build:
	docker build -t $(IMAGE_NAME) .

pull:
	docker pull $(IMAGE_NAME)

push:
	docker push $(IMAGE_NAME)
