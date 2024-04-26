package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/JulianVidal/tagger/internal/command"
	"github.com/JulianVidal/tagger/internal/serialize"
	"github.com/JulianVidal/tagger/internal/socket"
	"github.com/JulianVidal/tagger/server/internal/engine"
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

	data := eng.ToJson()

	err := os.WriteFile("engine.json", data, 0644)
	if err != nil {
		panic("Couldn't write engine to file")
	}

	file, err := os.ReadFile("engine.json")
	if err != nil {
		panic("Couldn't read file")
	}

	eng = engine.FromJson(file)

	eng.Print()

	run()
}

func run() {
	fmt.Println("Starting Server...")

	server, err := net.Listen(socket.TYPE, socket.ADDRESS)
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
		processClient(connection)
	}
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection")
		return
	}

	fmt.Println("Received: ", string(buffer[:mLen]))
	var com command.Command
	err = json.Unmarshal(buffer[:mLen], &com)
	if err != nil {
		fmt.Printf("Error unmarhsalling buffer into command. %s\n", err)
		return
	}

	fmt.Printf("Unmarshal: %+v\n", com)

	switch com.Noun {
	case command.Tag:
		var subject serialize.Tag
		err = json.Unmarshal(com.Subject, &subject)
		if err != nil {
			fmt.Printf("Error unmarshalling subject. %s\n", err)
		}
		fmt.Printf("Unmarshal Tag: %+v\n", subject)
	case command.Object:
		var subject serialize.Obj
		err = json.Unmarshal(com.Subject, &subject)
		if err != nil {
			fmt.Printf("Error unmarshalling subject. %s\n", err)
		}
		fmt.Printf("Unmarshal Object: %+v\n", subject)
	}
	if err != nil {
		fmt.Printf("Couldn't unmarshal subject. %s\n", err)
	}

	connection.Write([]byte("Command passed or failed"))
}
