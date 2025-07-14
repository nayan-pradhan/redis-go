# redis-go

A minimal Redis-like server implemented in Go, supporting basic key-value and hash operations with RESP protocol parsing and concurrent request handling.

## Features

- RESP protocol parsing (compatible with redis-cli)
- Basic commands: PING, SET, GET, DEL
- Hash commands: HSET, HGET, HGETALL
- Concurrency-safe with sync.RWMutex
- In-memory storage for strings and hashes
- Append Only File (AOF) persistence
- Extensible command handler map

## Getting Started

### Prerequisites

- Go 1.18+

### Build & Run

```bash
go run .
```

### Connect with redis-cli

```bash
redis-cli 
```

## Supported Commands

| Command   | Description                                      | Example Usage                |
|-----------|--------------------------------------------------|------------------------------|
| PING      | Test server connection                           | `PING`                       |
| SET       | Set key to value                                 | `SET mykey myvalue`          |
| GET       | Get value of key                                 | `GET mykey`                  |
| DEL       | Delete key from storage                          | `DEL mykey`                  |
| HSET      | Set field in hash                                | `HSET myhash field value`    |
| HGET      | Get field value from hash                        | `HGET myhash field`          |
| HGETALL   | Get all fields and values from hash              | `HGETALL myhash`             |

## Project Structure

- `main.go`      - Server loop, connection handling, command dispatch
- `handler.go`   - Command handlers (PING, SET, GET, DEL, HSET, HGET, HGETALL)
- `resp.go`      - RESP protocol parsing and marshalling
- `aof.go`       - Append Only File persistence logic
- `resp_test.go` - Unit tests for RESP parsing and marshalling

## Development

### Add a New Command

1. Implement a handler function in `handler.go`.
2. Add the handler to the `Handlers` map.
3. Update tests in `resp_test.go` if needed.

### Run Tests

```bash
go test
```
