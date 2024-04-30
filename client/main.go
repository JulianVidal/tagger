package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/JulianVidal/tagger/internal/command"
	"github.com/JulianVidal/tagger/internal/serialized"
	"github.com/JulianVidal/tagger/internal/socket"
)

// TODO: Move the tag and objects added in server and send them through commands
// TODO: Create CLI were you point it to a file, give tags and it should figure out the rest
func main() {
	fmt.Printf("Dialing: " + socket.ADDRESS + "\n")

	connection, err := net.Dial(socket.TYPE, socket.ADDRESS)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	obj := serialized.Obj{
		Name:   "temp",
		Format: "temp format",
	}

	com := command.Packet{
		Type: command.AddObj,
		Data: command.AddObjData{
			Obj: obj,
		},
	}

	sendCommand(connection, com)
}

func sendCommand(connection net.Conn, com command.Packet) error {
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
