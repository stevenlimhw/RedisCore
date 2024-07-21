package main

import "strconv"

type Value struct {
	typ   string
	str   string
	bulk  string
	array []*Value
	num   int
}

func (v *Value) Marshal() []byte {
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

func (v *Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, stringToken)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, bulkToken)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalArray() []byte {
	arrLen := len(v.array)
	var bytes []byte
	bytes = append(bytes, arrayToken)
	bytes = append(bytes, strconv.Itoa(arrLen)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < arrLen; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

func (v *Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, errorToken)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v *Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}
