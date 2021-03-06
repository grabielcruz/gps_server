package main

import (
	"fmt"
	"gps_server/connections"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8081"
	CONN_TYPE = "tcp"
)

func main() {
	go listenTCP()

	r := gin.Default()
	r.POST("/gps/:imei/*action", func(c *gin.Context) {
		imei := c.Param("imei")
		action := c.Param("action")
		singleConnections := connections.GetConnections()
		conn := singleConnections.Collection[imei]
		fmt.Println("Sending to connection ", action)
		conn.Write([]byte(action))
		c.String(http.StatusOK, "command", action, " received fro imei ", imei)
	})
	r.Run(":8080")
}

func listenTCP() {
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
	// reader := bufio.NewReader(os.Stdin)
	file, err := os.OpenFile("gps_log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
		file, err = os.Create("gps_log.txt")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer file.Close()

	for {
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error reading: ", err.Error())
			defer conn.Close()
			break
		}

		fmt.Println("reqLen: ", reqLen)
		text := string(buf[:reqLen])
		if reqLen == 26 {
			singleConnections := connections.GetConnections()
			keys := strings.Split(text, ",")
			imei := strings.Split(keys[1], ":")[1]
			fmt.Println("imei: ", imei)
			singleConnections.Collection[imei] = conn
			// conn.Write([]byte("LOAD"))
			fmt.Println("singleConnections: ", singleConnections)
		}

		// if reqLen == 16 {
		// 	conn.Write([]byte("ON"))
		// }

		fmt.Println(text)
		_, err = file.WriteString(text)
		if err != nil {
			fmt.Println(err.Error())
		}
		// input, err := reader.ReadString('\n')
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		// if len(input) > 0 {
		// 	conn.Write([]byte(input))
		// }
	}
}
