package main

import (
	"bufio"
	"fmt"
	"gps_server/connections"
	"net"
	"os"
	"strings"
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
	reader := bufio.NewReader(os.Stdin)
	stored := false
	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			defer conn.Close()
			break
		}

		fmt.Println("reqLen: ", reqLen)
		if !stored && reqLen == 26 {
			singleConnections := connections.GetConnections()
			keys := strings.Split(string(buf[:reqLen]), ",")
			imei := strings.Split(keys[1], ":")[1]
			fmt.Println("imei: ", imei)
			singleConnections.Collection[imei] = conn
			stored = true
			conn.Write([]byte("LOAD"))
			fmt.Println("singleConnections: ", singleConnections)
		}

		if reqLen == 16 {
			conn.Write([]byte("ON"))
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

		fmt.Println(string(buf[:reqLen]))
		// _, err = file.WriteString(text)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
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
