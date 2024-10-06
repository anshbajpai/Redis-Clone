package main

import (
	"bufio"
	"fmt"
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

func (resp *Resp) Read() (Value, error) {
	dataType, err := resp.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch dataType {
	case ARRAY:
		return resp.readArray()
	case BULK:
		return resp.readBulk()
	default:
		fmt.Printf("Unrecognized data type: %v", string(dataType))
		return Value{}, nil
	}
}

func (resp *Resp) readArray() (Value, error) {
	result := Value{}
	result.typ = "array"

	// Read the array length
	length, _, err := resp.readInteger()
	if err != nil {
		return result, err
	}

	// Initialize the array to hold the parsed values
	result.array = make([]Value, 0)
	for i := 0; i < length; i++ {
		parsedValue, err := resp.Read()
		if err != nil {
			return result, err
		}

		// Append each parsed value to the result array
		result.array = append(result.array, parsedValue)
	}

	return result, nil
}

func (resp *Resp) readBulk() (Value, error) {
	result := Value{}
	result.typ = "bulk"

	// Read the length of the bulk string
	length, _, err := resp.readInteger()
	if err != nil {
		return result, err
	}

	// Allocate a buffer for the bulk data
	bulkData := make([]byte, length)

	// Read the bulk data into the buffer
	resp.reader.Read(bulkData)

	// Convert the byte slice into a string
	result.bulk = string(bulkData)

	// Read and discard the trailing CRLF
	resp.readLine()

	return result, nil
}

func (v Value) Convert() []byte {
	var convertFunc func() []byte

	switch v.typ {
	case "array":
		convertFunc = v.convertArray
	case "bulk":
		convertFunc = v.convertBulk
	case "string":
		convertFunc = v.convertString
	case "null":
		convertFunc = v.convertNull
	case "error":
		convertFunc = v.convertError
	default:
		return []byte{}
	}

	return convertFunc()
}

func (v Value) convertString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) convertBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) convertArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Convert()...)
	}

	return bytes
}

func (v Value) convertError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) convertNull() []byte {
	return []byte("$-1\r\n")
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Convert()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
