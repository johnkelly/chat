package main

import (
	"net"
)

type Session struct {
	Registered bool
	UserName   string
}

func (sess *Session) Register(s string, conn net.Conn) {
	if !validUserName(s) {
		conn.Write([]byte("Invalid Username. Try again:\n"))
		return
	}
	sess.UserName = s
	sess.Registered = true

	userSet[sess.UserName] = true
	conn.Write([]byte("Registered. Start typing to send your messages\n"))
}

func validUserName(s string) bool {
	if _, ok := userSet[s]; ok {
		return false
	}
	return true
}
