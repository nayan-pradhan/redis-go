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

func (r *RESP) readInteger() (x int, n int, err error) {
	line, n, err := r.readValue() // get the line from the reader
	if err != nil {
		return 0, 0, err // return error if any
	}
	i64, err := strconv.ParseInt(string(line), 10, 64) // parse the line as an integer
	if err != nil {
		return 0, 0, err // return error if parsing fails
	}
	return int(i64), n, nil // return the parsed integer and number of bytes read
}

func (r *RESP) Read() (Value, error) {
	_type, err := r.reader.ReadByte() // first byte indicates the type of value
	if err != nil {
		return Value{}, err
	}
	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Println("Unknown type: ", string(_type))
		return Value{}, nil
	}
}

func (r *RESP) readArray() (Value, error) {
	v := Value{typ: "array"}

	len, _, err := r.readInteger() // read the length of the array
	if err != nil {
		return v, err // return error if reading length fails
	}
	v.array = make([]Value, 0)
	for range len {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val) // append the read value to the array
	}
	return v, nil
}

func (r *RESP) readBulk() (Value, error) {
	v := Value{typ: "bulk"}

	len, _, err := r.readInteger() // read the length of the bulk string
	if err != nil {
		return v, err // return error if reading length fails
	}
	bulk := make([]byte, len) // create a byte slice of the specified length
	r.reader.Read(bulk)       // read the bulk string into the byte slice
	v.bulk = string(bulk)     // convert the byte slice to a string
	r.readValue()             // read the trailing \r\n
	return v, nil
}
