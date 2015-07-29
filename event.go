package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"
)

func (irc *Connection) AddCallback(eventcode string, callback func(*Event)) string {
	eventcode = strings.ToUpper(eventcode)

	if _, ok := irc.events[eventcode]; !ok {
		irc.events[eventcode] = make(map[string]func(*Event))
	}
	h := sha1.New()
	rawId := []byte(fmt.Sprintf("%v%d", reflect.ValueOf(callback).Pointer(), rand.Int63()))
	h.Write(rawId)
	id := fmt.Sprintf("%x", h.Sum(nil))
	irc.events[eventcode][id] = callback
	return id
}

func (irc *Connection) RunCallbacks(event *Event) {
	msg := event.Message
	log.Println(msg.String())
	if callbacks, ok := irc.events[event.Code]; ok {

		for _, callback := range callbacks {
			go callback(event)
		}
	}
}
