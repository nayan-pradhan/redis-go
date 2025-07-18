package main

import (
	"fmt"
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
	"DEL":     del,
	"HDEL":    hdel,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{} // using sync.RWMutex because our server should handle requests concurrently
var HSETs = map[string]map[string]string{}
var HSETsMU = sync.RWMutex{} // using sync.RWMutex because our server should handle requests concurrently

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
	if len(args) != 2 {
		fmt.Println("Invalid. Args received: ", args)
		return Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'set' command, should receive key value!",
		}
	}
	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock() // blocks all other goroutines (read + write) until lock is released
	SETs[key] = value
	SETsMu.Unlock()
	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 { // argument should only be key to search
		fmt.Println("Invalid. Arg received: ", args)
		return Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'get' command, should receive a key to search!",
		}
	}
	key := args[0].bulk

	SETsMu.RLock() // multiple goroutines can read at same time as long as no goroutine holds write lock
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{
			typ: "null",
		}
	}
	return Value{
		typ:  "bulk",
		bulk: value,
	}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		fmt.Println("Invalid. Args received: ", args)
		return Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'hset' command, should receive a hash key value pair to store!",
		}
	}
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMU.Lock()
	if _, ok := HSETs[hash]; !ok { // checks if key exists in hash map, if not init empty hash value
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMU.Unlock()

	return Value{
		typ: "string",
		str: "OK",
	}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		fmt.Println("Invalid. Args received: ", args)
		return Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'hget' command, should receive a hash and key to get value",
		}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMU.RLock()
	value, ok := HSETs[hash][key]
	HSETsMU.RUnlock()

	if !ok {
		return Value{
			typ: "null",
		}
	}

	return Value{
		typ:  "bulk",
		bulk: value,
	}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		fmt.Println("Invalid. Args received: ", args)
		return Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'HGETALL', should receive a single key value",
		}
	}

	hash := args[0].bulk
	HSETsMU.RLock()
	value, ok := HSETs[hash]
	HSETsMU.RUnlock()
	if !ok {
		return Value{
			typ: "null",
		}
	}
	arrayValue := make([]Value, 0)
	for k, v := range value {
		arrayValue = append(arrayValue, Value{typ: "bulk", bulk: k})
		arrayValue = append(arrayValue, Value{typ: "bulk", bulk: v})
	}

	return Value{
		typ:   "array",
		array: arrayValue,
	}
}

func del(args []Value) Value {
	if len(args) < 1 {
		fmt.Println("Invalid number of arguments of 'del' command, should receive key/s to delete from storage.")
		return Value{
			typ: "error",
			str: "ERR wrong number of arguments for 'DEL', should receive a key/s",
		}
	}

	count := 0
	keys := args[0:]
	SETsMu.Lock()
	defer SETsMu.Unlock()

	for k := range keys {
		_, ok := SETs[keys[k].bulk]
		if !ok { // key does not exist in map
			fmt.Printf("%s does not exist in memory, no action taken.\n", keys[k].bulk)
		} else {
			delete(SETs, keys[k].bulk)
			count += 1
		}
	}

	return Value{
		typ: "integer",
		num: count,
	}
}

func hdel(args []Value) Value {
	if len(args) < 2 {
		return Value{
			typ: "error",
			str: "ERR wrong number of arguments of 'hdel' command, should receive hash key/s to delete from memory",
		}
	}

	count := 0
	hash := args[0].bulk
	keys := args[1:]

	HSETsMU.Lock()
	defer HSETsMU.Unlock()

	_, ok := HSETs[hash]
	if !ok {
		fmt.Printf("%s does not exist as hash.\n", hash)
	} else {
		for k := range keys {
			_, ok = HSETs[hash][keys[k].bulk]
			if !ok {
				fmt.Printf("%s does not exist in hash with key\n", keys[k].bulk)
			} else {
				delete(HSETs[hash], keys[k].bulk)
				count += 1
			}
		}
	}

	return Value{
		typ: "integer",
		num: count,
	}

}
