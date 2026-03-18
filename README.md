# Distributed Cache (Go)

A lightweight HTTP-based distributed cache service written in Go with:

- LRU eviction
- Optional TTL per key
- Save/load persistence to disk
- Basic replica forwarding support
- Consistent hash client library (`client/`)

## Run Locally

### Prerequisites
- Go 1.26+

### Start server
```bash
go run main.go -port 8000 -capacity 100
```

### Quick API test
```bash
# set key
curl "http://localhost:8000/set?key=user&value=john"

# get key
curl "http://localhost:8000/get?key=user"

# set with ttl (seconds)
curl "http://localhost:8000/set?key=temp&value=123&ttl=60"

# save cache to disk
curl "http://localhost:8000/save?file=cache.json"

# load cache from disk
curl "http://localhost:8000/load?file=cache.json"
```

## Deploy Using Docker

### Build image
```bash
docker build -t distributed-cache:latest .
```

### Run container
```bash
docker run --rm -p 8000:8000 distributed-cache:latest
```

Server will be available at:
- `http://localhost:8000/`
- Container runs as non-root user and includes a healthcheck on `/`.

## Deploy Using Docker Compose

```bash
docker compose up --build
```

This starts one service:
- `distributed-cache` on port `8000`
- with `restart: unless-stopped` and a built-in healthcheck.

## Project Structure

```text
.
|-- main.go
|-- Dockerfile
|-- docker-compose.yml
|-- cache/
|   |-- cache.go
|   |-- lru.go
|   `-- persistence.go
|-- server/
|   `-- server.go
|-- hash/
|   `-- consistent_hash.go
`-- client/
    `-- client.go
```

## Notes

- Keys set without `ttl` do not expire.
- `GET /get` returns `404` when key is not found or expired.
- Persistence files are created inside the running container/app working directory.
