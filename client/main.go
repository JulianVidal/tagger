package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/JulianVidal/tagger/internal/command"
	"github.com/JulianVidal/tagger/internal/serialize"
	"github.com/JulianVidal/tagger/internal/socket"
)

func main() {
	fmt.Printf("Dialing: " + socket.ADDRESS + "\n")

	connection, err := net.Dial(socket.TYPE, socket.ADDRESS)
	if err != nil {
		panic(err)
	}

	obj := serialize.Obj{
		Name:   "temp",
		Format: "temp format",
		Tags:   nil,
	}

	com, err := command.CreateCommand(command.Add, command.Object, obj)
	if err != nil {
		panic(err)
	}

	sendCommand(connection, com)

	defer connection.Close()
}

func sendCommand(connection net.Conn, com command.Command) error {
	commandJSON, err := json.Marshal(com)
	if err != nil {
		return fmt.Errorf("Couldn't encode command into JSON: %w", err)
	}

	fmt.Printf("Sent: %s\n", string(commandJSON))

	_, err = connection.Write(commandJSON)
	if err != nil {
		return fmt.Errorf("Couldn't write json to socket: %w", err)
	}

	buffer := make([]byte, 1024)

	mLen, err := connection.Read(buffer)
	if err != nil {
		return fmt.Errorf("Couldn't read data from buffer: %w", err)
	}

	fmt.Printf("Received: %s", string(buffer[:mLen]))

	return nil
}
