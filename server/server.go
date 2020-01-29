package server

import (
	"chat/common"
	"fmt"
	"net"
)

type message struct {
	Data []byte
	User string
}

type serverConn struct {
	net.Conn
	User string
}

func (conn *serverConn) writeStrf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(conn.Conn, format+string(common.DELIMITER), args...)
}

type connectionPool map[string]*serverConn

func (connPool *connectionPool) publish(mes string) {
	for _, conn := range *connPool {
		conn.writeStrf("%s", mes)
	}
}

func serveRoom(in <-chan message, newConn <-chan *serverConn) {
	connPool := make(connectionPool)
	for {
		select {
		case mes := <-in:
			connPool.publish(fmt.Sprintf("%s: %s", mes.User, string(mes.Data)))
		case conn := <-newConn:
			connPool[conn.User] = conn
			connPool.publish(fmt.Sprintf("%s joined!", conn.User))
		}
	}
}

func handle(conn *serverConn, connCh chan *serverConn, room chan message) {
	var mesBytes []byte
	var mes message

	buf := make([]byte, 32)
	// First read will be the user's name
	_, err := conn.Read(buf)
	if err != nil {
		return
	}
	buf = common.TrimBytes(buf)
	user := string(buf)

	mes.User = user
	conn.User = user
	connCh <- conn

	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Read error: %s", err.Error())
			return
		}
		mesBytes = append(mesBytes, buf[:n]...)
		if common.EndsWithDelimiter(mesBytes) {
			mes.Data = common.TrimBytes(mesBytes)
			room <- mes
			mesBytes = nil
		}
	}
}

func Serve(addr string) {
	listner, err := net.Listen("tcp", addr)
	defer listner.Close()
	if err != nil {
		panic("Cannot listen")
	}
	room := make(chan message)
	connCh := make(chan *serverConn)
	go serveRoom(room, connCh)
	for {
		co, err := listner.Accept()
		if err != nil {
			fmt.Println("Cannot establish connection")
			continue
		}
		defer co.Close()
		conn := serverConn{Conn: co, User: ""}
		go handle(&conn, connCh, room)
	}
}
