package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"

	"github.com/JulianVidal/tagger/internal/command"
	"github.com/JulianVidal/tagger/internal/socket"
	"github.com/JulianVidal/tagger/server/internal/engine"
)

func main() {
	// tagger.Execute()
	engine.InitEngine()

	engine.AddTag("BT", []string{})
	engine.AddTag("Sat", []string{"BT"})

	engine.AddObj("lab_doc.doc", "Word Document", []string{"Sat"})
	engine.AddObj("start.ptx", "Powerpoint Presentation", []string{"BT"})

	println("Query: BT")
	for _, obj := range engine.Query([]string{"BT"}) {
		obj.Print()
	}

	println("Query: Sat")
	for _, obj := range engine.Query([]string{"Sat"}) {
		obj.Print()
	}

	println("All:")

	engine.Print()

	fmt.Println("-----------------------------")

	data := engine.ToJson()

	err := os.WriteFile("engine.json", data, 0644)
	if err != nil {
		panic("Couldn't write engine to file")
	}

	file, err := os.ReadFile("engine.json")
	if err != nil {
		panic("Couldn't read file")
	}

	engine.FromJson(file)

	engine.Print()

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
	var com command.Packet
	err = json.Unmarshal(buffer[:mLen], &com)
	if err != nil {
		fmt.Printf("Error unmarhsalling buffer into command. %s\n", err)
		return
	}

	fmt.Printf("Unmarshal: %+v\n", com)

	err = runCommand(com)

	if err != nil {
		connection.Write([]byte("Command failed"))
	}
}

func runCommand(com command.Packet) error {
	switch com.Type {
	case command.AddTag:
	case command.DelTag:

	case command.AddObj:
		obj := com.Data.(command.AddObjData).Obj
		engine.AddObj(obj.Name, obj.Format, obj.Tags)
		// object := com.Data.(types.Object)
		// engine.AddObj(object.Nameobject.Format)
	case command.DelObj:

	case command.Query:

	default:
		return errors.New("Unknown command")
	}

	return nil
}
