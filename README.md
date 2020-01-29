# go-chat

## What's this

Tiny chat application built in golang on top of TCP.

Multiple users can join a room, and share messages.

## How to use

### Install

```bash
go get github.com/sabigara/go-chat
```

### Start server

```bash
go-chat server 127.0.0.1:1234
# addr & port arg is optional. default is localhost:9999
```

Make sure $GOPATH/bin is in your $PATH.

### Start client

Execute following command on multiple terminal windows.

Address and port must be the same as server.

```bash
go-chat client 127.0.0.1:1234
# addr & port arg is optional. default is localhost:9999
```

### Have a chat

Now you are ready for a chat!

But introduce yourself first. Your name is considered identifier, so provide different ones for each windows.
