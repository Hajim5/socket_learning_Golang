package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
)

const (
    HOST = "127.0.0.1"
    PORT = "65456"
)

func main() {
    fmt.Println("> udp-multithread echo-client is activated")

    // Resolve server UDP address
    serverAddr, err := net.ResolveUDPAddr("udp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("> ResolveUDPAddr() failed: %v", err)
    }

    // Create UDP connection 
    conn, err := net.DialUDP("udp", nil, serverAddr)
    if err != nil {
        log.Fatalf("> DialUDP() failed: %v", err)
    }
    defer conn.Close()

    // Channel to signal receiver goroutine to stop
    done := make(chan struct{})

    // Start receiver goroutine
    go receiveLoop(conn, done)

    // Main goroutine: send user input
    sendLoop(conn)

    // User typed "quit" → stop receiver
    close(done)

    fmt.Println("> udp-multithread echo-client is de-activated")
}

// Continuously receive data from server and print it
func receiveLoop(conn *net.UDPConn, done <-chan struct{}) {
    buf := make([]byte, 1024)

    for {
        select {
        case <-done:
            return
        default:
            n, _, err := conn.ReadFromUDP(buf)
            if err != nil {
                fmt.Println("> recvfrom() failed or server closed:", err)
                return
            }

            msg := strings.TrimSpace(string(buf[:n]))
            fmt.Println("> received:", msg)

            if msg == "quit" {
                return
            }
        }
    }
}

// Read from stdin and send datagrams to server
func sendLoop(conn *net.UDPConn) {
    stdin := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("> ")
        text, err := stdin.ReadString('\n')
        if err != nil {
            fmt.Println("> input error:", err)
            return
        }

        text = strings.TrimSpace(text)

        // Send as a UDP datagram (no newline needed)
        _, err = conn.Write([]byte(text))
        if err != nil {
            fmt.Println("> sendto() failed:", err)
            return
        }

        if text == "quit" {
            return
        }
    }
}
