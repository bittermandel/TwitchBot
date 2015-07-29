package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/sorcix/irc"
)

type Connection struct {
	socket   *irc.Conn
	username string
	oauthkey string

	out *irc.Encoder
	in  *irc.Decoder

	lastMessage time.Time
	events      map[string]map[string]func(*Event)
	commands    map[string]map[string]func(*Event)
	Log         *log.Logger

	active bool
}

type Event struct {
	Code       string
	Message    *irc.Message
	Connection *Connection
}

func (conn *Connection) sendMessage(rawMessage string) {
	message := irc.ParseMessage(rawMessage)
	conn.out.Encode(message)
}

func (conn *Connection) sendRawMessage(rawMessage string) {
	conn.out.Write([]byte(rawMessage))
}

func getCredentials() (username, password string) {
	dat, err := ioutil.ReadFile("./credentials")
	if err != nil {
		return "", ""
	}
	creds := strings.Split(string(dat), "\n")
	return creds[0], creds[1]

}

func startConnection() *Connection {
	c, _ := irc.Dial("irc.twitch.tv:6667")
	username, oauthkey := getCredentials()
	return &Connection{
		username: username,
		oauthkey: oauthkey,
		socket:   c,
		in:       &c.Decoder,
		out:      &c.Encoder,
	}
}
