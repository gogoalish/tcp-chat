package cmd

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

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
