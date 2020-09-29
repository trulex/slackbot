FROM golang:1.15-alpine

RUN apk add --no-cache git build-base

WORKDIR /go/src/app
COPY . .
COPY TOKEN /etc/slackbot/TOKEN

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["/go/bin/app"]