IMAGE_NAME := zero-downtime-training:test

build:
	docker build -t $(IMAGE_NAME) .

test: build
	docker run --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		$(IMAGE_NAME)
