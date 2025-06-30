package resp

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
