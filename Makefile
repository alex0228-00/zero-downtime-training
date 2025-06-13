
IMAGE_NAME = zero-downtime-training
NETWORK = zero-downtime-training

build-v1:
	docker build --target v1 -t $(IMAGE_NAME):v1 .

docker-create-network:
	docker network create $(NETWORK)

docker-deploy-mysql:
	docker run -d \
		--name mysql \
		--network $(NETWORK) \
		-e MYSQL_ROOT_PASSWORD=rootpwd \
		-e MYSQL_DATABASE=assets \
		-e MYSQL_USER=testuser \
		-e MYSQL_PASSWORD=testpassword \
		-p 3306:3306 \
		mysql:8.0