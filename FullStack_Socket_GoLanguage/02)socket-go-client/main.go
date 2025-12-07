package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

const (
    HOST = "127.0.0.1"
    PORT = "65456"
)

func main() {
    fmt.Println("> echo-client is activated")

    //Connect to a aserver
    // Creates a connection to the server in this example 127.0.0.1:65456
    conn, err := net.Dial("tcp", HOST+":"+PORT)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)
    input := bufio.NewReader(os.Stdin)

    for {
        fmt.Print("> ")
        sendMsg, _ := input.ReadString('\n')
        sendMsg = sendMsg[:len(sendMsg)-1] 

        //client sends data after user put the input
        _, err := conn.Write([]byte(sendMsg + "\n"))
        if err != nil {
            log.Println("Write error:", err)
            break
        }
        //client receives the server responses if the input was sent
        recv, _ := reader.ReadString('\n')
        fmt.Print("> received: ", recv)

        if sendMsg == "quit" {
            break
        }
    }

    fmt.Println("> echo-client is de-activated")
}
