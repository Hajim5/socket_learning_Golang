package main

import (
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

var (
    clients   = make(map[string]*net.UDPAddr) 
    clientsMu sync.Mutex
)

func main() {
    fmt.Println("> udp-chat-server is activated")

    // Bind UDP address
    addr, err := net.ResolveUDPAddr("udp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("> ResolveUDPAddr() failed: %v", err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatalf("> ListenUDP() failed: %v", err)
    }
    defer conn.Close()

    buf := make([]byte, 2048)

    for {
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("> recvfrom() failed:", err)
            continue
        }

        msg := strings.TrimSpace(string(buf[:n]))
        handleMessage(conn, msg, clientAddr)
    }
}

func handleMessage(conn *net.UDPConn, msg string, client *net.UDPAddr) {
    key := client.String()

    switch {
    case msg == "#REG":
        registerClient(client)
    case msg == "#DEREG" || msg == "quit":
        deregisterClient(client)
    default:
        // normal chat message
        clientsMu.Lock()
        count := len(clients)
        _, registered := clients[key]
        clientsMu.Unlock()

        if count == 0 {
            fmt.Println("> no clients to echo")
            return
        }
        if !registered {
            fmt.Printf("> unregistered client %s tried to send message: %s (ignored)\n",
                key, msg)
            return
        }

        fmt.Printf("> from %s: %s\n", key, msg)
        broadcast(conn, msg)
    }
}

func registerClient(client *net.UDPAddr) {
    clientsMu.Lock()
    defer clientsMu.Unlock()

    key := client.String()
    if _, exists := clients[key]; !exists {
        clients[key] = client
        fmt.Printf("> client registered: %s (total %d)\n", key, len(clients))
    } else {
        fmt.Printf("> client already registered: %s\n", key)
    }
}

func deregisterClient(client *net.UDPAddr) {
    clientsMu.Lock()
    defer clientsMu.Unlock()

    key := client.String()
    if _, exists := clients[key]; exists {
        delete(clients, key)
        fmt.Printf("> client deregistered: %s (total %d)\n", key, len(clients))
    } else {
        fmt.Printf("> client not found to deregister: %s\n", key)
    }
}

func broadcast(conn *net.UDPConn, msg string) {
    clientsMu.Lock()
    defer clientsMu.Unlock()

    for key, c := range clients {
        _, err := conn.WriteToUDP([]byte(msg), c)
        if err != nil {
            fmt.Printf("> sendto() failed to %s: %v\n", key, err)
        }
    }
}
