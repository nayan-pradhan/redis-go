package main

import (
	"bufio"
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
