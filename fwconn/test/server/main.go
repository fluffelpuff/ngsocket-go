package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fluffelpuff/ngsocket-go/fwconn"
)

func main() {
	// Erstelle einen TLS-Listener auf Port 4433
	listener, err := net.Listen("tcp", ":4433")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server läuft auf Port 4433...")

	// Akzeptiere eingehende Verbindungen
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Failed to accept connection: %v", err)
	}

	upgradedConn, err := fwconn.UpgradeConn(conn, fwconn.Server, []string{"1.0"})
	if err != nil {
		return
	}

	data, err := upgradedConn.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(len(data))
}
