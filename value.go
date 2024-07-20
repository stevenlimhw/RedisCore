package main

import "strconv"

type Value struct {

  // determine the data type carried by the value
  typ string

  // value of the string received from simple strings
  str string

  // value of the integer from integers
  num int

  // string received from RESP bulk strings
  bulk string

  // values received from RESP arrays
  array []Value

}

func (v Value) Marshal() []byte {
  switch v.typ {
  case "string":
    return v.marshalString()
  case "bulk":
    return v.marshalBulk()
  case "array":
    return v.marshalArray()
  case "null":
    return v.marshalNull()
  case "error":
    return v.marshalError()
  default:
    return []byte{}
  }
}

func (v Value) marshalString() []byte {
  var bytes []byte
  bytes = append(bytes, STRING)
  bytes = append(bytes, v.str...)
  bytes = append(bytes, '\r', '\n')
  return bytes
}

func (v Value) marshalBulk() []byte {
  var bytes []byte
  bytes = append(bytes, BULK)
  bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
  bytes = append(bytes, '\r', '\n')
  bytes = append(bytes, v.bulk...)
  bytes = append(bytes, '\r', '\n')
  return bytes
}

func (v Value) marshalArray() []byte {
  len := len(v.array)
  var bytes []byte
  bytes = append(bytes, ARRAY)
  bytes = append(bytes, strconv.Itoa(len)...)
  bytes = append(bytes, '\r', '\n')

  for i := 0; i < len; i++ {
    bytes = append(bytes, v.array[i].Marshal()...)
  }
  
  return bytes
}

func (v Value) marshalError() []byte {
  var bytes []byte
  bytes = append(bytes, ERROR)
  bytes = append(bytes, v.str...)
  bytes = append(bytes, '\r', '\n')
  return bytes
}

func (v Value) marshalNull() []byte {
  return []byte("$-1\r\n")
}

