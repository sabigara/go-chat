package main

import (
	"flag"

	"github.com/sabigara/go-chat/client"
	"github.com/sabigara/go-chat/server"
)

func main() {
	flag.Parse()
	args := flag.Args()
	cmd := args[0]
	var addr string
	if len(args) > 1 {
		addr = args[1]
	} else {
		addr = "localhost:9999"
	}

	switch cmd {
	case "server":
		server.Serve(addr)
	case "client":
		client.Run(addr)
	}
}
