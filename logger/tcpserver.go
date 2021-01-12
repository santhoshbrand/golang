package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
)

// ChangeLogLevelRequest defines the request type to change the log levels
type ChangeLogLevelRequest struct {
	NewLogLevel uint8 // New log level
}

func main() {
	var req ChangeLogLevelRequest

	arglen := len(os.Args)
	fmt.Println("Command line args : ", arglen)

	if arglen < 2 {
		fmt.Println("ERROR: Insufficient arguements.. Input loglevel missing")
		return
	}
	loglevel, err := strconv.Atoi(os.Args[arglen-1])
	fmt.Println("Input log level : ", loglevel)

	servAddr := "localhost:3333"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	req.NewLogLevel = uint8(loglevel)
	ec := gob.NewEncoder(conn)
	err = ec.Encode(&req)

	if err != nil {
		println("Write to server failed:", err.Error())
		conn.Close()
		os.Exit(1)
	}

	fmt.Println("written to server = ", req)
	conn.Close()
}
