FROM golang:alpine

COPY /bin/pigeon-mqtt /
WORKDIR /
EXPOSE 5151

CMD ["./pigeon-mqtt"]
