package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func LeaveChat(nickname string) {
	delete(clients, nickname)
	message := fmt.Sprintf("%v has left our chat...", nickname)
	for i := range clients {
		CleanAndPrint(i, message)
	}
	history += "\n" + message
}

func JoinChat(conn net.Conn, nickname string) {
	message := fmt.Sprintf("%v has joined our chat...", nickname)
	for i := range clients {
		if i != nickname {
			CleanAndPrint(i, message)
		}
	}
	clients[nickname] = conn
	if history != "" {
		fmt.Fprintln(conn, history[1:])
	}
	history += "\n" + message
}

func IsValidMsg(message string) bool {
	if message == "" {
		return false
	}
	for _, x := range message {
		if x < 32 || x > 126 {
			return false
		}
	}
	return true
}

func CleanAndPrint(name, message string) {
	prefix := fmt.Sprintf("[%0.19v][%v]:", time.Now(), name)
	fmt.Fprint(clients[name], "\r"+strings.Repeat(" ", len(prefix))+"\r")
	fmt.Fprintln(clients[name], message)
	PrintPrefix(name)
}

func PrintPrefix(name string) {
	time := (time.Now())
	prefix := fmt.Sprintf("[%0.19v][%v]:", time, name)
	fmt.Fprint(clients[name], prefix)
}

func PrintLogo(conn net.Conn) {
	data, err := os.ReadFile("./assets/logo.txt")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(data)
}
