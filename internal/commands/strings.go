package commands

import (
	"sync"

	"github.com/bobhonores/godis/internal/resp"
)

var hashMaps = map[string]string{}
var mutexHashMaps = sync.RWMutex{}

func set(args []resp.Token) resp.Token {
	if len(args) != 2 {
		return resp.Token{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'set' command",
		}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	mutexHashMaps.Lock()
	hashMaps[key] = value
	mutexHashMaps.Unlock()

	return resp.Token{
		Typ: "string",
		Str: "OK",
	}
}

func get(args []resp.Token) resp.Token {
	if len(args) != 1 {
		return resp.Token{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'get' command",
		}
	}

	key := args[0].Bulk

	mutexHashMaps.RLock()
	value, ok := hashMaps[key]
	mutexHashMaps.RUnlock()

	if !ok {
		return resp.Token{Typ: "null"}
	}

	return resp.Token{
		Typ:  "bulk",
		Bulk: value,
	}
}
