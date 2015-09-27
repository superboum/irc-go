package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "6667"
	CONN_TYPE = "tcp"
	SERV_NAME = ":irc.deuxfleurs.fr"
)

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

type Channel struct {
	clients []*Client
	w       *World
	name    string
}

func (ch *Channel) Send(c *Client, text string) {
	for i := 0; i < len(ch.clients); i++ {
		ch.clients[i].SendPrivMsg(c, ch.name, text)
	}
}

type Client struct {
	w        *World
	co       net.Conn
	nickname string
}

func (c *Client) HandleRequest() {
	buf := make([]byte, 1024)
	c.HelloMessage()

	for {
		ln, err := c.co.Read(buf)
		msg := string(buf[0:ln])
		if err != nil {
			fmt.Println("Error reading:", err.Error(), " - ", c.co.RemoteAddr())
			break
		}

		re := regexp.MustCompile("[\r\n]*([A-Z]+) ([^\n\r]+)")
		msg_dec := re.FindAllStringSubmatch(msg, -1)

		for i := 0; i < len(msg_dec); i++ {
			cmd := msg_dec[i][1]
			prop := msg_dec[i][2]

			fmt.Println("cmd:", cmd, "prop:", prop)
			switch cmd {
			case "NICK":
				c.SetNickname(prop)
			case "USER":
				c.SetUser(prop)
			case "JOIN":
				c.Join(prop)
			case "PING":
				c.Pong(prop)
			case "PRIVMSG":
				c.PrivMsg(prop)
			case "QUIT":
				c.Close(prop)
			default:
				fmt.Println("Error, unrecognized: " + msg)
			}
		}
		//c.co.Write([]byte(":qdufour MODE qdufour +i\n"))
	}
}

func (c *Client) HelloMessage() {
	c.co.Write([]byte(SERV_NAME + " NOTICE *  :*** HELLO WORLD \r\n"))
}

func (c *Client) SetNickname(prop string) {
	c.nickname = prop
	c.co.Write([]byte(SERV_NAME + " 001 " + c.nickname + " :Welcome to this test IRC go server " + c.nickname + "\r\n"))
}

func (c *Client) SetUser(prop string) {
	fmt.Println("Not yet implemented set user: ", prop)
}

func (c *Client) Join(prop string) {
	c.w.JoinChannel(c, prop)
	c.co.Write([]byte(":" + c.nickname + "!hostname JOIN " + prop + "\r\n"))
}

func (c *Client) PrivMsg(prop string) {
	fmt.Println(prop)
	re := regexp.MustCompile("(#[a-zA-Z0-9]+) :(.*)")
	prop_dec := re.FindAllStringSubmatch(prop, -1)
	fmt.Println(prop_dec)
	fmt.Println(prop_dec[0])

	c.w.SendToChannel(c, prop_dec[0][0], prop_dec[0][1])
}

func (c *Client) Pong(prop string) {
	c.co.Write([]byte("PONG " + prop + "\r\n"))
}

func (c *Client) Close(prop string) {
	c.co.Close()
}

func (c *Client) SendPrivMsg(from *Client, channel string, text string) {
	c.co.Write([]byte(":" + from.nickname + "!hostname PRIVMSG " + channel + " :" + text + "\r\n"))
}

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	w := new(World)
	w.Init()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		client := new(Client)
		client.co = conn
		w.RegisterClient(client)
		go client.HandleRequest()
	}
}
