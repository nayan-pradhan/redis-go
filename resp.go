package main

import (
	"bufio"
	"fmt"
	"io"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string  // Type of the value, eg STRING, ERROR, INTEGER, BULK, ARRAY
	str   string  // Holds value of string received from simple strings
	num   int     // Holds value of integer received from integers
	bulk  string  // Holds valueof strin received from bulk strings
	array []Value // Holds all value received from arrays
}

type RESP struct {
	reader *bufio.Reader // Buffered reader to read RESP data
}

func NewRESP(rd io.Reader) *RESP {
	return &RESP{reader: bufio.NewReader(rd)} // Create a new RESP instance with a buffered reader
}

func (r *RESP) readValue() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break // Check if the last two bytes are \r\n, indicating end of line
		}
		fmt.Println("Read byte:", b, "Current line:", string(line), "length:", len(line))
	}
	return line[:len(line)-2], n, nil
}
