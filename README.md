# slackbot
Bot for Slack

## Installing

To start using slackbot, install Go 1.7+ and run `go get`:

```sh
$ go get github.com/trulex/slackbot
```

This will install `slackbot` in to your `$GOBIN` path.

## Running

To start slackbot, run following: 
```sh
$ slackbot -tpath=/etc/slackbot/TOKEN -debug
```
Where `tpath` is location of your Slack auth token.
