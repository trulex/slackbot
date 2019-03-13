FROM golang:1.8

WORKDIR /go/src/app
COPY . .
COPY TOKEN /etc/slackbot/TOKEN

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["/go/src/app/slackbot"]