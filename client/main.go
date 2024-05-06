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

	// obj := serialized.Obj{
	// 	Name:   "temp",
	// 	Format: "temp format",
	// }
	//
	// com := command.Packet{
	// 	Type: command.AddObj,
	// 	Data: command.AddObjData{
	// 		Obj: obj,
	// 	},
	// }

	comAddBT := command.Packet{
		Type: command.AddTag,
		Data: command.AddTagData{
			Tag: serialized.Tag{
				Name: "BT",
			},
		},
	}

	comAddSat := command.Packet{
		Type: command.AddTag,
		Data: command.AddTagData{
			Tag: serialized.Tag{
				Name: "Sat",
				Tags: []string{"BT"},
			},
		},
	}

	comAddLab := command.Packet{
		Type: command.AddObj,
		Data: command.AddObjData{
			Obj: serialized.Obj{
				Name:   "lab_doc.doc",
				Format: "Word Document",
				Tags:   []string{"Sat"},
			},
		},
	}

	comAddStart := command.Packet{
		Type: command.AddObj,
		Data: command.AddObjData{
			Obj: serialized.Obj{
				Name:   "start.ptx",
				Format: "Powerpoint Presentation",
				Tags:   []string{"Sat"},
			},
		},
	}

	comQryBT := command.Packet{
		Type: command.Query,
		Data: command.QueryData{
			Tags: []string{"BT"},
		},
	}

	comQrySat := command.Packet{
		Type: command.Query,
		Data: command.QueryData{
			Tags: []string{"Sat"},
		},
	}

	sendCommand(connection, comAddBT)
	sendCommand(connection, comAddSat)
	sendCommand(connection, comAddLab)
	sendCommand(connection, comAddStart)

	_, err = sendCommand(connection, comQryBT)
	if err != nil {
		fmt.Printf("Couldn't get query due to error: '%s'\n", err)
	}

	_, err = sendCommand(connection, comQrySat)
	if err != nil {
		fmt.Printf("Couldn't get query due to error: '%s'\n", err)
	}

}

func sendCommand(connection net.Conn, com command.Packet) (string, error) {
	commandJSON, err := json.Marshal(com)
	if err != nil {
		return "", fmt.Errorf("Couldn't encode command into JSON: %w", err)
	}

	fmt.Printf("Sent: %s\n", string(commandJSON))

	_, err = connection.Write(commandJSON)
	if err != nil {
		return "", fmt.Errorf("Couldn't write json to socket: %w", err)
	}

	buffer := make([]byte, 1024)

	mLen, err := connection.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("Couldn't read data from buffer: %w", err)
	}

	fmt.Printf("Received: %s\n", string(buffer[:mLen]))

	return string(buffer[:mLen]), nil
}
