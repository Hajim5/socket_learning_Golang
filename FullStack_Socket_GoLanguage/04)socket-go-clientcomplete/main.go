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
	runClient()
	fmt.Println("> echo-client is de-activated")
}

func runClient() {

	conn, err := net.Dial("tcp", HOST+":"+PORT)
	if err != nil {
		fmt.Println("> connect() failed and program terminated")
		fmt.Println("> error:", err)
		return
	}
	defer conn.Close()

	serverReader := bufio.NewReader(conn)
	stdinReader := bufio.NewReader(os.Stdin)

	for {
		// [=start=]
		fmt.Print("> ")
		sendMsg, err := stdinReader.ReadString('\n')
		if err != nil {
			log.Println("> input failed by error:", err)
			return
		}
		// remove trailing '\n'
		if len(sendMsg) > 0 && sendMsg[len(sendMsg)-1] == '\n' {
			sendMsg = sendMsg[:len(sendMsg)-1]
		}

		// send data
		_, err = conn.Write([]byte(sendMsg + "\n"))
		if err != nil {
			fmt.Println("> send() failed and program terminated")
			fmt.Println("> error:", err)
			return
		}

		recvData, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("> recv() failed and program terminated")
			fmt.Println("> error:", err)
			return
		}

		if len(recvData) > 0 && recvData[len(recvData)-1] == '\n' {
			recvData = recvData[:len(recvData)-1]
		}

		fmt.Println("> received:", recvData)

		if sendMsg == "quit" {
			break
		}
	}
}
