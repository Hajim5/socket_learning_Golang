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
    fmt.Println("> multithread echo-client is activated")

    // Connect to server
    conn, err := net.Dial("tcp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("> connect() failed: %v", err)
    }
    defer conn.Close()

    // Channel to signal when we should stop
    done := make(chan struct{})

    // Start receiver goroutine (like Python's receive thread)
    go receiveLoop(conn, done)

    // Main goroutine: send user input
    sendLoop(conn)

    // If we typed "quit", tell receiver to stop
    close(done)

    fmt.Println("> multithread echo-client is de-activated")
}

// Continuously read data from server and print it
func receiveLoop(conn net.Conn, done <-chan struct{}) {
    reader := bufio.NewReader(conn)

    for {
        select {
        case <-done:
            // main decided to quit
            return
        default:
            // try to read from server
            msg, err := reader.ReadString('\n')
            if err != nil {
                // server closed or error
                fmt.Println("> recv() failed or server closed connection:", err)
                return
            }

            msg = strings.TrimRight(msg, "\r\n")
            fmt.Println("> received:", msg)

            if msg == "quit" {
                // server sent quit, just stop
                return
            }
        }
    }
}

// Read from stdin and send to server
func sendLoop(conn net.Conn) {
    stdin := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("> ")
        text, err := stdin.ReadString('\n')
        if err != nil {
            fmt.Println("> input error:", err)
            return
        }

        text = strings.TrimRight(text, "\r\n")

        // send to server
        _, err = conn.Write([]byte(text + "\n"))
        if err != nil {
            fmt.Println("> send() failed:", err)
            return
        }

        if text == "quit" {
            // tell server we are quitting
            return
        }
    }
}
