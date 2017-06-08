all: build docker

build:
	@GOOS=linux GOARCH=amd64 go build -o cloud-dyndns.linux-amd64

docker:
	docker build -t cloud-dyndns .
