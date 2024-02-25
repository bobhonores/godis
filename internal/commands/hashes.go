package commands

import (
	"sync"

	"github.com/bobhonores/godis/internal/resp"
)

var hashes = map[string]map[string]string{}
var mutexHashes = sync.RWMutex{}

func hset(args []resp.Token) resp.Token {
	if len(args) != 3 {
		return resp.Token{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'hset' command",
		}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	mutexHashes.Lock()
	if _, ok := hashes[hash]; !ok {
		hashes[hash] = map[string]string{}
	}
	hashes[hash][key] = value
	mutexHashes.Unlock()

	return resp.Token{
		Typ: "string",
		Str: "OK",
	}
}

func hget(args []resp.Token) resp.Token {
	if len(args) != 2 {
		return resp.Token{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'hget' command",
		}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	mutexHashes.RLock()
	value, ok := hashes[hash][key]
	mutexHashes.RUnlock()

	if !ok {
		return resp.Token{Typ: "null"}
	}

	return resp.Token{
		Typ:  "bulk",
		Bulk: value,
	}
}

func hgetAll(args []resp.Token) resp.Token {
	if len(args) != 1 {
		return resp.Token{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'hgetall' command",
		}
	}

	hash := args[0].Bulk
	content := make([]resp.Token, 0)

	mutexHashes.RLock()
	value, ok := hashes[hash]
	for k, v := range value {
		content = append(content, resp.Token{Typ: "bulk", Bulk: k})
		content = append(content, resp.Token{Typ: "bulk", Bulk: v})
	}
	mutexHashes.RUnlock()

	if !ok {
		return resp.Token{Typ: "null"}
	}

	return resp.Token{
		Typ:   "array",
		Array: content,
	}
}
