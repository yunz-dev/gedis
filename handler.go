package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) Value {
  "PING": ping,
  "SET": set,
  "DEL": del,
  "GET": get,
  "HSET": hset,
  "HGET": hget,
  "HDEL": hdel,
  "HGETALL": hgetall,
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


func del(args []Value) Value {
    key := args[0].bulk
    deletedCount := 0

    // Lock both maps at the same time to ensure atomicity
    SETsMu.Lock()
    if _, exists := SETs[key]; exists {
        delete(SETs, key)
        deletedCount++
    }
    SETsMu.Unlock()

    HSETsMu.Lock()
    if _, exists := HSETs[key]; exists {
        delete(HSETs, key)
        deletedCount++
    }
    HSETsMu.Unlock()

    return Value{typ: "integer", num: deletedCount}
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

func hdel(args []Value) Value {
  if len(args) < 2 {
    return Value{typ: "error", str: "ERR wrong number of arguments for HDEL command"}
  }

  hash := args[0].bulk
  keys := args[1:] // Remaining arguments are field keys to delete
  deletedCount := 0

  HSETsMu.Lock()
  if fields, exists := HSETs[hash]; exists {
    for _, key := range keys {
      if _, found := fields[key.bulk]; found {
        delete(fields, key.bulk)
        deletedCount++
      }
    }
    // If the hash is now empty, remove it entirely
    if len(fields) == 0 {
      delete(HSETs, hash)
    }
  }
  HSETsMu.Unlock()

  return Value{typ: "integer", num: deletedCount}
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

func hgetall(args []Value) Value {
  if len(args) != 1 {
    return Value{typ: "error", str: "ERR wrong number of arguments for HGETALL command"}
  }

  hash := args[0].bulk

  HSETsMu.RLock()
  fields, exists := HSETs[hash]
  HSETsMu.RUnlock()

  if !exists {
    return Value{typ: "array", array: []Value{}}
  }

  result := make([]Value, 0, len(fields)*2)
  for key, value := range fields {
    result = append(result, Value{typ: "string", str: key})
    result = append(result, Value{typ: "string", str: value})
  }

  return Value{typ: "array", array: result}
}
