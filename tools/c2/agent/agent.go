package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

const secret = "EnCryp!e0?"

type Agent struct {
	isAuth bool
	id     string
	secret string
	name   string
	ws     *websocket.Conn
}
type Payload struct {
	mode    string
	args    []string
	length  int
	content string
	from    string
	status  string
}

func NewAgent() *Agent {
	return &Agent{
		isAuth: false,
		id:     "",
		secret: secret,
		name:   "agent-" + RandomString(5),
	}
}

func NewPayload() *Payload {
	return &Payload{
		mode:    "",
		args:    []string{},
		length:  0,
		content: "",
		from:    "server",
		status:  "",
	}
}
func (a *Agent) connect(url string, origin string) {

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		failCount := 1
		log.Println("Connection failed, re-trying ...", failCount)
		for {
			ws, err = websocket.Dial(url, "", origin)
			if err != nil {
				failCount++
				time.Sleep(5 * time.Second)
				continue
			} else {
				break
			}
		}
	}
	a.ws = ws
	log.Printf("Connected to %s", url)
	if a.isAuth != true {
		a.authReq()
	}

	a.listenLoop()
}

func (a *Agent) listenLoop() {
	var buf = make([]byte, 1024)
	for {
		n, err := a.ws.Read(buf)

		if err != nil {
			log.Println("Error reading", err)
		}
		msg := buf[:n]

		a.msgHandler(string(msg))

		time.Sleep(1 * time.Second)
	}
}

func (a *Agent) msgHandler(msg string) {
	fmt.Println(string(msg))
	payload := NewPayload()
	payload.parser(msg)
	fmt.Println(*payload)
	if a.isAuth != true && payload.mode != "authResponse" {
		fmt.Println("----------------here---------------\n", a.isAuth, payload.mode)
		a.authReq()
	}
	switch payload.mode {
	case "cmd":
		a.handleCmd(payload)
	case "authResponse":
		a.handleAuthRes(payload)
	case "msg":
		a.handleMsg(payload)
	case "ping":
		a.handlePing(payload)
	}

}

func (a *Agent) authReq() {
	authCmd := "/auth:" + a.name + "@" + a.secret
	send(authCmd, a.ws)
}

func (a *Agent) handleAuthRes(payload *Payload) {
	if payload.mode != "authResponse" {
		return
	}
	switch payload.status {
	case "E001":
		fmt.Println("Retrying ...")
		time.Sleep(3 * time.Second)
		a.authReq()
	case "E403":
		a.authReq()
	case "S001":
		data := strings.Split(payload.content, "/")
		a.isAuth = true
		a.id = data[0]
		fmt.Println("+++++++++++++auth+++++++++++++\n", a.isAuth, payload.mode)
	}
}

func (a *Agent) handlePing(payload *Payload) {
	if payload.mode != "ping" {
		return
	}
	send("/pong:"+payload.from, a.ws)
}

func (a *Agent) handleCmd(payload *Payload) {
	if payload.mode != "cmd" {
		return
	}
	send("/msg: well recived -> "+payload.from, a.ws)
}

func (a *Agent) handleMsg(payload *Payload) {
	if payload.mode != "ping" {
		return
	}
	send("/msg: well recived -> "+payload.from, a.ws)
}

func main() {
	origin := "http://localhost/"
	url := "ws://localhost:3200/ws"
	agent := NewAgent()
	agent.connect(url, origin)
}

func send(str string, ws *websocket.Conn) bool {
	_, err := ws.Write([]byte(str))
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	return true
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (p *Payload) parser(str string) {
	args := strings.Split(str, " -> ")
	switch args[0] {
	case "[cmd]":
		p.mode = "cmd"
		p.cmdParser(args[1])
	case "[auth]":
		p.mode = "authResponse"
		p.authResParser(args[1])
	case "[ping]":
		p.mode = "ping"
		p.pingParser(args[1])
	default:
		p.mode = "msg"
		p.msgParser(str)
	}
}

func (p *Payload) cmdParser(str string) {
	args := strings.Split(str, " <- ")
	p.args = args
	p.content = args[0]
	p.from = args[1]
	p.length = len(args)
}

func (p *Payload) authResParser(str string) {
	args := strings.Split(str, ": ")
	p.args = args
	p.status = args[0]
	p.content = args[1]
	p.length = len(args)
}

func (p *Payload) pingParser(str string) {
	p.from = str
	p.length = 1

}

func (p *Payload) msgParser(str string) {
	args := strings.Split(str, " <- ")
	p.args = args
	p.length = len(args)
	p.content = args[0]
	p.from = "*"
	if len(args) > 1 {
		p.from = args[1]
	}
}
