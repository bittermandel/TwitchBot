package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"

	"github.com/xuyu/goredis"
)

func (irc *Connection) addCommand(commandcode string, callback func(*Event)) string {
	commandcode = strings.ToUpper(commandcode)

	if _, ok := irc.commands[commandcode]; !ok {
		irc.commands[commandcode] = make(map[string]func(*Event))
	}
	h := sha1.New()
	rawId := []byte(fmt.Sprintf("%v%d", reflect.ValueOf(callback).Pointer(), rand.Int63()))
	h.Write(rawId)
	id := fmt.Sprintf("%x", h.Sum(nil))
	irc.commands[commandcode][id] = callback
	return id
}

func (irc *Connection) initializeCallbacks(client *goredis.Redis) {
	irc.events = make(map[string]map[string]func(*Event))

	irc.AddCallback("PING", func(e *Event) {
		irc.sendMessage("PONG :" + e.Message.Trailing)
		log.Println("PONG :" + e.Message.Trailing)
	})

	irc.AddCallback("PRIVMSG", func(e *Event) {
		_, err := client.Incr("messages")
		if err != nil {
			panic(err)
		}
	})
}
