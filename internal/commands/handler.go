package commands

import (
	"sync"

	"github.com/bobhonores/slamigan/internal/resp"
)

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
	// "HGETALL": hgetall, // TODO: missing implementation
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{
			Typ: "string",
			Str: "PONG",
		}
	}

	return resp.Value{
		Typ: "string",
		Str: args[0].Bulk,
	}
}

var SETs = map[string]string{}
var mutexSet = sync.RWMutex{}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'set' command",
		}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	mutexSet.Lock()
	SETs[key] = value
	mutexSet.Unlock()

	return resp.Value{
		Typ: "string",
		Str: "OK",
	}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'get' command",
		}
	}

	key := args[0].Bulk

	mutexSet.RLock()
	value, ok := SETs[key]
	mutexSet.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{
		Typ:  "bulk",
		Bulk: value,
	}
}

var HSETs = map[string]map[string]string{}
var mutexHset = sync.RWMutex{}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'hset' command",
		}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	mutexHset.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	mutexHset.Unlock()

	return resp.Value{
		Typ: "string",
		Str: "OK",
	}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{
			Typ: "error",
			Str: "ERR wrong number of arguments for 'hget' command",
		}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	mutexHset.RLock()
	value, ok := HSETs[hash][key]
	mutexHset.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{
		Typ:  "bulk",
		Bulk: value,
	}
}
