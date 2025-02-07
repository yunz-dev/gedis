package main

import "sync"

var Handlers = map[string]func([]Value) Value {
  "PING": ping,
  "SET": set,
  "GET": get,
  "HSET": hget,
  "HGET": hget,
}

func ping(args []Value) Value {
  if len(args) == 0 {
    return Value{typ: "string", str: "PONG"}
  }
  return Value{typ: "string", str: args[0].bulk}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
  if len(args) != 2 {
    return Value{typ: "error", str: "ERR wrong number of arguments for SET command"}
  }

  key := args[0].bulk
  value := args[1].bulk

  SETsMu.Lock()
  SETs[key] = value
  SETsMu.Unlock()

  return Value{typ: "string", str: "OK"}
}


func get(args []Value) Value {
  if len(args) != 1 {
    return Value{typ: "error", str: "ERR wrong number of arguments for GET command"}
  }

  key := args[0].bulk

  SETsMu.RLock()
  value, ok := SETs[key]
  SETsMu.RUnlock()

  if !ok {
    return Value{typ: "null"}
  }

  return Value{typ: "bulk", bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
  if len(args) != 3 {
    return Value{typ: "error", str: "ERR wrong number of arguments for HSET command"}
  }

  hash := args[0].bulk
  key := args[1].bulk
  value := args[2].bulk

  HSETsMu.Lock()
  if _, ok := HSETs[hash]; !ok {
    HSETs[hash] = map[string]string{}
  }
  HSETs[hash][key] = value
  HSETsMu.Unlock()

  return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
  if len(args) != 2 {
    return Value{typ: "error", str: "ERR wrong number of arguments for HSET command"}
  }

  hash := args[0].bulk
  key := args[1].bulk

  HSETsMu.RLock()
  value, ok := HSETs[hash][key]
  HSETsMu.RUnlock()

  if !ok {
    return Value{typ: "null"}
  }

  return Value{typ: "bulk", bulk: value}
}
