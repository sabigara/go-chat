package client

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/sabigara/go-chat/common"
	"github.com/sabigara/go-chat/console"
)

type clientConn struct {
	net.Conn
	User string
}

func (conn *clientConn) writefStr(format string, args ...interface{}) error {
	_, err := fmt.Fprintf(conn, format+string(common.DELIMITER), args...)
	return err
}

func scan(userCh chan string) {
	for {
		userCh <- console.Scan()
	}
}

func handleInput(userIn, servIn chan string, conn *clientConn) {
	for {
		select {
		case in := <-userIn:
			console.CursorUp(1)
			console.EraseLine()
			conn.writefStr(in)
		case in := <-servIn:
			console.EraseLine()
			console.Writeln(in)
			console.Writef("%s: ", conn.User)
		}
	}
}

func listenServer(servIn chan string, conn *clientConn) {
	var mesBytes []byte
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("client: read error. %s\n", err.Error())
			return
		}
		mesBytes = append(mesBytes, buf[:n]...)
		if common.EndsWithDelimiter(mesBytes) {
			servIn <- string(common.TrimBytes(mesBytes))
			mesBytes = nil
		}
	}
}

func chat(conn *clientConn) {
	userCh := make(chan string)
	servCh := make(chan string)
	go handleInput(userCh, servCh, conn)
	go listenServer(servCh, conn)
	scan(userCh)
}

func Run(addr string) {
	sc := bufio.NewScanner(os.Stdin)
	console.SetScanner(sc)
	co, err := net.Dial("tcp", addr)
	if err != nil {
		console.Writef("client: dial error. %s", err.Error())
		return
	}
	defer co.Close()
	name := console.Prompt("Who are you?\n")
	conn := clientConn{Conn: co, User: name}
	console.Writef("Welcome, %s!\n", name)
	console.Writef("%s: ", name)
	conn.writefStr("%s", name)

	chat(&conn)
}
