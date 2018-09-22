all: build
iothost=$(IOT_HOST)
certificate=$(PIGEON_MQTT_CERTIFICATE)
privatekey=$(PIGEON_MQTT_PRIVATE_KEY)


.PHONY: build
build b: 
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate ./...
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build\
		-ldflags='-X main.iothost=${iothost} -X main.certificate=${certificate} -X main.privateKey=${privatekey}'\
		-o bin/pigeon-mqtt main.go
build-osx bx:
	@CGO_ENABLED=0 go generate ./...
	@CGO_ENABLED=0 GOARCH=amd64 go build\
		-ldflags='-X main.iothost=${iothost} -X main.certificate=${certificate} -X main.privateKey=${privatekey}'\
		-o bin/pigeon-mqtt main.go
deploy d: build
	@docker-compose build
	@docker-compose up
