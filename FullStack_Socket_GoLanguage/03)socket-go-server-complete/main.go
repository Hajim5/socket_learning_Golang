package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	HOST = "127.0.0.1"
	PORT = "65456"
)

//using runServer() function 

func main() {
	//Function for activating and deactivating the server
	fmt.Println("> echo-server is activated")
	runServer()
	fmt.Println("> echo-server is de-activated")
}

func runServer() {
	// Create TCP listener 
	ln, err := net.Listen("tcp", HOST+":"+PORT)
	if err != nil {
		// Equivalent to: bind() or listen() failed
		fmt.Println("> bind() or listen() failed and program terminated")
		fmt.Println("> error:", err)
		return
	}
	defer ln.Close()

	// Accept one client
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("> accept() failed and program terminated")
		fmt.Println("> error:", err)
		return
	}
	defer conn.Close()

	// Client info
	//Access IP and Port separately
	addr := conn.RemoteAddr().(*net.TCPAddr)
	fmt.Printf("> client connected by IP address %s with Port number %d\n",
		addr.IP.String(), addr.Port)

	reader := bufio.NewScanner(conn)

	for reader.Scan() {
		recv := reader.Text()
		fmt.Println("> echoed:", recv)

		// Echo back
		_, err := conn.Write([]byte(recv + "\n"))
		if err != nil {
			log.Println("> send() failed by error:", err)
			break
		}

		if recv == "quit" {
			break
		}
	}

	if err := reader.Err(); err != nil {
		fmt.Println("> recv() failed by error:", err)
	}
}
