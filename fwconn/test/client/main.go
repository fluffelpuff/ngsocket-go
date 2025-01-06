package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/fluffelpuff/ngsocket-go/fwconn"
)

// Generate1MBString erzeugt einen String mit einer Größe von 1 MB.
func Generate1MBString() string {
	// Der Buchstabe 'a' hat eine Größe von 1 Byte in UTF-8.
	return strings.Repeat("a", 5000)
}

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:4433")
	if err != nil {
		fmt.Println(err)
		return
	}
	upgradedConn, err := fwconn.UpgradeConn(conn, fwconn.Client, []string{"1.0"})
	if err != nil {
		return
	}

	err = upgradedConn.Write([]byte(Generate1MBString()))
	if err != nil {
		panic(err)
	}

}
