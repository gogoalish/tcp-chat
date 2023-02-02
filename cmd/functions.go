package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func LeaveChat(nickname string) {
	message := fmt.Sprintf("\n%v has left our chat...", nickname)
	mutex.Lock()
	for i := range clients {
		fmt.Fprint(clients[i], message)
	}
	history += message
	mutex.Unlock()
}

func JoinChat(nickname string) string {
	message := fmt.Sprintf("\n%v has joined our chat...", nickname)
	for i := range clients {
		fmt.Fprint(clients[i], message)
	}
	return message
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

func PrintLogo(conn net.Conn) {
	data, err := os.ReadFile("./assets/logo.txt")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(data)
}

func Clean(line string) string {
	return "\r" + strings.Repeat(" ", len(line)) + "\r"
}
