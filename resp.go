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
		// fmt.Println("Read byte:", b, "Current line:", string(line), "length:", len(line))
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

func (v Value) Marshal() []byte { // Marshal method that calls the specefied method based on the type of value
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshlNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte                  // create empty byte slice to hold marshaled output
	bytes = append(bytes, STRING)     // adds const prefix ('+') to start output
	bytes = append(bytes, v.str...)   // appends each byte of the string value to the output
	bytes = append(bytes, '\r', '\n') // appends suffix \r\n to indicate end of the string
	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte                                    // create empty byte slice to hold marshaled output
	bytes = append(bytes, BULK)                         // adds constant prefix ('$') to start output
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...) // appends the length of the bulk string
	bytes = append(bytes, '\r', '\n')                   // appends suffix \r\n to
	bytes = append(bytes, v.bulk...)                    // appends the bulk string itself
	bytes = append(bytes, '\r', '\n')                   // appends another \r\n to indicate end of the bulk string
	return bytes
}

func (v Value) marshalArray() []byte {
	len := len(v.array)                                 // get the length of the array
	var bytes []byte                                    // create empty byte slice to hold marshaled output
	bytes = append(bytes, ARRAY)                        // adds constant prefix ('*') to start output
	bytes = append(bytes, []byte(strconv.Itoa(len))...) // appends the length of the array as a string
	bytes = append(bytes, '\r', '\n')                   // appends suffix \r\n to indicate end of the length
	for i := 0; i < len; i++ {                          // iterate over each element in the array
		bytes = append(bytes, (v.array[i]).Marshal()...) // recursively marshal each element and append to the output
	}
	return bytes
}
