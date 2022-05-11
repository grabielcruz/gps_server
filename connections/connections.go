package connections

import (
	"fmt"
	"net"
	"sync"
)

var lock = &sync.Mutex{}

type ConnectionsMap struct {
	Collection map[string]net.Conn
}

var connections *ConnectionsMap

func GetConnections() *ConnectionsMap {
	if connections == nil {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("Creating single instance now.")
		connections = &ConnectionsMap{}
	} else {
		fmt.Println("Single instance already created.")
	}
	return connections
}
