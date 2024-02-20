package commands

import "github.com/bobhonores/slamigan/internal/resp"

func ping(args []resp.Token) resp.Token {
	if len(args) == 0 {
		return resp.Token{
			Typ: "string",
			Str: "PONG",
		}
	}

	return resp.Token{
		Typ: "string",
		Str: args[0].Bulk,
	}
}
