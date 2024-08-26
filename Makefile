docker-build:
	docker build -t payments-server .

docker-run:
	docker run -p 8080:8080 -t payments-server