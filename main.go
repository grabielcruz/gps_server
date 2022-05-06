package main

import (
	"fmt"
	"log"
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
	count := 0
	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			defer conn.Close()
			break
		}

		file, err := os.OpenFile("gps_log.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err.Error())
			file, err = os.Create("gps_log.txt")
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		defer file.Close()

		text := string(buf[:reqLen]) + "\n"
		fmt.Println(text)
		fmt.Println(reqLen)
		// _, err = file.WriteString(text)
		if err != nil {
			log.Fatal(err.Error())
		}
		if reqLen == 26 {
			conn.Write([]byte("LOAD"))
		} else if count%2 == 0 {
			fmt.Println("Sending this: **,imei:864035050161315,101,10s;")
			conn.Write([]byte("**,imei:864035050161315,101,10s;"))
		} else if reqLen == 16 {
			fmt.Println("Sending ON")
			conn.Write([]byte("ON"))
		}
		// conn.Close()
		count++
		fmt.Println("count: ", count)
	}
}
