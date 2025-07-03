package main

import (
	"reflect"
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
