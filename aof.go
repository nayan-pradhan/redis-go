package main

import (
	"bufio"
	"os"
	"sync"
)

// AOF: Append Only File - we record each command in a file as RESP.
// When a restart occurs, Redis reads all RESP commands from teh AOF and executes them in memory for data persistency.

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}
