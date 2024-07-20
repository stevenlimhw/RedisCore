package main

import (
	"bufio"
	"log/slog"
	"strconv"
)


type Resp struct {
  reader *bufio.Reader
}

func (resp *Resp) parseRespCommand() {
  // the first byte determines the data type
  b, _ := resp.reader.ReadByte()
  slog.Info(string(b))
  if b != '$' {
    slog.Error("Invalid type: expecting bulk strings only")
    //os.Exit(1)
    return
  }

  // the second byte represents the number of chars in the string
  sizeBytes, _ := resp.reader.ReadByte()
  size, _ := strconv.ParseInt(string(sizeBytes), 10, 64)
  slog.Info(strconv.Itoa(int(size)))

  // consume CR (i.e. \r)
  resp.reader.ReadByte()
  resp.reader.ReadByte()
  // consume LF (i.e. \n)
  resp.reader.ReadByte()
  resp.reader.ReadByte()

  // read the characters
  strContent := make([]byte, size)
  resp.reader.Read(strContent)
  slog.Info(string(strContent))
}


