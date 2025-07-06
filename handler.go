package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{} // using sync.RWMutex because our server should handle requests concurrently

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{
			typ: "string",
			str: "PONG",
		}
	}
	return Value{
		typ: "string",
		str: args[0].bulk, // Return the first argument as the response
	}
}

func set(args []Value) Value {
	return Value{}
}

func get(args []Value) Value {
	return Value{}
}
