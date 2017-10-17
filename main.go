package main

/*
TODO:
* Need to issue each session a random id
* Map the ids to channels and a way to iterate through all ids
to publish messages
*/

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	MAX_CLIENTS = 128
	PORT        = ":4000"
	PROTOCOL    = "tcp"
)

var userSet map[string]bool
var clients []net.Conn

func main() {
	userSet = make(map[string]bool)
	clients = make([]net.Conn, 0, MAX_CLIENTS)

	listener, err := net.Listen(PROTOCOL, PORT)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		if len(clients) == MAX_CLIENTS {
			conn.Write([]byte("Chat Room Full. Try again later!\n"))
			conn.Close()
			continue
		}

		clients = append(clients, conn)

		fmt.Printf("Connect: Local %s -> Remote %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	sess := &Session{
		Registered: false,
		UserName:   "Unknown",
	}

	conn.Write([]byte("Connected!\n"))
	conn.Write([]byte("Enter Username:\n"))

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		result := issueCommand(strings.TrimSpace(msg), conn, sess)
		if result {
			return
		}
	}
}

// Return true if the server should disconnect else false
func issueCommand(cmd string, conn net.Conn, sess *Session) bool {
	if cmd == "exit" {
		removeClient(conn)
		fmt.Printf("Disconnect: Local %s -> Remote %s\n", conn.LocalAddr().String(), conn.RemoteAddr().String())
		return true
	}
	if !sess.Registered {
		sess.Register(cmd, conn)
		return false
	}

	broadcastMsg(sess.UserName, cmd)

	return false
}

func broadcastMsg(sender, cmd string) {
	for i := 0; i < len(clients); i++ {
		clients[i].Write([]byte(sender))
		clients[i].Write([]byte(": "))

		clients[i].Write([]byte(cmd))
		clients[i].Write([]byte("\n"))
		clients[i].Write([]byte("\n"))
	}
}

// TODO: use a better algorithm
func removeClient(conn net.Conn) {
	for i := 0; i < len(clients); i++ {
		if clients[i] == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}
