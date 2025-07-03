package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestMarshalString(t *testing.T) {
	v := Value{typ: "string", str: "OK"}
	expected := []byte("+OK\r\n")
	result := v.marshalString()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("marshalString failed: got %q, want %q", result, expected)
	}
}

func TestMarshalBulk(t *testing.T) {
	v := Value{typ: "bulk", bulk: "hello"}
	expected := []byte("$5\r\nhello\r\n")
	result := v.marshalBulk()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("marshalBulk failed: got %q, want %q", result, expected)
	}
}

func TestMarshalInteger(t *testing.T) {
	v := Value{typ: "integer", num: 42}
	expected := []byte(":42\r\n")
	result := v.marshalInteger()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("marshalInteger failed: got %q, want %q", result, expected)
	}
}

func TestMarshalNull(t *testing.T) {
	v := Value{typ: "null"}
	expected := []byte("$-1\r\n")
	result := v.marshalNull()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("marshalNull failed: got %q, want %q", result, expected)
	}
}

func TestMarshalError(t *testing.T) {
	v := Value{typ: "error", str: "ERR something went wrong"}
	expected := []byte("-ERR something went wrong\r\n")
	result := v.marshalError()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("marshalError failed: got %q, want %q", result, expected)
	}
}

func TestMarshalArray(t *testing.T) {
	v := Value{
		typ: "array",
		array: []Value{
			{typ: "string", str: "foo"},
			{typ: "bulk", bulk: "bar"},
		},
	}
	expected := []byte("*2\r\n+foo\r\n$3\r\nbar\r\n")
	result := v.marshalArray()
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("marshalArray failed: got %q, want %q", result, expected)
	}
}

func TestRESPReadString(t *testing.T) {
	input := "+OK\r\n"
	r := NewRESP(strings.NewReader(input))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if v.typ != "string" || v.str != "OK" {
		t.Errorf("Read string: got typ=%q, str=%q; want typ=\"string\", str=\"OK\"", v.typ, v.str)
	}
}

func TestRESPReadInteger(t *testing.T) {
	input := ":123\r\n"
	r := NewRESP(strings.NewReader(input))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if v.typ != "integer" || v.num != 123 {
		t.Errorf("Read integer: got typ=%q, num=%d; want typ=\"integer\", num=123", v.typ, v.num)
	}
}

func TestRESPReadBulk(t *testing.T) {
	input := "$5\r\nhello\r\n"
	r := NewRESP(strings.NewReader(input))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if v.typ != "bulk" || v.bulk != "hello" {
		t.Errorf("Read bulk: got typ=%q, bulk=%q; want typ=\"bulk\", bulk=\"hello\"", v.typ, v.bulk)
	}
}

func TestRESPReadArray(t *testing.T) {
	input := "*2\r\n+foo\r\n$3\r\nbar\r\n"
	r := NewRESP(strings.NewReader(input))
	v, err := r.Read()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if v.typ != "array" || len(v.array) != 2 {
		t.Errorf("Read array: got typ=%q, len=%d; want typ=\"array\", len=2", v.typ, len(v.array))
	}
	if v.array[0].typ != "string" || v.array[0].str != "foo" {
		t.Errorf("Read array[0]: got typ=%q, str=%q; want typ=\"string\", str=\"foo\"", v.array[0].typ, v.array[0].str)
	}
	if v.array[1].typ != "bulk" || v.array[1].bulk != "bar" {
		t.Errorf("Read array[1]: got typ=%q, bulk=%q; want typ=\"bulk\", bulk=\"bar\"", v.array[1].typ, v.array[1].bulk)
	}
}
