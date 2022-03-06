package infra

import (
	"github.com/m-mizutani/ghnotify/pkg/infra/notify"
	"github.com/m-mizutani/opac"
)

type Clients struct {
	slackClient notify.SlackClient
	opac        opac.Client
}

func (x *Clients) Slack() notify.SlackClient {
	if x.slackClient == nil {
		panic("slackClient is not configured")
	}
	return x.slackClient
}

func (x *Clients) OPAC() opac.Client {
	if x.opac == nil {
		panic("opac is not configured")
	}
	return x.opac
}

type Option func(c *Clients)

func New(options ...Option) *Clients {
	clients := &Clients{}
	for _, opt := range options {
		opt(clients)
	}
	return clients
}

func WithSlack(client notify.SlackClient) Option {
	return func(c *Clients) {
		c.slackClient = client
	}
}

func WithOPAC(client opac.Client) Option {
	return func(c *Clients) {
		c.opac = client
	}
}
