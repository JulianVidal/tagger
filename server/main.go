package main

import (
	"fmt"
	"net"

	"github.com/JulianVidal/tagger/server/internal/engine"
)

const (
	SERVER_NETWORK = "localhost:0"
	SERVER_TYPE    = "tcp"
)

func main() {
	// tagger.Execute()
	eng := engine.NewEngine()

	eng.AddTag("BT", []string{})
	eng.AddTag("Sat", []string{"BT"})

	eng.AddObj("lab_doc.doc", "Word Document", []string{"Sat"})
	eng.AddObj("start.ptx", "Powerpoint Presentation", []string{"BT"})

	println("Query: BT")
	for _, obj := range eng.Query([]string{"BT"}) {
		obj.Print()
	}

	println("Query: Sat")
	for _, obj := range eng.Query([]string{"Sat"}) {
		obj.Print()
	}

	println("All:")

	eng.Print()

	fmt.Println("-----------------------------")

	eng.Json("engine.json")

	eng = engine.FromJson("engine.json")

	eng.Print()

	run()
}

func run() {
	fmt.Println("Starting Server...")

	server, err := net.Listen(SERVER_TYPE, SERVER_NETWORK)

	if err != nil {
		panic("Error listening: " + err.Error())
	}

	fmt.Printf("Server listening on network: %s\n", server.Addr())

	defer server.Close()

	for {
		connection, err := server.Accept()

		if err != nil {
			panic("Error accepting new connection")
		}

		fmt.Println("New Client accepted")
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connectoin")
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	connection.Write([]byte("Received message: " + string(buffer[:mLen])))
}
