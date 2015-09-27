package lib

import ()

type Channel struct {
	clients []*Client
	w       *World
	name    string
}

func (ch *Channel) Send(c *Client, text string) {
	for i := 0; i < len(ch.clients); i++ {
		if ch.clients[i] != c {
			ch.clients[i].SendPrivMsg(c, ch.name, text)
		}
	}
}
