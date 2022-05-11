package main

import (
	"fmt"
	"net"
	"sync"
)

var lock = &sync.Mutex{}

type connectionsMap struct {
	collection map[string]net.Conn
}

var connections *connectionsMap

func GetConnections() *connectionsMap {
	if connections == nil {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("Creating single instance now.")
		connections = &connectionsMap{}
	} else {
		fmt.Println("Single instance already created.")
	}
	return connections
}
