package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xuyu/goredis"
)

func main() {
	conn := startConnection()
	defer conn.socket.Close()
	client, _ := goredis.Dial(&goredis.DialConfig{Address: "127.0.0.1:6379"})

	conn.initializeCallbacks(client)

	channel := "bittermandel"

	conn.sendRawMessage("PASS oauth:avgzw6jychctbacx0kp1odifn1bzkr\r\n")
	conn.sendMessage("NICK bittermandel\r\n")
	conn.sendMessage(fmt.Sprintf("JOIN #%s\r\n", channel))
	log.Printf("Connected to IRC server\n")

	conn.addChatCommand("!test", "testing")

	conn.addTimedMessage(func() {
		conn.sendMessage("PRIVMSG #bittermandel :If you want more entertainment, please follow me here on Twitch and on Youtube. Content coming soon!")

	}, 5)

	for {
		line, err := conn.in.Decode()
		if err != nil {
			log.Fatal(err)
		}
		serialized, err := json.Marshal(line)
		client.LPush(fmt.Sprintf("messages:%s", channel), string(serialized))
		//log.Println(string(serialized))

		conn.RunCallbacks(&Event{
			Message:    line,
			Connection: conn,
			Code:       line.Command,
		})
	}
}
