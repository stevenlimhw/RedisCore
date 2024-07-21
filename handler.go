package main

// map with handler names as keys and the function as values
var Handlers = map[string]func([]Value) Value{
  "COMMAND": command,
  "PING": ping,
}

func command(args []Value) Value {
  return Value{typ: "string", str: ""}
}

func ping(args []Value) Value {
  if len(args) == 0 {
    return Value{typ: "string", str: "PONG"}
  }
  return Value{typ: "string", str: args[0].bulk}
}
