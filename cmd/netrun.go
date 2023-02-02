package cmd

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	HOST = "localhost"
	TYPE = "tcp"
)

var (
	PORT       = "8080"
	clients    = make(map[string]net.Conn)
	currClient net.Conn
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
	PrintLogo(conn)
	var message string
	nickname := GetNickname(conn)
	joinMessage := JoinChat(nickname)
	mutex.Lock()
	clients[nickname] = conn
	if history != "" {
		fmt.Fprintln(conn, history[1:])
	}
	history += joinMessage
	mutex.Unlock()
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
		currClient = conn
		mutex.Unlock()
		strchan <- message
	}
	delete(clients, nickname)
	LeaveChat(nickname)
}

func SendMessage() {
	for message := range strchan {
		time := time.Now()
		formatedM := fmt.Sprintf("\n[%0.19v]%v", time, message[:len(message)-1])
		mutex.Lock()
		for i := range clients {
			if clients[i] != currClient {
				fmt.Fprint(clients[i], formatedM)
			}
		}
		history += formatedM
		mutex.Unlock()
	}
}
