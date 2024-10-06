package main

import "sync"

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

var keyValueStore = make(map[string]string)
var storeMutex = sync.RWMutex{}

func set(arguments []Value) Value {
	if len(arguments) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := arguments[0].bulk
	value := arguments[1].bulk

	storeMutex.Lock()
	keyValueStore[key] = value
	storeMutex.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(params []Value) Value {
	if len(params) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	requestedKey := params[0].bulk

	storeMutex.RLock()
	retrievedValue, exists := keyValueStore[requestedKey]
	storeMutex.RUnlock()

	if !exists {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: retrievedValue}
}
