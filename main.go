package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/mitjaziv/slackbot/command"
	"github.com/mitjaziv/slackbot/service"
)

var commands = map[string]func() command.Command{
	"^echo ":  command.Echo,
	"^time$":  command.Time,
	"^hello$": command.Hello,
	"^menu":   command.Menu,
}

type slackBot struct {
	exit chan bool

	connector service.Connector
	commands  map[string]command.Command
	services  map[string]string
	sync.RWMutex
}

func newSlackBot(commands map[string]command.Command) *slackBot {
	commands["^help$"] = help(commands)

	return &slackBot{
		exit:      make(chan bool),
		commands:  commands,
		connector: service.NewConnector(),
		services:  make(map[string]string),
	}
}

func main() {
	// Token path
	tpath := flag.String("tpath", "/etc/slackbot/TOKEN", "token path")

	// Debug mode
	debug := flag.Bool("debug", false, "debug mode")

	// Parse flags
	flag.Parse()

	// Get token from file
	t, err := ioutil.ReadFile(*tpath)
	if err != nil {
		log.Printf("error reading slack token %v", err)
		os.Exit(1)
	}

	// Convert token to string
	token := strings.TrimSpace(string(t))

	// Create commands map
	cmds := make(map[string]command.Command)

	// register default built in commands
	for pattern, cmd := range commands {
		cmds[pattern] = cmd()
	}

	// Create bot
	b := newSlackBot(cmds)

	// Start bot
	if err := b.start(*debug, token); err != nil {
		log.Printf("error starting bot %v", err)
		os.Exit(1)
	}

	// Run bot
	if err := b.run(); err != nil {
		log.Printf("error running bot %v", err)
		os.Exit(1)
	}

	// Stop bot
	if err := b.stop(); err != nil {
		log.Printf("error stopping bot %v", err)
	}
}

func (b *slackBot) start(debug bool, token string) error {
	log.Println("[slackbot] starting")

	// Initialize connection
	if err := b.connector.Init(debug, token); err != nil {
		return err
	}

	// Start/Open connection
	if err := b.connector.Start(); err != nil {
		return err
	}

	return nil
}

func (b *slackBot) run() error {
	log.Println("[slackbot] connecting")

	c, err := b.connector.Stream()
	if err != nil {
		return err
	}

	// Process connection until exit event
	for {
		select {
		case <-b.exit:
			log.Println("[slackbot] closing connection")
			return c.Close()
		default:
			var recvEv service.Event
			// receive input
			if err := c.Recv(&recvEv); err != nil {
				return err
			}

			// only process TextEvent
			if recvEv.Type != service.TextEvent {
				continue
			}
			if len(recvEv.Data) == 0 {
				continue
			}
			if err := b.process(c, recvEv); err != nil {
				return err
			}
		}
	}
}

func (b *slackBot) stop() error {
	log.Println("[slackbot] stopping")
	close(b.exit)

	// closing connection
	if err := b.connector.Stop(); err != nil {
		log.Printf("[slackbot] %v", err)
	}

	return nil
}

func (b *slackBot) process(c *service.Connection, ev service.Event) error {
	args := strings.Split(string(ev.Data), " ")
	if len(args) == 0 {
		return nil
	}

	b.RLock()
	defer b.RUnlock()

	// try built in command or skip if it doesn't match
	for pattern, cmd := range b.commands {
		if m, err := regexp.Match(pattern, ev.Data); err != nil || !m {
			continue
		}
		// matched, exec command
		rsp, err := cmd.Exec(args...)
		if err != nil {
			rsp = []byte("error executing cmd: " + err.Error())
		}
		// send response
		return c.Send(&service.Event{
			Meta: ev.Meta,
			From: ev.To,
			To:   ev.From,
			Type: service.TextEvent,
			Data: rsp,
		})
	}

	return nil
}

// Crete help command
func help(commands map[string]command.Command) command.Command {
	var cmds []command.Command
	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}
	return command.NewCommand("help", func(args ...string) ([]byte, error) {
		response := []string{"\n"}
		for _, cmd := range cmds {
			response = append(response, cmd.Name())
		}
		return []byte(strings.Join(response, "\n")), nil
	})
}
