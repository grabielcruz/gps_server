package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8081"
	CONN_TYPE = "tcp"
)

func main() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			defer conn.Close()
			break
		}

		// file, err := os.OpenFile("gps_log.txt", os.O_APPEND|os.O_WRONLY, 0644)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	file, err = os.Create("gps_log.txt")
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 	}
		// }
		// defer file.Close()

		text := string(buf[:reqLen]) + "\n"
		fmt.Println(text)
		// _, err = file.WriteString(text)
		if err != nil {
			fmt.Println(err.Error())
		}
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		conn.Write([]byte(input))

		// conn.Close()

	}
}
