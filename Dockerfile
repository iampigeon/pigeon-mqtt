FROM golang:alpine

COPY /bin/pigeon-mqtt /
WORKDIR /
EXPOSE 9010

CMD ["./pigeon-mqtt"]
