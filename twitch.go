package main

import (
	"errors"
	"strings"
	"time"
)

func (conn *Connection) addChatCommand(command string, response string) error {
	if !strings.HasPrefix(command, "!") {
		return errors.New("Wrong format")
	}
	if strings.HasPrefix(command, "/") {
		return errors.New("Unacceptable")
	}
	conn.AddCallback("PRIVMSG", func(e *Event) {
		if strings.HasPrefix(e.Message.Trailing, command) {
			conn.sendMessage("PRIVMSG " + strings.Join(e.Message.Params, " ") + " " + response)
		}
	})
	return nil
}

func (conn *Connection) addTimedMessage(task func(), delay time.Duration) error {
	ticker := time.NewTicker(delay * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				task()
			}
		}
	}()
	return nil
}
