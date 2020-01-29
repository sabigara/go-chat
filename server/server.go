package server

import (
	"fmt"
	"net"

	"github.com/sabigara/go-chat/common"
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
			// Add connection to pool.
			// If the key that's already in pool sent, treat it as the signal of leave.
			if _, ok := connPool[conn.User]; ok {
				delete(connPool, conn.User)
				connPool.publish(fmt.Sprintf("%s left.", conn.User))
			} else {
				connPool[conn.User] = conn
				connPool.publish(fmt.Sprintf("%s joined!", conn.User))
			}
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
			fmt.Printf("server: read error %s\n", err.Error())
			connCh <- conn
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
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("server: cannot listen. %s\n", err.Error())
		return
	}
	defer listener.Close()
	room := make(chan message)
	connCh := make(chan *serverConn)
	go serveRoom(room, connCh)
	for {
		co, err := listener.Accept()
		if err != nil {
			fmt.Printf("server: cannot establish connection. %s\n", err.Error())
			continue
		}
		defer co.Close()
		conn := serverConn{Conn: co, User: ""}
		go handle(&conn, connCh, room)
	}
}
