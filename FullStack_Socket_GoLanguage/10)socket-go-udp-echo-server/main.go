package main

import (
    "fmt"
    "log"
    "net"
    "strings"
)

const (
    HOST = "127.0.0.1"
    PORT = "65456"
)

func main() {
    fmt.Println("> echo-server is activated")

    addr, err := net.ResolveUDPAddr("udp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("> ResolveUDPAddr() failed: %v", err)
    }

    // Listen for incoming UDP packets
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatalf("> ListenUDP() failed: %v", err)
    }
    defer conn.Close()

    buf := make([]byte, 1024)

    for {
        // Read a datagram
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("> recvfrom() failed:", err)
            continue
        }

        msg := strings.TrimSpace(string(buf[:n]))
        fmt.Println("> echoed:", msg)

        // Echo back to sender
        _, err = conn.WriteToUDP(buf[:n], clientAddr)
        if err != nil {
            fmt.Println("> sendto() failed:", err)
            continue
        }
    }

}
