FROM balenalib/raspberry-pi-debian-golang:latest

WORKDIR /go/src/app
COPY . .

RUN go build -ldflags "-X main.build=docker" -v ./...
RUN go install -v ./...

EXPOSE 8000

CMD ["luncher-api"]