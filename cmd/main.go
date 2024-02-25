package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/bobhonores/godis/internal/aof"
	"github.com/bobhonores/godis/internal/commands"
	"github.com/bobhonores/godis/internal/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Creating a server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Starting AOF
	aof, err := aof.NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(t resp.Token) {
		command := strings.ToUpper(t.Array[0].Bulk)
		args := t.Array[1:]

		handler, ok := commands.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			return
		}

		handler(args)
	})

	// Listening connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		reader := resp.NewReader(conn)
		value, err := reader.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		// TODO: maybe this could be a validation inside the module
		if value.Typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		// TODO: maybe this could be a validation inside the module
		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		writer := resp.NewWriter(conn)

		// TODO: maybe this could be simplify
		// the command is passed to a module and whether the
		// outcome is ok or not, return a value (token)
		// to expose outside
		handler, ok := commands.Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(resp.Token{Typ: "string", Str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
