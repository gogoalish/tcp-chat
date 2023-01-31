package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	HOST = "localhost"
	TYPE = "tcp"
)

var (
	nicknames  []string
	PORT       = "8080"
	clients    []net.Conn
	lastClient net.Conn
	strchan    = make(chan string)
	history    string
	mutex      sync.Mutex
)

func NetRun(PORT string) {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on the port :%v", PORT)
	defer listen.Close()
	go HandleConnections(listen)
	SendMessage()
}

func HandleConnections(listen net.Listener) {
	for {
		if len(clients) == 10 {
			continue
		}
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go WriteMessage(conn)
	}
}

func WriteMessage(conn net.Conn) {
	var message string
	nickname := GetNickname(conn)
	joinMessage := JoinChat(nickname)
	clients = append(clients, conn)
	if history != "" {
		fmt.Fprintln(conn, history[1:])
	}
	history += joinMessage
	defer func() {
		RemoveClient(conn, nickname)
		LeftChat(nickname)
	}()
	for {
		time := (time.Now())
		prefix := fmt.Sprintf("[%0.19v][%v]:", time, nickname)
		fmt.Fprint(conn, prefix)
		bufString, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		if !IsValidMsg(bufString) {
			continue
		}
		message = fmt.Sprintf("[%v]:%v", nickname, bufString)
		lastClient = conn
		strchan <- message
	}
}

func SendMessage() {
	for message := range strchan {
		time := time.Now()
		formatedM := fmt.Sprintf("\n[%0.19v]%v", time, message[:len(message)-1])
		for i := range clients {
			if clients[i] != lastClient {
				fmt.Fprint(clients[i], formatedM)
			}
		}
		history += formatedM
	}
}

func GetNickname(conn net.Conn) string {
	fmt.Fprint(conn, "Enter nickname: ")
	bufString, _ := bufio.NewReader(conn).ReadString('\n')
	for !IsValidName(bufString[:len(bufString)-1]) {
		fmt.Fprint(conn, "imya huevoe: ")
		bufString, _ = bufio.NewReader(conn).ReadString('\n')
	}
	for !IsUniqueName(bufString[:len(bufString)-1]) {
		fmt.Fprint(conn, "imya est uzhe: ")
		bufString, _ = bufio.NewReader(conn).ReadString('\n')
	}
	nicknames = append(nicknames, bufString[:len(bufString)-1])
	return bufString[:len(bufString)-1]
}

func IsValidName(bufString string) bool {
	for _, x := range strings.ToLower(bufString) {
		if x < 'a' || x > 'z' {
			return false
		}
	}
	return true
}

func IsUniqueName(bufString string) bool {
	for i := range nicknames {
		if nicknames[i] == bufString {
			return false
		}
	}
	return true
}

func RemoveClient(conn net.Conn, nickname string) {
	for i, client := range clients {
		if client == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
	for i := range nicknames {
		if nicknames[i] == nickname {
			nicknames = append(nicknames[:i], nicknames[i+1:]...)
			break
		}
	}
}

func LeftChat(nickname string) {
	message := fmt.Sprintf("\n%v has left our chat...", nickname)
	for i := range clients {
		fmt.Fprint(clients[i], message)
	}
	history += message
}

func JoinChat(nickname string) string {
	message := fmt.Sprintf("\n%v has joined our chat...", nickname)
	for i := range clients {
		fmt.Fprint(clients[i], message)
	}
	return message
}

func IsValidMsg(message string) bool {
	for _, x := range message {
		if x >= 32 {
			return true
		}
	}
	return false
}
