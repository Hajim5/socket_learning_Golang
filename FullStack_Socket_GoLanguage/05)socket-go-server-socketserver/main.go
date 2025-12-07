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
    fmt.Println("> echo-server (socketserver style) is activated")

    ln, err := net.Listen("tcp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    defer ln.Close()

    conn, err := ln.Accept()
    if err != nil {
        log.Fatalf("Failed to accept: %v", err)
    }
    defer conn.Close()

    //separates connection handling from server stup
    // = python's socketserver handler pattern
    //allow better structure for multithreading or chat or etc
    handleConnection(conn)

    fmt.Println("> echo-server is de-activated")
}
// extracting IP and port using type assertion
//gives IP and Port separately
func handleConnection(conn net.Conn) {
    addr := conn.RemoteAddr().(*net.TCPAddr)

    fmt.Printf(
        "> client connected by IP address %s with Port number %d\n",
        addr.IP.String(),
        addr.Port,
    )

    reader := bufio.NewScanner(conn)

    for reader.Scan() {
        recv := reader.Text()
        fmt.Println("> echoed:", recv)

        _, err := conn.Write([]byte(recv + "\n"))
        if err != nil {
            log.Println("> send failed:", err)
            return
        }

        if recv == "quit" {
            return
        }
    }

    if err := reader.Err(); err != nil {
        fmt.Println("> recv() failed by error:", err)
    }
}
