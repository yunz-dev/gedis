package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {
  Addr := ""
  flag.StringVar(&Addr, "addr", ":6397", "HTTP network address")
  flag.Parse()
  fmt.Println("Listening on port", Addr)
  // start TCP Listen
  l, err := net.Listen("tcp", Addr)
  if err != nil {
    fmt.Println(err)
    return
  }

  aof, err := NewAof("database.aof")
  if err != nil {
    fmt.Println(err)
    return
  }
  defer aof.Close()

  aof.Read(func(value Value) {
    command := strings.ToUpper(value.array[0].bulk)
    args := value.array[1:]

    handler, ok := Handlers[command]
    if !ok {
      fmt.Println("Invalid command: ", command)
      return
    }
    handler(args)
  })

  // Listen for connections
  conn, err := l.Accept()
  if err != nil {
    fmt.Println(err)
    return
  }

  defer conn.Close()

  for {
    resp := NewResp(conn)
    value, err := resp.Read()
    if err != nil {
      fmt.Println(err)
      return
    }

    if value.typ != "array" {
      fmt.Println("Invalid request, expected array")
      continue
    }

    if len(value.array) == 0 {
      fmt.Println("Invalid request, expected array length > 0")
      continue
    }

    command := strings.ToUpper(value.array[0].bulk)
    args := value.array[1:]

    writer := NewWriter(conn)

    handler, ok := Handlers[command]
    if !ok {
      fmt.Println("Invalid commands: ", command)
      writer.Write(Value{typ: "string", str: ""})
      continue
    }

    if command == "SET" || command == "HSET" || command == "HDEL" || command == "DEL" {
      aof.Write(value)
    }

    result := handler(args)
    writer.Write(result)

  }
}
