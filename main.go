package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8081"
	CONN_TYPE = "tcp"
)

var connections map[string]net.Conn

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
	text := ""
	buf := make([]byte, 1024)
	reader := bufio.NewReader(os.Stdin)
	reqLen := 0
	var err error
	stored := false
	for {
		reqLen, err = conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			defer conn.Close()
			break
		}

		fmt.Println("reqLen: ", reqLen)
		if !stored && reqLen == 26 {
			text = string(buf[:reqLen])
			keys := strings.Split(text, ",")
			imei := strings.Split(keys[1], ":")[1]
			fmt.Println("imei: ", imei)
			connections[imei] = conn
			fmt.Println("connections: ", connections)
			stored = true
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

		text = string(buf[:reqLen])
		fmt.Println(text)
		// _, err = file.WriteString(text)
		if err != nil {
			fmt.Println(err.Error())
		}
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(input) > 0 {
			conn.Write([]byte(input))
		}

		// conn.Close()

	}
}
