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
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		mutex.Lock()
		if len(clients) == 10 {
			fmt.Fprint(conn, "chat is full. Try again later.\n")
			conn.Close()
		}
		mutex.Unlock()
		go WriteMessage(conn)
	}
}

func WriteMessage(conn net.Conn) {
	var message string
	nickname := GetNickname(conn)
	joinMessage := JoinChat(nickname)
	mutex.Lock()
	clients = append(clients, conn)
	if history != "" {
		fmt.Fprintln(conn, history[1:])
	}
	history += joinMessage
	mutex.Unlock()
	defer func() {
		RemoveClient(conn, nickname)
		LeaveChat(nickname)
	}()
	for {
		time := (time.Now())
		prefix := fmt.Sprintf("[%0.19v][%v]:", time, nickname)
		fmt.Fprint(conn, prefix)
		bufString, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		bufString = strings.TrimSpace(bufString)
		if !IsValidMsg(bufString) {
			continue
		}
		message = fmt.Sprintf("[%v]:%v", nickname, bufString)
		mutex.Lock()
		lastClient = conn
		mutex.Unlock()
		strchan <- message
	}
}

func SendMessage() {
	for message := range strchan {
		time := time.Now()
		formatedM := fmt.Sprintf("\n[%0.19v]%v", time, message[:len(message)-1])
		mutex.Lock()
		for i := range clients {
			if clients[i] != lastClient {
				fmt.Fprint(clients[i], formatedM)
			}
		}
		history += formatedM
		mutex.Unlock()
	}
}

func RemoveClient(conn net.Conn, nickname string) {
	mutex.Lock()
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
	mutex.Unlock()
}

func LeaveChat(nickname string) {
	message := fmt.Sprintf("\n%v has left our chat...", nickname)
	for i := range clients {
		fmt.Fprint(clients[i], message)
	}
	mutex.Lock()
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
