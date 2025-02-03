# Gedis - A Lightweight Redis Clone

Gedis is a high-performance, in-memory key-value store designed for speed and simplicity. Inspired by Redis, Gedis supports basic data structures and commands while maintaining a lightweight footprint.

## Features
- In-memory key-value storage
- Support for strings, integers, arrays
- Simple and efficient command processing
- Fully compatible with REDIS clients
- Lightweight and fast

## Usage
Start the serer:
```bash
go run *.go
```

Connect with Redis Client:
```bash
redis-cli
```
Use basic commands
```bash
SET key value
GET key
DEL key
```
### License

MIT License
