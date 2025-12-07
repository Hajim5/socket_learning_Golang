package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
    "sync"
)

const (
    HOST = "127.0.0.1"
    PORT = "65456"
)

// Global list of clients (like group_queue in Python)
var (
    clients   = make(map[net.Conn]struct{})
    clientsMu sync.Mutex
)

func main() {
    fmt.Println("> multithread chat-server is activated")

    ln, err := net.Listen("tcp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    defer ln.Close()

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("> accept() failed:", err)
            continue
        }

        addClient(conn)
        go handleClient(conn)
    }
}

func addClient(conn net.Conn) {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    clients[conn] = struct{}{}

    addr := conn.RemoteAddr().(*net.TCPAddr)
    fmt.Printf("> client joined: %s:%d (total clients: %d)\n",
        addr.IP.String(), addr.Port, len(clients))
}

func removeClient(conn net.Conn) {
    clientsMu.Lock()
    defer clientsMu.Unlock()
    delete(clients, conn)

    addr := conn.RemoteAddr().(*net.TCPAddr)
    fmt.Printf("> client left: %s:%d (total clients: %d)\n",
        addr.IP.String(), addr.Port, len(clients))
}

// Broadcast message to all connected clients
func broadcast(sender net.Conn, msg string) {
    clientsMu.Lock()
    defer clientsMu.Unlock()

    for c := range clients {
        // Optionally skip sender if you don't want them to see their own message
        // if c == sender { continue }

        _, err := fmt.Fprintln(c, msg)
        if err != nil {
            fmt.Println("> send() failed to a client:", err)
        }
    }
}

func handleClient(conn net.Conn) {
    defer func() {
        removeClient(conn)
        conn.Close()
    }()

    addr := conn.RemoteAddr().(*net.TCPAddr)
    fmt.Printf("> client connected by IP address %s with Port number %d\n",
        addr.IP.String(), addr.Port)

    reader := bufio.NewScanner(conn)

    for reader.Scan() {
        text := reader.Text()
        text = strings.TrimRight(text, "\r\n")

        if text == "quit" {
            // Client wants to disconnect
            fmt.Printf("> client %s:%d requested quit\n", addr.IP.String(), addr.Port)
            return
        }

        fmt.Printf("> received from %s:%d: %s\n", addr.IP.String(), addr.Port, text)
        // Broadcast to all clients (chat)
        broadcast(conn, text)
    }

    if err := reader.Err(); err != nil {
        fmt.Println("> recv() failed by error:", err)
    }
}
