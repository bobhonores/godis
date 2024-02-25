package commands

import (
	"github.com/bobhonores/godis/internal/resp"
)

var Handlers = map[string]func([]resp.Token) resp.Token{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetAll,
}
