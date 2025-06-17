build:
	docker build -t zero-downtime-training .

test:
	docker run --rm zero-downtime-training