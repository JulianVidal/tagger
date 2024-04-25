package main

// import "github.com/JulianVidal/tagger/cmd/tagger"
import (
	"fmt"
	"net"
	"os"
)

const (
	SERVER_NETWORK = "localhost"
	SERVER_TYPE    = "tcp"
)

func main() {
	port := os.Args[1]
	connection, err := net.Dial(SERVER_TYPE, SERVER_NETWORK+":"+port)
	if err != nil {
		panic(err)
	}

	_, err = connection.Write([]byte("Test write to socket"))
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	defer connection.Close()
}
