IMAGE_NAME := zero-downtime-training:test

build:
	docker build -t $(IMAGE_NAME) .

test: build
	docker run --rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		$(IMAGE_NAME)

docker-clean:
	docker system prune -af --volumes

clean-servers:
	bash -c '\
	for image in $$(docker images zero-downtime-training --format "{{.Repository}}:{{.Tag}}"); do \
		containers=$$(docker ps -a --filter ancestor=$$image --format "{{.ID}}"); \
		if [ -n "$$containers" ]; then \
			echo "Stopping $$containers from image $$image"; \
			docker stop $$containers; \
			docker rm $$containers; \
		fi; \
	done'

clean: clean-servers
	docker container stop zero-downtime-training-mysql