//Listen on host 127.0.0.1 at port 65456
//Accept one client
//Receive data from the client and print it
//Send back the sama data(echo) at the client
//If quit = stop and exit

//Every go file start with package
//main = executable program 
//Go will search for main() inside a package to start the program
package main

import (
	"bufio" //easier to read from the connection line by line
	"fmt" //formatting I/O
	"log" //logging with timestamps and levels, good for errors
	"net" // networking package (TCP/UDP) replaces "socket" in Python
)
//using const because do not want to change the values
const (
	HOST = "127.0.0.1"
	PORT = "65456"
)
func main() { // = int main()
	fmt.Println("> echo-server is activated") // = print()
	//creating and binding listening socket
	ln, err := net.Listen("tcp", HOST+":"+PORT) // tcp = protocol type 
	if err != nil { // check if something went wrong
		log.Fatalf("Failed to listen: %v", err)
	}
	// quit = close this listener
	defer ln.Close()
	// aceepting a client
	conn, err := ln.Accept() //waits for one client to connect
	if err != nil {
		log.Fatalf("Failed to accept: %v", err)
	}
	// close client connection 
	defer conn.Close()
	//gets info about client (IP and port)
	fmt.Printf("> client connected: %s\n", conn.RemoteAddr().String())
	// Preparing to read data in a loop
	reader := bufio.NewScanner(conn)
	for reader.Scan() {
		recv := reader.Text()
		fmt.Println("> echoed:", recv)
		_, err := conn.Write([]byte(recv + "\n"))
		if err != nil {
			log.Println("Write error:", err)
			break
		}
		//stop the server
		if recv == "quit" {
			break
		}
	}
	fmt.Println("> echo-server is de-activated")
}
