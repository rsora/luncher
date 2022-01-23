FROM balenalib/raspberry-pi-debian-golang:latest

WORKDIR /go/src/app
COPY . .

RUN go build ./...
RUN go install -v ./...

CMD ["luncher"]