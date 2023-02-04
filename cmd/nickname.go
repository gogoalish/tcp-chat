package cmd

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func GetNickname(conn net.Conn) string {
	bufString, _ := bufio.NewReader(conn).ReadString('\n')
	CheckLength(&bufString, conn)
	for !IsValidName(bufString[:len(bufString)-1]) {
		fmt.Fprint(conn, "Name has not-valid characters. Allowed[0-9a-zA-Z]: ")
		bufString, _ = bufio.NewReader(conn).ReadString('\n')
		CheckLength(&bufString, conn)
	}
	for !IsUniqueName(bufString[:len(bufString)-1]) {
		fmt.Fprint(conn, "Name already exists: ")
		bufString, _ = bufio.NewReader(conn).ReadString('\n')
		CheckLength(&bufString, conn)
	}
	return bufString[:len(bufString)-1]
}

func IsValidName(bufString string) bool {
	if len(bufString) == 0 {
		return false
	}
	for _, x := range strings.ToLower(bufString) {
		if (x < 'a' || x > 'z') && (x < '0' || x > '9') {
			return false
		}
	}
	return true
}

func IsUniqueName(bufString string) bool {
	for nickname := range clients {
		if nickname == bufString {
			return false
		}
	}
	return true
}

func CheckLength(bufStringP *string, conn net.Conn) {
	bufString := *bufStringP
	for !(len(bufString) <= 21) || len(bufString) == 0 {
		fmt.Fprint(conn, "Length of name must contain at least 1 character and max 20 letters: ")
		bufString, _ = bufio.NewReader(conn).ReadString('\n')
	}
	*bufStringP = bufString
}
