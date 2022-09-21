package main

import (
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
)

type XmppConfig struct {
	Enabled  bool
	Server   string
	Jid      string
	Password string
	Room     string
	TLS      bool
	Nickname string
}

type HttpConfig struct {
	Listen_Address string
}

type AlertConfig struct {
	Template string
}

type Config struct {
	Xmpp     XmppConfig
	Http     HttpConfig
	Alerting AlertConfig
}

var GlobalConfig Config

func ParseConfig(path string) error {
	var cfg Config

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(bytes, &cfg)
	if err != nil {
		return err
	}

	GlobalConfig = cfg
	return nil
}
