all: build docker

build:
	@GOOS=linux GOARCH=amd64 go build -o cloud-dyndns.linux-amd64

docker:
	@docker build -t cloud-dyndns .

run:
	@docker run \
	--rm \
	--name=cloud-dyndns \
	-e "AWS_ACCESS_KEY=AKAAKSJHDAKSJHDK" \
	-e "AWS_SECRET_KEY=sdjfkhgasjdhfakjsdhgfkjahsdgfkjhag" \
	cloud-dyndns \
	-frequency=60m \
	-zoneid=XXXXXXXXXXXXX \
	-hostname=host.domain.com. \
	-ttl=900 \
	-debug=true
