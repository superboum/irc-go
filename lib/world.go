package lib

import ()

type World struct {
	clients  []*Client
	channels map[string]*Channel
}

func (w *World) Init() {
	w.channels = make(map[string]*Channel)
}

func (w *World) RegisterClient(c *Client) {
	c.w = w
	w.clients = append(w.clients, c)
}

func (w *World) JoinChannel(c *Client, name string) {
	if w.channels[name] == nil {
		ch := new(Channel)
		ch.name = name
		ch.w = w
		w.channels[name] = ch
	}
	w.channels[name].clients = append(w.channels[name].clients, c)
}

func (w *World) SendToChannel(c *Client, name string, text string) {
	w.channels[name].Send(c, text)
}
