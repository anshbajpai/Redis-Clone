package main

import (
	"bufio"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (resp *Resp) readLine() (content []byte, bytesRead int, err error) {
	for {
		byteRead, err := resp.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		bytesRead++
		content = append(content, byteRead)
		if len(content) >= 2 && content[len(content)-2] == '\r' {
			break
		}
	}
	return content[:len(content)-2], bytesRead, nil
}

func (resp *Resp) readInteger() (number int, bytesRead int, err error) {
	lineContent, bytesRead, err := resp.readLine()
	if err != nil {
		return 0, 0, err
	}
	parsedInt64, err := strconv.ParseInt(string(lineContent), 10, 64)
	if err != nil {
		return 0, bytesRead, err
	}
	return int(parsedInt64), bytesRead, nil
}
