package service

import (
	"errors"
	"sync"
	"github.com/nlopes/slack"
)

// Connector is structure of connector
type Connector struct {
	debug bool
	token string

	sync.Mutex
	running bool
	exit    chan bool

	api *slack.Client
}

// NewConnector is constructor for generating new connector
func NewConnector() Connector {
	return Connector{}
}

// Init is initializing new connector
func (c *Connector) Init(debug bool, token string) error {
	if len(token) == 0 {
		return errors.New("missing slack token")
	}

	c.debug = debug
	c.token = token

	return nil
}

// Stream func creates connection and auth
func (c *Connector) Stream() (*Connection, error) {
	c.Lock()
	defer c.Unlock()

	if !c.running {
		return nil, errors.New("not running")
	}

	// test auth
	auth, err := c.api.AuthTest()
	if err != nil {
		return nil, err
	}

	rtm := c.api.NewRTM()
	exit := make(chan bool)

	go rtm.ManageConnection()

	go func() {
		select {
		case <-c.exit:
			select {
			case <-exit:
				return
			default:
				close(exit)
			}
		case <-exit:
		}

		rtm.Disconnect()
	}()

	conn := &Connection{
		auth:  auth,
		rtm:   rtm,
		exit:  exit,
		names: make(map[string]string),
	}

	go conn.run()

	return conn, nil
}

// Start func opens slack connection
func (c *Connector) Start() error {
	if len(c.token) == 0 {
		return errors.New("missing slack token")
	}

	c.Lock()
	defer c.Unlock()

	if c.running {
		return nil
	}

	api := slack.New(c.token, slack.OptionDebug(c.debug))

	// test auth
	_, err := api.AuthTest()
	if err != nil {
		return err
	}

	c.api = api
	c.exit = make(chan bool)
	c.running = true

	return nil
}

// Stop func closes exit channels
func (c *Connector) Stop() error {
	c.Lock()
	defer c.Unlock()

	if !c.running {
		return nil
	}

	close(c.exit)
	c.running = false
	return nil
}
