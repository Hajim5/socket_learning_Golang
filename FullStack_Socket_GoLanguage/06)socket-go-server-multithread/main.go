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

func main() {
    fmt.Println("> multithread echo-server is activated")

    ln, err := net.Listen("tcp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    defer ln.Close()

    //Infinite loop for accepting multiple clients
    // == class ThreadedTCPServer()
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("> accept() failed:", err)
            continue
        }

        go handleClient(conn)
    }
}
// Handles one client session
func handleClient(conn net.Conn) {
    defer conn.Close()

    addr := conn.RemoteAddr().(*net.TCPAddr)
    fmt.Printf("> client connected by IP address %s with Port number %d\n",
        addr.IP.String(),
        addr.Port)

    reader := bufio.NewScanner(conn)

    for reader.Scan() {
        recv := reader.Text()
        fmt.Println("> echoed:", recv)

        _, err := conn.Write([]byte(recv + "\n"))
        if err != nil {
            fmt.Println("> send() failed:", err)
            return
        }

        if recv == "quit" {
            break
        }
    }

    if err := reader.Err(); err != nil {
        fmt.Println("> recv() failed:", err)
    }

    fmt.Printf("> client disconnected: %s:%d\n", addr.IP.String(), addr.Port)
}
