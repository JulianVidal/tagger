package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/JulianVidal/tagger/internal/command"
	"github.com/JulianVidal/tagger/internal/socket"
	"github.com/JulianVidal/tagger/server/internal/engine"
)

func main() {
	// tagger.Execute()
	engine.InitEngine()

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
			fmt.Printf("New Client rejected: %s\n", err)
			continue
		}

		fmt.Println("New Client accepted")
		processClient(connection)
	}
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	for {
		mLen, err := connection.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection")
			return
		}

		fmt.Println("Received: ", string(buffer[:mLen]))
		var com command.Packet
		err = json.Unmarshal(buffer[:mLen], &com)
		if err != nil {
			fmt.Printf("Error unmarhsalling into command: %s\n", err)
			return
		}

		fmt.Printf("Unmarshal: %+v\n", com)

		result, err := runCommand(com)

		if err != nil {
			msg := fmt.Sprintf("Command failed due to: %s\n", err)
			connection.Write([]byte(msg))
		} else {
			connection.Write([]byte(result))
			engine.Print()
		}
	}
}

func runCommand(com command.Packet) (string, error) {
	switch com.Type {
	case command.AddTag:
		tag := com.Data.(command.AddTagData).Tag
		err := engine.AddTag(tag.Name, tag.Tags)
		if err != nil {
			return "", fmt.Errorf("Couldn't add tag '%s' due to: %s", tag.Name, err)
		}
		return "Added tag", nil
	case command.DelTag:

	case command.AddObj:
		obj := com.Data.(command.AddObjData).Obj
		err := engine.AddObj(obj.Name, obj.Format, obj.Tags)
		if err != nil {
			return "", fmt.Errorf("Couldn't add object '%s' due to: %s", obj.Name, err)
		}
		return "Added object", nil
	case command.DelObj:

	case command.Query:
		query := com.Data.(command.QueryData).Tags
		objs, err := engine.Query(query)
		if err != nil {
			return "", fmt.Errorf("Coulnd't fullfill query '%s' due to : %s", query, err)
		}
		result := ""
		for _, obj := range objs {
			result += fmt.Sprintf("%s\n", obj)
		}
		return result, nil

	default:
		return "", errors.New("Unknown command")
	}

	return "Sucessfully ran command", nil
}
