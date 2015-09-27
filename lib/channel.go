package lib

import (
	"fmt"
)

type Channel struct {
	clients []*Client
	w       *World
	name    string
}

func (ch *Channel) Send(c *Client, text string) {
	fmt.Println(ch)
	for i := 0; i < len(ch.clients); i++ {
		ch.clients[i].SendPrivMsg(c, ch.name, text)
	}
}
