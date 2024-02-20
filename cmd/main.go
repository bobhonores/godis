package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/bobhonores/slamigan/internal/commands"
	"github.com/bobhonores/slamigan/internal/resp"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

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

		result := handler(args)
		writer.Write(result)
	}
}
