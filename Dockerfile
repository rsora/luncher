FROM balenalib/raspberry-pi-debian-golang:latest

WORKDIR /go/src/app
COPY . .

RUN go build ./...
RUN go install -v ./...

EXPOSE 8000

CMD ["luncher"]