A **distributed file storage system** written in Go. It's in **early development** — the core building blocks (storage, P2P transport, server) are being implemented incrementally. Think of it as a nascent BitTorrent-like or IPFS-like system where nodes can store and retrieve files across a peer-to-peer network.
---
## Core Components
### 1. Storage Layer (`storage.go`)
**`Store`** — handles local file storage with Content-Addressable Storage (CAS):
| Function | Purpose |
|---|---|
| `CASPathTransformFunc(key)` | SHA-1 hashes a key, splits into 5-char path segments to avoid filesystem limits |
| `DefaultPathTransformFunc` | Identity transform (passthrough) |
| `Write(key, reader)` | Stores data on disk under the hashed path |
| `Read(key)` | Returns file contents as a `bytes.Buffer` |
| `Has(key)` | Checks if a key exists on disk |
| `Delete(key)` | Removes the file and its top-level folder |
| `Clear()` | Wipes the entire storage root |
**Key types:**
- `PathTransformFunc` — `func(string) PathKey` — transforms a string key into a filesystem path
- `PathKey` — `{Pathname, Filename}` with `FullPath()` helper
### 2. Server Layer (`server.go`)
**`FileServer`** — the main application orchestrator:
| Field | Purpose |
|---|---|
| `store *Store` | Local storage instance |
| `Transport p2p.Transport` | Network transport (TCP currently) |
| `quitch chan struct{}` | Shutdown signal channel |
**Flow:** `Start()` → `Transport.ListenAndAccept()` → enters `loop()` which consumes incoming `RPC` messages from the transport channel until `Quit()` is called.
> **Note:** The message loop currently just `fmt.Println(msg)` — real message handling is **not yet implemented**.
### 3. P2P Network Layer (`p2p/`)
**Interfaces** (`transport.go`):
- `Transport` — `ListenAndAccept()`, `Consume() <-chan RPC`, `Close()`
- `Peer` — `Close()`
**TCP Implementation** (`tcp_transport.go`):
- `TCPTransport` — listens on a port, accepts connections, decodes incoming messages into `RPC` structs and sends them on the `rpcch` channel
- `TCPPeer` — wraps a `net.Conn` with an `outbound` flag
- `TCPTransportOpts` — config: `ListenAddr`, `HandshakeFunc`, `Decoder`, `OnPeer` callback
**Message** (`message.go`):
- `RPC` — `{From net.Addr, Payload []byte}` — minimal message envelope
**Encoding** (`encoding.go`):
- `Decoder` interface — `Decode(io.Reader, *RPC) error`
- `GOBDecoder` — uses Go's `encoding/gob`
- `DefaultDecoder` — raw byte read (1028-byte buffer)
**Handshake** (`handshake.go`):
- `HandshakeFunc` — `func(Peer) error`
- `NOPHandshakeFunc` — no-op (always succeeds)
### 4. Entry Point (`main.go`)
Wires everything together:
1. Creates `TCPTransport` on `:3000` with `NOPHandshakeFunc` and `DefaultDecoder`
2. Creates `FileServer` with storage root `/tmp_network` and `CASPathTransformFunc`
3. Starts the server
4. Auto-quits after 3 seconds (temporary dev behavior)
---
## Current State & What's Missing
This is **Phase 1 / early scaffold**. Here's what's incomplete:
1. **No message handling** — `server.go:57` just prints messages; no actual file store/retrieve/delete over the network
2. **No outbound connections** — `TCPTransport` only accepts incoming connections; there's no `Dial()` or `Connect()` method to reach other peers
3. **No peer tracking** — `peers` map is declared but never populated
4. **`OnPeer` callback is commented out** in `main.go`
5. **`GOBDecoder` is unused** — `DefaultDecoder` (raw bytes) is used instead
6. **Auto-quit after 3s** — the server isn't meant to run indefinitely yet
7. **No request/response protocol** — `RPC` is just `{From, Payload}` with no message type, key, or status fields
8. **No replication/sharding logic** — no concept of which node owns which keys
---
## How to Work With It
```bash
make build   # builds to bin/fs
make run     # builds and runs (auto-quits after 3s)
make test    # runs all tests with verbose output
Tests cover: CASPathTransformFunc hashing, Store write/read/delete lifecycle, and basic TCPTransport listening.
