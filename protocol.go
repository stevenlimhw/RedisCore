package main

import (
	"bufio"
	"io"
	"log/slog"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Resp struct {
	reader *bufio.Reader
}

func NewResp(reader io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(reader),
	}
}

func (resp *Resp) Parse() (Value, error) {
	// the first byte determines the data type
	_type, err := resp.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case BULK:
		return resp.readBulk()
	case ARRAY:
		return resp.readArray()
	default:
		slog.Error("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// Read one byte at a time until we reach end of line indicated by '\r'.
// Then return the line without the last 2 bytes, which are '\r\n'.
func (resp *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := resp.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (resp *Resp) readInteger() (x int, n int, err error) {
	line, n, err := resp.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// Read the RESP Array type starting from the second byte, since
// the first byte has already been read in the Read method.
func (resp *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read array length
	arrLen, _, err := resp.readInteger()
	if err != nil {
		return v, err
	}

	// parse and read the value for each subsequent lines
	v.array = make([]Value, 0)
	for i := 0; i < arrLen; i++ {
		val, err := resp.Parse()
		if err != nil {
			return v, err
		}

		// append parsed value to array
		v.array = append(v.array, val)
	}

	return v, nil
}

func (resp *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	bulkLen, _, err := resp.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, bulkLen)
	resp.reader.Read(bulk)
	v.bulk = string(bulk)

	// consume the trailing CRLF
	resp.readLine()

	return v, nil
}
