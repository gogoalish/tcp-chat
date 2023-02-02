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
	currClient string
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
	nickname := GetNickname(conn)
	mutex.Lock()
	JoinChat(conn, nickname)
	mutex.Unlock()
	for {
		PrintPrefix(nickname)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		if !IsValidMsg(message) {
			continue
		}
		mutex.Lock()
		currClient = nickname
		mutex.Unlock()
		strchan <- message
	}
	mutex.Lock()
	LeaveChat(nickname)
	mutex.Unlock()
}

func SendMessage() {
	for message := range strchan {
		message = fmt.Sprintf("[%0.19v][%v]:%v", time.Now(), currClient, message)
		mutex.Lock()
		for name := range clients {
			if name != currClient {
				CleanAndPrint(name, message)
			}
		}
		history += "\n" + message
		mutex.Unlock()
	}
}
